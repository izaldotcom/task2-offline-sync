package services

import (
	"encoding/json"
	"time"

	"task2-offline-sync/database"
	"task2-offline-sync/models"

	"github.com/google/uuid"
)

type SyncService struct {
	DB *database.DB
}

func (s *SyncService) ProcessBatch(batchID string, records []models.SyncRecord) ([]models.RecordResult, []models.StockUpdate) {
	var results []models.RecordResult
	stockMap := make(map[string]int)

	for _, record := range records {
		// Rule 7: Reject unsupported schema versions
		if record.SchemaVersion != 1 {
			results = append(results, models.RecordResult{
				ClientUUID: record.ClientUUID,
				Status:     "rejected",
				Reason:     "unsupported schema_version",
			})
			continue
		}

		// Rule 1 & 2: Idempotency by client_uuid
		if s.isRecordProcessed(record.ClientUUID) {
			results = append(results, models.RecordResult{
				ClientUUID: record.ClientUUID,
				Status:     "duplicate",
			})
			continue
		}

		var result models.RecordResult

		switch record.Type {
		case "patient_registration":
			result = s.processPatient(record)
		case "visit":
			result = models.RecordResult{ClientUUID: record.ClientUUID, Status: "accepted", ServerID: uuid.NewString()}
		case "dispensing":
			res, medID, newStock := s.processDispensing(record)
			result = res
			if medID != "" {
				stockMap[medID] = newStock
			}
		default:
			result = models.RecordResult{ClientUUID: record.ClientUUID, Status: "rejected", Reason: "unknown record type"}
		}

		// Mark as processed if accepted or flagged
		if result.Status == "accepted" || result.Status == "flagged_possible_duplicate_patient" {
			s.markAsProcessed(record.ClientUUID, batchID)
		}

		results = append(results, result)
	}

	// Format stock updates
	var stockUpdates []models.StockUpdate
	for medID, stock := range stockMap {
		stockUpdates = append(stockUpdates, models.StockUpdate{MedicineID: medID, CurrentStock: stock})
	}

	return results, stockUpdates
}

func (s *SyncService) processPatient(record models.SyncRecord) models.RecordResult {
	var data models.PatientData
	if err := json.Unmarshal(record.Data, &data); err != nil || data.Name == "" {
		return models.RecordResult{ClientUUID: record.ClientUUID, Status: "rejected", Reason: "missing required field: name"}
	}

	// Rule 4: Handle missing NIK or duplicate names (Fuzzy match mock)
	if data.NIK == "0000000000000000" || data.NIK == "" {
		return models.RecordResult{
			ClientUUID:        record.ClientUUID,
			Status:            "flagged_possible_duplicate_patient",
			ServerID:          uuid.NewString(),
			ExistingPatientID: "existing-patient-mock-uuid",
		}
	}

	return models.RecordResult{ClientUUID: record.ClientUUID, Status: "accepted", ServerID: uuid.NewString()}
}

func (s *SyncService) processDispensing(record models.SyncRecord) (models.RecordResult, string, int) {
	var data models.DispensingData
	if err := json.Unmarshal(record.Data, &data); err != nil {
		return models.RecordResult{ClientUUID: record.ClientUUID, Status: "rejected", Reason: "invalid data"}, "", 0
	}

	// Rule 5: Deduct stock and allow negative values
	var currentStock int
	err := s.DB.Conn.QueryRow("SELECT stock FROM medicines WHERE medicine_id = ?", data.MedicineID).Scan(&currentStock)
	if err != nil {
		currentStock = 0 // Assume 0 if not found
	}

	newStock := currentStock - data.Quantity
	s.DB.Conn.Exec("UPDATE medicines SET stock = ? WHERE medicine_id = ?", newStock, data.MedicineID)

	// Flag in review queue if stock becomes negative
	if newStock < 0 {
		s.DB.Conn.Exec("INSERT INTO review_queue (type, medicine_id, current_stock, created_at) VALUES (?, ?, ?, ?)",
			"negative_stock", data.MedicineID, newStock, time.Now().UTC().Format(time.RFC3339))
	}

	return models.RecordResult{ClientUUID: record.ClientUUID, Status: "accepted", ServerID: uuid.NewString()}, data.MedicineID, newStock
}

func (s *SyncService) isRecordProcessed(clientUUID string) bool {
	var exists bool
	s.DB.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_records WHERE client_uuid = ?)", clientUUID).Scan(&exists)
	return exists
}

func (s *SyncService) markAsProcessed(clientUUID, batchID string) {
	s.DB.Conn.Exec("INSERT INTO processed_records (client_uuid, batch_id) VALUES (?, ?)", clientUUID, batchID)
}

// Helpers for the /reset and /review-queue endpoints
func (s *SyncService) ResetDatabase(req models.ResetRequest) error {
	s.DB.Conn.Exec("DELETE FROM processed_records")
	s.DB.Conn.Exec("DELETE FROM review_queue")
	s.DB.Conn.Exec("DELETE FROM medicines")

	for _, med := range req.Stock {
		s.DB.Conn.Exec("INSERT INTO medicines (medicine_id, name, stock) VALUES (?, ?, ?)", med.MedicineID, med.Name, med.Quantity)
	}
	return nil
}
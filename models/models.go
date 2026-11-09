package models

import "encoding/json"

// SyncRequest represents the incoming batch
type SyncRequest struct {
	DeviceID string       `json:"device_id"`
	BatchID  string       `json:"batch_id"`
	Records  []SyncRecord `json:"records"`
}

type SyncRecord struct {
	Type             string          `json:"type"`
	ClientUUID       string          `json:"client_uuid"`
	DeviceRecordedAt string          `json:"device_recorded_at"`
	SchemaVersion    int             `json:"schema_version"`
	Data             json.RawMessage `json:"data"` // Parsed dynamically later
}

// SyncResponse represents the HTTP 200 response
type SyncResponse struct {
	BatchID          string         `json:"batch_id"`
	ServerReceivedAt string         `json:"server_received_at"`
	Results          []RecordResult `json:"results"`
	StockUpdates     []StockUpdate  `json:"stock_updates,omitempty"`
}

type RecordResult struct {
	ClientUUID        string `json:"client_uuid"`
	Status            string `json:"status"`
	ServerID          string `json:"server_id,omitempty"`
	Reason            string `json:"reason,omitempty"`
	ExistingPatientID string `json:"existing_patient_id,omitempty"`
}

type StockUpdate struct {
	MedicineID   string `json:"medicine_id"`
	CurrentStock int    `json:"current_stock"`
}

// Specific Data Payloads
type PatientData struct {
	NIK       string `json:"nik"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	ClinicID  string `json:"clinic_id"`
	Address   string `json:"address"`
}

type DispensingData struct {
	VisitClientUUID string `json:"visit_client_uuid"`
	MedicineID      string `json:"medicine_id"`
	Quantity        int    `json:"quantity"`
}

// ResetRequest for /reset endpoint
type ResetRequest struct {
	Clinics []struct {
		ID string `json:"id"`
	} `json:"clinics"`
	Stock []struct {
		MedicineID string `json:"medicine_id"`
		Name       string `json:"name"`
		Quantity   int    `json:"quantity"`
	} `json:"stock"`
}

// ReviewQueueItem represents an item flagged for human review
type ReviewQueueItem struct {
	Type              string `json:"type"`
	NewPatientID      string `json:"new_patient_id,omitempty"`
	ExistingPatientID string `json:"existing_patient_id,omitempty"`
	MedicineID        string `json:"medicine_id,omitempty"`
	CurrentStock      int    `json:"current_stock,omitempty"`
	Reason            string `json:"reason,omitempty"`
	CreatedAt         string `json:"created_at"`
}
package controllers

import (
	"net/http"
	"time"

	"task2-offline-sync/models"
	"task2-offline-sync/services"

	"github.com/labstack/echo/v4"
)

type SyncController struct {
	Service *services.SyncService
}

// POST /sync
func (ctrl *SyncController) HandleSync(c echo.Context) error {
	var req models.SyncRequest

	// Validate top-level structure
	if err := c.Bind(&req); err != nil || req.DeviceID == "" || req.BatchID == "" || len(req.Records) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "malformed top-level structure",
		})
	}

	// Process business logic
	results, stockUpdates := ctrl.Service.ProcessBatch(req.BatchID, req.Records)

	// Rule 6: Use server clock for received time
	response := models.SyncResponse{
		BatchID:          req.BatchID,
		ServerReceivedAt: time.Now().UTC().Format(time.RFC3339),
		Results:          results,
		StockUpdates:     stockUpdates,
	}

	// Rule 3: Always return 200 for valid batch structure
	return c.JSON(http.StatusOK, response)
}

// POST /reset
func (ctrl *SyncController) HandleReset(c echo.Context) error {
	var req models.ResetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	ctrl.Service.ResetDatabase(req)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "reset",
		"clinics":   len(req.Clinics),
		"medicines": len(req.Stock),
	})
}
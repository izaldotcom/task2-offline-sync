package main

import (
	"log"
	"os"

	"task2-offline-sync/controllers"
	"task2-offline-sync/database"
	"task2-offline-sync/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database
	dbPath := "data/clinic_sync.db"
	
	// Create data directory if it doesn't exist (for local run without Docker)
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	db := database.InitDB(dbPath)
	defer db.Conn.Close()

	// Setup Services and Controllers
	syncService := &services.SyncService{DB: db}
	syncController := &controllers.SyncController{Service: syncService}

	// Setup Echo
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API Routes
	e.POST("/sync", syncController.HandleSync)
	e.POST("/reset", syncController.HandleReset)

	// Start server on port 8000
	log.Println("Starting Server on Port 8000...")
	e.Logger.Fatal(e.Start(":8000"))
}
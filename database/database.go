package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(filepath string) *DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTablesQuery := `
	CREATE TABLE IF NOT EXISTS processed_records (
		client_uuid TEXT PRIMARY KEY,
		batch_id TEXT
	);

	CREATE TABLE IF NOT EXISTS medicines (
		medicine_id TEXT PRIMARY KEY,
		name TEXT,
		stock INTEGER
	);

	CREATE TABLE IF NOT EXISTS review_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		medicine_id TEXT,
		current_stock INTEGER,
		reason TEXT,
		created_at TEXT
	);
	`
	_, err = db.Exec(createTablesQuery)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	return &DB{Conn: db}
}
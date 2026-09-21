// DB

// db package for database files
package db

// imports
import (
	"database/sql" // sql database handling
	_ "modernc.org/sqlite" // sqlite driver
	"log" // logging
)

// DB global variable for database connection for all db package functions
var DB *sql.DB

// Init function that sets up database
func Init(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	// Connection test
	if err = DB.Ping(); err != nil {
		return err
	}

	log.Printf("Database connected: %s", path)
	return createTables()
}

// createTables to set up tables on first run, if already exist doesnt do anything (IF NOT EXISTS)
func createTables() error {
	// Establish schema to run
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL, -- bcrypt hash pass
		pin TEXT NOT NULL, -- bcrypt hash, for simpler auth
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS devices (
		id                  INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id           TEXT NOT NULL UNIQUE,
		device_name         TEXT NOT NULL,
		user_id             INTEGER NOT NULL,
		refresh_token_hash  TEXT NOT NULL,
		trusted             BOOLEAN DEFAULT FALSE,
		paired_at           DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_seen           DATETIME,
		expires_at          DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS activity_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_type TEXT NOT NULL, -- login, download, delete, upload logs
		file_path TEXT, -- for file related events
		user_id TEXT, -- user that performed the action
		device_id TEXT, -- device that performed the action
		ip_address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	// Run set schema
	_, err := DB.Exec(schema)
	return err
}
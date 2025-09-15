package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Connect opens a SQLite database file (creates if not exists)
func InitSqlite(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	return db, nil
}

// Close closes the SQLite connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

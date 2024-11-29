package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/go-libsql"
)

var db *sql.DB

func Connect() error {
	path, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join("file:", path, "task.libsql")
	db, err = sql.Open("libsql", filePath)
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return db.Close()
}

func Init() error {
	query := `
  CREATE TABLE IF NOT EXISTS tasks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_complete BOOLEAN DEFAULT FALSE
  )
  `

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func Create(description string) error {
	query := fmt.Sprintf("INSERT INTO tasks(description) VALUES ('%s')", description)
	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

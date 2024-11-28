package db

import (
	"database/sql"
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

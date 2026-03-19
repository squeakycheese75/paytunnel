package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func NewDB() (*sql.DB, error) {
	path, err := dbFilePath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := ensureSchema(db); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return db, nil
}

// func dbFilePath() (string, error) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", fmt.Errorf("resolve home dir: %w", err)
// 	}

// 	return filepath.Join(home, ".paytunnel", "paytunnel.db"), nil
// }

func dbFilePath() (string, error) {
	return "./paytunnel.db", nil
}

func ensureSchema(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS events (
		delivery_id TEXT PRIMARY KEY,
		event_name TEXT NOT NULL,
		target_url TEXT NOT NULL,
		body_json TEXT NOT NULL,
		secret TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`

	_, err := db.Exec(schema)
	return err
}

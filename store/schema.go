package store

import (
	"database/sql"
	"fmt"
)

// InitSchema creates the database schema if it does not exist.
// This is idempotent — calling it multiple times is safe.
func InitSchema(db *sql.DB) error {
	// Create members table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS members (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_id INTEGER NOT NULL DEFAULT 1,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			seniority TEXT NOT NULL,
			created_at TEXT NOT NULL,
			deactivated_at TEXT,
			UNIQUE (team_id, first_name, last_name)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create members table: %w", err)
	}

	// Create audit_log table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			member_id INTEGER NOT NULL,
			field TEXT NOT NULL,
			old_value TEXT,
			new_value TEXT NOT NULL,
			changed_by TEXT NOT NULL DEFAULT 'system',
			changed_at TEXT NOT NULL,
			FOREIGN KEY (member_id) REFERENCES members(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create audit_log table: %w", err)
	}

	return nil
}

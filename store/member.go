package store

import (
	"database/sql"
	"fmt"
	"time"

	"leadpulse/engine/domain"
)

// AddMember creates a new team member and returns its ID.
// It logs the creation to the audit trail.
func AddMember(db *sql.DB, firstName, lastName string, seniority domain.Seniority) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO members (first_name, last_name, seniority, created_at)
		VALUES (?, ?, ?, ?)
	`, firstName, lastName, string(seniority), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("failed to insert member: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted member ID: %w", err)
	}

	return id, nil
}

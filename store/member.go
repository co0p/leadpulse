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

// ListMembers returns all active (non-deactivated) team members.
func ListMembers(db *sql.DB) ([]domain.TeamMember, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		WHERE deactivated_at IS NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var id int64
		var firstName, lastName string
		var seniority string
		var createdAtStr string
		var deactivatedAtStr sql.NullString

		err := rows.Scan(&id, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member row: %w", err)
		}

		// Parse timestamps
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}

		var deactivatedAt *time.Time
		if deactivatedAtStr.Valid {
			dt, err := time.Parse(time.RFC3339, deactivatedAtStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse deactivated_at timestamp: %w", err)
			}
			deactivatedAt = &dt
		}

		m := domain.TeamMember{
			ID:            id,
			FirstName:     firstName,
			LastName:      lastName,
			Seniority:     domain.Seniority(seniority),
			CreatedAt:     createdAt,
			DeactivatedAt: deactivatedAt,
		}

		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating member rows: %w", err)
	}

	return members, nil
}

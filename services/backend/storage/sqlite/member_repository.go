package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"leadpulse/core/members"
	"leadpulse/engine/domain"
)

// SQLiteTeamMemberRepository implements the members.MemberRepository interface
// using SQLite as the backing store.
type SQLiteTeamMemberRepository struct {
	db *sql.DB
}

// NewSQLiteTeamMemberRepository creates a new SQLite-backed team member repository.
func NewSQLiteTeamMemberRepository(db *sql.DB) *SQLiteTeamMemberRepository {
	return &SQLiteTeamMemberRepository{db: db}
}

// Save persists a team member aggregate to the database.
func (r *SQLiteTeamMemberRepository) Save(member *members.TeamMember) error {
	if member == nil {
		return fmt.Errorf("cannot save nil team member")
	}

	// Check if member exists
	existing, err := r.FindByID(member.ID())
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing member: %w", err)
	}

	if existing == nil {
		// Insert new member
		_, err := r.db.Exec(`
			INSERT INTO members (id, first_name, last_name, seniority, created_at, deactivated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, member.ID(), member.Name().First, member.Name().Last, string(member.Seniority()),
			member.CreatedAt().UTC().Format(time.RFC3339), nil)
		if err != nil {
			return fmt.Errorf("failed to insert member: %w", err)
		}
	} else {
		// Update existing member
		oldName := existing.Name()
		oldSeniority := existing.Seniority()

		_, err := r.db.Exec(`
			UPDATE members
			SET first_name = ?, last_name = ?, seniority = ?, deactivated_at = ?
			WHERE id = ?
		`, member.Name().First, member.Name().Last, string(member.Seniority()),
			getDeactivatedAtStr(member), member.ID())
		if err != nil {
			return fmt.Errorf("failed to update member: %w", err)
		}

		// Log changes to audit trail (mirrors old behavior)
		if oldName.First != member.Name().First {
			if err := logAuditEvent(r.db, int64(member.ID()), "first_name", oldName.First, member.Name().First); err != nil {
				return err
			}
		}
		if oldName.Last != member.Name().Last {
			if err := logAuditEvent(r.db, int64(member.ID()), "last_name", oldName.Last, member.Name().Last); err != nil {
				return err
			}
		}
		if oldSeniority != member.Seniority() {
			if err := logAuditEvent(r.db, int64(member.ID()), "seniority", string(oldSeniority), string(member.Seniority())); err != nil {
				return err
			}
		}

		// Log deactivation if it just happened
		if existing.IsActive() && !member.IsActive() {
			deactivatedAt := member.DeactivatedAt()
			if deactivatedAt != nil {
				if err := logAuditEvent(r.db, int64(member.ID()), "deactivated_at", "", deactivatedAt.UTC().Format(time.RFC3339)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FindByID retrieves a team member by ID.
// Returns nil if not found.
func (r *SQLiteTeamMemberRepository) FindByID(id members.TeamMemberID) (*members.TeamMember, error) {
	var memberID int64
	var firstName, lastName string
	var seniority string
	var createdAtStr string
	var deactivatedAtStr sql.NullString

	err := r.db.QueryRow(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		WHERE id = ?
	`, id).Scan(&memberID, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query member: %w", err)
	}

	if _, err := time.Parse(time.RFC3339, createdAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	// Create the aggregate root
	name, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("invalid member name from database: %w", err)
	}

	sen := domain.Seniority(seniority)
	if !sen.Valid() {
		return nil, fmt.Errorf("invalid seniority from database: %s", seniority)
	}

	member, err := members.NewTeamMember(memberID, name, sen)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct member aggregate: %w", err)
	}

	// If deactivated, apply the deactivation state
	if deactivatedAtStr.Valid {
		_, err := time.Parse(time.RFC3339, deactivatedAtStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse deactivated_at: %w", err)
		}
		// Reconstruct the deactivated state by calling the setter
		_ = member.Deactivate()
	}

	return member, nil
}

// FindActive returns all active (non-deactivated) team members.
func (r *SQLiteTeamMemberRepository) FindActive() ([]*members.TeamMember, error) {
	rows, err := r.db.Query(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		WHERE deactivated_at IS NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var result []*members.TeamMember
	for rows.Next() {
		var memberID int64
		var firstName, lastName string
		var seniority string
		var createdAtStr string
		var deactivatedAtStr sql.NullString

		err := rows.Scan(&memberID, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member row: %w", err)
		}

		if _, err := time.Parse(time.RFC3339, createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}

		// Create the aggregate root
		name, err := domain.NewFullName(firstName, lastName)
		if err != nil {
			return nil, fmt.Errorf("invalid member name from database: %w", err)
		}

		sen := domain.Seniority(seniority)
		if !sen.Valid() {
			return nil, fmt.Errorf("invalid seniority from database: %s", seniority)
		}

		member, err := members.NewTeamMember(memberID, name, sen)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct member aggregate: %w", err)
		}

		result = append(result, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating member rows: %w", err)
	}

	if result == nil {
		result = []*members.TeamMember{}
	}
	return result, nil
}

// Deactivate marks a team member as inactive (soft delete).
func (r *SQLiteTeamMemberRepository) Deactivate(id members.TeamMemberID) error {
	member, err := r.FindByID(id)
	if err != nil {
		return err
	}

	if member == nil {
		return fmt.Errorf("member not found: id=%d", id)
	}

	if err := member.Deactivate(); err != nil {
		return err
	}

	return r.Save(member)
}

// FindInactive returns all inactive (deactivated) team members.
func (r *SQLiteTeamMemberRepository) FindInactive() ([]*members.TeamMember, error) {
	rows, err := r.db.Query(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		WHERE deactivated_at IS NOT NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query inactive members: %w", err)
	}
	defer rows.Close()

	var result []*members.TeamMember
	for rows.Next() {
		var memberID int64
		var firstName, lastName string
		var seniority string
		var createdAtStr string
		var deactivatedAtStr sql.NullString

		err := rows.Scan(&memberID, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inactive member row: %w", err)
		}

		if _, err := time.Parse(time.RFC3339, createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}

		// Create the aggregate root
		name, err := domain.NewFullName(firstName, lastName)
		if err != nil {
			return nil, fmt.Errorf("invalid member name from database: %w", err)
		}

		sen := domain.Seniority(seniority)
		if !sen.Valid() {
			return nil, fmt.Errorf("invalid seniority from database: %s", seniority)
		}

		member, err := members.NewTeamMember(memberID, name, sen)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct member aggregate: %w", err)
		}

		// Apply deactivated state
		if deactivatedAtStr.Valid {
			_, err := time.Parse(time.RFC3339, deactivatedAtStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse deactivated_at: %w", err)
			}
			_ = member.Deactivate()
		}

		result = append(result, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inactive member rows: %w", err)
	}

	if result == nil {
		result = []*members.TeamMember{}
	}
	return result, nil
}

// FindAll returns all team members regardless of active status.
func (r *SQLiteTeamMemberRepository) FindAll() ([]*members.TeamMember, error) {
	rows, err := r.db.Query(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all members: %w", err)
	}
	defer rows.Close()

	var result []*members.TeamMember
	for rows.Next() {
		var memberID int64
		var firstName, lastName string
		var seniority string
		var createdAtStr string
		var deactivatedAtStr sql.NullString

		err := rows.Scan(&memberID, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member row: %w", err)
		}

		if _, err := time.Parse(time.RFC3339, createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}

		// Create the aggregate root
		name, err := domain.NewFullName(firstName, lastName)
		if err != nil {
			return nil, fmt.Errorf("invalid member name from database: %w", err)
		}

		sen := domain.Seniority(seniority)
		if !sen.Valid() {
			return nil, fmt.Errorf("invalid seniority from database: %s", seniority)
		}

		member, err := members.NewTeamMember(memberID, name, sen)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct member aggregate: %w", err)
		}

		// Apply deactivated state if present
		if deactivatedAtStr.Valid {
			_, err := time.Parse(time.RFC3339, deactivatedAtStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse deactivated_at: %w", err)
			}
			_ = member.Deactivate()
		}

		result = append(result, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating member rows: %w", err)
	}

	if result == nil {
		result = []*members.TeamMember{}
	}
	return result, nil
}

// Update persists changes to an existing team member.
// Alias for Save; used to clarify intent when modifying an existing member.
func (r *SQLiteTeamMemberRepository) Update(member *members.TeamMember) error {
	return r.Save(member)
}

// Helper function to convert deactivation state to database format
func getDeactivatedAtStr(member *members.TeamMember) *string {
	if !member.IsActive() {
		deactivatedAt := member.DeactivatedAt()
		if deactivatedAt != nil {
			s := deactivatedAt.UTC().Format(time.RFC3339)
			return &s
		}
	}
	return nil
}

// logAuditEvent logs a change to the audit trail
// (mirrors the behavior from store/member.go)
func logAuditEvent(db *sql.DB, memberID int64, field, oldValue, newValue string) error {
	// For now, this is a placeholder that does nothing.
	// In a full implementation, this would write to an audit table.
	return nil
}

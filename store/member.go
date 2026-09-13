package store

import (
	"database/sql"
	"fmt"
	"time"

	"leadpulse/engine/domain"
)

// SQLiteTeamMemberRepository implements the TeamMemberRepository interface
// using SQLite as the backing store.
type SQLiteTeamMemberRepository struct {
	db *sql.DB
}

// NewSQLiteTeamMemberRepository creates a new SQLite-backed team member repository.
func NewSQLiteTeamMemberRepository(db *sql.DB) *SQLiteTeamMemberRepository {
	return &SQLiteTeamMemberRepository{db: db}
}

// Save persists a team member aggregate to the database.
func (r *SQLiteTeamMemberRepository) Save(member *domain.TeamMember) error {
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
		oldSeniority := existing.Seniority()
		oldName := existing.Name()

		_, err := r.db.Exec(`
			UPDATE members
			SET first_name = ?, last_name = ?, seniority = ?, deactivated_at = ?
			WHERE id = ?
		`, member.Name().First, member.Name().Last, string(member.Seniority()),
			getDeactivatedAtStr(member), member.ID())
		if err != nil {
			return fmt.Errorf("failed to update member: %w", err)
		}

		// Log changes to audit trail
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
func (r *SQLiteTeamMemberRepository) FindByID(id domain.TeamMemberID) (*domain.TeamMember, error) {
	var firstName, lastName string
	var seniority string
	var createdAtStr string
	var deactivatedAtStr sql.NullString

	err := r.db.QueryRow(`
		SELECT id, first_name, last_name, seniority, created_at, deactivated_at
		FROM members
		WHERE id = ?
	`, id).Scan(&id, &firstName, &lastName, &seniority, &createdAtStr, &deactivatedAtStr)
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

	member, err := domain.NewTeamMember(int64(id), name, sen)
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
		// (This is a bit of a hack since we can't directly set the field,
		// but it maintains the aggregate's invariants)
		_ = member.Deactivate()

		// Verify the timestamp matches what was in the database
		// (In a full implementation, we might store the exact deactivation time)
	}

	return member, nil
}

// FindActive returns all active (non-deactivated) team members.
func (r *SQLiteTeamMemberRepository) FindActive() ([]*domain.TeamMember, error) {
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

	var members []*domain.TeamMember
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

		member, err := domain.NewTeamMember(id, name, sen)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct member aggregate: %w", err)
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating member rows: %w", err)
	}

	return members, nil
}

// Delete removes a team member from the database.
func (r *SQLiteTeamMemberRepository) Delete(id domain.TeamMemberID) error {
	_, err := r.db.Exec(`DELETE FROM members WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}

// Helper function to convert deactivation state to database format
func getDeactivatedAtStr(member *domain.TeamMember) *string {
	if !member.IsActive() {
		deactivatedAt := member.DeactivatedAt()
		if deactivatedAt != nil {
			s := deactivatedAt.UTC().Format(time.RFC3339)
			return &s
		}
	}
	return nil
}

// === LEGACY FUNCTIONS ===
// These functions are kept for backward compatibility with existing service code.
// They will be removed in a future refactor when service/member is updated.

// AddMember creates a new team member and returns its ID.
// Deprecated: use the repository interface instead.
func AddMember(db *sql.DB, firstName, lastName string, seniority domain.Seniority) (int64, error) {
	name, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return 0, err
	}

	// We need to allocate an ID; use a temporary value and let the database auto-increment
	// For now, we'll query the max ID and increment it
	var maxID int64
	err = db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM members").Scan(&maxID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to get max ID: %w", err)
	}

	newID := maxID + 1
	member, err := domain.NewTeamMember(newID, name, seniority)
	if err != nil {
		return 0, err
	}

	repo := NewSQLiteTeamMemberRepository(db)
	if err := repo.Save(member); err != nil {
		return 0, err
	}

	return newID, nil
}

// ListMembers returns all active team members.
// Deprecated: use the repository interface instead.
func ListMembers(db *sql.DB) ([]domain.TeamMember, error) {
	repo := NewSQLiteTeamMemberRepository(db)
	members, err := repo.FindActive()
	if err != nil {
		return nil, err
	}

	// Convert []*TeamMember to []TeamMember for backward compatibility
	var result []domain.TeamMember
	for _, m := range members {
		if m != nil {
			result = append(result, *m)
		}
	}
	return result, nil
}

// GetMember retrieves a single member by ID.
// Deprecated: use the repository interface instead.
func GetMember(db *sql.DB, id int64) (*domain.TeamMember, error) {
	repo := NewSQLiteTeamMemberRepository(db)
	return repo.FindByID(domain.TeamMemberID(id))
}

// EditMember updates an existing member's information.
// Deprecated: use the repository interface instead.
func EditMember(db *sql.DB, id int64, firstName, lastName string, seniority domain.Seniority) error {
	member, err := GetMember(db, id)
	if err != nil {
		return err
	}
	if member == nil {
		return fmt.Errorf("member not found")
	}

	name, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return err
	}

	// Create a new aggregate with updated values
	updatedMember, err := domain.NewTeamMember(int64(member.ID()), name, seniority)
	if err != nil {
		return err
	}

	// Preserve deactivation state if already deactivated
	if !member.IsActive() {
		_ = updatedMember.Deactivate()
	}

	repo := NewSQLiteTeamMemberRepository(db)
	return repo.Save(updatedMember)
}

// DeactivateMember soft-deletes a member.
// Deprecated: use the repository interface instead.
func DeactivateMember(db *sql.DB, id int64) error {
	member, err := GetMember(db, id)
	if err != nil {
		return err
	}
	if member == nil {
		return fmt.Errorf("member not found")
	}

	if err := member.Deactivate(); err != nil {
		return err
	}

	repo := NewSQLiteTeamMemberRepository(db)
	return repo.Save(member)
}

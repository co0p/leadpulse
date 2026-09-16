package sqlite

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"leadpulse/core/members"
	"leadpulse/engine/domain"
)

// Helper function to create an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create the members table
	_, err = db.Exec(`
		CREATE TABLE members (
			id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			seniority TEXT NOT NULL,
			created_at TEXT NOT NULL,
			deactivated_at TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create members table: %v", err)
	}

	// Create the audit table (required by logAuditEvent)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			member_id INTEGER NOT NULL,
			field TEXT NOT NULL,
			old_value TEXT,
			new_value TEXT,
			changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create audit_log table: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// TestSQLiteTeamMemberRepository_Save_Inserts_NewMember tests that
// SQLiteTeamMemberRepository.Save inserts a new member.
func TestSQLiteTeamMemberRepository_Save_Inserts_NewMember(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	// Create and save a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := members.NewTeamMember(1, name, domain.SeniorityMid)

	err := repo.Save(member)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it was inserted
	retrieved, err := repo.FindByID(members.TeamMemberID(1))
	if err != nil {
		t.Fatalf("expected no error on retrieve, got %v", err)
	}

	if retrieved == nil {
		t.Fatalf("expected member to be retrieved")
	}

	if retrieved.Name().First != "Alice" {
		t.Errorf("expected FirstName='Alice', got '%s'", retrieved.Name().First)
	}
}

// TestSQLiteTeamMemberRepository_FindByID_Returns_SavedMember tests that
// FindByID retrieves a previously saved member.
func TestSQLiteTeamMemberRepository_FindByID_Returns_SavedMember(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	// Save a member
	name, _ := domain.NewFullName("Bob", "Jones")
	member, _ := members.NewTeamMember(42, name, domain.SeniorityJunior)
	repo.Save(member)

	// Retrieve it
	retrieved, err := repo.FindByID(members.TeamMemberID(42))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.Name().First != "Bob" {
		t.Errorf("expected FirstName='Bob', got '%s'", retrieved.Name().First)
	}

	if retrieved.Seniority() != domain.SeniorityJunior {
		t.Errorf("expected Seniority='Junior', got '%s'", retrieved.Seniority())
	}
}

// TestSQLiteTeamMemberRepository_FindActive_Excludes_DeactivatedMembers tests that
// FindActive only returns non-deactivated members.
func TestSQLiteTeamMemberRepository_FindActive_Excludes_DeactivatedMembers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	// Save two members
	name1, _ := domain.NewFullName("Alice", "Active")
	member1, _ := members.NewTeamMember(1, name1, domain.SeniorityMid)
	repo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Inactive")
	member2, _ := members.NewTeamMember(2, name2, domain.SeniorityMid)
	repo.Save(member2)

	// Deactivate the second member
	repo.Deactivate(members.TeamMemberID(2))

	// Retrieve active members
	active, err := repo.FindActive()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(active) != 1 {
		t.Errorf("expected 1 active member, got %d", len(active))
	}

	if active[0].Name().First != "Alice" {
		t.Errorf("expected active member to be Alice, got %s", active[0].Name().First)
	}
}

// TestSQLiteTeamMemberRepository_Update_ModifiesMember tests that
// Save updates an existing member's information.
func TestSQLiteTeamMemberRepository_Update_ModifiesMember(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	// Save a member
	name, _ := domain.NewFullName("Charlie", "Brown")
	member, _ := members.NewTeamMember(3, name, domain.SeniorityJunior)
	repo.Save(member)

	// Retrieve, modify, and save
	retrieved, _ := repo.FindByID(members.TeamMemberID(3))
	retrieved.UpdateName("Charles", "Browning")
	repo.Save(retrieved)

	// Verify the update
	updated, _ := repo.FindByID(members.TeamMemberID(3))
	if updated.Name().First != "Charles" {
		t.Errorf("expected FirstName='Charles' after update, got '%s'", updated.Name().First)
	}
}

// TestSQLiteTeamMemberRepository_Deactivate_MarksInactive tests that
// Deactivate marks a member as inactive.
func TestSQLiteTeamMemberRepository_Deactivate_MarksInactive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	// Save a member
	name, _ := domain.NewFullName("Diana", "Prince")
	member, _ := members.NewTeamMember(4, name, domain.SeniorityMid)
	repo.Save(member)

	// Deactivate
	err := repo.Deactivate(members.TeamMemberID(4))
	if err != nil {
		t.Fatalf("expected no error on deactivate, got %v", err)
	}

	// Verify it's marked as inactive
	retrieved, _ := repo.FindByID(members.TeamMemberID(4))
	if retrieved.IsActive() {
		t.Errorf("expected member to be inactive after deactivation")
	}

	if retrieved.DeactivatedAt() == nil {
		t.Errorf("expected DeactivatedAt to be set")
	}
}

// TestSQLiteTeamMemberRepository_FindByID_NotFound_ReturnsNil tests that
// FindByID returns nil when member is not found.
func TestSQLiteTeamMemberRepository_FindByID_NotFound_ReturnsNil(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLiteTeamMemberRepository(db)

	retrieved, err := repo.FindByID(members.TeamMemberID(999))
	if err != nil {
		t.Fatalf("expected no error for not-found member, got %v", err)
	}

	if retrieved != nil {
		t.Errorf("expected nil for not-found member, got %v", retrieved)
	}
}

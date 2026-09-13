package store

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"leadpulse/engine/domain"
)

// TestTeamMemberRepositorySave tests that Save creates a new member.
func TestTeamMemberRepositorySave(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTeamMemberRepository(db)

	// Create a team member aggregate
	name, _ := domain.NewFullName("Alice", "Chen")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)

	// Save the member
	err := repo.Save(member)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify it was persisted
	var firstName, lastName, seniority string
	err = db.QueryRow("SELECT first_name, last_name, seniority FROM members WHERE id = ?", 1).
		Scan(&firstName, &lastName, &seniority)
	if err != nil {
		t.Fatalf("member not found in database: %v", err)
	}

	if firstName != "Alice" {
		t.Errorf("expected firstName 'Alice', got %q", firstName)
	}
	if lastName != "Chen" {
		t.Errorf("expected lastName 'Chen', got %q", lastName)
	}
	if seniority != string(domain.SenioritySenior) {
		t.Errorf("expected seniority 'Senior', got %q", seniority)
	}
}

// TestTeamMemberRepositoryFindActive tests that FindActive returns only active members.
func TestTeamMemberRepositoryFindActive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTeamMemberRepository(db)

	// Create and save two members
	name1, _ := domain.NewFullName("Alice", "Chen")
	member1, _ := domain.NewTeamMember(1, name1, domain.SenioritySenior)
	repo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Smith")
	member2, _ := domain.NewTeamMember(2, name2, domain.SeniorityMid)
	repo.Save(member2)

	// List active members (none deactivated yet)
	members, err := repo.FindActive()
	if err != nil {
		t.Fatalf("FindActive failed: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}

	// Verify the members
	if members[0].ID() != 1 || members[0].Name().First != "Alice" {
		t.Errorf("first member mismatch")
	}
	if members[1].ID() != 2 || members[1].Name().First != "Bob" {
		t.Errorf("second member mismatch")
	}
}

// TestTeamMemberRepositoryFindByID tests that FindByID retrieves a member.
func TestTeamMemberRepositoryFindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTeamMemberRepository(db)

	// Create and save a member
	name, _ := domain.NewFullName("Bob", "Smith")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	// Retrieve it by ID
	retrieved, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("member not found")
	}

	if retrieved.ID() != 1 || retrieved.Name().First != "Bob" {
		t.Errorf("member mismatch")
	}
}

// TestTeamMemberRepositorySaveUpdate tests that Save updates an existing member.
func TestTeamMemberRepositorySaveUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTeamMemberRepository(db)

	// Create and save a member
	name1, _ := domain.NewFullName("Bob", "Smith")
	member, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	repo.Save(member)

	// Update the member
	name2, _ := domain.NewFullName("Robert", "Smith")
	updatedMember, _ := domain.NewTeamMember(1, name2, domain.SenioritySenior)
	repo.Save(updatedMember)

	// Verify the change
	var firstName, lastName, seniority string
	err := db.QueryRow("SELECT first_name, last_name, seniority FROM members WHERE id = ?", 1).
		Scan(&firstName, &lastName, &seniority)
	if err != nil {
		t.Fatalf("failed to query updated member: %v", err)
	}

	if firstName != "Robert" {
		t.Errorf("expected firstName 'Robert', got %q", firstName)
	}
	if seniority != string(domain.SenioritySenior) {
		t.Errorf("expected seniority 'Senior', got %q", seniority)
	}

	// Verify audit log entries were created
	var auditCount int
	err = db.QueryRow("SELECT COUNT(*) FROM audit_log WHERE member_id = ?", 1).Scan(&auditCount)
	if err != nil {
		t.Fatalf("failed to count audit entries: %v", err)
	}
	if auditCount < 2 {
		t.Errorf("expected at least 2 audit entries, got %d", auditCount)
	}
}

// TestTeamMemberRepositoryDeactivate tests that deactivation works correctly.
func TestTeamMemberRepositoryDeactivate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTeamMemberRepository(db)

	// Create and save a member
	name, _ := domain.NewFullName("Carol", "Davis")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityPrincipal)
	repo.Save(member)

	// Deactivate the member
	member.Deactivate()
	repo.Save(member)

	// Verify the member is no longer in FindActive
	members, _ := repo.FindActive()
	for _, m := range members {
		if m.ID() == 1 {
			t.Errorf("deactivated member should not appear in FindActive")
		}
	}

	// Verify the member's data is still in the database (soft delete)
	var firstName string
	var deactivatedAt sql.NullString
	err := db.QueryRow("SELECT first_name, deactivated_at FROM members WHERE id = ?", 1).
		Scan(&firstName, &deactivatedAt)
	if err != nil {
		t.Fatalf("deactivated member should still exist in database: %v", err)
	}

	if firstName != "Carol" {
		t.Errorf("member data should be preserved, got firstName %q", firstName)
	}
	if !deactivatedAt.Valid {
		t.Errorf("deactivated_at should be set")
	}
}

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	if err := InitSchema(db); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

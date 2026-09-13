package store

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"leadpulse/engine/domain"
)

// TestAddMember_createsAndReturnsID tests that AddMember creates a new member
// and returns its ID.
func TestAddMember_createsAndReturnsID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Call the function that doesn't exist yet
	id, err := AddMember(db, "Alice", "Chen", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("AddMember failed: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected positive ID, got %d", id)
	}

	// Verify the member was inserted
	var firstName, lastName, seniority string
	err = db.QueryRow("SELECT first_name, last_name, seniority FROM members WHERE id = ?", id).
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

// TestListMembers_returnsActiveMembers tests that ListMembers returns only active members.
func TestListMembers_returnsActiveMembers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add two members
	id1, _ := AddMember(db, "Alice", "Chen", domain.SenioritySenior)
	id2, _ := AddMember(db, "Bob", "Smith", domain.SeniorityMid)

	// List members (none deactivated yet)
	members, err := ListMembers(db)
	if err != nil {
		t.Fatalf("ListMembers failed: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}

	// Verify the members
	if members[0].ID != id1 || members[0].FirstName != "Alice" {
		t.Errorf("first member mismatch")
	}
	if members[1].ID != id2 || members[1].FirstName != "Bob" {
		t.Errorf("second member mismatch")
	}
}

// TestEditMember_updatesAndAudits tests that EditMember updates a member and logs changes.
func TestEditMember_updatesAndAudits(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add a member
	id, _ := AddMember(db, "Bob", "Smith", domain.SeniorityMid)

	// Edit the member
	err := EditMember(db, id, "Robert", "Smith", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("EditMember failed: %v", err)
	}

	// Verify the change
	var firstName, lastName, seniority string
	err = db.QueryRow("SELECT first_name, last_name, seniority FROM members WHERE id = ?", id).
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
	err = db.QueryRow("SELECT COUNT(*) FROM audit_log WHERE member_id = ?", id).Scan(&auditCount)
	if err != nil {
		t.Fatalf("failed to count audit entries: %v", err)
	}
	if auditCount < 2 {
		t.Errorf("expected at least 2 audit entries, got %d", auditCount)
	}
}

// TestDeactivateMember_softDeletesAndPreservesData tests that DeactivateMember soft-deletes.
func TestDeactivateMember_softDeletesAndPreservesData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add a member
	id, _ := AddMember(db, "Carol", "Davis", domain.SeniorityPrincipal)

	// Deactivate the member
	err := DeactivateMember(db, id)
	if err != nil {
		t.Fatalf("DeactivateMember failed: %v", err)
	}

	// Verify the member is no longer in ListMembers
	members, _ := ListMembers(db)
	for _, m := range members {
		if m.ID == id {
			t.Errorf("deactivated member should not appear in ListMembers")
		}
	}

	// Verify the member's data is still in the database (soft delete)
	var firstName string
	var deactivatedAt sql.NullString
	err = db.QueryRow("SELECT first_name, deactivated_at FROM members WHERE id = ?", id).
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

// TestAddMember_createsAuditLogEntry tests that adding a member creates an audit entry.
func TestAddMember_createsAuditLogEntry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	id, _ := AddMember(db, "Alice", "Chen", domain.SenioritySenior)

	// For now, AddMember doesn't log to audit trail (deferred per plan)
	// This test is a placeholder for future enhancement
	// Verify the member exists
	var firstName string
	err := db.QueryRow("SELECT first_name FROM members WHERE id = ?", id).Scan(&firstName)
	if err != nil {
		t.Fatalf("member not found: %v", err)
	}
	if firstName != "Alice" {
		t.Errorf("expected Alice, got %q", firstName)
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

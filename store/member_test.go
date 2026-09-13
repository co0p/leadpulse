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

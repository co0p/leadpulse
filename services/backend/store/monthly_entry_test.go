package store

import (
	"testing"

	_ "modernc.org/sqlite"

	"leadpulse/engine/domain"
)

// TestMonthlyEntryRepositorySave tests that Save creates and persists a monthly entry.
func TestMonthlyEntryRepositorySave(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteMonthlyEntryRepository(db)

	// Create a monthly entry aggregate
	signals, _ := domain.NewMonthlyRawSignals(
		3,  // morale
		85, // billability
		4,  // csat
		15, // net_margin
		5,  // positive_feedback
		1,  // critical_feedback
		4,  // overtime_hours
		90, // delivery_reliability
		2,  // mentoring_hours
		8,  // evidence_notes_count
	)

	entry, _ := domain.NewMonthlyEntry(1, "2024-10", signals)

	// Save the entry
	err := repo.Save(entry)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify it was persisted
	var month string
	var morale, billability int
	err = db.QueryRow(
		"SELECT month, morale, billability FROM monthly_entries WHERE member_id = ? AND month = ?",
		1, "2024-10",
	).Scan(&month, &morale, &billability)
	if err != nil {
		t.Fatalf("monthly entry not found in database: %v", err)
	}

	if month != "2024-10" {
		t.Errorf("expected month 2024-10, got %s", month)
	}
	if morale != 3 {
		t.Errorf("expected morale 3, got %d", morale)
	}
	if billability != 85 {
		t.Errorf("expected billability 85, got %d", billability)
	}
}

// TestMonthlyEntryRepositoryFindByID tests that FindByID retrieves a monthly entry.
func TestMonthlyEntryRepositoryFindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteMonthlyEntryRepository(db)

	// Create and save an entry
	signals, _ := domain.NewMonthlyRawSignals(
		4, 80, 5, 20, 3, 0, 2, 95, 1, 6,
	)
	entry, _ := domain.NewMonthlyEntry(1, "2024-11", signals)
	repo.Save(entry)

	// Retrieve it by ID
	id := domain.MonthlyEntryID{MemberID: domain.TeamMemberID(1), Month: "2024-11"}
	retrieved, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("entry not found")
	}

	if retrieved.ID().MemberID != 1 || retrieved.ID().Month != "2024-11" {
		t.Errorf("ID mismatch")
	}

	if retrieved.Signals().Morale == nil || *retrieved.Signals().Morale != 4 {
		t.Errorf("morale mismatch: expected 4, got %v", retrieved.Signals().Morale)
	}
}

// TestMonthlyEntryRepositoryFindByMember tests that FindByMember returns all entries for a member.
func TestMonthlyEntryRepositoryFindByMember(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteMonthlyEntryRepository(db)

	// Create and save multiple entries for the same member
	signals1, _ := domain.NewMonthlyRawSignals(3, 85, 4, 15, 5, 1, 4, 90, 2, 8)
	entry1, _ := domain.NewMonthlyEntry(1, "2024-09", signals1)
	repo.Save(entry1)

	signals2, _ := domain.NewMonthlyRawSignals(4, 80, 5, 20, 3, 0, 2, 95, 1, 6)
	entry2, _ := domain.NewMonthlyEntry(1, "2024-10", signals2)
	repo.Save(entry2)

	// Create an entry for a different member
	signals3, _ := domain.NewMonthlyRawSignals(2, 70, 3, 10, 2, 1, 0, 85, 0, 4)
	entry3, _ := domain.NewMonthlyEntry(2, "2024-10", signals3)
	repo.Save(entry3)

	// Retrieve all entries for member 1
	entries, err := repo.FindByMember(domain.TeamMemberID(1))
	if err != nil {
		t.Fatalf("FindByMember failed: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Verify ordering (DESC by month)
	if entries[0].ID().Month != "2024-10" || entries[1].ID().Month != "2024-09" {
		t.Errorf("entries not ordered by month DESC")
	}
}

// TestMonthlyEntryRepositoryRejectsInvalidState tests that invalid signals are rejected.
func TestMonthlyEntryRepositoryRejectsInvalidState(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_ = NewSQLiteMonthlyEntryRepository(db) // Verify repo creation works

	// Try to create an entry with invalid morale (out of range)
	_, err := domain.NewMonthlyRawSignals(
		6,  // Invalid: morale must be 0–5
		85, // billability
		4, 15, 5, 1, 4, 90, 2, 8,
	)
	if err == nil {
		t.Fatal("expected error for invalid morale, got none")
	}

	// Try to create an entry with invalid billability
	_, err = domain.NewMonthlyRawSignals(
		3,
		150, // Invalid: billability must be 0–100
		4, 15, 5, 1, 4, 90, 2, 8,
	)
	if err == nil {
		t.Fatal("expected error for invalid billability, got none")
	}
}

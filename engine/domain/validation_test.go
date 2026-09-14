package domain

import (
	"testing"
)

// TestValidateTeamMemberUniqueness_rejectsDuplicate tests that ValidationService
// rejects team members with duplicate names.
// Business rule: team member names must be unique within active members.
func TestValidateTeamMemberUniqueness_rejectsDuplicate(t *testing.T) {
	// Arrange
	memberRepo := NewInMemoryTeamMemberRepository()
	entryRepo := NewInMemoryMonthlyEntryRepository()
	service := NewValidationService(memberRepo, entryRepo)

	// Create first member with ID 1
	name1, _ := NewFullName("Alice", "Smith")
	member1, _ := NewTeamMember(1, name1, SeniorityMid)
	memberRepo.Save(member1)

	// Create second member with ID 2 but same name
	name2, _ := NewFullName("Alice", "Smith")
	member2, _ := NewTeamMember(2, name2, SeniorityMid)

	// Act
	err := service.ValidateTeamMemberUniqueness(member2)

	// Assert
	if err == nil {
		t.Error("ValidateTeamMemberUniqueness should reject duplicate names")
	}
}

// TestValidateEntryUniqueness_rejectsDuplicate tests that ValidationService
// rejects duplicate entries for the same member and month.
// Business rule: one entry per member per month.
func TestValidateEntryUniqueness_rejectsDuplicate(t *testing.T) {
	// Arrange
	memberRepo := NewInMemoryTeamMemberRepository()
	entryRepo := NewInMemoryMonthlyEntryRepository()
	service := NewValidationService(memberRepo, entryRepo)

	// Create and save an entry for member 1, month "2025-01"
	// First, create a minimal MonthlyRawSignals with at least one signal filled
	// morale=3 (0-5), csat=3 (1-5)
	signals, _ := NewMonthlyRawSignals(3, 0, 3, 0, 0, 0, 0, 0, 0, 0)

	// Create and save the entry
	entry, createErr := NewMonthlyEntry(1, "2025-01", signals)
	if createErr != nil {
		t.Fatalf("Failed to create entry: %v", createErr)
	}
	saveErr := entryRepo.Save(entry)
	if saveErr != nil {
		t.Fatalf("Failed to save entry: %v", saveErr)
	}

	// Act
	// Validation should fail because entry already exists
	err := service.ValidateEntryUniqueness(TeamMemberID(1), "2025-01")

	// Assert
	if err == nil {
		t.Error("ValidateEntryUniqueness should reject duplicate entries")
	}
}

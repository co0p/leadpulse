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

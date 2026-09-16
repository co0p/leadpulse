package members

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestDeactivateMemberUseCase_Success_DeactivatesMember tests that
// DeactivateMemberUseCase successfully deactivates an active member.
func TestDeactivateMemberUseCase_Success_DeactivatesMember(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	uc := NewDeactivateMemberUseCase(repo)
	input := DeactivateMemberInput{MemberID: 1}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatalf("expected output to be non-nil")
	}

	if output.Status != "Inactive" {
		t.Errorf("expected Status='Inactive', got '%s'", output.Status)
	}

	if output.DeactivatedAt == nil {
		t.Errorf("expected DeactivatedAt to be set")
	}

	// Verify persistence: member should no longer appear in FindActive
	active, _ := repo.FindActive()
	if len(active) != 0 {
		t.Errorf("expected no active members after deactivation, got %d", len(active))
	}
}

// TestDeactivateMemberUseCase_MemberNotFound_ReturnsError tests that
// DeactivateMemberUseCase returns an error when the member is not found.
func TestDeactivateMemberUseCase_MemberNotFound_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewDeactivateMemberUseCase(repo)

	input := DeactivateMemberInput{MemberID: 999}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for not-found member, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestDeactivateMemberUseCase_AlreadyDeactivated_ReturnsError tests that
// DeactivateMemberUseCase returns an error when trying to deactivate an already deactivated member.
func TestDeactivateMemberUseCase_AlreadyDeactivated_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member and deactivate it
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)
	repo.Deactivate(TeamMemberID(1))

	uc := NewDeactivateMemberUseCase(repo)
	input := DeactivateMemberInput{MemberID: 1}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for already-deactivated member, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestDeactivateMemberUseCase_InvalidMemberID_ReturnsError tests that
// DeactivateMemberUseCase returns an error for invalid (non-positive) member ID.
func TestDeactivateMemberUseCase_InvalidMemberID_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewDeactivateMemberUseCase(repo)

	input := DeactivateMemberInput{MemberID: -1}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for negative ID, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

package members

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestReactivateMemberUseCase_Success_ReactivatesMember tests that
// ReactivateMemberUseCase successfully reactivates a deactivated member.
func TestReactivateMemberUseCase_Success_ReactivatesMember(t *testing.T) {
	repo := newMockMemberRepository()

	// Add a member and deactivate it
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)
	repo.Deactivate(TeamMemberID(1))

	uc := NewReactivateMemberUseCase(repo)
	input := ReactivateMemberInput{MemberID: 1}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatalf("expected output to be non-nil")
	}

	if output.Status != "Active" {
		t.Errorf("expected Status='Active', got '%s'", output.Status)
	}

	if output.ID != 1 {
		t.Errorf("expected ID=1, got %d", output.ID)
	}
}

// TestReactivateMemberUseCase_AlreadyActive_ReturnsError tests that
// ReactivateMemberUseCase returns an error when trying to reactivate an already-active member.
func TestReactivateMemberUseCase_AlreadyActive_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()

	// Add an active member
	name, _ := domain.NewFullName("Bob", "Jones")
	member, _ := NewTeamMember(1, name, domain.SeniorityJunior)
	repo.Save(member)

	uc := NewReactivateMemberUseCase(repo)
	input := ReactivateMemberInput{MemberID: 1}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for already-active member, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestReactivateMemberUseCase_MemberNotFound_ReturnsError tests that
// ReactivateMemberUseCase returns an error when the member is not found.
func TestReactivateMemberUseCase_MemberNotFound_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewReactivateMemberUseCase(repo)

	input := ReactivateMemberInput{MemberID: 999}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for not-found member, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestReactivateMemberUseCase_InvalidMemberID_ReturnsError tests that
// ReactivateMemberUseCase returns an error for invalid (non-positive) member ID.
func TestReactivateMemberUseCase_InvalidMemberID_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewReactivateMemberUseCase(repo)

	input := ReactivateMemberInput{MemberID: -1}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for negative ID, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

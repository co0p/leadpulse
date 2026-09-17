package members

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestTeamMember_Reactivate_Success tests that Reactivate() successfully
// reactivates a deactivated member by setting deactivatedAt to nil.
func TestTeamMember_Reactivate_Success(t *testing.T) {
	// Arrange: create and deactivate a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)

	// Deactivate the member first
	_ = member.Deactivate()
	if member.IsActive() {
		t.Fatalf("precondition failed: member should be deactivated after Deactivate()")
	}

	// Act: reactivate the member
	err := member.Reactivate()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !member.IsActive() {
		t.Errorf("expected member to be active after Reactivate(), but IsActive() returned false")
	}

	if member.DeactivatedAt() != nil {
		t.Errorf("expected DeactivatedAt to be nil after Reactivate(), but got %v", member.DeactivatedAt())
	}
}

// TestTeamMember_Reactivate_AlreadyActive_ReturnsError tests that Reactivate()
// returns an error when called on an already-active member.
func TestTeamMember_Reactivate_AlreadyActive_ReturnsError(t *testing.T) {
	// Arrange: create an active member (default state)
	name, _ := domain.NewFullName("Bob", "Jones")
	member, _ := NewTeamMember(2, name, domain.SeniorityJunior)

	if !member.IsActive() {
		t.Fatalf("precondition failed: member should be active by default")
	}

	// Act: attempt to reactivate an already-active member
	err := member.Reactivate()

	// Assert
	if err == nil {
		t.Fatalf("expected error when reactivating an already-active member, but got nil")
	}

	if member.IsActive() != true {
		t.Errorf("expected member to still be active, but IsActive() returned false")
	}

	if member.DeactivatedAt() != nil {
		t.Errorf("expected DeactivatedAt to be nil, but got %v", member.DeactivatedAt())
	}
}

package member

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestAddMember_createsAndReturns tests that AddMember calls the repository and returns the created member.
func TestAddMember_createsAndReturns(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Call AddMember
	member, err := svc.AddMember("Alice", "Chen", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("AddMember failed: %v", err)
	}

	if member == nil {
		t.Fatalf("expected member, got nil")
	}

	if member.Name().First != "Alice" || member.Name().Last != "Chen" || member.Seniority() != domain.SenioritySenior {
		t.Errorf("member data mismatch: %+v", member)
	}

	if member.ID() <= 0 {
		t.Errorf("expected positive ID, got %d", member.ID())
	}
}

// TestAddMember_rejectsEmptyName tests that AddMember rejects empty names.
func TestAddMember_rejectsEmptyName(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Call AddMember with empty first name
	_, err := svc.AddMember("", "Chen", domain.SenioritySenior)
	if err == nil {
		t.Fatalf("expected error for empty first name")
	}
}

// TestAddMember_rejectsInvalidSeniority tests that AddMember rejects invalid seniority.
func TestAddMember_rejectsInvalidSeniority(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Call AddMember with invalid seniority
	_, err := svc.AddMember("Alice", "Chen", domain.Seniority("Invalid"))
	if err == nil {
		t.Fatalf("expected error for invalid seniority")
	}
}

// TestListMembers_returnsActiveMembers tests that ListMembers returns active members.
func TestListMembers_returnsActiveMembers(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Add members
	svc.AddMember("Alice", "Chen", domain.SenioritySenior)
	svc.AddMember("Bob", "Smith", domain.SeniorityMid)

	// List members
	members, err := svc.ListMembers()
	if err != nil {
		t.Fatalf("ListMembers failed: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

// TestEditMember_updatesFields tests that EditMember updates member fields.
func TestEditMember_updatesFields(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Add a member
	member, _ := svc.AddMember("Bob", "Smith", domain.SeniorityMid)

	// Edit the member
	updated, err := svc.EditMember(int64(member.ID()), "Robert", "Smith", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("EditMember failed: %v", err)
	}

	if updated.Name().First != "Robert" || updated.Seniority() != domain.SenioritySenior {
		t.Errorf("member not updated correctly: %+v", updated)
	}
}

// TestEditMember_rejectsInvalidSeniority tests that EditMember rejects invalid seniority.
func TestEditMember_rejectsInvalidSeniority(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Add a member
	member, _ := svc.AddMember("Bob", "Smith", domain.SeniorityMid)

	// Try to edit with invalid seniority
	_, err := svc.EditMember(int64(member.ID()), "Bob", "Smith", domain.Seniority("Invalid"))
	if err == nil {
		t.Fatalf("expected error for invalid seniority")
	}
}

// TestDeactivateMember_removes tests that DeactivateMember hides the member.
func TestDeactivateMember_removes(t *testing.T) {
	repo := domain.NewInMemoryTeamMemberRepository()
	svc := NewService(repo)

	// Add a member
	member, _ := svc.AddMember("Carol", "Davis", domain.SeniorityPrincipal)

	// Deactivate
	err := svc.DeactivateMember(int64(member.ID()))
	if err != nil {
		t.Fatalf("DeactivateMember failed: %v", err)
	}

	// Verify member is no longer in the list
	members, _ := svc.ListMembers()
	for _, m := range members {
		if m.ID() == member.ID() {
			t.Errorf("deactivated member should not appear in list")
		}
	}
}

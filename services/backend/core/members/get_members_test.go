package members

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestGetMembersUseCase_NoMembers_ReturnsEmptyList tests that
// GetMembersUseCase returns an empty list when no active members exist.
func TestGetMembersUseCase_NoMembers_ReturnsEmptyList(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewGetMembersUseCase(repo)

	output, err := uc.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatalf("expected output to be non-nil")
	}

	if len(output.Members) != 0 {
		t.Errorf("expected empty members list, got %d members", len(output.Members))
	}
}

// TestGetMembersUseCase_SingleMember_ReturnsOneItem tests that
// GetMembersUseCase returns one member when one active member exists.
func TestGetMembersUseCase_SingleMember_ReturnsOneItem(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SenioritySenior)
	repo.Save(member)

	uc := NewGetMembersUseCase(repo)
	output, err := uc.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(output.Members))
	}

	if output.Members[0].FirstName != "Alice" {
		t.Errorf("expected FirstName='Alice', got '%s'", output.Members[0].FirstName)
	}
}

// TestGetMembersUseCase_MultipleMembers_ReturnsAll tests that
// GetMembersUseCase returns all active members.
func TestGetMembersUseCase_MultipleMembers_ReturnsAll(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add three members
	for i := 1; i <= 3; i++ {
		name, _ := domain.NewFullName("Member", "Test"+string(rune('0'+i)))
		member, _ := NewTeamMember(int64(i), name, domain.SeniorityMid)
		repo.Save(member)
	}

	uc := NewGetMembersUseCase(repo)
	output, err := uc.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Members) != 3 {
		t.Errorf("expected 3 members, got %d", len(output.Members))
	}
}

// TestGetMembersUseCase_ExcludesInactiveMembers tests that
// GetMembersUseCase excludes deactivated members.
func TestGetMembersUseCase_ExcludesInactiveMembers(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add two members
	name1, _ := domain.NewFullName("Alice", "Active")
	member1, _ := NewTeamMember(1, name1, domain.SenioritySenior)
	repo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Inactive")
	member2, _ := NewTeamMember(2, name2, domain.SeniorityJunior)
	repo.Save(member2)

	// Deactivate the second member
	repo.Deactivate(TeamMemberID(2))

	uc := NewGetMembersUseCase(repo)
	output, err := uc.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Members) != 1 {
		t.Errorf("expected 1 active member, got %d", len(output.Members))
	}

	if output.Members[0].FirstName != "Alice" {
		t.Errorf("expected active member to be Alice, got %s", output.Members[0].FirstName)
	}
}

// TestGetFilteredMembersUseCase_Status_Active_ReturnsOnlyActive tests that
// GetFilteredMembersUseCase with status="active" returns only active members.
func TestGetFilteredMembersUseCase_Status_Active_ReturnsOnlyActive(t *testing.T) {
	repo := newMockMemberRepository()

	// Add two active members
	name1, _ := domain.NewFullName("Alice", "Active1")
	member1, _ := NewTeamMember(1, name1, domain.SenioritySenior)
	repo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Active2")
	member2, _ := NewTeamMember(2, name2, domain.SeniorityMid)
	repo.Save(member2)

	// Add one inactive member
	name3, _ := domain.NewFullName("Charlie", "Inactive")
	member3, _ := NewTeamMember(3, name3, domain.SeniorityJunior)
	repo.Save(member3)
	repo.Deactivate(TeamMemberID(3))

	uc := NewGetFilteredMembersUseCase(repo)
	input := GetFilteredMembersInput{Status: "active"}
	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatalf("expected output to be non-nil")
	}

	if len(output.Members) != 2 {
		t.Errorf("expected 2 active members, got %d", len(output.Members))
	}

	for _, member := range output.Members {
		if member.Status != "Active" {
			t.Errorf("expected Status='Active', got '%s' for member %s", member.Status, member.FirstName)
		}
	}
}

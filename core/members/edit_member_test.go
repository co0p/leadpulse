package members

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestEditMemberUseCase_UpdateAllFields_Success tests that
// EditMemberUseCase successfully updates all fields of a member.
func TestEditMemberUseCase_UpdateAllFields_Success(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	uc := NewEditMemberUseCase(repo)
	input := EditMemberInput{
		MemberID:  1,
		FirstName: "Alicia",
		LastName:  "Jones",
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.FirstName != "Alicia" {
		t.Errorf("expected FirstName='Alicia', got '%s'", output.FirstName)
	}

	if output.LastName != "Jones" {
		t.Errorf("expected LastName='Jones', got '%s'", output.LastName)
	}

	if output.Seniority != "Senior" {
		t.Errorf("expected Seniority='Senior', got '%s'", output.Seniority)
	}
}

// TestEditMemberUseCase_PartialUpdate_OnlyChangedFields tests that
// EditMemberUseCase supports partial updates (empty strings mean no change).
func TestEditMemberUseCase_PartialUpdate_OnlyChangedFields(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	uc := NewEditMemberUseCase(repo)
	input := EditMemberInput{
		MemberID:  1,
		FirstName: "",        // No change
		LastName:  "Jones",   // Change last name only
		Seniority: "",        // No change
	}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.FirstName != "Alice" {
		t.Errorf("expected FirstName='Alice' (unchanged), got '%s'", output.FirstName)
	}

	if output.LastName != "Jones" {
		t.Errorf("expected LastName='Jones', got '%s'", output.LastName)
	}

	if output.Seniority != "Mid" {
		t.Errorf("expected Seniority='Mid' (unchanged), got '%s'", output.Seniority)
	}
}

// TestEditMemberUseCase_MemberNotFound_ReturnsError tests that
// EditMemberUseCase returns an error when the member is not found.
func TestEditMemberUseCase_MemberNotFound_ReturnsError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewEditMemberUseCase(repo)

	input := EditMemberInput{
		MemberID:  999,
		FirstName: "Alice",
		LastName:  "Smith",
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error for not-found member, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestEditMemberUseCase_InvalidSeniority_ReturnsValidationError tests that
// EditMemberUseCase returns a validation error for invalid seniority.
func TestEditMemberUseCase_InvalidSeniority_ReturnsValidationError(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	uc := NewEditMemberUseCase(repo)
	input := EditMemberInput{
		MemberID:  1,
		FirstName: "",
		LastName:  "",
		Seniority: "InvalidSeniority",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestEditMemberUseCase_OnlyUpdateSeniority_KeepsNameUnchanged tests that
// EditMemberUseCase supports updating only seniority while keeping name unchanged.
func TestEditMemberUseCase_OnlyUpdateSeniority_KeepsNameUnchanged(t *testing.T) {
	repo := newMockMemberRepository()
	
	// Add a member
	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := NewTeamMember(1, name, domain.SeniorityMid)
	repo.Save(member)

	uc := NewEditMemberUseCase(repo)
	input := EditMemberInput{
		MemberID:  1,
		FirstName: "", // Empty first name
		LastName:  "", // Empty last name
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.FirstName != "Alice" {
		t.Errorf("expected FirstName='Alice' (unchanged), got '%s'", output.FirstName)
	}

	if output.LastName != "Smith" {
		t.Errorf("expected LastName='Smith' (unchanged), got '%s'", output.LastName)
	}

	if output.Seniority != "Senior" {
		t.Errorf("expected Seniority='Senior', got '%s'", output.Seniority)
	}
}

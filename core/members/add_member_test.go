package members

import (
	"fmt"
	"sync"
	"testing"
)

// mockMemberRepository is a test double implementing MemberRepository
// for use in unit tests of use cases.
type mockMemberRepository struct {
	mu      sync.RWMutex
	members map[int64]*TeamMember
}

func newMockMemberRepository() *mockMemberRepository {
	return &mockMemberRepository{
		members: make(map[int64]*TeamMember),
	}
}

func (m *mockMemberRepository) Save(member *TeamMember) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if member == nil {
		return nil
	}
	m.members[int64(member.ID())] = member
	return nil
}

func (m *mockMemberRepository) FindByID(id TeamMemberID) (*TeamMember, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if member, ok := m.members[int64(id)]; ok {
		return member, nil
	}
	return nil, nil
}

func (m *mockMemberRepository) FindActive() ([]*TeamMember, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var active []*TeamMember
	for _, member := range m.members {
		if member.IsActive() {
			active = append(active, member)
		}
	}
	if active == nil {
		active = []*TeamMember{}
	}
	return active, nil
}

func (m *mockMemberRepository) Deactivate(id TeamMemberID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if member, ok := m.members[int64(id)]; ok {
		return member.Deactivate()
	}
	return nil
}

// TestAddMemberUseCase_Success_CreatesAndPersists tests that AddMemberUseCase
// successfully creates a new team member and persists it to the repository.
func TestAddMemberUseCase_Success_CreatesAndPersists(t *testing.T) {
	// Arrange
	repo := newMockMemberRepository()
	uc := NewAddMemberUseCase(repo)

	input := AddMemberInput{
		FirstName: "Alice",
		LastName:  "Smith",
		Seniority: "Senior",
	}

	// Act
	output, err := uc.Execute(input)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatalf("expected output to be non-nil")
	}

	if output.FirstName != "Alice" {
		t.Errorf("expected FirstName='Alice', got '%s'", output.FirstName)
	}

	if output.LastName != "Smith" {
		t.Errorf("expected LastName='Smith', got '%s'", output.LastName)
	}

	if output.Seniority != "Senior" {
		t.Errorf("expected Seniority='Senior', got '%s'", output.Seniority)
	}

	if output.Status != "Active" {
		t.Errorf("expected Status='Active', got '%s'", output.Status)
	}

	if output.ID == 0 {
		t.Errorf("expected ID to be non-zero")
	}

	// Verify persistence: retrieve the member from the repository
	memberID := TeamMemberID(output.ID)
	retrievedMember, err := repo.FindByID(memberID)
	if err != nil {
		t.Fatalf("expected member to be persisted, got error: %v", err)
	}

	if retrievedMember == nil {
		t.Fatalf("expected retrieved member to be non-nil")
	}

	if retrievedMember.Name().First != "Alice" || retrievedMember.Name().Last != "Smith" {
		t.Errorf("expected persisted member name to be 'Alice Smith', got '%s %s'",
			retrievedMember.Name().First, retrievedMember.Name().Last)
	}
}

// TestAddMemberUseCase_EmptyFirstName_ReturnsValidationError tests that
// AddMemberUseCase returns an error when firstName is empty.
func TestAddMemberUseCase_EmptyFirstName_ReturnsValidationError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewAddMemberUseCase(repo)

	input := AddMemberInput{
		FirstName: "",
		LastName:  "Smith",
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestAddMemberUseCase_EmptyLastName_ReturnsValidationError tests that
// AddMemberUseCase returns an error when lastName is empty.
func TestAddMemberUseCase_EmptyLastName_ReturnsValidationError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewAddMemberUseCase(repo)

	input := AddMemberInput{
		FirstName: "Alice",
		LastName:  "",
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestAddMemberUseCase_InvalidSeniority_ReturnsValidationError tests that
// AddMemberUseCase returns an error when seniority is invalid.
func TestAddMemberUseCase_InvalidSeniority_ReturnsValidationError(t *testing.T) {
	repo := newMockMemberRepository()
	uc := NewAddMemberUseCase(repo)

	input := AddMemberInput{
		FirstName: "Alice",
		LastName:  "Smith",
		Seniority: "InvalidSeniority",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// TestAddMemberUseCase_RepositorySaveError_ReturnsError tests that
// AddMemberUseCase returns an error when the repository Save fails.
func TestAddMemberUseCase_RepositorySaveError_ReturnsError(t *testing.T) {
	// Create a mock that returns an error on Save
	repo := &errorMockMemberRepository{}
	uc := NewAddMemberUseCase(repo)

	input := AddMemberInput{
		FirstName: "Alice",
		LastName:  "Smith",
		Seniority: "Senior",
	}

	output, err := uc.Execute(input)

	if err == nil {
		t.Fatalf("expected error from repository, got nil")
	}

	if output != nil {
		t.Errorf("expected output to be nil on error, got %v", output)
	}
}

// errorMockMemberRepository is a test double that always returns an error on Save
type errorMockMemberRepository struct{}

func (e *errorMockMemberRepository) Save(member *TeamMember) error {
	return fmt.Errorf("repository error")
}

func (e *errorMockMemberRepository) FindByID(id TeamMemberID) (*TeamMember, error) {
	return nil, nil
}

func (e *errorMockMemberRepository) FindActive() ([]*TeamMember, error) {
	return []*TeamMember{}, nil
}

func (e *errorMockMemberRepository) Deactivate(id TeamMemberID) error {
	return nil
}

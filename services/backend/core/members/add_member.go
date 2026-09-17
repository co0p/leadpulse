package members

import (
	"fmt"
	"time"

	"leadpulse/engine/domain"
)

// AddMemberInput is the input struct for AddMemberUseCase.
type AddMemberInput struct {
	FirstName string
	LastName  string
	Seniority string
}

// AddMemberOutput is the output struct for AddMemberUseCase.
type AddMemberOutput struct {
	ID        int64
	FirstName string
	LastName  string
	Seniority string
	Status    string
	CreatedAt time.Time
}

// AddMemberUseCase handles the creation of new team members.
type AddMemberUseCase struct {
	repo MemberRepository
}

// NewAddMemberUseCase creates a new instance of AddMemberUseCase with an injected repository.
func NewAddMemberUseCase(repo MemberRepository) *AddMemberUseCase {
	return &AddMemberUseCase{repo: repo}
}

// Execute creates and persists a new team member.
// Validates input (non-empty names, valid seniority) before creating the aggregate.
// Returns an error if validation fails or persistence fails.
func (uc *AddMemberUseCase) Execute(input AddMemberInput) (*AddMemberOutput, error) {
	// Validation
	if input.FirstName == "" || input.LastName == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}

	seniority := domain.Seniority(input.Seniority)
	if !seniority.Valid() {
		return nil, fmt.Errorf("invalid seniority: %s", input.Seniority)
	}

	// Create aggregate (generates a new ID; for now, use a simple counter)
	// In production, this would be auto-incremented from the database
	// For testing, the repository will manage ID generation
	// Use timestamp-based ID to avoid collisions in tests
	id := int64(time.Now().UnixNano() / 1000000) % 1000000
	if id <= 0 {
		id = 1
	}

	name, err := domain.NewFullName(input.FirstName, input.LastName)
	if err != nil {
		return nil, err
	}

	member, err := NewTeamMember(id, name, seniority)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := uc.repo.Save(member); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	// Map to output
	output := &AddMemberOutput{
		ID:        int64(member.ID()),
		FirstName: member.Name().First,
		LastName:  member.Name().Last,
		Seniority: string(member.Seniority()),
		Status:    "Active",
		CreatedAt: member.CreatedAt(),
	}

	return output, nil
}

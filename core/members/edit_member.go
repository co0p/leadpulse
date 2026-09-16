package members

import (
	"fmt"

	"leadpulse/engine/domain"
)

// EditMemberInput is the input struct for EditMemberUseCase.
type EditMemberInput struct {
	MemberID  int64
	FirstName string // Empty string means no change
	LastName  string // Empty string means no change
	Seniority string // Empty string means no change
}

// EditMemberOutput is the output struct for EditMemberUseCase.
type EditMemberOutput struct {
	ID        int64
	FirstName string
	LastName  string
	Seniority string
	Status    string
	CreatedAt interface{} // time.Time; used generically here
}

// EditMemberUseCase handles updating an existing team member.
type EditMemberUseCase struct {
	repo MemberRepository
}

// NewEditMemberUseCase creates a new instance of EditMemberUseCase with an injected repository.
func NewEditMemberUseCase(repo MemberRepository) *EditMemberUseCase {
	return &EditMemberUseCase{repo: repo}
}

// Execute updates a team member with the provided input.
// Empty string fields are treated as "no change" (partial updates supported).
// Returns an error if the member is not found or if validation fails.
func (uc *EditMemberUseCase) Execute(input EditMemberInput) (*EditMemberOutput, error) {
	if input.MemberID <= 0 {
		return nil, fmt.Errorf("member ID must be positive")
	}

	// Retrieve the member
	member, err := uc.repo.FindByID(TeamMemberID(input.MemberID))
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, fmt.Errorf("member not found: id=%d", input.MemberID)
	}

	// Apply updates (empty strings mean no change)
	if input.FirstName != "" || input.LastName != "" {
		// Both firstName and lastName are required to update; if one is empty, use current value
		firstName := input.FirstName
		lastName := input.LastName
		if firstName == "" {
			firstName = member.Name().First
		}
		if lastName == "" {
			lastName = member.Name().Last
		}

		if err := member.UpdateName(firstName, lastName); err != nil {
			return nil, err
		}
	}

	if input.Seniority != "" {
		seniority := domain.Seniority(input.Seniority)
		if !seniority.Valid() {
			return nil, fmt.Errorf("invalid seniority: %s", input.Seniority)
		}
		if err := member.ChangeSeniority(seniority); err != nil {
			return nil, err
		}
	}

	// Persist changes
	if err := uc.repo.Save(member); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	// Map to output
	output := &EditMemberOutput{
		ID:        int64(member.ID()),
		FirstName: member.Name().First,
		LastName:  member.Name().Last,
		Seniority: string(member.Seniority()),
		Status:    "Active",
		CreatedAt: member.CreatedAt(),
	}

	return output, nil
}

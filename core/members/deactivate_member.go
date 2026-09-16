package members

import (
	"fmt"
)

// DeactivateMemberInput is the input struct for DeactivateMemberUseCase.
type DeactivateMemberInput struct {
	MemberID int64
}

// DeactivateMemberOutput is the output struct for DeactivateMemberUseCase.
type DeactivateMemberOutput struct {
	ID            int64
	FirstName     string
	LastName      string
	Seniority     string
	Status        string
	DeactivatedAt interface{} // *time.Time; used generically here
}

// DeactivateMemberUseCase handles deactivation of team members.
type DeactivateMemberUseCase struct {
	repo MemberRepository
}

// NewDeactivateMemberUseCase creates a new instance of DeactivateMemberUseCase with an injected repository.
func NewDeactivateMemberUseCase(repo MemberRepository) *DeactivateMemberUseCase {
	return &DeactivateMemberUseCase{repo: repo}
}

// Execute deactivates a team member by ID.
// Returns an error if the member is not found or already deactivated.
func (uc *DeactivateMemberUseCase) Execute(input DeactivateMemberInput) (*DeactivateMemberOutput, error) {
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

	// Deactivate
	if err := member.Deactivate(); err != nil {
		return nil, err
	}

	// Persist changes
	if err := uc.repo.Save(member); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	// Map to output
	output := &DeactivateMemberOutput{
		ID:            int64(member.ID()),
		FirstName:     member.Name().First,
		LastName:      member.Name().Last,
		Seniority:     string(member.Seniority()),
		Status:        "Inactive",
		DeactivatedAt: member.DeactivatedAt(),
	}

	return output, nil
}

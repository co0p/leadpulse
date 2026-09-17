package members

import "fmt"

// ReactivateMemberInput is the input struct for ReactivateMemberUseCase.
type ReactivateMemberInput struct {
	MemberID int64
}

// ReactivateMemberOutput is the output struct for ReactivateMemberUseCase.
type ReactivateMemberOutput struct {
	ID            int64
	FirstName     string
	LastName      string
	Seniority     string
	Status        string
	DeactivatedAt interface{} // *time.Time; used generically here
}

// ReactivateMemberUseCase handles reactivation of team members.
type ReactivateMemberUseCase struct {
	repo MemberRepository
}

// NewReactivateMemberUseCase creates a new instance of ReactivateMemberUseCase with an injected repository.
func NewReactivateMemberUseCase(repo MemberRepository) *ReactivateMemberUseCase {
	return &ReactivateMemberUseCase{repo: repo}
}

// Execute reactivates a team member by ID.
// Returns an error if the member is not found or if already active.
func (uc *ReactivateMemberUseCase) Execute(input ReactivateMemberInput) (*ReactivateMemberOutput, error) {
	// Validate input
	if input.MemberID <= 0 {
		return nil, fmt.Errorf("team member ID must be positive")
	}

	// Find the member
	member, err := uc.repo.FindByID(TeamMemberID(input.MemberID))
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, fmt.Errorf("team member not found")
	}

	// Reactivate
	if err := member.Reactivate(); err != nil {
		return nil, err
	}

	// Persist
	if err := uc.repo.Update(member); err != nil {
		return nil, err
	}

	// Return output
	output := &ReactivateMemberOutput{
		ID:            int64(member.ID()),
		FirstName:     member.Name().First,
		LastName:      member.Name().Last,
		Seniority:     string(member.Seniority()),
		Status:        "Active",
		DeactivatedAt: member.DeactivatedAt(),
	}

	return output, nil
}

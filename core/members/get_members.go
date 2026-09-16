package members

import (
	"time"
)

// GetMembersOutput is the output struct for GetMembersUseCase.
type GetMembersOutput struct {
	Members []MemberDTO
}

// MemberDTO is a data transfer object for representing a team member.
type MemberDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Seniority string
	Status    string
	CreatedAt time.Time
}

// GetMembersUseCase retrieves all active team members.
type GetMembersUseCase struct {
	repo MemberRepository
}

// NewGetMembersUseCase creates a new instance of GetMembersUseCase with an injected repository.
func NewGetMembersUseCase(repo MemberRepository) *GetMembersUseCase {
	return &GetMembersUseCase{repo: repo}
}

// Execute retrieves all active team members.
// Returns an empty list if no active members exist (not an error).
func (uc *GetMembersUseCase) Execute() (*GetMembersOutput, error) {
	members, err := uc.repo.FindActive()
	if err != nil {
		return nil, err
	}

	// Map to DTOs
	dtos := make([]MemberDTO, len(members))
	for i, member := range members {
		dtos[i] = MemberDTO{
			ID:        int64(member.ID()),
			FirstName: member.Name().First,
			LastName:  member.Name().Last,
			Seniority: string(member.Seniority()),
			Status:    "Active",
			CreatedAt: member.CreatedAt(),
		}
	}

	return &GetMembersOutput{Members: dtos}, nil
}

package members

import (
	"fmt"
	"time"
)

// GetMembersOutput is the output struct for GetMembersUseCase.
type GetMembersOutput struct {
	Members []MemberDTO
}

// GetFilteredMembersInput is the input struct for GetFilteredMembersUseCase.
type GetFilteredMembersInput struct {
	Status string // "active", "deactivated", or "all"
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

// GetFilteredMembersUseCase retrieves team members filtered by status.
type GetFilteredMembersUseCase struct {
	repo MemberRepository
}

// NewGetFilteredMembersUseCase creates a new instance of GetFilteredMembersUseCase with an injected repository.
func NewGetFilteredMembersUseCase(repo MemberRepository) *GetFilteredMembersUseCase {
	return &GetFilteredMembersUseCase{repo: repo}
}

// Execute retrieves team members filtered by status.
// Status can be "active", "deactivated", or "all".
// Returns an error if the status is invalid.
func (uc *GetFilteredMembersUseCase) Execute(input GetFilteredMembersInput) (*GetMembersOutput, error) {
	// Validate status
	if input.Status != "active" && input.Status != "deactivated" && input.Status != "all" {
		return nil, fmt.Errorf("invalid status: %s; must be one of: active, deactivated, all", input.Status)
	}

	var members []*TeamMember
	var err error

	// Fetch based on status
	switch input.Status {
	case "active":
		members, err = uc.repo.FindActive()
	case "deactivated":
		members, err = uc.repo.FindInactive()
	case "all":
		members, err = uc.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	// Map to DTOs
	dtos := make([]MemberDTO, len(members))
	for i, member := range members {
		status := "Active"
		if !member.IsActive() {
			status = "Inactive"
		}
		dtos[i] = MemberDTO{
			ID:        int64(member.ID()),
			FirstName: member.Name().First,
			LastName:  member.Name().Last,
			Seniority: string(member.Seniority()),
			Status:    status,
			CreatedAt: member.CreatedAt(),
		}
	}

	return &GetMembersOutput{Members: dtos}, nil
}

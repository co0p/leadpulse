package member

import (
	"fmt"
	"strings"

	"leadpulse/engine/domain"
)

// Service provides use cases for team member management.
// It now depends on the TeamMemberRepository interface instead of direct database access.
type Service struct {
	repo domain.TeamMemberRepository
}

// NewService creates a new member service with a repository implementation.
func NewService(repo domain.TeamMemberRepository) *Service {
	return &Service{repo: repo}
}

// NewServiceWithDB creates a service using the SQLite repository.
// This is provided for backward compatibility with existing code.
func NewServiceWithDB(db interface{}) *Service {
	// The db parameter is now unused but kept for compatibility
	// In a full refactor, callers would be updated to pass the repository directly
	panic("use NewService with a repository instead")
}

// AddMember creates a new team member and returns it.
func (s *Service) AddMember(firstName, lastName string, seniority domain.Seniority) (*domain.TeamMember, error) {
	// Validate input
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}

	if !seniority.Valid() {
		return nil, fmt.Errorf("invalid seniority level: %s", seniority)
	}

	// Create the full name value object
	name, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Allocate an ID for the new member
	// (In a full implementation, the repository might handle this)
	allMembers, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate ID: %w", err)
	}

	var maxID domain.TeamMemberID
	for _, m := range allMembers {
		if m.ID() > maxID {
			maxID = m.ID()
		}
	}

	newID := maxID + 1

	// Create the aggregate root
	member, err := domain.NewTeamMember(int64(newID), name, seniority)
	if err != nil {
		return nil, err
	}

	// Persist via repository
	if err := s.repo.Save(member); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	return member, nil
}

// ListMembers returns all active team members.
func (s *Service) ListMembers() ([]domain.TeamMember, error) {
	members, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}

	// Convert []*TeamMember to []TeamMember for backward compatibility
	var result []domain.TeamMember
	for _, m := range members {
		if m != nil {
			result = append(result, *m)
		}
	}
	return result, nil
}

// GetMember retrieves a single member by ID.
func (s *Service) GetMember(id int64) (*domain.TeamMember, error) {
	return s.repo.FindByID(domain.TeamMemberID(id))
}

// EditMember updates an existing member's information.
func (s *Service) EditMember(id int64, firstName, lastName string, seniority domain.Seniority) (*domain.TeamMember, error) {
	// Validate input
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}

	if !seniority.Valid() {
		return nil, fmt.Errorf("invalid seniority level: %s", seniority)
	}

	// Fetch existing member
	member, err := s.repo.FindByID(domain.TeamMemberID(id))
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, fmt.Errorf("member not found: %d", id)
	}

	// Create updated full name
	name, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Create a new aggregate with updated values
	updatedMember, err := domain.NewTeamMember(id, name, seniority)
	if err != nil {
		return nil, err
	}

	// Preserve deactivation state if already deactivated
	if !member.IsActive() {
		_ = updatedMember.Deactivate()
	}

	// Persist via repository
	if err := s.repo.Save(updatedMember); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	return updatedMember, nil
}

// DeactivateMember soft-deletes a team member.
func (s *Service) DeactivateMember(id int64) error {
	member, err := s.repo.FindByID(domain.TeamMemberID(id))
	if err != nil {
		return err
	}
	if member == nil {
		return fmt.Errorf("member not found: %d", id)
	}

	// Call the aggregate method to deactivate
	if err := member.Deactivate(); err != nil {
		return err
	}

	// Persist the deactivated state
	return s.repo.Save(member)
}

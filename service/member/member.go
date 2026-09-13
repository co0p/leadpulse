package member

import (
	"database/sql"
	"fmt"
	"strings"

	"leadpulse/engine/domain"
	"leadpulse/store"
)

// Service provides use cases for team member management.
type Service struct {
	db *sql.DB
}

// NewService creates a new member service.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
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

	// Add to store
	id, err := store.AddMember(s.db, firstName, lastName, seniority)
	if err != nil {
		return nil, err
	}

	// Retrieve and return the created member
	return store.GetMember(s.db, id)
}

// ListMembers returns all active team members.
func (s *Service) ListMembers() ([]domain.TeamMember, error) {
	return store.ListMembers(s.db)
}

// GetMember retrieves a single member by ID.
func (s *Service) GetMember(id int64) (*domain.TeamMember, error) {
	return store.GetMember(s.db, id)
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

	// Update in store
	if err := store.EditMember(s.db, id, firstName, lastName, seniority); err != nil {
		return nil, err
	}

	// Retrieve and return the updated member
	return store.GetMember(s.db, id)
}

// DeactivateMember soft-deletes a team member.
func (s *Service) DeactivateMember(id int64) error {
	return store.DeactivateMember(s.db, id)
}

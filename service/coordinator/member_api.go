package coordinator

import (
	"leadpulse/engine/domain"
	membersvc "leadpulse/service/member"
)

// MemberAPICoordinator owns all business logic for the Members API workflow.
// It manages member list state and member CRUD operations for the HTTP API.
// This coordinator is API-specific and supports both internal int64 IDs and external UUID representation.
// Coordinators have zero HTTP dependencies and are fully unit testable.
type MemberAPICoordinator struct {
	memberService *membersvc.Service
	members       []domain.TeamMember
}

// NewMemberAPICoordinator creates a new coordinator for the Members API workflow.
func NewMemberAPICoordinator(memberService *membersvc.Service) *MemberAPICoordinator {
	return &MemberAPICoordinator{
		memberService: memberService,
		members:       []domain.TeamMember{},
	}
}

// Load initializes the coordinator by loading all active team members.
func (c *MemberAPICoordinator) Load() error {
	members, err := c.memberService.ListMembers()
	if err != nil {
		return err
	}
	c.members = members
	return nil
}

// Validate checks if the current state is valid.
func (c *MemberAPICoordinator) Validate() error {
	// API validation: no rules at this layer.
	return nil
}

// GetMembers returns the current list of active members.
func (c *MemberAPICoordinator) GetMembers() []domain.TeamMember {
	return c.members
}

// AddMember creates a new team member and persists it.
// After success, the member list is reloaded to reflect the new member.
//
// Returns a ValidationError if:
// - firstName or lastName is empty
// - seniority is invalid
// - Member creation fails (wrapped with user-friendly message)
func (c *MemberAPICoordinator) AddMember(firstName, lastName string, seniority domain.Seniority) error {
	if firstName == "" || lastName == "" {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"First name and last name are required",
			nil,
		)
	}

	if !seniority.Valid() {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"Invalid seniority level",
			nil,
		)
	}

	_, err := c.memberService.AddMember(firstName, lastName, seniority)
	if err != nil {
		// Wrap service layer errors
		return WrapError(
			ValidationErrorKindDatabaseFailure,
			"Failed to add member. Please try again.",
			err,
		)
	}

	// Reload members list
	return c.Load()
}

// EditMember updates an existing team member's information.
// After success, the member list is reloaded.
//
// Only non-empty fields are updated. If only seniority is provided, first and last names 
// are retrieved from the existing member.
//
// At least one field must be non-empty for the update to proceed.
//
// Returns a ValidationError if:
// - memberID is invalid or member not found
// - firstName and lastName are both empty (and seniority is also empty)
// - seniority is invalid (if provided)
// - Member update fails (wrapped with user-friendly message)
func (c *MemberAPICoordinator) EditMember(memberID int64, firstName, lastName string, seniority domain.Seniority) error {
	// Validate that at least one field is being updated
	if firstName == "" && lastName == "" && seniority == "" {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"At least one field must be provided to update",
			nil,
		)
	}

	// Get existing member
	existingMember := c.GetMemberByID(memberID)
	if existingMember == nil {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"Member not found",
			nil,
		)
	}

	// Use provided values, or fall back to existing values
	finalFirstName := firstName
	finalLastName := lastName
	finalSeniority := seniority

	if finalFirstName == "" {
		finalFirstName = existingMember.Name().First
	}
	if finalLastName == "" {
		finalLastName = existingMember.Name().Last
	}
	if finalSeniority == "" {
		finalSeniority = existingMember.Seniority()
	}

	// Now validate the final values
	if finalFirstName == "" || finalLastName == "" {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"First name and last name are required",
			nil,
		)
	}

	if !finalSeniority.Valid() {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"Invalid seniority level",
			nil,
		)
	}

	_, err := c.memberService.EditMember(memberID, finalFirstName, finalLastName, finalSeniority)
	if err != nil {
		// Wrap service layer errors
		return WrapError(
			ValidationErrorKindDatabaseFailure,
			"Failed to edit member. Please try again.",
			err,
		)
	}

	// Reload members list
	return c.Load()
}

// DeactivateMember marks a team member as inactive.
// After success, the member list is reloaded.
//
// Returns a ValidationError if:
// - memberID is invalid or member not found
// - member deactivation fails (wrapped with user-friendly message)
func (c *MemberAPICoordinator) DeactivateMember(memberID int64) error {
	// Verify member exists before attempting deactivation
	existingMember := c.GetMemberByID(memberID)
	if existingMember == nil {
		return NewValidationError(
			ValidationErrorKindInvalidField,
			"Member not found",
			nil,
		)
	}

	err := c.memberService.DeactivateMember(memberID)
	if err != nil {
		// Wrap service layer errors
		return WrapError(
			ValidationErrorKindDatabaseFailure,
			"Failed to deactivate member. Please try again.",
			err,
		)
	}

	// Reload members list
	return c.Load()
}

// GetMemberByID finds a member in the current list by ID.
func (c *MemberAPICoordinator) GetMemberByID(memberID int64) *domain.TeamMember {
	for i := range c.members {
		if int64(c.members[i].ID()) == memberID {
			return &c.members[i]
		}
	}
	return nil
}

// GetMemberCount returns the number of members in the current list.
func (c *MemberAPICoordinator) GetMemberCount() int {
	return len(c.members)
}

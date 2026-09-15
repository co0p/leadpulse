package controllers

import (
	"leadpulse/engine/domain"
	membersvc "leadpulse/service/member"
)

// SettingsController owns all business logic for the Settings screen.
// It manages member list state and member CRUD operations.
// Controllers have zero Fyne dependencies and are fully unit testable.
type SettingsController struct {
	memberService *membersvc.Service
	members       []domain.TeamMember
}

// NewSettingsController creates a new controller for the Settings screen.
func NewSettingsController(memberService *membersvc.Service) *SettingsController {
	return &SettingsController{
		memberService: memberService,
		members:       []domain.TeamMember{},
	}
}

// Load initializes the controller by loading all active team members.
func (c *SettingsController) Load() error {
	members, err := c.memberService.ListMembers()
	if err != nil {
		return err
	}
	c.members = members
	return nil
}

// Validate checks if the current state is valid.
func (c *SettingsController) Validate() error {
	// Settings validation: no rules at this layer.
	return nil
}

// GetMembers returns the current list of members.
func (c *SettingsController) GetMembers() []domain.TeamMember {
	return c.members
}

// AddMember creates a new team member and persists it.
// After success, the member list is reloaded to reflect the new member.
//
// Returns a ValidationError if:
// - firstName or lastName is empty
// - Member creation fails (wrapped with user-friendly message)
func (c *SettingsController) AddMember(firstName, lastName string, seniority domain.Seniority) error {
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
// Returns a ValidationError if:
// - firstName or lastName is empty
// - Member update fails (wrapped with user-friendly message)
func (c *SettingsController) EditMember(memberID int64, firstName, lastName string, seniority domain.Seniority) error {
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

	_, err := c.memberService.EditMember(memberID, firstName, lastName, seniority)
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
// Returns a ValidationError if member deactivation fails (wrapped with user-friendly message).
func (c *SettingsController) DeactivateMember(memberID int64) error {
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
func (c *SettingsController) GetMemberByID(memberID int64) *domain.TeamMember {
	for i := range c.members {
		if int64(c.members[i].ID()) == memberID {
			return &c.members[i]
		}
	}
	return nil
}

// GetMemberCount returns the number of members in the current list.
func (c *SettingsController) GetMemberCount() int {
	return len(c.members)
}

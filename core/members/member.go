package members

import (
	"fmt"
	"time"

	"leadpulse/engine/domain"
)

// TeamMemberID is a value object for team member identity.
type TeamMemberID int64

// TeamMember is an aggregate root representing a team member with all their properties.
// All fields are private; access is controlled via methods.
type TeamMember struct {
	id            TeamMemberID
	name          domain.FullName
	seniority     domain.Seniority
	createdAt     time.Time
	deactivatedAt *time.Time
}

// NewTeamMember creates a new team member aggregate root.
// All constraints are enforced at construction.
func NewTeamMember(id int64, name domain.FullName, seniority domain.Seniority) (*TeamMember, error) {
	if id <= 0 {
		return nil, fmt.Errorf("team member ID must be positive")
	}

	if !seniority.Valid() {
		return nil, fmt.Errorf("invalid seniority: %s", seniority)
	}

	return &TeamMember{
		id:        TeamMemberID(id),
		name:      name,
		seniority: seniority,
		createdAt: time.Now(),
	}, nil
}

// ID returns the team member's immutable ID.
func (tm *TeamMember) ID() TeamMemberID {
	return tm.id
}

// Name returns the team member's full name.
func (tm *TeamMember) Name() domain.FullName {
	return tm.name
}

// Seniority returns the team member's seniority level.
func (tm *TeamMember) Seniority() domain.Seniority {
	return tm.seniority
}

// CreatedAt returns when the member was created.
func (tm *TeamMember) CreatedAt() time.Time {
	return tm.createdAt
}

// IsActive returns true if the member has not been deactivated.
func (tm *TeamMember) IsActive() bool {
	return tm.deactivatedAt == nil
}

// DeactivatedAt returns the deactivation timestamp, or nil if active.
func (tm *TeamMember) DeactivatedAt() *time.Time {
	return tm.deactivatedAt
}

// Deactivate marks the team member as inactive.
// Once deactivated, this cannot be undone (enforces business rule at aggregate).
func (tm *TeamMember) Deactivate() error {
	if !tm.IsActive() {
		return fmt.Errorf("team member is already deactivated")
	}
	now := time.Now()
	tm.deactivatedAt = &now
	return nil
}

// ChangeSeniority updates the seniority level.
// Encapsulation ensures only valid transitions are possible.
func (tm *TeamMember) ChangeSeniority(newSeniority domain.Seniority) error {
	if !newSeniority.Valid() {
		return fmt.Errorf("invalid seniority: %s", newSeniority)
	}
	tm.seniority = newSeniority
	return nil
}

// UpdateName updates both the first and last name of the team member.
// Both firstName and lastName must be non-empty.
func (tm *TeamMember) UpdateName(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return fmt.Errorf("first and last name are required")
	}
	newName, err := domain.NewFullName(firstName, lastName)
	if err != nil {
		return err
	}
	tm.name = newName
	return nil
}

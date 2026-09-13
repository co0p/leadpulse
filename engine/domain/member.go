package domain

import "time"

// TeamMember represents a team member in the scorecard system.
// It is an immutable value object. All fields are exported for persistence,
// but the type is meant to be created only via constructor functions or store layer.
type TeamMember struct {
	ID           int64
	FirstName    string
	LastName     string
	Seniority    Seniority
	CreatedAt    time.Time
	DeactivatedAt *time.Time
}

// FullName returns the concatenated first and last name.
func (tm TeamMember) FullName() string {
	return tm.FirstName + " " + tm.LastName
}

// IsActive returns true if the member has not been deactivated.
func (tm TeamMember) IsActive() bool {
	return tm.DeactivatedAt == nil
}

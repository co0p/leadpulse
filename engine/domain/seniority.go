package domain

// Seniority represents the seniority level of a team member.
type Seniority string

const (
	SeniorityJunior    Seniority = "Junior"
	SeniorityMid       Seniority = "Mid"
	SenioritySenior    Seniority = "Senior"
	SeniorityPrincipal Seniority = "Principal"
)

// Valid returns true if s is a valid Seniority value.
func (s Seniority) Valid() bool {
	switch s {
	case SeniorityJunior, SeniorityMid, SenioritySenior, SeniorityPrincipal:
		return true
	default:
		return false
	}
}

// AllSeniorities returns a slice of all valid Seniority values.
func AllSeniorities() []Seniority {
	return []Seniority{
		SeniorityJunior,
		SeniorityMid,
		SenioritySenior,
		SeniorityPrincipal,
	}
}

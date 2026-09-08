package performances

import "fmt"

// EvidenceCategory identifies the type of evidence backing a performance entry.
type EvidenceCategory string

const (
	EvidenceCategoryProjectOutcome     EvidenceCategory = "ProjectOutcome"
	EvidenceCategoryPeerFeedback       EvidenceCategory = "PeerFeedback"
	EvidenceCategoryBlockerResolution  EvidenceCategory = "BlockerResolution"
	EvidenceCategoryTechnicalAchievement EvidenceCategory = "TechnicalAchievement"
	EvidenceCategoryCollaboration      EvidenceCategory = "Collaboration"
	EvidenceCategoryOther              EvidenceCategory = "Other"
)

// EvidenceSource identifies where a piece of evidence originated.
type EvidenceSource string

const (
	EvidenceSourceManualEntry       EvidenceSource = "ManualEntry"
	EvidenceSourceGoogleDocImport   EvidenceSource = "GoogleDocImport"
	EvidenceSourceProjectLeadFeedback EvidenceSource = "ProjectLeadFeedback"
)

// EvidenceItem is a single dated piece of evidence linked to an entry.
type EvidenceItem struct {
	Category    EvidenceCategory `json:"category"`
	Source      EvidenceSource   `json:"source"`
	Description string           `json:"description"`
}

// Entry is a monthly performance snapshot for a single employee.
type Entry struct {
	ID            int64          `json:"id"`
	EmployeeID    int64          `json:"employeeId"`
	Year          int            `json:"year"`
	Month         int            `json:"month"`
	Morale        float64        `json:"morale"`
	Execution     float64        `json:"execution"`
	Impact        float64        `json:"impact"`
	Growth        float64        `json:"growth"`
	Culture       float64        `json:"culture"`
	Projects      []string       `json:"projects"`
	EvidenceItems []EvidenceItem `json:"evidenceItems"`
}

// Period returns the calendar period string for the entry.
func (e Entry) Period() string {
	return fmt.Sprintf("%04d-%02d", e.Year, e.Month)
}

// Multiplier returns the evidence-based 0-5 multiplier for the month.
// Five requires at least three evidence items; otherwise it is zero.
func (e Entry) Multiplier() float64 {
	if len(e.EvidenceItems) >= 3 {
		return 5
	}
	return 0
}

// HighCapacity returns true when the employee worked on at least three different projects.
func (e Entry) HighCapacity() bool {
	return len(e.Projects) >= 3
}

package overview

import "leadpulse/internal/employees"

// EmployeeHealth captures rolling three-month results for one employee.
type EmployeeHealth struct {
	EmployeeID    int64   `json:"employeeId"`
	Name          string  `json:"name"`
	Morale        float64 `json:"morale"`
	Execution     float64 `json:"execution"`
	Impact        float64 `json:"impact"`
	Multiplier    float64 `json:"multiplier"`
	Growth        float64 `json:"growth"`
	Culture       float64 `json:"culture"`
	ProjectCount  int     `json:"projectCount"`
	EvidenceCount int     `json:"evidenceCount"`
}

// AnomalyType identifies the kind of attention signal.
type AnomalyType string

const (
	AnomalyMoraleDrop   AnomalyType = "MoraleDrop"
	AnomalyHighCapacity AnomalyType = "HighCapacity"
	AnomalyBurnoutRisk  AnomalyType = "BurnoutRisk"
)

// Anomaly is an actionable signal for a specific employee.
type Anomaly struct {
	EmployeeID int64       `json:"employeeId"`
	Name       string      `json:"name"`
	Type       AnomalyType `json:"type"`
	Message    string      `json:"message"`
}

// EvidenceGap describes a metric whose latest rolling window lacks evidence.
type EvidenceGap struct {
	EmployeeID int64  `json:"employeeId"`
	Name       string `json:"name"`
	Metric     string `json:"metric"`
	Message    string `json:"message"`
}

// Overview is the team pulse read model shown on startup.
type Overview struct {
	Headcount         int              `json:"headcount"`
	RollingTeamMorale float64          `json:"rollingTeamMorale"`
	HighCapacityCount int              `json:"highCapacityCount"`
	TeamMultiplier    float64          `json:"teamMultiplier"`
	HealthDimensions  HealthDimensions `json:"healthDimensions"`
	EmployeeHealth    []EmployeeHealth `json:"employeeHealth"`
	Anomalies         []Anomaly        `json:"anomalies"`
	EvidenceGaps      []EvidenceGap    `json:"evidenceGaps"`
}

// HealthDimensions aggregates rolling three-month results across the team.
type HealthDimensions struct {
	Execution  float64 `json:"execution"`
	Impact     float64 `json:"impact"`
	Multiplier float64 `json:"multiplier"`
	Growth     float64 `json:"growth"`
	Culture    float64 `json:"culture"`
}

// Input dependencies for building an overview.
type Input struct {
	Employees    []employees.Employee
	Entries      []EntryAdapter
	CurrentYear  int
	CurrentMonth int
}

// EntryAdapter exposes the fields the overview calculator needs from any entry type.
type EntryAdapter struct {
	EmployeeID    int64
	Year          int
	Month         int
	Morale        float64
	Execution     float64
	Impact        float64
	Growth        float64
	Culture       float64
	Projects      []string
	EvidenceCount int
}

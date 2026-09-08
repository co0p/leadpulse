package overview

import (
	"math"
	"testing"

	"leadpulse/internal/employees"
)

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func TestCalculate_NoEmployees(t *testing.T) {
	result := Calculate(Input{CurrentYear: 2026, CurrentMonth: 9})
	if result.Headcount != 0 {
		t.Fatalf("expected headcount 0, got %d", result.Headcount)
	}
}

func TestCalculate_SingleEmployeeCompleteHistory(t *testing.T) {
	emp := employees.Employee{ID: 1, FirstName: "Grace", SecondName: "Hopper"}
	entries := []EntryAdapter{
		{EmployeeID: 1, Year: 2026, Month: 7, Morale: 4, Execution: 4, Impact: 4, Growth: 4, Culture: 4, Projects: []string{"Alpha"}, EvidenceCount: 3},
		{EmployeeID: 1, Year: 2026, Month: 8, Morale: 4, Execution: 4, Impact: 4, Growth: 4, Culture: 4, Projects: []string{"Alpha"}, EvidenceCount: 3},
		{EmployeeID: 1, Year: 2026, Month: 9, Morale: 5, Execution: 5, Impact: 5, Growth: 5, Culture: 5, Projects: []string{"Alpha"}, EvidenceCount: 3},
	}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      entries,
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	if result.Headcount != 1 {
		t.Fatalf("expected headcount 1, got %d", result.Headcount)
	}
	if !floatEquals(result.RollingTeamMorale, 4.3333) {
		t.Fatalf("expected team morale ~4.33, got %f", result.RollingTeamMorale)
	}
	if !floatEquals(result.TeamMultiplier, 5.0) {
		t.Fatalf("expected team multiplier 5, got %f", result.TeamMultiplier)
	}
	if result.HighCapacityCount != 0 {
		t.Fatalf("expected no high capacity, got %d", result.HighCapacityCount)
	}
}

func TestCalculate_IncompleteHistoryUsesZero(t *testing.T) {
	emp := employees.Employee{ID: 2, FirstName: "Ada", SecondName: "Lovelace"}
	entries := []EntryAdapter{
		{EmployeeID: 2, Year: 2026, Month: 9, Morale: 3, Execution: 3, Impact: 3, Growth: 3, Culture: 3, Projects: []string{"Alpha"}, EvidenceCount: 1},
	}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      entries,
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	if !floatEquals(result.RollingTeamMorale, 1.0) {
		t.Fatalf("expected team morale 1.0 with zero-fill, got %f", result.RollingTeamMorale)
	}
}

func TestCalculate_HighCapacityAndMoraleDrop(t *testing.T) {
	emp := employees.Employee{ID: 3, FirstName: "Alan", SecondName: "Turing"}
	entries := []EntryAdapter{
		{EmployeeID: 3, Year: 2026, Month: 7, Morale: 5, Execution: 4, Impact: 4, Growth: 4, Culture: 4, Projects: []string{"Alpha"}, EvidenceCount: 2},
		{EmployeeID: 3, Year: 2026, Month: 8, Morale: 4, Execution: 4, Impact: 4, Growth: 4, Culture: 4, Projects: []string{"Alpha", "Beta", "Gamma"}, EvidenceCount: 2},
		{EmployeeID: 3, Year: 2026, Month: 9, Morale: 3, Execution: 3, Impact: 3, Growth: 3, Culture: 3, Projects: []string{"Alpha", "Beta", "Gamma"}, EvidenceCount: 2},
	}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      entries,
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	if result.HighCapacityCount != 1 {
		t.Fatalf("expected high capacity count 1, got %d", result.HighCapacityCount)
	}
	if len(result.Anomalies) == 0 {
		t.Fatal("expected anomalies")
	}
	foundMoraleDrop := false
	foundHighCapacity := false
	for _, a := range result.Anomalies {
		if a.Type == AnomalyMoraleDrop {
			foundMoraleDrop = true
		}
		if a.Type == AnomalyHighCapacity {
			foundHighCapacity = true
		}
	}
	if !foundMoraleDrop {
		t.Fatalf("expected morale drop anomaly, got %+v", result.Anomalies)
	}
	if !foundHighCapacity {
		t.Fatalf("expected high capacity anomaly, got %+v", result.Anomalies)
	}
}

func TestCalculate_EvidenceGap(t *testing.T) {
	emp := employees.Employee{ID: 4, FirstName: "Linus", SecondName: "Torvalds"}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      []EntryAdapter{},
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	if len(result.EvidenceGaps) == 0 {
		t.Fatal("expected evidence gap for missing entries")
	}
}

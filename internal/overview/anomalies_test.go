package overview

import (
	"testing"

	"leadpulse/internal/employees"
)

func TestDetectBurnoutRisk_HighCapacityAndLowMorale(t *testing.T) {
	emp := employees.Employee{ID: 1, FirstName: "Alan", SecondName: "Turing"}
	entries := []EntryAdapter{
		{EmployeeID: 1, Year: 2026, Month: 8, Morale: 3, Projects: []string{"Alpha", "Beta", "Gamma"}, EvidenceCount: 2},
		{EmployeeID: 1, Year: 2026, Month: 9, Morale: 3, Projects: []string{"Alpha", "Beta", "Gamma"}, EvidenceCount: 2},
	}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      entries,
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	found := false
	for _, a := range result.Anomalies {
		if a.Type == AnomalyBurnoutRisk {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected burnout risk anomaly, got %+v", result.Anomalies)
	}
}

func TestDetectBurnoutRisk_NoHighCapacity(t *testing.T) {
	emp := employees.Employee{ID: 2, FirstName: "Grace", SecondName: "Hopper"}
	entries := []EntryAdapter{
		{EmployeeID: 2, Year: 2026, Month: 8, Morale: 3, Projects: []string{"Alpha"}, EvidenceCount: 3},
		{EmployeeID: 2, Year: 2026, Month: 9, Morale: 3, Projects: []string{"Alpha"}, EvidenceCount: 3},
	}
	result := Calculate(Input{
		Employees:    []employees.Employee{emp},
		Entries:      entries,
		CurrentYear:  2026,
		CurrentMonth: 9,
	})

	for _, a := range result.Anomalies {
		if a.Type == AnomalyBurnoutRisk {
			t.Fatalf("unexpected burnout risk anomaly: %+v", a)
		}
	}
}

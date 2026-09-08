package app

import (
	"testing"

	"leadpulse/internal/employees"
	"leadpulse/internal/performances"
)

func TestTeamPulseOverview(t *testing.T) {
	empRepo := employees.NewFakeRepository()
	perfRepo := performances.NewFakeRepository()
	app := New(empRepo, perfRepo)

	emp, err := empRepo.Add(employees.Input{FirstName: "Grace", SecondName: "Hopper", Seniority: employees.SeniorityPrincipal, StartDate: "2018-03-01"})
	if err != nil {
		t.Fatalf("add employee: %v", err)
	}
	_, err = perfRepo.Add(performances.Entry{
		EmployeeID: emp.ID,
		Year:       2026,
		Month:      9,
		Morale:     5,
		Execution:  5,
		Impact:     5,
		Growth:     5,
		Culture:    5,
		Projects:   []string{"Alpha"},
		EvidenceItems: []performances.EvidenceItem{
			{Category: performances.EvidenceCategoryProjectOutcome, Source: performances.EvidenceSourceManualEntry, Description: "A"},
			{Category: performances.EvidenceCategoryPeerFeedback, Source: performances.EvidenceSourceManualEntry, Description: "B"},
			{Category: performances.EvidenceCategoryTechnicalAchievement, Source: performances.EvidenceSourceManualEntry, Description: "C"},
		},
	})
	if err != nil {
		t.Fatalf("add performance entry: %v", err)
	}

	overview, err := app.TeamPulseOverview()
	if err != nil {
		t.Fatalf("get overview: %v", err)
	}
	if overview.Headcount != 1 {
		t.Fatalf("expected headcount 1, got %d", overview.Headcount)
	}
	if overview.TeamMultiplier != 5.0/3.0 {
		t.Fatalf("expected team multiplier ~1.67 with zero-fill, got %f", overview.TeamMultiplier)
	}
}

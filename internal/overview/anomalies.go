package overview

import (
	"fmt"
	"sort"

	"leadpulse/internal/employees"
)

// detectAnomalies returns actionable signals for the team.
func detectAnomalies(input Input, health []EmployeeHealth, entriesByEmployee map[int64][]EntryAdapter) []Anomaly {
	byID := make(map[int64]employees.Employee)
	for _, emp := range input.Employees {
		byID[emp.ID] = emp
	}

	anomalies := make([]Anomaly, 0)
	for _, h := range health {
		_ = byID[h.EmployeeID]
		if h.ProjectCount >= 3 {
			anomalies = append(anomalies, Anomaly{
				EmployeeID: h.EmployeeID,
				Name:       h.Name,
				Type:       AnomalyHighCapacity,
				Message:    fmt.Sprintf("%s is assigned to %d projects", h.Name, h.ProjectCount),
			})
		}

		moraleDrop := calculateMoraleDrop(entriesByEmployee[h.EmployeeID], input.CurrentYear, input.CurrentMonth)
		if moraleDrop >= 1.0 {
			anomalies = append(anomalies, Anomaly{
				EmployeeID: h.EmployeeID,
				Name:       h.Name,
				Type:       AnomalyMoraleDrop,
				Message:    fmt.Sprintf("%s's rolling morale dropped by %.1f", h.Name, moraleDrop),
			})
		}

		if detectBurnoutRisk(entriesByEmployee[h.EmployeeID], input.CurrentYear, input.CurrentMonth) {
			anomalies = append(anomalies, Anomaly{
				EmployeeID: h.EmployeeID,
				Name:       h.Name,
				Type:       AnomalyBurnoutRisk,
				Message:    fmt.Sprintf("%s shows burnout risk (low morale and overloaded)", h.Name),
			})
		}
	}
	return anomalies
}

// calculateMoraleDrop compares the most recent monthly morale value against
// the monthly morale value of the immediately preceding month.
func calculateMoraleDrop(entries []EntryAdapter, currentYear, currentMonth int) float64 {
	if len(entries) == 0 {
		return 0
	}
	sorted := make([]EntryAdapter, len(entries))
	copy(sorted, entries)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Year != sorted[j].Year {
			return sorted[i].Year < sorted[j].Year
		}
		return sorted[i].Month < sorted[j].Month
	})

	latest := sorted[len(sorted)-1]
	if len(sorted) < 2 {
		return 0
	}
	previous := sorted[len(sorted)-2]
	drop := previous.Morale - latest.Morale
	if drop < 0 {
		return 0
	}
	return drop
}

// detectBurnoutRisk returns true when any 60-day window (two most recent months)
// has high capacity and rolling morale <= 3.
func detectBurnoutRisk(entries []EntryAdapter, currentYear, currentMonth int) bool {
	if len(entries) == 0 {
		return false
	}
	sorted := make([]EntryAdapter, len(entries))
	copy(sorted, entries)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Year != sorted[j].Year {
			return sorted[i].Year < sorted[j].Year
		}
		return sorted[i].Month < sorted[j].Month
	})

	latest := sorted[len(sorted)-1]
	if !highCapacityEntry(latest) {
		return false
	}
	if len(sorted) == 1 {
		return latest.Morale <= 3
	}
	previous := sorted[len(sorted)-2]
	window := []EntryAdapter{previous, latest}
	var moraleSum float64
	for _, e := range window {
		moraleSum += e.Morale
	}
	rolling := moraleSum / float64(len(window))
	return rolling <= 3 && (highCapacityEntry(previous) || highCapacityEntry(latest))
}

func highCapacityEntry(entry EntryAdapter) bool {
	return len(uniqueProjects(entry.Projects)) >= 3
}

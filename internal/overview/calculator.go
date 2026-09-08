package overview

import (
	"fmt"
	"sort"

	"leadpulse/internal/employees"
)

// Calculate produces a team overview from employees and their entries.
// Missing months in the rolling three-month window are treated as zero.
func Calculate(input Input) Overview {
	if len(input.Employees) == 0 {
		return Overview{
			Headcount:        0,
			HealthDimensions: HealthDimensions{},
			EmployeeHealth:   make([]EmployeeHealth, 0),
			Anomalies:        make([]Anomaly, 0),
			EvidenceGaps:     make([]EvidenceGap, 0),
		}
	}

	entriesByEmployee := groupEntriesByEmployee(input.Entries)
	employeeHealth := make([]EmployeeHealth, 0, len(input.Employees))
	var teamMoraleSum, teamMultiplierSum, teamExecutionSum, teamImpactSum, teamGrowthSum, teamCultureSum float64
	highCapacityCount := 0

	for _, emp := range input.Employees {
		latest := latestRollingEntries(entriesByEmployee[emp.ID], input.CurrentYear, input.CurrentMonth)
		health := calculateEmployeeHealth(emp, latest)
		employeeHealth = append(employeeHealth, health)

		teamMoraleSum += health.Morale
		teamMultiplierSum += health.Multiplier
		teamExecutionSum += health.Execution
		teamImpactSum += health.Impact
		teamGrowthSum += health.Growth
		teamCultureSum += health.Culture

		if health.ProjectCount >= 3 {
			highCapacityCount++
		}
	}

	n := float64(len(input.Employees))
	overview := Overview{
		Headcount:         len(input.Employees),
		RollingTeamMorale: teamMoraleSum / n,
		HighCapacityCount: highCapacityCount,
		TeamMultiplier:    teamMultiplierSum / n,
		HealthDimensions: HealthDimensions{
			Execution:  teamExecutionSum / n,
			Impact:     teamImpactSum / n,
			Multiplier: teamMultiplierSum / n,
			Growth:     teamGrowthSum / n,
			Culture:    teamCultureSum / n,
		},
		EmployeeHealth: employeeHealth,
		Anomalies:      make([]Anomaly, 0),
		EvidenceGaps:   make([]EvidenceGap, 0),
	}

	overview.Anomalies = detectAnomalies(input, employeeHealth, entriesByEmployee)
	overview.EvidenceGaps = detectEvidenceGaps(input.Employees, employeeHealth, entriesByEmployee, input.CurrentYear, input.CurrentMonth)

	return overview
}

func groupEntriesByEmployee(entries []EntryAdapter) map[int64][]EntryAdapter {
	grouped := make(map[int64][]EntryAdapter)
	for _, entry := range entries {
		grouped[entry.EmployeeID] = append(grouped[entry.EmployeeID], entry)
	}
	for id := range grouped {
		sort.SliceStable(grouped[id], func(i, j int) bool {
			if grouped[id][i].Year != grouped[id][j].Year {
				return grouped[id][i].Year < grouped[id][j].Year
			}
			return grouped[id][i].Month < grouped[id][j].Month
		})
	}
	return grouped
}

// latestRollingEntries returns the three most recent months ending at currentYear/currentMonth,
// back-filling missing months with zero-valued entries.
func latestRollingEntries(entries []EntryAdapter, currentYear, currentMonth int) []EntryAdapter {
	end := toMonthIndex(currentYear, currentMonth)
	start := end - 2
	if start < 0 {
		start = 0
	}

	byIndex := make(map[int]EntryAdapter)
	for _, entry := range entries {
		idx := toMonthIndex(entry.Year, entry.Month)
		byIndex[idx] = entry
	}

	result := make([]EntryAdapter, 0, 3)
	for idx := start; idx <= end; idx++ {
		if entry, ok := byIndex[idx]; ok {
			result = append(result, entry)
		} else {
			result = append(result, zeroEntry(idx))
		}
	}
	return result
}

func toMonthIndex(year, month int) int {
	return year*12 + (month - 1)
}

func fromMonthIndex(idx int) (year, month int) {
	year = idx / 12
	month = (idx % 12) + 1
	return
}

func zeroEntry(idx int) EntryAdapter {
	year, month := fromMonthIndex(idx)
	return EntryAdapter{Year: year, Month: month}
}

func calculateEmployeeHealth(emp employees.Employee, latest []EntryAdapter) EmployeeHealth {
	var morale, execution, impact, multiplier, growth, culture float64
	var projectCount, evidenceCount int
	if len(latest) > 0 {
		last := latest[len(latest)-1]
		projectCount = len(uniqueProjects(last.Projects))
		evidenceCount = last.EvidenceCount
	}
	for _, entry := range latest {
		morale += entry.Morale
		execution += entry.Execution
		impact += entry.Impact
		multiplier += multiplierForEntry(entry)
		growth += entry.Growth
		culture += entry.Culture
	}
	n := float64(len(latest))
	return EmployeeHealth{
		EmployeeID:    emp.ID,
		Name:          fmt.Sprintf("%s %s", emp.FirstName, emp.SecondName),
		Morale:        morale / n,
		Execution:     execution / n,
		Impact:        impact / n,
		Multiplier:    multiplier / n,
		Growth:        growth / n,
		Culture:       culture / n,
		ProjectCount:  projectCount,
		EvidenceCount: evidenceCount,
	}
}

func multiplierForEntry(entry EntryAdapter) float64 {
	if entry.EvidenceCount >= 3 {
		return 5
	}
	return 0
}

func uniqueProjects(projects []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(projects))
	for _, p := range projects {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			result = append(result, p)
		}
	}
	return result
}

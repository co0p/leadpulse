package overview

import (
	"fmt"

	"leadpulse/internal/employees"
)

// detectEvidenceGaps flags metrics whose latest rolling window is missing or has weak evidence.
func detectEvidenceGaps(employeesList []employees.Employee, health []EmployeeHealth, entriesByEmployee map[int64][]EntryAdapter, currentYear, currentMonth int) []EvidenceGap {
	gaps := make([]EvidenceGap, 0)
	healthByID := make(map[int64]EmployeeHealth)
	for _, h := range health {
		healthByID[h.EmployeeID] = h
	}

	for _, emp := range employeesList {
		latest := latestRollingEntries(entriesByEmployee[emp.ID], currentYear, currentMonth)
		if len(latest) == 0 {
			gaps = append(gaps, EvidenceGap{
				EmployeeID: emp.ID,
				Name:       fmt.Sprintf("%s %s", emp.FirstName, emp.SecondName),
				Metric:     "all",
				Message:    "No performance entries for the last three months",
			})
			continue
		}

		coveredMonths := 0
		for _, entry := range latest {
			if entry.Year != 0 {
				coveredMonths++
			}
		}
		if coveredMonths < 3 {
			gaps = append(gaps, EvidenceGap{
				EmployeeID: emp.ID,
				Name:       fmt.Sprintf("%s %s", emp.FirstName, emp.SecondName),
				Metric:     "rolling window",
				Message:    fmt.Sprintf("Only %d of 3 months have performance data", coveredMonths),
			})
		}

		h := healthByID[emp.ID]
		if h.EvidenceCount == 0 {
			gaps = append(gaps, EvidenceGap{
				EmployeeID: emp.ID,
				Name:       h.Name,
				Metric:     "evidence",
				Message:    "Latest month has no recorded evidence items",
			})
		}
	}
	return gaps
}

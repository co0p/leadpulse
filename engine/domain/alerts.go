package domain

// AlertSeverity represents the severity level of a raised alert.
type AlertSeverity int

const (
	AlertSeverityNone AlertSeverity = iota
	AlertSeverityAmber
	AlertSeverityRed
)

// EvaluatePerformanceDeterioration evaluates the Performance Deterioration
// alert condition for a single member's Delta1 and Delta3 trend values.
// PRD 7.1: Red: Delta1 <= -10 OR Delta3 <= -15
func EvaluatePerformanceDeterioration(delta1, delta3 float64) AlertSeverity {
	if delta1 <= -10 || delta3 <= -15 {
		return AlertSeverityRed
	}
	if delta1 <= -6 || delta3 <= -10 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

// EvaluateMoraleRisk evaluates the Morale Risk alert condition for a
// member's current and prior month MoraleN normalized scores.
// PRD 7.1:
//   Amber: MoraleN < 50 for 1 month
//   Red: MoraleN < 50 for 2 consecutive months OR MoraleN < 35 current month
func EvaluateMoraleRisk(currentMoraleN, priorMoraleN float64, hasPriorMonth bool) AlertSeverity {
	if currentMoraleN < 35 {
		return AlertSeverityRed
	}
	if hasPriorMonth && currentMoraleN < 50 && priorMoraleN < 50 {
		return AlertSeverityRed
	}
	if currentMoraleN < 50 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

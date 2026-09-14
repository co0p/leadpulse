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
	return AlertSeverityNone
}

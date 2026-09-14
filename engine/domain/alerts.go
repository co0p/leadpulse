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

// EvaluateBurnoutRisk evaluates the Burnout Risk alert condition for a
// member's overtime, morale, performance delta, and delivery reliability.
// PRD 7.1:
//   Amber: OvertimeHours >= 20 AND MoraleN < 60
//   Red: OvertimeHours >= 25 AND (Delta1 < 0 OR DeliveryN < 60)
func EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN float64) AlertSeverity {
	if overtimeHours >= 25 && (delta1 < 0 || deliveryN < 60) {
		return AlertSeverityRed
	}
	if overtimeHours >= 20 && moraleN < 60 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

// EvaluateFeedbackRisk evaluates the Feedback Risk alert condition for a
// member's current and prior month critical feedback counts.
// PRD 7.1:
//   Amber: CriticalFeedbackCount >= 3
//   Red: CriticalFeedbackCount >= 4 AND declining (worsening, non-decreasing)
//        2-month critical trend
// "Declining 2-month critical trend" is interpreted as the problem not
// improving: current month's count is not lower than the prior month's.
// Without prior-month data, Red cannot be confirmed; max severity is Amber.
func EvaluateFeedbackRisk(currentCritical, priorCritical int, hasPriorMonth bool) AlertSeverity {
	if currentCritical >= 4 && hasPriorMonth && currentCritical >= priorCritical {
		return AlertSeverityRed
	}
	if currentCritical >= 3 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

// EvaluateCustomerBusinessRisk evaluates the Customer/Business Risk alert
// condition for a member's CSAT and net margin normalized scores.
// PRD 7.1:
//   Amber: CSATN < 60 OR MarginN < 45
//   Red: CSATN < 50 AND MarginN < 40
func EvaluateCustomerBusinessRisk(csatN, marginN float64) AlertSeverity {
	if csatN < 60 || marginN < 45 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

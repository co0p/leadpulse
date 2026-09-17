package domain

import "fmt"

// AlertSeverity represents the severity level of a raised alert.
type AlertSeverity int

const (
	AlertSeverityNone AlertSeverity = iota
	AlertSeverityAmber
	AlertSeverityRed
)

// AlertType identifies which alert condition triggered an Alert.
type AlertType int

const (
	AlertTypePerformanceDeterioration AlertType = iota
	AlertTypeMoraleRisk
	AlertTypeBurnoutRisk
	AlertTypeFeedbackRisk
	AlertTypeCustomerBusinessRisk
	AlertTypeDataQualityRisk
	AlertTypeTeamMoraleDrift
	AlertTypeTeamDeliveryDrift
	AlertTypeSystemicBurnout
	AlertTypeCalibrationRisk
)

// Alert represents a single triggered alert for a member, carrying enough
// detail to be understood without re-deriving it from raw scores.
type Alert struct {
	Type     AlertType
	Severity AlertSeverity
	MemberID TeamMemberID
}

// MemberAlertInputs bundles all the pre-computed values EvaluateMemberAlerts
// needs to evaluate all 6 individual alert conditions for one member.
type MemberAlertInputs struct {
	MemberID TeamMemberID

	// Performance Deterioration
	Delta1 float64
	Delta3 float64

	// Morale Risk
	CurrentMoraleN float64
	PriorMoraleN   float64
	HasPriorMonth  bool

	// Burnout Risk
	OvertimeHours float64
	DeliveryN     float64

	// Feedback Risk
	CurrentCritical int
	PriorCritical   int

	// Customer/Business Risk
	CSATN   float64
	MarginN float64

	// Data Quality Risk
	CompletenessPct float64
}

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
//
//	Amber: MoraleN < 50 for 1 month
//	Red: MoraleN < 50 for 2 consecutive months OR MoraleN < 35 current month
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
//
//	Amber: OvertimeHours >= 20 AND MoraleN < 60
//	Red: OvertimeHours >= 25 AND (Delta1 < 0 OR DeliveryN < 60)
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
//
//	Amber: CriticalFeedbackCount >= 3
//	Red: CriticalFeedbackCount >= 4 AND declining (worsening, non-decreasing)
//	     2-month critical trend
//
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
//
//	Amber: CSATN < 60 OR MarginN < 45
//	Red: CSATN < 50 AND MarginN < 40
func EvaluateCustomerBusinessRisk(csatN, marginN float64) AlertSeverity {
	if csatN < 50 && marginN < 40 {
		return AlertSeverityRed
	}
	if csatN < 60 || marginN < 45 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

// EvaluateDataQualityRisk evaluates the Data Quality Risk alert condition
// for a member's completeness percentage.
// PRD 7.1:
//
//	Amber: CompletenessPct < 85
//	Red: CompletenessPct < 70
func EvaluateDataQualityRisk(completenessPct float64) AlertSeverity {
	if completenessPct < 70 {
		return AlertSeverityRed
	}
	if completenessPct < 85 {
		return AlertSeverityAmber
	}
	return AlertSeverityNone
}

// EvaluateMemberAlerts evaluates all 6 individual alert conditions for one
// member and returns an Alert for each condition that triggers (Amber or
// Red). Conditions that evaluate to None are excluded from the result.
func EvaluateMemberAlerts(inputs MemberAlertInputs) []Alert {
	conditions := []struct {
		alertType AlertType
		severity  AlertSeverity
	}{
		{AlertTypePerformanceDeterioration, EvaluatePerformanceDeterioration(inputs.Delta1, inputs.Delta3)},
		{AlertTypeMoraleRisk, EvaluateMoraleRisk(inputs.CurrentMoraleN, inputs.PriorMoraleN, inputs.HasPriorMonth)},
		{AlertTypeBurnoutRisk, EvaluateBurnoutRisk(inputs.OvertimeHours, inputs.CurrentMoraleN, inputs.Delta1, inputs.DeliveryN)},
		{AlertTypeFeedbackRisk, EvaluateFeedbackRisk(inputs.CurrentCritical, inputs.PriorCritical, inputs.HasPriorMonth)},
		{AlertTypeCustomerBusinessRisk, EvaluateCustomerBusinessRisk(inputs.CSATN, inputs.MarginN)},
		{AlertTypeDataQualityRisk, EvaluateDataQualityRisk(inputs.CompletenessPct)},
	}

	var alerts []Alert
	for _, c := range conditions {
		if c.severity != AlertSeverityNone {
			alerts = append(alerts, Alert{Type: c.alertType, Severity: c.severity, MemberID: inputs.MemberID})
		}
	}
	return alerts
}

// EvaluateTeamMoraleDrift evaluates the Team Morale Drift alert condition
// for the percentage of team members with Morale Red status.
// PRD 7.2: Team Morale Drift Red: >=30% members have Morale Red
func EvaluateTeamMoraleDrift(pctMembersWithMoraleRed float64) AlertSeverity {
	if pctMembersWithMoraleRed >= 30 {
		return AlertSeverityRed
	}
	return AlertSeverityNone
}

// EvaluateTeamDeliveryDrift evaluates the Team Delivery Drift alert
// condition for the team's aggregate Delta3.
// PRD 7.2: Team Delivery Drift Red: team Delta3 <= -10
func EvaluateTeamDeliveryDrift(teamDelta3 float64) AlertSeverity {
	if teamDelta3 <= -10 {
		return AlertSeverityRed
	}
	return AlertSeverityNone
}

// EvaluateSystemicBurnout evaluates the Systemic Burnout alert condition
// for the percentage of team members with Burnout Red status.
// PRD 7.2: Systemic Burnout Red: >=25% members Burnout Red
func EvaluateSystemicBurnout(pctMembersWithBurnoutRed float64) AlertSeverity {
	if pctMembersWithBurnoutRed >= 25 {
		return AlertSeverityRed
	}
	return AlertSeverityNone
}

// EvaluateCalibrationRisk evaluates the Calibration Risk alert condition
// for the team's TII standard deviation history. Expects the 3 most recent
// monthly stddev values, oldest first or newest first (order does not
// matter — all 3 months must satisfy the threshold).
// PRD 7.2: Calibration Risk Amber: team TII stddev < 6 for 3 months
func EvaluateCalibrationRisk(stddevHistory []float64) (AlertSeverity, error) {
	if len(stddevHistory) < 3 {
		return AlertSeverityNone, fmt.Errorf("insufficient data for Calibration Risk: need 3 months, got %d", len(stddevHistory))
	}

	for _, stddev := range stddevHistory {
		if stddev >= 6 {
			return AlertSeverityNone, nil
		}
	}
	return AlertSeverityAmber, nil
}

// TeamAlertInputs bundles all the pre-computed values EvaluateTeamAlerts
// needs to evaluate all 4 team-level alert conditions.
type TeamAlertInputs struct {
	PctMembersWithMoraleRed  float64
	TeamDelta3               float64
	PctMembersWithBurnoutRed float64
	StddevHistory            []float64
}

// EvaluateTeamAlerts evaluates all 4 team-level alert conditions and
// returns an Alert for each condition that triggers (Amber or Red).
// Conditions that evaluate to None are excluded from the result.
func EvaluateTeamAlerts(inputs TeamAlertInputs) ([]Alert, error) {
	calibrationSeverity, err := EvaluateCalibrationRisk(inputs.StddevHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate calibration risk: %w", err)
	}

	conditions := []struct {
		alertType AlertType
		severity  AlertSeverity
	}{
		{AlertTypeTeamMoraleDrift, EvaluateTeamMoraleDrift(inputs.PctMembersWithMoraleRed)},
		{AlertTypeTeamDeliveryDrift, EvaluateTeamDeliveryDrift(inputs.TeamDelta3)},
		{AlertTypeSystemicBurnout, EvaluateSystemicBurnout(inputs.PctMembersWithBurnoutRed)},
		{AlertTypeCalibrationRisk, calibrationSeverity},
	}

	var alerts []Alert
	for _, c := range conditions {
		if c.severity != AlertSeverityNone {
			alerts = append(alerts, Alert{Type: c.alertType, Severity: c.severity})
		}
	}
	return alerts, nil
}

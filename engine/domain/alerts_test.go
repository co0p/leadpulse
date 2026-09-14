package domain

import "testing"

// TestEvaluatePerformanceDeterioration_redAtDelta1Boundary tests that
// EvaluatePerformanceDeterioration returns Red when Delta1 is exactly -10.
// PRD 7.1: Performance Deterioration Red: Delta1 <= -10 OR Delta3 <= -15
func TestEvaluatePerformanceDeterioration_redAtDelta1Boundary(t *testing.T) {
	// Arrange: Delta1 = -10 (Red boundary), Delta3 = 0 (not a factor)
	delta1 := -10.0
	delta3 := 0.0

	// Act
	severity := EvaluatePerformanceDeterioration(delta1, delta3)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluatePerformanceDeterioration(delta1=%v, delta3=%v) = %v, want %v", delta1, delta3, severity, AlertSeverityRed)
	}
}

// TestEvaluatePerformanceDeterioration_amberAtDelta1Boundary tests that
// EvaluatePerformanceDeterioration returns Amber when Delta1 is exactly -6
// (below the Amber threshold but not severe enough for Red).
// PRD 7.1: Amber: Delta1 <= -6 OR Delta3 <= -10
func TestEvaluatePerformanceDeterioration_amberAtDelta1Boundary(t *testing.T) {
	// Arrange: Delta1 = -6 (Amber boundary), Delta3 = 0 (not a factor)
	delta1 := -6.0
	delta3 := 0.0

	// Act
	severity := EvaluatePerformanceDeterioration(delta1, delta3)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluatePerformanceDeterioration(delta1=%v, delta3=%v) = %v, want %v", delta1, delta3, severity, AlertSeverityAmber)
	}
}

// TestEvaluatePerformanceDeterioration_amberAtDelta3Boundary tests that
// EvaluatePerformanceDeterioration returns Amber when Delta3 is exactly -10,
// even when Delta1 alone would not trigger any alert.
// PRD 7.1: Amber: Delta1 <= -6 OR Delta3 <= -10
func TestEvaluatePerformanceDeterioration_amberAtDelta3Boundary(t *testing.T) {
	// Arrange: Delta1 = 0 (not a factor), Delta3 = -10 (Amber boundary)
	delta1 := 0.0
	delta3 := -10.0

	// Act
	severity := EvaluatePerformanceDeterioration(delta1, delta3)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluatePerformanceDeterioration(delta1=%v, delta3=%v) = %v, want %v", delta1, delta3, severity, AlertSeverityAmber)
	}
}

// TestEvaluatePerformanceDeterioration_redAtDelta3Boundary tests that
// EvaluatePerformanceDeterioration returns Red when Delta3 is exactly -15,
// even when Delta1 alone would not trigger Red.
// PRD 7.1: Performance Deterioration Red: Delta1 <= -10 OR Delta3 <= -15
func TestEvaluatePerformanceDeterioration_redAtDelta3Boundary(t *testing.T) {
	// Arrange: Delta1 = 0 (not a factor), Delta3 = -15 (Red boundary)
	delta1 := 0.0
	delta3 := -15.0

	// Act
	severity := EvaluatePerformanceDeterioration(delta1, delta3)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluatePerformanceDeterioration(delta1=%v, delta3=%v) = %v, want %v", delta1, delta3, severity, AlertSeverityRed)
	}
}

// TestEvaluatePerformanceDeterioration_none tests that
// EvaluatePerformanceDeterioration returns None when neither Delta1 nor
// Delta3 cross any threshold.
// PRD 7.1: no alert when Delta1 > -6 AND Delta3 > -10
func TestEvaluatePerformanceDeterioration_none(t *testing.T) {
	// Arrange: healthy trend, no deterioration
	delta1 := 2.0
	delta3 := 5.0

	// Act
	severity := EvaluatePerformanceDeterioration(delta1, delta3)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluatePerformanceDeterioration(delta1=%v, delta3=%v) = %v, want %v", delta1, delta3, severity, AlertSeverityNone)
	}
}

// TestEvaluateMoraleRisk_redWhenCurrentBelow35 tests that EvaluateMoraleRisk
// returns Red when the current month's MoraleN is below 35, regardless of
// prior month history.
// PRD 7.1: Morale Risk Red: MoraleN < 35 current month
func TestEvaluateMoraleRisk_redWhenCurrentBelow35(t *testing.T) {
	// Arrange: current month severely low morale, no prior month data
	currentMoraleN := 30.0
	priorMoraleN := 80.0 // healthy prior month; should not prevent Red via severe-current rule
	hasPriorMonth := true

	// Act
	severity := EvaluateMoraleRisk(currentMoraleN, priorMoraleN, hasPriorMonth)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateMoraleRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentMoraleN, priorMoraleN, hasPriorMonth, severity, AlertSeverityRed)
	}
}

// TestEvaluateMoraleRisk_redWhenTwoConsecutiveMonthsBelow50 tests that
// EvaluateMoraleRisk returns Red when both current and prior month MoraleN
// are below 50, even though neither alone is severe enough (< 35) to
// trigger Red via the single-month rule.
// PRD 7.1: Morale Risk Red: MoraleN < 50 for 2 consecutive months
func TestEvaluateMoraleRisk_redWhenTwoConsecutiveMonthsBelow50(t *testing.T) {
	// Arrange: both months moderately low (< 50 but >= 35)
	currentMoraleN := 45.0
	priorMoraleN := 40.0
	hasPriorMonth := true

	// Act
	severity := EvaluateMoraleRisk(currentMoraleN, priorMoraleN, hasPriorMonth)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateMoraleRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentMoraleN, priorMoraleN, hasPriorMonth, severity, AlertSeverityRed)
	}
}

// TestEvaluateMoraleRisk_amberWhenSingleMonthBelow50 tests that
// EvaluateMoraleRisk returns Amber when only the current month's MoraleN is
// below 50 and there is no prior month data, so the two-consecutive-month
// Red rule does not apply.
// PRD 7.1: Morale Risk Amber: MoraleN < 50 for 1 month
func TestEvaluateMoraleRisk_amberWhenSingleMonthBelow50(t *testing.T) {
	// Arrange: current month moderately low, no prior month data
	currentMoraleN := 45.0
	priorMoraleN := 0.0
	hasPriorMonth := false

	// Act
	severity := EvaluateMoraleRisk(currentMoraleN, priorMoraleN, hasPriorMonth)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateMoraleRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentMoraleN, priorMoraleN, hasPriorMonth, severity, AlertSeverityAmber)
	}
}

// TestEvaluateMoraleRisk_amberNotRedAtExactly35 tests that
// EvaluateMoraleRisk returns Amber (not Red) when current month MoraleN is
// exactly 35 — the severe-current-month Red rule requires strictly below
// 35, and no prior month data exists to trigger the sustained-Red rule.
// PRD 7.1: Red requires MoraleN < 35 (strict); exactly 35 does not qualify
func TestEvaluateMoraleRisk_amberNotRedAtExactly35(t *testing.T) {
	// Arrange: current month exactly at the Red boundary (not below it)
	currentMoraleN := 35.0
	priorMoraleN := 0.0
	hasPriorMonth := false

	// Act
	severity := EvaluateMoraleRisk(currentMoraleN, priorMoraleN, hasPriorMonth)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateMoraleRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentMoraleN, priorMoraleN, hasPriorMonth, severity, AlertSeverityAmber)
	}
}

// TestEvaluateMoraleRisk_noneWhenAtOrAbove50 tests that EvaluateMoraleRisk
// returns None when the current month's MoraleN is at or above 50, with no
// prior month history triggering a sustained condition.
// PRD 7.1: no alert when MoraleN >= 50
func TestEvaluateMoraleRisk_noneWhenAtOrAbove50(t *testing.T) {
	// Arrange: healthy current month, no prior month data
	currentMoraleN := 60.0
	priorMoraleN := 0.0
	hasPriorMonth := false

	// Act
	severity := EvaluateMoraleRisk(currentMoraleN, priorMoraleN, hasPriorMonth)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateMoraleRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentMoraleN, priorMoraleN, hasPriorMonth, severity, AlertSeverityNone)
	}
}

// TestEvaluateBurnoutRisk_amberAtBoundary tests that EvaluateBurnoutRisk
// returns Amber at the exact boundary: OvertimeHours = 20 AND MoraleN < 60.
// PRD 7.1: Burnout Risk Amber: OvertimeHours >= 20 AND MoraleN < 60
func TestEvaluateBurnoutRisk_amberAtBoundary(t *testing.T) {
	// Arrange: overtime at Amber boundary, morale below 60
	overtimeHours := 20.0
	moraleN := 55.0
	delta1 := 0.0
	deliveryN := 100.0

	// Act
	severity := EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateBurnoutRisk(overtime=%v, morale=%v, delta1=%v, delivery=%v) = %v, want %v", overtimeHours, moraleN, delta1, deliveryN, severity, AlertSeverityAmber)
	}
}

// TestEvaluateBurnoutRisk_noneWhenMoraleAtBoundary60 tests that
// EvaluateBurnoutRisk returns None when overtime is high but MoraleN is
// exactly 60 — the Amber rule requires MoraleN strictly below 60.
// PRD 7.1: Amber requires MoraleN < 60 (strict); exactly 60 does not qualify
func TestEvaluateBurnoutRisk_noneWhenMoraleAtBoundary60(t *testing.T) {
	// Arrange: overtime high, morale exactly at boundary (not below it)
	overtimeHours := 22.0
	moraleN := 60.0
	delta1 := 0.0
	deliveryN := 100.0

	// Act
	severity := EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateBurnoutRisk(overtime=%v, morale=%v, delta1=%v, delivery=%v) = %v, want %v", overtimeHours, moraleN, delta1, deliveryN, severity, AlertSeverityNone)
	}
}

// TestEvaluateBurnoutRisk_redViaNegativeDelta1 tests that EvaluateBurnoutRisk
// returns Red when OvertimeHours >= 25 AND Delta1 < 0, even when
// DeliveryN is healthy.
// PRD 7.1: Burnout Risk Red: OvertimeHours >= 25 AND (Delta1 < 0 OR DeliveryN < 60)
func TestEvaluateBurnoutRisk_redViaNegativeDelta1(t *testing.T) {
	// Arrange: overtime at Red boundary, negative delta, healthy delivery
	overtimeHours := 25.0
	moraleN := 100.0 // healthy; should not prevent Red via delta1 rule
	delta1 := -1.0
	deliveryN := 100.0

	// Act
	severity := EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateBurnoutRisk(overtime=%v, morale=%v, delta1=%v, delivery=%v) = %v, want %v", overtimeHours, moraleN, delta1, deliveryN, severity, AlertSeverityRed)
	}
}

// TestEvaluateBurnoutRisk_redViaLowDeliveryN tests that EvaluateBurnoutRisk
// returns Red when OvertimeHours >= 25 AND DeliveryN < 60, even when
// Delta1 is non-negative.
// PRD 7.1: Burnout Risk Red: OvertimeHours >= 25 AND (Delta1 < 0 OR DeliveryN < 60)
func TestEvaluateBurnoutRisk_redViaLowDeliveryN(t *testing.T) {
	// Arrange: overtime at Red boundary, non-negative delta, low delivery
	overtimeHours := 25.0
	moraleN := 100.0
	delta1 := 0.0
	deliveryN := 50.0

	// Act
	severity := EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateBurnoutRisk(overtime=%v, morale=%v, delta1=%v, delivery=%v) = %v, want %v", overtimeHours, moraleN, delta1, deliveryN, severity, AlertSeverityRed)
	}
}

// TestEvaluateBurnoutRisk_noneWhenOvertimeLow tests that EvaluateBurnoutRisk
// returns None when overtime is below the Amber threshold, regardless of
// other factors.
// PRD 7.1: no alert when OvertimeHours < 20
func TestEvaluateBurnoutRisk_noneWhenOvertimeLow(t *testing.T) {
	// Arrange: low overtime, other factors would trigger if overtime were high
	overtimeHours := 10.0
	moraleN := 30.0
	delta1 := -5.0
	deliveryN := 20.0

	// Act
	severity := EvaluateBurnoutRisk(overtimeHours, moraleN, delta1, deliveryN)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateBurnoutRisk(overtime=%v, morale=%v, delta1=%v, delivery=%v) = %v, want %v", overtimeHours, moraleN, delta1, deliveryN, severity, AlertSeverityNone)
	}
}


// TestEvaluateFeedbackRisk_amberAtThreeCritical tests that
// EvaluateFeedbackRisk returns Amber when CriticalFeedbackCount is exactly 3.
// PRD 7.1: Feedback Risk Amber: CriticalFeedbackCount >= 3
func TestEvaluateFeedbackRisk_amberAtThreeCritical(t *testing.T) {
	// Arrange: exactly 3 critical feedback items, no prior month data
	currentCritical := 3
	priorCritical := 0
	hasPriorMonth := false

	// Act
	severity := EvaluateFeedbackRisk(currentCritical, priorCritical, hasPriorMonth)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateFeedbackRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentCritical, priorCritical, hasPriorMonth, severity, AlertSeverityAmber)
	}
}

// TestEvaluateFeedbackRisk_redWhenFourPlusAndNonDecreasing tests that
// EvaluateFeedbackRisk returns Red when CriticalFeedbackCount >= 4 and the
// current month's count is not lower than the prior month's (worsening or
// stagnant trend).
// PRD 7.1: Red: CriticalFeedbackCount >= 4 AND declining 2-month critical trend
// (interpreted as non-decreasing, i.e., the problem is not improving)
func TestEvaluateFeedbackRisk_redWhenFourPlusAndNonDecreasing(t *testing.T) {
	// Arrange: current month >= 4, count did not decrease from prior month
	currentCritical := 4
	priorCritical := 4
	hasPriorMonth := true

	// Act
	severity := EvaluateFeedbackRisk(currentCritical, priorCritical, hasPriorMonth)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateFeedbackRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentCritical, priorCritical, hasPriorMonth, severity, AlertSeverityRed)
	}
}

// TestEvaluateFeedbackRisk_amberNotRedWhenImproving tests that
// EvaluateFeedbackRisk returns Amber (not Red) when CriticalFeedbackCount
// is >= 4 but has decreased from the prior month (improving trend).
// PRD 7.1: Red requires a non-decreasing (worsening) trend; an improving
// trend does not qualify even with count >= 4.
func TestEvaluateFeedbackRisk_amberNotRedWhenImproving(t *testing.T) {
	// Arrange: current month >= 4 but lower than prior month (improving)
	currentCritical := 4
	priorCritical := 6
	hasPriorMonth := true

	// Act
	severity := EvaluateFeedbackRisk(currentCritical, priorCritical, hasPriorMonth)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateFeedbackRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentCritical, priorCritical, hasPriorMonth, severity, AlertSeverityAmber)
	}
}

// TestEvaluateFeedbackRisk_amberWhenNoPriorMonthData tests that
// EvaluateFeedbackRisk returns Amber (not Red) when CriticalFeedbackCount
// is >= 4 but there is no prior month data to confirm a worsening trend.
// This is a deliberate conservative default: Red requires 2-month history.
func TestEvaluateFeedbackRisk_amberWhenNoPriorMonthData(t *testing.T) {
	// Arrange: current month severe, but no prior month data exists
	currentCritical := 5
	priorCritical := 0
	hasPriorMonth := false

	// Act
	severity := EvaluateFeedbackRisk(currentCritical, priorCritical, hasPriorMonth)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateFeedbackRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentCritical, priorCritical, hasPriorMonth, severity, AlertSeverityAmber)
	}
}

// TestEvaluateFeedbackRisk_noneWhenBelowThree tests that EvaluateFeedbackRisk
// returns None when CriticalFeedbackCount is below 3.
// PRD 7.1: no alert when CriticalFeedbackCount < 3
func TestEvaluateFeedbackRisk_noneWhenBelowThree(t *testing.T) {
	// Arrange: healthy feedback count
	currentCritical := 1
	priorCritical := 2
	hasPriorMonth := true

	// Act
	severity := EvaluateFeedbackRisk(currentCritical, priorCritical, hasPriorMonth)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateFeedbackRisk(current=%v, prior=%v, hasPrior=%v) = %v, want %v", currentCritical, priorCritical, hasPriorMonth, severity, AlertSeverityNone)
	}
}

// TestEvaluateCustomerBusinessRisk_amberAtCSATBoundary tests that
// EvaluateCustomerBusinessRisk returns Amber when CSATN is exactly 59
// (below 60) and MarginN is healthy.
// PRD 7.1: Customer/Business Risk Amber: CSATN < 60 OR MarginN < 45
func TestEvaluateCustomerBusinessRisk_amberAtCSATBoundary(t *testing.T) {
	// Arrange: CSAT just below Amber threshold, margin healthy
	csatN := 59.0
	marginN := 100.0

	// Act
	severity := EvaluateCustomerBusinessRisk(csatN, marginN)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateCustomerBusinessRisk(csatN=%v, marginN=%v) = %v, want %v", csatN, marginN, severity, AlertSeverityAmber)
	}
}

// TestEvaluateCustomerBusinessRisk_amberAtMarginBoundary tests that
// EvaluateCustomerBusinessRisk returns Amber when MarginN is exactly 44
// (below 45), even when CSATN is healthy.
// PRD 7.1: Customer/Business Risk Amber: CSATN < 60 OR MarginN < 45
func TestEvaluateCustomerBusinessRisk_amberAtMarginBoundary(t *testing.T) {
	// Arrange: CSAT healthy, margin just below Amber threshold
	csatN := 100.0
	marginN := 44.0

	// Act
	severity := EvaluateCustomerBusinessRisk(csatN, marginN)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateCustomerBusinessRisk(csatN=%v, marginN=%v) = %v, want %v", csatN, marginN, severity, AlertSeverityAmber)
	}
}

// TestEvaluateCustomerBusinessRisk_redAtBothBoundaries tests that
// EvaluateCustomerBusinessRisk returns Red when both CSATN < 50 and
// MarginN < 40 at their exact boundaries.
// PRD 7.1: Customer/Business Risk Red: CSATN < 50 AND MarginN < 40
func TestEvaluateCustomerBusinessRisk_redAtBothBoundaries(t *testing.T) {
	// Arrange: both CSAT and margin below Red thresholds
	csatN := 49.0
	marginN := 39.0

	// Act
	severity := EvaluateCustomerBusinessRisk(csatN, marginN)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateCustomerBusinessRisk(csatN=%v, marginN=%v) = %v, want %v", csatN, marginN, severity, AlertSeverityRed)
	}
}

// TestEvaluateCustomerBusinessRisk_amberNotRedWhenOnlyOneConditionMet tests
// that EvaluateCustomerBusinessRisk returns Amber (not Red) when only the
// CSAT condition for Red is met but MarginN is healthy.
// PRD 7.1: Red requires BOTH CSATN < 50 AND MarginN < 40
func TestEvaluateCustomerBusinessRisk_amberNotRedWhenOnlyOneConditionMet(t *testing.T) {
	// Arrange: CSAT meets Red threshold, margin healthy
	csatN := 45.0
	marginN := 100.0

	// Act
	severity := EvaluateCustomerBusinessRisk(csatN, marginN)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateCustomerBusinessRisk(csatN=%v, marginN=%v) = %v, want %v", csatN, marginN, severity, AlertSeverityAmber)
	}
}

// TestEvaluateCustomerBusinessRisk_noneWhenBothHealthy tests that
// EvaluateCustomerBusinessRisk returns None when both CSATN and MarginN
// are healthy.
// PRD 7.1: no alert when CSATN >= 60 AND MarginN >= 45
func TestEvaluateCustomerBusinessRisk_noneWhenBothHealthy(t *testing.T) {
	// Arrange: both metrics healthy
	csatN := 80.0
	marginN := 70.0

	// Act
	severity := EvaluateCustomerBusinessRisk(csatN, marginN)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateCustomerBusinessRisk(csatN=%v, marginN=%v) = %v, want %v", csatN, marginN, severity, AlertSeverityNone)
	}
}

// TestEvaluateDataQualityRisk_amberBelow85 tests that EvaluateDataQualityRisk
// returns Amber when CompletenessPct is exactly 84 (below 85).
// PRD 7.1: Data Quality Risk Amber: CompletenessPct < 85
func TestEvaluateDataQualityRisk_amberBelow85(t *testing.T) {
	// Arrange: completeness just below Amber threshold
	completenessPct := 84.0

	// Act
	severity := EvaluateDataQualityRisk(completenessPct)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateDataQualityRisk(completenessPct=%v) = %v, want %v", completenessPct, severity, AlertSeverityAmber)
	}
}

// TestEvaluateDataQualityRisk_redBelow70 tests that EvaluateDataQualityRisk
// returns Red when CompletenessPct is exactly 69 (below 70).
// PRD 7.1: Data Quality Risk Red: CompletenessPct < 70
func TestEvaluateDataQualityRisk_redBelow70(t *testing.T) {
	// Arrange: completeness below Red threshold
	completenessPct := 69.0

	// Act
	severity := EvaluateDataQualityRisk(completenessPct)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateDataQualityRisk(completenessPct=%v) = %v, want %v", completenessPct, severity, AlertSeverityRed)
	}
}

// TestEvaluateDataQualityRisk_amberNotRedAtExactly70 tests that
// EvaluateDataQualityRisk returns Amber (not Red) when CompletenessPct is
// exactly 70 — the Red rule requires strictly below 70.
// PRD 7.1: Red requires CompletenessPct < 70 (strict); exactly 70 does not qualify
func TestEvaluateDataQualityRisk_amberNotRedAtExactly70(t *testing.T) {
	// Arrange: completeness exactly at Red boundary (not below it)
	completenessPct := 70.0

	// Act
	severity := EvaluateDataQualityRisk(completenessPct)

	// Assert
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateDataQualityRisk(completenessPct=%v) = %v, want %v", completenessPct, severity, AlertSeverityAmber)
	}
}

// TestEvaluateDataQualityRisk_noneAtExactly85 tests that
// EvaluateDataQualityRisk returns None when CompletenessPct is exactly 85 —
// the Amber rule requires strictly below 85.
// PRD 7.1: no alert when CompletenessPct >= 85
func TestEvaluateDataQualityRisk_noneAtExactly85(t *testing.T) {
	// Arrange: completeness exactly at Amber boundary (not below it)
	completenessPct := 85.0

	// Act
	severity := EvaluateDataQualityRisk(completenessPct)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateDataQualityRisk(completenessPct=%v) = %v, want %v", completenessPct, severity, AlertSeverityNone)
	}
}

// TestEvaluateMemberAlerts_returnsAllTriggeredAlerts tests that
// EvaluateMemberAlerts aggregates all 6 individual alert evaluators and
// returns an Alert for each condition that triggers (Amber or Red).
func TestEvaluateMemberAlerts_returnsAllTriggeredAlerts(t *testing.T) {
	// Arrange: inputs that trigger both Performance Deterioration (Red)
	// and Data Quality Risk (Amber), with all other conditions healthy
	memberID := TeamMemberID(1)
	inputs := MemberAlertInputs{
		MemberID:         memberID,
		Delta1:           -10, // triggers Performance Deterioration Red
		Delta3:           0,
		CurrentMoraleN:   100,
		PriorMoraleN:     0,
		HasPriorMonth:    false,
		OvertimeHours:    0,
		DeliveryN:        100,
		CurrentCritical:  0,
		PriorCritical:    0,
		CSATN:            100,
		MarginN:          100,
		CompletenessPct:  80, // triggers Data Quality Risk Amber
	}

	// Act
	alerts := EvaluateMemberAlerts(inputs)

	// Assert: exactly 2 alerts triggered
	if len(alerts) != 2 {
		t.Fatalf("EvaluateMemberAlerts() returned %d alerts, want 2: %+v", len(alerts), alerts)
	}
}

// TestEvaluateMemberAlerts_excludesNonTriggeredConditions tests that
// EvaluateMemberAlerts returns an empty slice when all 6 conditions
// evaluate to None (healthy member, no alerts).
func TestEvaluateMemberAlerts_excludesNonTriggeredConditions(t *testing.T) {
	// Arrange: all inputs healthy, no condition should trigger
	memberID := TeamMemberID(2)
	inputs := MemberAlertInputs{
		MemberID:        memberID,
		Delta1:          5,
		Delta3:          10,
		CurrentMoraleN:  90,
		PriorMoraleN:    90,
		HasPriorMonth:   true,
		OvertimeHours:   5,
		DeliveryN:       100,
		CurrentCritical: 0,
		PriorCritical:   0,
		CSATN:           100,
		MarginN:         100,
		CompletenessPct: 100,
	}

	// Act
	alerts := EvaluateMemberAlerts(inputs)

	// Assert
	if len(alerts) != 0 {
		t.Errorf("EvaluateMemberAlerts() returned %d alerts, want 0: %+v", len(alerts), alerts)
	}
}

// TestEvaluateMemberAlerts_detailIdentifiesConditionAndMember tests that
// each returned Alert carries enough detail (condition type, severity, and
// member ID) to be understood without re-deriving it from raw scores.
// PRD/increment AC-3: alert detail sufficiency.
func TestEvaluateMemberAlerts_detailIdentifiesConditionAndMember(t *testing.T) {
	// Arrange: trigger exactly one condition (Data Quality Risk, Red)
	memberID := TeamMemberID(7)
	inputs := MemberAlertInputs{
		MemberID:        memberID,
		Delta1:          5,
		Delta3:          10,
		CurrentMoraleN:  90,
		PriorMoraleN:    90,
		HasPriorMonth:   true,
		OvertimeHours:   5,
		DeliveryN:       100,
		CurrentCritical: 0,
		PriorCritical:   0,
		CSATN:           100,
		MarginN:         100,
		CompletenessPct: 50, // triggers Data Quality Risk Red
	}

	// Act
	alerts := EvaluateMemberAlerts(inputs)

	// Assert: exactly one alert, with full identifying detail
	if len(alerts) != 1 {
		t.Fatalf("EvaluateMemberAlerts() returned %d alerts, want 1: %+v", len(alerts), alerts)
	}
	alert := alerts[0]
	if alert.Type != AlertTypeDataQualityRisk {
		t.Errorf("alert.Type = %v, want %v", alert.Type, AlertTypeDataQualityRisk)
	}
	if alert.Severity != AlertSeverityRed {
		t.Errorf("alert.Severity = %v, want %v", alert.Severity, AlertSeverityRed)
	}
	if alert.MemberID != memberID {
		t.Errorf("alert.MemberID = %v, want %v", alert.MemberID, memberID)
	}
}

// TestEvaluateTeamMoraleDrift_redAt30Percent tests that
// EvaluateTeamMoraleDrift returns Red when exactly 30% of members have
// Morale Red.
// PRD 7.2: Team Morale Drift Red: >=30% members have Morale Red
func TestEvaluateTeamMoraleDrift_redAt30Percent(t *testing.T) {
	// Arrange: 30% of members have Morale Red
	pctMembersWithMoraleRed := 30.0

	// Act
	severity := EvaluateTeamMoraleDrift(pctMembersWithMoraleRed)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateTeamMoraleDrift(pct=%v) = %v, want %v", pctMembersWithMoraleRed, severity, AlertSeverityRed)
	}
}

// TestEvaluateTeamMoraleDrift_noneBelow30Percent tests that
// EvaluateTeamMoraleDrift returns None when the percentage of members with
// Morale Red is below 30%.
// PRD 7.2: no alert when pctMembersWithMoraleRed < 30
func TestEvaluateTeamMoraleDrift_noneBelow30Percent(t *testing.T) {
	// Arrange: below the Red threshold
	pctMembersWithMoraleRed := 29.0

	// Act
	severity := EvaluateTeamMoraleDrift(pctMembersWithMoraleRed)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateTeamMoraleDrift(pct=%v) = %v, want %v", pctMembersWithMoraleRed, severity, AlertSeverityNone)
	}
}

// TestEvaluateTeamDeliveryDrift_redAtNegativeTenDelta3 tests that
// EvaluateTeamDeliveryDrift returns Red when the team's Delta3 is exactly
// -10.
// PRD 7.2: Team Delivery Drift Red: team Delta3 <= -10
func TestEvaluateTeamDeliveryDrift_redAtNegativeTenDelta3(t *testing.T) {
	// Arrange: team Delta3 at Red boundary
	teamDelta3 := -10.0

	// Act
	severity := EvaluateTeamDeliveryDrift(teamDelta3)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateTeamDeliveryDrift(delta3=%v) = %v, want %v", teamDelta3, severity, AlertSeverityRed)
	}
}

// TestEvaluateTeamDeliveryDrift_noneAboveThreshold tests that
// EvaluateTeamDeliveryDrift returns None when the team's Delta3 is above
// the Red threshold.
// PRD 7.2: no alert when team Delta3 > -10
func TestEvaluateTeamDeliveryDrift_noneAboveThreshold(t *testing.T) {
	// Arrange: healthy team Delta3
	teamDelta3 := 5.0

	// Act
	severity := EvaluateTeamDeliveryDrift(teamDelta3)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateTeamDeliveryDrift(delta3=%v) = %v, want %v", teamDelta3, severity, AlertSeverityNone)
	}
}

// TestEvaluateSystemicBurnout_redAt25Percent tests that
// EvaluateSystemicBurnout returns Red when exactly 25% of members have
// Burnout Red.
// PRD 7.2: Systemic Burnout Red: >=25% members Burnout Red
func TestEvaluateSystemicBurnout_redAt25Percent(t *testing.T) {
	// Arrange: 25% of members have Burnout Red
	pctMembersWithBurnoutRed := 25.0

	// Act
	severity := EvaluateSystemicBurnout(pctMembersWithBurnoutRed)

	// Assert
	if severity != AlertSeverityRed {
		t.Errorf("EvaluateSystemicBurnout(pct=%v) = %v, want %v", pctMembersWithBurnoutRed, severity, AlertSeverityRed)
	}
}

// TestEvaluateSystemicBurnout_noneBelow25Percent tests that
// EvaluateSystemicBurnout returns None when the percentage of members with
// Burnout Red is below 25%.
// PRD 7.2: no alert when pctMembersWithBurnoutRed < 25
func TestEvaluateSystemicBurnout_noneBelow25Percent(t *testing.T) {
	// Arrange: below the Red threshold
	pctMembersWithBurnoutRed := 24.0

	// Act
	severity := EvaluateSystemicBurnout(pctMembersWithBurnoutRed)

	// Assert
	if severity != AlertSeverityNone {
		t.Errorf("EvaluateSystemicBurnout(pct=%v) = %v, want %v", pctMembersWithBurnoutRed, severity, AlertSeverityNone)
	}
}

// TestEvaluateCalibrationRisk_amberWhenStddevBelow6ForThreeMonths tests
// that EvaluateCalibrationRisk returns Amber when the team TII stddev has
// been below 6 for 3 consecutive months.
// PRD 7.2: Calibration Risk Amber: team TII stddev < 6 for 3 months
func TestEvaluateCalibrationRisk_amberWhenStddevBelow6ForThreeMonths(t *testing.T) {
	// Arrange: 3 months of low stddev (compressed team scores)
	stddevHistory := []float64{5.0, 4.5, 3.0}

	// Act
	severity, err := EvaluateCalibrationRisk(stddevHistory)

	// Assert
	if err != nil {
		t.Fatalf("EvaluateCalibrationRisk(%v) returned error: %v", stddevHistory, err)
	}
	if severity != AlertSeverityAmber {
		t.Errorf("EvaluateCalibrationRisk(%v) = %v, want %v", stddevHistory, severity, AlertSeverityAmber)
	}
}

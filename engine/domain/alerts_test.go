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

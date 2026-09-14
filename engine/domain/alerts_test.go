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


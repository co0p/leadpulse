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

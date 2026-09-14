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

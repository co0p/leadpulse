package scoring

import (
	"math"
	"testing"
)

// TestClamp_withinBounds tests that Clamp returns the value unchanged when it is within bounds.
// PRD 5.1: clamp(x, min, max) = min(max(x, min), max)
func TestClamp_withinBounds(t *testing.T) {
	result := Clamp(50.0, 0.0, 100.0)
	if result != 50.0 {
		t.Errorf("expected 50.0, got %v", result)
	}
}

// TestClamp_belowMin tests that Clamp returns min when value is below minimum.
func TestClamp_belowMin(t *testing.T) {
	result := Clamp(-10.0, 0.0, 100.0)
	if result != 0.0 {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestClamp_aboveMax tests that Clamp returns max when value is above maximum.
func TestClamp_aboveMax(t *testing.T) {
	result := Clamp(150.0, 0.0, 100.0)
	if result != 100.0 {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestNorm01_minBoundary tests that Norm01 returns 0 when x equals min.
// PRD 5.1: norm01(x, min, max) = clamp((x - min)/(max - min), 0, 1)
func TestNorm01_minBoundary(t *testing.T) {
	result := Norm01(-20.0, -20.0, 60.0)
	if !almostEqual(result, 0.0, 0.0001) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNorm01_maxBoundary tests that Norm01 returns 1 when x equals max.
func TestNorm01_maxBoundary(t *testing.T) {
	result := Norm01(60.0, -20.0, 60.0)
	if !almostEqual(result, 1.0, 0.0001) {
		t.Errorf("expected 1.0, got %v", result)
	}
}

// TestNorm01_midpoint tests that Norm01 returns 0.5 when x is at the midpoint.
func TestNorm01_midpoint(t *testing.T) {
	result := Norm01(20.0, -20.0, 60.0)
	if !almostEqual(result, 0.5, 0.0001) {
		t.Errorf("expected 0.5, got %v", result)
	}
}

// almostEqual compares two floats with a tolerance for floating-point precision.
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

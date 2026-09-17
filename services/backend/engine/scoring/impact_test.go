package scoring

import (
	"testing"
)

// TestComputeImpactWeighted_allZero tests that all-zero impact ratings produce 0 impact score.
// PRD 5.3: ImpactWeighted = 0.20*IG + 0.35*IP + 0.25*IT + 0.20*IO (0–5)
//
//	ImpactWeighted100 = ImpactWeighted * 20 (0–100)
func TestComputeImpactWeighted_allZero(t *testing.T) {
	result := ComputeImpactWeighted(0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeImpactWeighted_allMax tests that all-max (5) impact ratings produce 100 impact score.
func TestComputeImpactWeighted_allMax(t *testing.T) {
	result := ComputeImpactWeighted(5, 5, 5, 5)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeImpactWeighted_mixed tests a mixed impact rating case.
// IG=3, IP=4, IT=2, IO=3
// ImpactWeighted = 0.20*3 + 0.35*4 + 0.25*2 + 0.20*3 = 0.6 + 1.4 + 0.5 + 0.6 = 3.1
// ImpactWeighted100 = 3.1 * 20 = 62
func TestComputeImpactWeighted_mixed(t *testing.T) {
	result := ComputeImpactWeighted(3, 4, 2, 3)
	if !almostEqual(result, 62.0, 0.01) {
		t.Errorf("expected 62.0, got %v", result)
	}
}

// TestComputeImpactWeighted_typicalCaseProjects tests a typical project-focused case.
// IG=2, IP=4, IT=2, IO=2 (focus on project impact)
// ImpactWeighted = 0.20*2 + 0.35*4 + 0.25*2 + 0.20*2 = 0.4 + 1.4 + 0.5 + 0.4 = 2.7
// ImpactWeighted100 = 2.7 * 20 = 54
func TestComputeImpactWeighted_typicalCaseProjects(t *testing.T) {
	result := ComputeImpactWeighted(2, 4, 2, 2)
	if !almostEqual(result, 54.0, 0.01) {
		t.Errorf("expected 54.0, got %v", result)
	}
}

package scoring

import (
	"testing"
)

// TestComputeDimensionGrowth_allZero tests DG with all zero contributions.
// PRD 5.4: DG = Σ(weight_i * C_i) / 100
// With all contributions zero, DG = 0
func TestComputeDimensionGrowth_allZero(t *testing.T) {
	result := ComputeDimensionGrowth(0, 0, 0, 0, 0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeDimensionGrowth_allMax tests DG with all contributions at max (100).
func TestComputeDimensionGrowth_allMax(t *testing.T) {
	result := ComputeDimensionGrowth(100, 100, 100, 100, 100, 100, 100, 100)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeDimensionGrowth_typicalEntry tests DG with typical mixed values.
// DG weights per PRD 5.4: MoraleN 20, CriticalN 15, PositiveN 10, MentoringN 15, DeliveryN 10, BillabilityN 10, OvertimeN 15, EvidenceN 5
// Test with: Morale=80, Critical=70, Positive=60, Mentoring=90, Delivery=75, Billability=80, Overtime=50, Evidence=40
// DG = (20*80 + 15*70 + 10*60 + 15*90 + 10*75 + 10*80 + 15*50 + 5*40) / 100
//
//	= (1600 + 1050 + 600 + 1350 + 750 + 800 + 750 + 200) / 100
//	= 7100 / 100 = 71
func TestComputeDimensionGrowth_typicalEntry(t *testing.T) {
	result := ComputeDimensionGrowth(80, 70, 60, 90, 75, 80, 50, 40)
	if !almostEqual(result, 71.0, 0.01) {
		t.Errorf("expected 71.0, got %v", result)
	}
}

// TestComputeDimensionProject_allZero tests DP with all zero contributions.
func TestComputeDimensionProject_allZero(t *testing.T) {
	result := ComputeDimensionProject(0, 0, 0, 0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeDimensionProject_allMax tests DP with all contributions at max (100).
func TestComputeDimensionProject_allMax(t *testing.T) {
	result := ComputeDimensionProject(100, 100, 100, 100, 100, 100, 100)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeDimensionProject_typicalEntry tests DP with typical mixed values.
// DP weights per PRD 5.4: DeliveryN 25, CSATN 20, MarginN 20, BillabilityN 15, CriticalN 10, PositiveN 5, MoraleN 5
// Test with: Delivery=85, CSAT=80, Margin=70, Billability=75, Critical=60, Positive=70, Morale=75
// DP = (25*85 + 20*80 + 20*70 + 15*75 + 10*60 + 5*70 + 5*75) / 100
//
//	= (2125 + 1600 + 1400 + 1125 + 600 + 350 + 375) / 100
//	= 7575 / 100 = 75.75
func TestComputeDimensionProject_typicalEntry(t *testing.T) {
	result := ComputeDimensionProject(85, 80, 70, 75, 60, 70, 75)
	if !almostEqual(result, 75.75, 0.01) {
		t.Errorf("expected 75.75, got %v", result)
	}
}

// TestComputeDimensionTeam_allZero tests DT with all zero contributions.
func TestComputeDimensionTeam_allZero(t *testing.T) {
	result := ComputeDimensionTeam(0, 0, 0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeDimensionTeam_allMax tests DT with all contributions at max (100).
func TestComputeDimensionTeam_allMax(t *testing.T) {
	result := ComputeDimensionTeam(100, 100, 100, 100, 100, 100)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeDimensionTeam_typicalEntry tests DT with typical mixed values.
// DT weights per PRD 5.4: MentoringN 25, PositiveN 20, CriticalN 20, MoraleN 15, DeliveryN 10, OvertimeN 10
// Test with: Mentoring=80, Positive=70, Critical=75, Morale=65, Delivery=85, Overtime=60
// DT = (25*80 + 20*70 + 20*75 + 15*65 + 10*85 + 10*60) / 100
//
//	= (2000 + 1400 + 1500 + 975 + 850 + 600) / 100
//	= 7325 / 100 = 73.25
func TestComputeDimensionTeam_typicalEntry(t *testing.T) {
	result := ComputeDimensionTeam(80, 70, 75, 65, 85, 60)
	if !almostEqual(result, 73.25, 0.01) {
		t.Errorf("expected 73.25, got %v", result)
	}
}

// TestComputeDimensionOrg_allZero tests DO with all zero contributions.
func TestComputeDimensionOrg_allZero(t *testing.T) {
	result := ComputeDimensionOrg(0, 0, 0, 0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeDimensionOrg_allMax tests DO with all contributions at max (100).
func TestComputeDimensionOrg_allMax(t *testing.T) {
	result := ComputeDimensionOrg(100, 100, 100, 100, 100, 100, 100)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeDimensionOrg_typicalEntry tests DO with typical mixed values.
// DO weights per PRD 5.4: MarginN 25, CSATN 20, BillabilityN 15, DeliveryN 15, MentoringN 10, PositiveN 10, EvidenceN 5
// Test with: Margin=80, CSAT=75, Billability=70, Delivery=85, Mentoring=65, Positive=70, Evidence=60
// DO = (25*80 + 20*75 + 15*70 + 15*85 + 10*65 + 10*70 + 5*60) / 100
//
//	= (2000 + 1500 + 1050 + 1275 + 650 + 700 + 300) / 100
//	= 7475 / 100 = 74.75
func TestComputeDimensionOrg_typicalEntry(t *testing.T) {
	result := ComputeDimensionOrg(80, 75, 70, 85, 65, 70, 60)
	if !almostEqual(result, 74.75, 0.01) {
		t.Errorf("expected 74.75, got %v", result)
	}
}

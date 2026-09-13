package scoring

import (
	"testing"
)

// TestNormalizeMorale_zeroScore tests that NormalizeMorale returns 0 when morale score is 0.
// PRD 5.2: MoraleN = (Morale / 5) * 100
// When Morale = 0, MoraleN = 0
func TestNormalizeMorale_zeroScore(t *testing.T) {
	result := NormalizeMorale(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNormalizeMorale_maxScore tests that NormalizeMorale returns 100 when morale score is 5.
func TestNormalizeMorale_maxScore(t *testing.T) {
	result := NormalizeMorale(5)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestNormalizeMorale_midRange tests that NormalizeMorale returns 60 when morale score is 3.
func TestNormalizeMorale_midRange(t *testing.T) {
	result := NormalizeMorale(3)
	if !almostEqual(result, 60.0, 0.01) {
		t.Errorf("expected 60.0, got %v", result)
	}
}

// Billability tests
// TestNormalizeBillability_optimum75Percent tests that billability is 100 at the optimal 75%.
// PRD 5.2: BillabilityN = clamp(100 - (abs(75 - 75) / 15) * 100, 0, 100) = 100
func TestNormalizeBillability_optimum75Percent(t *testing.T) {
	result := NormalizeBillability(75)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0 at 75%%, got %v", result)
	}
}

// TestNormalizeBillability_veryLow tests billability at 0% (far below optimal).
// Deviation = 75, reduced score = 100 - (75/15)*100 = 100 - 500 = -400, clamped to 0
func TestNormalizeBillability_veryLow(t *testing.T) {
	result := NormalizeBillability(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0 at 0%%, got %v", result)
	}
}

// TestNormalizeBillability_veryHigh tests billability at 100% (far above optimal).
// Deviation = 25, reduced score = 100 - (25/15)*100 = 100 - 166.67 = -66.67, clamped to 0
func TestNormalizeBillability_veryHigh(t *testing.T) {
	result := NormalizeBillability(100)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0 (clamped) at 100%%, got %v", result)
	}
}

// CSAT tests
// TestNormalizeCSAT_minScore1 tests that CSAT=1 normalizes to 0.
// PRD 5.2: CSATN = ((1 - 1) / 4) * 100 = 0
func TestNormalizeCSAT_minScore1(t *testing.T) {
	result := NormalizeCSAT(1)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNormalizeCSAT_maxScore5 tests that CSAT=5 normalizes to 100.
// PRD 5.2: CSATN = ((5 - 1) / 4) * 100 = 100
func TestNormalizeCSAT_maxScore5(t *testing.T) {
	result := NormalizeCSAT(5)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// Margin tests
// TestNormalizeMargin_belowMin tests margin at -30 (below -20 minimum).
// norm01(-30, -20, 60) = clamp((-30 - (-20)) / (60 - (-20)), 0, 1) = clamp(-10/80, 0, 1) = 0
// Result = 0 * 100 = 0
func TestNormalizeMargin_belowMin(t *testing.T) {
	result := NormalizeMargin(-30)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0 (clamped), got %v", result)
	}
}

// TestNormalizeMargin_aboveMax tests margin at 80 (above +60 maximum).
// norm01(80, -20, 60) = clamp((80 - (-20)) / (60 - (-20)), 0, 1) = clamp(100/80, 0, 1) = 1
// Result = 1 * 100 = 100
func TestNormalizeMargin_aboveMax(t *testing.T) {
	result := NormalizeMargin(80)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0 (clamped), got %v", result)
	}
}

// TestNormalizeMargin_within20to60 tests margin at 20 (midpoint of the expected range).
// norm01(20, -20, 60) = (20 - (-20)) / (60 - (-20)) = 40/80 = 0.5
// Result = 0.5 * 100 = 50
func TestNormalizeMargin_within20to60(t *testing.T) {
	result := NormalizeMargin(20)
	if !almostEqual(result, 50.0, 0.01) {
		t.Errorf("expected 50.0, got %v", result)
	}
}

// Positive Feedback tests
// TestNormalizePositive_zeroFeedback tests zero positive feedback.
// PositiveN = 0 / 8 * 100 = 0
func TestNormalizePositive_zeroFeedback(t *testing.T) {
	result := NormalizePositive(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNormalizePositive_cappedAt8 tests that 8+ positive feedback counts as 100.
// PositiveN = min(8, 8) / 8 * 100 = 100
func TestNormalizePositive_cappedAt8(t *testing.T) {
	result := NormalizePositive(8)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// Critical Feedback tests
// TestNormalizeCritical_zeroFeedback tests zero critical feedback.
// CriticalN = 100 - (0 / 6 * 100) = 100
func TestNormalizeCritical_zeroFeedback(t *testing.T) {
	result := NormalizeCritical(0)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestNormalizeCritical_cappedAt6 tests that 6+ critical feedback counts as 0.
// CriticalN = 100 - (min(6, 6) / 6 * 100) = 100 - 100 = 0
func TestNormalizeCritical_cappedAt6(t *testing.T) {
	result := NormalizeCritical(6)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// Overtime tests
// TestNormalizeOvertime_zeroHours tests zero overtime.
// OvertimeN = 100 - (0 / 30 * 100) = 100
func TestNormalizeOvertime_zeroHours(t *testing.T) {
	result := NormalizeOvertime(0)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestNormalizeOvertime_cappedAt30 tests that 30+ hours counts as 0.
// OvertimeN = 100 - (min(30, 30) / 30 * 100) = 100 - 100 = 0
func TestNormalizeOvertime_cappedAt30(t *testing.T) {
	result := NormalizeOvertime(30)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// Delivery tests
// TestNormalizeDelivery_clampsTo0_100 tests that delivery is clamped to [0, 100].
func TestNormalizeDelivery_clampsTo0_100(t *testing.T) {
	// Valid case
	result := NormalizeDelivery(85)
	if !almostEqual(result, 85.0, 0.01) {
		t.Errorf("expected 85.0, got %v", result)
	}
	// Below bounds
	resultLow := NormalizeDelivery(-10)
	if !almostEqual(resultLow, 0.0, 0.01) {
		t.Errorf("expected 0.0 (clamped), got %v", resultLow)
	}
	// Above bounds
	resultHigh := NormalizeDelivery(150)
	if !almostEqual(resultHigh, 100.0, 0.01) {
		t.Errorf("expected 100.0 (clamped), got %v", resultHigh)
	}
}

// Mentoring tests
// TestNormalizeMentoring_zeroHours tests zero mentoring hours.
// MentoringN = 0 / 12 * 100 = 0
func TestNormalizeMentoring_zeroHours(t *testing.T) {
	result := NormalizeMentoring(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNormalizeMentoring_cappedAt12 tests that 12+ hours counts as 100.
// MentoringN = min(12, 12) / 12 * 100 = 100
func TestNormalizeMentoring_cappedAt12(t *testing.T) {
	result := NormalizeMentoring(12)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// Evidence tests
// TestNormalizeEvidence_zeroNotes tests zero evidence notes.
// EvidenceN = 0 / 6 * 100 = 0
func TestNormalizeEvidence_zeroNotes(t *testing.T) {
	result := NormalizeEvidence(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestNormalizeEvidence_cappedAt6 tests that 6+ notes counts as 100.
// EvidenceN = min(6, 6) / 6 * 100 = 100
func TestNormalizeEvidence_cappedAt6(t *testing.T) {
	result := NormalizeEvidence(6)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

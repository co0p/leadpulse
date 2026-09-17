package scoring

import (
	"testing"
)

// TestComputeTII_allZeroDimensions tests TII with all dimension scores at zero.
// PRD 5.5: TII = 0.20*DG + 0.35*DP + 0.25*DT + 0.20*DO
func TestComputeTII_allZeroDimensions(t *testing.T) {
	result := ComputeTII(0, 0, 0, 0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeTII_allMaxDimensions tests TII with all dimension scores at max (100).
func TestComputeTII_allMaxDimensions(t *testing.T) {
	result := ComputeTII(100, 100, 100, 100)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeTII_typicalMixedDimensions tests TII with typical mixed dimension values.
// DG=70, DP=75, DT=73, DO=74
// TII = 0.20*70 + 0.35*75 + 0.25*73 + 0.20*74 = 14 + 26.25 + 18.25 + 14.8 = 73.3
func TestComputeTII_typicalMixedDimensions(t *testing.T) {
	result := ComputeTII(70, 75, 73, 74)
	if !almostEqual(result, 73.3, 0.01) {
		t.Errorf("expected 73.3, got %v", result)
	}
}

// TestComputeCompleteness_zeroFields tests completeness with zero fields filled.
// PRD 4.3: CompletenessPct = (FilledRequiredFields / 11) * 100
func TestComputeCompleteness_zeroFields(t *testing.T) {
	result := ComputeCompleteness(0)
	if !almostEqual(result, 0.0, 0.01) {
		t.Errorf("expected 0.0, got %v", result)
	}
}

// TestComputeCompleteness_allFieldsFilled tests completeness with all 11 fields filled.
func TestComputeCompleteness_allFieldsFilled(t *testing.T) {
	result := ComputeCompleteness(11)
	if !almostEqual(result, 100.0, 0.01) {
		t.Errorf("expected 100.0, got %v", result)
	}
}

// TestComputeCompleteness_partiallyFilled tests completeness with 7 of 11 fields filled.
// CompletenessPct = (7 / 11) * 100 ≈ 63.6%
func TestComputeCompleteness_partiallyFilled(t *testing.T) {
	result := ComputeCompleteness(7)
	if !almostEqual(result, 63.636, 0.1) {
		t.Errorf("expected 63.636, got %v", result)
	}
}

// TestComputeScores_endToEndTypicalMember tests a full scoring pipeline with typical input.
// Morale=4, Billability=75%, CSAT=4, Margin=20%, Positive=2, Critical=1,
// Overtime=10, Delivery=90%, Mentoring=4, Evidence=2
// All impact ratings = 3 (moderate importance across all layers)
// Expected: normalized scores, dimension scores, TII computed correctly
func TestComputeScores_endToEndTypicalMember(t *testing.T) {
	// Compute normalized scores
	moraleN := NormalizeMorale(4)            // 80
	billabilityN := NormalizeBillability(75) // 100
	csatN := NormalizeCSAT(4)                // 75
	marginN := NormalizeMargin(20)           // 50
	positiveN := NormalizePositive(2)        // 25
	criticalN := NormalizeCritical(1)        // 83.33
	overtimeN := NormalizeOvertime(10)       // 66.67
	deliveryN := NormalizeDelivery(90)       // 90
	mentoringN := NormalizeMentoring(4)      // 33.33
	evidenceN := NormalizeEvidence(2)        // 33.33

	// Compute impact-weighted scores (all impact ratings = 3)
	moraleImpact := ComputeImpactWeighted(3, 3, 3, 3)      // 60
	billabilityImpact := ComputeImpactWeighted(3, 3, 3, 3) // 60
	csatImpact := ComputeImpactWeighted(3, 3, 3, 3)        // 60
	marginImpact := ComputeImpactWeighted(3, 3, 3, 3)      // 60
	positiveImpact := ComputeImpactWeighted(3, 3, 3, 3)    // 60
	criticalImpact := ComputeImpactWeighted(3, 3, 3, 3)    // 60
	overtimeImpact := ComputeImpactWeighted(3, 3, 3, 3)    // 60
	deliveryImpact := ComputeImpactWeighted(3, 3, 3, 3)    // 60
	mentoringImpact := ComputeImpactWeighted(3, 3, 3, 3)   // 60
	evidenceImpact := ComputeImpactWeighted(3, 3, 3, 3)    // 60

	// Compute contribution scores
	moraleC := (moraleN * moraleImpact) / 100                // 48
	billabilityC := (billabilityN * billabilityImpact) / 100 // 60
	csatC := (csatN * csatImpact) / 100                      // 45
	marginC := (marginN * marginImpact) / 100                // 30
	positiveC := (positiveN * positiveImpact) / 100          // 15
	criticalC := (criticalN * criticalImpact) / 100          // 50
	overtimeC := (overtimeN * overtimeImpact) / 100          // 40
	deliveryC := (deliveryN * deliveryImpact) / 100          // 54
	mentoringC := (mentoringN * mentoringImpact) / 100       // 20
	evidenceC := (evidenceN * evidenceImpact) / 100          // 20

	// Compute dimension scores
	dg := ComputeDimensionGrowth(moraleC, criticalC, positiveC, mentoringC, deliveryC, billabilityC, overtimeC, evidenceC)
	dp := ComputeDimensionProject(deliveryC, csatC, marginC, billabilityC, criticalC, positiveC, moraleC)
	dt := ComputeDimensionTeam(mentoringC, positiveC, criticalC, moraleC, deliveryC, overtimeC)
	do := ComputeDimensionOrg(marginC, csatC, billabilityC, deliveryC, mentoringC, positiveC, evidenceC)

	// Compute TII
	tii := ComputeTII(dg, dp, dt, do)

	// Verify TII is in a reasonable range (not 0 or 100, but somewhere in middle)
	if tii < 30.0 || tii > 70.0 {
		t.Errorf("expected TII in range [30, 70] for typical member, got %v", tii)
	}

	// Compute completeness (10 fields filled + at least one impact rating = 11)
	completeness := ComputeCompleteness(11)
	if !almostEqual(completeness, 100.0, 0.01) {
		t.Errorf("expected completeness 100.0, got %v", completeness)
	}
}

// TestComputeScores_endToEndStrongPerformer tests scoring for a high-performing member.
// High morale, good billability, high CSAT, positive margin, lots of positive feedback,
// minimal critical feedback, reasonable overtime, excellent delivery, mentoring hours.
func TestComputeScores_endToEndStrongPerformer(t *testing.T) {
	// Normalized scores
	moraleN := NormalizeMorale(5)            // 100
	billabilityN := NormalizeBillability(75) // 100
	csatN := NormalizeCSAT(5)                // 100
	marginN := NormalizeMargin(50)           // 100
	positiveN := NormalizePositive(8)        // 100
	criticalN := NormalizeCritical(0)        // 100
	overtimeN := NormalizeOvertime(5)        // 83.33
	deliveryN := NormalizeDelivery(95)       // 95
	mentoringN := NormalizeMentoring(12)     // 100
	evidenceN := NormalizeEvidence(6)        // 100

	// Impact-weighted (all ratings = 5, max impact)
	maxImpact := ComputeImpactWeighted(5, 5, 5, 5) // 100

	// Contribution scores (all at max)
	moraleC := (moraleN * maxImpact) / 100           // 100
	billabilityC := (billabilityN * maxImpact) / 100 // 100
	csatC := (csatN * maxImpact) / 100               // 100
	marginC := (marginN * maxImpact) / 100           // 100
	positiveC := (positiveN * maxImpact) / 100       // 100
	criticalC := (criticalN * maxImpact) / 100       // 100
	overtimeC := (overtimeN * maxImpact) / 100       // 83.33
	deliveryC := (deliveryN * maxImpact) / 100       // 95
	mentoringC := (mentoringN * maxImpact) / 100     // 100
	evidenceC := (evidenceN * maxImpact) / 100       // 100

	// Dimension scores
	dg := ComputeDimensionGrowth(moraleC, criticalC, positiveC, mentoringC, deliveryC, billabilityC, overtimeC, evidenceC)
	dp := ComputeDimensionProject(deliveryC, csatC, marginC, billabilityC, criticalC, positiveC, moraleC)
	dt := ComputeDimensionTeam(mentoringC, positiveC, criticalC, moraleC, deliveryC, overtimeC)
	do := ComputeDimensionOrg(marginC, csatC, billabilityC, deliveryC, mentoringC, positiveC, evidenceC)

	// TII
	tii := ComputeTII(dg, dp, dt, do)

	// Strong performer should have TII near or at max
	if tii < 90.0 {
		t.Errorf("expected TII > 90 for strong performer, got %v", tii)
	}
}

// TestComputeScores_endToEndStrugglingMember tests scoring for a struggling member.
// Low morale, low billability, low CSAT, negative margin, little positive feedback,
// multiple critical feedback items, high overtime, low delivery, no mentoring.
func TestComputeScores_endToEndStrugglingMember(t *testing.T) {
	// Normalized scores
	moraleN := NormalizeMorale(1)            // 20
	billabilityN := NormalizeBillability(50) // 66.67 (50% → 25 deviation → reduced)
	csatN := NormalizeCSAT(1)                // 0
	marginN := NormalizeMargin(-15)          // ~25
	positiveN := NormalizePositive(0)        // 0
	criticalN := NormalizeCritical(4)        // 33.33
	overtimeN := NormalizeOvertime(25)       // 16.67
	deliveryN := NormalizeDelivery(50)       // 50
	mentoringN := NormalizeMentoring(0)      // 0
	evidenceN := NormalizeEvidence(0)        // 0

	// Impact-weighted (all ratings = 1, minimal impact)
	minImpact := ComputeImpactWeighted(1, 1, 1, 1) // 20

	// Contribution scores (all low)
	moraleC := (moraleN * minImpact) / 100           // 4
	billabilityC := (billabilityN * minImpact) / 100 // 13.33
	csatC := (csatN * minImpact) / 100               // 0
	marginC := (marginN * minImpact) / 100           // 5
	positiveC := (positiveN * minImpact) / 100       // 0
	criticalC := (criticalN * minImpact) / 100       // 6.67
	overtimeC := (overtimeN * minImpact) / 100       // 3.33
	deliveryC := (deliveryN * minImpact) / 100       // 10
	mentoringC := (mentoringN * minImpact) / 100     // 0
	evidenceC := (evidenceN * minImpact) / 100       // 0

	// Dimension scores
	dg := ComputeDimensionGrowth(moraleC, criticalC, positiveC, mentoringC, deliveryC, billabilityC, overtimeC, evidenceC)
	dp := ComputeDimensionProject(deliveryC, csatC, marginC, billabilityC, criticalC, positiveC, moraleC)
	dt := ComputeDimensionTeam(mentoringC, positiveC, criticalC, moraleC, deliveryC, overtimeC)
	do := ComputeDimensionOrg(marginC, csatC, billabilityC, deliveryC, mentoringC, positiveC, evidenceC)

	// TII
	tii := ComputeTII(dg, dp, dt, do)

	// Struggling member should have low TII
	if tii > 30.0 {
		t.Errorf("expected TII < 30 for struggling member, got %v", tii)
	}
}

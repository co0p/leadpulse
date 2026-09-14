package domain

import (
	"testing"
)

// TestComputeImpactRating_moraleSignal tests that ScoringService.ComputeImpactRating
// correctly converts a morale signal and its impact rating into a weighted score.
// PRD 5.2: impact ratings (IG, IP, IT, IO) are each 0-5, weighted 20%, 35%, 25%, 20%
// and converted to 0-100 scale by multiplying by 20.
func TestComputeImpactRating_moraleSignal(t *testing.T) {
	// Arrange
	service := NewScoringService()
	morale, _ := NewMoraleScore(4)
	impactRating, _ := NewImpactRating(5, 4, 3, 2) // IG=5, IP=4, IT=3, IO=2

	// Act
	result := service.ComputeImpactRating(morale, impactRating)

	// Assert
	expectedWeightedScore := impactRating.ToWeightedScore()
	if result != expectedWeightedScore {
		t.Errorf("ComputeImpactRating() = %v, want %v", result, expectedWeightedScore)
	}
}

// TestComputeScores_validScores tests that ScoringService.ComputeScores
// returns all dimension scores within the valid range (0–100).
// This test verifies the structure and range constraints of a complete scoring result.
func TestComputeScores_validScores(t *testing.T) {
	// Arrange
	service := NewScoringService()

	// Act
	result := service.ComputeScores()

	// Assert
	// Verify all dimension scores are in range 0–100
	if result.DimensionScores.DG < 0 || result.DimensionScores.DG > 100 {
		t.Errorf("DimensionScores.DG out of range: %v", result.DimensionScores.DG)
	}
	if result.DimensionScores.DP < 0 || result.DimensionScores.DP > 100 {
		t.Errorf("DimensionScores.DP out of range: %v", result.DimensionScores.DP)
	}
	if result.DimensionScores.DT < 0 || result.DimensionScores.DT > 100 {
		t.Errorf("DimensionScores.DT out of range: %v", result.DimensionScores.DT)
	}
	if result.DimensionScores.DO < 0 || result.DimensionScores.DO > 100 {
		t.Errorf("DimensionScores.DO out of range: %v", result.DimensionScores.DO)
	}
}

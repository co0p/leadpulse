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

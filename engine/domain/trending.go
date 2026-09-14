package domain

import (
	"fmt"
	"math"
)

// TrendService encapsulates all trend calculation logic (moving averages, deltas, volatility).
type TrendService struct {
	entryRepo MonthlyEntryRepository
}

// NewTrendService creates a new TrendService with a monthly entry repository.
func NewTrendService(entryRepo MonthlyEntryRepository) *TrendService {
	return &TrendService{
		entryRepo: entryRepo,
	}
}

// ComputeMA3 computes the three-month moving average from a slice of TII scores.
// Expects at least 3 months of data.
func (t *TrendService) ComputeMA3(scores []float64) (float64, error) {
	if len(scores) < 3 {
		return 0, fmt.Errorf("insufficient data for MA3: need 3 months, got %d", len(scores))
	}
	return (scores[0] + scores[1] + scores[2]) / 3, nil
}

// ComputeDelta computes Delta1 (one-month change) and Delta3 (three-month change).
// Expects at least 4 months of data (to compute Delta3).
// Returns Delta1 and Delta3 for the last month.
func (t *TrendService) ComputeDelta(scores []float64) (float64, float64, error) {
	if len(scores) < 4 {
		return 0, 0, fmt.Errorf("insufficient data for Delta3: need 4 months, got %d", len(scores))
	}

	// Last month index
	n := len(scores) - 1

	// Delta1 = TII_t - TII_{t-1}
	delta1 := scores[n] - scores[n-1]

	// Delta3 = TII_t - TII_{t-3}
	delta3 := scores[n] - scores[n-3]

	return delta1, delta3, nil
}

// ComputeVolatility computes the three-month volatility (standard deviation).
// Expects at least 3 months of data.
func (t *TrendService) ComputeVolatility(scores []float64) (float64, error) {
	if len(scores) < 3 {
		return 0, fmt.Errorf("insufficient data for Vol3: need 3 months, got %d", len(scores))
	}

	// Take the last 3 months for volatility calculation
	n := len(scores)
	recent := scores[n-3 : n]

	// Calculate mean
	mean := (recent[0] + recent[1] + recent[2]) / 3

	// Calculate variance
	var sumSquaredDev float64
	for _, v := range recent {
		diff := v - mean
		sumSquaredDev += diff * diff
	}
	variance := sumSquaredDev / 3

	// Return standard deviation
	return math.Sqrt(variance), nil
}

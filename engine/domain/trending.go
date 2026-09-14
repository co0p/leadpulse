package domain

import "fmt"

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

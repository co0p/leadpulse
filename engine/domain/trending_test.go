package domain

import (
	"testing"
)

// TestComputeMA3_threeMonths tests that TrendService.ComputeMA3
// correctly computes the three-month moving average.
// PRD: MA3 = (TII_t + TII_{t-1} + TII_{t-2}) / 3
func TestComputeMA3_threeMonths(t *testing.T) {
	// Arrange
	repo := NewInMemoryMonthlyEntryRepository()
	service := NewTrendService(repo)

	// Mock monthly entry data with TII scores for 3 months
	memberID := TeamMemberID(1)

	// For now, test that ComputeMA3 returns a valid average
	scores := []float64{30.0, 50.0, 70.0} // 3 months of TII scores

	// Act
	ma3, err := service.ComputeMA3(scores)

	// Assert
	if err != nil {
		t.Errorf("ComputeMA3() returned error: %v", err)
	}

	expectedMA3 := (30.0 + 50.0 + 70.0) / 3 // 50.0
	if ma3 != expectedMA3 {
		t.Errorf("ComputeMA3() = %v, want %v", ma3, expectedMA3)
	}

	_ = memberID // silence unused variable
}

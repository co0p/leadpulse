package domain

import (
	"math"
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

// TestComputeMA3_shortHistory tests that TrendService.ComputeMA3
// returns an error when there are fewer than 3 months of history.
func TestComputeMA3_shortHistory(t *testing.T) {
	// Arrange
	repo := NewInMemoryMonthlyEntryRepository()
	service := NewTrendService(repo)

	shortScores := []float64{30.0, 50.0} // only 2 months of data

	// Act
	_, err := service.ComputeMA3(shortScores)

	// Assert
	if err == nil {
		t.Error("ComputeMA3() with 2 months should return an error")
	}
}

// TestComputeDelta_onlyAndThreeMonth tests that TrendService.ComputeDelta
// correctly computes Delta1 (one-month change) and Delta3 (three-month change).
// PRD: Delta1 = TII_t - TII_{t-1}, Delta3 = TII_t - TII_{t-3}
func TestComputeDelta_onlyAndThreeMonth(t *testing.T) {
	// Arrange
	repo := NewInMemoryMonthlyEntryRepository()
	service := NewTrendService(repo)

	// 4 months of scores: month 0=30, month 1=50, month 2=70, month 3=80
	scores := []float64{30.0, 50.0, 70.0, 80.0}

	// Act
	delta1, delta3, err := service.ComputeDelta(scores)

	// Assert
	if err != nil {
		t.Errorf("ComputeDelta() returned error: %v", err)
	}

	// For month 3: Delta1 = 80 - 70 = 10, Delta3 = 80 - 30 = 50
	expectedDelta1 := 10.0
	expectedDelta3 := 50.0

	if delta1 != expectedDelta1 {
		t.Errorf("Delta1 = %v, want %v", delta1, expectedDelta1)
	}
	if delta3 != expectedDelta3 {
		t.Errorf("Delta3 = %v, want %v", delta3, expectedDelta3)
	}
}

// TestComputeVolatility_threeMonth tests that TrendService.ComputeVolatility
// correctly computes the three-month volatility (standard deviation).
// PRD: Vol3 = stddev(TII_t, TII_{t-1}, TII_{t-2})
func TestComputeVolatility_threeMonth(t *testing.T) {
	// Arrange
	repo := NewInMemoryMonthlyEntryRepository()
	service := NewTrendService(repo)

	// Scores with known volatility: 50, 50, 50 should have Vol3=0
	scores := []float64{50.0, 50.0, 50.0}

	// Act
	vol, err := service.ComputeVolatility(scores)

	// Assert
	if err != nil {
		t.Errorf("ComputeVolatility() returned error: %v", err)
	}

	// Expected: stddev of [50, 50, 50] = 0
	expectedVol := 0.0
	if !almostEqual(vol, expectedVol, 0.0001) {
		t.Errorf("Vol3 = %v, want %v", vol, expectedVol)
	}

	// Test with varying data: 40, 50, 60
	varyingScores := []float64{40.0, 50.0, 60.0}
	vol2, err := service.ComputeVolatility(varyingScores)

	if err != nil {
		t.Errorf("ComputeVolatility() returned error: %v", err)
	}

	// Expected: stddev of [40, 50, 60] ≈ 8.165
	expectedVol2 := math.Sqrt(200.0 / 3) // sqrt(sum of squared deviations / n)
	if !almostEqual(vol2, expectedVol2, 0.01) {
		t.Errorf("Vol3(varying) = %v, want %v", vol2, expectedVol2)
	}
}

// almostEqual checks if two floats are approximately equal within tolerance.
func almostEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

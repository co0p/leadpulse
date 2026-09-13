package scoring

// Clamp returns x constrained to the range [min, max].
// Formula per PRD 5.1: clamp(x, min, max) = min(max(x, min), max)
// Used in normalization to ensure scores stay within valid ranges.
func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// Norm01 normalizes x to the range [0, 1] based on the interval [min, max].
// Formula per PRD 5.1: norm01(x, min, max) = clamp((x - min)/(max - min), 0, 1)
// Used in several normalization functions (e.g., MarginN) to map arbitrary ranges to 0-1.
func Norm01(x, min, max float64) float64 {
	if max == min {
		return 0 // Edge case: if range is zero, return 0
	}
	return Clamp((x-min)/(max-min), 0, 1)
}

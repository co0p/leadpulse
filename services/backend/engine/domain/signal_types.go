package domain

import "fmt"

// === RAW SIGNAL VALUE OBJECTS ===
// Each raw signal type enforces its domain constraints at construction time.

// MoraleScore represents morale self-assessment (0–5).
type MoraleScore int

// NewMoraleScore creates a validated morale score.
// Must be in the range 0–5.
func NewMoraleScore(v int) (MoraleScore, error) {
	if v < 0 || v > 5 {
		return 0, fmt.Errorf("morale score must be 0–5, got %d", v)
	}
	return MoraleScore(v), nil
}

// BillabilityPercent represents billability (0–100).
type BillabilityPercent int

// NewBillabilityPercent creates a validated billability percentage.
// Must be in the range 0–100.
func NewBillabilityPercent(v int) (BillabilityPercent, error) {
	if v < 0 || v > 100 {
		return 0, fmt.Errorf("billability must be 0–100, got %d", v)
	}
	return BillabilityPercent(v), nil
}

// CSATScore represents customer satisfaction (1–5).
type CSATScore int

// NewCSATScore creates a validated CSAT score.
// Must be in the range 1–5.
func NewCSATScore(v int) (CSATScore, error) {
	if v < 1 || v > 5 {
		return 0, fmt.Errorf("CSAT must be 1–5, got %d", v)
	}
	return CSATScore(v), nil
}

// NetMarginPercent represents net margin (-20 to +60).
type NetMarginPercent int

// NewNetMarginPercent creates a validated net margin percentage.
// Must be in the range -20 to +60.
func NewNetMarginPercent(v int) (NetMarginPercent, error) {
	if v < -20 || v > 60 {
		return 0, fmt.Errorf("net margin must be -20 to +60, got %d", v)
	}
	return NetMarginPercent(v), nil
}

// FeedbackCount represents positive or critical feedback count (>= 0).
type FeedbackCount int

// NewFeedbackCount creates a validated feedback count.
// Must be >= 0.
func NewFeedbackCount(v int) (FeedbackCount, error) {
	if v < 0 {
		return 0, fmt.Errorf("feedback count must be >= 0, got %d", v)
	}
	return FeedbackCount(v), nil
}

// Hours represents hours worked (>= 0).
type Hours int

// NewHours creates a validated hours value.
// Must be >= 0.
func NewHours(v int) (Hours, error) {
	if v < 0 {
		return 0, fmt.Errorf("hours must be >= 0, got %d", v)
	}
	return Hours(v), nil
}

// Percent represents a generic percentage (0–100).
type Percent int

// NewPercent creates a validated percentage.
// Must be in the range 0–100.
func NewPercent(v int) (Percent, error) {
	if v < 0 || v > 100 {
		return 0, fmt.Errorf("percent must be 0–100, got %d", v)
	}
	return Percent(v), nil
}

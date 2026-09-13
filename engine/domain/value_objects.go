package domain

import (
	"fmt"
	"math"
)

// === COMPOSITE VALUE OBJECTS ===

// FullName is a value object representing a person's name.
type FullName struct {
	First string
	Last  string
}

// NewFullName creates a validated full name.
// Both first and last name must be non-empty.
func NewFullName(first, last string) (FullName, error) {
	if first == "" || last == "" {
		return FullName{}, fmt.Errorf("first and last name are required")
	}
	return FullName{First: first, Last: last}, nil
}

// String returns the full name as "First Last".
func (fn FullName) String() string {
	return fn.First + " " + fn.Last
}

// ImpactRating represents one signal's impact across all four dimensions.
// Each dimension (Growth, Project, Team, Organization) is rated 0–5.
type ImpactRating struct {
	Growth       int // IG: 0–5
	Project      int // IP: 0–5
	Team         int // IT: 0–5
	Organization int // IO: 0–5
}

// NewImpactRating creates a validated impact rating.
// All dimensions must be in the range 0–5.
func NewImpactRating(growth, project, team, org int) (ImpactRating, error) {
	for name, v := range map[string]int{
		"growth": growth, "project": project, "team": team, "organization": org,
	} {
		if v < 0 || v > 5 {
			return ImpactRating{}, fmt.Errorf("impact rating %s must be 0–5, got %d", name, v)
		}
	}
	return ImpactRating{
		Growth:       growth,
		Project:      project,
		Team:         team,
		Organization: org,
	}, nil
}

// ToWeightedScore converts the impact rating to a 0–100 scale using layer weights.
// Weights: Growth 20%, Project 35%, Team 25%, Organization 20%.
func (ir ImpactRating) ToWeightedScore() float64 {
	return (0.20*float64(ir.Growth) + 0.35*float64(ir.Project) +
		0.25*float64(ir.Team) + 0.20*float64(ir.Organization)) * 20
}

// === COMPUTED SCORE VALUE OBJECTS ===

// NormalizedScore represents a signal after normalization (0–100).
type NormalizedScore float64

// NewNormalizedScore creates a validated normalized score.
// Must be in the range 0–100 and must be a valid number (not NaN or Inf).
func NewNormalizedScore(v float64) (NormalizedScore, error) {
	if v < 0 || v > 100 || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("normalized score must be 0–100, got %f", v)
	}
	return NormalizedScore(v), nil
}

// ContributionScore represents the contribution of a signal (normalized × impact) (0–100).
type ContributionScore float64

// NewContributionScore creates a validated contribution score.
// Computed as (normalized * impactWeighted) / 100.
func NewContributionScore(normalized float64, impactWeighted float64) (ContributionScore, error) {
	result := (normalized * impactWeighted) / 100
	if result < 0 || result > 100 || math.IsNaN(result) {
		return 0, fmt.Errorf("contribution score out of range: %f", result)
	}
	return ContributionScore(result), nil
}

// ValidateDimensionScores checks that all dimension scores are in valid range (0–100).
// This is a helper for aggregate creation; the existing DimensionScores struct in scoring.go
// is used for data transfer, but we enforce constraints via this validator.
func ValidateDimensionScores(dg, dp, dt, d_o float64) error {
	for name, v := range map[string]float64{
		"growth": dg, "project": dp, "team": dt, "org": d_o,
	} {
		if v < 0 || v > 100 || math.IsNaN(v) {
			return fmt.Errorf("dimension %s score out of range: %f", name, v)
		}
	}
	return nil
}

// ValidateTotalImpactIndex checks that TII is in valid range (0–100).
func ValidateTotalImpactIndex(v float64) error {
	if v < 0 || v > 100 || math.IsNaN(v) {
		return fmt.Errorf("total impact index must be 0–100, got %f", v)
	}
	return nil
}

// CompletenessPercent represents the percentage of required fields that are filled (0–100).
type CompletenessPercent float64

// NewCompletenessPercent creates a completeness percentage from a signal count.
// Computed as (signalCount / 11) * 100, capped at 100.
func NewCompletenessPercent(signalCount int) CompletenessPercent {
	pct := (float64(signalCount) / 11) * 100
	if pct > 100 {
		pct = 100
	}
	return CompletenessPercent(pct)
}

// ConfidenceScore represents the confidence in the computed scores (0–100).
type ConfidenceScore float64

// NewConfidenceScore creates a validated confidence score.
// Computed as 0.7*completeness + 0.3*evidenceScore.
func NewConfidenceScore(completeness float64, evidenceScore float64) (ConfidenceScore, error) {
	conf := 0.7*completeness + 0.3*evidenceScore
	if conf < 0 || conf > 100 || math.IsNaN(conf) {
		return 0, fmt.Errorf("confidence score out of range: %f", conf)
	}
	return ConfidenceScore(conf), nil
}

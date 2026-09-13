package domain

// ScoringResult contains all computed scores for a monthly entry.
// Produced by the scoring engine and persisted by the store layer.
type ScoringResult struct {
	// Normalized signal scores (0–100) — result of applying each signal's normalization formula
	NormalizedScores NormalizedScores

	// Impact-weighted scores (0–100) — each signal's impact rating converted to 0–100 scale
	ImpactWeightedScores ImpactWeightedScores

	// Contribution scores (0–100) — combined effect of normalized score and impact rating
	ContributionScores ContributionScores

	// Dimension scores (0–100) — weighted sum of contributions per impact layer
	DimensionScores DimensionScores

	// Overall metrics
	TII                 float64 // Total Impact Index (0–100)
	CompletenessPct     float64 // (FilledRequiredFields / 11) * 100
	Confidence          float64 // Quality confidence badge (pending PRD clarification)
}

// NormalizedScores holds the 0–100 normalized value for each raw signal.
// Each score is computed via a signal-specific normalization formula from PRD 5.2.
type NormalizedScores struct {
	MoraleN      float64
	BillabilityN float64
	CSATN        float64
	MarginN      float64
	PositiveN    float64
	CriticalN    float64
	OvertimeN    float64
	DeliveryN    float64
	MentoringN   float64
	EvidenceN    float64
}

// ImpactWeightedScores holds the 0–100 weighted impact rating for each signal.
// Each score combines the four impact ratings (IG, IP, IT, IO) using layer weights:
// Growth 20%, Project 35%, Team 25%, Organization 20%.
// Formula: ImpactWeighted100 = (0.20*IG + 0.35*IP + 0.25*IT + 0.20*IO) * 20
type ImpactWeightedScores struct {
	MoraleImpactWeighted    float64
	BillabilityImpactWeighted float64
	CSATImpactWeighted      float64
	MarginImpactWeighted    float64
	PositiveImpactWeighted  float64
	CriticalImpactWeighted  float64
	OvertimeImpactWeighted  float64
	DeliveryImpactWeighted  float64
	MentoringImpactWeighted float64
	EvidenceImpactWeighted  float64
}

// ContributionScores holds the combined signal contribution (normalized × impact) for each signal.
// Formula: Cs = (Ns * Is) / 100, where Ns = normalized score, Is = impact-weighted score.
type ContributionScores struct {
	MoraleC      float64
	BillabilityC float64
	CSATC        float64
	MarginC      float64
	PositiveC    float64
	CriticalC    float64
	OvertimeC    float64
	DeliveryC    float64
	MentoringC   float64
	EvidenceC    float64
}

// DimensionScores holds the four impact dimension scores (0–100 each).
// Each dimension is a weighted sum of contribution scores for signals relevant to that layer.
// Dimensions: Personal Growth (DG), Project Impact (DP), Team Impact (DT), Organization Impact (DO).
type DimensionScores struct {
	DG float64 // Personal Growth
	DP float64 // Project Impact
	DT float64 // Team Impact
	DO float64 // Organization Impact
}

// TrendMetrics holds historical trend calculations for a team member.
// Computed when 2+ months of history is available.
type TrendMetrics struct {
	Delta1 float64 // TII_t - TII_{t-1}, month-over-month change
	Delta3 float64 // TII_t - TII_{t-3}, three-month change
	MA3    float64 // (TII_t + TII_{t-1} + TII_{t-2}) / 3, three-month moving average
	Vol3   float64 // stddev(TII_t, TII_{t-1}, TII_{t-2}), three-month volatility

	// Dimension-level trends (same pattern as TII)
	DG_Delta3 float64
	DP_Delta3 float64
	DT_Delta3 float64
	DO_Delta3 float64
}

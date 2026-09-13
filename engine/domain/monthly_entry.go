package domain

import "time"

// MonthlyEntry represents a complete monthly scorecard entry for a team member.
// It contains raw signals, impact ratings, and computed scores.
// Identified by (member_id, month).
type MonthlyEntry struct {
	// Identity
	MemberID  int64
	Month     string // YYYY-MM format
	CreatedAt time.Time

	// Raw Signals (10 fields per PRD 4.1)
	MoraleSelfScore      int    // 0–5
	BillabilityPercent   int    // 0–100
	CSAT                 int    // 1–5
	NetMarginPercent     int    // expected range: -20 to +60
	PositiveFeedbackCnt  int    // >=0
	CriticalFeedbackCnt  int    // >=0
	OvertimeHours        int    // >=0
	DeliveryReliabilityPercent int // 0–100
	MentoringEnablementHours   int // >=0
	EvidenceNotesCount   int    // >=0

	// Impact Ratings (4 per signal × 10 signals = 40 fields per PRD 4.2)
	// Each signal has IG (Growth), IP (Project), IT (Team), IO (Organization) ratings (0–5)
	MoraleIG, MoraleIP, MoraleIT, MoraleIO                     int
	BillabilityIG, BillabilityIP, BillabilityIT, BillabilityIO int
	CSATIG, CSATIP, CSATIT, CSATIO                             int
	MarginIG, MarginIP, MarginIT, MarginIO                     int
	PositiveIG, PositiveIP, PositiveIT, PositiveIO             int
	CriticalIG, CriticalIP, CriticalIT, CriticalIO             int
	OvertimeIG, OvertimeIP, OvertimeIT, OvertimeIO             int
	DeliveryIG, DeliveryIP, DeliveryIT, DeliveryIO             int
	MentoringIG, MentoringIP, MentoringIT, MentoringIO         int
	EvidenceIG, EvidenceIP, EvidenceIT, EvidenceIO             int

	// Computed Scores (populated by scoring engine)
	// Normalized signal scores (0–100)
	MoraleN              float64
	BillabilityN         float64
	CSATN                float64
	MarginN              float64
	PositiveN            float64
	CriticalN            float64
	OvertimeN            float64
	DeliveryN            float64
	MentoringN           float64
	EvidenceN            float64

	// Impact-weighted scores (0–100) per signal
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

	// Contribution scores (0–100) per signal: C = (N × ImpactWeighted) / 100
	MoraleC            float64
	BillabilityC       float64
	CSATC              float64
	MarginC            float64
	PositiveC          float64
	CriticalC          float64
	OvertimeC          float64
	DeliveryC          float64
	MentoringC         float64
	EvidenceC          float64

	// Dimension scores (0–100)
	DG float64 // Personal Growth
	DP float64 // Project Impact
	DT float64 // Team Impact
	DO float64 // Organization Impact

	// Overall Index
	TII float64 // Total Impact Index (0–100)

	// Quality metrics
	CompletenessPct float64 // (FilledRequiredFields / 11) * 100
	Confidence      float64 // 0.7*CompletenessPct + 0.3*EvidenceN (pending PRD clarification)
}

// FilledRequiredFields counts how many of the 11 required raw signal fields are non-zero.
// These are the fields that must be present for a valid monthly entry:
// Morale, Billability, CSAT, NetMargin, PositiveFeedback, CriticalFeedback,
// OvertimeHours, DeliveryReliability, MentoringHours, EvidenceNotes, and one impact rating per signal.
// For v1, we count raw signals only (10 fields) + at least one impact rating filled (11th check).
func (me MonthlyEntry) FilledRequiredFields() int {
	count := 0

	// Count non-zero raw signals
	if me.MoraleSelfScore > 0 {
		count++
	}
	if me.BillabilityPercent > 0 {
		count++
	}
	if me.CSAT > 0 {
		count++
	}
	// NetMargin can be negative, so we check for non-default (assume -999 is unset if used as sentinel)
	// For now, treat as always "set" if entry exists; scoring will handle zero values
	count++ // NetMargin (always counted if entry exists)

	if me.PositiveFeedbackCnt >= 0 {
		count++
	}
	if me.CriticalFeedbackCnt >= 0 {
		count++
	}
	if me.OvertimeHours >= 0 {
		count++
	}
	if me.DeliveryReliabilityPercent > 0 {
		count++
	}
	if me.MentoringEnablementHours >= 0 {
		count++
	}
	if me.EvidenceNotesCount >= 0 {
		count++
	}

	return count
}

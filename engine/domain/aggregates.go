package domain

import (
	"fmt"
	"time"
)

// === AGGREGATE ROOTS ===
// Aggregates are clusters of entities and value objects that are treated as a single unit.
// All invariants and constraints are enforced within the aggregate.

// TeamMemberID is a value object for team member identity.
type TeamMemberID int64

// TeamMember is an aggregate root representing a team member with all their properties.
// All fields are private; access is controlled via methods.
type TeamMember struct {
	id            TeamMemberID
	name          FullName
	seniority     Seniority
	createdAt     time.Time
	deactivatedAt *time.Time
}

// NewTeamMember creates a new team member aggregate root.
// All constraints are enforced at construction.
func NewTeamMember(id int64, name FullName, seniority Seniority) (*TeamMember, error) {
	if id <= 0 {
		return nil, fmt.Errorf("team member ID must be positive")
	}

	if !seniority.Valid() {
		return nil, fmt.Errorf("invalid seniority: %s", seniority)
	}

	return &TeamMember{
		id:        TeamMemberID(id),
		name:      name,
		seniority: seniority,
		createdAt: time.Now(),
	}, nil
}

// ID returns the team member's immutable ID.
func (tm *TeamMember) ID() TeamMemberID {
	return tm.id
}

// Name returns the team member's full name.
func (tm *TeamMember) Name() FullName {
	return tm.name
}

// Seniority returns the team member's seniority level.
func (tm *TeamMember) Seniority() Seniority {
	return tm.seniority
}

// CreatedAt returns when the member was created.
func (tm *TeamMember) CreatedAt() time.Time {
	return tm.createdAt
}

// IsActive returns true if the member has not been deactivated.
func (tm *TeamMember) IsActive() bool {
	return tm.deactivatedAt == nil
}

// DeactivatedAt returns the deactivation timestamp, or nil if active.
func (tm *TeamMember) DeactivatedAt() *time.Time {
	return tm.deactivatedAt
}

// Deactivate marks the team member as inactive.
// Once deactivated, this cannot be undone (enforces business rule at aggregate).
func (tm *TeamMember) Deactivate() error {
	if !tm.IsActive() {
		return fmt.Errorf("team member is already deactivated")
	}
	now := time.Now()
	tm.deactivatedAt = &now
	return nil
}

// ChangeSeniority updates the seniority level.
// Encapsulation ensures only valid transitions are possible.
func (tm *TeamMember) ChangeSeniority(newSeniority Seniority) error {
	if !newSeniority.Valid() {
		return fmt.Errorf("invalid seniority: %s", newSeniority)
	}
	tm.seniority = newSeniority
	return nil
}

// === Monthly Entry Aggregate Root ===
// The monthly entry is a complete aggregate containing:
// - Raw signals (as a sub-aggregate)
// - Impact ratings for each signal
// - Computed scores (produced by scoring engine)

// MonthlyEntryID is the composite identifier for a monthly entry.
type MonthlyEntryID struct {
	MemberID TeamMemberID
	Month    string // YYYY-MM format
}

// MonthlyRawSignals is a sub-aggregate containing all raw input signals for a month.
type MonthlyRawSignals struct {
	Morale               int // 0–5
	Billability          int // 0–100
	CSAT                 int // 1–5
	NetMargin            int // -20 to +60
	PositiveFeedback     int // >= 0
	CriticalFeedback     int // >= 0
	OvertimeHours        int // >= 0
	DeliveryReliability  int // 0–100
	MentoringHours       int // >= 0
	EvidenceNotesCount   int // >= 0
}

// NewMonthlyRawSignals validates and creates a raw signals sub-aggregate.
// All signal values are validated against their type constraints.
func NewMonthlyRawSignals(
	morale, billability, csat, netMargin, positiveFeedback, criticalFeedback,
	overtimeHours, deliveryReliability, mentoringHours, evidenceNotesCount int,
) (MonthlyRawSignals, error) {
	// Validate each signal individually
	if morale < 0 || morale > 5 {
		return MonthlyRawSignals{}, fmt.Errorf("morale must be 0–5, got %d", morale)
	}
	if billability < 0 || billability > 100 {
		return MonthlyRawSignals{}, fmt.Errorf("billability must be 0–100, got %d", billability)
	}
	if csat < 1 || csat > 5 {
		return MonthlyRawSignals{}, fmt.Errorf("CSAT must be 1–5, got %d", csat)
	}
	if netMargin < -20 || netMargin > 60 {
		return MonthlyRawSignals{}, fmt.Errorf("net margin must be -20 to +60, got %d", netMargin)
	}
	if positiveFeedback < 0 {
		return MonthlyRawSignals{}, fmt.Errorf("positive feedback must be >= 0, got %d", positiveFeedback)
	}
	if criticalFeedback < 0 {
		return MonthlyRawSignals{}, fmt.Errorf("critical feedback must be >= 0, got %d", criticalFeedback)
	}
	if overtimeHours < 0 {
		return MonthlyRawSignals{}, fmt.Errorf("overtime hours must be >= 0, got %d", overtimeHours)
	}
	if deliveryReliability < 0 || deliveryReliability > 100 {
		return MonthlyRawSignals{}, fmt.Errorf("delivery reliability must be 0–100, got %d", deliveryReliability)
	}
	if mentoringHours < 0 {
		return MonthlyRawSignals{}, fmt.Errorf("mentoring hours must be >= 0, got %d", mentoringHours)
	}
	if evidenceNotesCount < 0 {
		return MonthlyRawSignals{}, fmt.Errorf("evidence notes count must be >= 0, got %d", evidenceNotesCount)
	}

	return MonthlyRawSignals{
		Morale:              morale,
		Billability:         billability,
		CSAT:                csat,
		NetMargin:           netMargin,
		PositiveFeedback:    positiveFeedback,
		CriticalFeedback:    criticalFeedback,
		OvertimeHours:       overtimeHours,
		DeliveryReliability: deliveryReliability,
		MentoringHours:      mentoringHours,
		EvidenceNotesCount:  evidenceNotesCount,
	}, nil
}

// FilledSignalCount returns how many non-zero signals are present.
func (mrs MonthlyRawSignals) FilledSignalCount() int {
	count := 0
	if mrs.Morale > 0 {
		count++
	}
	if mrs.Billability > 0 {
		count++
	}
	if mrs.CSAT > 0 {
		count++
	}
	// NetMargin: only count if explicitly set (non-zero or negative)
	// For this prototype, 0 is default/unset, so we don't count it
	if mrs.PositiveFeedback > 0 {
		count++
	}
	if mrs.CriticalFeedback > 0 {
		count++
	}
	if mrs.OvertimeHours > 0 {
		count++
	}
	if mrs.DeliveryReliability > 0 {
		count++
	}
	if mrs.MentoringHours > 0 {
		count++
	}
	if mrs.EvidenceNotesCount > 0 {
		count++
	}
	return count
}

// ComputedScores holds all derived scores (produced by scoring engine).
type ComputedScores struct {
	// Normalized signal scores (0–100)
	NormalizedScores NormalizedScores

	// Impact-weighted scores (0–100)
	ImpactWeightedScores ImpactWeightedScores

	// Contribution scores (0–100)
	ContributionScores ContributionScores

	// Dimension scores (0–100)
	DimensionScores DimensionScores

	// Overall metrics
	TII             float64
	CompletenessPct float64
	Confidence      float64
}

// MonthlyEntry is an aggregate root representing a complete monthly scorecard.
type MonthlyEntry struct {
	id           MonthlyEntryID
	signals      MonthlyRawSignals
	impacts      map[string]ImpactRating // one per signal
	computed     *ComputedScores         // nil until scored
	createdAt    time.Time
	computedAt   *time.Time // timestamp of last computation
}

// NewMonthlyEntry creates a new monthly entry aggregate.
// Enforces all constraints at construction.
func NewMonthlyEntry(memberID int64, month string, signals MonthlyRawSignals) (*MonthlyEntry, error) {
	if memberID <= 0 {
		return nil, fmt.Errorf("member ID must be positive")
	}

	// Validate month format (YYYY-MM)
	if len(month) != 7 || month[4] != '-' {
		return nil, fmt.Errorf("month must be in YYYY-MM format")
	}

	// At least one signal must be filled
	if signals.FilledSignalCount() < 1 {
		return nil, fmt.Errorf("monthly entry must have at least one signal filled")
	}

	return &MonthlyEntry{
		id: MonthlyEntryID{
			MemberID: TeamMemberID(memberID),
			Month:    month,
		},
		signals:   signals,
		impacts:   make(map[string]ImpactRating),
		createdAt: time.Now(),
	}, nil
}

// ID returns the monthly entry's immutable ID.
func (me *MonthlyEntry) ID() MonthlyEntryID {
	return me.id
}

// Signals returns the raw signals (read-only).
func (me *MonthlyEntry) Signals() MonthlyRawSignals {
	return me.signals
}

// SetImpactRating sets the impact rating for a signal.
// Enforces that the rating is valid before accepting it.
func (me *MonthlyEntry) SetImpactRating(signalName string, rating ImpactRating) error {
	if signalName == "" {
		return fmt.Errorf("signal name is required")
	}

	me.impacts[signalName] = rating
	return nil
}

// GetImpactRating retrieves the impact rating for a signal.
func (me *MonthlyEntry) GetImpactRating(signalName string) (ImpactRating, bool) {
	rating, ok := me.impacts[signalName]
	return rating, ok
}

// SetComputedScores updates the computed scores (called by scoring engine).
// All constraints are validated before acceptance.
func (me *MonthlyEntry) SetComputedScores(scores ComputedScores) error {
	// Validate dimension scores
	if err := ValidateDimensionScores(
		scores.DimensionScores.DG,
		scores.DimensionScores.DP,
		scores.DimensionScores.DT,
		scores.DimensionScores.DO,
	); err != nil {
		return fmt.Errorf("invalid dimension scores: %w", err)
	}

	// Validate TII
	if err := ValidateTotalImpactIndex(scores.TII); err != nil {
		return fmt.Errorf("invalid TII: %w", err)
	}

	me.computed = &scores
	now := time.Now()
	me.computedAt = &now
	return nil
}

// ComputedScores returns the computed scores, or nil if not yet computed.
func (me *MonthlyEntry) ComputedScores() *ComputedScores {
	return me.computed
}

// HasComputedScores returns true if the entry has been scored.
func (me *MonthlyEntry) HasComputedScores() bool {
	return me.computed != nil
}

// FilledSignalCount returns how many required signals are filled.
func (me *MonthlyEntry) FilledSignalCount() int {
	count := me.signals.FilledSignalCount()
	if len(me.impacts) > 0 {
		count++ // At least one impact rating filled
	}
	if count > 11 {
		count = 11
	}
	return count
}

// CreatedAt returns when the entry was created.
func (me *MonthlyEntry) CreatedAt() time.Time {
	return me.createdAt
}

// ComputedAt returns when the entry was last scored, or nil if not scored.
func (me *MonthlyEntry) ComputedAt() *time.Time {
	return me.computedAt
}

// SetCreatedAt is used by the repository layer to restore the creation timestamp
// from persistent storage (internal use only).
func (me *MonthlyEntry) SetCreatedAt(t time.Time) {
	me.createdAt = t
}

// SetComputedAt is used by the repository layer to restore the computed timestamp
// from persistent storage (internal use only).
func (me *MonthlyEntry) SetComputedAt(t time.Time) {
	me.computedAt = &t
}

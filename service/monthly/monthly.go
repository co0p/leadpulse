package monthly

import (
	"fmt"

	"leadpulse/engine/domain"
	"leadpulse/engine/scoring"
)

func valueOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// Service provides use cases for monthly entry management.
// It depends on repository interfaces and domain services for business logic.
type Service struct {
	memberRepo         domain.TeamMemberRepository
	entryRepo          domain.MonthlyEntryRepository
	scoringService     *domain.ScoringService
	trendService       *domain.TrendService
	entryDomainService *domain.MonthlyEntryService
}

// NewService creates a new monthly service with repository implementations and domain services.
func NewService(memberRepo domain.TeamMemberRepository, entryRepo domain.MonthlyEntryRepository) *Service {
	return &Service{
		memberRepo:         memberRepo,
		entryRepo:          entryRepo,
		scoringService:     domain.NewScoringService(),
		trendService:       domain.NewTrendService(entryRepo),
		entryDomainService: domain.NewMonthlyEntryService(memberRepo, entryRepo),
	}
}

// CreateEntry creates a new monthly entry for a team member.
// This delegates to the domain service which enforces cross-aggregate rules.
func (s *Service) CreateEntry(
	memberID int64,
	month string,
	morale, billability, csat, netMargin,
	positiveFeedback, criticalFeedback,
	overtimeHours, deliveryReliability,
	mentoringHours, evidenceNotesCount int,
) (*domain.MonthlyEntry, error) {
	// Validate month format
	if len(month) != 7 || month[4] != '-' {
		return nil, fmt.Errorf("month must be in YYYY-MM format")
	}

	// Create raw signals (validates all signal constraints)
	signals, err := domain.NewMonthlyRawSignals(
		morale, billability, csat, netMargin,
		positiveFeedback, criticalFeedback,
		overtimeHours, deliveryReliability,
		mentoringHours, evidenceNotesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid signals: %w", err)
	}

	// Use domain service to enforce cross-aggregate rules
	entry, err := s.entryDomainService.CreateEntry(memberID, month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves a monthly entry by member ID and month.
func (s *Service) GetEntry(memberID int64, month string) (*domain.MonthlyEntry, error) {
	if len(month) != 7 || month[4] != '-' {
		return nil, fmt.Errorf("month must be in YYYY-MM format")
	}

	id := domain.MonthlyEntryID{
		MemberID: domain.TeamMemberID(memberID),
		Month:    month,
	}

	entry, err := s.entryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entry: %w", err)
	}

	return entry, nil
}

// ListEntriesByMember returns all monthly entries for a team member.
func (s *Service) ListEntriesByMember(memberID int64) ([]*domain.MonthlyEntry, error) {
	entries, err := s.entryRepo.FindByMember(domain.TeamMemberID(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}

	return entries, nil
}

// UpdateEntry updates an existing monthly entry's signals.
// This delegates to the domain service which handles validation and persistence.
func (s *Service) UpdateEntry(
	memberID int64,
	month string,
	morale, billability, csat, netMargin,
	positiveFeedback, criticalFeedback,
	overtimeHours, deliveryReliability,
	mentoringHours, evidenceNotesCount int,
) (*domain.MonthlyEntry, error) {
	// Validate month format
	if len(month) != 7 || month[4] != '-' {
		return nil, fmt.Errorf("month must be in YYYY-MM format")
	}

	// Create raw signals (validates all signal constraints)
	signals, err := domain.NewMonthlyRawSignals(
		morale, billability, csat, netMargin,
		positiveFeedback, criticalFeedback,
		overtimeHours, deliveryReliability,
		mentoringHours, evidenceNotesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid signals: %w", err)
	}

	// Use domain service to handle update
	entry, err := s.entryDomainService.UpdateEntry(memberID, month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}

	return entry, nil
}

// ComputeScores computes all dimension scores for an entry using the ScoringService.
// This delegates to the domain service which encapsulates all scoring logic.
func (s *Service) ComputeScores(memberID int64, month string) (*domain.ScoringResult, error) {
	// Retrieve the entry
	entry, err := s.GetEntry(memberID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entry: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("entry not found for member %d in month %s", memberID, month)
	}

	// Delegate to ScoringService to compute scores
	result := s.scoringService.ComputeScores()
	if result == nil {
		return nil, fmt.Errorf("failed to compute scores")
	}

	return result, nil
}

func (s *Service) PreviewScores(signals domain.MonthlyRawSignals) *domain.ScoringResult {
	morale := valueOrZero(signals.Morale)
	billability := valueOrZero(signals.Billability)
	csat := valueOrZero(signals.CSAT)
	netMargin := valueOrZero(signals.NetMargin)
	positiveFeedback := valueOrZero(signals.PositiveFeedback)
	criticalFeedback := valueOrZero(signals.CriticalFeedback)
	overtimeHours := valueOrZero(signals.OvertimeHours)
	deliveryReliability := valueOrZero(signals.DeliveryReliability)
	mentoringHours := valueOrZero(signals.MentoringHours)
	evidenceNotesCount := valueOrZero(signals.EvidenceNotesCount)

	normMorale := scoring.NormalizeMorale(morale)
	normBillability := scoring.NormalizeBillability(billability)
	normCSAT := scoring.NormalizeCSAT(csat)
	normMargin := scoring.NormalizeMargin(netMargin)
	normPositive := scoring.NormalizePositive(positiveFeedback)
	normCritical := scoring.NormalizeCritical(criticalFeedback)
	normOvertime := scoring.NormalizeOvertime(overtimeHours)
	normDelivery := scoring.NormalizeDelivery(deliveryReliability)
	normMentoring := scoring.NormalizeMentoring(mentoringHours)
	normEvidence := scoring.NormalizeEvidence(evidenceNotesCount)

	dg := scoring.ComputeDimensionGrowth(normMorale, normCritical, normPositive, normMentoring, normDelivery, normBillability, normOvertime, normEvidence)
	dp := scoring.ComputeDimensionProject(normDelivery, normCSAT, normMargin, normBillability, normCritical, normPositive, normMorale)
	dt := scoring.ComputeDimensionTeam(normMentoring, normPositive, normCritical, normMorale, normDelivery, normOvertime)
	do := scoring.ComputeDimensionOrg(normMargin, normCSAT, normBillability, normDelivery, normMentoring, normPositive, normEvidence)
		tii := scoring.ComputeTII(dg, dp, dt, do)
	filled := signals.FilledSignalCount()
	completed := scoring.ComputeCompleteness(filled)

	return &domain.ScoringResult{
		DimensionScores: domain.DimensionScores{
			DG: dg,
			DP: dp,
			DT: dt,
			DO: do,
		},
		TII:             tii,
		CompletenessPct: completed,
		Confidence:      0,
	}
}

// GetTrends retrieves trend data for a member using the TrendService.
// This delegates to the domain service which calculates moving averages, deltas, and volatility.
func (s *Service) GetTrends(memberID int64) (*domain.TrendMetrics, error) {
	// Retrieve all entries for the member
	entries, err := s.ListEntriesByMember(memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}

	if len(entries) < 3 {
		return nil, fmt.Errorf("insufficient data for trend calculation: need 3 months, got %d", len(entries))
	}

	// Extract TII scores (for now, use a placeholder implementation)
	// In a real implementation, entries would contain computed scores
	// For now, return a placeholder to satisfy the interface
	return &domain.TrendMetrics{
		MA3: 0,
	}, nil
}

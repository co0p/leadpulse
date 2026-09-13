package monthly

import (
	"fmt"

	"leadpulse/engine/domain"
)

// Service provides use cases for monthly entry management.
// It depends on repository interfaces to abstract persistence concerns.
type Service struct {
	memberRepo domain.TeamMemberRepository
	entryRepo  domain.MonthlyEntryRepository
}

// NewService creates a new monthly service with repository implementations.
func NewService(memberRepo domain.TeamMemberRepository, entryRepo domain.MonthlyEntryRepository) *Service {
	return &Service{
		memberRepo: memberRepo,
		entryRepo:  entryRepo,
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
	domainSvc := domain.NewMonthlyEntryService(s.memberRepo, s.entryRepo)
	entry, err := domainSvc.CreateEntry(memberID, month, signals)
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
	domainSvc := domain.NewMonthlyEntryService(s.memberRepo, s.entryRepo)
	entry, err := domainSvc.UpdateEntry(memberID, month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}

	return entry, nil
}

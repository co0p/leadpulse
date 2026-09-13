package domain

import "fmt"

// === DOMAIN SERVICES ===
// Services that operate on aggregates but require repository access.
// Domain services enforce business rules that span multiple aggregates.

// MonthlyEntryService provides domain operations on monthly entries.
// It enforces cross-aggregate rules such as:
// - Only active members can have entries
// - One entry per member per month (uniqueness constraint)
type MonthlyEntryService struct {
	memberRepo TeamMemberRepository
	entryRepo  MonthlyEntryRepository
}

// NewMonthlyEntryService creates a new domain service.
func NewMonthlyEntryService(memberRepo TeamMemberRepository, entryRepo MonthlyEntryRepository) *MonthlyEntryService {
	return &MonthlyEntryService{
		memberRepo: memberRepo,
		entryRepo:  entryRepo,
	}
}

// CreateEntry creates and persists a new monthly entry for a member.
// This is a domain service because it enforces business rules across aggregates:
// - The member must exist and be active
// - Only one entry per member per month is allowed
func (s *MonthlyEntryService) CreateEntry(
	memberID int64,
	month string,
	signals MonthlyRawSignals,
) (*MonthlyEntry, error) {
	// Rule 1: Member must exist and be active
	member, err := s.memberRepo.FindByID(TeamMemberID(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to check member: %w", err)
	}
	if member == nil {
		return nil, fmt.Errorf("member not found: %d", memberID)
	}
	if !member.IsActive() {
		return nil, fmt.Errorf("cannot create entry for deactivated member")
	}

	// Create the entry aggregate
	entry, err := NewMonthlyEntry(memberID, month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}

	// Rule 2: One entry per member per month
	existing, err := s.entryRepo.FindByID(entry.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing entry: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("entry already exists for member %d in month %s", memberID, month)
	}

	// Persist the entry
	if err := s.entryRepo.Save(entry); err != nil {
		return nil, fmt.Errorf("failed to save entry: %w", err)
	}

	return entry, nil
}

// UpdateEntry updates an existing entry's signals and invalidates computed scores.
// This enforces the rule that computed scores are immutable until the signals change.
func (s *MonthlyEntryService) UpdateEntry(
	memberID int64,
	month string,
	signals MonthlyRawSignals,
) (*MonthlyEntry, error) {
	// Fetch the existing entry
	id := MonthlyEntryID{
		MemberID: TeamMemberID(memberID),
		Month:    month,
	}
	entry, err := s.entryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entry: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("entry not found for member %d in month %s", memberID, month)
	}

	// Create a new entry with updated signals
	// (This re-validates all constraints)
	updatedEntry, err := NewMonthlyEntry(memberID, month, signals)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}

	// Preserve impact ratings from the old entry
	for signalName, rating := range entry.impacts {
		_ = updatedEntry.SetImpactRating(signalName, rating)
	}

	// Persist the updated entry
	if err := s.entryRepo.Save(updatedEntry); err != nil {
		return nil, fmt.Errorf("failed to save updated entry: %w", err)
	}

	return updatedEntry, nil
}

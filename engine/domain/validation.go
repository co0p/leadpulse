package domain

import "fmt"

// ValidationService enforces team-wide rules and historical consistency.
type ValidationService struct {
	memberRepo TeamMemberRepository
	entryRepo  MonthlyEntryRepository
}

// NewValidationService creates a new ValidationService with repository dependencies.
func NewValidationService(memberRepo TeamMemberRepository, entryRepo MonthlyEntryRepository) *ValidationService {
	return &ValidationService{
		memberRepo: memberRepo,
		entryRepo:  entryRepo,
	}
}

// ValidateTeamMemberUniqueness checks that the given member's name is unique
// within all active team members.
func (v *ValidationService) ValidateTeamMemberUniqueness(member *TeamMember) error {
	active, err := v.memberRepo.FindActive()
	if err != nil {
		return fmt.Errorf("failed to retrieve active members: %w", err)
	}

	// Check if another active member has the same name
	for _, existing := range active {
		if existing.Name().String() == member.Name().String() {
			return fmt.Errorf("team member with name %s already exists", member.Name().String())
		}
	}

	return nil
}

// ValidateEntryUniqueness checks that no entry already exists for the given
// member and month combination.
func (v *ValidationService) ValidateEntryUniqueness(memberID TeamMemberID, month string) error {
	id := MonthlyEntryID{
		MemberID: memberID,
		Month:    month,
	}

	existing, err := v.entryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to check entry existence: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("entry for member %d in month %s already exists", memberID, month)
	}

	return nil
}

// ValidateEntry checks that an entry can be created for the given member.
// Verifies that the member exists and is active.
func (v *ValidationService) ValidateEntry(memberID TeamMemberID) error {
	member, err := v.memberRepo.FindByID(memberID)
	if err != nil {
		return fmt.Errorf("failed to retrieve member: %w", err)
	}

	if member == nil {
		return fmt.Errorf("member with ID %d does not exist", memberID)
	}

	if !member.IsActive() {
		return fmt.Errorf("member with ID %d is not active", memberID)
	}

	return nil
}

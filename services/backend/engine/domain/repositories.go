package domain

import (
	"fmt"
	"sync"
)

// === REPOSITORY INTERFACES ===
// Repositories abstract away storage concerns from the domain layer.
// The domain layer depends on these interfaces, not on concrete storage.

// TeamMemberRepository defines persistence operations for TeamMember aggregates.
type TeamMemberRepository interface {
	// Save persists a team member aggregate.
	Save(member *TeamMember) error

	// FindByID retrieves a team member by ID.
	// Returns nil if not found.
	FindByID(id TeamMemberID) (*TeamMember, error)

	// FindActive returns all active (non-deactivated) team members.
	FindActive() ([]*TeamMember, error)

	// Delete removes a team member.
	Delete(id TeamMemberID) error
}

// MonthlyEntryRepository defines persistence operations for MonthlyEntry aggregates.
type MonthlyEntryRepository interface {
	// Save persists a monthly entry aggregate.
	Save(entry *MonthlyEntry) error

	// FindByID retrieves a monthly entry by its composite ID.
	// Returns nil if not found.
	FindByID(id MonthlyEntryID) (*MonthlyEntry, error)

	// FindByMember retrieves all monthly entries for a member.
	FindByMember(memberID TeamMemberID) ([]*MonthlyEntry, error)

	// Delete removes a monthly entry.
	Delete(id MonthlyEntryID) error
}

// === IN-MEMORY REPOSITORY IMPLEMENTATIONS ===
// For testing and development, in-memory repositories provide fast, isolated storage.

// InMemoryTeamMemberRepository is a simple in-memory implementation.
type InMemoryTeamMemberRepository struct {
	mu      sync.RWMutex
	members map[TeamMemberID]*TeamMember
}

// NewInMemoryTeamMemberRepository creates a new in-memory team member repository.
func NewInMemoryTeamMemberRepository() *InMemoryTeamMemberRepository {
	return &InMemoryTeamMemberRepository{
		members: make(map[TeamMemberID]*TeamMember),
	}
}

// Save persists a team member.
func (r *InMemoryTeamMemberRepository) Save(member *TeamMember) error {
	if member == nil {
		return fmt.Errorf("cannot save nil team member")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.members[member.ID()] = member
	return nil
}

// FindByID retrieves a team member by ID.
func (r *InMemoryTeamMemberRepository) FindByID(id TeamMemberID) (*TeamMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	member, ok := r.members[id]
	if !ok {
		return nil, nil
	}
	return member, nil
}

// FindActive returns all active team members.
func (r *InMemoryTeamMemberRepository) FindActive() ([]*TeamMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var active []*TeamMember
	for _, member := range r.members {
		if member.IsActive() {
			active = append(active, member)
		}
	}
	return active, nil
}

// Delete removes a team member.
func (r *InMemoryTeamMemberRepository) Delete(id TeamMemberID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, id)
	return nil
}

// InMemoryMonthlyEntryRepository is a simple in-memory implementation.
type InMemoryMonthlyEntryRepository struct {
	mu      sync.RWMutex
	entries map[string]*MonthlyEntry // key: "memberID:month"
}

// NewInMemoryMonthlyEntryRepository creates a new in-memory monthly entry repository.
func NewInMemoryMonthlyEntryRepository() *InMemoryMonthlyEntryRepository {
	return &InMemoryMonthlyEntryRepository{
		entries: make(map[string]*MonthlyEntry),
	}
}

// keyFor generates a composite key from ID.
func (r *InMemoryMonthlyEntryRepository) keyFor(id MonthlyEntryID) string {
	return fmt.Sprintf("%d:%s", id.MemberID, id.Month)
}

// Save persists a monthly entry.
func (r *InMemoryMonthlyEntryRepository) Save(entry *MonthlyEntry) error {
	if entry == nil {
		return fmt.Errorf("cannot save nil monthly entry")
	}

	// Validate entry state
	if entry.FilledSignalCount() < 1 {
		return fmt.Errorf("monthly entry must have at least one signal filled")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.keyFor(entry.ID())
	r.entries[key] = entry
	return nil
}

// FindByID retrieves a monthly entry by its composite ID.
func (r *InMemoryMonthlyEntryRepository) FindByID(id MonthlyEntryID) (*MonthlyEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.entries[r.keyFor(id)]
	if !ok {
		return nil, nil
	}
	return entry, nil
}

// FindByMember retrieves all monthly entries for a member.
func (r *InMemoryMonthlyEntryRepository) FindByMember(memberID TeamMemberID) ([]*MonthlyEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var entries []*MonthlyEntry
	for _, entry := range r.entries {
		if entry.ID().MemberID == memberID {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

// Delete removes a monthly entry.
func (r *InMemoryMonthlyEntryRepository) Delete(id MonthlyEntryID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.entries, r.keyFor(id))
	return nil
}

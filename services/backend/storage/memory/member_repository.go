package memory

import (
	"fmt"
	"sync"

	"leadpulse/core/members"
)

// InMemoryMemberRepository is an in-memory implementation of the MemberRepository port.
// It stores members in a Go map and is suitable for testing.
// All operations are synchronous and thread-safe via mutex.
type InMemoryMemberRepository struct {
	mu      sync.RWMutex
	members map[int64]*members.TeamMember
	nextID  int64
}

// NewInMemoryMemberRepository creates a new in-memory member repository.
// The nextID is initialized to 1; each Save of a member without an ID will auto-increment.
func NewInMemoryMemberRepository() *InMemoryMemberRepository {
	return &InMemoryMemberRepository{
		members: make(map[int64]*members.TeamMember),
		nextID:  1,
	}
}

// Save persists a team member to memory.
// If the member ID is 0, a new ID is assigned.
// If the member ID already exists, the member is updated.
func (r *InMemoryMemberRepository) Save(member *members.TeamMember) error {
	if member == nil {
		return fmt.Errorf("cannot save nil team member")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := int64(member.ID())
	if id <= 0 {
		return fmt.Errorf("team member must have a positive ID")
	}

	// Store a copy to avoid external mutations
	r.members[id] = member
	return nil
}

// FindByID retrieves a team member by ID.
// Returns nil and an error if not found.
func (r *InMemoryMemberRepository) FindByID(id members.TeamMemberID) (*members.TeamMember, error) {
	if id <= 0 {
		return nil, fmt.Errorf("member ID must be positive")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	member, exists := r.members[int64(id)]
	if !exists {
		return nil, fmt.Errorf("member not found: id=%d", id)
	}
	return member, nil
}

// FindActive retrieves all active team members.
// Returns an empty slice (not nil) if no active members exist.
func (r *InMemoryMemberRepository) FindActive() ([]*members.TeamMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var active []*members.TeamMember
	for _, member := range r.members {
		if member.IsActive() {
			active = append(active, member)
		}
	}

	// Return empty slice instead of nil for consistency
	if active == nil {
		active = []*members.TeamMember{}
	}
	return active, nil
}

// Deactivate marks a team member as inactive by ID.
// Loads the member, calls Deactivate() on the aggregate to enforce business rules,
// then saves the updated member.
func (r *InMemoryMemberRepository) Deactivate(id members.TeamMemberID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	member, exists := r.members[int64(id)]
	if !exists {
		return fmt.Errorf("member not found: id=%d", id)
	}

	// Aggregate enforces business rule: cannot deactivate already-inactive member
	if err := member.Deactivate(); err != nil {
		return err
	}

	// Member is already in the map; mutation is reflected
	return nil
}

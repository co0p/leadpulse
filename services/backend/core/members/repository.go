package members

// MemberRepository is a port (interface) that defines the contract for member persistence.
// All storage implementations (SQLite, in-memory, etc.) must conform to this interface.
//
// The repository is responsible for:
// - Saving new and updated team members to persistent storage
// - Retrieving individual team members by ID
// - Retrieving all active team members
// - Marking team members as inactive (soft delete)
//
// The repository interface lives in the core package because it is infrastructure-agnostic.
// Implementations (adapters) live in storage/ packages and depend inward on this interface.
type MemberRepository interface {
	// Save persists a team member aggregate to storage.
	// If the member does not exist, it is created.
	// If the member already exists (by ID), it is updated.
	// Returns an error if persistence fails.
	Save(member *TeamMember) error

	// Update persists changes to an existing team member.
	// Alias for Save; used to clarify intent when modifying an existing member.
	// Returns an error if persistence fails.
	Update(member *TeamMember) error

	// FindByID retrieves a team member by their ID.
	// Returns the team member if found, or nil and an error if not found or if retrieval fails.
	FindByID(id TeamMemberID) (*TeamMember, error)

	// FindActive retrieves all active team members (those with deactivatedAt == nil).
	// Returns an empty slice (not nil) if no active members exist.
	// Returns an error if retrieval fails.
	FindActive() ([]*TeamMember, error)

	// FindInactive retrieves all inactive team members (those with deactivatedAt != nil).
	// Returns an empty slice (not nil) if no inactive members exist.
	// Returns an error if retrieval fails.
	FindInactive() ([]*TeamMember, error)

	// FindAll retrieves all team members regardless of active status.
	// Returns an empty slice (not nil) if no members exist.
	// Returns an error if retrieval fails.
	FindAll() ([]*TeamMember, error)

	// Deactivate marks a team member as inactive by ID.
	// Returns an error if the member is not found or if the operation fails.
	// Note: The aggregate enforces the business rule that already-inactive members cannot be deactivated again.
	Deactivate(id TeamMemberID) error
}

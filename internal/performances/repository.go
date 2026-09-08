package performances

// Repository defines the persistence contract for performance entries.
type Repository interface {
	// Add creates a new performance entry and returns it with an assigned ID.
	Add(input Entry) (Entry, error)

	// ListForEmployee returns all entries for an employee, ordered by period ascending.
	ListForEmployee(employeeID int64) ([]Entry, error)

	// ListForAllEmployees returns all entries for all employees, unordered.
	ListForAllEmployees() ([]Entry, error)
}

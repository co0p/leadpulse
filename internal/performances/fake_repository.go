package performances

// FakeRepository is an in-memory test double implementing Repository.
type FakeRepository struct {
	entries []Entry
	nextID  int64
}

// NewFakeRepository creates an empty fake repository.
func NewFakeRepository() *FakeRepository {
	return &FakeRepository{entries: make([]Entry, 0)}
}

// Add stores an entry in memory and assigns a synthetic ID.
func (f *FakeRepository) Add(input Entry) (Entry, error) {
	f.nextID++
	entry := input
	entry.ID = f.nextID
	f.entries = append(f.entries, entry)
	return entry, nil
}

// ListForEmployee returns entries for the requested employee.
func (f *FakeRepository) ListForEmployee(employeeID int64) ([]Entry, error) {
	result := make([]Entry, 0)
	for _, entry := range f.entries {
		if entry.EmployeeID == employeeID {
			result = append(result, entry)
		}
	}
	return result, nil
}

// ListForAllEmployees returns all stored entries.
func (f *FakeRepository) ListForAllEmployees() ([]Entry, error) {
	result := make([]Entry, 0, len(f.entries))
	result = append(result, f.entries...)
	return result, nil
}

// Reset clears all stored entries.
func (f *FakeRepository) Reset() {
	f.entries = make([]Entry, 0)
	f.nextID = 0
}

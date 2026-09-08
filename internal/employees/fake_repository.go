package employees

// FakeRepository is an in-memory test double implementing Repository.
type FakeRepository struct {
	employees []Employee
	nextID    int64
}

// NewFakeRepository creates an empty fake employee repository.
func NewFakeRepository() *FakeRepository {
	return &FakeRepository{employees: make([]Employee, 0)}
}

// Add stores an employee and assigns a synthetic ID.
func (f *FakeRepository) Add(input Input) (Employee, error) {
	f.nextID++
	emp := Employee{
		ID:         f.nextID,
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		Seniority:  input.Seniority,
		StartDate:  input.StartDate,
	}
	f.employees = append(f.employees, emp)
	return emp, nil
}

// List returns all stored employees.
func (f *FakeRepository) List() ([]Employee, error) {
	result := make([]Employee, 0, len(f.employees))
	result = append(result, f.employees...)
	return result, nil
}

// Remove deletes an employee by ID.
func (f *FakeRepository) Remove(id int64) error {
	filtered := make([]Employee, 0, len(f.employees))
	for _, emp := range f.employees {
		if emp.ID != id {
			filtered = append(filtered, emp)
		}
	}
	f.employees = filtered
	return nil
}

// Reset clears all stored employees.
func (f *FakeRepository) Reset() {
	f.employees = make([]Employee, 0)
	f.nextID = 0
}

// Compile-time assertion that FakeRepository implements Repository.
var _ Repository = (*FakeRepository)(nil)

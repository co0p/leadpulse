package employees

// Repository defines the persistence contract for employee records.
type Repository interface {
	Add(input Input) (Employee, error)
	List() ([]Employee, error)
	Remove(id int64) error
}

// Compile-time assertion that Store implements Repository.
var _ Repository = (*Store)(nil)

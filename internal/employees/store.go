package employees

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	SeniorityJunior    = "Junior"
	SeniorityMidlevel  = "Midlevel"
	SenioritySenior    = "Senior"
	SeniorityPrincipal = "Principal"
)

type Input struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Seniority  string `json:"seniority"`
	StartDate  string `json:"startDate"`
}

type Employee struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Seniority  string `json:"seniority"`
	StartDate  string `json:"startDate"`
}

type Store struct {
	db *sql.DB
}

func NewStore(databasePath string) (*Store, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		second_name TEXT NOT NULL,
		seniority TEXT NOT NULL,
		start_date TEXT NOT NULL
	)`)
	return err
}

func (s *Store) Add(input Input) (Employee, error) {
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.SecondName) == "" || strings.TrimSpace(input.StartDate) == "" {
		return Employee{}, fmt.Errorf("employee names and start date are required")
	}
	if !isValidSeniority(input.Seniority) {
		return Employee{}, fmt.Errorf("invalid seniority: %s", input.Seniority)
	}
	result, err := s.db.Exec(`INSERT INTO employees (first_name, second_name, seniority, start_date) VALUES (?, ?, ?, ?)`, input.FirstName, input.SecondName, input.Seniority, input.StartDate)
	if err != nil {
		return Employee{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Employee{}, err
	}
	return Employee{ID: id, FirstName: input.FirstName, SecondName: input.SecondName, Seniority: input.Seniority, StartDate: input.StartDate}, nil
}

func (s *Store) List() ([]Employee, error) {
	rows, err := s.db.Query(`SELECT id, first_name, second_name, seniority, start_date FROM employees ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Employee, 0)
	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.SecondName, &employee.Seniority, &employee.StartDate); err != nil {
			return nil, err
		}
		result = append(result, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) Remove(id int64) error {
	_, err := s.db.Exec(`DELETE FROM employees WHERE id = ?`, id)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}

func isValidSeniority(seniority string) bool {
	switch seniority {
	case SeniorityJunior, SeniorityMidlevel, SenioritySenior, SeniorityPrincipal:
		return true
	default:
		return false
	}
}

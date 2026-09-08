package performances

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store is the SQLite implementation of Repository.
type Store struct {
	db *sql.DB
}

// NewStore opens or creates the performance store backed by databasePath.
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
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS performance_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		employee_id INTEGER NOT NULL,
		year INTEGER NOT NULL,
		month INTEGER NOT NULL,
		morale REAL NOT NULL,
		execution REAL NOT NULL,
		impact REAL NOT NULL,
		growth REAL NOT NULL,
		culture REAL NOT NULL,
		projects TEXT NOT NULL,
		evidence_items TEXT NOT NULL,
		FOREIGN KEY (employee_id) REFERENCES employees(id),
		UNIQUE (employee_id, year, month)
	)`)
	return err
}

// Add creates a new performance entry.
func (s *Store) Add(input Entry) (Entry, error) {
	projectsJSON, err := json.Marshal(input.Projects)
	if err != nil {
		return Entry{}, err
	}
	evidenceJSON, err := json.Marshal(input.EvidenceItems)
	if err != nil {
		return Entry{}, err
	}
	result, err := s.db.Exec(
		`INSERT INTO performance_entries (employee_id, year, month, morale, execution, impact, growth, culture, projects, evidence_items)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.EmployeeID, input.Year, input.Month, input.Morale, input.Execution, input.Impact, input.Growth, input.Culture, projectsJSON, evidenceJSON,
	)
	if err != nil {
		return Entry{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Entry{}, err
	}
	return Entry{
		ID:            id,
		EmployeeID:    input.EmployeeID,
		Year:          input.Year,
		Month:         input.Month,
		Morale:        input.Morale,
		Execution:     input.Execution,
		Impact:        input.Impact,
		Growth:        input.Growth,
		Culture:       input.Culture,
		Projects:      input.Projects,
		EvidenceItems: input.EvidenceItems,
	}, nil
}

// ListForEmployee returns entries for one employee ordered by year and month ascending.
func (s *Store) ListForEmployee(employeeID int64) ([]Entry, error) {
	rows, err := s.db.Query(
		`SELECT id, employee_id, year, month, morale, execution, impact, growth, culture, projects, evidence_items
		FROM performance_entries WHERE employee_id = ? ORDER BY year ASC, month ASC`,
		employeeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEntries(rows)
}

// ListForAllEmployees returns every entry in the store.
func (s *Store) ListForAllEmployees() ([]Entry, error) {
	rows, err := s.db.Query(
		`SELECT id, employee_id, year, month, morale, execution, impact, growth, culture, projects, evidence_items
		FROM performance_entries ORDER BY year ASC, month ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEntries(rows)
}

func scanEntries(rows *sql.Rows) ([]Entry, error) {
	result := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		var projectsJSON, evidenceJSON string
		if err := rows.Scan(&entry.ID, &entry.EmployeeID, &entry.Year, &entry.Month, &entry.Morale, &entry.Execution, &entry.Impact, &entry.Growth, &entry.Culture, &projectsJSON, &evidenceJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(projectsJSON), &entry.Projects); err != nil {
			return nil, fmt.Errorf("unmarshal projects: %w", err)
		}
		if err := json.Unmarshal([]byte(evidenceJSON), &entry.EvidenceItems); err != nil {
			return nil, fmt.Errorf("unmarshal evidence items: %w", err)
		}
		result = append(result, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

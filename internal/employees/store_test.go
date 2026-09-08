package employees

import (
	"path/filepath"
	"testing"
)

func TestStoreAddsEmployeeWithApprovedFields(t *testing.T) {
	store := newTestStore(t)

	employee, err := store.Add(Input{FirstName: "Ada", SecondName: "Lovelace", Seniority: SenioritySenior, StartDate: "2020-01-15"})
	if err != nil {
		t.Fatal(err)
	}
	if employee.FirstName != "Ada" || employee.SecondName != "Lovelace" || employee.Seniority != SenioritySenior || employee.StartDate != "2020-01-15" {
		t.Fatalf("unexpected employee: %+v", employee)
	}
}

func TestStoreRejectsInvalidSeniority(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.Add(Input{FirstName: "Ada", SecondName: "Lovelace", Seniority: "Staff", StartDate: "2020-01-15"}); err == nil {
		t.Fatal("expected invalid seniority error")
	}
}

func TestStorePersistsAndRemovesEmployees(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "leadpulse.db")
	store, err := NewStore(databasePath)
	if err != nil {
		t.Fatal(err)
	}

	first, err := store.Add(Input{FirstName: "Ada", SecondName: "Lovelace", Seniority: SenioritySenior, StartDate: "2020-01-15"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Add(Input{FirstName: "Grace", SecondName: "Hopper", Seniority: SeniorityPrincipal, StartDate: "2018-03-01"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewStore(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	defer reloaded.Close()

	employees, err := reloaded.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(employees) != 2 {
		t.Fatalf("expected 2 employees after reload, got %d", len(employees))
	}
	if err := reloaded.Remove(first.ID); err != nil {
		t.Fatal(err)
	}
	employees, err = reloaded.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(employees) != 1 || employees[0].ID != second.ID {
		t.Fatalf("unexpected employees after removal: %+v", employees)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := NewStore(filepath.Join(t.TempDir(), "leadpulse.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

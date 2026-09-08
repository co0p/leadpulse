package performances

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_AddAndList(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	defer store.Close()

	entry := Entry{
		EmployeeID: 1,
		Year:       2026,
		Month:      8,
		Morale:     4,
		Execution:  4,
		Impact:     4,
		Growth:     4,
		Culture:    4,
		Projects:   []string{"Alpha"},
		EvidenceItems: []EvidenceItem{
			{Category: EvidenceCategoryProjectOutcome, Source: EvidenceSourceManualEntry, Description: "Delivered Alpha"},
		},
	}
	created, err := store.Add(entry)
	if err != nil {
		t.Fatalf("add entry: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected assigned id")
	}

	list, err := store.ListForEmployee(1)
	if err != nil {
		t.Fatalf("list entries: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
	if list[0].Morale != 4 {
		t.Fatalf("expected morale 4, got %f", list[0].Morale)
	}
}

func TestStore_UsesEnvironmentPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "env.db")
	os.Setenv("LEADPULSE_DATABASE_PATH", path)
	defer os.Unsetenv("LEADPULSE_DATABASE_PATH")

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	defer store.Close()

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("database file should exist: %v", err)
	}
}

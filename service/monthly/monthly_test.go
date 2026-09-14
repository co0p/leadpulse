package monthly

import (
	"testing"

	"leadpulse/engine/domain"
)

// TestCreateMonthlyEntryForActiveMember tests creating an entry for an active member.
func TestCreateMonthlyEntryForActiveMember(t *testing.T) {
	// Create in-memory repositories for testing
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member
	name, _ := domain.NewFullName("Alice", "Chen")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)
	memberRepo.Save(member)

	// Create the monthly service
	svc := NewService(memberRepo, entryRepo)

	// Create a monthly entry
	// Use pointer-style args to match new signature
	m, b, c, n := 3, 85, 4, 15
	pf, cf, ot := 5, 1, 4
	dr, mh, en := 90, 2, 8
	entry, err := svc.CreateEntry(
		1, "2024-10",
		&m, &b, &c, &n,
		&pf, &cf,
		&ot, &dr,
		&mh, &en,
	)

	if err != nil {
		t.Fatalf("CreateEntry failed: %v", err)
	}

	if entry == nil {
		t.Fatal("entry is nil")
	}

	if entry.ID().MemberID != 1 || entry.ID().Month != "2024-10" {
		t.Errorf("entry ID mismatch")
	}

	if entry.Signals().Morale == nil || *entry.Signals().Morale != 3 {
		t.Errorf("expected morale 3, got %v", entry.Signals().Morale)
	}
}

// TestCreateMonthlyEntryForInactiveMemberFails tests that entries cannot be created for deactivated members.
func TestCreateMonthlyEntryForInactiveMemberFails(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create and deactivate a team member
	name, _ := domain.NewFullName("Bob", "Smith")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	member.Deactivate()
	memberRepo.Save(member)

	// Create the monthly service
	svc := NewService(memberRepo, entryRepo)

	// Try to create an entry for the deactivated member
	// call using pointers where possible to match signature; use nil to test error path
	a, b, c, d, e, f, g, h, i, j := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	_, err := svc.CreateEntry(1, "2024-10", &a, &b, &c, &d, &e, &f, &g, &h, &i, &j)

	if err == nil {
		t.Fatal("expected error for deactivated member, got none")
	}

	if err.Error() != "failed to create entry: cannot create entry for deactivated member" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestCreateMonthlyEntryRejectsInvalidInput tests that out-of-range signal values are rejected.
func TestCreateMonthlyEntryRejectsInvalidInput(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Davis")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityJunior)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Use pointers to call CreateEntry; this path should still return validation error
	a, b, c, d, e, f, g, h, i, j := 6, 85, 4, 15, 5, 1, 4, 90, 2, 8
	_, err := svc.CreateEntry(1, "2024-10", &a, &b, &c, &d, &e, &f, &g, &h, &i, &j)

	if err == nil {
		t.Fatal("expected error for invalid morale, got none")
	}

	if err.Error() != "invalid signals: morale must be 0–5, got 6" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestCreateMonthlyEntryRejectsDuplicate tests that duplicate entries are rejected.
func TestCreateMonthlyEntryRejectsDuplicate(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create an active team member
	name, _ := domain.NewFullName("Charlie", "Davis")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityJunior)
	memberRepo.Save(member)

	// Create the monthly service
	svc := NewService(memberRepo, entryRepo)

	// Create the first entry using pointers
	a1, b1, c1, d1, e1, f1, g1, h1, i1, j1 := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	_, err := svc.CreateEntry(1, "2024-10", &a1, &b1, &c1, &d1, &e1, &f1, &g1, &h1, &i1, &j1)
	if err != nil {
		t.Fatalf("first CreateEntry failed: %v", err)
	}

	// Try to create a duplicate entry for the same member and month
	a2, b2, c2, d2, e2, f2, g2, h2, i2, j2 := 4, 80, 5, 20, 3, 0, 2, 95, 1, 6
	_, err = svc.CreateEntry(1, "2024-10", &a2, &b2, &c2, &d2, &e2, &f2, &g2, &h2, &i2, &j2)

	if err == nil {
		t.Fatal("expected error for duplicate entry, got none")
	}

	if err.Error() != "failed to create entry: entry already exists for member 1 in month 2024-10" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestGetMonthlyEntry tests retrieving a monthly entry by ID.
func TestGetMonthlyEntry(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member and entry
	name, _ := domain.NewFullName("Diana", "Evans")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Create an entry
	aX, bX, cX, dX, eX, fX, gX, hX, iX, jX := 4, 80, 5, 20, 3, 0, 2, 95, 1, 6
	created, _ := svc.CreateEntry(1, "2024-11", &aX, &bX, &cX, &dX, &eX, &fX, &gX, &hX, &iX, &jX)

	// Retrieve it
	retrieved, err := svc.GetEntry(1, "2024-11")

	if err != nil {
		t.Fatalf("GetEntry failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("entry not found")
	}

	if retrieved.ID().Month != created.ID().Month {
		t.Errorf("month mismatch")
	}

	if retrieved.Signals().Morale == nil || *retrieved.Signals().Morale != 4 {
		t.Errorf("expected morale 4, got %v", retrieved.Signals().Morale)
	}
}

// TestListEntriesByMember tests retrieving all entries for a member.
func TestListEntriesByMember(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member
	name, _ := domain.NewFullName("Eve", "Franklin")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Create multiple entries
	// Create entries using pointer args
	a1, b1, c1, d1, e1, f1, g1, h1, i1, j1 := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	a2, b2, c2, d2, e2, f2, g2, h2, i2, j2 := 4, 80, 5, 20, 3, 0, 2, 95, 1, 6
	a3, b3, c3, d3, e3, f3, g3, h3, i3, j3 := 5, 90, 4, 25, 4, 1, 3, 92, 3, 9
	svc.CreateEntry(1, "2024-09", &a1, &b1, &c1, &d1, &e1, &f1, &g1, &h1, &i1, &j1)
	svc.CreateEntry(1, "2024-10", &a2, &b2, &c2, &d2, &e2, &f2, &g2, &h2, &i2, &j2)
	svc.CreateEntry(1, "2024-11", &a3, &b3, &c3, &d3, &e3, &f3, &g3, &h3, &i3, &j3)

	// List entries
	entries, err := svc.ListEntriesByMember(1)

	if err != nil {
		t.Fatalf("ListEntriesByMember failed: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

// TestUpdateEntry tests updating an existing monthly entry.
func TestUpdateEntry(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member
	name, _ := domain.NewFullName("Frank", "Garcia")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Create an entry
	a4, b4, c4, d4, e4, f4, g4, h4, i4, j4 := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	svc.CreateEntry(1, "2024-10", &a4, &b4, &c4, &d4, &e4, &f4, &g4, &h4, &i4, &j4)

	// Update it with new signals
	u1, u2, u3, u4, u5, u6, u7, u8, u9, u10 := 5, 95, 5, 30, 8, 0, 5, 98, 4, 12
	updated, err := svc.UpdateEntry(1, "2024-10", &u1, &u2, &u3, &u4, &u5, &u6, &u7, &u8, &u9, &u10)

	if err != nil {
		t.Fatalf("UpdateEntry failed: %v", err)
	}

	if updated.Signals().Morale == nil || *updated.Signals().Morale != 5 {
		t.Errorf("expected morale 5, got %v", updated.Signals().Morale)
	}

	if updated.Signals().Billability == nil || *updated.Signals().Billability != 95 {
		t.Errorf("expected billability 95, got %v", updated.Signals().Billability)
	}
}

// TestComputeScores_delegatesToScoringService tests that ComputeScores delegates to ScoringService.
func TestService_PreviewScores_FullSignalsReturnComputedTII(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Grace", "Harris")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	signals, err := domain.NewMonthlyRawSignals(
		3, 85, 4, 15,
		5, 1, 4,
		90, 2, 8,
	)
	if err != nil {
		t.Fatalf("NewMonthlyRawSignals failed: %v", err)
	}

	result := svc.PreviewScores(signals)
	if result == nil {
		t.Fatal("preview result is nil")
	}
	if result.TII <= 0 {
		t.Fatalf("expected positive TII, got %v", result.TII)
	}
	if result.CompletenessPct != 90.9090909090909 {
		t.Fatalf("expected 90.9090909090909 completeness for 10 filled signals, got %v", result.CompletenessPct)
	}
}

func TestComputeScores_delegatesToScoringService(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member
	name, _ := domain.NewFullName("Grace", "Harris")
	member, _ := domain.NewTeamMember(1, name, domain.SenioritySenior)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Create an entry
	a5, b5, c5, d5, e5, f5, g5, h5, i5, j5 := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	svc.CreateEntry(1, "2024-10", &a5, &b5, &c5, &d5, &e5, &f5, &g5, &h5, &i5, &j5)

	// Compute scores
	result, err := svc.ComputeScores(1, "2024-10")

	if err != nil {
		t.Fatalf("ComputeScores failed: %v", err)
	}

	if result == nil {
		t.Fatal("scoring result is nil")
	}

	// Verify dimension scores are in valid range
	if result.DimensionScores.DG < 0 || result.DimensionScores.DG > 100 {
		t.Errorf("DG score out of range: %v", result.DimensionScores.DG)
	}
}

// TestGetTrends_delegatesToTrendService tests that GetTrends delegates to TrendService.
func TestGetTrends_delegatesToTrendService(t *testing.T) {
	// Create in-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a team member
	name, _ := domain.NewFullName("Henry", "Jackson")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	svc := NewService(memberRepo, entryRepo)

	// Create entries for 3 months (minimum for trend calculation)
	b1a, b1b, b1c, b1d, b1e, b1f, b1g, b1h, b1i, b1j := 3, 85, 4, 15, 5, 1, 4, 90, 2, 8
	b2a, b2b, b2c, b2d, b2e, b2f, b2g, b2h, b2i, b2j := 4, 80, 5, 20, 3, 0, 2, 95, 1, 6
	b3a, b3b, b3c, b3d, b3e, b3f, b3g, b3h, b3i, b3j := 5, 90, 4, 25, 4, 1, 3, 92, 3, 9
	svc.CreateEntry(1, "2024-09", &b1a, &b1b, &b1c, &b1d, &b1e, &b1f, &b1g, &b1h, &b1i, &b1j)
	svc.CreateEntry(1, "2024-10", &b2a, &b2b, &b2c, &b2d, &b2e, &b2f, &b2g, &b2h, &b2i, &b2j)
	svc.CreateEntry(1, "2024-11", &b3a, &b3b, &b3c, &b3d, &b3e, &b3f, &b3g, &b3h, &b3i, &b3j)

	// Get trends
	trends, err := svc.GetTrends(1)

	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}

	if trends == nil {
		t.Fatal("trend metrics is nil")
	}
}

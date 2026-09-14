package screens

import (
	"testing"
	"time"

	"leadpulse/engine/domain"
	monthlysvc "leadpulse/service/monthly"
)

// TestSaveEnableDisable verifies Save button enablement follows completeness threshold.
func TestSaveEnableDisable(t *testing.T) {
	// Build a minimal member service and monthly service using in-memory repos
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	// Add one member
	name, _ := domain.NewFullName("Test", "User")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	// Create the monthly service
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// We cannot easily introspect the fyne widget tree in unit tests without
	// the actual app driver. Instead, assert that PreviewScores behaves as
	// expected for empty vs partially-filled signals which drive Save enablement.

	// Empty signals -> completeness 0 -> Save should be disabled by policy
	emptySignals, _ := domain.NewMonthlyRawSignalsFromPointers(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	res := monthlyService.PreviewScores(emptySignals)
	if res == nil {
		t.Fatal("expected result")
	}
	if res.CompletenessPct != 0 {
		t.Fatalf("expected completeness 0 for empty signals, got %v", res.CompletenessPct)
	}

	// Enough fields (10 filled) -> completeness 90.909... -> Save enabled
	full, _ := domain.NewMonthlyRawSignals(3, 85, 4, 15, 5, 1, 4, 90, 2, 8)
	r2 := monthlyService.PreviewScores(full)
	if r2 == nil || r2.CompletenessPct < 70 {
		t.Fatalf("expected completeness >=70 for full signals, got %v", r2)
	}
}

func TestCopyFromPreviousMonth_PopulatesPreview(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create member and service
	name, _ := domain.NewFullName("Prev", "User")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Create an entry for previous month
	prevMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
	a, b, c, d, e, f, g, h, i, j := 2, 70, 3, 10, 1, 0, 1, 85, 1, 2
	_, err := monthlyService.CreateEntry(1, prevMonth, &a, &b, &c, &d, &e, &f, &g, &h, &i, &j)
	if err != nil {
		t.Fatalf("failed to create previous month entry: %v", err)
	}

	// Retrieve and preview
	entry, err := monthlyService.GetEntry(1, prevMonth)
	if err != nil || entry == nil {
		t.Fatalf("expected prior entry, got err=%v entry=%v", err, entry)
	}
	res := monthlyService.PreviewScores(entry.Signals())
	if res == nil {
		t.Fatal("expected preview result from previous month signals")
	}
	if res.CompletenessPct == 0 {
		t.Fatalf("unexpected completeness 0 for populated previous month: %v", res)
	}
}

// stubMemberService implements minimal interface used by screen constructor
type stubMemberService struct {
	members []domain.TeamMember
}

func (s *stubMemberService) ListMembers() ([]domain.TeamMember, error) { return s.members, nil }

package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

// TestMainWindowStartsStops is a smoke test that constructs the full main
// window using NewMainWindow and ensures it can show and close without panicking.
func TestMainWindowStartsStops(t *testing.T) {
	// Prepare in-memory repos and services
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Add one active member so the screen renders the full form
	name, _ := domain.NewFullName("Integration", "User")
	memberAgg, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	if err := memberRepo.Save(memberAgg); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	services := &ApplicationServices{
		MemberService:  memberService,
		MonthlyService: monthlyService,
	}

	app := test.NewApp()
	// create window via NewMainWindow
	w := NewMainWindow(app, services)
	w.Show()

	// Let the event loop run briefly to ensure rendering occurs
	time.Sleep(50 * time.Millisecond)

	// Close window and quit app
	w.Close()
	time.Sleep(10 * time.Millisecond)
	app.Quit()
}

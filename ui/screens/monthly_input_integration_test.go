package screens

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

// TestRunFyneAppStartsAndStops creates a fyne app, mounts the Monthly Input
// screen, shows the window briefly and then closes it. This is a smoke test
// that the screen can be constructed and rendered without panicking.
func TestRunFyneAppStartsAndStops(t *testing.T) {
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

	// Build the screen
	content := NewMonthlyInputScreen(memberService, monthlyService)

	// Start fyne test app and show a window with the content
	app := test.NewApp()
	w := app.NewWindow("integration-test")
	w.SetContent(content)
	w.Show()

	// Let the event loop run briefly to ensure rendering occurs
	time.Sleep(50 * time.Millisecond)

	// Close window and quit app
	w.Close()
	// ensure the app has time to process the close
	time.Sleep(10 * time.Millisecond)
	app.Quit()
}

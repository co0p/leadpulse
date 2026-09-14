package ui

import (
	"fyne.io/fyne/v2"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
	"leadpulse/ui/screens"
)

// ApplicationServices holds all service interfaces needed by the UI.
// The UI depends on service interfaces, not on database or repository concerns.
type ApplicationServices struct {
	MemberService  *member.Service
	MonthlyService *monthly.Service
}

// NewMainWindow creates and returns the main application window.
// The UI accepts only service interfaces; no database coupling.
func NewMainWindow(app fyne.App, services *ApplicationServices) fyne.Window {
	window := app.NewWindow("Team Impact Scorecard")
	window.Resize(fyne.NewSize(1200, 800))

	// Create screen registry with injected services
	screenRegistry := screens.NewScreenRegistry(services.MemberService)

	window.SetContent(screenRegistry.SettingsScreen())

	return window
}

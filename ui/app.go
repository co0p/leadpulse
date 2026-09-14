package ui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"leadpulse/service/member"
	"leadpulse/store"
	"leadpulse/ui/screens"
)

// NewMainWindow creates and returns the main application window.
func NewMainWindow(app fyne.App, db *sql.DB) fyne.Window {
	window := app.NewWindow("Team Impact Scorecard")
	window.Resize(fyne.NewSize(1200, 800))

	// Create repository implementations
	memberRepo := store.NewSQLiteTeamMemberRepository(db)
	entryRepo := store.NewSQLiteMonthlyEntryRepository(db)

	// Create services with injected repositories
	memberService := member.NewService(memberRepo, entryRepo)

	// Create screen registry
	screenRegistry := screens.NewScreenRegistry(memberService)

	window.SetContent(screenRegistry.SettingsScreen())

	return window
}

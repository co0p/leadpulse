package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"fyne.io/fyne/v2/app"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
	"leadpulse/store"
	"leadpulse/ui"
)

func main() {
	// Initialize database
	dbPath, err := getDatabasePath()
	if err != nil {
		log.Fatalf("Failed to get database path: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := store.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize repositories (persistence layer)
	memberRepo := store.NewSQLiteTeamMemberRepository(db)
	entryRepo := store.NewSQLiteMonthlyEntryRepository(db)

	// Initialize application services (use case layer)
	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthly.NewService(memberRepo, entryRepo)

	// Create service container
	services := &ui.ApplicationServices{
		MemberService:  memberService,
		MonthlyService: monthlyService,
	}

	// Create Fyne app
	fyneApp := app.New()

	// Create and show the main window (UI accepts services only)
	window := ui.NewMainWindow(fyneApp, services)
	window.ShowAndRun()
}

// getDatabasePath returns the path to the SQLite database file.
func getDatabasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, "leadpulse")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "data.db"), nil
}

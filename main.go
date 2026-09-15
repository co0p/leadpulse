package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"leadpulse/server"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
	"leadpulse/store"
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

	// Start HTTP server with coordinators available for HTTP handlers
	// (Coordinators can be injected into handlers via dependency injection)
	_ = memberService  // Available for HTTP handlers
	_ = monthlyService // Available for HTTP handlers

	// Boot HTTP server (blocks indefinitely)
	if err := server.Start("localhost:8080"); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
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

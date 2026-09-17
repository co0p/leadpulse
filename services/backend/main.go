package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"leadpulse/core/members"
	"leadpulse/server"
	"leadpulse/storage/sqlite"
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

	// Initialize storage adapters (persistence layer)
	memberRepo := sqlite.NewSQLiteTeamMemberRepository(db)

	// Initialize use cases (application layer)
	addMemberUC := members.NewAddMemberUseCase(memberRepo)
	getMembersUC := members.NewGetMembersUseCase(memberRepo)
	getFilteredMembersUC := members.NewGetFilteredMembersUseCase(memberRepo)
	editMemberUC := members.NewEditMemberUseCase(memberRepo)
	deactivateMemberUC := members.NewDeactivateMemberUseCase(memberRepo)
	reactivateMemberUC := members.NewReactivateMemberUseCase(memberRepo)

	// Boot HTTP server with use cases for API handlers (blocks indefinitely)
	if err := server.Start("0.0.0.0:8080", addMemberUC, getMembersUC, getFilteredMembersUC, editMemberUC, deactivateMemberUC, reactivateMemberUC); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

// getDatabasePath returns the path to the SQLite database file.
// In Docker, uses /data/data.db (volume mount).
// For local development, falls back to user config directory.
func getDatabasePath() (string, error) {
	// Check if running in Docker (volume mount available)
	if _, err := os.Stat("/data"); err == nil {
		return "/data/data.db", nil
	}

	// Fallback to local user config directory for development
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

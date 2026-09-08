package main

import (
	"embed"
	"os"
	"path/filepath"

	application "leadpulse/internal/app"
	"leadpulse/internal/employees"
	"leadpulse/internal/performances"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dataDirectory, err := os.UserConfigDir()
	if err != nil {
		println("Error:", err.Error())
		return
	}
	databasePath := filepath.Join(dataDirectory, "LeadPulse", "leadpulse.db")
	if databaseOverride := os.Getenv("LEADPULSE_DATABASE_PATH"); databaseOverride != "" {
		databasePath = databaseOverride
	}
	if err := os.MkdirAll(filepath.Dir(databasePath), 0755); err != nil {
		println("Error:", err.Error())
		return
	}
	employeeStore, err := employees.NewStore(databasePath)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	defer employeeStore.Close()

	performanceStore, err := performances.NewStore(databasePath)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	defer performanceStore.Close()

	app := application.New(employeeStore, performanceStore)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "leadpulse",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
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

	registry := screens.NewScreenRegistry(services.MemberService, services.MonthlyService)

	// Main content area — swapped when nav selection changes.
	content := container.NewStack()

	// Screen definitions: label, nav key, and loader.
	type navItem struct {
		label  string
		loader func() fyne.CanvasObject
	}

	navItems := []navItem{
		{"Overview", func() fyne.CanvasObject { return screens.PlaceholderScreen("Overview") }},
		{"Input", registry.MonthlyInputScreen},
		{"Alerts", func() fyne.CanvasObject { return screens.PlaceholderScreen("Alerts") }},
		{"Review", func() fyne.CanvasObject { return screens.PlaceholderScreen("Review") }},
	}

	settingsItem := navItem{"Settings", registry.SettingsScreen}

	// setScreen replaces the main content area.
	setScreen := func(loader func() fyne.CanvasObject) {
		content.Objects = []fyne.CanvasObject{loader()}
		content.Refresh()
	}

	// Start on Input (the current increment focus).
	setScreen(registry.MonthlyInputScreen)

	// Sidebar background.
	sidebarBg := canvas.NewRectangle(color.NRGBA{R: 241, G: 245, B: 249, A: 255}) // slate-100

	// Nav buttons for primary items.
	var navButtons []*widget.Button
	for _, item := range navItems {
		item := item // capture
		btn := widget.NewButton(item.label, func() {
			setScreen(item.loader)
		})
		btn.Alignment = widget.ButtonAlignLeading
		navButtons = append(navButtons, btn)
	}

	settingsBtn := widget.NewButton(settingsItem.label, func() {
		setScreen(settingsItem.loader)
	})
	settingsBtn.Alignment = widget.ButtonAlignLeading

	// Sidebar layout: app name at top, nav items, spacer, settings at bottom.
	appTitle := widget.NewRichTextFromMarkdown("**TIS**")

	primaryNav := container.NewVBox()
	for _, btn := range navButtons {
		primaryNav.Add(btn)
	}

	sidebar := container.NewStack(
		sidebarBg,
		container.NewBorder(
			container.NewVBox(appTitle, widget.NewSeparator()),
			container.NewVBox(widget.NewSeparator(), settingsBtn),
			nil, nil,
			container.NewVBox(primaryNav),
		),
	)
	sidebar.Resize(fyne.NewSize(160, 800))

	// Header: cycle context + separator.
	cycle := currentCycleLabel()
	header := container.NewBorder(
		nil,
		widget.NewSeparator(),
		nil,
		nil,
		container.NewHBox(
			widget.NewIcon(theme.InfoIcon()),
			widget.NewLabel("Cycle: "+cycle),
			layout.NewSpacer(),
		),
	)

	// Shell: sidebar left, header top, content fills remainder.
	shell := container.NewBorder(
		header,
		nil,
		sidebar,
		nil,
		content,
	)

	window.SetContent(shell)
	return window
}

// currentCycleLabel returns the current month formatted as "Month YYYY".
func currentCycleLabel() string {
	return time.Now().Format("January 2006")
}

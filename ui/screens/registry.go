package screens

import (
	"fyne.io/fyne/v2"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
)

// ScreenRegistry manages all application screens.
type ScreenRegistry struct {
	memberService *member.Service
	monthlyService *monthly.Service
}

// NewScreenRegistry creates a new screen registry.
func NewScreenRegistry(memberService *member.Service, monthlyService *monthly.Service) *ScreenRegistry {
	return &ScreenRegistry{
		memberService: memberService,
		monthlyService: monthlyService,
	}
}

// MonthlyInputScreen returns the Monthly Input screen (Screen B).
func (r *ScreenRegistry) MonthlyInputScreen() fyne.CanvasObject {
	return NewMonthlyInputScreen(r.memberService, r.monthlyService)
}

// SettingsScreen returns the Settings screen (Screen F).
func (r *ScreenRegistry) SettingsScreen() fyne.CanvasObject {
	return NewSettingsScreen(r.memberService)
}

// PlaceholderScreen returns a simple placeholder for screens not yet implemented.
func PlaceholderScreen(title string) fyne.CanvasObject {
	return newPlaceholderScreen(title)
}

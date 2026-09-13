package screens

import (
	"fyne.io/fyne/v2"
	"leadpulse/service/member"
)

// ScreenRegistry manages all application screens.
type ScreenRegistry struct {
	memberService *member.Service
}

// NewScreenRegistry creates a new screen registry.
func NewScreenRegistry(memberService *member.Service) *ScreenRegistry {
	return &ScreenRegistry{
		memberService: memberService,
	}
}

// SettingsScreen returns the Settings screen (Screen F).
func (r *ScreenRegistry) SettingsScreen() fyne.CanvasObject {
	return NewSettingsScreen(r.memberService)
}

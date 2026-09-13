package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"leadpulse/service/member"
)

// NewSettingsScreen creates the Settings screen (Screen F).
func NewSettingsScreen(memberService *member.Service) fyne.CanvasObject {
	// Title
	title := widget.NewRichTextFromMarkdown("# Settings")

	// Member management section
	memberTitle := widget.NewRichTextFromMarkdown("## Team Members")

	// Create the member list
	memberList := widget.NewList(
		func() int {
			members, _ := memberService.ListMembers()
			return len(members)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewButton("Edit", func() {}),
				widget.NewButton("Deactivate", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			hbox := item.(*fyne.Container)
			members, _ := memberService.ListMembers()
			if id < len(members) {
				member := members[id]
				hbox.Objects[0].(*widget.Label).SetText(member.FullName())
				hbox.Objects[1].(*widget.Label).SetText(string(member.Seniority))
			}
		},
	)

	// Add Member button
	addButton := widget.NewButton("Add Member", func() {
		// Placeholder for dialog
	})

	// Button bar for member actions
	memberButtonBar := container.NewHBox(addButton)

	// Combine member section
	memberSection := container.NewVBox(
		memberTitle,
		memberList,
		memberButtonBar,
	)

	// Settings content
	settingsContent := container.NewVBox(
		title,
		widget.NewSeparator(),
		memberSection,
	)

	// Scroll container for long content
	scrollContainer := container.NewScroll(settingsContent)

	return scrollContainer
}

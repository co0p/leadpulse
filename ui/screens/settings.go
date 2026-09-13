package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	domain "leadpulse/engine/domain"
	"leadpulse/service/member"
)

// NewSettingsScreen creates the Settings screen (Screen F).
func NewSettingsScreen(memberService *member.Service) fyne.CanvasObject {
	// Title bar
	title := widget.NewRichTextFromMarkdown("# Settings")
	subtitle := widget.NewLabel("Manage team members")

	// Member list — refreshed after mutations
	var memberList *widget.List
	memberList = widget.NewList(
		func() int {
			members, _ := memberService.ListMembers()
			return len(members)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
				layout.NewSpacer(),
				widget.NewButton("Edit", func() {}),
				widget.NewButton("Deactivate", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			hbox := item.(*fyne.Container)
			members, _ := memberService.ListMembers()
			if id >= len(members) {
				return
			}
			m := members[id]
			hbox.Objects[0].(*widget.Label).SetText(m.Name().String())
			hbox.Objects[1].(*widget.Label).SetText(string(m.Seniority()))

			// Wire Edit button
			hbox.Objects[3].(*widget.Button).OnTapped = func() {
				showEditDialog(memberService, int64(m.ID()), m.Name().First, m.Name().Last, m.Seniority(), memberList)
			}
			// Wire Deactivate button
			hbox.Objects[4].(*widget.Button).OnTapped = func() {
				_ = memberService.DeactivateMember(int64(m.ID()))
				memberList.Refresh()
			}
		},
	)

	// Add Member button
	addButton := widget.NewButton("+ Add Member", func() {
		showAddDialog(memberService, memberList)
	})

	// Layout: title block at top, button bar at bottom, list fills middle
	top := container.NewVBox(title, subtitle, widget.NewSeparator())
	bottom := container.NewHBox(addButton)

	return container.NewBorder(top, bottom, nil, nil, memberList)
}

// showAddDialog opens the Add Member form dialog.
func showAddDialog(memberService *member.Service, list *widget.List) {
	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetPlaceHolder("First name")

	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetPlaceHolder("Last name")

	seniorities := seniorityStrings()
	senioritySelect := widget.NewSelect(seniorities, func(string) {})
	senioritySelect.SetSelected(seniorities[0])

	formItems := []*widget.FormItem{
		{Text: "First Name", Widget: firstNameEntry},
		{Text: "Last Name", Widget: lastNameEntry},
		{Text: "Seniority", Widget: senioritySelect},
	}

	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowForm("Add Member", "Add", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}
		_, _ = memberService.AddMember(
			firstNameEntry.Text,
			lastNameEntry.Text,
			domain.Seniority(senioritySelect.Selected),
		)
		list.Refresh()
	}, w)
}

// showEditDialog opens the Edit Member form dialog pre-filled with current values.
func showEditDialog(memberService *member.Service, id int64, firstName, lastName string, seniority domain.Seniority, list *widget.List) {
	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetText(firstName)

	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetText(lastName)

	seniorities := seniorityStrings()
	senioritySelect := widget.NewSelect(seniorities, func(string) {})
	senioritySelect.SetSelected(string(seniority))

	formItems := []*widget.FormItem{
		{Text: "First Name", Widget: firstNameEntry},
		{Text: "Last Name", Widget: lastNameEntry},
		{Text: "Seniority", Widget: senioritySelect},
	}

	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowForm("Edit Member", "Save", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}
		_, _ = memberService.EditMember(
			id,
			firstNameEntry.Text,
			lastNameEntry.Text,
			domain.Seniority(senioritySelect.Selected),
		)
		list.Refresh()
	}, w)
}

func seniorityStrings() []string {
	all := domain.AllSeniorities()
	out := make([]string, len(all))
	for i, s := range all {
		out[i] = string(s)
	}
	return out
}

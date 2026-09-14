package screens

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"leadpulse/engine/domain"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
	"leadpulse/ui/controllers"
)

// MonthlyInputScreen wraps the split container and exposes test getters.
type MonthlyInputScreen struct {
	*container.Split
	controller *controllers.MonthlyInputController
	fields     map[string]*widget.Entry
	memberList *widget.List
	saveButton *widget.Button
}

func NewMonthlyInputScreen(memberService *member.Service, monthlyService *monthly.Service) fyne.CanvasObject {
	if memberService == nil {
		return widget.NewLabel("Monthly Input")
	}

	// Create controller
	controller := controllers.NewMonthlyInputController(memberService, monthlyService)
	if err := controller.Load(); err != nil {
		return widget.NewLabel("Failed to load team members")
	}

	members := controller.GetMembers()
	if len(members) == 0 {
		return widget.NewLabel("No active team members")
	}

	labelFor := func(field string) string {
		switch field {
		case "morale":
			return "Morale (0-5)"
		case "billability":
			return "Billability % (0-100)"
		case "csat":
			return "CSAT (1-5)"
		case "net_margin":
			return "Net Margin % (-20 to +60)"
		case "positive_feedback":
			return "Positive Feedback (>=0)"
		case "critical_feedback":
			return "Critical Feedback (>=0)"
		case "overtime_hours":
			return "Overtime Hours (>=0)"
		case "delivery_reliability":
			return "Delivery Reliability % (0-100)"
		case "mentoring_hours":
			return "Mentoring Hours (>=0)"
		case "evidence_notes_count":
			return "Evidence Notes Count (>=0)"
		default:
			return field
		}
	}

	fieldNames := []string{"morale", "billability", "csat", "net_margin", "positive_feedback", "critical_feedback", "overtime_hours", "delivery_reliability", "mentoring_hours", "evidence_notes_count"}
	fields := map[string]*widget.Entry{}
	for _, name := range fieldNames {
		entry := widget.NewEntry()
		entry.SetPlaceHolder(labelFor(name))
		fields[name] = entry
	}

	statusLabel := widget.NewLabel("")
	statusLabel.Hide()
	previewLabel := widget.NewLabel("TII: 0.00")
	previewLabel.TextStyle = fyne.TextStyle{Bold: true}
	previewMeta := widget.NewLabel("Completeness: 0%")

	// Helper: populate form fields from controller state
	populateFormFields := func() {
		state := controller.GetFormState()
		if state.Morale != nil {
			fields["morale"].SetText(strconv.Itoa(*state.Morale))
		} else {
			fields["morale"].SetText("")
		}
		if state.Billability != nil {
			fields["billability"].SetText(strconv.Itoa(*state.Billability))
		} else {
			fields["billability"].SetText("")
		}
		if state.CSAT != nil {
			fields["csat"].SetText(strconv.Itoa(*state.CSAT))
		} else {
			fields["csat"].SetText("")
		}
		if state.NetMargin != nil {
			fields["net_margin"].SetText(strconv.Itoa(*state.NetMargin))
		} else {
			fields["net_margin"].SetText("")
		}
		if state.PositiveFeedback != nil {
			fields["positive_feedback"].SetText(strconv.Itoa(*state.PositiveFeedback))
		} else {
			fields["positive_feedback"].SetText("")
		}
		if state.CriticalFeedback != nil {
			fields["critical_feedback"].SetText(strconv.Itoa(*state.CriticalFeedback))
		} else {
			fields["critical_feedback"].SetText("")
		}
		if state.OvertimeHours != nil {
			fields["overtime_hours"].SetText(strconv.Itoa(*state.OvertimeHours))
		} else {
			fields["overtime_hours"].SetText("")
		}
		if state.DeliveryReliability != nil {
			fields["delivery_reliability"].SetText(strconv.Itoa(*state.DeliveryReliability))
		} else {
			fields["delivery_reliability"].SetText("")
		}
		if state.MentoringHours != nil {
			fields["mentoring_hours"].SetText(strconv.Itoa(*state.MentoringHours))
		} else {
			fields["mentoring_hours"].SetText("")
		}
		if state.EvidenceNotesCount != nil {
			fields["evidence_notes_count"].SetText(strconv.Itoa(*state.EvidenceNotesCount))
		} else {
			fields["evidence_notes_count"].SetText("")
		}
	}

	// Helper: parse a string field to *int, or nil if empty
	parsePtr := func(text string) (*int, error) {
		if text == "" {
			return nil, nil
		}
		val, err := strconv.Atoi(text)
		return &val, err
	}

	// Helper: read form fields and update controller state
	readFormFields := func() error {
		for _, fieldName := range fieldNames {
			val, err := parsePtr(fields[fieldName].Text)
			if err != nil {
				return fmt.Errorf("invalid number in %s", labelFor(fieldName))
			}
			controller.SetFormField(fieldName, val)
		}
		return nil
	}

	// Helper: update preview based on current form
	var updatePreview func()
	var saveButton *widget.Button

	updatePreview = func() {
		if err := readFormFields(); err != nil {
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}

		// Clear validation message now that parsing succeeded
		statusLabel.Hide()

		state := controller.GetFormState()
		signals, err := domain.NewMonthlyRawSignalsFromPointers(
			state.Morale, state.Billability, state.CSAT, state.NetMargin,
			state.PositiveFeedback, state.CriticalFeedback, state.OvertimeHours, state.DeliveryReliability, state.MentoringHours, state.EvidenceNotesCount,
		)
		if err != nil {
			// Domain-level validation failed (e.g. out-of-range)
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}

		res := monthlyService.PreviewScores(signals)
		if res == nil {
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			saveButton.Disable()
			return
		}
		previewLabel.SetText(fmt.Sprintf("TII: %.2f", res.TII))
		previewMeta.SetText(fmt.Sprintf("Completeness: %.0f%%", res.CompletenessPct))
		// Always allow Save; completeness is shown for reference only
		saveButton.Enable()
	}

	// loadMemberData loads the form fields for the given member via controller
	var loadMemberData func(widget.ListItemID)
	loadMemberData = func(id widget.ListItemID) {
		if id >= len(members) {
			return
		}
		m := members[id]
		statusLabel.SetText(fmt.Sprintf("Member: %s", m.Name().String()))
		statusLabel.Show()

		// Use controller to select member and load data
		if err := controller.SelectMember(int(id)); err != nil {
			statusLabel.SetText(err.Error())
			return
		}

		// Populate form from controller state
		populateFormFields()

		// Update the preview for the newly-selected member
		updatePreview()
	}

	var memberList *widget.List
	memberList = widget.NewList(
		func() int { return len(members) },
		func() fyne.CanvasObject {
			// Simple HBox with two labels: name + status
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(members) {
				return
			}
			m := members[id]
			row := item.(*fyne.Container)
			nameLabel := row.Objects[0].(*widget.Label)
			statusLabel := row.Objects[1].(*widget.Label)
			nameLabel.SetText(m.Name().String())
			if int(id) == controller.GetSelectedIndex() {
				nameLabel.SetText(m.Name().String() + " (selected)")
			}
			// Check if entry exists by trying to load it
			entry, _ := monthlyService.GetEntry(int64(m.ID()), controller.GetCurrentMonth())
			if entry != nil {
				statusLabel.SetText("Done")
			} else {
				statusLabel.SetText("Draft")
			}
		},
	)

	memberList.OnSelected = func(id widget.ListItemID) {
		if id >= len(members) {
			return
		}
		memberList.Refresh()
		loadMemberData(id)
	}

	saveButton = widget.NewButton("Save", func() {
		// Read form fields into controller state
		if err := readFormFields(); err != nil {
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			return
		}

		// Save via controller
		if err := controller.SaveMember(); err != nil {
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			return
		}

		statusLabel.SetText("Saved")
		statusLabel.Show()

		// Refresh the member list to update Draft/Done badges
		memberList.Refresh()

		// Reload the form for the current member to show persisted data
		memberIndex := controller.GetSelectedIndex()
		if memberIndex >= 0 {
			memberList.UnselectAll()
			memberList.Select(widget.ListItemID(memberIndex))
		}
	})
	// Initially enable Save; allow partial data entry
	saveButton.Enable()

	copyButton := widget.NewButton("Copy from previous month", func() {
		if err := controller.CopyFromPreviousMonth(); err != nil {
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			return
		}
		populateFormFields()
		statusLabel.SetText("Copied from previous month")
		statusLabel.Show()
		updatePreview()
	})

	// Wire live preview: update preview on any field change
	for _, e := range fields {
		e.OnChanged = func(_ string) {
			updatePreview()
		}
	}

	formItems := container.NewVBox()
	for _, name := range fieldNames {
		formItems.Add(widget.NewLabel(labelFor(name)))
		formItems.Add(fields[name])
	}

	rightPanel := container.NewVBox(
		widget.NewLabel("Monthly Input"),
		previewLabel,
		previewMeta,
		formItems,
		statusLabel,
		copyButton,
		saveButton,
	)

	leftPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Team Members"),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		memberList,
	)

	split := container.NewHSplit(leftPanel, container.NewVScroll(rightPanel))
	split.Offset = 0.4

	// Wrap in MonthlyInputScreen to provide test getters
	screen := &MonthlyInputScreen{
		Split:      split,
		controller: controller,
		fields:     fields,
		memberList: memberList,
		saveButton: saveButton,
	}
	return screen
}

// ========== Test-Only Getters ==========
// These methods are for testing only and should NOT be used in production code.
// They allow tests to access internal screen components without reflection.

// GetFormFieldForTesting returns an Entry widget for a specific signal field by name.
// fieldName should be one of: "morale", "billability", "csat", etc.
// Returns nil if the field doesn't exist. Only use this in tests.
func (s *MonthlyInputScreen) GetFormFieldForTesting(fieldName string) *widget.Entry {
	if s.fields == nil {
		return nil
	}
	return s.fields[fieldName]
}

// GetMemberListForTesting returns the member List widget.
// Only use this in tests.
func (s *MonthlyInputScreen) GetMemberListForTesting() *widget.List {
	return s.memberList
}

// GetSaveButtonForTesting returns the Save button widget.
// Only use this in tests.
func (s *MonthlyInputScreen) GetSaveButtonForTesting() *widget.Button {
	return s.saveButton
}

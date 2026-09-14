package screens

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"leadpulse/engine/domain"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
)

func NewMonthlyInputScreen(memberService *member.Service, monthlyService *monthly.Service) fyne.CanvasObject {
	if memberService == nil {
		return widget.NewLabel("Monthly Input")
	}

	members, err := memberService.ListMembers()
	if err != nil || len(members) == 0 {
		return widget.NewLabel("No active team members")
	}

	currentMonth := time.Now().Format("2006-01")
	selected := 0

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

	// saveButton is declared here so closures defined earlier can reference it;
	// it will be assigned after its construction further down the function.
	var saveButton *widget.Button

	// updatePreview reads current form values, attempts to create a pointer-backed
	// MonthlyRawSignals via the domain constructor that accepts pointers, and
	// then asks the monthly service for a preview scoring result. UI labels are
	// updated to reflect the result or any validation errors.
	updatePreview := func() {
		var (
			pMorale, pBillability, pCSAT, pNetMargin   *int
			pPositive, pCritical, pOvertime, pDelivery *int
			pMentoring, pEvidence                      *int
		)

		parsePtr := func(text string) (*int, error) {
			if text == "" {
				return nil, nil
			}
			v, err := strconv.Atoi(text)
			if err != nil {
				return nil, err
			}
			return &v, nil
		}

		var err error
		if pMorale, err = parsePtr(fields["morale"].Text); err != nil {
			statusLabel.SetText("Invalid number in Morale")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pBillability, err = parsePtr(fields["billability"].Text); err != nil {
			statusLabel.SetText("Invalid number in Billability")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pCSAT, err = parsePtr(fields["csat"].Text); err != nil {
			statusLabel.SetText("Invalid number in CSAT")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pNetMargin, err = parsePtr(fields["net_margin"].Text); err != nil {
			statusLabel.SetText("Invalid number in Net Margin")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pPositive, err = parsePtr(fields["positive_feedback"].Text); err != nil {
			statusLabel.SetText("Invalid number in Positive Feedback")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pCritical, err = parsePtr(fields["critical_feedback"].Text); err != nil {
			statusLabel.SetText("Invalid number in Critical Feedback")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pOvertime, err = parsePtr(fields["overtime_hours"].Text); err != nil {
			statusLabel.SetText("Invalid number in Overtime Hours")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pDelivery, err = parsePtr(fields["delivery_reliability"].Text); err != nil {
			statusLabel.SetText("Invalid number in Delivery Reliability")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pMentoring, err = parsePtr(fields["mentoring_hours"].Text); err != nil {
			statusLabel.SetText("Invalid number in Mentoring Hours")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}
		if pEvidence, err = parsePtr(fields["evidence_notes_count"].Text); err != nil {
			statusLabel.SetText("Invalid number in Evidence Notes")
			statusLabel.Show()
			previewLabel.SetText("TII: —")
			previewMeta.SetText("Completeness: —")
			return
		}

		// Clear validation message now that parsing succeeded
		statusLabel.Hide()

		signals, err := domain.NewMonthlyRawSignalsFromPointers(
			pMorale, pBillability, pCSAT, pNetMargin,
			pPositive, pCritical, pOvertime, pDelivery, pMentoring, pEvidence,
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
		// Enable Save only when completeness >= 70% (per plan thresholds)
		if res.CompletenessPct >= 70.0 {
			saveButton.Enable()
		} else {
			saveButton.Disable()
		}
	}

	// (helper function ParseSignalsFromStrings removed; tests use a local parser)

	refreshMemberList := func(list *widget.List) {
		list.Refresh()
	}

	memberList := widget.NewList(
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
			if id == selected {
				nameLabel.SetText(m.Name().String() + " (selected)")
			}
			if _, err := monthlyService.GetEntry(int64(m.ID()), currentMonth); err == nil {
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
		selected = int(id)
		memberList.Refresh()
		m := members[id]
		statusLabel.SetText(fmt.Sprintf("Member: %s", m.Name().String()))
		statusLabel.Show()
		for _, entry := range fields {
			entry.SetText("")
		}
		if currentEntry, err := monthlyService.GetEntry(int64(m.ID()), currentMonth); err == nil && currentEntry != nil {
			signals := currentEntry.Signals()
			if signals.Morale != nil {
				fields["morale"].SetText(strconv.Itoa(*signals.Morale))
			}
			if signals.Billability != nil {
				fields["billability"].SetText(strconv.Itoa(*signals.Billability))
			}
			if signals.CSAT != nil {
				fields["csat"].SetText(strconv.Itoa(*signals.CSAT))
			}
			if signals.NetMargin != nil {
				fields["net_margin"].SetText(strconv.Itoa(*signals.NetMargin))
			}
			if signals.PositiveFeedback != nil {
				fields["positive_feedback"].SetText(strconv.Itoa(*signals.PositiveFeedback))
			}
			if signals.CriticalFeedback != nil {
				fields["critical_feedback"].SetText(strconv.Itoa(*signals.CriticalFeedback))
			}
			if signals.OvertimeHours != nil {
				fields["overtime_hours"].SetText(strconv.Itoa(*signals.OvertimeHours))
			}
			if signals.DeliveryReliability != nil {
				fields["delivery_reliability"].SetText(strconv.Itoa(*signals.DeliveryReliability))
			}
			if signals.MentoringHours != nil {
				fields["mentoring_hours"].SetText(strconv.Itoa(*signals.MentoringHours))
			}
			if signals.EvidenceNotesCount != nil {
				fields["evidence_notes_count"].SetText(strconv.Itoa(*signals.EvidenceNotesCount))
			}
		}
		// Update the preview for the newly-selected member
		updatePreview()
	}

	saveButton = widget.NewButton("Save", func() {
		m := members[selected]
		values := map[string]*int{}
		for _, fieldName := range fieldNames {
			text := fields[fieldName].Text
			if text == "" {
				values[fieldName] = nil
				continue
			}
			parsed, err := strconv.Atoi(text)
			if err != nil {
				statusLabel.SetText("Please enter valid numeric values")
				statusLabel.Show()
				return
			}
			values[fieldName] = &parsed
		}
		_ = m
		// Prepare pointer args directly from values map
		ptr := func(v *int) *int { return v }
		if _, err := monthlyService.GetEntry(int64(m.ID()), currentMonth); err == nil {
			_, err = monthlyService.UpdateEntry(
				int64(m.ID()), currentMonth,
				ptr(values["morale"]), ptr(values["billability"]), ptr(values["csat"]), ptr(values["net_margin"]),
				ptr(values["positive_feedback"]), ptr(values["critical_feedback"]),
				ptr(values["overtime_hours"]), ptr(values["delivery_reliability"]),
				ptr(values["mentoring_hours"]), ptr(values["evidence_notes_count"]),
			)
		} else {
			_, err = monthlyService.CreateEntry(
				int64(m.ID()), currentMonth,
				ptr(values["morale"]), ptr(values["billability"]), ptr(values["csat"]), ptr(values["net_margin"]),
				ptr(values["positive_feedback"]), ptr(values["critical_feedback"]),
				ptr(values["overtime_hours"]), ptr(values["delivery_reliability"]),
				ptr(values["mentoring_hours"]), ptr(values["evidence_notes_count"]),
			)
		}
		if err != nil {
			statusLabel.SetText(err.Error())
			statusLabel.Show()
			return
		}
		statusLabel.SetText("Saved")
		statusLabel.Show()
		refreshMemberList(memberList)
	})
	// Initially decide whether Save should be enabled
	// Disable when completeness < 70% per plan; updatePreview will enable/disable
	saveButton.Disable()

	copyButton := widget.NewButton("Copy from previous month", func() {
		m := members[selected]
		prevMonth := previousMonth(currentMonth)
		entry, err := monthlyService.GetEntry(int64(m.ID()), prevMonth)
		if err != nil || entry == nil {
			statusLabel.SetText("No prior month entry")
			statusLabel.Show()
			return
		}
		signals := entry.Signals()
		if signals.Morale != nil {
			fields["morale"].SetText(strconv.Itoa(*signals.Morale))
		}
		if signals.Billability != nil {
			fields["billability"].SetText(strconv.Itoa(*signals.Billability))
		}
		if signals.CSAT != nil {
			fields["csat"].SetText(strconv.Itoa(*signals.CSAT))
		}
		if signals.NetMargin != nil {
			fields["net_margin"].SetText(strconv.Itoa(*signals.NetMargin))
		}
		if signals.PositiveFeedback != nil {
			fields["positive_feedback"].SetText(strconv.Itoa(*signals.PositiveFeedback))
		}
		if signals.CriticalFeedback != nil {
			fields["critical_feedback"].SetText(strconv.Itoa(*signals.CriticalFeedback))
		}
		if signals.OvertimeHours != nil {
			fields["overtime_hours"].SetText(strconv.Itoa(*signals.OvertimeHours))
		}
		if signals.DeliveryReliability != nil {
			fields["delivery_reliability"].SetText(strconv.Itoa(*signals.DeliveryReliability))
		}
		if signals.MentoringHours != nil {
			fields["mentoring_hours"].SetText(strconv.Itoa(*signals.MentoringHours))
		}
		if signals.EvidenceNotesCount != nil {
			fields["evidence_notes_count"].SetText(strconv.Itoa(*signals.EvidenceNotesCount))
		}
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

	leftPanel := container.NewBorder(nil, nil, nil, nil, memberList)
	leftPanel.Objects = []fyne.CanvasObject{widget.NewLabel("Team Members"), memberList}
	return container.NewHSplit(leftPanel, rightPanel)
}

func previousMonth(month string) string {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return month
	}
	return t.AddDate(0, -1, 0).Format("2006-01")
}

func valueOrNil(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

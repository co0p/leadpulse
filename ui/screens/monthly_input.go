package screens

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	refreshMemberList := func(list *widget.List) {
		list.Refresh()
	}

	memberList := widget.NewList(
		func() int { return len(members) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil, container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
			))
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
			if signals.Morale != nil { fields["morale"].SetText(strconv.Itoa(*signals.Morale)) }
			if signals.Billability != nil { fields["billability"].SetText(strconv.Itoa(*signals.Billability)) }
			if signals.CSAT != nil { fields["csat"].SetText(strconv.Itoa(*signals.CSAT)) }
			if signals.NetMargin != nil { fields["net_margin"].SetText(strconv.Itoa(*signals.NetMargin)) }
			if signals.PositiveFeedback != nil { fields["positive_feedback"].SetText(strconv.Itoa(*signals.PositiveFeedback)) }
			if signals.CriticalFeedback != nil { fields["critical_feedback"].SetText(strconv.Itoa(*signals.CriticalFeedback)) }
			if signals.OvertimeHours != nil { fields["overtime_hours"].SetText(strconv.Itoa(*signals.OvertimeHours)) }
			if signals.DeliveryReliability != nil { fields["delivery_reliability"].SetText(strconv.Itoa(*signals.DeliveryReliability)) }
			if signals.MentoringHours != nil { fields["mentoring_hours"].SetText(strconv.Itoa(*signals.MentoringHours)) }
			if signals.EvidenceNotesCount != nil { fields["evidence_notes_count"].SetText(strconv.Itoa(*signals.EvidenceNotesCount)) }
		}
	}

	saveButton := widget.NewButton("Save", func() {
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
		if _, err := monthlyService.GetEntry(int64(m.ID()), currentMonth); err == nil {
			_, err = monthlyService.UpdateEntry(
				int64(m.ID()), currentMonth,
				valueOrNil(values["morale"]), valueOrNil(values["billability"]), valueOrNil(values["csat"]), valueOrNil(values["net_margin"]),
				valueOrNil(values["positive_feedback"]), valueOrNil(values["critical_feedback"]),
				valueOrNil(values["overtime_hours"]), valueOrNil(values["delivery_reliability"]),
				valueOrNil(values["mentoring_hours"]), valueOrNil(values["evidence_notes_count"]),
			)
		} else {
			_, err = monthlyService.CreateEntry(
				int64(m.ID()), currentMonth,
				valueOrNil(values["morale"]), valueOrNil(values["billability"]), valueOrNil(values["csat"]), valueOrNil(values["net_margin"]),
				valueOrNil(values["positive_feedback"]), valueOrNil(values["critical_feedback"]),
				valueOrNil(values["overtime_hours"]), valueOrNil(values["delivery_reliability"]),
				valueOrNil(values["mentoring_hours"]), valueOrNil(values["evidence_notes_count"]),
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
		if signals.Morale != nil { fields["morale"].SetText(strconv.Itoa(*signals.Morale)) }
		if signals.Billability != nil { fields["billability"].SetText(strconv.Itoa(*signals.Billability)) }
		if signals.CSAT != nil { fields["csat"].SetText(strconv.Itoa(*signals.CSAT)) }
		if signals.NetMargin != nil { fields["net_margin"].SetText(strconv.Itoa(*signals.NetMargin)) }
		if signals.PositiveFeedback != nil { fields["positive_feedback"].SetText(strconv.Itoa(*signals.PositiveFeedback)) }
		if signals.CriticalFeedback != nil { fields["critical_feedback"].SetText(strconv.Itoa(*signals.CriticalFeedback)) }
		if signals.OvertimeHours != nil { fields["overtime_hours"].SetText(strconv.Itoa(*signals.OvertimeHours)) }
		if signals.DeliveryReliability != nil { fields["delivery_reliability"].SetText(strconv.Itoa(*signals.DeliveryReliability)) }
		if signals.MentoringHours != nil { fields["mentoring_hours"].SetText(strconv.Itoa(*signals.MentoringHours)) }
		if signals.EvidenceNotesCount != nil { fields["evidence_notes_count"].SetText(strconv.Itoa(*signals.EvidenceNotesCount)) }
		statusLabel.SetText("Copied from previous month")
		statusLabel.Show()
	})

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

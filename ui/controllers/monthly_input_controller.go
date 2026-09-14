package controllers

import (
	"fmt"
	"time"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
)

// MonthlyInputController owns all business logic for the Monthly Input screen.
// It manages form state, member selection, data persistence, and validation.
// Controllers have zero Fyne dependencies and are fully unit testable.
type MonthlyInputController struct {
	memberService  *member.Service
	monthlyService *monthly.Service

	// State
	currentMonth    string
	members         []domain.TeamMember
	selectedIndex   int
	currentMember   *domain.TeamMember
	formState       MonthlyInputState
}

// MonthlyInputState holds the current form field values.
type MonthlyInputState struct {
	Morale               *int
	Billability          *int
	CSAT                 *int
	NetMargin            *int
	PositiveFeedback     *int
	CriticalFeedback     *int
	OvertimeHours        *int
	DeliveryReliability  *int
	MentoringHours       *int
	EvidenceNotesCount   *int
}

// NewMonthlyInputController creates a new controller for the Monthly Input screen.
func NewMonthlyInputController(memberService *member.Service, monthlyService *monthly.Service) *MonthlyInputController {
	return &MonthlyInputController{
		memberService:  memberService,
		monthlyService: monthlyService,
		currentMonth:   time.Now().Format("2006-01"),
		selectedIndex:  -1,
		formState:      MonthlyInputState{},
	}
}

// Load initializes the controller by loading all team members.
func (c *MonthlyInputController) Load() error {
	members, err := c.memberService.ListMembers()
	if err != nil {
		return err
	}
	c.members = members
	return nil
}

// Validate checks if the form state is valid (all numeric fields parse correctly).
func (c *MonthlyInputController) Validate() error {
	// Form validation: no rules at this layer.
	// Individual field validation is done at the UI layer (parsePtr).
	// Data persistence validation is done by the service layer.
	return nil
}

// SelectMember loads data for the given member (by list index).
// Clears form fields before loading new member's data.
func (c *MonthlyInputController) SelectMember(index int) error {
	if index < 0 || index >= len(c.members) {
		return fmt.Errorf("invalid member index: %d", index)
	}

	c.selectedIndex = index
	c.currentMember = &c.members[index]

	// Clear form
	c.formState = MonthlyInputState{}

	// Load existing entry for this member in the current month
	entry, err := c.monthlyService.GetEntry(int64(c.currentMember.ID()), c.currentMonth)
	if err != nil || entry == nil {
		// No existing entry; form stays empty
		return nil
	}

	// Populate form from loaded entry
	signals := entry.Signals()
	c.formState.Morale = signals.Morale
	c.formState.Billability = signals.Billability
	c.formState.CSAT = signals.CSAT
	c.formState.NetMargin = signals.NetMargin
	c.formState.PositiveFeedback = signals.PositiveFeedback
	c.formState.CriticalFeedback = signals.CriticalFeedback
	c.formState.OvertimeHours = signals.OvertimeHours
	c.formState.DeliveryReliability = signals.DeliveryReliability
	c.formState.MentoringHours = signals.MentoringHours
	c.formState.EvidenceNotesCount = signals.EvidenceNotesCount

	return nil
}

// SetFormField updates a single form field value.
func (c *MonthlyInputController) SetFormField(fieldName string, value *int) error {
	switch fieldName {
	case "morale":
		c.formState.Morale = value
	case "billability":
		c.formState.Billability = value
	case "csat":
		c.formState.CSAT = value
	case "net_margin":
		c.formState.NetMargin = value
	case "positive_feedback":
		c.formState.PositiveFeedback = value
	case "critical_feedback":
		c.formState.CriticalFeedback = value
	case "overtime_hours":
		c.formState.OvertimeHours = value
	case "delivery_reliability":
		c.formState.DeliveryReliability = value
	case "mentoring_hours":
		c.formState.MentoringHours = value
	case "evidence_notes_count":
		c.formState.EvidenceNotesCount = value
	default:
		return fmt.Errorf("unknown field: %s", fieldName)
	}
	return nil
}

// GetFormState returns a copy of the current form state.
func (c *MonthlyInputController) GetFormState() MonthlyInputState {
	return c.formState
}

// ClearForm resets all form fields to empty.
func (c *MonthlyInputController) ClearForm() {
	c.formState = MonthlyInputState{}
}

// SaveMember persists the current form state to the service.
// Creates new entry if one doesn't exist; updates if it does.
// FIX for pointer aliasing bug: Creates a copy of each value to avoid
// multiple pointers referencing the same memory location.
func (c *MonthlyInputController) SaveMember() error {
	if c.currentMember == nil {
		return fmt.Errorf("no member selected")
	}

	// Helper: creates a copy of the pointer to avoid aliasing
	// (original bug was: values[fieldName] = &parsed where parsed is loop-scoped)
	copyPtr := func(v *int) *int {
		if v == nil {
			return nil
		}
		copy := *v
		return &copy
	}

	memberID := int64(c.currentMember.ID())

	// Check if entry exists
	existing, _ := c.monthlyService.GetEntry(memberID, c.currentMonth)
	var err error
	if existing != nil {
		// Update existing
		_, err = c.monthlyService.UpdateEntry(
			memberID, c.currentMonth,
			copyPtr(c.formState.Morale),
			copyPtr(c.formState.Billability),
			copyPtr(c.formState.CSAT),
			copyPtr(c.formState.NetMargin),
			copyPtr(c.formState.PositiveFeedback),
			copyPtr(c.formState.CriticalFeedback),
			copyPtr(c.formState.OvertimeHours),
			copyPtr(c.formState.DeliveryReliability),
			copyPtr(c.formState.MentoringHours),
			copyPtr(c.formState.EvidenceNotesCount),
		)
	} else {
		// Create new
		_, err = c.monthlyService.CreateEntry(
			memberID, c.currentMonth,
			copyPtr(c.formState.Morale),
			copyPtr(c.formState.Billability),
			copyPtr(c.formState.CSAT),
			copyPtr(c.formState.NetMargin),
			copyPtr(c.formState.PositiveFeedback),
			copyPtr(c.formState.CriticalFeedback),
			copyPtr(c.formState.OvertimeHours),
			copyPtr(c.formState.DeliveryReliability),
			copyPtr(c.formState.MentoringHours),
			copyPtr(c.formState.EvidenceNotesCount),
		)
	}

	return err
}

// GetMembers returns the list of all team members.
func (c *MonthlyInputController) GetMembers() []domain.TeamMember {
	return c.members
}

// GetSelectedMemberID returns the ID of the currently selected member, or -1 if none selected.
func (c *MonthlyInputController) GetSelectedMemberID() int64 {
	if c.currentMember == nil {
		return -1
	}
	return int64(c.currentMember.ID())
}

// GetSelectedIndex returns the index of the currently selected member in the members list.
func (c *MonthlyInputController) GetSelectedIndex() int {
	return c.selectedIndex
}

// GetCurrentMonth returns the current month in YYYY-MM format.
func (c *MonthlyInputController) GetCurrentMonth() string {
	return c.currentMonth
}

// CopyFromPreviousMonth loads data from the previous month and populates the form.
func (c *MonthlyInputController) CopyFromPreviousMonth() error {
	if c.currentMember == nil {
		return fmt.Errorf("no member selected")
	}

	prevMonth := previousMonth(c.currentMonth)
	entry, err := c.monthlyService.GetEntry(int64(c.currentMember.ID()), prevMonth)
	if err != nil || entry == nil {
		return fmt.Errorf("no entry found for previous month")
	}

	signals := entry.Signals()
	c.formState.Morale = signals.Morale
	c.formState.Billability = signals.Billability
	c.formState.CSAT = signals.CSAT
	c.formState.NetMargin = signals.NetMargin
	c.formState.PositiveFeedback = signals.PositiveFeedback
	c.formState.CriticalFeedback = signals.CriticalFeedback
	c.formState.OvertimeHours = signals.OvertimeHours
	c.formState.DeliveryReliability = signals.DeliveryReliability
	c.formState.MentoringHours = signals.MentoringHours
	c.formState.EvidenceNotesCount = signals.EvidenceNotesCount

	return nil
}

// previousMonth returns the month before the given month (YYYY-MM format).
func previousMonth(month string) string {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return month
	}
	return t.AddDate(0, -1, 0).Format("2006-01")
}

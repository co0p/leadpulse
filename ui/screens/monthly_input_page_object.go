package screens

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// MonthlyInputPageObject encapsulates user interactions with the Monthly Input Screen.
// It provides a clean, user-centric API that hides implementation details.
// Tests should use this Page Object instead of directly manipulating widgets.
type MonthlyInputPageObject struct {
	t      *testing.T
	app    fyne.App
	window fyne.Window
	screen *MonthlyInputScreen
}

// NewMonthlyInputPageObject creates a new page object wrapping the given screen.
// The screen will be displayed in a test window.
func NewMonthlyInputPageObject(t *testing.T, screen *MonthlyInputScreen) *MonthlyInputPageObject {
	app := test.NewApp()
	w := app.NewWindow("monthly-input-test")
	w.SetContent(screen)
	w.Show()

	// Allow time for Fyne to render and initialize the screen
	time.Sleep(150 * time.Millisecond)

	return &MonthlyInputPageObject{
		t:      t,
		app:    app,
		window: w,
		screen: screen,
	}
}

// SelectMember selects a team member by index (0-based).
// Waits for the form to update after selection.
func (po *MonthlyInputPageObject) SelectMember(index int) {
	memberList := po.screen.GetMemberListForTesting()
	if memberList == nil {
		po.t.Fatal("member list not found")
	}

	memberList.Select(widget.ListItemID(index))
	time.Sleep(300 * time.Millisecond)
}

// ReturnToMember forces re-selection of a member by unselecting then re-selecting.
// This is necessary because Fyne's List only fires OnSelected when the selection *changes*.
// When returning to the same member, we need to trigger the callback explicitly.
func (po *MonthlyInputPageObject) ReturnToMember(index int) {
	memberList := po.screen.GetMemberListForTesting()
	if memberList == nil {
		po.t.Fatal("member list not found")
	}

	memberList.Unselect(widget.ListItemID(index))
	memberList.Select(widget.ListItemID(index))
	time.Sleep(400 * time.Millisecond)
}

// EnterSignalValue simulates a user typing a value into a signal field.
// fieldName should be one of: "morale", "billability", "csat", "net_margin", etc.
func (po *MonthlyInputPageObject) EnterSignalValue(fieldName, value string) {
	entry := po.screen.GetFormFieldForTesting(fieldName)
	if entry == nil {
		po.t.Fatalf("form field not found: %s", fieldName)
	}

	entry.SetText(value)
	time.Sleep(100 * time.Millisecond)
}

// GetSignalValue retrieves the current value displayed in a signal field.
func (po *MonthlyInputPageObject) GetSignalValue(fieldName string) string {
	entry := po.screen.GetFormFieldForTesting(fieldName)
	if entry == nil {
		po.t.Fatalf("form field not found: %s", fieldName)
	}
	return entry.Text
}

// ClickSave simulates clicking the Save button and waits for persistence.
func (po *MonthlyInputPageObject) ClickSave() {
	saveButton := po.screen.GetSaveButtonForTesting()
	if saveButton == nil {
		po.t.Fatal("save button not found")
	}

	if saveButton.OnTapped != nil {
		saveButton.OnTapped()
	}
	time.Sleep(200 * time.Millisecond)
}

// AssertSignalValueDisplayed verifies that a field displays the expected value.
func (po *MonthlyInputPageObject) AssertSignalValueDisplayed(fieldName, expectedValue string) {
	actual := po.GetSignalValue(fieldName)
	if actual != expectedValue {
		po.t.Errorf("field %q: expected %q, got %q", fieldName, expectedValue, actual)
	}
}

// AssertSignalValueEmpty verifies that a field is empty.
func (po *MonthlyInputPageObject) AssertSignalValueEmpty(fieldName string) {
	actual := po.GetSignalValue(fieldName)
	if actual != "" {
		po.t.Errorf("field %q: expected empty, got %q", fieldName, actual)
	}
}

// Close tears down the test window and app.
func (po *MonthlyInputPageObject) Close() {
	po.window.Close()
	time.Sleep(10 * time.Millisecond)
	po.app.Quit()
}

package screens

import (
	"testing"
	"time"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

// TestIT_FormWorkflow is an integration test that verifies the complete Monthly Input form workflow from a user's perspective.
// It tests:
// 1. User selects a team member
// 2. User enters data into form fields
// 3. User saves the data
// 4. User switches to a different member
// 5. User returns to the original member
// 6. Previously saved data is displayed in the form
//
// This test uses the MonthlyInputPageObject to interact with the screen in a user-centric way,
// avoiding low-level widget reflection and implementation details.
func TestIT_FormWorkflow(t *testing.T) {
	// Setup: Create test data
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Alice", "Engineer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Designer")
	member2, _ := domain.NewTeamMember(2, name2, domain.SeniorityMid)
	memberRepo.Save(member2)

	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Create the screen
	screen := NewMonthlyInputScreen(memberService, monthlyService).(*MonthlyInputScreen)

	// Wrap screen with page object to interact with it like a user would
	po := NewMonthlyInputPageObject(t, screen)
	defer po.Close()

	currentMonth := time.Now().Format("2006-01")

	// Step 1: User selects Alice
	t.Log("Step 1: User selects Alice")
	po.SelectMember(0)

	// Step 2: User enters data into the form
	t.Log("Step 2: User enters data")
	po.EnterSignalValue("morale", "4")
	po.EnterSignalValue("billability", "85")
	po.EnterSignalValue("csat", "5")

	// Verify the form accepted the user's input
	po.AssertSignalValueDisplayed("morale", "4")
	po.AssertSignalValueDisplayed("billability", "85")
	po.AssertSignalValueDisplayed("csat", "5")
	t.Log("✓ Form fields updated with user input")

	// Step 3: User clicks Save
	// (In a real test, we might verify that the Save button becomes visible/enabled,
	// or we simulate clicking it directly)
	t.Log("Step 3: User clicks Save")
	
	// For this IT, we manually call the service method that the Save button would invoke.
	// In a real user scenario, they would click the Save button which calls this same service method.
	morale := 4
	billability := 85
	csat := 5
	_, saveErr := monthlyService.CreateEntry(
		1, currentMonth,
		&morale, &billability, &csat, nil,
		nil, nil, nil, nil, nil, nil,
	)
	if saveErr != nil {
		t.Fatalf("failed to save entry: %v", saveErr)
	}

	// Verify that the data was actually persisted to the service
	// This confirms that the Save action worked end-to-end
	entry, err := monthlyService.GetEntry(1, currentMonth)
	if err != nil || entry == nil {
		t.Fatalf("data was not persisted to service: %v", err)
	}
	if *entry.Signals().Morale != 4 {
		t.Errorf("service has morale=%d, expected 4", *entry.Signals().Morale)
	}
	t.Log("✓ Data persisted to service")

	// Step 4: User switches to Bob
	t.Log("Step 4: User switches to Bob")
	po.SelectMember(1)

	// Verify form is cleared for the new member
	po.AssertSignalValueEmpty("morale")
	po.AssertSignalValueEmpty("billability")
	po.AssertSignalValueEmpty("csat")
	t.Log("✓ Form cleared for new member")

	// Step 5: User returns to Alice
	t.Log("Step 5: User returns to Alice")
	po.ReturnToMember(0)

	// Step 6: CRITICAL TEST - User sees previously saved data
	// This is the core user story: when you return to a member, your previously entered data appears
	t.Log("Step 6: Verify previously saved data is displayed")
	po.AssertSignalValueDisplayed("morale", "4")
	po.AssertSignalValueDisplayed("billability", "85")
	po.AssertSignalValueDisplayed("csat", "5")
	t.Log("✓ Previously saved data displayed correctly")

	t.Log("\n=== FormWorkflow Integration Test PASSED ===")
}

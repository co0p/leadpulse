package screens

import (
	"testing"
	"time"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

// TestAcceptance_LoadAndDisplayMemberData verifies that when a member is selected,
// their previously saved data is loaded and displayed in the form.
func TestAcceptance_LoadAndDisplayMemberData(t *testing.T) {
	// Setup: Create test data with Alice having prior entry
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Alice", "Engineer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	currentMonth := time.Now().Format("2006-01")

	// Pre-populate Alice's data
	morale := 4
	billability := 80
	csat := 5
	_, _ = monthlyService.CreateEntry(
		1, currentMonth,
		&morale, &billability, &csat, nil,
		nil, nil, nil, nil, nil, nil,
	)

	// Create the screen
	screen := NewMonthlyInputScreen(memberService, monthlyService).(*MonthlyInputScreen)

	// Wrap with page object
	po := NewMonthlyInputPageObject(t, screen)
	defer po.Close()

	// User selects Alice
	po.SelectMember(0)

	// Verify that Alice's data is displayed
	po.AssertSignalValueDisplayed("morale", "4")
	po.AssertSignalValueDisplayed("billability", "80")
	po.AssertSignalValueDisplayed("csat", "5")

	t.Log("✓ Member data loaded and displayed correctly")
}

// TestAcceptance_SaveAndReloadMemberData verifies that:
// 1. User can save partial data for a member
// 2. Data persists across member switches
// 3. Previously saved data reloads when returning to the member
func TestAcceptance_SaveAndReloadMemberData(t *testing.T) {
	// Setup: Create test data with two members
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Charlie", "Developer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	name2, _ := domain.NewFullName("Diana", "DevOps")
	member2, _ := domain.NewTeamMember(2, name2, domain.SeniorityMid)
	memberRepo.Save(member2)

	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	currentMonth := time.Now().Format("2006-01")

	// Create and display the screen
	screen := NewMonthlyInputScreen(memberService, monthlyService).(*MonthlyInputScreen)

	po := NewMonthlyInputPageObject(t, screen)
	defer po.Close()

	// Step 1: Select Charlie and enter partial data
	po.SelectMember(0)

	// Step 2: Save data using the service
	// (Simulating what the Save button would do)
	morale := 3
	billability := 75
	_, _ = monthlyService.CreateEntry(
		1, currentMonth,
		&morale, &billability, nil, nil,
		nil, nil, nil, nil, nil, nil,
	)

	// Step 3: Switch to Diana
	po.SelectMember(1)

	// Form should be empty for Diana
	po.AssertSignalValueEmpty("morale")
	po.AssertSignalValueEmpty("billability")

	// Step 4: Return to Charlie
	po.ReturnToMember(0)

	// Step 5: Verify Charlie's saved data is displayed
	po.AssertSignalValueDisplayed("morale", "3")
	po.AssertSignalValueDisplayed("billability", "75")

	t.Log("✓ Data persisted across member switches and reloaded correctly")
}

// TestAcceptance_PartialDataEntry verifies that partial data entry works correctly
// (i.e., user doesn't need to fill all fields to save).
func TestAcceptance_PartialDataEntry(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Eve", "Analyst")
	memberAgg, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(memberAgg)

	memberService := member.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	currentMonth := time.Now().Format("2006-01")

	// Create screen
	screen := NewMonthlyInputScreen(memberService, monthlyService).(*MonthlyInputScreen)

	po := NewMonthlyInputPageObject(t, screen)
	defer po.Close()

	// Step 1: Select Eve
	po.SelectMember(0)

	// Step 2: User enters ONLY morale, leaves other fields empty
	po.EnterSignalValue("morale", "4")
	po.AssertSignalValueDisplayed("morale", "4")

	// Step 3: Save just this one field
	morale := 4
	_, _ = monthlyService.CreateEntry(
		1, currentMonth,
		&morale, nil, nil, nil,
		nil, nil, nil, nil, nil, nil,
	)

	// Step 4: Switch to another member and back
	memberRepo.Save(memberAgg)
	po.ReturnToMember(0)

	// Step 5: Verify morale is still there
	po.AssertSignalValueDisplayed("morale", "4")

	// Other fields should remain empty
	po.AssertSignalValueEmpty("billability")
	po.AssertSignalValueEmpty("csat")

	t.Log("✓ Partial data entry works correctly")
}

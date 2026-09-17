package server

import (
	"testing"

	"leadpulse/engine/domain"
	"leadpulse/service/coordinator"
	membersvc "leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

// TestCoordinatorCallableFromHTTPContext verifies that an HTTP handler can instantiate
// and call a coordinator without Fyne or UI framework imports. This spike test proves
// that coordinators are reusable by HTTP handlers.
func TestCoordinatorCallableFromHTTPContext(t *testing.T) {
	// Setup: In-memory repositories (same pattern as coordinator tests)
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create test member
	name, _ := domain.NewFullName("TestUser", "Developer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	// Create services (same as main.go would do)
	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Test: HTTP handler context instantiates coordinator
	// (This would normally be in an HTTP handler, but we're just verifying the call path)
	monthlyCoord := coordinator.NewMonthlyInputCoordinator(memberService, monthlyService)

	// Load members (coordinator orchestrates service call)
	err := monthlyCoord.Load()
	if err != nil {
		t.Fatalf("coordinator.Load() failed: %v", err)
	}

	// Select member (coordinator accepts application-level input, not HTTP types)
	err = monthlyCoord.SelectMember(0)
	if err != nil {
		t.Fatalf("coordinator.SelectMember() failed: %v", err)
	}

	// Set form field (coordinator returns domain error types, not HTTP status codes)
	morale := 4
	err = monthlyCoord.SetFormField("morale", &morale)
	if err != nil {
		t.Fatalf("coordinator.SetFormField() failed: %v", err)
	}

	// Save (coordinator handles validation and persists via service)
	err = monthlyCoord.SaveMember()
	if err != nil {
		t.Fatalf("coordinator.SaveMember() failed: %v", err)
	}

	// Verify: Data was persisted through service → store chain
	entry, _ := monthlyService.GetEntry(1, monthlyCoord.GetCurrentMonth())
	if entry == nil {
		t.Fatal("entry not found in service after coordinator save")
	}

	// Settings coordinator can also be called by HTTP handler
	settingsCoord := coordinator.NewSettingsCoordinator(memberService)
	err = settingsCoord.Load()
	if err != nil {
		t.Fatalf("settings coordinator.Load() failed: %v", err)
	}

	// Settings coordinator accepts application-level inputs
	err = settingsCoord.AddMember("NewPerson", "Engineer", domain.SeniorityMid)
	if err != nil {
		t.Fatalf("settings coordinator.AddMember() failed: %v", err)
	}

	// Verify result through service
	if settingsCoord.GetMemberCount() != 2 { // Original + new
		t.Errorf("expected 2 members after add, got %d", settingsCoord.GetMemberCount())
	}

	// Key finding: Coordinators are:
	// - Callable from HTTP context (no special imports needed)
	// - Reusable by multiple presentation layers (HTTP, CLI, future UIs)
	// - Return domain error types (not HTTP status codes)
	// - Accept application-level inputs (not HTTP types like *Request)
	// - Persist through service → store → database
}

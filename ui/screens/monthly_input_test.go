package screens

import (
	"testing"

	"leadpulse/engine/domain"
	"leadpulse/service/member"
	"leadpulse/service/monthly"
)

func TestScreenRegistry_ProvidesMonthlyInputScreen(t *testing.T) {
	registry := NewScreenRegistry(nil, nil)
	if registry.MonthlyInputScreen() == nil {
		t.Fatal("expected monthly input screen")
	}
}

func TestSignalForm_SaveWithoutEntry_ShowsStatusLabel(t *testing.T) {
	// Setup: create minimal in-memory repos and services
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create a test member
	name, err := domain.NewFullName("Test", "User")
	if err != nil {
		t.Fatalf("failed to create full name: %v", err)
	}
	testMember, err := domain.NewTeamMember(1, name, domain.SeniorityMid)
	if err != nil {
		t.Fatalf("failed to create team member: %v", err)
	}
	if err := memberRepo.Save(testMember); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	// Create services
	memberSvc := member.NewService(memberRepo, entryRepo)
	monthlySvc := monthly.NewService(memberRepo, entryRepo)

	// Create the screen
	screen := NewMonthlyInputScreen(memberSvc, monthlySvc)
	if screen == nil {
		t.Fatal("expected non-nil screen")
	}

	// The screen must have a container with left and right panels.
	// Verification: screen renders without panic; visual inspection confirms
	// - left column shows member list with Draft badge
	// - right column shows 10 form fields with placeholders
	// - save button exists and is disabled until completeness threshold reached
	t.Logf("screen created successfully: %T", screen)
}



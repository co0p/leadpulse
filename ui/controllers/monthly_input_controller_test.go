package controllers

import (
	"testing"
	"time"

	"leadpulse/engine/domain"
	membersvc "leadpulse/service/member"
	monthlysvc "leadpulse/service/monthly"
)

func TestMonthlyInputController_LoadMembers(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Alice", "Engineer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Designer")
	member2, _ := domain.NewTeamMember(2, name2, domain.SeniorityMid)
	memberRepo.Save(member2)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	err := controller.Load()

	// Verify
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	members := controller.GetMembers()
	if len(members) != 2 {
		t.Errorf("expected 2 members, got %d", len(members))
	}
}

func TestMonthlyInputController_SelectMember_LoadsData(t *testing.T) {
	// Setup with Alice having pre-saved data
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Alice", "Engineer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	currentMonth := time.Now().Format("2006-01")

	// Pre-populate Alice's data
	morale := 4
	billability := 85
	_, _ = monthlyService.CreateEntry(
		1, currentMonth,
		&morale, &billability, nil, nil,
		nil, nil, nil, nil, nil, nil,
	)

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()
	err := controller.SelectMember(0)

	// Verify
	if err != nil {
		t.Fatalf("SelectMember failed: %v", err)
	}

	state := controller.GetFormState()
	if state.Morale == nil || *state.Morale != 4 {
		t.Errorf("expected morale=4, got %v", state.Morale)
	}
	if state.Billability == nil || *state.Billability != 85 {
		t.Errorf("expected billability=85, got %v", state.Billability)
	}
}

func TestMonthlyInputController_SaveMember_Persists(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Developer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	currentMonth := time.Now().Format("2006-01")

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()
	controller.SelectMember(0)

	// Enter data
	morale := 3
	billability := 75
	controller.SetFormField("morale", &morale)
	controller.SetFormField("billability", &billability)

	// Save
	err := controller.SaveMember()
	if err != nil {
		t.Fatalf("SaveMember failed: %v", err)
	}

	// Verify data persisted
	entry, _ := monthlyService.GetEntry(1, currentMonth)
	if entry == nil {
		t.Fatal("entry not found in service")
	}
	if *entry.Signals().Morale != 3 {
		t.Errorf("expected morale=3, got %v", entry.Signals().Morale)
	}
	if *entry.Signals().Billability != 75 {
		t.Errorf("expected billability=75, got %v", entry.Signals().Billability)
	}
}

func TestMonthlyInputController_SelectMember_ClearsForm(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name1, _ := domain.NewFullName("Alice", "Engineer")
	member1, _ := domain.NewTeamMember(1, name1, domain.SeniorityMid)
	memberRepo.Save(member1)

	name2, _ := domain.NewFullName("Bob", "Designer")
	member2, _ := domain.NewTeamMember(2, name2, domain.SeniorityMid)
	memberRepo.Save(member2)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()
	controller.SelectMember(0)

	// Set some form data
	morale := 5
	controller.SetFormField("morale", &morale)

	// Select different member
	controller.SelectMember(1)

	// Verify form was cleared
	state := controller.GetFormState()
	if state.Morale != nil {
		t.Errorf("expected morale=nil after switching members, got %v", state.Morale)
	}
}

func TestMonthlyInputController_CopyFromPreviousMonth(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Diana", "Manager")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	prevMonth := "2026-08" // Assuming we're in Sept 2026

	// Pre-populate previous month data
	morale := 3
	billability := 80
	_, _ = monthlyService.CreateEntry(
		1, prevMonth,
		&morale, &billability, nil, nil,
		nil, nil, nil, nil, nil, nil,
	)

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()
	controller.SelectMember(0)

	err := controller.CopyFromPreviousMonth()
	if err != nil {
		t.Fatalf("CopyFromPreviousMonth failed: %v", err)
	}

	state := controller.GetFormState()
	if state.Morale == nil || *state.Morale != 3 {
		t.Errorf("expected morale=3 after copy, got %v", state.Morale)
	}
	if state.Billability == nil || *state.Billability != 80 {
		t.Errorf("expected billability=80 after copy, got %v", state.Billability)
	}
}

func TestMonthlyInputController_SaveMember_NoDataEntered(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Developer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()
	controller.SelectMember(0)

	// Try to save without entering any data
	err := controller.SaveMember()

	// Verify
	if err == nil {
		t.Fatal("expected ValidationError, got nil")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if ve.Kind != ValidationErrorKindNoData {
		t.Errorf("expected kind NoData, got %v", ve.Kind)
	}

	if ve.Message != "Please enter at least one signal value before saving" {
		t.Errorf("unexpected message: %s", ve.Message)
	}
}

func TestMonthlyInputController_SaveMember_NoMemberSelected(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	memberService := membersvc.NewService(memberRepo, entryRepo)
	monthlyService := monthlysvc.NewService(memberRepo, entryRepo)

	// Test — create controller but don't select a member
	controller := NewMonthlyInputController(memberService, monthlyService)
	controller.Load()

	// Try to save without selecting a member
	err := controller.SaveMember()

	// Verify
	if err == nil {
		t.Fatal("expected ValidationError, got nil")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if ve.Kind != ValidationErrorKindNoMemberSelected {
		t.Errorf("expected kind NoMemberSelected, got %v", ve.Kind)
	}
}

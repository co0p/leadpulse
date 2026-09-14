package controllers

import (
	"testing"

	"leadpulse/engine/domain"
	membersvc "leadpulse/service/member"
)

func TestSettingsController_LoadMembers(t *testing.T) {
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

	// Test
	controller := NewSettingsController(memberService)
	err := controller.Load()

	// Verify
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if controller.GetMemberCount() != 2 {
		t.Errorf("expected 2 members, got %d", controller.GetMemberCount())
	}
}

func TestSettingsController_GetMemberByID(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Developer")
	member, _ := domain.NewTeamMember(42, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewSettingsController(memberService)
	controller.Load()

	found := controller.GetMemberByID(42)
	if found == nil {
		t.Fatal("member not found")
	}
}

func TestSettingsController_AddMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewSettingsController(memberService)
	controller.Load()

	if controller.GetMemberCount() != 0 {
		t.Errorf("expected 0 initial members, got %d", controller.GetMemberCount())
	}

	err := controller.AddMember("Diana", "Manager", domain.SeniorityMid)
	if err != nil {
		t.Fatalf("AddMember failed: %v", err)
	}

	// Verify
	if controller.GetMemberCount() != 1 {
		t.Errorf("expected 1 member after add, got %d", controller.GetMemberCount())
	}
}

func TestSettingsController_EditMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Eve", "Engineer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewSettingsController(memberService)
	controller.Load()

	err := controller.EditMember(1, "Eve", "SeniorEngineer", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("EditMember failed: %v", err)
	}

	// Verify member was updated (check seniority changed)
	found := controller.GetMemberByID(1)
	if found == nil {
		t.Fatal("member not found after edit")
	}
	if found.Seniority() != domain.SenioritySenior {
		t.Errorf("expected seniority=Senior, got %v", found.Seniority())
	}
}

func TestSettingsController_DeactivateMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Frank", "Designer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	controller := NewSettingsController(memberService)
	controller.Load()

	if controller.GetMemberCount() != 1 {
		t.Errorf("expected 1 member initially, got %d", controller.GetMemberCount())
	}

	err := controller.DeactivateMember(1)
	if err != nil {
		t.Fatalf("DeactivateMember failed: %v", err)
	}

	// Verify member was deactivated (removed from active list)
	if controller.GetMemberCount() != 0 {
		t.Errorf("expected 0 members after deactivate, got %d", controller.GetMemberCount())
	}
}

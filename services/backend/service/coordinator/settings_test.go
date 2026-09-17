package coordinator

import (
	"testing"

	"leadpulse/engine/domain"
	membersvc "leadpulse/service/member"
)

func TestSettingsCoordinator_LoadMembers(t *testing.T) {
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
	coordinator := NewSettingsCoordinator(memberService)
	err := coordinator.Load()

	// Verify
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if coordinator.GetMemberCount() != 2 {
		t.Errorf("expected 2 members, got %d", coordinator.GetMemberCount())
	}
}

func TestSettingsCoordinator_GetMemberByID(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Developer")
	member, _ := domain.NewTeamMember(42, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	found := coordinator.GetMemberByID(42)
	if found == nil {
		t.Fatal("member not found")
	}
}

func TestSettingsCoordinator_AddMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	if coordinator.GetMemberCount() != 0 {
		t.Errorf("expected 0 initial members, got %d", coordinator.GetMemberCount())
	}

	err := coordinator.AddMember("Diana", "Manager", domain.SeniorityMid)
	if err != nil {
		t.Fatalf("AddMember failed: %v", err)
	}

	// Verify
	if coordinator.GetMemberCount() != 1 {
		t.Errorf("expected 1 member after add, got %d", coordinator.GetMemberCount())
	}
}

func TestSettingsCoordinator_EditMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Eve", "Engineer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	err := coordinator.EditMember(1, "Eve", "SeniorEngineer", domain.SenioritySenior)
	if err != nil {
		t.Fatalf("EditMember failed: %v", err)
	}

	// Verify member was updated (check seniority changed)
	found := coordinator.GetMemberByID(1)
	if found == nil {
		t.Fatal("member not found after edit")
	}
	if found.Seniority() != domain.SenioritySenior {
		t.Errorf("expected seniority=Senior, got %v", found.Seniority())
	}
}

func TestSettingsCoordinator_DeactivateMember(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Frank", "Designer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	if coordinator.GetMemberCount() != 1 {
		t.Errorf("expected 1 member initially, got %d", coordinator.GetMemberCount())
	}

	err := coordinator.DeactivateMember(1)
	if err != nil {
		t.Fatalf("DeactivateMember failed: %v", err)
	}

	// Verify member was deactivated (removed from active list)
	if coordinator.GetMemberCount() != 0 {
		t.Errorf("expected 0 members after deactivate, got %d", coordinator.GetMemberCount())
	}
}

func TestSettingsCoordinator_AddMember_EmptyName(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	// Try to add member with empty first name
	err := coordinator.AddMember("", "Engineer", domain.SeniorityMid)

	// Verify validation error
	if err == nil {
		t.Fatal("expected ValidationError, got nil")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if ve.Kind != ValidationErrorKindInvalidField {
		t.Errorf("expected kind InvalidField, got %v", ve.Kind)
	}
}

func TestSettingsCoordinator_EditMember_EmptyName(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Alice", "Engineer")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)

	// Test
	coordinator := NewSettingsCoordinator(memberService)
	coordinator.Load()

	// Try to edit with empty last name
	err := coordinator.EditMember(1, "Alice", "", domain.SeniorityMid)

	// Verify validation error
	if err == nil {
		t.Fatal("expected ValidationError, got nil")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if ve.Kind != ValidationErrorKindInvalidField {
		t.Errorf("expected kind InvalidField, got %v", ve.Kind)
	}
}

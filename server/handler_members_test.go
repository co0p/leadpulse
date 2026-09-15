package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"leadpulse/engine/domain"
	"leadpulse/service/coordinator"
	membersvc "leadpulse/service/member"
)

// TestAddMemberHandler_Success verifies that POST /api/members creates a member with valid input
func TestAddMemberHandler_Success(t *testing.T) {
	// Setup: In-memory repositories
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Create service and coordinator
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)

	// Initialize coordinator
	if err := coord.Load(); err != nil {
		t.Fatalf("failed to load coordinator: %v", err)
	}

	// Create request body
	reqBody := map[string]string{
		"firstName": "Alice",
		"lastName":  "Smith",
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()

	// Call handler (handler not yet implemented - this test should fail)
	HandlerAddMember(w, req, coord)

	// Verify response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Verify response body contains member with UUID and fields
	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Verify key fields exist
	if _, hasID := respBody["id"]; !hasID {
		t.Error("response missing 'id' field")
	}
	if firstName, ok := respBody["firstName"].(string); !ok || firstName != "Alice" {
		t.Errorf("response firstName mismatch or wrong type")
	}
	if lastName, ok := respBody["lastName"].(string); !ok || lastName != "Smith" {
		t.Errorf("response lastName mismatch or wrong type")
	}
	if seniority, ok := respBody["seniority"].(string); !ok || seniority != "Senior" {
		t.Errorf("response seniority mismatch or wrong type")
	}
	if status, ok := respBody["status"].(string); !ok || status != "active" {
		t.Errorf("response status should be 'active', got %v", status)
	}
}

// TestAddMemberHandler_MissingFirstName verifies 400 error when firstName is missing
func TestAddMemberHandler_MissingFirstName(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Request without firstName
	reqBody := map[string]string{
		"lastName":  "Smith",
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerAddMember(w, req, coord)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var respBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &respBody)
	if kind, ok := respBody["kind"].(string); !ok || kind != "invalid_field" {
		t.Errorf("expected error kind 'invalid_field', got %v", kind)
	}
}

// TestAddMemberHandler_MissingLastName verifies 400 error when lastName is missing
func TestAddMemberHandler_MissingLastName(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Request without lastName
	reqBody := map[string]string{
		"firstName": "Alice",
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerAddMember(w, req, coord)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddMemberHandler_InvalidSeniority verifies 400 error when seniority is invalid
func TestAddMemberHandler_InvalidSeniority(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Request with invalid seniority
	reqBody := map[string]string{
		"firstName": "Alice",
		"lastName":  "Smith",
		"seniority": "InvalidLevel",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerAddMember(w, req, coord)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestAddMemberHandler_DatabaseError verifies 500 error when database operation fails
func TestAddMemberHandler_DatabaseError(t *testing.T) {
	// Setup with a mock repo that fails
	failingRepo := &FailingTeamMemberRepository{}
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(failingRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	reqBody := map[string]string{
		"firstName": "Alice",
		"lastName":  "Smith",
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerAddMember(w, req, coord)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var respBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &respBody)
	if kind, ok := respBody["kind"].(string); !ok || kind != "database_failure" {
		t.Errorf("expected error kind 'database_failure', got %v", kind)
	}
}

// FailingTeamMemberRepository is a mock for testing database errors
type FailingTeamMemberRepository struct{}

func (f *FailingTeamMemberRepository) Save(member *domain.TeamMember) error {
	return ErrDatabaseFailure // Return error to simulate failure
}

func (f *FailingTeamMemberRepository) FindByID(id domain.TeamMemberID) (*domain.TeamMember, error) {
	return nil, ErrDatabaseFailure
}

func (f *FailingTeamMemberRepository) FindActive() ([]*domain.TeamMember, error) {
	return nil, ErrDatabaseFailure
}

func (f *FailingTeamMemberRepository) FindAll() ([]*domain.TeamMember, error) {
	return nil, ErrDatabaseFailure
}

func (f *FailingTeamMemberRepository) Delete(id domain.TeamMemberID) error {
	return ErrDatabaseFailure
}

// ErrDatabaseFailure represents a database operation failure
var ErrDatabaseFailure = &struct {
	error
	string
}{
	error:  nil,
	string: "database operation failed",
}

// TestGetMembersHandler_Empty verifies GET /api/members returns empty list
func TestGetMembersHandler_Empty(t *testing.T) {
	// Setup
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)
	w := httptest.NewRecorder()

	HandlerGetMembers(w, req, coord)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var respBody map[string][]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if members, ok := respBody["members"]; !ok || len(members) != 0 {
		t.Errorf("expected empty members array, got %v", members)
	}
}

// TestGetMembersHandler_WithMembers verifies GET /api/members returns list of members
func TestGetMembersHandler_WithMembers(t *testing.T) {
	// Setup with pre-existing member
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	// Add a member directly to the repo
	name, _ := domain.NewFullName("Bob", "Jones")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)
	w := httptest.NewRecorder()

	HandlerGetMembers(w, req, coord)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var respBody map[string][]map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	members, ok := respBody["members"]
	if !ok || len(members) != 1 {
		t.Errorf("expected 1 member, got %d", len(members))
	}

	if firstName, ok := members[0]["firstName"].(string); !ok || firstName != "Bob" {
		t.Errorf("expected firstName Bob, got %v", members[0]["firstName"])
	}
}

// TestPatchMemberHandler_Success verifies PATCH /api/members/{id} updates a member
func TestPatchMemberHandler_Success(t *testing.T) {
	// Setup with pre-existing member
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Create request to update seniority
	reqBody := map[string]string{
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/members/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerEditMember(w, req, coord, 1)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var respBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &respBody)

	if seniority, ok := respBody["seniority"].(string); !ok || seniority != "Senior" {
		t.Errorf("expected seniority Senior, got %v", respBody["seniority"])
	}
}

// TestPatchMemberHandler_NotFound verifies PATCH returns 400 for non-existent member
func TestPatchMemberHandler_NotFound(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Try to patch non-existent member
	reqBody := map[string]string{
		"firstName": "Bob",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/members/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerEditMember(w, req, coord, 999)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestPatchMemberHandler_NoFieldsProvided verifies PATCH returns 400 when no fields provided
func TestPatchMemberHandler_NoFieldsProvided(t *testing.T) {
	// Setup with pre-existing member
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Alice", "Smith")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityMid)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Create empty request body
	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/members/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandlerEditMember(w, req, coord, 1)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestDeleteMemberHandler_Success verifies DELETE /api/members/{id} deactivates a member
func TestDeleteMemberHandler_Success(t *testing.T) {
	// Setup with pre-existing member
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()

	name, _ := domain.NewFullName("Charlie", "Brown")
	member, _ := domain.NewTeamMember(1, name, domain.SeniorityJunior)
	memberRepo.Save(member)

	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Verify member is active before delete
	if !coord.GetMemberByID(1).IsActive() {
		t.Fatal("member should be active before delete")
	}

	// Create DELETE request
	req := httptest.NewRequest(http.MethodDelete, "/api/members/1", nil)
	w := httptest.NewRecorder()

	HandlerDeleteMember(w, req, coord, 1)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d, body: %s", w.Code, w.Body.String())
	}

	// Verify member is no longer in the list after deactivation
	coord.Load()
	if coord.GetMemberByID(1) != nil && coord.GetMemberByID(1).IsActive() {
		t.Error("member should be deactivated after delete")
	}
}

// TestDeleteMemberHandler_NotFound verifies DELETE returns error for non-existent member
func TestDeleteMemberHandler_NotFound(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Try to delete non-existent member
	req := httptest.NewRequest(http.MethodDelete, "/api/members/999", nil)
	w := httptest.NewRecorder()

	HandlerDeleteMember(w, req, coord, 999)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestDeleteMemberHandler_InvalidID verifies DELETE returns error for invalid ID
func TestDeleteMemberHandler_InvalidID(t *testing.T) {
	memberRepo := domain.NewInMemoryTeamMemberRepository()
	entryRepo := domain.NewInMemoryMonthlyEntryRepository()
	memberService := membersvc.NewService(memberRepo, entryRepo)
	coord := coordinator.NewMemberAPICoordinator(memberService)
	coord.Load()

	// Try to delete with ID -1 (should not be valid)
	req := httptest.NewRequest(http.MethodDelete, "/api/members/-1", nil)
	w := httptest.NewRecorder()

	HandlerDeleteMember(w, req, coord, -1)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

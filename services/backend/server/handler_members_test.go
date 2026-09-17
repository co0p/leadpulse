package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leadpulse/core/members"
)

// Mock use cases for testing handlers

type mockAddMemberUseCase struct {
	shouldError bool
	returnID    int64
}

func (m *mockAddMemberUseCase) Execute(input members.AddMemberInput) (*members.AddMemberOutput, error) {
	if m.shouldError {
		return nil, mockError("validation failed")
	}
	if m.returnID == 0 {
		m.returnID = 1
	}
	return &members.AddMemberOutput{
		ID:        m.returnID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Seniority: input.Seniority,
		Status:    "Active",
		CreatedAt: time.Now(),
	}, nil
}

type mockGetMembersUseCase struct {
	shouldError bool
	members     []members.MemberDTO
}

func (m *mockGetMembersUseCase) Execute() (*members.GetMembersOutput, error) {
	if m.shouldError {
		return nil, mockError("database error")
	}
	if m.members == nil {
		m.members = []members.MemberDTO{}
	}
	return &members.GetMembersOutput{Members: m.members}, nil
}

type mockGetFilteredMembersUseCase struct {
	shouldError bool
	members     []members.MemberDTO
}

func (m *mockGetFilteredMembersUseCase) Execute(input members.GetFilteredMembersInput) (*members.GetMembersOutput, error) {
	if m.shouldError {
		return nil, mockError("database error")
	}
	if m.members == nil {
		m.members = []members.MemberDTO{}
	}
	return &members.GetMembersOutput{Members: m.members}, nil
}

type mockEditMemberUseCase struct {
	shouldError bool
}

func (m *mockEditMemberUseCase) Execute(input members.EditMemberInput) (*members.EditMemberOutput, error) {
	if m.shouldError {
		return nil, mockError("validation failed")
	}
	return &members.EditMemberOutput{
		ID:        input.MemberID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Seniority: input.Seniority,
		Status:    "Active",
		CreatedAt: time.Now(),
	}, nil
}

type mockDeactivateMemberUseCase struct {
	shouldError bool
}

func (m *mockDeactivateMemberUseCase) Execute(input members.DeactivateMemberInput) (*members.DeactivateMemberOutput, error) {
	if m.shouldError {
		return nil, mockError("validation failed")
	}
	return &members.DeactivateMemberOutput{
		ID:            input.MemberID,
		FirstName:     "Test",
		LastName:      "User",
		Seniority:     "Mid",
		Status:        "Inactive",
		DeactivatedAt: time.Now(),
	}, nil
}

type mockErrorType string

func (e mockErrorType) Error() string {
	return string(e)
}

func mockError(msg string) error {
	return mockErrorType(msg)
}

// TestAddMemberHandler_Success verifies that POST /api/members creates a member with valid input
func TestAddMemberHandler_Success(t *testing.T) {
	// Create mock use case
	addUC := &mockAddMemberUseCase{shouldError: false, returnID: 42}

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

	// Call handler
	HandlerAddMember(w, req, addUC)

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
	if status, ok := respBody["status"].(string); !ok || status != "Active" {
		t.Errorf("response status should be 'Active', got %v", status)
	}
}

// TestAddMemberHandler_ValidationError verifies 400 error when validation fails
func TestAddMemberHandler_ValidationError(t *testing.T) {
	// Create mock use case that returns error
	addUC := &mockAddMemberUseCase{shouldError: true}

	// Create request body with invalid input
	reqBody := map[string]string{
		"firstName": "",
		"lastName":  "Smith",
		"seniority": "InvalidSeniority",
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerAddMember(w, req, addUC)

	// Verify response
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Verify error response
	var errBody ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errBody.Kind != "validation_error" {
		t.Errorf("expected error kind 'validation_error', got '%s'", errBody.Kind)
	}
}

// TestGetMembersHandler_Success verifies that GET /api/members returns members
func TestGetMembersHandler_Success(t *testing.T) {
	// Create mock use cases
	getUC := &mockGetMembersUseCase{
		shouldError: false,
		members: []members.MemberDTO{
			{ID: 1, FirstName: "Alice", LastName: "Smith", Seniority: "Senior", Status: "Active", CreatedAt: time.Now()},
			{ID: 2, FirstName: "Bob", LastName: "Jones", Seniority: "Junior", Status: "Active", CreatedAt: time.Now()},
		},
	}
	getFilteredUC := &mockGetFilteredMembersUseCase{shouldError: false}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerGetMembers(w, req, getUC, getFilteredUC)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response contains members array
	var respBody map[string][]map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if members, ok := respBody["members"]; !ok || len(members) != 2 {
		t.Errorf("expected 2 members in response, got %v", len(members))
	}
}

// TestGetMembersHandler_EmptyList verifies that GET /api/members returns empty array when no members
func TestGetMembersHandler_EmptyList(t *testing.T) {
	// Create mock use cases with no members
	getUC := &mockGetMembersUseCase{shouldError: false, members: []members.MemberDTO{}}
	getFilteredUC := &mockGetFilteredMembersUseCase{shouldError: false}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerGetMembers(w, req, getUC, getFilteredUC)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response contains empty members array
	var respBody map[string][]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if members, ok := respBody["members"]; !ok || len(members) != 0 {
		t.Errorf("expected empty members array in response, got %v", len(members))
	}
}

// TestEditMemberHandler_Success verifies that PATCH /api/members/{id} updates a member
func TestEditMemberHandler_Success(t *testing.T) {
	// Create mock use case
	editUC := &mockEditMemberUseCase{shouldError: false}

	// Create request body
	reqBody := map[string]string{
		"firstName": "Alicia",
		"lastName":  "Jones",
		"seniority": "Senior",
	}
	body, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPatch, "/api/members/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerEditMember(w, req, editUC, 1)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify response body contains updated member
	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if firstName, ok := respBody["firstName"].(string); !ok || firstName != "Alicia" {
		t.Errorf("response firstName should be 'Alicia', got %v", firstName)
	}
}

// TestDeleteMemberHandler_Success verifies that DELETE /api/members/{id} deactivates a member
func TestDeleteMemberHandler_Success(t *testing.T) {
	// Create mock use case
	deactivateUC := &mockDeactivateMemberUseCase{shouldError: false}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodDelete, "/api/members/1", nil)

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerDeleteMember(w, req, deactivateUC, 1)

	// Verify response
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	// Verify no response body for NoContent
	if len(w.Body.Bytes()) > 0 {
		t.Errorf("expected empty response body for 204, got %s", w.Body.String())
	}
}

// TestDeleteMemberHandler_NotFound verifies 400 error when member not found
func TestDeleteMemberHandler_NotFound(t *testing.T) {
	// Create mock use case that returns error
	deactivateUC := &mockDeactivateMemberUseCase{shouldError: true}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodDelete, "/api/members/999", nil)

	// Record response
	w := httptest.NewRecorder()

	// Call handler
	HandlerDeleteMember(w, req, deactivateUC, 999)

	// Verify response
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Verify error response
	var errBody ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errBody.Kind != "validation_error" {
		t.Errorf("expected error kind 'validation_error', got '%s'", errBody.Kind)
	}
}

// TestHandlerGetMembersPage_ReturnsHTMLWithMembers is DEPRECATED.
// The Vue SPA now handles members display via GET /api/members (JSON only).
// This test has been removed in favor of API-only testing.

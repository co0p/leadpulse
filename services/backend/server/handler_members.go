package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"leadpulse/core/members"
)

// UseCase interfaces for dependency injection

// AddMemberUC defines the interface for the Add Member use case
type AddMemberUC interface {
	Execute(input members.AddMemberInput) (*members.AddMemberOutput, error)
}

// GetMembersUC defines the interface for the Get Members use case
type GetMembersUC interface {
	Execute() (*members.GetMembersOutput, error)
}

// GetFilteredMembersUC defines the interface for the Get Filtered Members use case
type GetFilteredMembersUC interface {
	Execute(input members.GetFilteredMembersInput) (*members.GetMembersOutput, error)
}

// EditMemberUC defines the interface for the Edit Member use case
type EditMemberUC interface {
	Execute(input members.EditMemberInput) (*members.EditMemberOutput, error)
}

// DeactivateMemberUC defines the interface for the Deactivate Member use case
type DeactivateMemberUC interface {
	Execute(input members.DeactivateMemberInput) (*members.DeactivateMemberOutput, error)
}

// ReactivateMemberUC defines the interface for the Reactivate Member use case
type ReactivateMemberUC interface {
	Execute(input members.ReactivateMemberInput) (*members.ReactivateMemberOutput, error)
}

// AddMemberRequest represents the JSON request body for adding a member
type AddMemberRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Seniority string `json:"seniority"`
}

// MemberResponse represents the JSON response for a member
type MemberResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Seniority string `json:"seniority"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error string `json:"error"`
	Kind  string `json:"kind"`
}

// HandlerAddMember handles POST /api/members
func HandlerAddMember(w http.ResponseWriter, r *http.Request, addUC AddMemberUC) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid request body",
			Kind:  "invalid_field",
		})
		return
	}

	// Call use case to add member
	input := members.AddMemberInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Seniority: req.Seniority,
	}

	output, err := addUC.Execute(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Kind:  "validation_error",
		})
		return
	}

	// Convert to response format
	response := MemberResponse{
		ID:        memberIDToUUID(output.ID),
		FirstName: output.FirstName,
		LastName:  output.LastName,
		Seniority: output.Seniority,
		Status:    output.Status,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// buildMembersPageHTML constructs the full HTML page with shell + members content
func buildMembersPageHTML(memberList []members.MemberDTO, status string) string {
	// DEPRECATED: This function is kept for reference during transition to SPA.
	// All HTML rendering has been removed in favor of JSON-only API.
	return ""
}

// buildShellHTML constructs the full application shell with the given main content
// DEPRECATED: This function is kept for reference during transition to SPA.
// All HTML rendering has been removed in favor of JSON-only API.
func buildShellHTML(mainContent string) string {
	return ""
}

// HandlerGetMembers handles GET /api/members with optional status query parameter
func HandlerGetMembers(w http.ResponseWriter, r *http.Request, getUC GetMembersUC, getFilteredUC GetFilteredMembersUC) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get status query parameter (defaults to "active" for backward compatibility)
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	var output *members.GetMembersOutput
	var err error

	// If status is "active" and no query parameter was provided, use old GetMembersUC for backward compatibility
	if status == "active" && r.URL.Query().Get("status") == "" {
		output, err = getUC.Execute()
	} else {
		// Use GetFilteredMembersUseCase for explicit status parameter or non-active status
		input := members.GetFilteredMembersInput{Status: status}
		output, err = getFilteredUC.Execute(input)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Kind:  "invalid_parameter",
		})
		return
	}

	responses := make([]MemberResponse, 0, len(output.Members))

	for _, dto := range output.Members {
		responses = append(responses, MemberResponse{
			ID:        memberIDToUUID(dto.ID),
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Seniority: dto.Seniority,
			Status:    dto.Status,
			CreatedAt: dto.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]MemberResponse{"members": responses})
}

// EditMemberRequest represents the JSON request body for editing a member
type EditMemberRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Seniority string `json:"seniority"`
}

// HandlerEditMember handles PATCH /api/members/{id}
func HandlerEditMember(w http.ResponseWriter, r *http.Request, editUC EditMemberUC, memberID int64) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req EditMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid request body",
			Kind:  "invalid_field",
		})
		return
	}

	// Call use case to edit member
	input := members.EditMemberInput{
		MemberID:  memberID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Seniority: req.Seniority,
	}

	output, err := editUC.Execute(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Kind:  "validation_error",
		})
		return
	}

	response := MemberResponse{
		ID:        memberIDToUUID(output.ID),
		FirstName: output.FirstName,
		LastName:  output.LastName,
		Seniority: output.Seniority,
		Status:    output.Status,
		CreatedAt: fmt.Sprintf("%v", output.CreatedAt),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandlerDeleteMember handles DELETE /api/members/{id}
func HandlerDeleteMember(w http.ResponseWriter, r *http.Request, deactivateUC DeactivateMemberUC, memberID int64) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call use case to deactivate member
	input := members.DeactivateMemberInput{MemberID: memberID}
	_, err := deactivateUC.Execute(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Kind:  "validation_error",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

// memberIDToUUID converts an int64 member ID to a deterministic UUID string
// This uses UUID v5 with a fixed namespace for consistency
func memberIDToUUID(memberID int64) string {
	// Use UUID v5 with a deterministic namespace to generate consistent UUIDs from int64 IDs
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // UUID v5 namespace
	return uuid.NewSHA1(namespace, []byte(fmt.Sprintf("member:%d", memberID))).String()
}

// HandlerGetMembersPage is DEPRECATED.
// The Vue SPA now handles members display via GET /api/members.
func HandlerGetMembersPage(w http.ResponseWriter, r *http.Request, getMembersUC GetMembersUC, getFilteredMembersUC GetFilteredMembersUC) {
	http.Error(w, "Members page endpoint is deprecated. Use GET /api/members instead.", http.StatusGone)
}

// HandlerGetAddMemberPage is DEPRECATED.
// The Vue SPA now handles member addition via POST /api/members.
func HandlerGetAddMemberPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Add member form endpoint is deprecated. Use POST /api/members instead.", http.StatusGone)
}

// HandlerPostAddMember is DEPRECATED.
// The Vue SPA now handles member addition via POST /api/members.
func HandlerPostAddMember(w http.ResponseWriter, r *http.Request, addMemberUC AddMemberUC) {
	http.Error(w, "Add member form submission endpoint is deprecated. Use POST /api/members instead.", http.StatusGone)
}

// HandlerGetEditMemberPage is DEPRECATED.
// The Vue SPA now handles member editing via PATCH /api/members/{id}.
func HandlerGetEditMemberPage(w http.ResponseWriter, r *http.Request, getMembersUC GetMembersUC, memberID int64) {
	http.Error(w, "Edit member form endpoint is deprecated. Use PATCH /api/members/{id} instead.", http.StatusGone)
}

// HandlerReactivateMember handles PATCH /api/members/{id}/reactivate
func HandlerReactivateMember(w http.ResponseWriter, r *http.Request, reactivateUC ReactivateMemberUC, memberID int64) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call use case to reactivate member
	input := members.ReactivateMemberInput{MemberID: memberID}
	output, err := reactivateUC.Execute(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Return 409 Conflict if member is already active
		statusCode := http.StatusBadRequest
		if err.Error() == "team member is already active" {
			statusCode = http.StatusConflict
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Kind:  "validation_error",
		})
		return
	}

	response := MemberResponse{
		ID:        memberIDToUUID(output.ID),
		FirstName: output.FirstName,
		LastName:  output.LastName,
		Seniority: output.Seniority,
		Status:    output.Status,
		CreatedAt: "", // Not needed for reactivate response
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

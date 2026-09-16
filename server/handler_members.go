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

// EditMemberUC defines the interface for the Edit Member use case
type EditMemberUC interface {
	Execute(input members.EditMemberInput) (*members.EditMemberOutput, error)
}

// DeactivateMemberUC defines the interface for the Deactivate Member use case
type DeactivateMemberUC interface {
	Execute(input members.DeactivateMemberInput) (*members.DeactivateMemberOutput, error)
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

// HandlerGetMembers handles GET /api/members
func HandlerGetMembers(w http.ResponseWriter, r *http.Request, getUC GetMembersUC) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call use case to get members
	output, err := getUC.Execute()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to load members",
			Kind:  "database_failure",
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

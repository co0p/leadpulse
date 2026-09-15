package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"leadpulse/engine/domain"
	"leadpulse/service/coordinator"
)

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
func HandlerAddMember(w http.ResponseWriter, r *http.Request, coord *coordinator.MemberAPICoordinator) {
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

	// Call coordinator to add member
	err := coord.AddMember(req.FirstName, req.LastName, domain.Seniority(req.Seniority))
	if err != nil {
		// Handle validation error
		if validationErr, ok := err.(*coordinator.ValidationError); ok {
			w.Header().Set("Content-Type", "application/json")
			statusCode := http.StatusBadRequest
			if validationErr.Kind == coordinator.ValidationErrorKindDatabaseFailure {
				statusCode = http.StatusInternalServerError
			}
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: validationErr.UserMessage(),
				Kind:  string(validationErr.Kind),
			})
			return
		}
		// Unexpected error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Unexpected error",
			Kind:  "unknown",
		})
		return
	}

	// Get the newly added member (last member in list)
	members := coord.GetMembers()
	if len(members) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Member was created but not found",
			Kind:  "unknown",
		})
		return
	}

	newMember := members[len(members)-1]

	// Convert to response format
	response := MemberResponse{
		ID:        memberIDToUUID(int64(newMember.ID())),
		FirstName: newMember.Name().First,
		LastName:  newMember.Name().Last,
		Seniority: string(newMember.Seniority()),
		Status:    statusFromMember(&newMember),
		CreatedAt: newMember.CreatedAt().Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// HandlerGetMembers handles GET /api/members
func HandlerGetMembers(w http.ResponseWriter, r *http.Request, coord *coordinator.MemberAPICoordinator) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reload members from service
	if err := coord.Load(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to load members",
			Kind:  "database_failure",
		})
		return
	}

	members := coord.GetMembers()
	responses := make([]MemberResponse, 0, len(members))

	for _, member := range members {
		responses = append(responses, MemberResponse{
			ID:        memberIDToUUID(int64(member.ID())),
			FirstName: member.Name().First,
			LastName:  member.Name().Last,
			Seniority: string(member.Seniority()),
			Status:    statusFromMember(&member),
			CreatedAt: member.CreatedAt().Format("2006-01-02T15:04:05Z"),
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
func HandlerEditMember(w http.ResponseWriter, r *http.Request, coord *coordinator.MemberAPICoordinator, memberID int64) {
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

	// Call coordinator to edit member
	err := coord.EditMember(memberID, req.FirstName, req.LastName, domain.Seniority(req.Seniority))
	if err != nil {
		if validationErr, ok := err.(*coordinator.ValidationError); ok {
			w.Header().Set("Content-Type", "application/json")
			statusCode := http.StatusBadRequest
			if validationErr.Kind == coordinator.ValidationErrorKindDatabaseFailure {
				statusCode = http.StatusInternalServerError
			}
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: validationErr.UserMessage(),
				Kind:  string(validationErr.Kind),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Unexpected error",
			Kind:  "unknown",
		})
		return
	}

	// Get the updated member
	member := coord.GetMemberByID(memberID)
	if member == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Member was updated but not found",
			Kind:  "unknown",
		})
		return
	}

	response := MemberResponse{
		ID:        memberIDToUUID(int64(member.ID())),
		FirstName: member.Name().First,
		LastName:  member.Name().Last,
		Seniority: string(member.Seniority()),
		Status:    statusFromMember(member),
		CreatedAt: member.CreatedAt().Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandlerDeleteMember handles DELETE /api/members/{id}
func HandlerDeleteMember(w http.ResponseWriter, r *http.Request, coord *coordinator.MemberAPICoordinator, memberID int64) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call coordinator to deactivate member
	err := coord.DeactivateMember(memberID)
	if err != nil {
		if validationErr, ok := err.(*coordinator.ValidationError); ok {
			w.Header().Set("Content-Type", "application/json")
			statusCode := http.StatusBadRequest
			if validationErr.Kind == coordinator.ValidationErrorKindDatabaseFailure {
				statusCode = http.StatusInternalServerError
			}
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: validationErr.UserMessage(),
				Kind:  string(validationErr.Kind),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Unexpected error",
			Kind:  "unknown",
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

// statusFromMember converts a member's active state to a status string
func statusFromMember(member *domain.TeamMember) string {
	if member.IsActive() {
		return "active"
	}
	return "inactive"
}

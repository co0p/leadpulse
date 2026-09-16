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

// HandlerGetMembersPage handles GET /members (HTML page)
func HandlerGetMembersPage(w http.ResponseWriter, r *http.Request, getMembersUC GetMembersUC) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call use case to get all active members
	output, err := getMembersUC.Execute()
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><body><h1>Error</h1><p>Failed to load members: %s</p></body></html>", err.Error())
		return
	}

	// Simple HTML response with member list
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Build HTML response
	html := "<html><head><title>Members</title></head><body>"
	html += "<h1>Members</h1>"

	if len(output.Members) == 0 {
		html += "<p>No members found</p>"
	} else {
		html += "<table border='1'><tr><th>Name</th><th>Seniority</th></tr>"
		for _, m := range output.Members {
			html += fmt.Sprintf("<tr><td>%s %s</td><td>%s</td></tr>", m.FirstName, m.LastName, m.Seniority)
		}
		html += "</table>"
	}

	html += "</body></html>"

	fmt.Fprint(w, html)
}

// HandlerGetAddMemberPage handles GET /members/add (add member form page)
func HandlerGetAddMemberPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Simple HTML form
	html := `
	<!DOCTYPE html>
	<html>
	<head><title>Add Member</title></head>
	<body>
		<h1>Add Member</h1>
		<form method="POST" action="/members">
			<div>
				<label>First Name:</label>
				<input type="text" name="firstName" required>
			</div>
			<div>
				<label>Last Name:</label>
				<input type="text" name="lastName" required>
			</div>
			<div>
				<label>Seniority:</label>
				<select name="seniority" required>
					<option value="">-- Select --</option>
					<option value="junior">Junior</option>
					<option value="mid">Mid-Level</option>
					<option value="senior">Senior</option>
					<option value="lead">Lead</option>
				</select>
			</div>
			<button type="submit">Add Member</button>
			<a href="/members">Cancel</a>
		</form>
	</body>
	</html>
	`
	fmt.Fprint(w, html)
}

// HandlerPostAddMember handles POST /members (form submission for add)
func HandlerPostAddMember(w http.ResponseWriter, r *http.Request, addMemberUC AddMemberUC) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><body><h1>Error</h1><p>Invalid form data</p></body></html>")
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	seniority := r.FormValue("seniority")

	// Call use case
	output, err := addMemberUC.Execute(members.AddMemberInput{
		FirstName: firstName,
		LastName:  lastName,
		Seniority: seniority,
	})

	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>Add Member</title></head>
		<body>
			<h1>Add Member</h1>
			<p style="color: red;">Error: %s</p>
			<form method="POST" action="/members">
				<div>
					<label>First Name:</label>
					<input type="text" name="firstName" value="%s" required>
				</div>
				<div>
					<label>Last Name:</label>
					<input type="text" name="lastName" value="%s" required>
				</div>
				<div>
					<label>Seniority:</label>
					<select name="seniority" required>
						<option value="">-- Select --</option>
						<option value="junior" %s>Junior</option>
						<option value="mid" %s>Mid-Level</option>
						<option value="senior" %s>Senior</option>
						<option value="lead" %s>Lead</option>
					</select>
				</div>
				<button type="submit">Add Member</button>
				<a href="/members">Cancel</a>
			</form>
		</body>
		</html>
		`, err.Error(), firstName, lastName,
			map[bool]string{true: "selected"}[seniority == "junior"],
			map[bool]string{true: "selected"}[seniority == "mid"],
			map[bool]string{true: "selected"}[seniority == "senior"],
			map[bool]string{true: "selected"}[seniority == "lead"])
		fmt.Fprint(w, html)
		return
	}

	// Redirect to members list
	w.Header().Set("Location", "/members")
	w.WriteHeader(http.StatusSeeOther)

	_ = output // silence unused warning
}

// HandlerGetEditMemberPage handles GET /members/{id}/edit (edit member form page)
func HandlerGetEditMemberPage(w http.ResponseWriter, r *http.Request, getMembersUC GetMembersUC, memberID int64) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call use case to get all members
	output, err := getMembersUC.Execute()
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><body><h1>Error</h1><p>Failed to load members</p></body></html>")
		return
	}

	// Find member by ID
	var foundMember *members.MemberDTO
	for i := range output.Members {
		if output.Members[i].ID == memberID {
			foundMember = &output.Members[i]
			break
		}
	}

	if foundMember == nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<html><body><h1>Not Found</h1><p>Member not found</p><a href='/members'>Back to members</a></body></html>")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Simple HTML form with pre-filled values
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head><title>Edit Member</title></head>
	<body>
		<h1>Edit Member</h1>
		<form method="POST" action="/api/members/%d" onsubmit="handleEdit(event)">
			<div>
				<label>First Name:</label>
				<input type="text" name="firstName" value="%s" required>
			</div>
			<div>
				<label>Last Name:</label>
				<input type="text" name="lastName" value="%s" required>
			</div>
			<div>
				<label>Seniority:</label>
				<select name="seniority" required>
					<option value="">-- Select --</option>
					<option value="junior" %s>Junior</option>
					<option value="mid" %s>Mid-Level</option>
					<option value="senior" %s>Senior</option>
					<option value="lead" %s>Lead</option>
				</select>
			</div>
			<button type="submit">Save Member</button>
			<a href="/members">Cancel</a>
		</form>
		<script>
			async function handleEdit(event) {
				event.preventDefault();
				const firstName = document.querySelector('input[name="firstName"]').value;
				const lastName = document.querySelector('input[name="lastName"]').value;
				const seniority = document.querySelector('select[name="seniority"]').value;
				
				const response = await fetch('/api/members/%d', {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ firstName, lastName, seniority })
				});
				
				if (response.ok) {
					window.location.href = '/members';
				} else {
					alert('Error saving member');
				}
			}
		</script>
	</body>
	</html>
	`, memberID, foundMember.FirstName, foundMember.LastName,
		map[bool]string{true: "selected"}[foundMember.Seniority == "junior"],
		map[bool]string{true: "selected"}[foundMember.Seniority == "mid"],
		map[bool]string{true: "selected"}[foundMember.Seniority == "senior"],
		map[bool]string{true: "selected"}[foundMember.Seniority == "lead"],
		memberID)

	fmt.Fprint(w, html)
}

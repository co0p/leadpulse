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
	// Build the members table rows
	membersHTML := ""
	
	if len(memberList) == 0 {
		membersHTML = `<div style="text-align: center; padding: 2rem; color: #7a7a7a;">
			<p>No members found</p>
		</div>`
	} else {
		membersHTML = `<table class="table is-striped is-hoverable is-fullwidth">
			<thead>
				<tr>
					<th>Name</th>
					<th>Seniority</th>
					<th>Status</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>`

		for _, m := range memberList {
			statusBadgeClass := ""
			if m.Status == "Inactive" {
				statusBadgeClass = "badge-inactive"
			}
			
			rowClass := ""
			if m.Status == "Inactive" {
				rowClass = ` class="inactive-row"`
			}

			actionButtons := ""
			if m.Status == "Inactive" {
				actionButtons = fmt.Sprintf(
					`<button type="button" class="button is-small is-success" onclick="reactivateMember(%d)">
						<span class="icon"><i class="fas fa-redo"></i></span>
						<span>Reactivate</span>
					</button>`,
					m.ID,
				)
			} else {
				actionButtons = fmt.Sprintf(
					`<a href="/members/%d/edit" class="button is-small is-info">Edit</a>
					<button type="button" class="button is-small is-danger" onclick="deactivateMember(%d)">
						<span class="icon"><i class="fas fa-trash"></i></span>
						<span>Deactivate</span>
					</button>`,
					m.ID, m.ID,
				)
			}

			membersHTML += fmt.Sprintf(
				`<tr%s>
					<td>%s %s</td>
					<td>%s</td>
					<td><span class="badge %s">%s</span></td>
					<td>%s</td>
				</tr>`,
				rowClass,
				m.FirstName, m.LastName,
				m.Seniority,
				statusBadgeClass, m.Status,
				actionButtons,
			)
		}

		membersHTML += `</tbody></table>`
	}

	// Build the full page with shell
	mainContent := fmt.Sprintf(`<section class="section">
	<div class="container">
		<div class="level">
			<div class="level-left">
				<div class="level-item">
					<h1 class="title">Team Members</h1>
				</div>
			</div>
			<div class="level-right">
				<div class="level-item">
					<a href="/members/add" class="button is-primary">
						<span class="icon"><i class="fas fa-plus"></i></span>
						<span>Add Member</span>
					</a>
				</div>
			</div>
		</div>

		<!-- Bulma Tabs -->
		<div class="tabs">
			<ul>
				<li class="%s">
					<a href="/members?status=active">Active</a>
				</li>
				<li class="%s">
					<a href="/members?status=deactivated">Deactivated</a>
				</li>
				<li class="%s">
					<a href="/members?status=all">All</a>
				</li>
			</ul>
		</div>

		<div class="box">
			%s
		</div>
	</div>
</section>

<style>
	.inactive-row {
		color: #999;
	}
	.badge {
		display: inline-block;
		padding: 2px 6px;
		border-radius: 3px;
		font-size: 12px;
	}
	.badge-inactive {
		background-color: #f5f5f5;
		color: #666;
	}
</style>

<script>
	async function reactivateMember(memberId) {
		if (!confirm('Reactivate this member?')) return;
		
		const response = await fetch('/api/members/' + memberId + '/reactivate', {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' }
		});
		
		if (response.ok) {
			window.location.reload();
		} else {
			alert('Error reactivating member');
		}
	}

	async function deactivateMember(memberId) {
		if (!confirm('Deactivate this member?')) return;
		
		const response = await fetch('/api/members/' + memberId, {
			method: 'DELETE',
			headers: { 'Content-Type': 'application/json' }
		});
		
		if (response.ok) {
			window.location.reload();
		} else {
			alert('Error deactivating member');
		}
	}
</script>`,
		map[bool]string{true: "is-active"}[status == "active"],
		map[bool]string{true: "is-active"}[status == "deactivated"],
		map[bool]string{true: "is-active"}[status == "all"],
		membersHTML,
	)

	// Return full shell + content
	return buildShellHTML(mainContent)
}

// buildShellHTML constructs the full application shell with the given main content
func buildShellHTML(mainContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Team Impact Scorecard</title>
    
    <!-- Bulma CSS -->
    <link rel="stylesheet" href="/dist/css/bulma.min.css">
    
    <!-- Font Awesome for icons -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    
    <!-- Shell styles -->
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        html, body {
            height: 100vh;
            overflow: hidden;
        }

        body {
            display: flex;
            flex-direction: column;
        }

        .app-container {
            display: flex;
            flex: 1;
            overflow: hidden;
        }

        .topbar {
            height: 56px;
            background-color: #fff;
            border-bottom: 1px solid #e8e8e8;
            display: flex;
            align-items: center;
            padding: 0 1.25rem;
            gap: 1rem;
            flex-shrink: 0;
            z-index: 10;
        }

        .topbar-left {
            display: flex;
            align-items: center;
            gap: 1rem;
            min-width: 200px;
        }

        .topbar-center {
            flex: 1;
            display: flex;
            justify-content: center;
            min-width: 0;
        }

        .topbar-center .field {
            width: 100%%;
            max-width: 400px;
        }

        .topbar-right {
            display: flex;
            align-items: center;
            min-width: 120px;
        }

        .sidebar {
            width: 260px;
            background-color: #f5f5f5;
            border-right: 1px solid #e8e8e8;
            display: flex;
            flex-direction: column;
            overflow-y: auto;
            flex-shrink: 0;
            z-index: 5;
        }

        .sidebar-header {
            padding: 1.5rem 1.25rem;
            border-bottom: 1px solid #e8e8e8;
            flex-shrink: 0;
        }

        .sidebar-header .logo h2 {
            font-size: 1.25rem;
            font-weight: 600;
            color: #2c3e50;
            margin: 0;
        }

        .sidebar-nav {
            flex: 1;
            overflow-y: auto;
            padding: 1rem 0;
        }

        .menu-list {
            list-style: none;
            padding: 0 0.5rem;
        }

        .menu-list li {
            margin: 0;
        }

        .menu-list a {
            display: block;
            padding: 0.75rem 1rem;
            color: #4a4a4a;
            text-decoration: none;
            border-radius: 4px;
            transition: background-color 0.2s;
        }

        .menu-list a:hover {
            background-color: #ebebeb;
        }

        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow-y: auto;
            background-color: #fafafa;
        }

        .main-content .section {
            padding: 1rem 1.25rem;
        }

        .main-content .container {
            margin: 0;
            max-width: none;
            width: 100%%;
        }

        .breadcrumb-area {
            display: flex;
            align-items: center;
            font-size: 0.95rem;
            color: #7a7a7a;
        }

        .sidebar-toggle {
            padding: 0.5rem;
            min-width: 44px;
        }

        /* Responsive */
        @media screen and (max-width: 1023px) {
            .sidebar {
                position: absolute;
                left: 0;
                top: 56px;
                height: calc(100vh - 56px);
                transform: translateX(-100%%);
                transition: transform 0.3s ease;
                box-shadow: 2px 0 5px rgba(0,0,0,0.1);
            }

            .sidebar.is-active {
                transform: translateX(0);
            }
        }

        @media screen and (max-width: 767px) {
            .topbar {
                padding: 0 0.75rem;
                gap: 0.5rem;
            }

            .main-content .section {
                padding: 0.75rem;
            }
        }
    </style>
</head>
<body x-data="shellState()" @keydown.escape="closeSidebar()">
    <header class="topbar" role="banner">
        <div class="topbar-left">
            <button class="button is-white sidebar-toggle" type="button" aria-label="Toggle sidebar" :aria-expanded="sidebarOpen" @click="toggleSidebar()">
                <span class="icon">
                    <i class="fas fa-bars"></i>
                </span>
            </button>
            <div class="breadcrumb-area">
                <span class="breadcrumb-text">Dashboard</span>
            </div>
        </div>

        <div class="topbar-center">
            <div class="field has-addons is-fullwidth">
                <p class="control is-expanded has-icons-left">
                    <input class="input" type="text" placeholder="Search..." aria-label="Search members and content" @focus="searchFocused = true" @blur="searchFocused = false">
                    <span class="icon is-left">
                        <i class="fas fa-search"></i>
                    </span>
                </p>
            </div>
        </div>

        <div class="topbar-right">
            <div class="buttons">
                <button class="button is-ghost" type="button" aria-label="Notifications" title="Notifications">
                    <span class="icon">
                        <i class="fas fa-bell"></i>
                    </span>
                </button>
            </div>
        </div>
    </header>

    <div class="app-container">
        <aside class="sidebar" :class="{ 'is-active': sidebarOpen }" role="navigation">
            <div class="sidebar-header">
                <div class="logo">
                    <h2>Leadpulse</h2>
                </div>
            </div>
            
            <nav class="sidebar-nav">
                <ul class="menu-list">
                    <li><a href="/" tabindex="0">Home</a></li>
                    <li><a href="/members" tabindex="0">Members</a></li>
                    <li><a href="/alerts" tabindex="0">Alerts</a></li>
                    <li><a href="/reports" tabindex="0">Reports</a></li>
                </ul>
            </nav>
        </aside>

        <main class="main-content" role="main">
            %s
        </main>
    </div>

    <!-- Alpine.js -->
    <script src="/dist/js/alpine.min.js" defer></script>

    <!-- Shell state and interactivity -->
    <script>
        function shellState() {
            return {
                sidebarOpen: window.innerWidth >= 1024,
                searchFocused: false,

                toggleSidebar() {
                    if (window.innerWidth < 1024) {
                        this.sidebarOpen = !this.sidebarOpen;
                    }
                },

                closeSidebar() {
                    if (window.innerWidth < 1024) {
                        this.sidebarOpen = false;
                    }
                },

                init() {
                    // Update sidebar state on window resize
                    window.addEventListener('resize', () => {
                        if (window.innerWidth >= 1024) {
                            this.sidebarOpen = true;
                        } else {
                            this.sidebarOpen = false;
                        }
                    });
                }
            };
        }
    </script>
</body>
</html>`, mainContent)
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

// HandlerGetMembersPage handles GET /members (HTML page with status filtering and tabs)
// This handler composes the shell template with the members content template
func HandlerGetMembersPage(w http.ResponseWriter, r *http.Request, getMembersUC GetMembersUC, getFilteredMembersUC GetFilteredMembersUC) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get status query parameter (defaults to "active")
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	// Validate status parameter
	if status != "active" && status != "deactivated" && status != "all" {
		status = "active"
	}

	var output *members.GetMembersOutput
	var err error

	// Fetch members based on status
	if status == "active" && r.URL.Query().Get("status") == "" {
		output, err = getMembersUC.Execute()
	} else {
		input := members.GetFilteredMembersInput{Status: status}
		output, err = getFilteredMembersUC.Execute(input)
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><body><h1>Error</h1><p>Failed to load members: %s</p></body></html>", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Build data for template
	data := map[string]interface{}{
		"Members": output.Members,
		"Status":  status,
	}

	// Render main content first, then wrap in shell
	// For now, return HTML (templates will be properly loaded in server.go)
	html := buildMembersPageHTML(output.Members, status)
	fmt.Fprint(w, html)

	_ = data
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

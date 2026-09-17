package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
)

// Assets holds the embedded web directory.
//go:embed dist templates
var Assets embed.FS

// Start boots the HTTP server on the given address.
// It serves the shell template on /, static assets from the embedded FS under /dist/,
// and registers API endpoints for team member CRUD operations.
// The function blocks indefinitely while the server runs.
func Start(addr string, addMemberUC AddMemberUC, getMembersUC GetMembersUC, getFilteredMembersUC GetFilteredMembersUC, editMemberUC EditMemberUC, deactivateMemberUC DeactivateMemberUC, reactivateMemberUC ReactivateMemberUC) error {
	mux := http.NewServeMux()

	// Register Members API endpoints
	if addMemberUC != nil && getMembersUC != nil && getFilteredMembersUC != nil && editMemberUC != nil && deactivateMemberUC != nil && reactivateMemberUC != nil {
		// JSON API endpoints
		mux.HandleFunc("POST /api/members", func(w http.ResponseWriter, r *http.Request) {
			HandlerAddMember(w, r, addMemberUC)
		})

		mux.HandleFunc("GET /api/members", func(w http.ResponseWriter, r *http.Request) {
			HandlerGetMembers(w, r, getMembersUC, getFilteredMembersUC)
		})

		mux.HandleFunc("PATCH /api/members/{id}", func(w http.ResponseWriter, r *http.Request) {
			// Extract memberID from path parameter
			idStr := r.PathValue("id")
			memberID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, `{"error":"Invalid member ID","kind":"validation_error"}`, http.StatusBadRequest)
				return
			}
			HandlerEditMember(w, r, editMemberUC, memberID)
		})

		mux.HandleFunc("DELETE /api/members/{id}", func(w http.ResponseWriter, r *http.Request) {
			// Extract memberID from path parameter
			idStr := r.PathValue("id")
			memberID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, `{"error":"Invalid member ID","kind":"validation_error"}`, http.StatusBadRequest)
				return
			}
			HandlerDeleteMember(w, r, deactivateMemberUC, memberID)
		})

		mux.HandleFunc("PATCH /api/members/{id}/reactivate", func(w http.ResponseWriter, r *http.Request) {
			// Extract memberID from path parameter
			idStr := r.PathValue("id")
			memberID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, `{"error":"Invalid member ID","kind":"validation_error"}`, http.StatusBadRequest)
				return
			}
			HandlerReactivateMember(w, r, reactivateMemberUC, memberID)
		})

		// HTML page endpoints
		mux.HandleFunc("GET /members", func(w http.ResponseWriter, r *http.Request) {
			HandlerGetMembersPage(w, r, getMembersUC, getFilteredMembersUC)
		})

		mux.HandleFunc("GET /members/add", func(w http.ResponseWriter, r *http.Request) {
			HandlerGetAddMemberPage(w, r)
		})

		mux.HandleFunc("POST /members", func(w http.ResponseWriter, r *http.Request) {
			HandlerPostAddMember(w, r, addMemberUC)
		})

		mux.HandleFunc("GET /members/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
			// Extract memberID from path parameter
			idStr := r.PathValue("id")
			memberID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid member ID", http.StatusBadRequest)
				return
			}
			HandlerGetEditMemberPage(w, r, getMembersUC, memberID)
		})
	}

	// Parse shell template
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		return fmt.Errorf("failed to parse shell template: %w", err)
	}

	// Health check endpoint
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Root handler: serve shell template
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		}
	})

	// Extract the dist subdirectory from the embedded FS
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		return fmt.Errorf("failed to extract dist directory from embedded FS: %w", err)
	}

	// Serve static assets from /dist/
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/dist/", http.StripPrefix("/dist/", fileServer))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Starting HTTP server on http://%s\n", addr)
	return server.ListenAndServe()
}

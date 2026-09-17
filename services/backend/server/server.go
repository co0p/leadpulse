package server

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"strconv"
)

// Assets holds the embedded web directory (Vue SPA build output).
//go:embed dist
var Assets embed.FS

// Start boots the HTTP server on the given address.
// It serves the Vue SPA on /, static assets from the embedded Vue dist/,
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
	}

	// Health check endpoint
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Extract the dist subdirectory from the embedded FS (Vue SPA build output)
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		return fmt.Errorf("failed to extract dist directory from embedded FS: %w", err)
	}

	// Serve static assets from /dist/ (CSS, JS, fonts, etc.)
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/dist/", http.StripPrefix("/dist/", fileServer))

	// SPA bootstrap handler: serve index.html for all non-API, non-static routes
	// This allows Vue Router to handle client-side routing
	mux.HandleFunc("/", serveSPA(distFS))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Starting HTTP server on http://%s\n", addr)
	return server.ListenAndServe()
}

// serveSPA returns a handler that serves the Vue SPA index.html for all routes
// except those handled by other more specific handlers.
// This enables Vue Router's client-side routing to work correctly.
func serveSPA(distFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reject paths that should not reach SPA bootstrap
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the requested file first (for assets not caught by /dist/ handler)
		file, err := distFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err == nil {
			defer file.Close()
			// File exists; let it be served
			fileServer := http.FileServer(http.FS(distFS))
			fileServer.ServeHTTP(w, r)
			return
		}

		// File not found; serve index.html for SPA routing
		indexFile, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, indexFile); err != nil {
			fmt.Printf("Error serving index.html: %v\n", err)
		}
	}
}

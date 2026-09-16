package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
)

// Assets holds the embedded web directory.
//go:embed dist
var Assets embed.FS

// Start boots the HTTP server on the given address.
// It serves static assets from the embedded FS under dist/.
// It also registers API endpoints for team member CRUD operations.
// The function blocks indefinitely while the server runs.
func Start(addr string, addMemberUC AddMemberUC, getMembersUC GetMembersUC, editMemberUC EditMemberUC, deactivateMemberUC DeactivateMemberUC) error {
	mux := http.NewServeMux()

	// Register Members API endpoints
	if addMemberUC != nil && getMembersUC != nil && editMemberUC != nil && deactivateMemberUC != nil {
		mux.HandleFunc("POST /api/members", func(w http.ResponseWriter, r *http.Request) {
			HandlerAddMember(w, r, addMemberUC)
		})

		mux.HandleFunc("GET /api/members", func(w http.ResponseWriter, r *http.Request) {
			HandlerGetMembers(w, r, getMembersUC)
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
	}

	// Extract the dist subdirectory from the embedded FS
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		return fmt.Errorf("failed to extract dist directory from embedded FS: %w", err)
	}

	// Serve embedded assets with a root index handler
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Starting HTTP server on http://%s\n", addr)
	return server.ListenAndServe()
}

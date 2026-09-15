package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

// Assets holds the embedded web directory.
//go:embed dist
var Assets embed.FS

// Start boots the HTTP server on the given address.
// It serves static assets from the embedded FS under dist/.
// The function blocks indefinitely while the server runs.
func Start(addr string) error {
	mux := http.NewServeMux()

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

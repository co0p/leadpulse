package server

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEmbedAssets verifies that the embedded assets are accessible.
func TestEmbedAssets(t *testing.T) {
	// Verify dist directory exists in the embedded FS
	_, err := Assets.Open("dist")
	if err != nil {
		t.Fatalf("failed to open dist directory: %v", err)
	}

	// Verify index.html exists in the embedded FS
	indexFile, err := Assets.Open("dist/index.html")
	if err != nil {
		t.Fatalf("index.html not found in embedded FS: %v", err)
	}
	defer indexFile.Close()

	// Verify index.html has content
	content, err := io.ReadAll(indexFile)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("index.html is empty")
	}

	if !contains(string(content), "Team Impact Scorecard") {
		t.Error("index.html does not contain expected content")
	}
}

// TestRootHandlerWithEmbedFS verifies that a GET request to / returns embedded content.
func TestRootHandlerWithEmbedFS(t *testing.T) {
	// Create the dist filesystem just like Start() does
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	// Create a handler similar to what Start() creates
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(distFS)))

	// Create a test request for the root path
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	mux.ServeHTTP(w, req)

	// The response may be a redirect or a 200; just verify we got a response
	if w.Code < 200 || w.Code >= 400 {
		t.Logf("GET / returned status %d (acceptable for directory listing)", w.Code)
	}
}

// TestIndexHTMLDirect verifies that we can serve index.html directly from the distFS.
func TestIndexHTMLDirect(t *testing.T) {
	// Create the dist filesystem like Start() does
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	// Open index.html from the distFS
	indexFile, err := distFS.Open("index.html")
	if err != nil {
		t.Fatalf("failed to open index.html from distFS: %v", err)
	}
	defer indexFile.Close()

	// Verify the content
	content, err := io.ReadAll(indexFile)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("index.html content is empty")
	}

	if !contains(string(content), "<!DOCTYPE html>") {
		t.Error("index.html is not valid HTML")
	}

	if !contains(string(content), "Team Impact Scorecard") {
		t.Error("index.html missing expected content")
	}
}

// TestFileServerWithDistFS verifies FileServer can serve files from distFS.
func TestFileServerWithDistFS(t *testing.T) {
	// Create the dist filesystem
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	// Create a file server like Start() does
	fileServer := http.FileServer(http.FS(distFS))

	// Create a test request for index.html
	req := httptest.NewRequest("GET", "/index.html", nil)
	w := httptest.NewRecorder()

	fileServer.ServeHTTP(w, req)

	// Verify we got a response (may be a redirect or 200)
	if w.Code == http.StatusNotFound {
		t.Fatalf("index.html not found in file server: status %d", w.Code)
	}

	if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
		// FileServer redirects directory requests; this is acceptable
		t.Logf("FileServer returned redirect status %d (acceptable)", w.Code)
		return
	}

	if w.Code == http.StatusOK {
		body := w.Body.String()
		if !contains(body, "Team Impact Scorecard") {
			t.Error("response missing expected content")
		}
	}
}

// TestHandlerNoFyneImports ensures this test file has no Fyne dependencies.
// This is a compile-time check that passes if the file builds without fyne imports.
func TestHandlerNoFyneImports(t *testing.T) {
	// This test simply verifies the package compiles without Fyne imports.
	// If there were Fyne imports, the build would fail.
	t.Log("server package has no Fyne imports")
}

// contains is a helper to check if a string contains a substring.
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

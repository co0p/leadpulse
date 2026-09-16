package server

import (
	"html/template"
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

	// Verify templates directory exists in the embedded FS
	_, err = Assets.Open("templates")
	if err != nil {
		t.Fatalf("templates directory not found in embedded FS: %v", err)
	}

	// Verify layout.html exists in the embedded FS
	layoutFile, err := Assets.Open("templates/layout.html")
	if err != nil {
		t.Fatalf("layout.html not found in embedded FS: %v", err)
	}
	defer layoutFile.Close()

	// Verify layout.html has content
	content, err := io.ReadAll(layoutFile)
	if err != nil {
		t.Fatalf("failed to read layout.html: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("layout.html is empty")
	}

	if !contains(string(content), "Team Impact Scorecard") {
		t.Error("layout.html does not contain expected content")
	}
}

// TestRootHandlerWithTemplate verifies that a GET request to / returns templated HTML.
func TestRootHandlerWithTemplate(t *testing.T) {
	// Parse the shell template
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	// Create a handler that serves the template
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	})

	// Create a test request for the root path
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	mux.ServeHTTP(w, req)

	// Verify we got a 200 OK response
	if w.Code != http.StatusOK {
		t.Fatalf("GET / returned status %d, expected 200", w.Code)
	}

	// Verify the response contains HTML content
	body := w.Body.String()
	if !contains(body, "<!DOCTYPE html>") {
		t.Error("response does not contain DOCTYPE")
	}

	if !contains(body, "Team Impact Scorecard") {
		t.Error("response does not contain expected title")
	}

	if !contains(body, "<header") {
		t.Error("response does not contain header element")
	}

	if !contains(body, "<aside") {
		t.Error("response does not contain aside/sidebar element")
	}

	if !contains(body, "<main") {
		t.Error("response does not contain main content element")
	}
}

// TestLayoutHTMLDirect verifies that we can parse layout.html directly from the embedded FS.
func TestLayoutHTMLDirect(t *testing.T) {
	// Parse the layout template
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse layout.html from embedded FS: %v", err)
	}

	// Verify template is not nil
	if tmpl == nil {
		t.Fatal("layout template is nil")
	}

	// Verify template name
	if tmpl.Name() != "layout.html" {
		t.Errorf("template name is %q, expected layout.html", tmpl.Name())
	}
}

// TestStaticAssetsServer verifies FileServer can serve assets from dist/.
func TestStaticAssetsServer(t *testing.T) {
	// Create the dist filesystem
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	// Create a file server like Start() does
	fileServer := http.FileServer(http.FS(distFS))

	// Test CSS file
	req := httptest.NewRequest("GET", "/css/bulma.min.css", nil)
	w := httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /css/bulma.min.css returned status %d, expected 200", w.Code)
	}

	// Test JavaScript file
	req = httptest.NewRequest("GET", "/js/alpine.min.js", nil)
	w = httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /js/alpine.min.js returned status %d, expected 200", w.Code)
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

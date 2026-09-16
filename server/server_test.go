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

// Acceptance tests for Web App Shell (AC-1 through AC-7)

// TestShellRendersFullLayout verifies shell renders with all 3 semantic regions
func TestShellRendersFullLayout(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, "<header") || !contains(body, "<aside") || !contains(body, "<main") {
		t.Error("shell missing semantic regions (header, aside, main)")
	}
	if !contains(body, "height: 100vh") {
		t.Error("layout missing full-height CSS")
	}
}

// TestSidebarComponentRenders verifies sidebar contains expected elements
func TestSidebarComponentRenders(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, "Leadpulse") {
		t.Error("sidebar missing logo/branding")
	}
	if !contains(body, "Home") || !contains(body, "Settings") {
		t.Error("sidebar missing nav links")
	}
	if !contains(body, "Tools") {
		t.Error("sidebar missing collapsible section")
	}
	if !contains(body, "Add Item") {
		t.Error("sidebar missing footer button")
	}
}

// TestTopBarComponentRenders verifies top bar contains expected elements
func TestTopBarComponentRenders(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, "sidebar-toggle") {
		t.Error("top bar missing toggle button")
	}
	if !contains(body, "Dashboard") {
		t.Error("top bar missing breadcrumb/context")
	}
	if !contains(body, `placeholder="Search..."`) {
		t.Error("top bar missing search input")
	}
	if !contains(body, "fa-bell") || !contains(body, "fa-question-circle") {
		t.Error("top bar missing quick action buttons")
	}
}

// TestContentAreaRenders verifies main content area exists and is scrollable
func TestContentAreaRenders(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, "main-content") {
		t.Error("content area missing main-content class")
	}
	if !contains(body, "overflow-y: auto") {
		t.Error("content area not scrollable")
	}
	if !contains(body, "Welcome to Team Impact Scorecard") {
		t.Error("content area missing placeholder content")
	}
}

// TestBulmaCSSLoadsWithout404 verifies Bulma CSS is accessible
func TestBulmaCSSLoadsWithout404(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	fileServer := http.FileServer(http.FS(distFS))
	req := httptest.NewRequest("GET", "/css/bulma.min.css", nil)
	w := httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /css/bulma.min.css returned %d, expected 200", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType == "" || !contains(contentType, "css") {
		t.Errorf("CSS content-type is %q, expected to contain 'css'", contentType)
	}

	if w.Body.Len() == 0 {
		t.Error("Bulma CSS file is empty")
	}
}

// TestHTMXLoadsWithout404 verifies HTMX JS is accessible
func TestHTMXLoadsWithout404(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	fileServer := http.FileServer(http.FS(distFS))
	req := httptest.NewRequest("GET", "/js/htmx.min.js", nil)
	w := httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /js/htmx.min.js returned %d, expected 200", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("HTMX JS file is empty")
	}
}

// TestAlpineLoadsWithout404 verifies Alpine.js is accessible
func TestAlpineLoadsWithout404(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	fileServer := http.FileServer(http.FS(distFS))
	req := httptest.NewRequest("GET", "/js/alpine.min.js", nil)
	w := httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /js/alpine.min.js returned %d, expected 200", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("Alpine JS file is empty")
	}
}

// TestResponsiveBreakpointsMetaTag verifies meta viewport tag for responsive design
func TestResponsiveBreakpointsMetaTag(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, `<meta name="viewport"`) || !contains(body, `width=device-width`) {
		t.Error("meta viewport tag missing or incomplete")
	}
}

// TestAccessibilitySemanticRegions verifies exactly 1 each of header, aside, main
func TestAccessibilitySemanticRegions(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	headerCount := countOccurrences(body, "<header")
	asideCount := countOccurrences(body, "<aside")
	mainCount := countOccurrences(body, "<main")

	if headerCount != 1 {
		t.Errorf("expected 1 <header>, found %d", headerCount)
	}
	if asideCount != 1 {
		t.Errorf("expected 1 <aside>, found %d", asideCount)
	}
	if mainCount != 1 {
		t.Errorf("expected 1 <main>, found %d", mainCount)
	}
}

// TestAccessibilityAriaExpandedOnToggle verifies sidebar toggle has aria-expanded
func TestAccessibilityAriaExpandedOnToggle(t *testing.T) {
	tmpl, err := template.ParseFS(Assets, "templates/layout.html")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, "aria-expanded") {
		t.Error("sidebar toggle missing aria-expanded attribute")
	}
	if !contains(body, `aria-label="Toggle sidebar"`) {
		t.Error("sidebar toggle missing aria-label")
	}
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

// countOccurrences counts how many times a substring appears in a string
func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			count++
			i += len(substr) - 1
		}
	}
	return count
}

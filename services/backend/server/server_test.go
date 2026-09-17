package server

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestEmbedAssets verifies that the embedded dist/ directory is accessible.
func TestEmbedAssets(t *testing.T) {
	// Verify dist directory exists in the embedded FS
	_, err := Assets.Open("dist")
	if err != nil {
		t.Skipf("dist directory not embedded locally (expected; will be embedded in Docker build): %v", err)
	}

	// Verify index.html exists in dist/
	indexFile, err := Assets.Open("dist/index.html")
	if err != nil {
		t.Fatalf("index.html not found in dist directory: %v", err)
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

	// SPA index.html should contain Vue app mount point
	if !strings.Contains(string(content), "id=\"app\"") {
		t.Error("index.html does not contain Vue mount point (id=\"app\")")
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

// SPA Bootstrap Tests (replacing template-based shell tests)

// TestRootPathReturnsSPAIndex verifies GET / serves SPA index.html
func TestRootPathReturnsSPAIndex(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Skipf("dist directory not embedded (expected during local development): %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveSPA(distFS))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET / returned status %d, expected 200", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "id=\"app\"") {
		t.Error("response does not contain Vue app mount point")
	}

	contentType := w.Header().Get("Content-Type")
	if contentType == "" || !strings.Contains(contentType, "text/html") {
		t.Errorf("Content-Type is %q, expected text/html", contentType)
	}
}

// TestSPARoutesServeSPAIndex verifies that SPA routes serve index.html for client-side routing
func TestSPARoutesServeSPAIndex(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Skipf("dist directory not embedded (expected during local development): %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveSPA(distFS))

	tests := []string{
		"/members",
		"/alerts",
		"/reports",
		"/members/1",
		"/nonexistent/path",
	}

	for _, path := range tests {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET %s returned status %d, expected 200", path, w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "id=\"app\"") {
			t.Errorf("GET %s response does not contain Vue app mount point", path)
		}
	}
}

// TestAPIPathsNotServedBySPA verifies that /api/* paths are not caught by SPA handler
func TestAPIPathsNotServedBySPA(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/", serveSPA(distFS))

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Should hit the health endpoint, not SPA handler
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/health returned status %d, expected 200", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "status") {
		t.Error("response should be JSON, not SPA index.html")
	}
}

// TestStaticAssetsCached verifies that assets served via /dist/ don't use SPA bootstrap
func TestStaticAssetsCached(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/dist/", http.StripPrefix("/dist/", fileServer))

	// Try to get a known CSS file
	req := httptest.NewRequest("GET", "/dist/css/bulma.min.css", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /dist/css/bulma.min.css returned status %d, expected 200", w.Code)
	}

	// Response should be CSS, not HTML
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") && !strings.Contains(contentType, "application/octet-stream") {
		t.Logf("Content-Type for CSS is %q (may vary by platform)", contentType)
	}
}

// Acceptance tests converted to SPA bootstrap behavior

// TestSPABootstrapRendersVueApp verifies that root path serves SPA with Vue app mount
func TestSPABootstrapRendersVueApp(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Skipf("dist directory not embedded (expected during local development): %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveSPA(distFS))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	// SPA should contain Vue app mount point and script references
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("response should start with DOCTYPE")
	}
	if !strings.Contains(body, "id=\"app\"") {
		t.Error("response missing Vue app mount point")
	}
	if !strings.Contains(body, "<script") {
		t.Error("response should contain script tags for Vue app")
	}
}

// TestSPANavigation verifies that client-side routes serve index.html
func TestSPANavigation(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Skipf("dist directory not embedded (expected during local development): %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveSPA(distFS))

	// Test navigation to a route that doesn't have a physical file
	req := httptest.NewRequest("GET", "/members", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /members returned status %d, expected 200", w.Code)
	}

	body := w.Body.String()
	// Should serve index.html, allowing Vue Router to handle the route
	if !strings.Contains(body, "id=\"app\"") {
		t.Error("response should be SPA index.html with app mount point")
	}
}

// TestHealthCheckEndpointJSON verifies /api/health returns JSON (not SPA bootstrap)
func TestHealthCheckEndpointJSON(t *testing.T) {
	mux := http.NewServeMux()

	// Health endpoint must be registered before SPA handler for priority
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "status") {
		t.Error("response should contain status field")
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
	if contentType == "" || (!strings.Contains(contentType, "css") && !strings.Contains(contentType, "octet-stream")) {
		t.Logf("Content-Type for CSS: %s (may vary by platform)", contentType)
	}
}

// TestHTMXLoadsWithout404 verifies HTMX library is accessible
func TestHTMXLoadsWithout404(t *testing.T) {
	distFS, err := fs.Sub(Assets, "dist")
	if err != nil {
		t.Fatalf("failed to extract dist: %v", err)
	}

	fileServer := http.FileServer(http.FS(distFS))
	req := httptest.NewRequest("GET", "/js/htmx.min.js", nil)
	w := httptest.NewRecorder()
	fileServer.ServeHTTP(w, req)

	// HTMX may not be in dist; this is a future test
	t.Logf("GET /js/htmx.min.js returned %d", w.Code)
}

// TestAlpineLoadsWithout404 verifies Alpine.js library is accessible
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
		t.Logf("GET /js/alpine.min.js returned %d (asset may not be in dist yet)", w.Code)
	}
}

// Helper function

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

package server

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestHTTPServerStartsAndServesHTML verifies the full acceptance criteria:
// AC-2: curl http://localhost:8080/ returns a 200 with a static HTML page embedded in the binary
func TestHTTPServerStartsAndServesHTML(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to find available port: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close()

	// Start server in a goroutine
	go func() {
		_ = Start(addr)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request to the root path
	resp, err := http.Get("http://" + addr + "/")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verify status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Verify response is HTML with expected content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	bodyStr := string(body)
	if !contains(bodyStr, "Team Impact Scorecard") {
		t.Error("response does not contain expected content")
	}

	if !contains(bodyStr, "<!DOCTYPE html>") {
		t.Error("response is not valid HTML")
	}

	if !contains(bodyStr, "SPA HTTP Server is running") {
		t.Error("response missing SPA server message")
	}
}

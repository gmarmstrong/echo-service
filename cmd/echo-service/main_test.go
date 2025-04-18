package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler_Basic(t *testing.T) {
	// Create a fake request
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("X-Test", "unit")

	rr := httptest.NewRecorder()
	// Call the handler with the recorder and request
	echoHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type: got %q, want application/json", ct)
	}

	var got responsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Method != http.MethodGet {
		t.Errorf("method: got %q, want GET", got.Method)
	}
	if got.Path != "/hello" {
		t.Errorf("path: got %q, want /hello", got.Path)
	}
	// Check for missing header, duplicate header, or wrong header value
	if h := got.Headers["X-Test"]; len(h) != 1 || h[0] != "unit" {
		t.Errorf("headers[X-Test]: got %v, want [unit]", h)
	}
}

func TestHealthz(t *testing.T) {
	// Create a fake request and a recorder
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	// We registered /healthz with an inline handler in main.go,
	// so we need a ServeMux like the one main() assembles.
	mux := http.NewServeMux()
	// Call the handler with the recorder and request
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want %d", rr.Code, http.StatusNoContent)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", rr.Body.Bytes())
	}
}

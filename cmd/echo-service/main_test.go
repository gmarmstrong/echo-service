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

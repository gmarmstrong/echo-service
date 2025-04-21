package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEchoHandler_Basic(t *testing.T) {
	// Create a synthetic request
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

func TestEchoHandler_MultiValueHeaderAndPOST(t *testing.T) {
	// Build a synthetic POST request with a repeated X-Multi header
	req := httptest.NewRequest(http.MethodPost, "/anything", nil)
	req.Header.Add("X-Multi", "foo")
	req.Header.Add("X-Multi", "bar") // second value for the same header

	rr := httptest.NewRecorder()
	// Call the handler with the recorder and request
	echoHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rr.Code, http.StatusOK)
	}

	var got responsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// Method should be POST
	if got.Method != http.MethodPost {
		t.Errorf("method: got %q, want POST", got.Method)
	}

	// Header slice should preserve both values in order
	want := []string{"foo", "bar"}
	if !reflect.DeepEqual(got.Headers["X-Multi"], want) {
		t.Errorf("headers[X-Multi]: got %v, want %v", got.Headers["X-Multi"], want)
	}
}

func TestHealthzHandler(t *testing.T) {
	// Create a synthetic request and a recorder
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	// Call the healthz handler with the recorder and request
	healthzHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want %d", rr.Code, http.StatusNoContent)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", rr.Body.Bytes())
	}
}

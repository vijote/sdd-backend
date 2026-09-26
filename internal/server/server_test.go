package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vijote/sdd-backend/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}

func TestUnknownRouteReturns404(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestShortenRouteWired(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	// Invalid URL proves the route is wired and reaches the service layer
	// (a real DB handle is required for a successful 201, which needs MySQL).
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"ftp://example.com"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid url") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

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
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"ftp://example.com"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid url") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestOldApiShortenRouteRemoved(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	// The dev ingress strips one /api prefix, so the backend must NOT register
	// /api/shorten — the old prefixed route must 404.
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestRedirectRouteWired(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	// A malformed code proves the route is wired: it 404s without touching the
	// (nil) database, since ValidateCode rejects it before the repo lookup.
	req := httptest.NewRequest(http.MethodGet, "/short", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if !strings.Contains(rec.Body.String(), "not found") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestRedirectUnknownCodeRouteWired(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	// A well-formed but unknown code reaches the repo; with a nil DB the repo
	// returns a database error, which the handler maps to 404.
	req := httptest.NewRequest(http.MethodGet, "/zzzzzzz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestRedirectPostMethodNotAllowed(t *testing.T) {
	router := NewRouter(&config.Config{Port: 8080, LogLevel: "info"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/aB3xK9m", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

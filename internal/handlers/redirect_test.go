package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/vijote/sdd-backend/internal/shortener"
)

// fakeRedirecter is an in-memory Redirecter for handler tests.
type fakeRedirecter struct {
	links    map[string]*shortener.Link // by code
	calls    int                        // number of Resolve invocations
	lastCode string
	fail     bool // return a non-ErrNotFound error
}

func (f *fakeRedirecter) Resolve(_ context.Context, code string) (*shortener.Link, error) {
	f.calls++
	f.lastCode = code
	if f.fail {
		return nil, errors.New("boom")
	}
	if link, ok := f.links[code]; ok {
		return link, nil
	}
	return nil, shortener.ErrNotFound
}

func getRedirect(t *testing.T, svc Redirecter, path string) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	r.Get("/{code}", Redirect(svc))
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestRedirectMovedPermanently(t *testing.T) {
	svc := &fakeRedirecter{links: map[string]*shortener.Link{
		"aB3xK9m": {Code: "aB3xK9m", LongURL: "https://example.com/path"},
	}}
	rec := getRedirect(t, svc, "/aB3xK9m")
	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/path" {
		t.Fatalf("Location = %q, want %q", loc, "https://example.com/path")
	}
	if svc.calls != 1 {
		t.Fatalf("expected 1 resolve call, got %d", svc.calls)
	}
}

func TestRedirectUnknownCode(t *testing.T) {
	svc := &fakeRedirecter{}
	rec := getRedirect(t, svc, "/zzzzzzz")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["error"] != "not found" {
		t.Errorf("error = %q", body["error"])
	}
}

func TestRedirectMalformedCode(t *testing.T) {
	svc := &fakeRedirecter{}
	rec := getRedirect(t, svc, "/short")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	// The handler delegates to the service; the service's ValidateCode rejects
	// malformed codes before any repository call (asserted in the shortener
	// package by TestServiceResolveMalformedSkipsRepo).
	if svc.calls != 1 {
		t.Fatalf("expected 1 resolve call, got %d", svc.calls)
	}
}

func TestRedirectServiceErrorIs404(t *testing.T) {
	// Any service error (not just ErrNotFound) surfaces as 404 per contract.
	svc := &fakeRedirecter{fail: true}
	rec := getRedirect(t, svc, "/aB3xK9m")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestRedirectHeadRequest(t *testing.T) {
	svc := &fakeRedirecter{links: map[string]*shortener.Link{
		"aB3xK9m": {Code: "aB3xK9m", LongURL: "https://example.com/path"},
	}}
	r := chi.NewRouter()
	r.Head("/{code}", Redirect(svc))
	req := httptest.NewRequest(http.MethodHead, "/aB3xK9m", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301 for HEAD, got %d", rec.Code)
	}
}

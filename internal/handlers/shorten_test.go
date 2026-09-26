package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vijote/sdd-backend/internal/shortener"
)

// fakeShortener is an in-memory Shortener for handler tests.
type fakeShortener struct {
	links map[string]*shortener.Link // by long URL
}

func newFakeShortener() *fakeShortener {
	return &fakeShortener{links: map[string]*shortener.Link{}}
}

func (f *fakeShortener) Shorten(_ context.Context, longURL string) (*shortener.Link, error) {
	if err := shortener.ValidateURL(longURL); err != nil {
		return nil, err
	}
	if link, ok := f.links[longURL]; ok {
		return link, nil
	}
	link := &shortener.Link{ID: uint(len(f.links) + 1), Code: "aB3xK9m", LongURL: longURL}
	f.links[longURL] = link
	return link, nil
}

func postShorten(t *testing.T, svc Shortener, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	rec := httptest.NewRecorder()
	Shorten(svc).ServeHTTP(rec, req)
	return rec
}

func TestShortenCreated(t *testing.T) {
	rec := postShorten(t, newFakeShortener(), `{"url":"https://example.com/path"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["code"] != "aB3xK9m" {
		t.Errorf("code = %q", body["code"])
	}
	if body["long_url"] != "https://example.com/path" {
		t.Errorf("long_url = %q", body["long_url"])
	}
}

func TestShortenInvalidURL(t *testing.T) {
	rec := postShorten(t, newFakeShortener(), `{"url":"ftp://example.com"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid url") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestShortenInvalidJSON(t *testing.T) {
	rec := postShorten(t, newFakeShortener(), `{not json`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid json body") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestShortenDuplicateSameCode(t *testing.T) {
	svc := newFakeShortener()
	first := postShorten(t, svc, `{"url":"https://example.com/path"}`)
	second := postShorten(t, svc, `{"url":"https://example.com/path"}`)
	if first.Code != http.StatusCreated || second.Code != http.StatusCreated {
		t.Fatalf("expected 201s, got %d and %d", first.Code, second.Code)
	}
	if first.Body.String() != second.Body.String() {
		t.Fatalf("expected identical bodies for duplicate, got %q vs %q", first.Body.String(), second.Body.String())
	}
}

func TestShortenInternalError(t *testing.T) {
	svc := failingShortener{}
	rec := postShorten(t, svc, `{"url":"https://example.com/path"}`)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "internal error") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

// failingShortener always errors with something non-validation.
type failingShortener struct{}

func (failingShortener) Shorten(_ context.Context, _ string) (*shortener.Link, error) {
	return nil, errors.New("boom")
}

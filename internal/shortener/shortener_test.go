package shortener

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
)

var base62Pattern = regexp.MustCompile(`^[0-9a-zA-Z]+$`)

func TestGenerateCodeDeterministic(t *testing.T) {
	first := GenerateCode("https://example.com/path", 0)
	second := GenerateCode("https://example.com/path", 0)
	if first != second {
		t.Fatalf("expected deterministic code, got %q vs %q", first, second)
	}
}

func TestGenerateCodeShape(t *testing.T) {
	code := GenerateCode("https://example.com/path", 0)
	if len(code) != CodeLength {
		t.Fatalf("expected %d chars, got %d (%q)", CodeLength, len(code), code)
	}
	if !base62Pattern.MatchString(code) {
		t.Fatalf("expected base62 code, got %q", code)
	}
}

func TestGenerateCodeSaltDistinct(t *testing.T) {
	first := GenerateCode("https://example.com/path", 0)
	second := GenerateCode("https://example.com/path", 1)
	if first == second {
		t.Fatalf("expected distinct codes for different salts, both %q", first)
	}
}

func TestValidateURL(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"valid https", "https://example.com/a/b", false},
		{"valid http", "http://example.com", false},
		{"empty", "", true},
		{"whitespace", "   ", true},
		{"no scheme", "example.com", true},
		{"ftp scheme", "ftp://example.com", true},
		{"javascript scheme", "javascript:alert(1)", true},
		{"no host", "http://", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.raw)
			if tc.wantErr && !errors.Is(err, ErrInvalidURL) {
				t.Fatalf("expected ErrInvalidURL, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected valid, got %v", err)
			}
		})
	}
}

// fakeRepo is an in-memory Repository for service tests.
type fakeRepo struct {
	links     map[string]*Link // by code
	byLongURL map[string]*Link
	createErr error // returned on the first create only
	creates   int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{links: map[string]*Link{}, byLongURL: map[string]*Link{}}
}

func (f *fakeRepo) Create(_ context.Context, link *Link) error {
	f.creates++
	if f.createErr != nil {
		err := f.createErr
		f.createErr = nil // fail once, then succeed
		return err
	}
	if _, dup := f.links[link.Code]; dup {
		return ErrDuplicateCode
	}
	stored := *link
	f.links[link.Code] = &stored
	f.byLongURL[link.LongURL] = &stored
	return nil
}

func (f *fakeRepo) FindByLongURL(_ context.Context, longURL string) (*Link, error) {
	if link, ok := f.byLongURL[longURL]; ok {
		return link, nil
	}
	return nil, ErrNotFound
}

func (f *fakeRepo) FindByCode(_ context.Context, code string) (*Link, error) {
	if link, ok := f.links[code]; ok {
		return link, nil
	}
	return nil, ErrNotFound
}

func TestServiceShortenCreatesLink(t *testing.T) {
	repo := newFakeRepo()
	link, err := NewService(repo).Shorten(context.Background(), "https://example.com/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(link.Code) != CodeLength {
		t.Fatalf("expected %d char code, got %q", CodeLength, link.Code)
	}
	if link.LongURL != "https://example.com/path" {
		t.Fatalf("unexpected long URL %q", link.LongURL)
	}
}

func TestServiceShortenIdempotent(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo)
	first, err := svc.Shorten(context.Background(), "https://example.com/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := svc.Shorten(context.Background(), "https://example.com/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.Code != second.Code {
		t.Fatalf("expected same code for duplicate, got %q vs %q", first.Code, second.Code)
	}
	if repo.creates != 1 {
		t.Fatalf("expected 1 create, got %d", repo.creates)
	}
}

func TestServiceShortenCollisionRetry(t *testing.T) {
	repo := newFakeRepo()
	repo.createErr = ErrDuplicateCode // fail once, then succeed
	svc := NewService(repo)
	link, err := svc.Shorten(context.Background(), "https://example.com/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.creates < 2 {
		t.Fatalf("expected retry after collision, got %d creates", repo.creates)
	}
	if len(link.Code) != CodeLength {
		t.Fatalf("expected %d char code, got %q", CodeLength, link.Code)
	}
}

func TestServiceShortenInvalidURL(t *testing.T) {
	repo := newFakeRepo()
	_, err := NewService(repo).Shorten(context.Background(), "ftp://example.com")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
	if !strings.Contains(err.Error(), "invalid url") {
		t.Fatalf("unexpected error text %q", err.Error())
	}
}

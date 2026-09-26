package shortener

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/vijote/sdd-backend/internal/database"
)

// CodeLength is the exact number of base62 characters in a short code.
const CodeLength = 7

// MaxSalt is the highest salt value tried during collision retry.
const MaxSalt = 1000

// ErrInvalidURL is returned when the long URL fails validation.
var ErrInvalidURL = errors.New("invalid url")

// ErrNotFound is returned by Repository lookups when no link exists.
var ErrNotFound = errors.New("link not found")

// ErrCodeExhausted is returned when collision retry exceeds MaxSalt.
var ErrCodeExhausted = errors.New("code generation exhausted")

// Link is the persisted short link (model owned by internal/database to keep
// migrate.go registration cycle-free).
type Link = database.Link

// Repository persists and queries links.
type Repository interface {
	Create(ctx context.Context, link *Link) error
	FindByLongURL(ctx context.Context, longURL string) (*Link, error)
	FindByCode(ctx context.Context, code string) (*Link, error)
}

// Service shortens long URLs.
type Service interface {
	Shorten(ctx context.Context, longURL string) (*Link, error)
}

// ValidateURL checks that raw is non-empty, parseable, and http/https.
func ValidateURL(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return ErrInvalidURL
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ErrInvalidURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidURL
	}
	if parsed.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

// GenerateCode hashes the long URL (with an incrementing salt) and returns
// the first CodeLength base62 characters. Deterministic for a given (url, salt).
func GenerateCode(longURL string, salt int) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", longURL, salt)))
	return base62Encode(hash[:])[:CodeLength]
}

// base62Encode encodes bytes as a base62 string ([0-9a-zA-Z]).
func base62Encode(data []byte) string {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	num := new(big.Int).SetBytes(data)
	base := big.NewInt(62)
	mod := new(big.Int)
	var out strings.Builder
	for num.Sign() > 0 {
		num.DivMod(num, base, mod)
		out.WriteByte(alphabet[mod.Int64()])
	}
	encoded := out.String()
	// Pad with leading zeros if the number was small.
	for len(encoded) < CodeLength {
		encoded = "0" + encoded
	}
	return encoded
}

// ErrDuplicateCode is returned by Repository.Create when the code already exists.
var ErrDuplicateCode = errors.New("duplicate code")

// service is the default Service implementation.
type service struct {
	repo Repository
}

// NewService creates a Service backed by the given repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Shorten validates the URL, returns an existing link for a duplicate long
// URL, or generates a code (with collision retry) and persists a new link.
func (s *service) Shorten(ctx context.Context, longURL string) (*Link, error) {
	if err := ValidateURL(longURL); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByLongURL(ctx, longURL)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	link := &Link{LongURL: longURL}
	for salt := 0; salt <= MaxSalt; salt++ {
		link.Code = GenerateCode(longURL, salt)
		if err := s.repo.Create(ctx, link); err == nil {
			return link, nil
		} else if !errors.Is(err, ErrDuplicateCode) {
			return nil, err
		}
	}
	return nil, ErrCodeExhausted
}

package shortener

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/vijote/sdd-backend/internal/database"
)

// gormRepository implements Repository over a GORM MySQL handle.
type gormRepository struct {
	db *database.DB
}

// NewRepository creates a Repository backed by the given database handle.
func NewRepository(db *database.DB) Repository {
	return &gormRepository{db: db}
}

// Create persists a new link. Uniqueness violations on code or long_url map
// to ErrDuplicateCode.
func (r *gormRepository) Create(ctx context.Context, link *Link) error {
	if r.db == nil {
		return fmt.Errorf("shortener create: %w", errors.New("database unavailable"))
	}
	err := r.db.WithContext(ctx).Create(link).Error
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicateCode
		}
		return fmt.Errorf("shortener create: %w", err)
	}
	return nil
}

// FindByLongURL returns the link with the given long URL, or ErrNotFound.
func (r *gormRepository) FindByLongURL(ctx context.Context, longURL string) (*Link, error) {
	if r.db == nil {
		return nil, fmt.Errorf("shortener find by long url: %w", errors.New("database unavailable"))
	}
	var link Link
	err := r.db.WithContext(ctx).Where("long_url = ?", longURL).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("shortener find by long url: %w", err)
	}
	return &link, nil
}

// FindByCode returns the link with the given code, or ErrNotFound.
func (r *gormRepository) FindByCode(ctx context.Context, code string) (*Link, error) {
	if r.db == nil {
		return nil, fmt.Errorf("shortener find by code: %w", errors.New("database unavailable"))
	}
	var link Link
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("shortener find by code: %w", err)
	}
	return &link, nil
}

// isDuplicateKeyError reports whether err is a MySQL duplicate-key error
// (error 1062).
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Error 1062")
}

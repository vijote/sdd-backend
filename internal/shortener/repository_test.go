package shortener

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/vijote/sdd-backend/internal/config"
	"github.com/vijote/sdd-backend/internal/database"
)

// startMySQL spins up a throwaway MySQL container and returns config pointing
// at it. Tests are skipped when Docker is unavailable (e.g. local CI-less runs).
func startMySQL(t *testing.T) *config.Config {
	t.Helper()

	ctx := context.Background()
	container, err := mysql.Run(ctx,
		"mysql:8.4",
		mysql.WithDatabase("sdd_backend"),
		mysql.WithUsername("root"),
		mysql.WithPassword("testpass"),
		testcontainers.WithEnv(map[string]string{
			"MYSQL_ROOT_HOST": "%",
		}),
	)
	if err != nil {
		t.Skipf("skip: could not start MySQL container (Docker unavailable?): %v", err)
	}
	t.Cleanup(func() {
		if termErr := container.Terminate(ctx); termErr != nil {
			t.Errorf("terminate container: %v", termErr)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}

	portNum, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("parse port %q: %v", port.Port(), err)
	}

	return &config.Config{
		DBHost:     host,
		DBPort:     portNum,
		DBUser:     "root",
		DBPassword: "testpass",
		DBName:     "sdd_backend",
	}
}

// setupRepo starts a MySQL container, migrates the schema, and returns a
// ready Repository. Skips when Docker is unavailable.
func setupRepo(t *testing.T) Repository {
	t.Helper()

	cfg := startMySQL(t)

	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("database.Open() error: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("database.Migrate() error: %v", err)
	}
	return NewRepository(db)
}

func TestIntegrationCreateAndFind(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupRepo(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	link := &Link{Code: "aB3xK9m", LongURL: "https://example.com/integration"}
	if err := repo.Create(ctx, link); err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if link.ID == 0 {
		t.Fatal("expected non-zero ID after Create")
	}

	byCode, err := repo.FindByCode(ctx, "aB3xK9m")
	if err != nil {
		t.Fatalf("FindByCode() error: %v", err)
	}
	if byCode.LongURL != "https://example.com/integration" {
		t.Errorf("FindByCode() long URL = %q", byCode.LongURL)
	}

	byLong, err := repo.FindByLongURL(ctx, "https://example.com/integration")
	if err != nil {
		t.Fatalf("FindByLongURL() error: %v", err)
	}
	if byLong.Code != "aB3xK9m" {
		t.Errorf("FindByLongURL() code = %q", byLong.Code)
	}
}

func TestIntegrationDuplicateLongURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupRepo(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	first := &Link{Code: "aaa1111", LongURL: "https://example.com/dup"}
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("first Create() error: %v", err)
	}

	// Same long URL, different code → duplicate-key error.
	second := &Link{Code: "bbb2222", LongURL: "https://example.com/dup"}
	err := repo.Create(ctx, second)
	if !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("expected ErrDuplicateCode, got %v", err)
	}

	// Same code, different long URL → duplicate-key error.
	third := &Link{Code: "aaa1111", LongURL: "https://example.com/other"}
	if err := repo.Create(ctx, third); !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("expected ErrDuplicateCode, got %v", err)
	}
}

func TestIntegrationNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := setupRepo(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := repo.FindByCode(ctx, "zzzzzzz"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if _, err := repo.FindByLongURL(ctx, "https://missing.example.com"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

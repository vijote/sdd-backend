package database

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/vijote/sdd-backend/internal/config"
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

func TestIntegrationOpenAndPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	cfg := startMySQL(t)

	db, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		t.Errorf("Ping() returned error: %v", err)
	}
}

func TestIntegrationMigrateIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	cfg := startMySQL(t)

	db, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}

	// First run migrates; second run must be a no-op that still succeeds.
	for run := 1; run <= 2; run++ {
		if err := Migrate(db); err != nil {
			t.Fatalf("Migrate() run %d returned error: %v", run, err)
		}
	}
}

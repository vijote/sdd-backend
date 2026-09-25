package config

import (
	"os"
	"testing"
)

// unsetEnv removes the given env vars for the duration of the test,
// restoring any previous values on cleanup.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		if old, ok := os.LookupEnv(key); ok {
			t.Cleanup(func() {
				if err := os.Setenv(key, old); err != nil {
					t.Errorf("restore %s: %v", key, err)
				}
			})
			if err := os.Unsetenv(key); err != nil {
				t.Fatalf("unset %s: %v", key, err)
			}
		}
	}
}

func TestLoadDefaults(t *testing.T) {
	unsetEnv(t, "PORT", "LOG_LEVEL", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "localhost")
	}
	if cfg.DBPort != 3306 {
		t.Errorf("DBPort = %d, want 3306", cfg.DBPort)
	}
	if cfg.DBUser != "root" {
		t.Errorf("DBUser = %q, want %q", cfg.DBUser, "root")
	}
	if cfg.DBPassword != "" {
		t.Errorf("DBPassword = %q, want empty", cfg.DBPassword)
	}
	if cfg.DBName != "sdd_backend" {
		t.Errorf("DBName = %q, want %q", cfg.DBName, "sdd_backend")
	}
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want 9090", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}

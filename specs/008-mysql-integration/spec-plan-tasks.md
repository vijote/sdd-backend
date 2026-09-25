---
name: 008-mysql-integration
description: GORM+MySQL persistence layer with discrete env config, startup migration via init container contract, and split liveness/readiness endpoints.
date: 2026-09-24
status: Implemented
---

# Spec · Plan · Tasks: MySQL Integration

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: `internal/database` package owns connection lifecycle, migration entry point, and health ping.
- **API (go-chi)**: `GET /healthz` stays liveness-only (no DB). New `GET /readyz` returns 200 when DB ping succeeds, 503 otherwise.
- **Persistence (GORM)**: `gorm.io/gorm` + `gorm.io/driver/mysql`. DSN composed from discrete env fields. No AutoMigrate at app startup — migrations run out-of-band (init container / K8s Job, managed in the infra repo); this repo exposes a `migrate` subcommand.
- **Kubernetes**: none in this repo. Migration execution contract: `sdd-backend migrate` (one-shot, exit 0 on success) to be wrapped by the infra repo's init container/Job.
- **Integration**: `testcontainers-go` spins a MySQL container for integration tests.

### Go Interfaces & Models
```go
// internal/database/database.go
type DB struct {
    *gorm.DB
}

// Open connects using discrete config fields; pings to verify.
func Open(cfg *config.Config) (*DB, error)

// Ping verifies connectivity (used by /readyz).
func (d *DB) Ping(ctx context.Context) error

// internal/database/migrate.go
// Migrate runs schema migrations. Idempotent; safe to run once per deploy.
func Migrate(db *DB) error

// internal/config/config.go — added fields
type Config struct {
    Port       int    `envconfig:"PORT" default:"8080"`
    LogLevel   string `envconfig:"LOG_LEVEL" default:"info"`
    DBHost     string `envconfig:"DB_HOST" default:"localhost"`
    DBPort     int    `envconfig:"DB_PORT" default:"3306"`
    DBUser     string `envconfig:"DB_USER" default:"root"`
    DBPassword string `envconfig:"DB_PASSWORD" default:""`
    DBName     string `envconfig:"DB_NAME" default:"sdd_backend"`
}
```

### API Contracts
```http
GET /healthz   → 200 {} (liveness only, unchanged)
GET /readyz    → 200 {} when DB ping OK
               → 503 {"error":"database unavailable"} when ping fails
```

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: `go build -v ./...` compiles without errors.
- [x] AC-002: `golangci-lint run` reports no issues.
- [x] AC-003: `go test -v ./...` passes all unit tests.
- [x] AC-004: `go test -v ./internal/database -run TestIntegration` passes with testcontainers MySQL (Open + Ping + Migrate).
- [x] AC-005: `go test -v ./internal/handlers` proves `/readyz` returns 200 with a live DB handle and 503 when ping fails (injected fake).
- [x] AC-006: `go run . migrate` exits 0 against a migrated testcontainers MySQL and is idempotent on second run.

### Assumptions & Constraints
- MySQL driver only (constitution).
- No migration lock in-app: single migration runner (init container) is the infra repo's responsibility.
- App process never mutates schema at startup.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `internal/config/config.go` | Update | Add DB_* fields |
| `internal/config/config_test.go` | Update | Cover DB_* defaults |
| `internal/database/database.go` | Create | Open/Ping via GORM |
| `internal/database/migrate.go` | Create | Migrate entry point |
| `internal/database/database_test.go` | Create | testcontainers integration tests |
| `internal/handlers/health.go` | Update | Add Ready handler with injected DB pinger |
| `internal/handlers/health_test.go` | Create | /readyz 200/503 tests |
| `internal/server/server.go` | Update | Wire /readyz route, accept DB dependency |
| `cmd/server/main.go` | Update | Open DB, wire into server; `migrate` subcommand |

### Architectural Layers
1. **Transport** — chi routes `/healthz`, `/readyz`; handlers depend on a `Pinger` interface (fake-able).
2. **Data** — `internal/database`: GORM open/ping/migrate.
3. **App** — `cmd/server/main.go`: subcommand dispatch (`serve` default, `migrate` one-shot).

### Verification Gates
- `go mod tidy && go build -v ./...`
- `golangci-lint run`
- `go test -v -cover ./...`

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in path (Depends on T###, ...)`

### Stage 1: Config & Data Layer
- [x] T001 [Stage 1: Config] Add DB_* fields in `internal/config/config.go` + update `internal/config/config_test.go`
- [x] T002 [Stage 1: DB] Create Open/Ping in `internal/database/database.go` (Depends on T001)
- [x] T003 [Stage 1: DB] Create Migrate in `internal/database/migrate.go` (Depends on T002)
- [x] T004 [Stage 1: DB] testcontainers integration tests in `internal/database/database_test.go` (Depends on T003)

### Stage 2: API & Wiring
- [x] T005 [Stage 2: API] Add Ready handler + Pinger interface in `internal/handlers/health.go` (Depends on T002)
- [x] T006 [Stage 2: API] /readyz 200/503 tests in `internal/handlers/health_test.go` (Depends on T005)
- [x] T007 [Stage 2: App] Wire DB + /readyz in `internal/server/server.go` and `cmd/server/main.go`; add `migrate` subcommand (Depends on T003, T005)

### Stage 3: Validation
- [x] T008 [Stage 3: Validate] Run `go build -v ./...` (Depends on T007)
- [x] T009 [Stage 3: Validate] Run `golangci-lint run` (Depends on T007)
- [x] T010 [Stage 3: Validate] Run `go test -v ./...` (Depends on T007)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Go structs, GORM tags, and interfaces clearly defined?
- [ ] CHK002 HTTP endpoints, verbs, and expected JSON structures fully documented?
- [ ] CHK003 Migration execution contract (subcommand + exit codes) specified for infra repo?

### Code Quality & Security
- [ ] CHK004 Proper error handling and bubbling (no silent failures)?
- [ ] CHK005 Context properly passed from HTTP handler to DB queries?
- [ ] CHK006 DB credentials only via env vars, never logged?
- [ ] CHK007 Environment variables parsed using `envconfig` with defaults?

### Machine-Verifiability
- [ ] CHK008 Corresponding `*_test.go` files for core logic?
- [ ] CHK009 Database interactions validated via `testcontainers-go`?
- [ ] CHK010 All ACs rely on `go test`, `go build`, or `golangci-lint`?

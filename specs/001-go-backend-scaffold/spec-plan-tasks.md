# Spec · Plan · Tasks: Go Backend Scaffold

**Branch**: `001-go-backend-scaffold` | **Date**: 2026-09-20 | **Status**: Approved

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: `internal/config` (envconfig config), `internal/server` (chi router + http.Server), `internal/handlers` (HTTP handlers)
- **API (go-chi)**: `GET /healthz` → 200 `{"status":"ok"}`
- **Persistence (GORM)**: Out of scope (follow-up spec)
- **Kubernetes (client-go)**: Out of scope (follow-up spec)
- **Integration**: Out of scope (follow-up spec)

### Go Interfaces & Models
```go
// internal/config/config.go
type Config struct {
	Port     int    `envconfig:"PORT" default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}
func Load() (*Config, error) // wraps envconfig.Process("", &c)

// internal/server/server.go
type Server struct {
	cfg *config.Config
}
func New(cfg *config.Config) *Server
func (s *Server) Start() error // http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: NewRouter(cfg)}
func NewRouter(cfg *config.Config) http.Handler // chi.NewRouter(); r.Get("/healthz", handlers.Health)

// internal/handlers/health.go
func Health(w http.ResponseWriter, r *http.Request) // 200, application/json, {"status":"ok"}
```

### API Contracts
```http
GET /healthz
```
**Response (200 OK):**
```json
{"status":"ok"}
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: `go build -v ./...` exits 0.
- [ ] AC-002: `golangci-lint run` exits 0.
- [ ] AC-003: `go test -v ./...` passes (config defaults/override, health handler, router 404).
- [ ] AC-004: `go build -o bin/server ./cmd/server` exits 0; `PORT=18080 ./bin/server &` then `curl -sf localhost:18080/healthz` returns `{"status":"ok"}`.
- [ ] AC-005: Push to `main` triggers `.github/workflows/ci.yml`; build, lint, and test steps all succeed.

### Assumptions & Constraints
- Module path `github.com/vijote/sdd-backend` (from git remote).
- Go 1.22+, chi v5, envconfig (no viper).
- JSON over HTTP.
- CI fully automated, no manual gates.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `go.mod` | Create | Module definition (go 1.22) |
| `.golangci.yml` | Create | Linter config (govet, errcheck, staticcheck) |
| `.gitignore` | Modify | Add `bin/`, `*.out`, `coverage.txt` |
| `internal/config/config.go` | Create | envconfig config struct + `Load()` |
| `internal/config/config_test.go` | Create | Unit tests: defaults, env override |
| `internal/handlers/health.go` | Create | Health handler |
| `internal/server/server.go` | Create | chi router + http.Server |
| `internal/server/server_test.go` | Create | httptest: `/healthz` 200, unknown route 404 |
| `cmd/server/main.go` | Create | Entrypoint wiring config → server |
| `.github/workflows/ci.yml` | Create | Automated CI on push/PR to main |

### Architectural Layers
1. **Transport** — chi router, JSON handlers.
2. **App** — `cmd/server/main.go` wires config → server.
3. **Config** — envconfig with typed defaults.

### Verification Gates
- `go mod tidy && go build ./...`
- `golangci-lint run`
- `go test -v -cover ./...`

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`
> Rules: 1 task = 1 file edit or 1 verification command. Same-stage independent tasks run in parallel.

### Stage 1: Module & Tooling
- [x] T001 [Stage 1: Module] `go mod init github.com/vijote/sdd-backend`
- [x] T002 [Stage 1: Tooling] Create `.golangci.yml`
- [x] T003 [Stage 1: Tooling] Update `.gitignore` with `bin/`, `*.out`, `coverage.txt`

### Stage 2: Config
- [x] T004 [Stage 2: Config] Create `internal/config/config.go` (Depends on T001)
- [x] T005 [Stage 2: Config] `go get github.com/kelseyhightower/envconfig` (Depends on T004)
- [x] T006 [Stage 2: Config] Create `internal/config/config_test.go` (Depends on T004)

### Stage 3: HTTP Layer
- [x] T007 [Stage 3: API] Create `internal/handlers/health.go` (Depends on T001)
- [x] T008 [Stage 3: API] Create `internal/server/server.go` (Depends on T007, T004)
- [x] T009 [Stage 3: API] `go get github.com/go-chi/chi/v5` (Depends on T008)
- [x] T010 [Stage 3: API] Create `internal/server/server_test.go` (Depends on T008)

### Stage 4: App Wiring & Validation
- [x] T011 [Stage 4: App] Create `cmd/server/main.go` (Depends on T008)
- [x] T012 [Stage 4: CI] Create `.github/workflows/ci.yml` (Depends on T002)
- [x] T013 [Stage 4: Validate] Run `go build -v ./...` (Depends on T011)
- [x] T014 [Stage 4: Validate] Run `golangci-lint run` (Depends on T011)
- [x] T015 [Stage 4: Validate] Run `go test -v ./...` (Depends on T011)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Go structs and function signatures defined with types and defaults?
- [ ] CHK002 HTTP endpoint, verb, and JSON response documented?

### Code Quality & Security
- [ ] CHK003 Error handling bubbled from `Load`/`Start` to `main` (no silent failures)?
- [ ] CHK004 Env vars parsed via envconfig with defaults?

### Machine-Verifiability
- [ ] CHK005 `*_test.go` files exist for config and server?
- [ ] CHK006 All ACs rely on `go build`, `golangci-lint`, `go test`, or curl?

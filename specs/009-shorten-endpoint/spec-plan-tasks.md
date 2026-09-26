---
name: 009-shorten-endpoint
description: POST /api/shorten endpoint — validated long URL, SHA-256→7 base62 code with collision retry, idempotent persistence via GORM Link model.
date: 2026-09-26
status: Implemented
---

# Spec · Plan · Tasks: Shorten Endpoint

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: `internal/shortener` package owns code generation (SHA-256 of long URL, truncated to 7 base62 chars; on collision re-hash with incrementing salt) and the service layer (validate → generate → persist, idempotent on duplicate long URL).
- **API (go-chi)**: `POST /api/shorten` accepts `{"url":"..."}`; returns `201 {"code":"aB3xK9m","long_url":"..."}`. Invalid body/URL → `400 {"error":"..."}`. No `BASE_URL` config — frontend derives the full URL from its own origin.
- **Persistence (GORM)**: `Link` model with unique index on `code` and `long_url`. Idempotency: lookup by `long_url` first; return existing code on duplicate. Registered in `internal/database/migrate.go` `models()`.
- **Kubernetes**: none.
- **Integration**: `testcontainers-go` MySQL for repository integration tests.

### Go Interfaces & Models
```go
// internal/shortener/shortener.go
type Link struct {
    ID        uint      `gorm:"primaryKey"`
    Code      string    `gorm:"size:7;uniqueIndex;not null"`
    LongURL   string    `gorm:"size:2048;uniqueIndex;not null"`
    CreatedAt time.Time
}

type Repository interface {
    Create(ctx context.Context, link *Link) error
    FindByLongURL(ctx context.Context, longURL string) (*Link, error)
    FindByCode(ctx context.Context, code string) (*Link, error)
}

type Service interface {
    Shorten(ctx context.Context, longURL string) (*Link, error)
}

// GenerateCode(longURL string, salt int) string — SHA-256 → first 7 base62 chars
// ValidateURL(raw string) error — non-empty, parseable, scheme http/https
```

### API Contracts
```http
POST /api/shorten
Content-Type: application/json

{"url": "https://example.com/some/long/path"}
```
**Response (201 Created):**
```json
{"code": "aB3xK9m", "long_url": "https://example.com/some/long/path"}
```
**Response (400 Bad Request):** `{"error":"invalid url"}` — empty, unparseable, or scheme not http/https.

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: `go build -v ./...` compiles without errors.
- [x] AC-002: `golangci-lint run` reports no issues.
- [x] AC-003: `go test -v ./...` passes all unit tests.
- [x] AC-004: `go test -v ./internal/shortener -run TestIntegration` passes with testcontainers MySQL (Create + FindByLongURL + FindByCode + duplicate idempotency).
- [x] AC-005: `go test -v ./internal/handlers` proves POST /api/shorten returns 201 with `{"code","long_url"}` on valid input, 400 on invalid URL, and the same code for a duplicate long URL.
- [x] AC-006: `go test -v ./internal/shortener` proves code generation is deterministic, 7 chars, base62, and collision retry with salt produces a distinct code.

### Assumptions & Constraints
- MySQL driver only; GORM (constitution).
- Code: exactly 7 base62 chars (`[0-9a-zA-Z]`).
- App process never mutates schema — `migrate` subcommand picks up `Link` via `models()`.
- No soft delete, no expiry in this spec.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `internal/shortener/shortener.go` | Create | Link model, Repository/Service interfaces, codegen, service impl |
| `internal/shortener/shortener_test.go` | Create | Codegen + validation + service unit tests (fake repo) |
| `internal/shortener/repository.go` | Create | GORM repository implementation |
| `internal/shortener/repository_test.go` | Create | testcontainers MySQL integration tests |
| `internal/database/migrate.go` | Update | Register `Link` in `models()` |
| `internal/handlers/shorten.go` | Create | POST /api/shorten handler (thin, delegates to Service) |
| `internal/handlers/shorten_test.go` | Create | 201/400/duplicate handler tests (fake service) |
| `internal/server/server.go` | Update | Wire `/api/shorten` route with Service dependency |
| `internal/server/server_test.go` | Update | Route wiring tests |

### Architectural Layers
1. **Transport** — chi route `POST /api/shorten`; handler depends on `shortener.Service` interface (fake-able).
2. **Service** — `internal/shortener`: validation, code generation + collision retry, idempotent persistence.
3. **Data** — GORM repository over `Link`; `models()` registration for the `migrate` subcommand.

### Verification Gates
- `go mod tidy && go build -v ./...`
- `golangci-lint run`
- `go test -v -cover ./...`

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in path (Depends on T###, ...)`

### Stage 1: Models & Data Layer
- [x] T001 [Stage 1: Core] Create Link model, codegen (SHA-256→7 base62, salt retry), validation in `internal/shortener/shortener.go`
- [x] T002 [Stage 1: Core] Codegen + validation unit tests in `internal/shortener/shortener_test.go` (Depends on T001)
- [x] T003 [Stage 1: DB] Create GORM repository in `internal/shortener/repository.go` (Depends on T001)
- [x] T004 [Stage 1: DB] Register Link in `models()` in `internal/database/migrate.go` (Depends on T001)
- [x] T005 [Stage 1: DB] testcontainers integration tests in `internal/shortener/repository_test.go` (Depends on T003)

### Stage 2: Service & API
- [x] T006 [Stage 2: Core] Implement Service (validate → generate → persist, idempotent) in `internal/shortener/shortener.go` (Depends on T003)
- [x] T007 [Stage 2: Core] Service unit tests with fake repo in `internal/shortener/shortener_test.go` (Depends on T006)
- [x] T008 [Stage 2: API] Create POST /api/shorten handler in `internal/handlers/shorten.go` (Depends on T006)
- [x] T009 [Stage 2: API] Handler 201/400/duplicate tests in `internal/handlers/shorten_test.go` (Depends on T008)
- [x] T010 [Stage 2: App] Wire route + Service in `internal/server/server.go` + `internal/server/server_test.go` (Depends on T008)

### Stage 3: Validation
- [x] T011 [Stage 3: Validate] Run `go build -v ./...` (Depends on T010)
- [x] T012 [Stage 3: Validate] Run `golangci-lint run` (Depends on T010)
- [x] T013 [Stage 3: Validate] Run `go test -v ./...` (Depends on T010)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Go structs, GORM tags, and interfaces clearly defined?
- [ ] CHK002 HTTP endpoint, verb, and expected JSON structures fully documented?
- [ ] CHK003 Code generation contract (7 base62, salt retry) specified?

### Code Quality & Security
- [ ] CHK004 Proper error handling and bubbling (no silent failures)?
- [ ] CHK005 Context properly passed from HTTP handler to DB queries?
- [ ] CHK006 Long URL validated (scheme allowlist http/https) before persistence?
- [ ] CHK007 No new env vars needed (defaults only)?

### Machine-Verifiability
- [ ] CHK008 Corresponding `*_test.go` files for core logic?
- [ ] CHK009 Database interactions validated via `testcontainers-go`?
- [ ] CHK010 All ACs rely on `go test`, `go build`, or `golangci-lint`?

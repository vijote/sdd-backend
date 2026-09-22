---
name: [NNN-feature-name]
description: [One-sentence overview of the spec]
date: [DATE]
status: Draft
---

# Spec · Plan · Tasks: [FEATURE_NAME]

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: [Domain models / Business logic / Packages]
- **API (go-chi)**: [HTTP routes / Handlers / Middleware / JSON Requests/Responses]
- **Persistence (GORM)**: [Struct tags / MySQL DDL / Migrations / Queries]
- **Kubernetes (client-go)**: [Client init / Resource interactions (Pods, ConfigMaps, etc.)]
- **Integration**: [GitHub Webhook validation (go-github) / External API calls]

### Go Interfaces & Models
```go
// Example interface
type [InterfaceName] interface {
	[MethodName](ctx context.Context, req [ReqStruct]) ([RespStruct], error)
}

// Example model
type [ModelName] struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Additional fields
}
```

### API Contracts
```http
[METHOD] /api/v1/[resource]
Authorization: Bearer [token]

{
  "field": "value"
}
```
**Response (200 OK):**
```json
{
  "data": { ... }
}
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: `go build -v ./...` compiles without errors.
- [ ] AC-002: `golangci-lint run` reports no warnings or errors.
- [ ] AC-003: `go test -v ./...` passes all unit tests.
- [ ] AC-004: DB schema migration tests pass successfully using `testcontainers-go`.
- [ ] AC-005: API endpoints return HTTP 200/201 on valid inputs and correct error codes on invalid inputs.

### Assumptions & Constraints
- Database driver is MySQL.
- Testing uses `testcontainers-go` for any DB operations.
- API is JSON over HTTP.
- GitHub webhook verification relies on HMAC SHA-256 signatures.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `cmd/[app]/main.go` | Modify | Entry point, server setup |
| `internal/[domain]/model.go` | Create | GORM domain models |
| `internal/[domain]/repository.go` | Create | Database interaction layer |
| `internal/[domain]/service.go` | Create | Business logic |
| `internal/api/handler.go` | Create | HTTP transport / chi router |
| `internal/[domain]/[file]_test.go` | Create | Unit / Integration tests |

### Architectural Layers
1. **Transport Layer** — `go-chi` routers, JSON marshaling, middleware, GitHub webhook validation.
2. **Service Layer** — Business logic, Go interfaces, orchestrating DB and K8s.
3. **Data Layer** — GORM models, repository pattern, MySQL interactions.
4. **Platform Layer** — Kubernetes interactions via `client-go`.

### Verification Gates
- `go mod tidy && go build ./...`
- `golangci-lint run`
- `go test -v -cover ./...`

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`
> Rules: 1 task = 1 file edit or 1 verification command. Same-stage independent tasks run in parallel. Stage N+1 requires full Stage N completion.

### Stage 1: Models & Data Layer
- [ ] T001 [Stage 1: DB] Create GORM struct in `internal/[domain]/model.go`
- [ ] T002 [Stage 1: DB] Create Repository interface and implementation in `internal/[domain]/repository.go` (Depends on T001)
- [ ] T003 [Stage 1: DB] Write integration test using `testcontainers-go` in `internal/[domain]/repository_test.go` (Depends on T002)

### Stage 2: Business Logic
- [ ] T004 [Stage 2: Core] Define Service interface in `internal/[domain]/service.go` (Depends on T002)
- [ ] T005 [Stage 2: Core] Implement business logic in `internal/[domain]/service.go` (Depends on T004)
- [ ] T006 [Stage 2: Core] Write unit tests for service in `internal/[domain]/service_test.go` (Depends on T005)

### Stage 3: API & Routing
- [ ] T007 [Stage 3: API] Implement HTTP handlers in `internal/api/handler.go` (Depends on T004)
- [ ] T008 [Stage 3: API] Register chi routes in `internal/api/router.go` (Depends on T007)
- [ ] T009 [Stage 3: API] Add webhook validation middleware (if applicable) in `internal/api/middleware.go`
- [ ] T010 [Stage 3: API] Write API handler tests in `internal/api/handler_test.go` (Depends on T007)

### Stage 4: App Wiring & Validation
- [ ] T011 [Stage 4: App] Wire dependencies and start chi server in `cmd/app/main.go` (Depends on T008)
- [ ] T012 [Stage 4: Validate] Run `go build -v ./...` (Depends on T011)
- [ ] T013 [Stage 4: Validate] Run `golangci-lint run` (Depends on T011)
- [ ] T014 [Stage 4: Validate] Run `go test -v ./...` (Depends on T011)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Go structs, GORM tags, and interfaces clearly defined?
- [ ] CHK002 HTTP endpoints, verbs, and expected JSON structures fully documented?
- [ ] CHK003 Kubernetes `client-go` RBAC and API dependencies specified?

### Code Quality & Security
- [ ] CHK004 Proper error handling and bubbling (no silent failures)?
- [ ] CHK005 Context (`context.Context`) properly passed down from HTTP handler to DB queries?
- [ ] CHK006 GitHub Webhooks validated via HMAC SHA-256?
- [ ] CHK007 Environment variables parsed using `envconfig` or `viper` with defaults?

### Machine-Verifiability
- [ ] CHK008 Are there corresponding `*_test.go` files for core logic?
- [ ] CHK009 Is the database interacting via `testcontainers-go` for real-world integration validation?
- [ ] CHK010 Do all Acceptance Criteria rely on `go test`, `go build`, or `golangci-lint`?

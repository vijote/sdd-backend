---
name: 009-1-backend-route-prefix
description: Register backend routes without the /api prefix so the dev ingress prefix-strip yields the intended public URL /api/shorten.
date: 2026-09-26
status: Implemented
---

# Spec · Plan · Tasks: Backend Route Prefix Fix

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No domain change. The dev ingress (`sdd-infra-v2`, path `/api(/|$)(.*)`, rewrite-target `/$2`) strips one `/api` before forwarding to `app-backend`. Backend must therefore register `POST /shorten` (not `/api/shorten`) so the public contract is `https://demo.vijote.dev/api/shorten`.
- **API (go-chi)**: `POST /shorten` — same request/response contract as 009 (`{"url"}` → `201 {"code","long_url"}` / `400 {"error"}`). `/healthz` and `/readyz` unchanged (probed in-cluster at unprefixed paths; publicly reachable at `/api/healthz`, `/api/readyz`).
- **Persistence (GORM)**: none.
- **Kubernetes**: none.
- **Integration**: none.

### API Contracts
```http
POST /api/shorten                      # public (through ingress)
→ POST /shorten                        # in-cluster (backend route)
Content-Type: application/json

{"url": "https://example.com/path"}
```
**Response (201 Created):** `{"code":"aB3xK9m","long_url":"https://example.com/path"}` — identical to 009.

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: `go build -v ./...` compiles without errors.
- [x] AC-002: `golangci-lint run` reports no issues.
- [x] AC-003: `go test -v ./...` passes all unit tests.
- [x] AC-004: `go test -v ./internal/server` proves `POST /shorten` returns 400 `{"error":"invalid url"}` for an invalid URL (route wired at unprefixed path).
- [x] AC-005: `go test -v ./internal/server` proves `POST /api/shorten` returns 404 (old prefixed route removed).

### Assumptions & Constraints
- Clean break: no alias for the old `/api/shorten` backend route (dev-only, no clients depend on it).
- Ingress manifests in `sdd-infra-v2` are NOT touched (decision A).

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `internal/server/server.go` | Update | Register `POST /shorten` instead of `POST /api/shorten` |
| `internal/server/server_test.go` | Update | Rewire `TestShortenRouteWired` to `/shorten`; add 404 test for old `/api/shorten` |
| `internal/handlers/shorten_test.go` | Update | Test path literal `/api/shorten` → `/shorten` (cosmetic; handler is path-agnostic) |

### Architectural Layers
1. **Transport** — chi route rename only; handler, service, repository untouched.

### Verification Gates
- `go build -v ./...`
- `golangci-lint run`
- `go test -v ./...`
- Live: `curl -sk -X POST https://demo.vijote.dev/api/shorten -H 'Content-Type: application/json' -d '{"url":"https://example.com"}'` → 201 (after deploy chain completes)

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in path (Depends on T###, ...)`

### Stage 1: Route Rename
- [x] T001 [Stage 1: Core] Rename route to `POST /shorten` in `internal/server/server.go`
- [x] T002 [Stage 1: Core] Update `TestShortenRouteWired` to `/shorten`; add `TestOldApiShortenRouteRemoved` (404) in `internal/server/server_test.go` (Depends on T001)
- [x] T003 [Stage 1: Core] Update path literal in `internal/handlers/shorten_test.go` (Depends on T001)

### Stage 2: Validation
- [x] T004 [Stage 2: Validate] Run `go build -v ./...` (Depends on T003)
- [x] T005 [Stage 2: Validate] Run `golangci-lint run` (Depends on T003)
- [x] T006 [Stage 2: Validate] Run `go test -v ./...` (Depends on T003)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Public vs in-cluster URL mapping documented?
- [ ] CHK002 Health probe paths explicitly unchanged?

### Code Quality & Security
- [ ] CHK003 No handler/service/repository changes (rename only)?
- [ ] CHK004 Old route removal covered by a 404 test?

### Machine-Verifiability
- [ ] CHK005 All ACs rely on `go test`, `go build`, or `golangci-lint`?
- [ ] CHK006 Live curl gate listed for post-deploy verification?

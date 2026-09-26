---
name: 010-redirect-endpoint
description: GET /{code} resolves a stored link and issues an HTTP 301 to the long URL; unknown or malformed codes return 404.
date: 2026-09-26
status: Implemented
---

# Spec · Plan · Tasks: Redirect Endpoint

---

## [SPEC] Redirect Endpoint

*   **[SPEC]** `GET /{code}` resolves a short code and issues an HTTP **301** redirect to the stored long URL. Unknown codes and malformed codes return **404** with a JSON error body. No DB query is made for malformed codes.
    *   **[PLAN]** `[UPDATE]` `internal/shortener/shortener.go` — add `Resolve(ctx, code)` to `Service` interface; implement via `repo.FindByCode`; add `ValidateCode` (7-char base62 check). `[CREATE]` `internal/handlers/redirect.go` — handler maps `ErrNotFound`/invalid code → 404 JSON, success → 301 + `Location`. `[UPDATE]` `internal/server/server.go` — register `r.Get("/{code}", ...)`. `[UPDATE]` tests: `internal/shortener/shortener_test.go`, `internal/handlers/redirect_test.go` (new), `internal/server/server_test.go`.
        *   **[TASK]** T001 [x] Add `ValidateCode(code string) bool` and `Resolve(ctx context.Context, code string) (*Link, error)` to shortener service (Depends on: none).
        *   **[TASK]** T002 [x] Unit tests for `ValidateCode` and `Resolve` (found / not-found / malformed) in `shortener_test.go` (Depends on T001).
        *   **[TASK]** T003 [x] Create `handlers.Redirect(svc)` returning 301 with `Location` header, empty body; 404 `{"error":"not found"}` for unknown/malformed (Depends on T001).
        *   **[TASK]** T004 [x] Handler tests in `internal/handlers/redirect_test.go`: 301 + Location on hit; 404 on `ErrNotFound`; 404 on malformed code with no repo call (Depends on T003).
        *   **[TASK]** T005 [x] Register `r.Get("/{code}", handlers.Redirect(...))` in `NewRouter`; router test: GET known code → 301, GET unknown → 404, POST /{code} → 405 (Depends on T003).
        *   **[TASK]** T006 [x] Full verification: `go build -v ./...` && `golangci-lint run` && `go test -v ./...` (Depends on T002, T004, T005).

### Go Contracts
```go
// internal/shortener/shortener.go
type Service interface {
    Shorten(ctx context.Context, longURL string) (*Link, error)
    Resolve(ctx context.Context, code string) (*Link, error) // ErrNotFound if unknown
}

// ValidateCode reports whether code is exactly CodeLength base62 chars.
func ValidateCode(code string) bool
```

### API Contract
```http
GET /{code}          # registered unprefixed; public URL: /api/{code} via ingress strip
```
**Response (301 Moved Permanently):** `Location: <long_url>`, empty body.
**Response (404 Not Found):** `{"error":"not found"}`

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: `go build -v ./...` compiles without errors.
- [ ] AC-002: `golangci-lint run` reports no warnings or errors.
- [ ] AC-003: `go test -v ./...` passes all unit tests.
- [ ] AC-004: Handler test asserts 301 + correct `Location` for a stored code.
- [ ] AC-005: Handler test asserts 404 for unknown code and for malformed (non-7-char base62) code, the latter without any repository call.

### Assumptions & Constraints
- 301 chosen over 302 (user decision); browsers may cache permanently — accepted.
- Code format is exactly 7 base62 chars (`[0-9a-zA-Z]{7}`); anything else is 404 pre-DB.
- Route registered unprefixed per 009-1 ingress strip convention.

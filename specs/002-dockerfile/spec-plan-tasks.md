---
name: 002-dockerfile
description: Multi-stage Dockerfile (golang:1.27-alpine → distroless) + .dockerignore to containerize the Go backend.
date: 2026-09-21
status: Approved
---

# Spec · Plan · Tasks: Dockerfile for Go Backend

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. Container packaging of the existing `cmd/server` binary.
- **Image**: Multi-stage build — `golang:1.27-alpine` (build) → `gcr.io/distroless/static-debian12` (runtime).
- **Binary**: `CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server`
- **Runtime**: `USER nonroot`, `EXPOSE 8080`, `ENTRYPOINT ["/server"]`.
- **Out of scope**: In-image `HEALTHCHECK` (distroless has no shell; K8s probes later), CI docker job, registry push, K8s deployment (follow-up specs).

### Dockerfile Contract
```dockerfile
# Build stage
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/server /server
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### .dockerignore Contract
```
.git/
.coda/
.github/
specs/
bin/
*.out
coverage.txt
.golangci.yml
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: `docker build -t sdd-backend:local .` exits 0.
- [ ] AC-002: `docker run -d --name sdd -p 18080:8080 sdd-backend:local && curl -sf localhost:18080/healthz` returns `{"status":"ok"}`.
- [ ] AC-003: `docker inspect --format '{{.Config.User}}' sdd-backend:local` returns `nonroot`.
- [ ] AC-004: `docker run -d --name sdd2 -e PORT=9090 -p 18090:9090 sdd-backend:local && curl -sf localhost:18090/healthz` returns `{"status":"ok"}` (env-configured port works).

**Cleanup after AC-002/AC-004:** `docker rm -f sdd sdd2`

### Assumptions & Constraints
- Go toolchain version from `go.mod` (`go 1.27.1`) → base image `golang:1.27-alpine`.
- Default `PORT=8080` (envconfig default in `internal/config/config.go`).
- Docker available locally for verification.
- No secrets in build context (`.dockerignore` excludes `.coda/`, `.git/`).

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `Dockerfile` | Create | Multi-stage build: alpine build → distroless static runtime |
| `.dockerignore` | Create | Minimal build context; exclude VCS, specs, tooling, artifacts |

### Architectural Layers
1. **Build stage** — `golang:1.27-alpine`: module cache via `go mod download`, static binary with stripped symbols.
2. **Runtime stage** — `gcr.io/distroless/static-debian12`: single binary, no shell, non-root user.

### Verification Gates
- `docker build -t sdd-backend:local .`
- `docker run` + `curl` against `/healthz` (default and env-configured port)
- `docker inspect` for non-root user

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`
> Rules: 1 task = 1 file edit or 1 verification command. Same-stage independent tasks run in parallel.

### Stage 1: Container Files
- [x] T001 [Stage 1: Docker] Create `Dockerfile` build stage
- [x] T002 [Stage 1: Docker] Add `Dockerfile` runtime stage (Depends on T001)
- [x] T003 [Stage 1: Docker] Create `.dockerignore`

### Stage 2: Validation
- [ ] T004 [Stage 2: Validate] Run `docker build -t sdd-backend:local .` (Depends on T002, T003)
- [ ] T005 [Stage 2: Validate] Run container, `curl -sf localhost:18080/healthz` (Depends on T004)
- [ ] T006 [Stage 2: Validate] `docker inspect` user check (Depends on T004)
- [ ] T007 [Stage 2: Validate] Run container with `PORT=9090`, `curl -sf localhost:18090/healthz` (Depends on T004)
- [ ] T008 [Stage 2: Validate] Cleanup: `docker rm -f sdd sdd2` (Depends on T005, T007)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Dockerfile stages, base images, and entrypoint documented?
- [ ] CHK002 `.dockerignore` excludes VCS, specs, and build artifacts?

### Code Quality & Security
- [ ] CHK003 Runtime image is distroless static with `USER nonroot`?
- [ ] CHK004 Binary built with `CGO_ENABLED=0` and stripped symbols (`-ldflags="-s -w"`)?

### Machine-Verifiability
- [ ] CHK005 All ACs rely on `docker build`, `docker run`, `curl`, or `docker inspect`?
- [ ] CHK006 No CI/registry/deployment scope leaked into this spec?

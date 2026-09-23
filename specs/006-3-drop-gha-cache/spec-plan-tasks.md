---
name: 006-3-drop-gha-cache
description: Remove GHA layer cache from the CI docker build step (docker driver does not support GHA cache export).
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Drop GHA Cache from Docker Build

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. One CI workflow fix.
- **Problem**: CI `docker` job fails at `Build and push` with `ERROR: failed to build: Cache export is not supported for the docker driver. Switch to a different driver, or turn on the containerd image store, and try again.` Spec 006 added `cache-from: type=gha` / `cache-to: type=gha,mode=max`, which require buildx's `docker-container` driver; `docker/build-push-action@v6` defaults to the plain `docker` driver.
- **Fix**: Remove both `cache-*` lines from the `Build and push` step. No driver change (user decision: drop the cache rather than switch drivers).
- **Out of scope**: buildx driver config, containerd image store, other steps.

### ci.yml Contract (delta — Build and push step)
```yaml
      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ${{ steps.ecr-login.outputs.registry }}/${{ env.AWS_ECR_REPOSITORY }}:${{ github.sha }}
            ${{ steps.ecr-login.outputs.registry }}/${{ env.AWS_ECR_REPOSITORY }}:latest
```

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: CI `docker` job green on push to `main` (no `Cache export is not supported` error; build + push succeed).

### Assumptions & Constraints
- Build time without layer cache is acceptable for this repo size (single small Go binary).
- All other `docker` job wiring (two-role auth, ecr-login, tags) unchanged.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Remove `cache-from` / `cache-to` from the `Build and push` step |

### Verification Gates
- CI run on `main`: `docker` job green.

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Fix
- [x] T001 [Stage 1: CI] Remove `cache-from: type=gha` and `cache-to: type=gha,mode=max` from the `Build and push` step in `.github/workflows/ci.yml`

### Stage 2: Validation
- [x] T002 [Stage 2: Validate] Push to `main`, confirm CI `docker` job green (AC-001) (Depends on T001)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Exact lines to remove documented?

### Machine-Verifiability
- [ ] CHK002 AC verifiable via CI run status?

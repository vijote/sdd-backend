---
name: 006-1-ecr-login-v2
description: Fix CI docker job by pinning aws-actions/amazon-ecr-login to v2 (v4 does not exist; latest is v2.1.7).
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Pin amazon-ecr-login to v2

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. One-line CI workflow fix.
- **Problem**: CI `docker` job fails with `Error: Unable to resolve action 'aws-actions/amazon-ecr-login@v4'; latest release is v2.1.7`. Spec 006 pinned `@v4` from knowledge (web search unavailable at implementation time); the action's latest major is v2.
- **Fix**: Change `aws-actions/amazon-ecr-login@v4` → `aws-actions/amazon-ecr-login@v2` in `.github/workflows/ci.yml`.
- **Out of scope**: other action versions, Dockerfile, ECR role/variables.

### ci.yml Contract (delta)
```yaml
      - name: Login to Amazon ECR
        uses: aws-actions/amazon-ecr-login@v2
        id: ecr-login
```

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: CI `docker` job green on push to `main` (no `Unable to resolve action` error; build + push succeed).

### Assumptions & Constraints
- `amazon-ecr-login@v2` exposes the same `registry` output used by the build-push step (v2 is the current major).
- All other spec 006 wiring (OIDC assume, tags, cache) unchanged.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | `amazon-ecr-login@v4` → `@v2` |

### Verification Gates
- CI run on `main`: `docker` job green.

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Fix
- [x] T001 [Stage 1: CI] Change `aws-actions/amazon-ecr-login@v4` to `aws-actions/amazon-ecr-login@v2` in `.github/workflows/ci.yml`

### Stage 2: Validation
- [x] T002 [Stage 2: Validate] Push to `main`, confirm CI `docker` job green (AC-001) (Depends on T001)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Exact action version documented?

### Machine-Verifiability
- [ ] CHK002 AC verifiable via CI run status?

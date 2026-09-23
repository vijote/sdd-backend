---
name: 004-golangci-lint-action-v7
description: Bump golangci-lint-action from v6 to v7 in CI to support the golangci-lint v2.13.2 binary pinned in spec 003.
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Bump golangci-lint-action to v7

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. CI workflow fix only.
- **Problem**: CI Lint step fails with `invalid version string 'v2.13.2', golangci-lint v2 is not supported by golangci-lint-action v6, you must update to golangci-lint-action v7`. Spec 003 pinned golangci-lint `v2.13.2` (v2-series), but `golangci-lint-action@v6` only supports v1.x binaries.
- **Fix**: Bump `golangci/golangci-lint-action` from `@v6` to `@v7` in `.github/workflows/ci.yml`.
- **Out of scope**: golangci-lint version change (stays `v2.13.2`), `.golangci.yml` changes, local verification (CI run is the acceptance gate).

### CI Workflow Contract (`.github/workflows/ci.yml`)
```yaml
      - name: Lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: "v2.13.2"
```

### Acceptance Criteria (machine-verifiable)
- [x] AC-001: CI `test` job green on push to `main` (Build, Lint, Test steps all pass; no `invalid version string` error).

### Assumptions & Constraints
- `golangci-lint` stays pinned at `v2.13.2` (spec 003).
- `golangci-lint-action@v7` supports golangci-lint v2.x binaries (per the error message's own instruction).
- `setup-go` stays at `1.27.1` (spec 003).

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Bump `golangci-lint-action` `@v6` → `@v7` |

### Verification Gates
- CI `test` job on `main` (Build → Lint → Test).

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: CI Workflow
- [x] T001 [Stage 1: CI] Change `uses: golangci/golangci-lint-action@v6` to `uses: golangci/golangci-lint-action@v7` in `.github/workflows/ci.yml`

### Stage 2: Validation
- [x] T002 [Stage 2: Validate] Push to `main`, confirm CI `test` job green (AC-001) (Depends on T001)

---

## [CHK] Technical Quality Gates

### Machine-Verifiability
- [ ] CHK001 AC-001 relies on the CI pipeline (build, lint, test) — no manual steps.

### Scope
- [ ] CHK002 Only `.github/workflows/ci.yml` is modified.
- [ ] CHK003 No Go source, `go.mod`, `.golangci.yml`, or golangci-lint version changes.

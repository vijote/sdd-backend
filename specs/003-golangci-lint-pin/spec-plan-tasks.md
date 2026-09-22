---
name: 003-golangci-lint-pin
description: Pin golangci-lint to v2.13.2 and setup-go to 1.27.1 in CI to fix the Go version mismatch lint failure.
date: 2026-09-22
status: Approved
---

# Spec · Plan · Tasks: Pin golangci-lint v2.13.2 in CI

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. CI workflow fix only.
- **Problem**: CI Lint step fails with `can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.27.1)`. `golangci-lint-action` with `version: latest` resolved to a binary built with Go 1.24; `go.mod` targets `go 1.27.1`.
- **Fix**: Pin `golangci-lint` to `v2.13.2` (built with Go ≥ 1.27.1) and align `setup-go` with `go.mod`.
- **Out of scope**: `.golangci.yml` changes (already `version: "2"` format), local verification (CI run is the acceptance gate).

### CI Workflow Contract (`.github/workflows/ci.yml`)
```yaml
      - uses: actions/setup-go@v5
        with:
          go-version: "1.27.1"

      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: "v2.13.2"
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: CI `test` job green on push to `main` (Build, Lint, Test steps all pass; no `can't load config` error).

### Assumptions & Constraints
- `go.mod` stays at `go 1.27.1`.
- `golangci-lint` v2.13.2 is built with Go ≥ 1.27.1 (user-confirmed version).
- `.golangci.yml` is already v2-format; no migration required.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Pin golangci-lint `v2.13.2`; bump `setup-go` to `1.27.1` |

### Verification Gates
- CI `test` job on `main` (Build → Lint → Test).

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: CI Workflow
- [x] T001 [Stage 1: CI] Set `version: "v2.13.2"` on the `golangci/golangci-lint-action@v6` step in `.github/workflows/ci.yml`
- [x] T002 [Stage 1: CI] Set `go-version: "1.27.1"` on the `actions/setup-go@v5` step in `.github/workflows/ci.yml`

### Stage 2: Validation
- [ ] T003 [Stage 2: Validate] Push to `main`, confirm CI `test` job green (AC-001) (Depends on T001, T002)

---

## [CHK] Technical Quality Gates

### Machine-Verifiability
- [ ] CHK001 AC-001 relies on the CI pipeline (build, lint, test) — no manual steps.

### Scope
- [ ] CHK002 Only `.github/workflows/ci.yml` is modified.
- [ ] CHK003 No Go source, `go.mod`, or `.golangci.yml` changes.

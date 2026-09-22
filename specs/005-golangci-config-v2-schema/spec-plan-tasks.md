---
name: 005-golangci-config-v2-schema
description: Migrate .golangci.yml to the v2 schema by moving issues.exclude-dirs to linters.exclusions.paths so golangci-lint config verify passes.
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Migrate .golangci.yml to v2 schema

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. Linter config migration only.
- **Problem**: CI Lint step fails `golangci-lint config verify` with `jsonschema: "issues" does not validate ... additional properties 'exclude-dirs' not allowed`. Spec 001's `.golangci.yml` uses `issues.exclude-dirs`, a v1-era key removed in the v2 schema. The v7 action (spec 004) now runs config verification, surfacing the invalid key.
- **Fix**: Per the v2 config docs, path exclusion lives at `linters.exclusions.paths` (regex patterns). Remove the `issues:` block; add `linters.exclusions.paths: ["^bin/"]` to preserve the original intent (no reported issues from `bin/`).
- **Out of scope**: linter set changes, `.github/workflows/ci.yml` changes, local verification (CI run is the acceptance gate).

### .golangci.yml Contract (v2 schema)
```yaml
version: "2"

run:
  timeout: 5m

linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - unused
    - ineffassign
  settings:
    govet:
      enable-all: true
      disable:
        - fieldalignment
  exclusions:
    paths:
      - ^bin/
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: CI `test` job green on push to `main` (Build, Lint, Test steps all pass; no `config verify` / `additional properties` error).

### Assumptions & Constraints
- golangci-lint stays pinned at `v2.13.2` (spec 003); action stays `@v7` (spec 004).
- `linters.exclusions.paths` entries are regex patterns matched against file paths; `^bin/` targets the top-level `bin/` directory only.
- Linter set and govet settings unchanged.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.golangci.yml` | Modify | Remove `issues.exclude-dirs`; add `linters.exclusions.paths: ["^bin/"]` |

### Verification Gates
- CI `test` job on `main` (Build → Lint → Test).

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Config Migration
- [x] T001 [Stage 1: Config] In `.golangci.yml`, delete the `issues:` block (`exclude-dirs: [bin]`) and add `exclusions.paths: ["^bin/"]` under `linters:`

### Stage 2: Validation
- [ ] T002 [Stage 2: Validate] Push to `main`, confirm CI `test` job green (AC-001) (Depends on T001)

---

## [CHK] Technical Quality Gates

### Machine-Verifiability
- [ ] CHK001 AC-001 relies on the CI pipeline (build, lint, test) — no manual steps.

### Scope
- [ ] CHK002 Only `.golangci.yml` is modified.
- [ ] CHK003 Linter set, govet settings, and CI workflow unchanged.

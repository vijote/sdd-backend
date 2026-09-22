# Current Session State

**Current Spec:**
- `specs/005-golangci-config-v2-schema` (T001–T002 `[x]` — DONE, CI green)

**Objective:**
- Migrate `.golangci.yml` to the v2 schema: remove `issues.exclude-dirs`, add `linters.exclusions.paths: ["^bin/"]` (fixes `config verify` failure under golangci-lint v2.13.2 / action v7).

**Context (Why):**
- CI Lint step failed `golangci-lint config verify`: `additional properties 'exclude-dirs' not allowed` — v1-era key removed in the v2 schema.

**Modified/Uncommitted Files:**
- `specs/003-golangci-lint-pin/spec-plan-tasks.md` (T003 + AC-001 marked done)
- `specs/004-golangci-lint-action-v7/spec-plan-tasks.md` (T002 + AC-001 marked done)
- `specs/005-golangci-config-v2-schema/spec-plan-tasks.md` (T002 + AC-001 marked done)

**Blockers/Unresolved Bugs:**
- None.

**Next Immediate Steps:**
- Commit the task-completion marks (specs 003/004/005).
- Spec 002 (dockerfile) T004–T008 validation still pending (docker build/run/curl/inspect).
- Start next spec: GORM/MySQL persistence (or CI docker job / GHCR push).

# Current Session State

**Current Spec:**
- `specs/005-golangci-config-v2-schema` (T001 `[x]`; T002 validation pending)

**Objective:**
- Migrate `.golangci.yml` to the v2 schema: remove `issues.exclude-dirs`, add `linters.exclusions.paths: ["^bin/"]` (fixes `config verify` failure under golangci-lint v2.13.2 / action v7).

**Context (Why):**
- CI Lint step failed `golangci-lint config verify`: `additional properties 'exclude-dirs' not allowed` — v1-era key removed in the v2 schema.

**Modified/Uncommitted Files:**
- `.golangci.yml` (issues.exclude-dirs → linters.exclusions.paths)
- `specs/005-golangci-config-v2-schema/spec-plan-tasks.md` (new spec)
- `.coda/feature.json` (active spec → 005)

**Blockers/Unresolved Bugs:**
- None.

**Next Immediate Steps:**
- T002: Push to `main`, confirm CI `test` job green (AC-001).
- Spec 002 (dockerfile) T004–T008 validation still pending (docker build/run/curl/inspect).
- Start next spec: GORM/MySQL persistence (or CI docker job / GHCR push).

# Current Session State

**Current Spec:**
- `specs/006-3-drop-gha-cache` (T001 `[x]` — DONE, validation T002 pending)

**Objective:**
- Remove GHA layer cache from the CI `docker` build step — the plain `docker` driver (build-push-action default) does not support GHA cache export.

**Context (Why):**
- Spec 006 added `cache-from: type=gha` / `cache-to: type=gha,mode=max`. CI failed: `ERROR: failed to build: Cache export is not supported for the docker driver.` User chose to drop the cache rather than switch to the `docker-container` driver.

**Modified/Uncommitted Files:**
- `.github/workflows/ci.yml` (removed `cache-from` / `cache-to` from Build and push step)
- `specs/006-3-drop-gha-cache/spec-plan-tasks.md` (new spec, status Implemented)
- `.coda/feature.json` (now points at 006-3)
- `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None.

**Next Immediate Steps:**
- T002: Push to `main`, confirm CI `docker` job green (AC-001).
- Then: K8s manifests / cluster deployment (follow-up spec).

# Current Session State

**Current Spec:**
- `specs/006-1-ecr-login-v2` (T001 `[x]` — DONE, validation T002 pending)

**Objective:**
- Fix CI `docker` job: `aws-actions/amazon-ecr-login@v4` does not exist (latest is v2.1.7); pin to `@v2`.

**Context (Why):**
- Spec 006 pinned `amazon-ecr-login@v4` from knowledge (web search unavailable at implementation time). First CI run failed: `Error: Unable to resolve action 'aws-actions/amazon-ecr-login@v4'; latest release is v2.1.7`. Follow-on fix spec per project convention.

**Modified/Uncommitted Files:**
- `.github/workflows/ci.yml` (ecr-login v4 → v2)
- `specs/006-1-ecr-login-v2/spec-plan-tasks.md` (new spec, status Implemented)
- `.coda/feature.json` (now points at 006-1)
- `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None.

**Next Immediate Steps:**
- T002: Push to `main`, confirm CI `docker` job green (AC-001).
- Then T005 of spec 006: `aws ecr describe-images` for the pushed SHA tag.
- Then: K8s manifests / cluster deployment (follow-up spec).

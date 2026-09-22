# Current Session State

**Current Spec:**
- `specs/006-2-ecr-two-role-auth` (T001–T002 `[x]` — DONE, validation T003 pending)

**Objective:**
- Switch CI `docker` job AWS auth to two-role OIDC chaining: bootstrap role → ECR target role (`role-chaining: true`), matching the user's Terraform project pattern.

**Context (Why):**
- Spec 006 assumed a single ECR role directly from OIDC. The user's standard pattern chains a shared bootstrap role → target role; this repo now matches it. ECR login/build-push/tags/cache unchanged.

**Modified/Uncommitted Files:**
- `.github/workflows/ci.yml` (env + bootstrap var; split assume step into bootstrap + chained target)
- `specs/006-2-ecr-two-role-auth/spec-plan-tasks.md` (new spec, status Implemented)
- `.coda/feature.json` (now points at 006-2)
- `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None. Prereqs for T003 (user-side, not spec tasks): new GitHub repo variable `AWS_BOOTSTRAP_ROLE_ARN`; bootstrap role trusts OIDC (`sub: repo:vijote/sdd-backend:ref:refs/heads/main`) + has `sts:AssumeRole` on the ECR target role; ECR target role trust policy changed from OIDC to role-to-role (bootstrap role) + 7 ECR push permissions.

**Next Immediate Steps:**
- T003: Push to `main`, confirm CI `docker` job green (AC-001).
- Then: K8s manifests / cluster deployment (follow-up spec).

# Current Session State

**Current Spec:**
- `specs/006-ecr-push` (T001–T003 `[x]` — DONE, validation T004–T005 pending)

**Objective:**
- Add a `docker` job to CI that builds the existing image (spec 002 Dockerfile) and pushes it to a pre-existing ECR repository via OIDC role assumption.

**Context (Why):**
- The image was never published to a registry, so the app could not be deployed to the k8s cluster. This spec closes the build→push gap. Auth uses OIDC (`id-token: write`) assuming a dedicated IAM role (`vars.AWS_ECR_ROLE_ARN`), no long-lived secrets.

**Modified/Uncommitted Files:**
- `.github/workflows/ci.yml` (permissions, env, workflow_dispatch, docker job)
- `specs/006-ecr-push/spec-plan-tasks.md` (new spec, status Implemented)
- `.coda/feature.json` (now points at 006)
- `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None. Prereqs for T004 (user-side, not spec tasks): GitHub repo variables `AWS_REGION`, `AWS_ECR_ROLE_ARN`, `AWS_ECR_REPOSITORY`; the dedicated IAM role with OIDC trust policy (`sub: repo:vijote/sdd-backend:ref:refs/heads/main`) + ECR push permissions; a `production` GitHub environment.
- Action majors pinned from knowledge (web search unavailable): `configure-aws-credentials@v6`, `amazon-ecr-login@v4`, `docker/build-push-action@v6` — verify latest majors before relying on them.

**Next Immediate Steps:**
- T004: Push to `main`, confirm CI `test` + `docker` jobs green (AC-001, AC-002).
- T005: `aws ecr describe-images` for the pushed SHA tag (AC-003).
- Then: K8s manifests / cluster deployment (follow-up spec).

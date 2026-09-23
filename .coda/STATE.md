# Current Session State

**Current Spec:**
- `specs/006-3-drop-gha-cache` (T001–T002 `[x]` — DONE, CI green, image pushed to ECR)

**Objective:**
- ECR push pipeline complete: CI `docker` job builds the image (spec 002 Dockerfile) and pushes to ECR via two-role OIDC auth. All specs in the 006 chain implemented and validated.

**Context (Why):**
- The image was never published to a registry, so the app could not be deployed to the k8s cluster. The 006 chain closed the build→push gap:
  - 006 — docker job in CI (build + push, SHA + latest tags)
  - 006-1 — pin `amazon-ecr-login@v2` (v4 does not exist)
  - 006-2 — two-role OIDC auth (bootstrap → ECR target, `role-chaining: true`)
  - 006-3 — drop GHA cache (plain docker driver does not support cache export)
- Final state: CI `test` + `docker` jobs green on `main`; image in ECR `sdd-k8s-platform/backend` tagged with commit SHA + `latest`.

**Modified/Uncommitted Files:**
- `specs/006-ecr-push/spec-plan-tasks.md` (T004–T005 + ACs checked)
- `specs/006-1-ecr-login-v2/spec-plan-tasks.md` (T002 + AC checked)
- `specs/006-2-ecr-two-role-auth/spec-plan-tasks.md` (T003 + AC checked)
- `specs/006-3-drop-gha-cache/spec-plan-tasks.md` (T002 + AC checked)
- `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None.

**Next Immediate Steps:**
- Commit the spec checkbox updates (chore commit).
- Start next feature spec: K8s manifests / cluster deployment (Deployment + Service + probes, image from ECR).

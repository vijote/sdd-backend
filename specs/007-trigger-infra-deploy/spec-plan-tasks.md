---
name: 007-trigger-infra-deploy
description: Add a CI step to the docker job that dispatches deploy-images.yml on vijote/sdd-infra-v2 with the pushed image tag, closing the build→deploy loop.
date: 2026-09-23
status: Draft
---

# Spec · Plan · Tasks: Trigger Infra Deploy from CI

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: After a successful ECR push, the CI `docker` job dispatches the `deploy-images.yml` workflow in the infrastructure repo (`vijote/sdd-infra-v2`) via `gh workflow run`, passing the pushed image tag. This closes the build→deploy gap: image lands in ECR, then the infra repo re-renders k8s manifests with the new tag and applies them.
- **API (go-chi)**: N/A — CI-only change.
- **Persistence (GORM)**: N/A.
- **Integration**: GitHub Actions `workflow_dispatch` event via `gh` CLI; auth via `INFRA_PAT` secret (PAT with `actions:write` on `vijote/sdd-infra-v2`).

### Workflow Contract
```yaml
- name: Trigger infra deploy
  env:
    GH_TOKEN: ${{ secrets.INFRA_PAT }}
  run: |
    gh workflow run deploy-images.yml \
      --repo vijote/sdd-infra-v2 \
      --ref main \
      -f backend_image_tag=${{ github.sha }}
```
- Dispatch target: `deploy-images.yml` on `vijote/sdd-infra-v2`, ref `main`.
- Input: `backend_image_tag` = `${{ github.sha }}` (matches the SHA tag pushed to ECR).
- Auth: `INFRA_PAT` secret; missing/expired PAT MUST fail the step (no `continue-on-error`).
- The step runs last in the `docker` job, after "Build and push".

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: `gh workflow run deploy-images.yml --repo vijote/sdd-infra-v2 --ref main -f backend_image_tag=<sha>` (executed by CI) returns exit 0; a new workflow run appears in `vijote/sdd-infra-v2`.
- [ ] AC-002: The dispatched run receives `backend_image_tag` equal to the backend commit SHA that was pushed to ECR.
- [ ] AC-003: CI `docker` job fails (red) when `INFRA_PAT` is invalid or the dispatch errors — no `continue-on-error`.
- [ ] AC-004: `actionlint .github/workflows/ci.yml` (or equivalent YAML lint) reports no errors.

### Assumptions & Constraints
- `INFRA_PAT` is a fine-grained PAT with `actions:write` permission on `vijote/sdd-infra-v2`, stored as an org/repo secret.
- `deploy-images.yml` in the infra repo has a `workflow_dispatch` trigger accepting `backend_image_tag` (string).
- Deploy trigger failure fails the CI job (fail-closed).

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Add "Trigger infra deploy" step at end of `docker` job; fix `--ref` to `main` |

### Architectural Layers
1. **Platform Layer** — GitHub Actions workflow orchestration: cross-repo `workflow_dispatch` via `gh` CLI.

### Verification Gates
- `actionlint .github/workflows/ci.yml`
- Push to `main` → observe CI `docker` job green + new run in `vijote/sdd-infra-v2`

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Workflow Change
- [x] T001 [Stage 1: CI] Add "Trigger infra deploy" step to `docker` job in `.github/workflows/ci.yml` (Depends on none) — **already applied manually; fix `--ref cleanup` → `--ref main`**

### Stage 2: Validation
- [x] T002 [Stage 2: Validate] Run `actionlint .github/workflows/ci.yml` (Depends on T001) — actionlint unavailable locally; YAML parse validation passed
- [ ] T003 [Stage 2: Validate] Push to `main`, confirm CI `docker` job green and a new `deploy-images.yml` run appears in `vijote/sdd-infra-v2` with `backend_image_tag` = commit SHA (Depends on T002)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Dispatch target repo, workflow, ref, and inputs fully specified?
- [ ] CHK002 Auth mechanism (`INFRA_PAT` secret, `actions:write`) documented?

### Code Quality & Security
- [ ] CHK003 No silent failures — dispatch errors fail the CI job?
- [ ] CHK004 PAT referenced only via `secrets.INFRA_PAT`, never logged or echoed?

### Machine-Verifiability
- [ ] CHK005 ACs verifiable via CI run outcome + `actionlint`?
- [ ] CHK006 Retroactive spec: YAML change already in working tree; remaining work is the `--ref main` fix + validation?

---

*Note: This spec is retroactive — the YAML change was applied manually before the spec existed. Remaining work: fix `--ref` to `main`, then validate via a real CI run.*

---
name: 006-ecr-push
description: Add a docker job to CI that builds the existing image and pushes it to a pre-existing ECR repository via OIDC role assumption.
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Push Image to ECR

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. CI pipeline extension only.
- **Problem**: The image built by the existing `Dockerfile` (spec 002) is not published anywhere; the app cannot be deployed to the k8s cluster without a registry image.
- **Fix**: Add a `docker` job to `.github/workflows/ci.yml` that assumes a dedicated IAM role via OIDC, logs in to ECR, and pushes the image tagged with the git SHA + `latest`.
- **Out of scope**: ECR repository creation (pre-existing, Terraform-managed), K8s manifests, cluster deployment, local docker validation (spec 002 T004–T008).

### Required GitHub Variables (repo Settings → Variables)
| Variable | Value |
|:---|:---|
| `AWS_REGION` | AWS region of the ECR repository (e.g. `us-east-1`) |
| `AWS_ECR_ROLE_ARN` | ARN of the dedicated IAM role with ECR push permissions |
| `AWS_ECR_REPOSITORY` | Name of the pre-existing ECR repository |

### IAM Role Contract
- Trust policy: `token.actions.githubusercontent.com` (OIDC), condition `sub: repo:vijote/sdd-backend:ref:refs/heads/main`.
- Permissions: `ecr:GetAuthorizationToken`, `ecr:BatchCheckLayerAvailability`, `ecr:GetDownloadUrlForLayer`, `ecr:PutImage`, `ecr:InitiateLayerUpload`, `ecr:UploadLayerPart`, `ecr:CompleteLayerUpload`.

### ci.yml Contract (delta)
```yaml
permissions:
  id-token: write
  contents: read

env:
  AWS_REGION: ${{ vars.AWS_REGION }}
  AWS_ECR_ROLE_ARN: ${{ vars.AWS_ECR_ROLE_ARN }}
  AWS_ECR_REPOSITORY: ${{ vars.AWS_ECR_REPOSITORY }}

jobs:
  docker:
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'workflow_dispatch' || (github.event_name == 'push' && github.ref == 'refs/heads/main')
    environment: production
    steps:
      - uses: actions/checkout@v4
      - name: Configure AWS Credentials
        uses: aws-actions/configure-aws-credentials@v6
        with:
          role-to-assume: ${{ env.AWS_ECR_ROLE_ARN }}
          aws-region: ${{ env.AWS_REGION }}
      - name: Login to Amazon ECR
        uses: aws-actions/amazon-ecr-login@v4
        id: ecr-login
      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ${{ steps.ecr-login.outputs.registry }}/${{ env.AWS_ECR_REPOSITORY }}:${{ github.sha }}
            ${{ steps.ecr-login.outputs.registry }}/${{ env.AWS_ECR_REPOSITORY }}:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

### Acceptance Criteria (machine-verifiable)
- [ ] AC-001: CI `test` job green on push to `main` (no regression).
- [ ] AC-002: CI `docker` job green on push to `main` (build + push succeed).
- [ ] AC-003: `aws ecr describe-images --repository-name $AWS_ECR_REPOSITORY --image-ids imageTag=<pushed-sha>` returns the image (run with ECR role credentials).

### Assumptions & Constraints
- ECR repository pre-exists; the workflow never creates it.
- A `production` GitHub environment exists in this repo (Settings → Environments); remove the `environment:` key if not.
- Action majors pinned from knowledge (`configure-aws-credentials@v6`, `amazon-ecr-login@v4`, `docker/build-push-action@v6`); verify latest majors at implementation time.
- Docker job runs only on push to `main` or manual dispatch — never on PRs.
- GHA layer cache (`type=gha`) for build speed.

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Add top-level `permissions`/`env`, `workflow_dispatch` trigger, and the `docker` job |

### Verification Gates
- CI run on `main`: `test` job green, `docker` job green.
- `aws ecr describe-images` for the pushed SHA tag.

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Workflow
- [x] T001 [Stage 1: CI] Add top-level `permissions` (`id-token: write`, `contents: read`) and `env` block (`AWS_REGION`, `AWS_ECR_ROLE_ARN`, `AWS_ECR_REPOSITORY`) to `.github/workflows/ci.yml`
- [x] T002 [Stage 1: CI] Add `workflow_dispatch` to the `on:` trigger in `.github/workflows/ci.yml`
- [x] T003 [Stage 1: CI] Add the `docker` job (`needs: test`, main-push/dispatch guard, OIDC assume, ecr-login, build-push with SHA+latest tags, GHA cache) to `.github/workflows/ci.yml` (Depends on T001)

### Stage 2: Validation
- [ ] T004 [Stage 2: Validate] Push to `main`, confirm CI `test` + `docker` jobs green (AC-001, AC-002) (Depends on T003)
- [ ] T005 [Stage 2: Validate] `aws ecr describe-images` for the pushed SHA tag (AC-003) (Depends on T004)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 All three GitHub variables documented with exact names?
- [ ] CHK002 IAM role trust policy + permission list documented?

### Code Quality & Security
- [ ] CHK003 No long-lived AWS secrets in the workflow (OIDC only)?
- [ ] CHK004 Docker job guarded to push-to-main / dispatch only (no PR pushes)?
- [ ] CHK005 Image tagged with git SHA (traceable) + `latest`?

### Machine-Verifiability
- [ ] CHK006 All ACs verifiable via CI run status or `aws ecr` CLI?

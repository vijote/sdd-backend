---
name: 006-2-ecr-two-role-auth
description: Switch CI docker job AWS auth to two-role OIDC chaining (bootstrap role → ECR target role with role-chaining).
date: 2026-09-22
status: Implemented
---

# Spec · Plan · Tasks: Two-Role OIDC Auth for ECR Push

---

## [SPEC] Technical Scope & Contracts

- **Core Logic**: No Go code changes. CI workflow auth change only.
- **Problem**: Spec 006 assumes a single ECR role directly from OIDC. The user's standard pattern (Terraform project) chains a bootstrap role → target role; this repo should match it.
- **Fix**: In the `docker` job, replace the single `configure-aws-credentials` step with two: assume `AWS_BOOTSTRAP_ROLE_ARN` (OIDC), then assume `AWS_ECR_ROLE_ARN` with `role-chaining: true`.
- **Out of scope**: ECR login/build-push steps, tags, cache, Dockerfile, other jobs.

### Required GitHub Variables (repo Settings → Variables)
| Variable | Value |
|:---|:---|
| `AWS_REGION` | AWS region of the ECR repository (unchanged) |
| `AWS_BOOTSTRAP_ROLE_ARN` | ARN of the shared OIDC-trusted bootstrap role (NEW) |
| `AWS_ECR_ROLE_ARN` | ARN of the ECR target role with push permissions (unchanged) |
| `AWS_ECR_REPOSITORY` | Name of the pre-existing ECR repository (unchanged) |

### IAM Role Contracts
- **Bootstrap role**: trust policy `token.actions.githubusercontent.com` (OIDC), condition `sub: repo:vijote/sdd-backend:ref:refs/heads/main`; permission `sts:AssumeRole` on the ECR target role ARN.
- **ECR target role**: trust policy = bootstrap role (role-to-role, `aws:PrincipalArn`); permissions `ecr:GetAuthorizationToken`, `ecr:BatchCheckLayerAvailability`, `ecr:GetDownloadUrlForLayer`, `ecr:PutImage`, `ecr:InitiateLayerUpload`, `ecr:UploadLayerPart`, `ecr:CompleteLayerUpload`.

### ci.yml Contract (delta — docker job)
```yaml
env:
  AWS_REGION: ${{ vars.AWS_REGION }}
  AWS_BOOTSTRAP_ROLE_ARN: ${{ vars.AWS_BOOTSTRAP_ROLE_ARN }}
  AWS_ECR_ROLE_ARN: ${{ vars.AWS_ECR_ROLE_ARN }}
  AWS_ECR_REPOSITORY: ${{ vars.AWS_ECR_REPOSITORY }}

jobs:
  docker:
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Configure AWS Bootstrap Credentials
        uses: aws-actions/configure-aws-credentials@v6
        with:
          role-to-assume: ${{ env.AWS_BOOTSTRAP_ROLE_ARN }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Assume ECR Target Role
        uses: aws-actions/configure-aws-credentials@v6
        with:
          role-to-assume: ${{ env.AWS_ECR_ROLE_ARN }}
          aws-region: ${{ env.AWS_REGION }}
          role-chaining: true

      - name: Login to Amazon ECR
        uses: aws-actions/amazon-ecr-login@v2
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
- [x] AC-001: CI `docker` job green on push to `main` (bootstrap assume → target assume → ECR login → build + push all succeed).

### Assumptions & Constraints
- `AWS_BOOTSTRAP_ROLE_ARN` is a new repo variable; the other three already exist.
- ECR target role trust policy changes from OIDC to role-to-role (bootstrap role) — user-side IAM change.
- `amazon-ecr-login@v2` and `docker/build-push-action@v6` unchanged (spec 006-1 fix preserved).

---

## [PLAN] Architecture Delta & File Impact

| File Path | Op | Purpose |
|:---|:---|:---|
| `.github/workflows/ci.yml` | Modify | Add `AWS_BOOTSTRAP_ROLE_ARN` to env; split single assume step into bootstrap + chained target assume |

### Verification Gates
- CI run on `main`: `docker` job green.

---

## [TASKS] Execution DAG

> Format: `- [ ] T### [Stage N: Label] Action in \`path\` (Depends on T###, ...)`

### Stage 1: Workflow
- [x] T001 [Stage 1: CI] Add `AWS_BOOTSTRAP_ROLE_ARN: ${{ vars.AWS_BOOTSTRAP_ROLE_ARN }}` to the top-level `env` block in `.github/workflows/ci.yml`
- [x] T002 [Stage 1: CI] In the `docker` job, replace the single `Configure AWS Credentials` step with `Configure AWS Bootstrap Credentials` (assume `AWS_BOOTSTRAP_ROLE_ARN`) + `Assume ECR Target Role` (assume `AWS_ECR_ROLE_ARN`, `role-chaining: true`) (Depends on T001)

### Stage 2: Validation
- [x] T003 [Stage 2: Validate] Push to `main`, confirm CI `docker` job green (AC-001) (Depends on T002)

---

## [CHK] Technical Quality Gates

### Contract Completeness
- [ ] CHK001 Both role ARN variables documented with exact names?
- [ ] CHK002 Bootstrap + target role trust policies and permissions documented?

### Code Quality & Security
- [ ] CHK003 No long-lived AWS secrets (OIDC + role chaining only)?
- [ ] CHK004 ECR target role no longer trusts OIDC directly (least privilege via bootstrap)?

### Machine-Verifiability
- [ ] CHK005 AC verifiable via CI run status?

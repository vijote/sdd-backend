---
name: specify
description: A condensed, interview-driven workflow to specify, plan, and task a feature in one go, designed for smaller context windows. Always writes the consolidated plan to specs/NNN-feature-name/spec-plan-tasks.md.
---

# specify: The Condensed Planning Workflow

You are executing a condensed, end-to-end specification and planning workflow. This skill combines the iterative interviewing of "grilling" with the structured output of `speckit-specify`, `speckit-plan`, and `speckit-tasks`.

## Phase 1: Load Governing Context

Before the interview begins, load `.coda/memory/constitution.md` as the immutable constraint set for this session.

Extract and apply the following from the constitution (silently — do not narrate to the user):
- **Zero Narrative Policy**: no conversational filler, marketing language, or verbose rationale at any point.
- **Token Efficiency**: total artifact length MUST be under 200 lines.
- **Machine-Verifiable Acceptance Gates**: every acceptance criterion must map to an executable CLI command (`go build`, `golangci-lint run`, `go test`, health checks, etc.).
- **Architecture First**: contracts (Go interfaces, REST schemas, DB models) must be explicit before tasks are listed.
- Any additional project-specific constraints defined in the constitution.

Use the constitution to **pre-answer** any interview questions it already settles (e.g., testing policy, CIDR ranges, IAM boundaries, state backend). Do not ask the user about things the constitution has already decided.

If `.coda/memory/constitution.md` does not exist, continue but note the absence in the completion report.

## Phase 2: The Iterative Interview (Grilling)

Your goal is to reach a shared understanding of the user's request, **within the boundaries set by the constitution**.
1. **The Design Tree**: Map the user's request as a tree of decisions.
2. **The Frontier**: Identify the "frontier" of questions — things not answered by the constitution or the codebase that you need to know *now*.
3. **Ask & Recommend**: Ask the frontier questions one round at a time. Number them and provide your recommended answer for each. Wait for the user's response before asking the next round.
4. **Codebase Context**: If a question can be answered by exploring the codebase, use your read tools to find the answer instead of asking the user.
5. **State Tracking (Small Context Optimization)**: Do not rely on implicit chat history. Keep a running, concise internal summary of "Settled Decisions" to avoid context bloat.

## Phase 3: The Nested Artifact Generation

Once the frontier is empty (all decisions settled), generate a **single, consolidated execution plan** as a markdown Artifact. Do **not** proceed to implementation yet.

**Always write the artifact to a spec file** — do not deliver it only via stdout. Create `specs/NNN-feature-name/spec-plan-tasks.md` where `NNN` is the next sequential number (scan `specs/` for the highest existing prefix) and `feature-name` is a short kebab-case slug. Follow the structure in `.coda/templates/spec-template.md` (frontmatter + SPEC / PLAN / TASKS / CHK sections). The frontmatter is required: `name` (spec id), `description` (one-sentence overview), `date`, `status` — it is the trace agents use to get a whole-picture view of all specs.

Format the artifact as a nested bulleted list to ensure perfect traceability. Each top-level bullet is a requirement (Spec), with its corresponding architectural changes (Plan) and steps (Tasks) nested beneath it.

**Format Template:**
*   **[SPEC]** <Requirement or Feature Point>
    *   **[PLAN]** <Files to `[CREATE]`, `[UPDATE]`, or `[DELETE]` and the architectural rationale>
    *   **[TASK]** <Actionable, step-by-step implementation tasks with `T###` IDs and `Depends on` annotations>

*Example:*
*   **[SPEC]** Add user authentication via JWT.
    *   **[PLAN]** `[CREATE]` `src/auth.js` for token validation logic. `[UPDATE]` `src/app.js` to apply the auth middleware.
    *   **[TASK]** T001 Install `jsonwebtoken` library. T002 Write `verifyToken` middleware in `auth.js` (Depends on T001). T003 Apply middleware to protected routes in `app.js` (Depends on T002).

Apply the constitution constraints throughout: no narrative, <200 lines total, all ACs as CLI commands.

## Phase 4: Human Confirmation
After writing the spec file, report its path and pause, explicitly asking the user for confirmation (e.g., "Spec written to `specs/NNN-feature-name/spec-plan-tasks.md`. Does this plan look correct? Should we proceed with implementation?"). Do not write any code until the user approves.

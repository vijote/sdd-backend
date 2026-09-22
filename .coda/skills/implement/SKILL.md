---
name: implement
description: >-
  Implement the pending tasks of the active spec. Triggered by chat — load this skill when the user says "implement" (bare or in phrases like "let's implement", "yes implement", "do it", "go ahead") with no other feature specified. Parses [SPEC], [PLAN], [TASKS] from the active spec-plan-tasks.md and executes the unchecked tasks.
compatibility: "Requires specs/<feature>/ directory with a spec-plan-tasks.md file"
---

# Implement Skill

Execute the unchecked tasks of the active spec. No ceremony — just implement.

## User Input

```text
$ARGUMENTS
```

`$ARGUMENTS` is optional. When the user triggers this skill via chat ("implement", "let's implement", "do it"), the argument is empty — that is the normal case.

## Steps

1. **Resolve the spec**:
   - If `$ARGUMENTS` names a feature folder or a `spec-plan-tasks.md` path, use it.
   - Otherwise read `.coda/feature.json` (`feature_directory`) and use that spec.
   - Set `ARTIFACT = <feature_dir>/spec-plan-tasks.md`. If it does not exist, stop and report the path.

2. **Collect pending tasks**: Parse `## [TASKS]` and take every unchecked line (`- [ ] T###`). If there are none, report "nothing to implement" and stop.

3. **Execute in dependency order** (respect `Depends on T###` annotations; same-stage independent tasks may be done together):
   - For each task, apply the governing contract from `## [SPEC]` and `## [PLAN]` (exact file paths, types, versions, values).
   - Create or edit exactly the file the task names. Match existing code style.
   - For verification tasks (push/CI/docker): do not auto-execute — surface the exact command for the user.
   - After each task, mark its line `- [x]` in `ARTIFACT`.
   - On failure: stop, report the task ID and error, wait for direction.

4. **Update tracking**: Set the spec frontmatter `status` to `Implemented` and update `.coda/STATE.md` (current spec, completed tasks, next steps).

## Completion Report

- **Spec**: feature directory
- **Tasks completed**: `T### — <description>` per `[x]` item
- **Files created/modified**: flat list of relative paths
- **Validation left to the user**: the AC-### commands from `## [SPEC]` as a copy-paste block
- **Suggested commit message**: one-liner, e.g. `feat(<feature-name>): <short description>`

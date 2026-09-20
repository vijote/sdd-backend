---
name: "speckit-app-implement"
description: "Execute implementation from a single spec-plan-tasks.md file (app-infra workflow). Parses [SPEC], [PLAN], [TASKS], and [CHK] sections from one consolidated artifact instead of separate spec/plan/tasks/checklist files."
compatibility: "Requires specs/<feature>/ directory with a spec-plan-tasks.md file"
metadata:
  author: "sdd-infra-v2"
  source: ".coda/templates/app-infra/spec-plan-tasks.md"
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).
The argument should be a feature folder name (e.g., `003-mysql-statefulset`) or a path to a `spec-plan-tasks.md` file.

---

## Pre-Execution Checks

1. **Resolve the feature directory**:
   - If `$ARGUMENTS` is a folder name, resolve it as `specs/<feature>/`.
   - If `$ARGUMENTS` is a file path ending in `spec-plan-tasks.md`, use its parent as `FEATURE_DIR`.
   - If `$ARGUMENTS` is empty, scan `specs/` for any folder containing `spec-plan-tasks.md` with unchecked `[ ]` tasks; if exactly one is found use it, otherwise ask the user to specify.
   - Set `ARTIFACT = FEATURE_DIR/spec-plan-tasks.md`. If it does not exist, stop and report the path.

2. **Parse sections** from `ARTIFACT` into four in-memory blocks:
   - `SPEC_BLOCK`  — content under `## [SPEC]`
   - `PLAN_BLOCK`  — content under `## [PLAN]`
   - `TASKS_BLOCK` — content under `## [TASKS]`
   - `CHK_BLOCK`   — content under `## [CHK]`

3. **Quality Gate — CHK block**:
   - Scan `CHK_BLOCK` for any unchecked item (`- [ ] CHK###`).
   - If blocking unchecked items exist, list them and ask the user to confirm before proceeding.
   - If the user declines, stop.

4. **Check for extension hooks** (`before_implement`):
   - If `.coda/extensions.yml` exists, read `hooks.before_implement`.
   - Skip hooks where `enabled: false`. Treat missing `enabled` as `true`.
   - Skip hooks with a non-empty `condition` (leave evaluation to HookExecutor).
   - Replace `.` with `-` in command names when constructing slash commands.
   - For optional hooks: display the hook block and let the user invoke manually.
   - For mandatory hooks: emit the hook block and invoke it; wait for completion before continuing.
   - If no hooks or no file: skip silently.

---

## Execution Loop

Load `.coda/memory/constitution.md` once as the immutable context for all tasks.

For each task in `TASKS_BLOCK` (in dependency order — respect `Depends on T###` annotations):

### Per-Task Isolation Protocol

1. **Extract task payload**:
   - Task line: `- [ ] T### [Stage N: Label] <action> in <path> (Depends on ...)`
   - `TARGET_FILE` = the file path named in the task line.
   - `CONTRACT_SNIPPET` = the relevant excerpt from `SPEC_BLOCK` (Go struct/interface, API route, DB model, AC-### item) and `PLAN_BLOCK` (file impact row for `TARGET_FILE`) that governs this task.

2. **Execute 1:1 implementation**:
   - Create or edit exactly `TARGET_FILE` (Go `.go`, DB `.sql`, shell script, etc.).
   - Apply `CONTRACT_SNIPPET` as the authoritative schema: variable types, API versions, resource names, namespaces, StorageClass, credentials references, CIDR values.
   - Follow standard Go formatting. Run `go fmt` mentally on Go output.
   - For verification tasks (no `TARGET_FILE`): output the exact CLI command(s) to run and their expected output. Do **not** auto-execute unless the user has explicitly granted permission.

3. **Mark progress**:
   - Update the task line in `ARTIFACT` from `- [ ]` to `- [x]`.
   - On failure: stop immediately, report the error with the task ID and file path, and wait for user direction before resuming.

4. **Parallelism**:
   - Tasks within the same Stage that share no `Depends on` overlap may be implemented in a single response block.
   - Stage N+1 tasks must not start until all Stage N tasks are `[x]`.

---

## Completion

1. Verify all `T###` items in `TASKS_BLOCK` are `[x]`.
2. Verify all `AC-###` items in `SPEC_BLOCK` have been addressed (not necessarily passed — validation is the user's responsibility per constitution).
3. Update `ARTIFACT` with all `[x]` marks persisted to disk.

**Do not run acceptance criteria commands automatically.** Surface them as a copy-paste block for the user.

---

## Post-Execution Hooks

Check `.coda/extensions.yml` for `hooks.after_implement` and process with the same rules as pre-execution hooks.

---

## Completion Report

Output a final summary containing:
- **Feature**: `FEATURE_DIR`
- **Tasks completed**: `T### — <description>` for each `[x]` item
- **Files created/modified**: flat list with relative paths
- **Acceptance criteria to validate** (copy-paste CLI block from `SPEC_BLOCK` AC-### items)
- **Suggested commit message**:
  ```
  feat(<feature-name>): implement <short description>

  Co-Authored-By: CODA <coda@globant.com>
  ```

All runtime validation and acceptance testing is handled by the user directly against AWS infrastructure.

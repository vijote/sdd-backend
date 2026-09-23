<!--
Sync Impact Report:
Version change: 3.0.0 → 4.0.0 (MAJOR: Pivot to URL shortener backend objective)
Modified principles: Core Principles 2, 7
Added sections: None
Removed sections: None
Follow-up TODOs: None
-->

# System & LLM Execution Directives

## Core Principles

### 1. Zero Narrative Policy
Skip all introductory conversational filler, user personas, marketing justifications, and high-level product narratives. Go directly to technical engineering contracts, Go interfaces, API specifications, and machine-verifiable acceptance criteria.

### 2. Architecture First & Explicit Engineering Contracts
Use explicit engineering jargon, precise file paths, and exact Go package/struct names. Every specification and plan MUST define:
- Exact Go structs, interfaces, and function signatures (with types, defaults, and constraints)
- HTTP API schemas using `go-chi/chi/v5` router paths and methods
- Database models and persistence logic using GORM (`gorm.io/gorm` and `gorm.io/driver/mysql`)
- Configuration structures via `envconfig` or `viper`
- Short code generation: SHA-256 hash of the long URL, truncated to 7 base62 characters, with collision detection and retry (re-hash with an incrementing salt) before persistence
- Redirect contract: `GET /{code}` MUST issue an HTTP 301/302 to the stored long URL; unknown codes MUST return 404

### 3. Payload & Token Efficiency
Keep spec, plan, and architecture delta artifacts strictly below 200 lines. Use compact markdown tables, bullet points, and code blocks. Non-frontier and local LLMs (e.g. GLM-4.6, Qwen-Coder) must not experience reasoning degradation from bloated context windows.

### 4. Granular Dependency Tree (Micro-DAG)
Write `tasks.md` as an acyclic dependency graph (DAG) where every task corresponds to a 1:1 file edit, Go package, API endpoint, or testing step with explicit dependency pointers:
`- [ ] T001 [Stage] Task description in path/to/file (Depends on Txxx)`

### 5. Machine-Verifiable Acceptance Gates
Never use vague adjectives ("robust", "scalable", "fast"). All acceptance criteria MUST be machine-verifiable through automated CI/CD commands:
- Go build and compilation (`go build -v ./...`)
- Linting (`golangci-lint run`)
- Unit and Integration tests (`go test -v ./...`)
- API endpoints returning expected HTTP status codes and JSON structures
- Database queries verifying expected persistence states

### 6. Testing Policy
- **Automated Testing Required**: Generate native Go unit tests (`*_test.go`) for business logic and routing.
- **Integration Testing**: Utilize `testcontainers-go` for database integration testing (e.g., spinning up a MySQL container for GORM operations).
- **Validation Steps**: API logic should be verifiable through Go's native testing package.

### 7. CI/CD Automation Policy
- **No Manual Approval Gates**: All deployments and test suites must run fully automated on main branch push.
- **Input Validation via Automation**: Long URL validation (scheme allowlist: `http`/`https`, non-empty, parseable) must be rigorously tested and validated in automation.

### 8. Prefer Standard Go Tooling Over Custom Scripts
- **Tooling-First Execution**: Use `go test`, `go build`, `go mod`, and `golangci-lint` as the primary mechanisms for validation and state management.

## Session Isolation Protocol
When implementing tasks with LLM agents, load only the minimal context payload:
1. Directive: `.coda/memory/constitution.md`
2. Active Single Task: Target file, action goal, and exact data/infrastructure contract
3. Instruction: Output ONLY the implementation code and verification command. No conversational text outside code blocks.

## Governance
This constitution is the non-negotiable governing standard for all artifacts in this repository. All PRs, plans, specifications, and task graphs must strictly adhere to these directives.

**Version**: 4.0.0 | **Ratified**: 2026-09-19 | **Last Amended**: 2026-09-23

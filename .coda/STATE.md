# Current Session State

**Current Spec:**
- `specs/009-1-backend-route-prefix` (T001–T006 `[x]` — DONE, build/lint/test green locally)

**Objective:**
- Register backend routes without the `/api` prefix (`POST /shorten`) so the dev ingress prefix-strip yields the intended public URL `https://demo.vijote.dev/api/shorten`. Old `/api/shorten` backend route removed (404-tested).

**Context (Why):**
- Follow-on fix to 009: live smoke test showed the ingress strips one `/api`, so the backend's `/api/shorten` was only reachable at `/api/api/shorten`. Decision A: fix backend-side, ingress untouched.

**Modified/Uncommitted Files:**
- `internal/server/server.go` (route rename)
- `internal/server/server_test.go` (`/shorten` wired test + old-route 404 test)
- `internal/handlers/shorten_test.go` (path literal)
- `specs/009-1-backend-route-prefix/spec-plan-tasks.md` (all tasks/ACs checked)
- `.coda/feature.json`, `.coda/STATE.md`

**Blockers/Unresolved Bugs:**
- None. (golangci-lint not on PATH — use `~/go/bin/golangci-lint`.)

**Next Immediate Steps:**
- Commit spec 009-1 (one-liner, fix group) and push — deploy chain picks it up.
- Live verify after deploy: `curl -sk -X POST https://demo.vijote.dev/api/shorten -H 'Content-Type: application/json' -d '{"url":"https://example.com"}'` → 201.
- `/specify` for 010 — redirect endpoint (GET /{code} → 301/302 to stored long URL, 404 unknown).

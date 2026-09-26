## Stable Project Preferences
- Before making edits, suggest creating a new spec using `/specify`. Do not make edits or push fixes on your own without direct user instruction.
- Commit style: one-liner commit messages; group files by feature — files from the same spec or prompt go in the same commit.
- Bug/finding handling: when a bug or gap is found during implementation or verification, fix it in a NEW follow-on spec with the NEXT sequential number (e.g., `003-4`, `003-5`), NOT by editing already-implemented specs to match code. Reserve `00x-0-*` strictly for bootstrap/prerequisite fixes.

## Important Decisions

## Learned Facts & Compatibility Rules

## Known Gotchas
- sdd-infra-v2 dev ingress strips one /api prefix (rewrite-target /$2 on path /api(/|$)(.*)), so backend routes registered under /api/... are only reachable publicly at /api/api/... — backend should register routes without the /api prefix (e.g. /shorten); live domain: https://demo.vijote.dev
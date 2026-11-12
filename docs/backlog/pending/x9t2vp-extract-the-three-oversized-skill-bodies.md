---
title: 'Extract the three oversized skill bodies'
priority: p3
section: ralvaskills
topics:
    - skills
    - tokens
---
The 2026-11-11 split audit found no skill that should become a second skill: median `SKILL.md` is 112 lines, one exceeds the 250-line budget, and cross-skill duplication tops out at 5% of shared 8-grams. What it did find is one reference manual hiding inside a body and three sections that restate a sibling skill. Extracting them is the whole split-shaped change; do not create new skills.

- **`go-architect`** (29.6 KB, 3× the 10 KB target): §14 Tooling is 7.5 KB — a linter catalogue, not a convention — so it moves whole to `GO_TOOLING.md`. §16 `go.work` (3.3 KB) joins the existing `LAYOUTS.md`. §12 sqlc/sqlx (1.7 KB) becomes a pointer to `sql-architect` with the code in `RECIPES.md`. Target ~17 KB. Bump 1.20.0 → 1.21.0.
- **`mcp-architect`** (32.2 KB, the largest body): §9 Authorization → `AUTH.md`, mirroring `rest-api-architect/AUTH_PATTERNS.md`; §13 Testing → `RECIPES.md`. Target ~27 KB. Bump 2.4.0 → 2.5.0.
- **`repo-tooling-architect`**: §8 CI (1.6 KB) restates `ci-cd-architect`; replace it with a pointer. Bump 1.6.0 → 1.7.0.

Rejected and not to be revisited without new information: splitting any stack skill into `-standards` (26 of them; no reader gain, every bundle breaks), `go-modules` (duplicates `LAYOUTS.md`), merging `performance-reviewer` + `backend-performance-architect` (deliberate report/design pair), and merging `demo-script-*` + `demo-deck-*` (different artifacts).

Done when each body is at its target size, the moved files are `UPPER_SNAKE_CASE` inside the skill folder, `task tokens` is regenerated, and the three minor releases are published.

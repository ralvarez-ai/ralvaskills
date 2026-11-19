---
title: 'MCP surface: 3 tools, no split — the always-on cost here is the catalogue'
priority: p3
section: ralvaskills
topics:
    - mcp-surface
---
Measured 2026-11-12 against `mcp-architect` 2.5.0 §3c / RECIPES §14.

- `rsk mcp`: 3 tools (`project_profile`, `search_skills`, `install_skills`, the last destructive) plus the `rsk://catalog` resource. Schemas are inline in Go; there are no JSON contract files to measure. That is ~3k tokens at worst.
- **Split?** No, and it should stay one server: the surface is three tools.

The token weight of this repo is elsewhere, and §3c does not apply to it: 67 skill descriptions load on every turn if all are installed (~7.6k tokens), the `global` bundle baseline is 1226/1500, and each project's `.rsk/SKILLS.md` has a 2000-token budget. Those are tracked by `docs/TOKEN_COUNTS.md`, `rsk status` and the existing digest work.

Exit: none for the MCP. Re-open only if the server passes ~20 tools (it would need a new capability family), and then apply the same audit: catalog file, grouped registration and the merge check.

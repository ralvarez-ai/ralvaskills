---
title: Remove the rsk MCP entry from .mcp.json on destroy
priority: p2
section: ralvaskills
topics:
    - rsk
    - mcp
---
`rsk destroy --yes` in a project created with `new --for claude-code` removes
`.rsk/`, drops the `@.rsk/CLAUDE.md` import and cleans `opencode.json`, but
leaves `.mcp.json` in place with the `rsk` server still registered. The command
reports "removed .rsk/ and cleaned up tool configs", so the leftover is silent.

Reproduce:

    rsk new --for claude-code && rsk destroy --yes && ls -A

Expected: no `.mcp.json`, or a `.mcp.json` with the rsk entry removed (other
servers preserved). Found while writing the binary-level e2e tier; the tier
asserts the `.rsk` and `CLAUDE.md` half of this and should gain the `.mcp.json`
assertion when it is fixed.

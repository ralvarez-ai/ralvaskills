---
title: 'Fix the code drifts the landing page work found'
priority: p2
section: ralvaskills
topics:
    - code
    - docs
    - owner-decision
---
Found while reading the code, recorded in the design doc's flags, not fixed because some touch published or adapter code (AGENTS.md: ask first):

- README says v3.3.0; the CLI is 3.4.1.
- The `rsk://catalog` description and three tool-input descriptions are in Spanish.
- `llm-app` names two skills that do not exist (`llm-app-architect`, `agent-architect`).
- AGENTS.md says `.claude/settings.json`, but the code writes `.mcp.json`.
- `rsk init --force` help mentions `rsk.mod`.
- `rsk new` answers yes without a TTY.
- The 9 personal skill names are public in the registry index.

Done when each is fixed or explicitly accepted, and `gen-data.mjs` and the page text still agree with the code.

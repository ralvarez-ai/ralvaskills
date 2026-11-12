---
title: 'Bring the global and bundle description budgets back under their limits'
priority: p2
section: ralvaskills
topics:
    - tokens
    - catalog
---
`docs/TOKEN_COUNTS.md` (regenerated 2026-11-06) shows the global baseline at about 1740 of 1500 description tokens, so it is over budget, and the `fastapi`, `go-grpc`, `gin` and `nethttp` bundles go over the 2200 session budget. Shorten descriptions (they load every turn for every installed skill) or rebalance what the `global` bundle contains. `token-economy` was added to `global` and accounts for part of the overage.

Done when `task tokens` reports no `OVER` row, or the budgets are changed on purpose and the reason is written down.

2026-11-12: closed by shortening descriptions, not by moving the budget. Ten of the 17 global descriptions were rewritten to the same what-and-when in fewer words (`token-economy` 230 → 73, `rsk-guide` 199 → 71, and eight more), which took the baseline from ~1867 to **1226 / 1500 (`ok`)** and put every bundle under the 2200 session budget — the worst is `nethttp` at 1839, where seven were over. `task tokens` regenerates the table with no `OVER` row. The bodies are untouched; only `description:` lines and patch bumps changed.

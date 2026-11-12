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

2026-11-12: `task tokens` regenerated the table after the ten renames and after the eight missing `global` members were installed. The global baseline is now **17 skills, ~1867 / 1500 desc tokens** (was ~1740), so the overage grew by about 127 tokens a turn; `go-cli`, `llm-app`, `python-grpc`, `fastapi`, `go-grpc`, `gin` and `nethttp` are over the 2200 session budget. Shortening the 17 global descriptions is now the cheapest fix — the bundle is complete and the page states the number.

---
title: 'Warn when .rsk/SKILLS.md exceeds its token budget'
priority: p3
section: ralvaskills
topics:
    - cli
    - tokens
---
`digest.Sync` returns the estimated size of the generated file but nothing uses it. `digest.BudgetTokens` is 2000 (about 8 KB; two typical digests come to about 840 tokens, so roughly five pinned skills reach it). Print a warning from `rsk pin`, `install` and `update` when the estimate passes the budget, naming the largest digests, and show the size in `rsk status`.

Done when pinning a sixth typical skill prints the warning on stderr and `--output json` stays clean.

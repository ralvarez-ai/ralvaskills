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

Resolved 2026-11-11: `digest.Measure` renders the entries `Sync` would write and returns a `digest.Size` — the whole-file estimate plus each digest's own, largest first — so nothing is written to measure. `cmd/rsk/budget.go` turns that into one line; `rsk pin`, both project install paths (`reportProjectInstall`) and `rsk update` (after it re-renders) print it to stderr when it passes `BudgetTokens`, and `rsk status` shows it under the tables, muted inside the budget and as a warning (with the three largest digests) past it. Verified with a built binary: pinning a sixth ~620-token skill printed `⚠ .rsk/SKILLS.md ~3782 tokens of the 2000-token budget — largest: alpha 620, bravo 620, charlie 620` on stderr while stdout kept only the success line, `rsk status -o json` wrote nothing to stderr and kept the budget line out of the payload, and a one-skill project printed the muted `~41 tokens of the 2000-token budget`.

This repo's own file is now known to be over budget: `rsk status` reports ~3046 tokens, largest `repo-tooling-architect` 745, `docs-architect` 570, `token-economy` 526, so the warning fires on every `rsk update` here until the digests are trimmed (`yg8j26` audits their content).

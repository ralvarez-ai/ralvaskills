---
title: 'Decide the bundle for backend-performance-architect'
priority: p2
section: ralvaskills
topics:
    - skills
    - owner-decision
---
`backend-performance-architect` 1.0.0 is published but belongs to no bundle, so `rsk install code-review` or any existing bundle leaves it out. Two shapes:

- **Join `code-review`**, which currently holds `security-reviewer`, `api-contract-reviewer` and `performance-reviewer` and is described as "cross-language reviewers". It would be the only non-reviewer there, and the description would need to widen.
- **Create a `backend` bundle** with `sql-architect`, `go-architect`, `backend-performance-architect` and `performance-reviewer`, so "make this service fast" installs as one unit.

Whichever wins, `internal/config/catalog.toml` changes and the bundle count and description budgets (`8k19zz`) still apply. The landing page counts bundles (`15`) in `gen-data.mjs` and `facts.ts`.

Done when the bundle is chosen, added, and the counts and description updated.

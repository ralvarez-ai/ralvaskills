---
title: 'Publish backend-performance-architect 1.0.0 and restore the site count gate'
priority: p1
section: ralvaskills
topics:
    - registry
    - skills
    - site
---
`backend-performance-architect` 1.0.0 is committed and pushed (`e3e812d`). The tree holds 61 skills while `registry/index.json` still lists 60 until the publish workflow commits the regenerated index, and `site/scripts/gen-data.mjs` asserts that the on-disk skills equal the index — so `task site:check` fails the count gate meanwhile, as it did for `mockup-architect` (`sykdrq`).

- `task registry:check` before the push: one new skill, `v1.0.0`, all version formats OK.
- `task tokens` regenerated `docs/TOKEN_COUNTS.md`: 61 skills, the new one at ~3413 body tokens.
- The `Publish Registry` run for `e3e812d` is the one to watch: it releases `backend-performance-architect@v1.0.0`, backs it up to R2 and commits the index.
- Bundle: deliberately left out of `code-review`, which holds reviewers only. Decide whether it joins that bundle or whether a new `backend` bundle groups `sql-architect`, `go-architect`, `backend-performance-architect` and `performance-reviewer`.
- The skill pairs with `performance-reviewer` (findings report) the way `observability-architect` pairs with `grafana-architect` (produce vs consume); if either side gains a section the other cites, keep the anchors real.

Done when the index lists 61 skills, `backend-performance-architect@v1.0.0` is published, `task site:check` passes on the updated `main`, and the bundle decision is recorded.

Resolved 2026-11-08: the `Publish Registry` run for `e3e812d` (42503851339) finished green, committed the regenerated index (`0d0cce7`) with 61 skills, and published `backend-performance-architect@v1.0.0` with its tarball. `gen-data.mjs`'s index assertion no longer fails on this skill, so the count gate is satisfied; the bundle decision moved to `k2v7dq`. Local `site:check` still stops earlier on an unrelated untracked `skills/frameworks/wails-architect/`, which is not part of this item.

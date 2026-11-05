---
title: 'Merge or push feat/web-ember-palette'
priority: p1
section: landing
topics:
    - palette
    - git
---
Three signed commits on top of `main` (the palette proposals, the Ember palette applied, the docs and the `MOCK` comparison option). Nothing under `skills/**`, `registry/index.json`, `cmd/` or `internal/` changes, so merging publishes nothing. The owner chooses how: a pull request, or a local merge commit with a detailed body as was done for the landing page. Pushing `main` runs `Check site`, which fails until the `NODE_AUTH_TOKEN` secret exists.

Done when the branch is in `main` and `git log main..feat/web-ember-palette` is empty.

Resolved 2026-11-04: merged into `main` with a signed merge commit (97da0aa). Not pushed yet.

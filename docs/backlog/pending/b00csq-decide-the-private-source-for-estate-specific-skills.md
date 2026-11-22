---
title: 'Decide the private source for estate-specific skills'
priority: p1
section: ralvaskills
topics:
    - skills
    - private
---
`docs/designs/2026-11-21-private-and-remote-skill-sources.md` says the author's personal skills were migrated to a private repository, but this review could not find its name or path. Decide the repo and the skill name prefix (`ralva-`), and register it as a source. Acceptance: path and name recorded here; `rsk` can install from it.

Source: `claude-design/docs/reviews/2026-11-21-landing-pages-estate-review.md` and the skills plan in this thread's answer (private skills source is an open decision).

Update 2026-11-22: the private source is `ralvarezdev/skills` (GitHub, private, branch `main`, skills under `skills/<domain>/...`). It holds `uru-paper-format`, `uru-thesis-defense-designer`, `uru-thesis-reviewer`, `uru-thesis-writer`, `go-cli-house`, `go-house-libs`, `go-lib-module`, `go-library-builder`, `go-estate-surfaces` and `pi-iteration-workflow`; the names carry no `ralva-` prefix. Those skills were dropped from `registry/index.json` the same day (index only; releases and history stay public). Still open: register it as a source and prove `rsk install` from it (needs `rsk` 3.10.0, `rsk source` does not exist in 3.9.1), and move `uns-platform`, whose lock names four of them.

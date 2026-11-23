---
title: Publish go-architect 1.21.0 and verify the registry
priority: p2
section: ralvaskills
topics:
    - go-architect
---
Bumped `go-architect` to 1.21.0 with the co-located `gen/` rule: generated contract code lives in `gen/` beside its consumer, mixed packages are split, generated packages are leaves, and sqlc stays with its context.

Commit `d977641` on `main`; `task lint` passed and `task registry:check` reported `v1.20.3 -> v1.21.0` ready. The push triggered Publish Registry run `45275217475`. Verify the run and the published index, then close this item.

**2026-11-23: release and index published.** `go-architect@v1.21.0` exists as a GitHub release (`2026-11-23T22:00:10Z`), and the bot committed `chore(registry): update index [skip ci]` as `64c007d`. The Pages deploy failed in both runs (`45275217475`, `45276049512`) with `No artifacts named "github-pages" were found`, while `Upload Pages artifact` reported success; that is a separate registry-publish bug, tracked in `5cwwfy`.

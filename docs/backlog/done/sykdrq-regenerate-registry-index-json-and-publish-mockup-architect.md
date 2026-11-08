---
title: 'Regenerate registry/index.json and publish mockup-architect 1.0.0'
priority: p1
section: ralvaskills
topics:
    - registry
    - skills
---
Until the index lists 60 skills, `task site:check` and the full site build fail the count gate. Regenerate with `task registry:generate` (ask first, per AGENTS.md), run `task lint` and the registry checks, then merge to publish. Design: `docs/designs/2026-11-05-family-chrome-and-mockup-bar.md`.

Update 2026-11-06: the index on `origin/main` lists 59 skills (it has `token-economy`, not `mockup-architect`), and the count is now 60. `mockup-architect` is committed, so the publish workflow adds it on the next push; `task registry:generate` is not needed. To run `task site:dev` before that, `task registry:sync-index` works locally; revert it with `git checkout registry/index.json` before pushing. The push itself is tracked in the release item.

Resolved 2026-11-08: `registry/index.json` on `main` lists 60 skills and includes `mockup-architect`, so the count gate this item tracked is satisfied. The next count move (60 → 61) is tracked in `b8p4rw`.

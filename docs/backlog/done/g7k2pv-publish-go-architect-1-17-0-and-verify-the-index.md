---
title: 'Publish go-architect 1.17.0 and verify the index'
priority: p2
section: ralvaskills
topics:
    - registry
    - skills
---
`go-architect` 1.17.0 adds the shared-package guidance: name a package for what it
provides, keep one tiny cross-context vocabulary package (`money`, `permission`
or the domain word) and split it when it grows unrelated clusters, and never name
another context's exported type in a port signature (enforce it with an AST test).
`SKILL.md`, `LAYOUTS.md` and `DIGEST.md` carry the rules; the frontmatter is
1.17.0.

- `task registry:check` before the push: `go-architect` v1.16.0 → v1.17.0, all
  version formats OK.
- `task tokens` regenerated `docs/TOKEN_COUNTS.md`; it was stale (61 skills, no
  `wails-architect`) and now says 62.
- The `Publish Registry` run for this push is the one to watch: it releases
  `go-architect@v1.17.0`, backs it up to R2 and commits the regenerated index.

Done when `registry/index.json` lists `go-architect` 1.17.0, the
`go-architect@v1.17.0` release exists, and `task registry:check` says there is
nothing left to publish on the updated `main`. Consumer projects pin the skill, so
the re-pin is tracked in their own backlog.

Resolved 2026-11-09: the `Publish Registry` run for `5a60c9f` (42642109810) finished green,
published `go-architect@v1.17.0` (15:11:44Z) and committed the regenerated index (`2e5825b`);
`registry/index.json` lists `latest: 1.17.0`; `task registry:check` has nothing left to publish for this
skill (it also reports the concurrent `remotion-architect` skill, which is not part of this release).
Consumer re-pins stay in their own backlogs.

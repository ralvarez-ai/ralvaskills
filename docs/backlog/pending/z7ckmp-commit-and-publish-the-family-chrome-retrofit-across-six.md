---
title: 'Commit and publish the family chrome retrofit across six repos'
priority: p1
section: ralvaskills
topics:
    - git
    - family
    - site
---
Only `claude-design`'s first commit is in. Working trees: ralvaskills (family credit, `mockup-architect` 1.0.0,
design doc, count 58 → 59), ralvafi (preview bar in `a-palettes.html` + generator), ralvapass (9 pages), ralvasend
(`palettes.html`, `views.html`), ralvadocket (1 page), ralvarez-brand-kit (`CHANGELOG.md`). ralvagym is clean: its
review and changelog are committed, and its local `main` is 32 commits ahead of origin.

Per-repo go-ahead is the owner's; nothing is committed or pushed without it. The trees hold work from earlier
sessions too, so stage by path, never `git add -A`: ralvafi already had `Taskfile.yml`,
`internal/financemcp/install.go` and `spikes/wails-desktop/build/android/overlay.json` modified; ralvapass,
ralvasend and ralvadocket carry unrelated site and mockup edits. Generated files travel with their generator
(`a-palettes.html`, ralvapass `index.html`, palettes `index.html`).

In this repo, regenerate `registry/index.json` first (`sykdrq`): a push to `main` publishes `skills/**`. Commits
are signed, author Ramón Álvarez, body explains why, no co-author trailer.

Done when each repo has one commit that `git status` shows clean for, and ralvaskills is pushed after the index
regeneration.

Update 2026-11-06: the ralvaskills tree is committed (family credit, `mockup-architect`, design doc, count 60, the digests and the compaction pass, one commit per concern, no trailer); only the push is left, tracked in the release item. The other five repos are unchanged.

---
title: 'Update the `go-house-libs` references to the `go.ralva.dev` paths and publish the skill'
priority: p2
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
    - skills
---
`skills/personal/go-house-libs/references/{MCPKIT,SVCKIT,TERMKIT}.md` still say `github.com/ralvarezdev/<kit>`. They were left out of the 2026-11-07 move on purpose: a `skills/**` edit publishes on push to `main`. Change the three paths to `go.ralva.dev/<kit>` (and the version lines to the tags in `changes/2026-11-07-go-module-host-go-ralva-dev.md`), bump the skill version in its frontmatter, then push. Ask before editing: it ships to every consumer that installs the skill.

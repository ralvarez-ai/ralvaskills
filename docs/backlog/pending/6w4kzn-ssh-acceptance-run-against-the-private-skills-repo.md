---
title: 'Run the SSH acceptance test against the private skills repository'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#9d3mhw
---
Every real-remote check so far used public HTTPS (`cursor/plugins`, `github/gitignore`, `anthropics/skills`). The SSH path, a private repository and a credential helper are untested, as is fetching a commit by id from a private host.

**Steps.** With the owner present: register `@ralvarezdev/skills` with a sandboxed `RSK_CONFIG_HOME`, trust it from a terminal, run `rsk source test`, `catalog`, a project install, `sync` after wiping the cache, and `update`. Needs the owner's SSH key and an explicit go-ahead; never run it unattended.

**Done when.** the run passes, or each failure is recorded as its own item. Blocks tagging the first CLI release.

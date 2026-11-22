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

Update 2026-11-21 (after step 1f): the run now also covers the implicit path. Add to the steps: `rsk install @ralvarezdev/skills/<name>` unregistered at a terminal (SSH first, the prompt shows both URLs), the same with `--trust @ralvarezdev/skills` and no terminal, `--save` then `rsk source trust`, and `rsk sync --trust` on a second sandbox from the lock. Note that two earlier 1f runs already used the owner's SSH key against a public repository (`vercel-labs/agent-skills`, read-only) before `ssh` was stubbed out: ssh finds the key through the passwd home, so an empty `$HOME` does not hide it. For a run that must not use the key, put a failing `ssh` first on `PATH`.

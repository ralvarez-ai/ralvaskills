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

Resolved 2026-11-21, with the owner present and after their go-ahead, in sandboxed `RSK_CONFIG_HOME`s, against `git@github.com:ralvarezdev/skills.git` (HEAD `cb88a6e`, ten skills): `source add`, `source test` refused before trust and reachable after it, the grant made on a pty; `catalog --source git`; a project `install --pin` that showed the digest, refused without `--yes` and wrote nothing, then pinned with it (lock v2 with url, commit, tree and path); `sync` after wiping the cache and the links, which fetched the locked commit by id over SSH from the private host; `update` (nothing to update, so drift detection on a private source is not exercised); an unregistered `@ralvarezdev/skills/<name>` approved on a pty, and the same with `--trust` and no terminal, neither writing `sources.json` or `trust.json`; `sync --trust` on a third machine from that lock. Failure paths with a failing `ssh` first on `PATH`: `catalog` warned per source and still listed the registry, a bare-name install failed naming the unreachable source and linked nothing, and a wrong implicit repository failed naming both the SSH and HTTPS attempts. The key was used only to read that repository (ls-remote and fetch); nothing was pushed. One finding, recorded as `frqzcr`: the runner honors the owner's `url.*.insteadOf`, so an HTTPS URL can travel over SSH.


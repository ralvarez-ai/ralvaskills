---
title: 'Refresh the consumers onto rsk 3.7.1 and repair their lock drift'
priority: p2
section: ralvaskills
topics:
    - rsk
    - skills
related:
    - k3vx7p
    - m4pnr7
    - 8k19zz
---
The digest audit and the trim published 31 skill releases, and the CLI is at `v3.7.1`, but the consumers are behind on both. Sixty-three repos still pin `github:ralvarez-ai/ralvaskills = "3.7.0"` in `mise.toml` (only this repository is on 3.7.1), and the 110 refreshed projects still hold the pre-trim digest text — harmless in rules, since no rule changed, but stale.

The lock drift is the part that matters. `ralvarez-go/corekit` was found with its lock at `go-architect` 1.20.1 while the link pointed at 1.20.2; `rsk update` before 3.7.1 called that current and never repaired it. 3.7.1 does repair it, so a bare `rsk update` per project both pulls the newer digests and settles any drift. The four skills whose descriptions `8k19zz` shortened (`adhd-mode` 1.0.2, `repo-tooling-architect` 1.6.3, `rsk-guide` 0.5.3, `token-economy` 1.1.2) come along with it.

One trap: `mise exec -- rsk` resolves the pin of the directory it runs in, so inside a consumer it still runs the old binary. Bump the pin first, or call the installed 3.7.1 by absolute path.

Done when every repo that pins rsk is on 3.7.1, a bare `rsk update` in each project reports no drift, the resulting `.rsk` changes are committed with explicit pathspecs (one `chore(rsk)` commit per repo), and the eleven repos whose commits cannot be pushed are recorded in `w8nq2s` rather than retried blindly.

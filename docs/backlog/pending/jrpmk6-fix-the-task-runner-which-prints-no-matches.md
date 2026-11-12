---
title: 'Fix the task runner, which prints No matches'
priority: p3
section: ralvaskills
topics:
    - tooling
---
`task --list` prints "No matches." in this environment, so the `web:*` and `ci:local` Taskfile entries were never run through Task; the underlying commands were run by hand. Find why (version, shim or Taskfile parsing), fix it, then run `task site:check` and `task ci:local` for real.

Reproduced 2026-11-11 while gating an `rsk update` change: `task lint` and `task --verbose lint` print only "No matches." and `task --list` prints the same, so the Go gate is affected too and had to be run as bare `golangci-lint run --config .golangci.yml ./...` (0 issues). The shim is `mise` (`/mnt/secondary/dev-cache/mise/shims/task`, Task 3.4.2).

Done when `task --list` lists the tasks and both run.

Root cause found 2026-11-12: `task` is not pinned in `mise.toml`, so the mise shim resolves the system `/usr/bin/task` 3.4.2, which cannot parse this Taskfile; `/usr/bin/task --list` fails the same way, so it is not a shim bug. `mise x task@3.53.1 -- task --list` lists every task correctly. The fix is one line in `mise.toml` (`task = "3.53.1"`, the newest installed), then run the gates as tasks again. Task 3.51–3.54 are already in the mise cache.

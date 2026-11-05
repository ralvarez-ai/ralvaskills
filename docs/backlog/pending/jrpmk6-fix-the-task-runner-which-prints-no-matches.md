---
title: 'Fix the task runner, which prints No matches'
priority: p3
section: ralvaskills
topics:
    - tooling
---
`task --list` prints "No matches." in this environment, so the `web:*` and `ci:local` Taskfile entries were never run through Task; the underlying commands were run by hand. Find why (version, shim or Taskfile parsing), fix it, then run `task site:check` and `task ci:local` for real.

Done when `task --list` lists the tasks and both run.

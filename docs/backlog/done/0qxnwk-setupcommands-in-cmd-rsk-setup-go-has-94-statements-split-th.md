---
title: '`setupCommands` in `cmd/rsk/setup.go` has 94 statements'
priority: p3
section: ralvaskills
topics:
    - ralvaskills
---

split the registration groups when touching the command tree

Resolved 2026-11-22: `setupCommands` now only registers the groups and calls seven registration functions (catalog, tools and sources, lifecycle, install, views, state, row actions) before `markDestructive` and `assignGroups`. Checked by diffing `rsk --help` and the help of every command and subcommand before and after: 1652 lines, identical.

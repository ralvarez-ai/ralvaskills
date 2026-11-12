---
title: 'Refresh the global skill installs'
priority: p3
section: ralvaskills
topics:
    - rsk
    - global
---
The 27 machine-wide skills for Claude Code and OpenCode are old (for example `docs-architect` 1.0.0, `go-architect` 1.10.0). Eight members of the `global` bundle are not installed: `code-design-refactor`, `grill-with-docs`, `token-economy`, `feature-planner`, `security-reviewer`, `design-patterns`, `rsk-guide` and `skill-builder`. Global installs were not touched in the 2026-11-06 pass. Related: `8k19zz` (the description budget is already over).

Done when the global installs are on current versions, the bundle members are installed or dropped from the bundle on purpose, and the description budget still holds.

2026-11-12: the four global links moved to `hexagonal-architecture@2.0.0` and `demo-script-designer@2.0.0`, so the rename half is done. `rsk` v3.7.0 (the alias-capable release) exists and `mise.toml` pins `github:ralvarez-ai/ralvaskills = "3.7.0"`, but it is not installed — every `mise` command warns `missing: github:ralvarez-ai/ralvaskills@3.7.0` and `rsk` still resolves to 3.6.0. `mise install` first, then `rsk list --global` to see which of the 27 are still old.

2026-11-12 (later): done and verified. `mise install` put 3.7.0 on disk and `mise exec -- rsk --version` reports `3.7.0 (rev 92e460f)`; a shell opened before the install still resolves 3.6.0 from its frozen PATH, so trust `mise which rsk` over `which rsk`. `rsk update --global` then moved all 26 installed skills off their 1.0.0/1.10.0 pins (`docs-architect` 1.2.1, `go-architect` 1.20.1). What is left is the other half of this item: the 8 `global` bundle members that are not installed (`code-design-refactor`, `grill-with-docs`, `token-economy`, `feature-planner`, `security-reviewer`, `design-patterns`, `rsk-guide`, `skill-builder`), and the description budget.

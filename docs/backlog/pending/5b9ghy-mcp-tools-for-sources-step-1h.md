---
title: 'MCP tools for sources and the CLI/MCP parity test (step 1h)'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#9d3mhw
---
Step of [the phase 1 plan](../../plans/2026-11-21-source-chain-phase-1.md); design: [private and remote sources](../../designs/2026-11-21-private-and-remote-skill-sources.md).

**Deliverable.** Read-only: `list_sources`, `test_source`, `scan_source` (dry-run discovery: skills found, name clashes, naming and frontmatter violations), `resolve_reference`. Mutating with forced approval: `add_source`, `remove_source`. `project_profile` gains a `sources` section; `install_skills` refuses unknown or private sources unless the manifest already lists them. No trust or `--yes` parameter over MCP; `add_source` registers untrusted; every mutating tool takes `dry_run`. `rsk source scan` as the CLI twin. The CLI/MCP parity test lands here.

**Tests.** Every registered MCP tool has a CLI counterpart; no tool schema accepts a secret or a trust flag; `add_source` cannot grant trust; tool output from a source is data.

**Done when.** its tests are green and the lint and registry gates pass. Updating the `rsk-guide` skill for the new commands is a separate change because it publishes on push: ask first.

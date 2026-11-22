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

Update 2026-11-21 (after step 1f), what 1h can reuse and must keep:
- `internal/resolve.Reference(ctx, ref, chain, Options{Bundles, ImplicitHint})` is the core of `resolve_reference`: plain structs with json tags, no printing. `Chain.Survey` asks every source and keeps each answer. Parity: the CLI is `rsk resolve`.
- `source.Path.Skipped` and the git index's skipped list are the raw material for `scan_source` and `rsk source scan` (name clashes already fail as `ErrDuplicate`; the skipped reasons are the symlink, bad-name, oversized and missing-SKILL.md cases). `cursor/plugins` is the real clash to test against.
- The MCP chain never sets `ChainOptions.Implicit`, and `TestInstallSkillsRefusesAnUnregisteredSourceAndNeverContactsIt` pins that `install_skills` refuses an unregistered source. Keep it, and extend it to private sources when 1g lands.
- Open: the MCP pin (`install_skills` with `pin`) shows no digest text (trust rule 7 is enforced in the CLI only); the tool's forced approval does not show it either. Decide here whether the tool returns the text and refuses to pin until the user confirms it.
- `project_profile` `sources` section: path sources show trust "not required"; an implicit source is never listed (nothing is registered).
- No trust or `--yes` parameter; `add_source` must not take `--trust`, and must never write `trust.json` (`--save` in the CLI already follows this).

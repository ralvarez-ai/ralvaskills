---
title: 'private flag and the local overlay (step 1g)'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#9d3mhw
---
Step of [the phase 1 plan](../../plans/2026-11-21-source-chain-phase-1.md); design: [private and remote sources](../../designs/2026-11-21-private-and-remote-skill-sources.md).

**Deliverable.** Honor `Private` (currently refused by `source.New` and `rsk source add --private`): excluded from `catalog`, MCP resources, `rsk://catalog` and `search_skills` unless `--include private`. Pins from a private source go to gitignored `.rsk/rsk.local.mod`, `rsk.local.lock` and `SKILLS.local.md`, merged with the committed files on read, appended to `.git/info/exclude`; `--shared` opts in to the committed files. Adapters: Claude Code reads the local digest through `CLAUDE.local.md`; OpenCode and Codex need an equivalent that is not committed, and whether each client honors an uncommitted import is unverified and is the first thing to test.

**Boundary.** This step touches `internal/tool/*` adapter config-writing: ask before changing it (AGENTS.md).

**Tests.** A private pin never reaches a committed file; both layers merge on read; `--shared` writes the committed ones; a teammate without the local files sees a clean project; per-client check of the uncommitted import.

**Done when.** its tests are green, the lint and registry gates pass, and the three adapters' behaviour is recorded.

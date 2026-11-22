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

Update 2026-11-21 (after step 1f): `--save` registers a source untrusted and `rsk install --pin` of a non-default source now shows the digest first, so the private overlay has two inputs to respect: a pin from a private source must skip the committed `SKILLS.md` even when `--yes` accepted its text, and `rsk pin` (which has no install step) must route by the skill's origin. Path-source skills (`skill.SourcePath`) never reach a lock, so only git sources need the overlay. An implicit install of a private repository (`@owner/repo/name` over SSH) currently pins into the committed files; decide whether implicit sources default to the overlay or ask.

Update 2026-11-22: the owner approved changing `internal/tool/*` for this step ("go ahead with the 4 items"). Not started: the first thing it must do is test, in the real clients, whether Claude Code (`CLAUDE.local.md`), OpenCode and Codex honor an uncommitted import, which cannot be done from the unit tests. Where a client cannot be verified, refuse to pin a private skill for it instead of writing its digest to a committed file. Start it in its own session; `internal/digest` now reads a versionless skill's description, and `internal/sourceops` and `internal/resolve` already carry the source views a `private` flag would filter.

Resolved 2026-11-22 (`0de382e` the private flag, `d8c717f` the manifest layers, `f2d3ed9` the local digest and the Claude Code and OpenCode imports, `cdf1bf9` install, pin and `--shared`, plus the follow-up fixes). The decisions and what each client did are in the plan, section 1g. Not covered: an implicit source is never private.


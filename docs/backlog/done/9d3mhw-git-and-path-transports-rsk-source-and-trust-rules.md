---
title: 'Git and path transports, rsk source and the trust rules (phase 1)'
priority: p1
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 1 of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

Plan: [phase 1 steps 1a to 1h](../../plans/2026-11-21-source-chain-phase-1.md).

Progress 2026-11-21: steps 1a to 1e are built and committed, not pushed (`4781793` hardened git runner, `3696b4c` `rsk source`, `ec6ef87` git transport for declared sources, `ac9ffd4` lock, tag pins and update for git skills, `2fdb002` lock by commit, `997870e` official source on the git transport). Step 1f is built and committed too (2026-11-21, `4513836` declared path sources, `198a35d` `rsk resolve`, `b03eaea` implicit sources in `internal/sources`, `61a746f` install from an unregistered source and `sync --trust`, `5dd1413` digest confirmation, `cc2fc81` `--from`; see `3k7vqe` in done and the plan, section 1f). Follow-ups found on the way: `g1zk6n`. The SSH acceptance run passed on 2026-11-21 (`6w4kzn` in done), so the first CLI release is unblocked. Step 1h is built too (2026-11-22, `5b9ghy` in done). Remaining, one backlog item each: 1g `private` and the local overlay (touches `internal/tool/*`, needs approval); 1h the MCP tools. No CLI release is tagged before 1f.

**Deliverable.** Extend `internal/git` (fetch, archive, hardening, redaction); git and path transports; `--from`; `rsk source`, `rsk resolve`; implicit sources with the trust prompt; the local overlay (5.7) and its adapter checks; read-only MCP tools (`list_sources`, `test_source`, `scan_source`, `resolve_reference`) and `add_source` / `remove_source`.

**Tests.** Bare-repo fixture, private repo over SSH, offline cache hit, URL allowlist, `roots` traversal.

**Risk.** A new cache root.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects. The trust rules of design section 3 ship in this same release, not after it. Depends on phase 0b.

Update 2026-11-22: steps 1a to 1f and 1h are done and the SSH acceptance run passed; the versionless digest fix and the `AGENTS.md` update are in. `v3.10.0` is tagged locally at `869bedc` and not pushed (`c3qwq1`). Remaining: 1g (`8d2nwx`, approved, own session), the `rsk-guide` 0.6.0 update (`7bxgw3`, blocked on a permission rule), the consumer roll-out (`9m3rtc`, after the release) and `g1zk6n`.

Update 2026-11-22 (later): step 1g is built too (`8d2nwx` in done): `--private`, the local layers, `SKILLS.local.md`, the Claude Code and OpenCode imports (verified in the real clients) and `--shared`. All of phase 1 is built. Open: the push of these commits, the `rsk-guide` update, the consumer roll-out and the follow-ups in `g1zk6n`.

Resolved 2026-11-22: phase 1 is built, verified and released as `v3.10.0` (tag at `2c1cc10`, six assets, checksum and `--version` checked). Every exit criterion of the plan holds, including the SSH acceptance run, and the plan is marked done. What is left is its own work: the consumer roll-out (`9m3rtc`), the `rsk-guide` 0.6.0 update (`7bxgw3`, blocked on a permission rule) and the follow-ups in `g1zk6n`.


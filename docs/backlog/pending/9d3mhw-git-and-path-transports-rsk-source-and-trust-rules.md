---
title: 'Git and path transports, rsk source and the trust rules (phase 1)'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 1 of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

**Deliverable.** Extend `internal/git` (fetch, archive, hardening, redaction); git and path transports; `--from`; `rsk source`, `rsk resolve`; implicit sources with the trust prompt; the local overlay (5.7) and its adapter checks; read-only MCP tools (`list_sources`, `test_source`, `scan_source`, `resolve_reference`) and `add_source` / `remove_source`.

**Tests.** Bare-repo fixture, private repo over SSH, offline cache hit, URL allowlist, `roots` traversal.

**Risk.** A new cache root.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects. The trust rules of design section 3 ship in this same release, not after it. Depends on phase 0b.

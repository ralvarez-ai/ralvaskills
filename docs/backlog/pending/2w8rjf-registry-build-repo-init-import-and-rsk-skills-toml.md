---
title: 'Registry build, repo init/import, rsk-skills.toml (phase 3)'
priority: p3
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 3 of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

**Deliverable.** `rsk registry build` and `deprecate`; `rsk repo init` and `import` with their MCP tools; the MCP/CLI parity test; `rsk-skills.toml`; docs; a migration guide; replace path-based `personal`.

**Tests.** Golden index, nested-layout fixture, duplicate-name failure.

**Risk.** A publishing workflow rewrite.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects. Depends on phase 1.

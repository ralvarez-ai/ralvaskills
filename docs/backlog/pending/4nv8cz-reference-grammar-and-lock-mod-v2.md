---
title: 'Reference grammar and lock/mod v2 (phase 0b)'
priority: p1
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 0b of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

**Deliverable.** `@owner/repo/name` grammar; qualified keys; `mod` and `lock` v2 written only when needed; `schema.Check` accepts a range.

**Tests.** v1 to v2 migration, old-binary error message, `ParseNameVersion` with `@`.

**Risk.** The schema change touches every command.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects. Depends on phase 0a.

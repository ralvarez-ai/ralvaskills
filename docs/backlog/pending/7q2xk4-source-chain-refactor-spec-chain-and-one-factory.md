---
title: 'Source chain refactor: Spec, Chain and one factory (phase 0a)'
priority: p1
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 0a of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

**Deliverable.** `Spec`, `Chain` and a factory collapsing the call sites that construct sources; no schema change.

**Tests.** Resolution order, ambiguity, fail-closed, behaviour identical to today.

**Risk.** The refactor touches about 40 `Source` references.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects.

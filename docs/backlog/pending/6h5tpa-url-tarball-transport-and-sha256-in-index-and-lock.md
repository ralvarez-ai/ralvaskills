---
title: 'URL and tarball transport, SHA-256 in index and lock (phase 2)'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
---
Phase 2 of [private and remote skill sources](../../designs/2026-11-21-private-and-remote-skill-sources.md) (section 10), decided in [ADR 0003](../../adr/0003-source-chain-and-reference-grammar.md).

**Deliverable.** Public URL and tarball transport; SHA-256 in the index and the lock, plus the backfill of existing versions.

**Tests.** Digest mismatch aborts, backfilled digests match the release assets, tarballs rebuild deterministically.

**Risk.** The backfill touches every published version.

**Done when.** its tests are green, `task lint` and `task registry:check` pass, v1 files still round-trip, and the behaviour of the default registry is byte-identical for existing projects. The SHA-256 and backfill half does not depend on the chain and may ship earlier.

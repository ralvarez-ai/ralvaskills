---
title: 'Extract the shared landing tooling once a second landing exists'
priority: p3
section: landing
topics:
    - reuse
    - tooling
---
The gate (`verify.mjs`), the comparison tools, `CodeRegion`, the film engine and the Biome config will repeat across landing pages. Do not extract speculatively: when a second landing is built, list what is really duplicated and propose a shared package next to `@ralvarez-ts/tokens`; the owner decides. If a landing lives next to an application, audit which application components are portable first (no providers, router or state; no global CSS).

Done when the duplication is measured and the owner has accepted or rejected the extraction.

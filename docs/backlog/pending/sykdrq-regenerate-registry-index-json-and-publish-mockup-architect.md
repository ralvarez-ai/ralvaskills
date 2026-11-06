---
title: 'Regenerate registry/index.json and publish mockup-architect 1.0.0'
priority: p1
section: ralvaskills
topics:
    - registry
    - skills
---
Until the index lists 59 skills, `task site:check` and the full site build fail the count gate. Regenerate with `task registry:generate` (ask first, per AGENTS.md), run `task lint` and the registry checks, then merge to publish. Design: `docs/designs/2026-11-05-family-chrome-and-mockup-bar.md`.

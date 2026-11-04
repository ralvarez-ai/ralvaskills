---
status: done
project: ralvaskills
date: 2026-11-04
related: [designs/2026-11-04-ralvaskills-landing-page.md]
---

# The landing page now matches the mockup in use, and CI checks it

## Why
The `web/` build had been declared equal to the combined mockup from still
screenshots at a loose pixel threshold. In use it was not: both films never
played, the section rail never showed, the header nav never marked the current
section, and one section was tinted. None of it could fail a gate that only
looked at stills.

## Impact and migration
- `web/` only; nothing under `skills/**`, `registry/index.json`, `cmd/` or
  `internal/` changed, so nothing publishes.
- `verify.mjs` is stricter (films, rail, nav) and waits for films before axe.
- New workflow `.github/workflows/check-site.yml` runs the build and the gate on
  pull requests touching `web/`. It needs a `NODE_AUTH_TOKEN` repository secret
  with `packages:read` to install `@ralvarez-ts/tokens`; until it exists the
  install step fails. The workflow has not been run on GitHub.
- Detail and the corrected claim in design doc round 7.

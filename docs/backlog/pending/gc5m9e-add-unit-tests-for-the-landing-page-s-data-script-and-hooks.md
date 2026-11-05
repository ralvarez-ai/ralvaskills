---
title: 'Add unit tests for the landing page''s data script and hooks'
priority: p3
section: landing
topics:
    - tests
---
The gate is Playwright end to end; there are no unit tests. The pieces worth testing directly are `gen-data.mjs` (the drift check and the parsing of the Go sources), the film engine and the reveal and count hooks. Choose a runner that fits Vite (Vitest) and keep it out of the build.

Done when those three have tests that fail when their behaviour is broken, wired into `pnpm check` and the workflow.

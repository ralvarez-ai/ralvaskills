---
title: 'Decide whether web/ is renamed to site/'
priority: p2
section: landing
topics:
    - structure
    - owner-decision
---
In other projects the landing lives in `site/` and the application stays in `web/`; here `web/` is only the landing, so it would be renamed for consistency. A rename touches the Taskfile (`web:*`, `ci:local`), `check-site.yml` and its path filters, AGENTS.md, README, `biome.json`, `.gitignore`, the comparison scripts' relative paths to the mockups, the docs and the act setup. Do it as one commit and re-run the gate and `act` afterwards.

Done when the owner has decided; if renamed, `grep -rn "web/"` finds nothing stale and `verify.mjs` and `act` pass.

Resolved 2026-11-04: renamed with `git mv` in d08a248. Taskfile tasks are now `site:*`; `check-site.yml`, AGENTS.md and the working directory follow. `pnpm lint`, `pnpm build` and `verify.mjs` pass from `site/`; `act` was not re-run, so the first CI run after the push is the check.

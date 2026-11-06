---
title: 'Add the NODE_AUTH_TOKEN secret so the Check site workflow can pass'
priority: p1
section: landing
topics:
    - ci
    - github
---
`.github/workflows/check-site.yml` fails at Install with `401 Unauthorized` from `npm.pkg.github.com`: `@ralvarez-ts/tokens` is on GitHub Packages and the repository has no `NODE_AUTH_TOKEN` secret with `packages:read`. It is an outward-facing change that needs a token chosen by the owner, so it is not done by the agent.

Done when `Check site` is green on a pull request and on `main`. That run is also the first time the hosted-Chrome branch of the Browser step is exercised (it has only been run under `act`, where Chromium is installed instead).

Checked 2026-11-04: the push of `main` (97d0cb2, the first run with the `site/**` paths) failed the same way as 30a9ee4, so the rename itself is still unproven in CI.

Checked 2026-11-04 after the final push: the last `Check site` run is still the failed one on 97d0cb2. Later pushes (42628d7, CHANGELOG only) do not match its path filters, so it will not run again until something under `site/**` changes or it is started by hand with `workflow_dispatch`.

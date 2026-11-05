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

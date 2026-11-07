---
title: 'Deploy the landing page to Cloudflare Workers at rsk.ralvarez.dev'
priority: p1
section: landing
topics:
    - deploy
    - cloudflare
blocked_by:
    - 6azjzw
---
Decided 2026-11-06 (supersedes the Pages decision below): a Cloudflare Worker `ralvaskills` with static assets (`site/wrangler.jsonc`), deployed with `wrangler deploy`, separate from any application and not embedded in one. The registry stays on GitHub Pages at skills.ralvarez.dev (`internal/config/fs.go` hardcodes that URL).

To do: create the project and the `rsk` CNAME; secrets `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` (least privilege, only for this workflow); `deploy-site.yml` (install, lint, build, gate, `wrangler pages deploy dist`) with path filters for the landing only; `public/_headers`; an `og:image`. `index.html` already points canonical and OG at https://rsk.ralvarez.dev/. Prove the workflow with `act` first. Deploying is outward-facing: confirm with the owner before running it.

Done when the URL serves the build, `verify.mjs` passes against the real URL, and a registry publish and a site deploy do not trigger each other.

Progress 2026-11-06: in the repo, `site/wrangler.jsonc` (static assets, custom domain `rsk.ralvarez.dev`), `site/.env.example`, `site/public/_headers`, the `og:image` card (`site/public/og.png`, drawn by `site/scripts/og.mjs`), the Teal & Sand `theme-color` metas, `task site:deploy` (gate, then `wrangler deploy`) and `.github/workflows/deploy-site.yml` (manual dispatch, environment `cloudflare`) exist; `wrangler deploy --dry-run` reads `dist/`. Not done, all outward-facing and the owner's: a token with Account > Workers Scripts: Edit and Zone > Workers Routes: Edit on ralvarez.dev; `site/.env` or the secrets `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID` and `NODE_AUTH_TOKEN`; then the first deploy (the custom domain is created by it, no Pages project and no CNAME step), `verify.mjs` against the real URL and a check that the registry workflow and this one do not trigger each other. The workflow has not been proven with `act`. The sibling sites follow this pattern: `ralvarez.dev`, `kaucrow`, `ralvafi`, `ralvapass`, `ralvasend`, `ralvagym`, `ralvadocket`.

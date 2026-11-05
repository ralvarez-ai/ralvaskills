---
title: 'Deploy the landing page to Cloudflare Pages at rsk.ralvarez.dev'
priority: p1
section: landing
topics:
    - deploy
    - cloudflare
blocked_by:
    - 6azjzw
---
Decided: Cloudflare Pages project `rsk`, Direct Upload, separate from any application and not embedded in one. The registry stays on GitHub Pages at skills.ralvarez.dev (`internal/config/fs.go` hardcodes that URL).

To do: create the project and the `rsk` CNAME; secrets `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` (least privilege, only for this workflow); `deploy-site.yml` (install, lint, build, gate, `wrangler pages deploy dist`) with path filters for the landing only; `public/_headers`; an `og:image`. `index.html` already points canonical and OG at https://rsk.ralvarez.dev/. Prove the workflow with `act` first. Deploying is outward-facing: confirm with the owner before running it.

Done when the URL serves the build, `verify.mjs` passes against the real URL, and a registry publish and a site deploy do not trigger each other.

---
title: 'Landing deploy follow-ups'
priority: p2
section: landing
topics:
    - deploy
    - cloudflare
    - cleanup
---
Left after the landings went live on Workers (2026-11-07):

- **Gate against the real URL:** `verify.mjs` was only run against the local `dist`, never against the live hosts. Run it (or a smoke check) against `rsk`, `fi`, `pass`, `send`, `gym`, `docket` on `.ralva.dev`.
- **ralvapass dirty files:** `Taskfile.yml`, `.github/workflows/deploy-site.yml`, `check-site.yml`, `CHANGELOG.md`, `site/README.md`, `site/public/_headers` and `site/scripts/headers.mjs` hold the owner's changes mixed with an agent's; commit them hunk by hunk. Its `site:deploy` edit (SITE_HOST, `deploy.mjs`) is uncommitted.
- **ralvasend:** `CHANGELOG.md` has an uncommitted "Landing hostname" line next to owner changes.
- **Headers:** `fi` and `gym` send 2 of the 3 basic security headers, `docket` none and has no `_headers` nor social card; check what they are meant to send.
- **Pages leftovers:** remove the custom domains and, once happy, the Pages projects `kaucrow-dev` and `ralvarez-dev`; enable Web Analytics for `kaucrow.dev` and `ralvarez.dev` if used (Pages injected its beacon, Workers does not). Their READMEs still describe the migration as pending.
- **Workflows:** to deploy from GitHub, set the repository variable `SITE_HOST` (a variable, not a secret) and the secrets `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID`, `NODE_AUTH_TOKEN` per repo; the `deploy-site.yml` files read the token, not `CLOUDFLARE_API_TOKEN_FILE`. Not proven with `act` (this repo's `task ci:local`).
- **Tokens:** `ralva-dev-workers`, `ralvarez-dev-workers`, `kaucrow-dev-workers` (in `~/Secrets/cloudflare/workers`, mode 0600) still carry Pages permissions from before; trim to Workers Scripts, Workers Routes and DNS on their zone. The DDNS token can drop `ralvarez.dev`.
- **Hard-coded app hosts:** done for ralvagym on 2026-11-07 (`VITE_APP_ORIGIN`, required by the build, `verify.mjs` reads it). Left: `ralvafi` `spikes/wails-desktop/main.go:23` defaults to `https://fi.ralvarez.app`, and its `web/index.html` has the host in `og:url` and `og:image` (see `homelab#hhndfy`).
- **Cleanup in `ralvarez.dev`:** the profile path fix touched eight scripts; its own `docs/` and skills still say `profile.md` in prose.

- **Registry host on the page:** `VITE_REGISTRY_HOST` (CI: the repository variable `REGISTRY_HOST`, now `skills.ralva.dev`) is required by the rsk build. Set the same variable wherever the landing is built outside this machine.
- **Local commits not pushed:** `1798c8d` (the refreshed `rsk.lock`) is committed here and not pushed; v3.6.0 and everything before it is on `origin`.

Done when each bullet is closed or consciously dropped.

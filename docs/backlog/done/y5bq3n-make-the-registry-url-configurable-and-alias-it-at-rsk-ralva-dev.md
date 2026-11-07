---
title: 'Make the registry URL configurable and alias it at rsk.ralva.dev/registry'
priority: p2
section: ralvaskills
topics:
    - registry
    - config
    - cloudflare
---
`DefaultRegistryURL = "https://skills.ralvarez.dev"` is a literal in `internal/config/fs.go`, compiled into every released `rsk`. `skills.ralvarez.dev` stays for ever (installed CLIs depend on it), but the URL must not be hard-coded: add an override order `RSK_REGISTRY_URL`, then `registry_url` in `~/.config/rsk/config.json`, then the default; `rsk init` already asks for it.

A nested alias on the landing host is viable: `source.NewRegistry` fetches `<baseURL>/index.json` and the index carries absolute `archive_url`s, so only `index.json` (82 KB) is needed at `https://rsk.ralva.dev/registry`. Decide how to serve it: a Worker route that proxies `/registry/*` to `skills.ralvarez.dev` (GitHub Pages stays the single publisher, no CI change, a few lines of Worker code and `run_worker_first` for that path), or copy `registry/index.json` into `site/dist/` and redeploy after every publish (no Worker code, but `publish-registry.yml` needs a deploy step and Cloudflare secrets). GitHub Pages cannot serve two custom domains, so the alias cannot be a second CNAME.

Done when `RSK_REGISTRY_URL=https://rsk.ralva.dev/registry rsk update` works, the default can change without touching more than the config, and a test covers the override order. Publishing a CLI release is outward-facing: confirm first. See `docs/plans/2026-11-06-ralva-domains-migration.md` in the homelab repo.

Closed 2026-11-07, with a different outcome from the alias: `RSK_REGISTRY_URL` (then `registry_url`, then the default) is in `Config.EffectiveRegistryURL`, with a test, and shipped in v3.6.0. The registry itself moved instead: GitHub Pages serves it at `skills.ralva.dev`, `skills.ralvarez.dev` is a 301 Redirect Rule to it, and the compiled default is the new host. The `/registry` proxy on the landing was built, deployed and removed again because it put traffic through a Worker for nothing. Checked: an installed 3.2.2 binary and a v3.6.0 build both installed a skill through the new and the old host.

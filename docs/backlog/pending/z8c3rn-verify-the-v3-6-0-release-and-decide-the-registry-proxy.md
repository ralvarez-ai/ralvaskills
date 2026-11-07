---
title: 'Verify the v3.6.0 release and decide how skills.ralva.dev is proxied'
priority: p3
section: ralvaskills
topics:
    - release
    - registry
---
v3.6.0 is published (6 assets, workflow green) but nobody has installed those assets: download the binary for this machine from the GitHub release, run it with an empty `RSK_CONFIG_HOME` against `skills.ralva.dev` (no `RSK_REGISTRY_URL`, to check the new compiled default) and install one skill. A build from the tag was already checked, the released binary was not.

Also decide the DNS mode of `skills.ralva.dev`: today the CNAME to `ralvarez-ai.github.io` is proxied by Cloudflare (it resolves to Cloudflare addresses) and GitHub shows no certificate of its own, while HTTPS works through Cloudflare. If a grey cloud is preferred, GitHub must issue its certificate first, which can take up to an hour; the redirect from `skills.ralvarez.dev` is unaffected either way.

Done means: the released binary installs a skill with no override, and the DNS mode is recorded in `docs/guides/registry-flow.md`. Local commits in this repo (`1798c8d` and the backlog commits after v3.6.0) are not pushed yet.

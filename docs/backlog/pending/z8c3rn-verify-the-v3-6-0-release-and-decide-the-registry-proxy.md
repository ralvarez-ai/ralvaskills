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

2026-11-12: superseded by `v3.7.0`, which carries the same work plus aliases and the budget warning; v3.6.1 was tagged but never released as a binary. The released `linux-amd64` binary was run on `ralvarez-go/corekit`: `rsk update docs-architect` applied 1.2.0 → 1.2.1 and re-rendered `.rsk/SKILLS.md` on its own, and the budget warning printed the three largest digests. One blocker remains for the "new compiled default" half: `~/.config/rsk/config.json` still sets `registry_url` to `https://skills.ralvarez.dev`, which outranks `config.DefaultRegistryURL`, so no run yet exercises `skills.ralva.dev` — drop or migrate that key before claiming the default is live. The DNS mode of `skills.ralva.dev` is still undecided.

2026-11-12 (later): `v3.7.1` followed the same day with the lock repair, and its released `linux-amd64` binary was exercised on `ralvarez-go/corekit` — it repaired a drifted lock and re-rendered `.rsk/SKILLS.md`. The `registry_url` override and the DNS decision are still open.

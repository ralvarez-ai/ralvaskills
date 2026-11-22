---
title: 'Push main and v3.10.0, then run the Release CLI workflow'
priority: p1
section: ralvaskills
topics:
    - rsk
    - release
related:
    - ralvaskills#9d3mhw
    - ralvaskills#9m3rtc
---
Owner-only: a hook in the agent session denies `git push *`, so the agent could not do it. State on 2026-11-22: `main` is about 30 commits ahead of `origin/main`, and a signed tag `v3.10.0` exists locally at `869bedc` (the changelog section is there). Every gate passed on that commit: `task lint`, `lint:doccheck`, `registry:check`, `go test ./...` and `task site:check`.

**Steps.** `git push origin main`, `git push origin v3.10.0`, then run the Release CLI workflow from the Actions tab with tag `v3.10.0` (it is manual-dispatch only). Pushing runs `publish-registry.yml` because `internal/**` changed; nothing under `skills/**` changed since the last push, so it should publish no skill: check its log says so. The site deploy is also manual (`deploy-site.yml`): the landing page now says nine tools, so deploy it after the release if you want the page to match.

**Done when.** the release assets exist and `rsk --version` from the released binary prints v3.10.0. Then start `9m3rtc` (the consumers' `mise.toml` pins must move before any project that uses official skills runs `rsk update`).

Resolved 2026-11-22: `main` and `v3.10.0` pushed (the tag at `2c1cc10`, after 1g), Release CLI run `44988887705` succeeded with six assets, the linux-amd64 checksum verifies and `rsk --version` prints `rsk 3.10.0 (commit 2c1cc10)`. The registry workflow published no skill. The site deploy is still manual. The consumer roll-out (`9m3rtc`) starts from here.

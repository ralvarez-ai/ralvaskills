---
title: 'Decide the module host of the other repositories under `ralvarezdev`'
priority: p3
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
Apps now use `go.ralva.dev/<name>`. Repositories still under the `ralvarezdev` account keep `github.com/ralvarezdev/...` or `github.com/planttalk/...`: `uns-platform` (20+ submodules), `m2p`, `strucgo`, `strucgo-runtime`, `klevor` and the `planttalk` ones. Decide per repository: move to `go.ralva.dev` (then add it to `MODULES` and `OWNER_OVERRIDES` in `go-vanity`), or leave as it is. Do not rename `uns-platform` before its consumers of the libraries are migrated.

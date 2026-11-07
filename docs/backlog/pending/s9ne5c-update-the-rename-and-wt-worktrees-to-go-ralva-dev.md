---
title: 'Update the `.rename` and `.wt` worktrees to `go.ralva.dev`'
priority: p3
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
`ralvapass.rename`, `ralvasend.rename` and `ralvadocket.wt/{generic-decisions,generic-ingest,split-backlog}` still require library tags that declare `go.ralvarez.dev`. Rebase or merge each onto its main branch once the app commits are pushed; otherwise repeat the `require` bump and `go mod tidy`.

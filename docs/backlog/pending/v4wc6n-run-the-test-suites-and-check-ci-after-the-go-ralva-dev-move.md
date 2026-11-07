---
title: 'Run the test suites and check CI after the `go.ralva.dev` move'
priority: p3
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
The move was verified with `go build ./...` and `go vet ./...` only (plus `go test ./...` in `ralvaskills`). Run `go test ./...` in every library and app, then check that each repository's CI resolves the private `go.ralva.dev/*` modules: `GOPRIVATE`, `GONOSUMDB` and git credentials in the workflows and in `mise.toml`/`Taskfile.yml` still name `go.ralvarez.dev` in places the text sweep did not reach (secrets and repository variables).

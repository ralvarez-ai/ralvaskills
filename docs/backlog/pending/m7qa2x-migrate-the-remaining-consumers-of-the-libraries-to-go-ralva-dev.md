---
title: 'Migrate the remaining consumers of the Go libraries from `github.com/ralvarezdev/*` to `go.ralva.dev`'
priority: p2
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
The libraries and nine apps moved on 2026-11-07, but `ralvaskills` and `autoreports` followed the same day (build and vet pass, uncommitted). These modules still require the libraries by their old GitHub path (older tags, still resolvable): `uns-platform` (about 25 modules: `svckit`, `corekit`, `identity`, `pgkit`, `grpckit`) and `vtitan` (`teamvoltimor`: `corekit`, `grpckit`, `resilience`, `svckit`, `termkit`, `tick`). Per module: point the `require` at the first `go.ralva.dev` tag of each library, rewrite the imports, `go mod tidy`, `go build ./... && go vet ./...`. Tags are in the release table of the 2026-11-07 change note. `uns-platform` and `vtitan` were left out on purpose: do them only when asked.

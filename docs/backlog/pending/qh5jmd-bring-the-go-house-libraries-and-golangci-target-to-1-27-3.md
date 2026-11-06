---
title: 'Bring the Go house libraries and golangci target to 1.27.3'
priority: p3
section: ralvaskills
topics:
    - go
    - skills
---
The Go skills and this repo now target 1.27.3, but `go-house-libs` still records the Go version each released module declares (1.27.1, and 1.27.2 for `contracts` and `connectrpckit`) and its "Go policy" line says 1.27.1. Move the modules, then update `STACK.md` and the policy line. Also confirm on 1.27.3 what `go-library-builder` ships: `assets/mise.toml`, `assets/Taskfile.yml`, `renovate.json` and the golangci-lint target its `STACK.md` now states, since nobody ran them on 1.27.3.

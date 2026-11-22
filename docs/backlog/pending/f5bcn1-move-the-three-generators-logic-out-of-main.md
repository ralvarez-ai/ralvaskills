---
title: 'move the three generators'' logic out of main'
priority: p3
section: ralvaskills
topics:
    - ralvaskills
---
`cmd/generate-registry/main.go` 383 lines, `cmd/count-tokens` 368 and `cmd/registry-remove` 294 each hold their whole
implementation in `main` — registry rendering, token counting, file removal — so none of it can be tested or reused. Move
the bodies into `internal/<tool>` packages and leave `main` a flag-parsing shim.

Module path `go.ralva.dev/ralvaskills/v3` is idiomatic (major-version suffix on a released CLI).

Note: this item was written by hand because `l9c2xm` carried an `l`, which is outside the tool's alphabet and made `rdk`
refuse the whole backlog; it is `k9c2xm` now and `rdk check` passes.

Update 2026-11-22: `count-tokens` is done (`bgmb76`, `internal/tokensreport`). `generate-registry` (436 lines) and `registry-remove` (296) remain; `generate-registry` feeds the published index, so check its move with `task registry:check` byte-identical.

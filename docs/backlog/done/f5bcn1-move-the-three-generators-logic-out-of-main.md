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

Update 2026-11-22: `registry-remove` is done: `cmd/registry-remove` is a flag shim over `internal/registryremove` (`Run(Options)`), now with unit tests (selection, validation, dry run, index-only removal). Only `generate-registry` (436 lines) remains here.

Note 2026-11-22: `generate-registry` is the publishing path (`publish-registry.yml` runs it; its path filter already includes `internal/**`, so a move needs no workflow edit). Its `main` and `runCheck` call `os.Exit` on many paths, which must become returned errors, and the move is verified by diffing `task registry:check` and an `--output-dir` run byte for byte before and after.

Resolved 2026-11-22: all three generators are shims now. `generate-registry` moved to `internal/registrygen` (`Run(Options)`, the `os.Exit` calls became returned errors, the shim prints `error: <err>` and exits 1; a few messages lost their `error ` prefix wording). Checked against a baseline taken before the move: a full generation from an empty index gives the same 65 tarballs (identical combined sha256), the same `index.json` and `new-versions.json` apart from timestamps, and `task registry:check` prints the same line. New tests cover generation, skipping a published version, `--check` writing nothing and rejecting a bad version, and a missing alias table.

---
title: Move `cmd/count-tokens/render.go` (451) into `internal/tokensreport`; leave `main.go` as the CLI
priority: p2
section: ralvaskills
topics:
    - ralvaskills
---

Resolved 2026-11-22: `cmd/count-tokens` is now a flag shim over `internal/tokensreport` (`Run(format, sortBy)`), with `render.go` and the unit tests moved along. Regeneration is byte-identical apart from the date stamp; `task lint` and `task tokens:lint` pass and the Taskfile paths include the new package.

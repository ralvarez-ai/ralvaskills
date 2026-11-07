---
title: 'Land the in-progress work that overlaps the module move in `autoreports`, `ralvapass`, `ralvasend` and `ralvaskills`'
priority: p2
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
`autoreports`: the module move and the rename of `automation-labs-work-hour-reports` to `autoreports` are committed and pushed for the files they touch (110, `HEAD` builds from a clean checkout); 11 modified files (docs rewrites) and 22 untracked ones are still unstaged. `ralvapass` (26 files plus `CHANGELOG`), `ralvasend` (58 files plus `CHANGELOG`) and `ralvaskills` (1 file) have unrelated edits in files the move also touched, left unstaged, so their commits may not build on their own. Commit the owner's work, then verify `go build ./...` from a clean checkout of each `HEAD`.

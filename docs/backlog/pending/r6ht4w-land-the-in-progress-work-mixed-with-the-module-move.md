---
title: 'Land the in-progress work that overlaps the module move in `autoreports`, `ralvapass`, `ralvasend` and `ralvaskills`'
priority: p2
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
`autoreports`: the working tree renames `automation-labs-work-hour-reports` to `autoreports` in 102 files; the module move sits on top and nothing is committed, so `HEAD` still imports the old path. `ralvapass` (26 files plus `CHANGELOG`), `ralvasend` (58 files plus `CHANGELOG`) and `ralvaskills` (1 file) have unrelated edits in files the move also touched, left unstaged, so their commits may not build on their own. Commit the owner's work, then verify `go build ./...` from a clean checkout of each `HEAD`.

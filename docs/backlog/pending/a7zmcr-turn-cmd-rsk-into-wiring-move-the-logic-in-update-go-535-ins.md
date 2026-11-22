---
title: 'Turn `cmd/rsk` into wiring: move the logic in `update.go` (535), `install.go` (426), `status.go` (424),'
priority: p2
section: ralvaskills
topics:
    - ralvaskills
---

`claude.go` (398), `catalog.go` (331) and `new.go` (304) into `internal/<area>`, leaving `cmd/rsk` with cobra registration only

Plan 2026-11-22: [plans/2026-11-22-cmd-rsk-into-wiring.md](../../plans/2026-11-22-cmd-rsk-into-wiring.md), six phases from `catalog` to `update`; three questions wait on the owner. Nothing is moved yet.

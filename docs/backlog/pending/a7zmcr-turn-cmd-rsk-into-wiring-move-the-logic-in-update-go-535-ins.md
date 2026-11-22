---
title: 'Turn `cmd/rsk` into wiring: move the logic in `update.go` (535), `install.go` (426), `status.go` (424),'
priority: p2
section: ralvaskills
topics:
    - ralvaskills
---

`claude.go` (398), `catalog.go` (331) and `new.go` (304) into `internal/<area>`, leaving `cmd/rsk` with cobra registration only

Plan 2026-11-22: [plans/2026-11-22-cmd-rsk-into-wiring.md](../../plans/2026-11-22-cmd-rsk-into-wiring.md), six phases from `catalog` to `update`; three questions wait on the owner. Nothing is moved yet.

Update 2026-11-22: phases 1 to 6 of the plan ran (catalog, claude, new, status, install, update); see "Where it ended" in the plan. The six files lost about 880 lines, less than hoped; the rest is prompts, rendering and the local-clone update path. Left pending for the remainder and the patch release.

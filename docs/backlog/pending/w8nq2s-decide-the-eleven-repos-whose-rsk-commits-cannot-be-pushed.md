---
title: 'Decide the eleven repos whose rsk commits cannot be pushed'
priority: p3
section: ralvaskills
topics:
    - rsk
    - git
related:
    - ralvaskills#m4pnr7
---
Two sweeps on 2026-11-12 committed `.rsk/SKILLS.md`, `.rsk/rsk.lock` and `mise.toml` in 145 repos, and 127 pushes landed. Eleven could not be pushed, for two reasons: six have no upstream at all — `ralvarez-go/ports`, `automation-labs-iec-contracts`, `automation-labs-iec-core` (on branch `spike/langium-st`), `automation-labs-plcverif-st-plugin`, `automation-labs-sparkov`, `automation-labs-openevolve-plc-adapter` — and five point at remotes that no longer exist — `notifier` and `protokit` (both archived), `desaparecidos-venezuela`, `strucgo`, `strucgo-runtime`, the last three carrying 6–14 unpushed commits of older work. Two more, `autoreports` and `uns-platform`, are committed but diverged while the rename migration was being re-pushed.

Every sweep rediscovers this, so it costs a pass each time. Done when each of the eleven is decided — re-point the remote, archive the repo, or declare it local-only and stop trying — and the decision is recorded here or in the item that owns that repo. `strucgo` and `strucgo-runtime` are probably the same decision as the module move to `go.ralva.dev`.

Note 2026-11-22: the roll-out of 3.10.0 will add to this list any repository whose push is refused (`9m3rtc`).

Note 2026-11-22 (night): the roll-out survey found eight more active main checkouts with no upstream, so a pin commit there stays local: `autobacklog`, `ralvadocket`, `autoreports`, `automation-labs-iec-core`, `automation-labs-coding-agents-benchmarks`, `automation-labs-openevolve-plc`, `automation-labs-openevolve-plc-adapter`, `ports`. Also 24 archived `ralvarez-go/*` repos hold an uncommitted `mise.toml` edit (3.7.0 to 3.9.1) and uncommitted `opencode.json`/`.opencode/` changes from an earlier sweep: decide commit, discard or leave before any bump there. No push has been attempted by the roll-out yet.

Note 2026-11-22 (batch 1): `ralvarez-ts/ui`, `ralvagym`, `ralvafi` and `automation-labs-openhands-plc` all have an upstream; their pin commits wait for a push after about 2026-11-23T05:21Z. `automation-labs-coding-agents-benchmarks` was not migrated: its main checkout sits on `feat/label-required` with no upstream.

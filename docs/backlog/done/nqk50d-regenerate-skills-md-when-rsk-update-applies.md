---
title: 'Regenerate .rsk/SKILLS.md when rsk update applies changes'
priority: p2
section: ralvaskills
topics:
    - cli
    - skills
related:
    - jefdb1
    - afwhb0
---
`rsk update` re-linked a skill and bumped `rsk.lock`, but never re-rendered `.rsk/SKILLS.md`, the file the tool adapters actually import. The skills are symlinks, so a pulled `DIGEST.md` was live, but the generated index kept the previous rules until the next `rsk pin` or `rsk install`. Both the header of `.rsk/SKILLS.md` and `docs/designs/2026-11-06-skill-digests-and-skills-md.md` claimed `rsk update` regenerated it, so nothing in review pointed at the gap.

Found 2026-11-11 while checking whether a digest edit reaches an already-pinned skill: `digest.Sync` was called only from `rsk new`, `rsk pin`/`unpin`, `rsk uninstall` and `install.Apply` (`internal/install/resolve.go:238`), and `cmd/rsk/update.go` did not import the package at all.

Done when `rsk update` re-renders the file in both local and registry mode, a test fails if it stops, and dry-run and an aborted confirmation still write nothing.

Resolved 2026-11-11 (`f439d52`): `refreshPinnedDigest` in `cmd/rsk/update.go` reads `rsk.mod` and calls `digest.Sync` after the lock bump, in `applyLocalUpdate` (post-pull) and in `runUpdateRegistry` (post re-link), both at project scope only. A write failure warns instead of failing an update that already succeeded, matching the lock bump. `TestRefreshPinnedDigest` pins the behaviour: a rewritten `DIGEST.md` moves into `SKILLS.md`, and the call is silent outside a project. Left over: the token-budget warning for the file is still `jefdb1`, and a registry-mode run where a skill fails to re-link returns before the refresh (the lock bump is skipped there too). Still unshipped: no CLI release has been cut since `f439d52` (the installed 3.6.0 is `6698ec8`), so consumers on 3.6.0 need an unpin/pin of one pinned skill to re-render `.rsk/SKILLS.md` — that is how the 2026-11-11 digest audit refreshed 110 projects.

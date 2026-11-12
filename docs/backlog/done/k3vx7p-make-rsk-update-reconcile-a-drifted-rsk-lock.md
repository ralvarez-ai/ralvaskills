---
title: 'Make `rsk update` reconcile a drifted `rsk.lock`'
priority: p3
section: ralvaskills
topics:
    - rsk
    - cli
related:
    - nqk50d
    - z8c3rn
---
`rsk update` decides whether a skill is current from the link in `.claude/skills/`, not from the version recorded in `.rsk/rsk.lock`. Found 2026-11-12 while verifying v3.7.0 on `ralvarez-go/corekit`: setting `docs-architect` to 1.2.0 in `rsk.lock` while the link still pointed at 1.2.1 made `rsk update docs-architect` print "All installed skills are up to date" and leave the stale lock in place. Re-pointing the symlink at 1.2.0 was what made the update apply.

Consequence: a lock that disagrees with the link — hand-edited, restored from an old commit, or left half-applied by an aborted update — never repairs itself, and `rsk status` reads the link, so the drift is invisible. The digest re-render shipped in v3.7.0 hangs off an update that never happens, so `.rsk/SKILLS.md` stays stale too.

Done when `rsk update` compares the lock as well as the link, applies whichever is behind, and a test covers a lock that disagrees with the link in both directions.

Resolved 2026-11-12 (`dacad1f`): `runUpdateRegistry` reads the project lock (`projectLockedVersions`, nil for `--global` or outside a project) and `findOutdatedSkills` treats a skill as outdated when either the link or the lock disagrees with the registry's latest (`needsUpdate`). A bare `rsk update` also checks lock entries whose link went missing (`withLockedNames`); a named update stays scoped to the names given — the first cut appended every locked skill and swept four unrelated updates into a one-skill run. `TestNeedsUpdate` covers six cases (including a lock ahead of latest), plus `TestWithLockedNames` and `TestProjectLockedVersions`. Verified end to end on `ralvarez-go/corekit`: lock at 1.20.0 with a current link, `rsk update go-architect` rewrote the lock to 1.20.2 and re-rendered `.rsk/SKILLS.md`. Shipped the same day in `v3.7.1`, and verified with the released binary on `ralvarez-go/corekit`, which was already drifted: its lock said `go-architect` 1.20.1 while the link pointed at 1.20.2, so the old binary called it current and 3.7.1 repaired it (`1.20.1 → 1.20.2`, lock rewritten, `.rsk/SKILLS.md` re-rendered) without touching another skill. Left over: the consumer side, which is `r9tk2v`.

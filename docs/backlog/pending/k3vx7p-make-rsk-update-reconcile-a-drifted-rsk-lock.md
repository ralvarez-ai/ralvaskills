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

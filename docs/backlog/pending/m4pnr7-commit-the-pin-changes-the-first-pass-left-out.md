---
title: 'Commit the pin changes the first pass left out'
priority: p2
section: ralvaskills
topics:
    - rsk
    - git
---
On 2026-11-06 one commit per repo, `chore(rsk): pin the shared skills and generate SKILLS.md`, went into 78 of the 79 repos whose `.rsk/rsk.mod` lists `commit-author` and whose `.git` is a directory (script: `/tmp/commit_pins.py`, list `/tmp/pinned.txt`; both are temporary). Left out: `automation-labs-titanic` ignores `.rsk` in git, so the pin never reaches history there. The earlier run pinned 87 projects, so about 8 repos (worktrees, nested repos, or `.git` files) are not covered. The untracked `.claude/skills/*` symlinks were not committed anywhere; the unrelated edits in `AGENTS.md` (`ratelimit`, `autoreports`) and `.gitignore` (`strucgo`) were left alone.

Done when the 8 uncovered pins are found and committed or ruled out, `automation-labs-titanic` is decided (track `.rsk` or accept it is local), and nothing pin-related is left uncommitted.

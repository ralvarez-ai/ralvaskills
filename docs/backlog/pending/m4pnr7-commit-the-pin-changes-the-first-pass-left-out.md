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

2026-11-12: the rename migration committed `.rsk/rsk.mod` + `.rsk/rsk.lock` in 30 repos (ADR 0002 step 3). `automation-labs-iec-core` was the one where a pin was already uncommitted — `go-architect` added to `pinned` — and it was left uncommitted as found, so it is still this item's work. The `@SKILLS.md` import drift those repos also carry is tracked separately in `h7v3cw`.

2026-11-12 (digest audit follow-up): the digest refresh produced two commit sweeps of its own — 82 repos got `chore(rsk): refresh the shared skill digests` (`.rsk/SKILLS.md` + `.rsk/rsk.lock` only, verified: no commit carried an unexpected path) and 63 got `chore(cli): pin rsk 3.7.0` (`mise.toml` only). 127 of those 145 pushes landed. Blocked: 11 repos, all for the same reason — no upstream at all (`ports`, `automation-labs-iec-core` on `spike/langium-st`, `sparkov`, `automation-labs-openevolve-plc-adapter`, `iec-contracts`, `plcverif-st-plugin`) or a remote that no longer exists (`notifier`, `protokit`, `desaparecidos-venezuela`, `strucgo`, `strucgo-runtime`, the last three carrying 6–14 unpushed commits of older work). `autoreports` and `uns-platform` are committed locally but diverged while the rename migration was being re-pushed; left alone rather than rebased under someone else. `automation-labs-titanic` still gitignores `.rsk`, unchanged.

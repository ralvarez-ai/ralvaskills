---
title: 'Stop the adapters importing whole skill bodies every session'
priority: p2
section: ralvaskills
topics:
    - cli
    - adapters
    - tokens
---
The Claude adapter writes `.rsk/CLAUDE.md` as `@../.claude/skills/<name>/SKILL.md` imports, and the OpenCode adapter lists the same files in `opencode.json` `instructions`. Both load every pinned skill's full body at session start, which defeats progressive disclosure (only the description should load until the skill activates). In this repo that is docs-architect + rsk-guide, about 4.9k tokens, plus `AGENTS.md` at about 1.9k, on every session.

**Update 2026-11-06:** the imports do not load at all. Claude Code expands an `@` import only when its target resolves inside the project; `.claude/skills/<name>` is a symlink into `~/.config/rsk/cache/`, so the body is skipped silently (tested with `claude -p`: a symlink to an outside directory and a plain outside file both failed; a real file, a nested import and a symlink to an inside target all loaded). Pinning therefore delivers nothing today, and the `rsk-guide` text saying pinned skills "load every turn" is wrong. Chosen direction: per-skill `DIGEST.md` plus a generated `.rsk/SKILLS.md` that the adapters import, see [the design](../../designs/2026-11-06-skill-digests-and-skills-md.md) and [the review](../../reviews/2026-11-06-skill-digests-and-compaction.md). All 60 digests are drafted; remaining work is the lint, `SPECS.md` and `skill-builder` changes, the adapter and `SKILLS.md` generation (ask first), and the release.

Original options: drop the body imports and rely on the native skill index (`.claude/skills/` is already discovered; check what OpenCode discovers natively), or import a generated index of names and descriptions instead. Changes `internal/tool/*`, so it touches every consumer project; needs a decision before work starts. Found while designing `token-economy`, whose audit mode flags this pattern.

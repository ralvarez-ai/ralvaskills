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

Options: drop the body imports and rely on the native skill index (`.claude/skills/` is already discovered; check what OpenCode discovers natively), or import a generated index of names and descriptions instead. Changes `internal/tool/*`, so it touches every consumer project; needs a decision before work starts. Found while designing `token-economy`, whose audit mode flags this pattern.

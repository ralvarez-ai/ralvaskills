---
title: 'Publish declared permissions in the registry index and `rsk mcp`'
priority: p3
section: ralvaskills
topics:
    - rsk
    - permissions
    - registry
related:
    - ralvaskills#qk4m2v
---
Phase 5 of [the permission plan](../../plans/2026-11-17-per-skill-permissions.md): a skill's `allowed-tools` is visible only after it is installed, because the registry index (`registry/index.json`) carries name, description, version and source but no permissions. `rsk catalog` therefore cannot answer "what will this skill ask me to allow?" before `rsk install` writes the preview, and `rsk mcp`'s `install_skills` returns the skills it installed without the rules they need.

Why it is not just one field: the declaration is client rule syntax (`Bash(...)`, `Read(...)`) while the index is meant to be client-neutral, so decide whether to publish the raw entries (simple, matches the tarball) or a normalized `{action, resource}` list (portable, needs a translation the registry does not have). Old binaries must keep reading an index that has the new field: it goes in `cmd/generate-registry` and `cmd/registry-remove` with the same compatibility care as `aliases` in ADR 0002, and `--check` should reject an entry that would fail `skill.ValidateAllowedTools`.

Done when a published skill's index entry lists its permissions, `rsk catalog` shows them (and still works against an index without the field), `install_skills` includes them in its result while still forcing user interaction, and `registry:check` fails on an invalid declaration instead of publishing it.

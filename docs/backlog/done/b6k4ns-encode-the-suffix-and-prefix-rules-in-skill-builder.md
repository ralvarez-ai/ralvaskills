---
title: 'Encode the suffix and prefix rules in skill-builder'
priority: p3
section: ralvaskills
topics:
    - skills
    - skill-naming
---
`skill-builder` scaffolds new skills from templates; it does not yet carry the naming rules [ADR 0002](../../adr/0002-rename-ten-skills-with-registry-aliases.md) settled, which is how the catalogue drifted to 44 `-architect` names in the first place.

Add to its checklist: the suffix table from the spec's Naming Convention (with `-designer`, `-writer`, `-format` and their meaning), the rule that **`-architect` requires a stack or protocol prefix** — no stack, no suffix — and the permanence rule (a published name is the registry key; a rename needs an ADR and an `[aliases]` entry). Point at `registry/aliases.toml` and `internal/skill/names.go` so the scaffolded skill is checked by `task skills:check` on the first run.

Done when a scaffolded skill with a retired or illegal name fails `task skills:check` before it is published, and `skill-builder` (1.1.0 → 1.2.0) is released.

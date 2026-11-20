---
title: Add a why skill that answers from the commit and PR record
priority: p3
section: ralvaskills
topics:
    - skills
    - pstack
---
Idea from the pstack review: a skill that answers "why is this code like this?" from the commit and
PR history (blame, linked PRs, review comments) before guessing. CodeGraph covers "how", so `how`
was left out of [the pstack adoption plan](../../plans/2026-11-20-adopt-pstack-skills.md); `why`
has no catalogue equivalent.

Upstream: `cursor/plugins` at `0b46ec399d55f2829e88822998f53b84944837d3`, `pstack/skills/why`
(MIT, Lauren Tan). Adapting it means the same attribution as that plan: pinned footer,
`LICENSE.upstream`, a `THIRD_PARTY_NOTICES.md` row. Needs its own plan; not built here.

---
title: 'Document DIGEST.md and lint it'
priority: p2
section: ralvaskills
topics:
    - skills
    - docs
    - lint
---
`DIGEST.md` exists in all 60 skills but is not in the specification, not scaffolded by `skill-builder`, and not checked. Add the entry to `docs/designs/2026-05-20-technical-specification.md`, scaffold it in `skill-builder`, and add a lint (the registry generator or `task lint`) that fails when a digest is missing, over 40 lines, has frontmatter, or does not end with `Full guide: invoke the <name> skill.`.

Also correct the `rsk-guide` body: it says pinned skills are imported and load every turn; they now contribute their digest through `.rsk/SKILLS.md`, which `rsk` generates. Skills publish on push, so bump the touched skills.

Done when `task lint` rejects a bad digest and `skill-builder` produces a valid one.

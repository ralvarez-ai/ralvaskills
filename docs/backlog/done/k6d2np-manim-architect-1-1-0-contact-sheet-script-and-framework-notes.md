---
title: '`manim-architect` 1.1.0: contact sheet script and the library''s own tools'
priority: p3
section: ralvaskills
topics:
    - manim-architect
    - explainer-architect
---

Review of a video version meant playing it. The skill named a contact sheet but only as a snippet in
`explainer-architect`, and it said nothing about what Manim itself provides.

**2026-11-11: done, published on push.** `manim-architect` is 1.1.0: `scripts/contact_sheet.py` (stdlib and ffmpeg) builds
one labelled image from a video and its sections file (a frame at 85% of each scene, or one every N seconds without a
sections file), one sheet per version. `FRAMEWORK.md` covers scene classes, formula parts and `TransformMatching*`,
animation families, `ValueTracker` with `always_redraw`, the built-in `next_section` / `--save_sections` and
`add_sound`, and the CLI flags. SKILL.md, DIGEST.md and STACK.md point to both. The same script and a library reference
also ship in the team's `video-narrator` skill (codex-atlas).

Open: `explainer-architect` RECIPES section 2 still carries the older unlabelled snippet; point it at this script once
that skill's uncommitted edits land.

# skill-builder digest

Scaffolds a new ralvaskills skill under `skills/<category>/<name>/`: SKILL.md skeleton, optional STACK and RECIPES files, SPECS.md updates.

## Always-on rules

- SPECS.md is canonical. Read it before guessing.
- Interview first: ask in order, wait for each answer. Order: topic and suffix, category, version target, bundle membership, cross-skill ties, attribution, body shape preview, decisions list.
- Name must match the SPECS suffix table. If no suffix fits, propose a new one and update SPECS. Do not shoehorn.
- Path is `skills/<category>/<name>/SKILL.md`. Never nest a SKILL.md.
- Description: one YAML scalar line, WHAT first (comma-separated), WHEN last, no filler, no `>` or `|`.
- Body: rules before code. Code blocks ~15 lines max. Body under ~10 KB / ~250 lines.
- Long code goes to the skill's RECIPES.md behind a one-line pointer.
- STACK: fetch live versions and pin them.
- A new skill needs 4-6 SPECS.md edits, applied in one pass.
- Add the attribution footer when adapted from an external source.
- Run the RECIPES validation checklist before declaring done.
- Hand-off: summarise what was generated and where, then ask whether to fill in content now.

## Hard constraints

- Never scaffold without an interview.
- Never guess. Read SPECS.md first.
- Never pin `latest`.
- Never leave a `📋 planned` row for a skill that now exists.

Invoke the full skill when: authoring a new skill, checking an existing skill folder against conventions, or moving or renaming a skill.

Full guide: invoke the skill-builder skill.

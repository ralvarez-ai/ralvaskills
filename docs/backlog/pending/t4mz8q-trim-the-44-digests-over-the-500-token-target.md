---
title: 'Trim the 44 digests that are over the ~500-token target'
priority: p3
section: ralvaskills
topics:
    - skills
    - tokens
related:
    - yg8j26
---
The 2026-11-11 digest audit ([review](../../reviews/2026-11-11-digest-audit.md)) put accuracy first and added about 220 tokens across the catalogue, so 44 of 67 `DIGEST.md` files are over the design's ~500-token target (measured at ~4 bytes/token; ~36k tokens in total). Largest: `go-architect` 955, `backend-performance-architect` 790, `repo-tooling-architect` 772, `wails-architect` 756, `python-architect` 732, `mcp-architect` 710, `go-library-builder` 704, `fastapi-architect` 695, `grpc-architect` 652, `uru-thesis-reviewer` 628.

Why it matters: `.rsk/SKILLS.md` concatenates the digests of the pinned skills, and `rsk` warns once that file passes its 2000-token budget. This repository's own file was ~3046 tokens before the 2026-11-11 refresh, and a project pinning six heavy skills can exceed the budget on digests alone. The warning itself is `jefdb1`.

Trimming is a content decision, not a format one: a rule is in a digest because it is always-on, so cutting one means deciding it belongs only in `SKILL.md`. The obvious candidates are the enumeration-heavy digests, where the digest repeats a list the body already has in a table — the naming and layout lists in `go-architect`, the transport menu in `backend-performance-architect`, the tooling defaults in `repo-tooling-architect`.

Constraints that still hold: plain markdown, at most 40 lines, no frontmatter, last line exactly `Full guide: invoke the <name> skill.`, extract-only from the skill's `SKILL.md`/`RECIPES.md`/`REFERENCE.md`/`STACK.md`. A digest edit publishes nothing without a `SKILL.md` version bump, and the digest must stay accurate to the body — the failure the audit fixed.

Done when every digest is at or under ~500 tokens, or the ones above it are recorded as deliberately over with the reason (a skill whose always-on rules genuinely need the room), and the target is restated in [the digest design](../../designs/2026-11-06-skill-digests-and-skills-md.md).

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

Constraints that still hold: plain markdown, at most 40 lines, no frontmatter, last line exactly `Full guide: invoke the <name> skill.`, extract-only from the skill's `SKILL.md`/`RECIPES.md`/`REFERENCE.md`/`STACK.md`. A digest edit publishes nothing without a `SKILL.md` version bump, and the shortened digest must stay faithful to the body — the defect class the audit was about.

Done when every digest is inside the band, or carries a recorded reason for sitting above it, and no digest keeps a justification clause, an enumeration or an example that the body already holds. Decided 2026-11-12: trim the fat and keep every rule — 500 stays the target, and a dense skill that genuinely carries 25–30 always-on rules lands at 700–750, which the design now accepts ([band](../../designs/2026-11-06-skill-digests-and-skills-md.md)) instead of paying for it with a dropped rule.

Resolved 2026-11-12: measured the catalogue at ~4 bytes/token and found the premise was the target, not the prose — 1,178 always-on rules in ~36k tokens, 26–42 tokens per rule, so a 500-token digest holds 12–21 rules. Trimming the fat (justification clauses, examples, enumerations the body already holds) buys 5–8%, not 35%. Decision recorded in [the design](../../designs/2026-11-06-skill-digests-and-skills-md.md): 500 stays the target and the band runs to 750, with every rule kept. After the trims the catalogue is 23 digests at or under 500, 41 inside the band, and one above it: `go-architect` at 877, which is its floor with 29 rules (~30 tokens each) and is now the recorded ceiling. Trimmed and published as patch releases: `go-architect` 1.20.2 (955 → 877), `backend-performance-architect` 1.0.2 (790 → 751), `repo-tooling-architect` 1.6.2 (772 → 730), `wails-architect` 1.0.2 (756 → 750). No rule was dropped from any digest.

The 110 consumer projects still hold the pre-trim digest text. No rule changed, so nothing is wrong there, but a `rsk update` in a project that pins these four now pulls the shorter versions.

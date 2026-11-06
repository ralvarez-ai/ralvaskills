---
title: 'Audit the 48 digests that were only checked for format'
priority: p3
section: ralvaskills
topics:
    - skills
    - tokens
---
Twelve digests were compared line by line with their skills (11 needed fixes: invented or over-strong wording, omitted hard constraints). The other 48 were drafted by a cheap model and checked mechanically only. Run the same audit (digest against `SKILL.md` and supporting files: unsupported claims, constraints worded stronger than the source, omitted never/always rules, paths that mislead without the skill) and fix what it finds.

Done when every digest has been audited; record the result in `docs/reviews/2026-11-06-skill-digests-and-compaction.md` or a new review.

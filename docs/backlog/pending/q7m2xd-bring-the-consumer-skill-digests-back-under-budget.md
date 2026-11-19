---
title: Bring the consumer skill digests back under the 2000-token budget
priority: p2
section: ralvaskills
topics:
    - rsk
    - digests
---
Pinning `test-designer` in the 13 library repos and this one pushed every
generated `.rsk/SKILLS.md` past the budget `rsk` warns at: 2709 tokens here
against a 2000-token budget, 4625 in `repokit`. The warning fires on every
`rsk install` and `rsk pin`, so it is now noise in 14 repos.

Decide between three options before adding another skill anywhere:

- prune the pinned set (the digest is per pinned skill, and `repo-tooling-architect`
  alone contributes 766);
- make `DIGEST.md` shorter per skill, which is a registry-wide change;
- raise the budget, which needs a measurement of what the extra context costs.

The warning is correct today: the file is loaded on every session in every
consumer, so this is a token cost paid on each start.

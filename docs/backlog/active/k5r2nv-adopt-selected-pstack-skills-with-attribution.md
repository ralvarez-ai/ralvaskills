---
title: Adopt selected pstack skills, with attribution
priority: p2
section: ralvaskills
topics:
    - skills
    - pstack
    - attribution
related:
    - ralvaskills#w3h8kd
---

Implement [the pstack adoption plan](../../plans/2026-11-20-adopt-pstack-skills.md). Nothing is
pushed; pushing `main` publishes every changed skill.

- [x] Phase 1: attribution scheme, notices, RECIPES section 9, `why` idea filed (`954b895`)
- [x] Phase 2: `blast-radius-reviewer`, `unslop`, `guardrail-builder` (`2072d2b`)
- [x] Phase 3: fold-ins into `performance-reviewer`, `token-economy` and the three cross-references (`f47ae7e`)
- [ ] Phase 4: retrofit `LICENSE.upstream` and pinned footers on the five adapted skills (staged, awaiting go-ahead to commit)
- [ ] Phase 5: counts (the plan says 67 to 70, but `task tokens` reports 72), `registry:generate`, `tokens`, CHANGELOG, `docs/README.md` row, `task site:check`
- [ ] Phase 6: dash sweep by script, `lint:dashes` task wired into `task lint`, one patch bump per touched skill
- [ ] Close: plan `status: done`, report each exit criterion with its proving command

Decided 2026-11-20: names as proposed; `guardrail-builder` in `global`; `unslop` rule 13 uses hyphens.
Open for the user: `caveman` ships `LICENSE-MIT` next to Apache `LICENSE` and `NOTICE`; keep or drop.

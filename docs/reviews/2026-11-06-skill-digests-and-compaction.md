---
status: done
project: ralvaskills
date: 2026-11-06
related:
  - designs/2026-11-06-skill-digests-and-skills-md.md
  - backlog/pending/q8vm2c-stop-the-adapters-importing-whole-skill-bodies-every-session.md
---

# Skill digests and lossless-compaction pass

Cheap agents drafted a `DIGEST.md` for all 60 skills (extract-only, 40 lines at most), listed
lossless-compaction candidates while reading each body, and then applied them. Nothing is
committed or version-bumped. Figures below are the agents' own reports unless marked measured.

## Findings
- **info** All 60 digests exist, are at most 40 lines, end with the exact pointer line and carry
  no frontmatter (measured). Together about 124 KB, roughly 520 tokens per skill.
- **info** An audit of 12 flagged digests found invented or over-strong wording in 11 and
  omitted hard constraints in most; all 11 were fixed. The other 48 digests were only checked
  mechanically.
- **info** Compaction: `SKILL.md` bytes went from about 615 KB to 586 KB (measured, about 5%).
  Most skills were already tight. Largest cuts: `repo-tooling-architect` 403 to 245 lines,
  `go-architect` 252 to 181, `website-concept-architect` 211 to 161, `grill-with-docs` 89 to 63,
  `astro-architect` 251 to 205, `python-architect` 146 to 119. Most savings are bulk moved to
  `RECIPES.md` or `REFERENCE.md` with a pointer left behind.
- **info** Medium-risk moves were skipped by the agents (for example `ui-ux-architect` tokens and
  ladder, `latex-architect` preamble list, `uru-scientific-paper-architect` typography). They
  remain candidates.
- **medium** `rsk-guide` says "Draft (v0.3.0)" while its frontmatter is 0.4.2, and its body
  claims pinned skills "load every turn", which the import test disproved.
- **medium** `mcp-architect` description is 1092 characters, over the 1024 limit.
- **medium** Broken links: `code-design-refactor` (`../tdd`, `../commit-author` should be under
  `../../workflows/`), `feature-planner` (`../grill-with-docs/CONTEXT.md` is missing),
  `skill-builder` (`../../../docs/SPECS.md` twice), `website-concept-architect` (links a
  missing skill), `hexagonal-arch` line 79 cites a section with no primary-port guidance. The
  tree had 39 broken links before this work and has the same set after, except one anchor in
  `nethttp-architect` that the pass broke and fixed (measured).
- **medium** Contradictions to decide: `commit-author` bans a `Co-Authored-By: Claude` trailer
  while the harness requires one; `caveman` ("terse for replies") against `token-economy`
  (never caveman in user replies); `nethttp-architect` puts Recover outermost while
  `gin-architect` has it third; `go-library-builder` pins Go 1.26 while `go-house-libs` says
  1.27.1; `go-architect` lines 16/181 (decoder lint) and 230/245 (`app/<service>` root) disagree;
  `nextjs-architect` says API routes for writes at line 54 and server actions at line 70;
  `ros2-architect` and its `STACK.md` word "default" differently.
- **medium** Possible factual errors to verify: `rest-api-architect` and
  `api-contract-reviewer` cite RFC 8594 for the `Deprecation` header (an agent says RFC 9745);
  `grafana-architect` line 91 error-budget formula looks wrong.
- **low** `demo-presentation-architect` cites a reference deck that is not in the folder.
- **low** `repo-tooling-architect` has `act` gotchas that may be specific to one project (a tunnel
  setup); check before this ships, since bundled skills must stay generic.
- **low** `ml-conference-paper-architect`, `hugo-architect` and `ui-ux-architect` repeat a line
  in their own `STACK.md`; layout bullets repeat across `gin`, `nethttp` and `fastapi`.

## Follow-ups
- Review the larger moves for information loss (done by independent agents on 18 skills, see the
  backlog item), then regenerate `TOKEN_COUNTS.md` with `task tokens`.
- Decide the contradictions and fix the broken links above in a separate change.
- Version bumps and the `SKILLS.md` generation follow the design's rollout; nothing ships until
  both are done.

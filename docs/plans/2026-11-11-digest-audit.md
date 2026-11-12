---
status: done
project: ralvaskills
date: 2026-11-11
related:
  - backlog/done/yg8j26-audit-the-48-digests-that-were-only-checked-for-format.md
  - reviews/2026-11-11-digest-audit.md
  - designs/2026-11-06-skill-digests-and-skills-md.md
---

# Digest audit: the 48 that were only checked for format

## Why

`docs/reviews/2026-11-06-skill-digests-and-compaction.md` records that 12 flagged digests were
compared line by line with their skills (11 had invented or over-strong wording and were fixed),
while the other 48 were drafted by a cheap model and only checked mechanically — 40 lines or
fewer, the pointer line, no frontmatter. That check cannot see a rule the body never states, a
`never` the body qualifies, or a constraint the digest silently drops.

## Method

Each digest is compared with its `SKILL.md` and the supporting files it leans on
(`STACK.md`, `RECIPES.md`, `REFERENCE.md`, `FRONTMATTER.md`, `TASKFILE.md`, `assets/`). Only
defects count:

1. **unsupported** — the skill's files never say it.
2. **over-strong** — `always`/`never`/`must` where the body allows an exception.
3. **omitted** — a hard constraint (never/always/ask-first/warn) the body states and the digest drops.
4. **misleading** — a path, command, flag or symbol that cannot work without the full skill.
5. **stale** — a version, name, flag or file the body contradicts.

Fix = edit the digest to match the body (extract-only: add nothing the body does not say) and bump
the skill's patch version, because a digest change publishes nothing until the version moves. Judged
and worded by hand or by an auditing agent, never generated from a summary; every suggested fix was
checked against the quoted source line before it was applied.

## Batches

All five batches were read in full on 2026-11-11: 67 digests, 26 skills fixed.

| Batch | Skills | Result |
|---|---|---|
| 0 — this repo's pinned set | `docs-architect`, `rsk-guide`, `token-economy`, `commit-author`, `repo-tooling-architect`, `adhd-mode` | 6 fixed, 8 defects |
| 1 — data, protocol, frontend | `ml-conference-paper-architect`, `study-artifact-builder`, `sql-architect`, `ddd-architect`, `hexagonal-arch`, `protobuf-architect`, `astro-architect`, `fastapi-architect`, `gin-architect`, `hugo-architect`, `manim-architect`, `nethttp-architect`, `nextjs-architect`, `react-architect`, `remotion-architect`, `reveal-js-architect` | 5 fixed, 6 defects |
| 2 — design, CI, infra | `wails-architect`, `demo-presentation-architect`, `demo-script-architect`, `explainer-architect`, `landing-page-designer`, `mockup-architect`, `narration-architect`, `ui-ux-architect`, `website-concept-architect`, `ci-cd-architect`, `docker-architect`, `grafana-architect`, `observability-architect`, `secrets-architect`, `go-architect`, `latex-architect` | 7 fixed, 7 defects; `go-architect` also trimmed 43 → 40 lines |
| 3 — languages, academic | `python-architect`, `event-driven-architect`, `caveman`, `skill-builder`, `go-house-libs`, `go-library-builder`, `pi-iteration-workflow`, `uru-scientific-paper-architect`, `uru-thesis-architect`, `uru-thesis-defense-architect`, `uru-thesis-reviewer`, `connectrpc-architect`, `grpc-architect`, `mcp-architect`, `rest-api-architect`, `api-contract-reviewer` | 5 fixed, 5 defects |
| 4 — reviews, workflow | `backend-performance-architect`, `codebase-hygiene-reviewer`, `performance-reviewer`, `security-reviewer`, `code-design-refactor`, `design-patterns`, `improve-codebase-architecture`, `logic-cleaner`, `ros2-architect`, `cli-tool-architect`, `feature-planner`, `grill-with-docs`, `tdd` | 3 fixed, 3 defects |

## Findings

Per skill, with the source line and the wording applied:
[reviews/2026-11-11-digest-audit.md](../reviews/2026-11-11-digest-audit.md). Totals: 29 defects
across 26 skills — 13 over-strong, 10 omitted, 3 misleading, none unsupported or stale — and 40
skills came back clean.

## Size and the pinned file

Batch 0 added about 180 tokens to the six digests, and the whole pass about 220 across the
catalogue. This repository pins from the **registry** (`config.repo_path` is empty, so
`.claude/skills/<name>` points into `~/.config/rsk/cache/registry/<name>/<version>`), which means
the change only reaches its own `.rsk/SKILLS.md` — today ~3046 of the 2000-token budget — after
these versions are published and `rsk update` re-links them; the working-tree edit alone changes
nothing here. Size is the separate question: 44 digests are still over the design's ~500-token
target, and only trimming rules that are not truly always-on brings a pinned file under budget.

## Remaining

1. Commit and publish: 27 skills carry digest edits and patch bumps; pushing publishes one release
   per skill, and CI regenerates `registry/index.json`.
2. The size question above — a content decision, not an audit fix.

# Documentation

Index and conventions for this folder. One doc type per folder; the folder and the
frontmatter `status` tell you the state.

| Folder | Holds |
|---|---|
| adr/ | Architecture decisions, numbered, immutable once accepted |
| backlog/ | project.md, one file per item in pending/, active/ and done/ |
| ideas/ plans/ designs/ | The path from thought to build |
| reviews/ changes/ | Dated records; never rewritten |
| guides/ methodology/ | Living procedures and working rules |
| templates/ archive/ | Doc templates; finished or obsolete docs |

Rules: frontmatter on every doc; relative links; drafts older than 2 months are promoted
or deleted; moving a file changes its state.
Routine: capture daily, triage weekly, prune monthly.
Check: `python3 <skill>/scripts/docs_lint.py docs`

## Index

## adr/

| No. | Doc | Status |
|---|---|---|
| 0001 | [Skill names keep `-architect` as the conventions marker](adr/0001-keep-the-architect-skill-suffix.md) | superseded by 0002 |
| 0002 | [Rename ten skills, and make a rename resolvable](adr/0002-rename-ten-skills-with-registry-aliases.md) | accepted |

## designs/

| Date | Doc | Status |
|---|---|---|
| 2026-05-20 | [ralvaskills — Technical Specification](designs/2026-05-20-technical-specification.md) | active |
| 2026-11-04 | [Landing page — round 1 and combined page](designs/2026-11-04-ralvaskills-landing-page.md) | draft |
| 2026-11-11 | [Skill rename aliases and manifest migration](designs/2026-11-11-skill-rename-aliases.md) | active |

## guides/

| Date | Doc | Status |
|---|---|---|
| 2026-05-22 | [Registry Flow](guides/registry-flow.md) | active |
| 2026-05-22 | [Registry Setup](guides/registry-setup.md) | active |

## plans/

| Date | Doc | Status |
|---|---|---|
| 2026-10-02 | [MCP Server Plan — `rsk mcp`](plans/2026-10-02-mcp-server-plan.md) | active |
| 2026-11-04 | [Landing page — work brief](plans/2026-11-04-landing-page-brief.md) | draft |

## reviews/

| Date | Doc | Status |
|---|---|---|
| 2026-11-06 | [Skill digests and compaction](reviews/2026-11-06-skill-digests-and-compaction.md) | active |
| 2026-11-11 | [Skill naming review — all 67 names](reviews/2026-11-11-skill-naming-review.md) | active |


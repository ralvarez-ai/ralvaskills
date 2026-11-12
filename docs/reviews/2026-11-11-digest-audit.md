---
status: done
project: ralvaskills
date: 2026-11-11
related:
  - plans/2026-11-11-digest-audit.md
  - reviews/2026-11-06-skill-digests-and-compaction.md
  - designs/2026-11-06-skill-digests-and-skills-md.md
---

# Digest audit: all 67 digests read against their skills

The 2026-11-06 pass drafted digests with cheap agents and audited 12 by hand; the other 48 were
only checked mechanically (40 lines, pointer line, no frontmatter). This pass read all 67 against
`SKILL.md` and the files a rule depends on (`STACK.md`, `RECIPES.md`, `REFERENCE.md`,
`FRONTMATTER.md`, `TASKFILE.md`, `assets/`), looking for five defects: unsupported, over-strong,
omitted hard constraint, misleading path, stale name. Method and batch roster:
[plans/2026-11-11-digest-audit.md](../plans/2026-11-11-digest-audit.md).

## Totals

- **67 audited, 29 defects fixed across 26 skills**; every fix is extract-only and each touched
  skill got a patch bump.
- **40 skills came back clean**; one of them, `go-architect`, needed a format fix only.
- Defect mix: 13 over-strong, 10 omitted, 3 misleading, 0 unsupported, 0 stale. Nothing was
  invented outright — the failure mode from the first pass had already been corrected on the 12
  that were audited then.
- A second pass over the *same* claim was worth it: three of the thirteen `over-strong` findings
  were rules stated as absolutes where the body qualifies them ("almost always", "in most cases",
  "unless an existing composite index covers it"), which is exactly the class a format check
  cannot see.

## Findings

### This repository's pinned six

| Skill | Defect | Fix |
|---|---|---|
| `docs-architect` | omitted: stale work is pruned, not archived (`SKILL.md:66,68`) | prune rule added |
| `rsk-guide` | omitted: the 2000-token budget (`SKILL.md:46`) | pin rule names the budget and the warning |
| `token-economy` | omitted: configuration beats instructions (`SKILL.md:46-50`) | "prefer a setting, then a hook, then a sentence" |
| `commit-author` | over-strong merge rule (`SKILL.md:38`); omitted harness-trailer rule (`SKILL.md:48`) | merge body made conditional; trailer clause added |
| `repo-tooling-architect` | omitted: nested `AGENTS.md` only for a package with its own commands (`SKILL.md:219`) | folded into the `AGENTS.md` rule |
| `adhd-mode` | omitted: the debug-spiral stop rule (`SKILL.md:41`); "cut words, never a verification" (`SKILL.md:49`) | both added |

### Data, protocol, frontend

| Skill | Defect | Fix |
|---|---|---|
| `astro-architect` | over-strong: schema is optional (`SKILL.md:48`); "rendered once per page" is the `src/components/` case (`SKILL.md:85`) | both qualified |
| `hugo-architect` | omitted: never rely on the global scope inside a partial (`SKILL.md:35`) | added |
| `study-artifact-builder` | over-strong: a question is skipped when the user already said it (`SKILL.md:24`) | exception added |
| `nextjs-architect` | over-strong: Playwright may run against `next dev` in CI (`SKILL.md:165`) | alternative added |
| `nethttp-architect` | over-strong: kin-openapi is intended but a committed `openapi.yaml` is allowed (`SKILL.md:106`) | reworded |

### Design, CI, infra

| Skill | Defect | Fix |
|---|---|---|
| `wails-architect` | over-strong: a private-repo GitHub provider needs a token in the binary (`SKILL.md:77`) | caveat and `endpoint` alternative added |
| `demo-presentation-architect` | over-strong: dividers skipped for decks of 8 slides or fewer (`SKILL.md:52`) | exception added |
| `explainer-architect` | over-strong: with audio the length is `voice + lead-in + tail + hold` (`SKILL.md:27`) | both formulas stated |
| `landing-page-designer` | misleading: the script is `verify.js` (`SKILL.md:99`) | path fixed |
| `mockup-architect` | over-strong: only pages that compare options are retrofitted (`SKILL.md:48`) | scope added |
| `grafana-architect` | over-strong: anonymous access off where data is not public (`SKILL.md:79`) | reworded |
| `latex-architect` | omitted: `\\` to end a paragraph is banned (`SKILL.md:136`) | added to the banned list |

### Languages, academic

| Skill | Defect | Fix |
|---|---|---|
| `event-driven-architect` | over-strong: several handlers on a command is *almost always* wrong (`SKILL.md:19`) | reworded |
| `uru-scientific-paper-architect` | over-strong: by method and intent, not just subject matter (`RECIPES.md:105`) | reworded |
| `uru-thesis-defense-architect` | misleading: the logo is at `assets/uru-logo.png` (`SKILL.md:15,42`) | path fixed |
| `mcp-architect` | omitted: servers MUST implement `server/discover` (`SKILL.md:145`) | added |
| `rest-api-architect` | omitted: HTTPS always, no plaintext in production (`SKILL.md:83`) | added |

### Reviews, workflow

| Skill | Defect | Fix |
|---|---|---|
| `backend-performance-architect` | over-strong: "every FK" unless a composite covers it (`RECIPES.md:77`) | qualified |
| `codebase-hygiene-reviewer` | omitted: dead only when every consumer reports it unreachable (`RECIPES.md:31`) | added |
| `cli-tool-architect` | over-strong: "no more than 2 positional args **in most cases**" (`SKILL.md:23`) | reworded |

### Format

`go-architect`'s digest was 43 lines against the 40-line cap after 1.19.0 and 1.20.0. Three lines
restated a rule already in the file — host:port via `fmt.Sprintf`, hand-editing generated sqlc
output, libraries importing `fx`/`wire` — so they were folded into the rule above them: 40 lines,
no rule lost.

## What the audit did not settle

- **Size.** 44 of 67 digests are still over the design's ~500-token target (largest now
  `go-architect` 952, `backend-performance-architect` 786, `repo-tooling-architect` 772). This pass
  added about 220 tokens across the catalogue, because accuracy came first. Bringing a project's
  pinned file under the 2000-token budget is a content decision — which rules are not truly
  always-on — not an audit fix.
- **The `rsk-guide` body** still calls itself a draft while the CLI is shipped, and the 2026-11-06
  review's other open items (contradictions, broken links, RFC 8594 vs 9745, the Grafana
  error-budget formula) are untouched here.
- **Publishing.** 27 skills carry digest edits and patch bumps; nothing is committed or pushed, and
  a push publishes one release per skill with the index regenerated by CI.

## Follow-ups

- Decide and publish the bump batch (or split it), then `rsk update` in the projects that pin these
  skills so their `.rsk/SKILLS.md` picks up the corrected digests.
- The size question above, as its own change.

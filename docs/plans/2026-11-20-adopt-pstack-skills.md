---
status: done
project: ralvaskills
date: 2026-11-20
related:
  - designs/2026-05-20-technical-specification.md
  - adr/0002-rename-ten-skills-with-registry-aliases.md
  - guides/registry-flow.md
---

# Adopt selected pstack skills, with attribution

[pstack](https://github.com/cursor/plugins/tree/0b46ec399d55f2829e88822998f53b84944837d3/pstack)
is Lauren Tan's (poteto) skill pack for Cursor, MIT-licensed, 58 skills. Most of it overlaps
this catalogue or depends on Cursor features. This plan takes the parts that fill real gaps:
three new skills and three fold-ins into existing skills. It also sets up attribution that ships
inside each published tarball, and fixes that gap for the five skills we already adapted from
other sources.

Upstream pin: `cursor/plugins@0b46ec399d55f2829e88822998f53b84944837d3` (2026-11-19). Every
link and diff in this plan uses that commit, not `main`.

## Scope

**In: new skills**

| pstack source | New skill | Domain | Bundle | What it adds |
|---|---|---|---|---|
| `skills/blast-radius` | `blast-radius-reviewer` | `quality/` | `code-review` | A five-level proof ladder (claimed, cited, reasoned, ran, reproduced). Name the one fact the change is safe because of, and prove it with a script. |
| `skills/unslop` | `unslop` | `meta/` | `docs` | Numbered rules that strip AI patterns from prose. Rule ids are stable so other skills can cite them. |
| `skills/correct` | `guardrail-builder` | `workflows/` | `global` | Turn a repeated agent mistake into enforcement, in order: architecture, types, a lint whose error names the fix, a test, docs last. Keep a rule-to-enforcer table in `AGENTS.md`. |

**In: fold-ins (patch or minor bump on an existing skill)**

| pstack source | Target skill | Change |
|---|---|---|
| `skills/benchmark-checklist`, `skills/principle-explain-the-number` | `performance-reviewer` (new `BENCHMARKS.md` reference) | The seven questions (why not double, tuned, limits, errors, reproduces, matters, happened), alternating A/B runs, and the report shape with run count, range and limiter. |
| `skills/principle-attack-the-premise` | `token-economy` | Extend the two-failure rule: before a third fix, write down the premise the failed fixes share and test it. |
| `skills/unslop` (cross-reference only) | `docs-architect`, `commit-author`, `narration-writer` | One line that points to `unslop` for prose review. These skills take no pstack text, so they get no attribution change. |

**Out**

- `poteto-mode`, `setup-pstack`, `arena`, `swarm`, `automations/benny`: these rely on Cursor's
  multi-model routing, model settings or automations, or encode one person's working style.
- `tdd`, `interrogate`, `architect`, `principle-test-behavior-not-implementation`,
  `principle-laziness-protocol`, `principle-guard-the-context-window`: already covered by `tdd`,
  `grill-with-docs` and `security-reviewer`, `feature-planner`, `test-designer`, `logic-cleaner`
  and `token-economy`.
- `how`: CodeGraph answers most of it. `why` (answers from the commit and PR record) is a
  candidate for a later plan, filed as a backlog idea and not built here.
- `statechart` (TLA+), `create-*-verification-skill`, `ketchup*`, `bro`, `teach`,
  `typescript-best-practices`, the remaining `principle-*` skills: too niche, or no gap to fill.
- Syncing with upstream. These are one-time adaptations pinned to one commit. Pulling a later
  upstream change means a new diff against the pin and a normal version bump.

## Attribution

The repository is GPL-3.0. Including MIT code is allowed, but the MIT license requires keeping
the copyright and permission notice in every copy or substantial portion. A footer link alone
does not satisfy that, and each skill is published as its own tarball. So the notice has to live
in the skill directory. `cmd/generate-registry` `createTarball` packs every file in that
directory (`cmd/generate-registry/main.go:166`), so a file placed there ships with every install.

Every skill that takes pstack text, new or fold-in, gets all three of the following:

1. **Footer** at the end of `SKILL.md`, following
   [skill-builder RECIPES §9](../../skills/meta/skill-builder/RECIPES.md) and the existing
   `tdd` and `caveman` footers, linked to the pinned commit:

   ```markdown
   ---

   _Adapted from [cursor/plugins pstack `blast-radius`](https://github.com/cursor/plugins/tree/0b46ec399d55f2829e88822998f53b84944837d3/pstack/skills/blast-radius) by Lauren Tan (MIT). Changed for this workspace._
   ```

   A fold-in uses `_Includes material adapted from ..._` and lists each pstack source skill.

2. **`LICENSE.upstream`** in the skill directory. It holds a header naming the source path and
   commit, then the upstream MIT text unchanged, including `Copyright (c) 2026 Lauren Tan`. A
   skill with more than one upstream gets one section per upstream in the same file.

3. **Row in `THIRD_PARTY_NOTICES.md`** at the repository root. This file (it already exists, with
   rows for the five adapted skills) lists every
   upstream: project, URL at the pinned commit, license, copyright holder, and the skills that
   use it. It covers pstack, `mattpocock/skills`, `JuliusBrussee/caveman` and
   `ayghri/i-have-adhd`.

`DIGEST.md` files carry no attribution. The digest is a summary that loads into context, and the
attribution lives in the files that ship.

### Existing adapted skills

| Skill | Upstream | License | What is missing |
|---|---|---|---|
| `tdd` | mattpocock/skills | MIT | `LICENSE.upstream`; footer links `main`, not a commit |
| `grill-with-docs` | mattpocock/skills | MIT | same |
| `improve-codebase-architecture` | mattpocock/skills | MIT | same |
| `caveman` | JuliusBrussee/caveman | Apache-2.0 | `LICENSE.upstream` with the Apache-2.0 text, plus the upstream `NOTICE` file if it has one (Apache §4d) |
| `adhd-mode` | ayghri/i-have-adhd | MIT | `LICENSE.upstream`. The copyright line is already inline, but the permission notice is not. Footer links the repo, not a commit. |

`hugo-architect` (credits the Hugo documentation) and `narration-writer` (credits an idea from a
post by Andrej Karpathy) credit ideas, not copied text, so their footers are enough.

Fixing these is a patch bump on each, so all five republish. That happens in phase 4, as its
own commit.

## Phases

1. **Prepare.** Add `THIRD_PARTY_NOTICES.md`. Add a one-line rule to skill-builder RECIPES §9:
   pin the footer link to a commit, and add `LICENSE.upstream` whenever upstream text is copied
   or closely paraphrased. File the `why` idea in `docs/backlog/`. No skill is published.
   *Outcome:* the attribution scheme is written down before the first skill uses it.

2. **New skills.** Build all three with `skill-builder`: `SKILL.md` at `version: 1.0.0`,
   `DIGEST.md`, footer, `LICENSE.upstream`. Changes from upstream:
   - `blast-radius-reviewer`: replace the `arena` step with an optional second-opinion
     subagent; replace Cursor `Task` and `readonly` wording with neutral language; point the
     callers step at CodeGraph when `.codegraph/` exists; cite `unslop` for the writeup. Keep
     the proof ladder and the hand-back format as written.
   - `unslop`: keep upstream rule numbers and leave the gaps, so a later upstream diff stays
     readable. Rewrite rule 13 as: "No em dashes or en dashes. Prefer a period or a comma;
     where a dash reads better, use a spaced hyphen (` - `). Ranges use an unspaced hyphen
     (`2-4`)." Upstream also bans hyphens used as dashes; this catalogue allows them.
   - `guardrail-builder`: keep the four-level order. Tie the rule table to
     `repo-tooling-architect`'s `AGENTS.md` rules, and tie the "configuration beats
     instructions" line to `token-economy`.
   Add each skill to its bundle in `internal/config/catalog.toml`.
   *Outcome:* three new skill directories that pass `task registry:lint` and `task tokens:lint`.

3. **Fold-ins.** `performance-reviewer` gets a new `BENCHMARKS.md` reference with the seven
   questions, plus a pointer from `SKILL.md` (minor bump). `token-economy` gets the
   attack-the-premise lines in `SKILL.md` and `DIGEST.md` (patch bump). Each gets a footer and a
   `LICENSE.upstream`. Add the `unslop` cross-reference line to `docs-architect`,
   `commit-author` and `narration-writer` (patch bumps).
   *Outcome:* the changed skills publish with their notices.

4. **Retrofit existing attributions.** Add `LICENSE.upstream` to `tdd`, `grill-with-docs`,
   `improve-codebase-architecture`, `caveman` and `adhd-mode`. Re-point their footers to the commit each was
   adapted from, or to the current upstream commit if that is unknown, and say which in the
   commit message. Patch bump each.
   *Outcome:* every adapted skill in the catalogue ships its upstream license.

5. **Counts and records.** The landing-page gate changes from 67 to 70 skills
   (`site/scripts/gen-data.mjs:8,78`, `site/src/data/facts.ts:177`, and any copy that states
   the count). The bundle count stays at 15. Run `task registry:generate` and `task tokens`.
   Add a `CHANGELOG.md` entry and an index row in `docs/README.md`.
   *Outcome:* `task site:check` passes against the new tree.

6. **Dash sweep.** Replace every em dash and en dash in the catalogue with a hyphen, per the
   rewritten `unslop` rule 13. As of 2026-11-20 that is 2,643 occurrences across nearly every
   skill directory, plus 11 non-Markdown files. Do it with a script, not by hand: ` — ` becomes
   ` - `, and `–` between numbers becomes `-`. Review the diff by hand only where a dash sits
   inside a code fence, an output sample or a non-Markdown file, because those may be quoting
   real output. Patch bump every touched skill. Add a `lint:dashes` task that fails when `—` or
   `–` appears under `skills/**`, and add it to `task lint`, so the rule is enforced and not
   just written down. The new skills from phases 2 and 3 are written with hyphens from the
   start.
   *Outcome:* no em or en dash under `skills/**`, and `task lint` keeps it that way. This
   republishes almost the whole catalogue in one push, so it gets its own commit and its own
   go-ahead.

Phases 2 to 4 and 6 edit `skills/**` and publish on push to `main`.

**One bump per skill per push.** Phases 3, 4 and 6 can touch the same skill (`token-economy`,
`tdd` and others). Before bumping, compare the skill's `SKILL.md` version with its entry in
`registry/index.json`. If it is already ahead, an earlier phase bumped it, so leave it. Use a
minor bump when any phase in the push adds content (phase 3), and a patch bump otherwise. Each phase is one or more commits,
and each needs its own go-ahead before it is pushed.

## Open decisions

1. ~~**`unslop` rule 13 (no em dashes).**~~ Decided 2026-11-20: keep the rule, use hyphens
   instead (phase 2 wording, phase 6 sweep).
2. ~~**Names.**~~ Decided 2026-11-20: `blast-radius-reviewer`, `unslop`, `guardrail-builder`
   (permanent registry keys, ADR 0002). pstack's own name `correct` is too generic.
3. ~~**`guardrail-builder` bundle.**~~ Decided 2026-11-20: `global`. It adds about one digest's
   worth of tokens wherever it is pinned.

## Exit criteria

- [x] `blast-radius-reviewer`, `unslop` and `guardrail-builder` exist with `version: 1.0.0`,
      `DIGEST.md`, a pinned footer and `LICENSE.upstream`
- [x] `performance-reviewer` and `token-economy` carry the fold-ins, a footer and
      `LICENSE.upstream`
- [x] `tdd`, `grill-with-docs`, `improve-codebase-architecture`, `caveman` and `adhd-mode` ship
      `LICENSE.upstream`, and their footers link a commit
- [x] `grep -rE '—|–' skills` finds nothing, and `task lint:dashes` runs as part of `task lint`
- [x] `THIRD_PARTY_NOTICES.md` lists every upstream and every skill that uses it
- [x] An extracted tarball of each touched skill contains `LICENSE.upstream` (check with
      `tar tzf` on the `task registry:generate --force` output; the task's `sources` glob is
      `skills/**/*.md`, so without `--force` it can skip a change that only adds
      `LICENSE.upstream`)
- [x] `task lint`, `task registry:check`, `task tokens:lint` and `task site:check` pass
- [x] `CHANGELOG.md` and the `docs/README.md` index are updated

---
status: active
project: ralvaskills
date: 2026-11-06
related:
  - backlog/done/q8vm2c-stop-the-adapters-importing-whole-skill-bodies-every-session.md
---

# Skill digests and a generated `SKILLS.md`

## Problem

`rsk pin` writes `@../.claude/skills/<name>/SKILL.md` imports into `.rsk/CLAUDE.md` (and the
same paths into `opencode.json` `instructions`). Two defects, measured 2026-11-06 with
throwaway projects and `claude -p`:

1. **The imports do not load.** Claude Code expands an `@` import only when the target resolves
   inside the project. `.claude/skills/<name>` is a symlink into `~/.config/rsk/cache/`, outside
   every project, so the body is skipped silently. Controls: a real file, a nested import and a
   symlink whose target is inside the project all loaded; a symlink to an outside directory and
   a plain outside file did not. Transcripts of real sessions confirm no pinned body is in
   context. So "pinned" currently delivers nothing.
2. **If they did load, they would cost too much.** A full body is about 5k tokens per skill,
   paid on every turn, which defeats progressive disclosure.

## Approach

Pinning becomes the single switch for "always in context", and what loads is a short digest,
not the body.

- Each skill may ship a **`DIGEST.md`** beside `SKILL.md`.
- rsk generates **`.rsk/SKILLS.md`** from the digests of the pinned skills. It lives inside the
  project, so imports of it load.
- The adapters import `.rsk/SKILLS.md` instead of the per-skill `SKILL.md` files.
- The full body still loads on demand through the Skill tool; every installed skill's
  description is still in the skill list. Unpinned skills cost only their description.

The skill author does not decide whether a skill is always loaded: that would charge every
consumer for the author's choice. The project decides, by pinning.

## Components and interfaces

### `DIGEST.md` (per skill, optional)

- Plain markdown, no frontmatter, hard cap of about 40 lines (about 500 tokens). The size is a
  band rather than a number: 500 is the target, and a dense skill that carries 25–30 always-on
  rules lands at 700–750, which is accepted (2026-11-12) rather than paid for by dropping a rule.
  What is not accepted is prose that restates the body — a justification clause, an enumeration
  the body already has in a table, an example — so trimming is always the first move.
- Contents: the always-on rules, the hard constraints ("never X without asking"), and one line
  saying when to invoke the full skill. Extracted from the body, nothing new.
- Last line: `Full guide: invoke the <name> skill.`
- Optional opt-out for skills with nothing worth loading every session: frontmatter
  `metadata.digest: none`.
- No `DIGEST.md` and no opt-out: rsk falls back to `- <name>: <description>`.
- Added to `SPECS.md` as a supporting file; `skill-builder` scaffolds it; `task lint` checks
  the cap and that the pointer line exists.

### `.rsk/SKILLS.md` (generated)

- Written by `rsk pin`, `rsk update` and `rsk new`; never edited by hand. Header says so.
- One `## <name> v<version>` section per pinned skill, in `rsk.mod` order, holding that
  skill's digest verbatim.
- Deterministic: same pins and versions produce the same bytes. No model is involved at
  generation time.
- A total budget (proposed 2k tokens) is checked on generation; rsk warns when pins exceed it.

### Adapters (`internal/tool/*`)

- Claude: `.rsk/CLAUDE.md` contains the single line `@SKILLS.md`.
- OpenCode: `opencode.json` `instructions` lists `.rsk/SKILLS.md`. `AGENTS.md` itself is not
  used for this, because OpenCode does not expand `@` imports in it.
- `AGENTS.md` stays human-written. rsk keeps seeding it once and never edits an existing one.

### Authoring the first digests

An agent drafts each `DIGEST.md` once, extract-only: copy rules that already exist in the body,
add none. A pilot on three skills (`token-economy`, `docs-architect`, `rsk-guide`) is reviewed
by hand against their bodies before the other skills are generated. Output is checked by the
lint, not by trust.

## Rollout

1. This design and an update to backlog item `q8vm2c`.
2. Pilot digests for the three pinned skills; review.
3. Digests for the remaining skills, in batches.
4. Lint and `SPECS.md`/`skill-builder` changes.
5. Adapter and `SKILLS.md` generation (ask first: it changes every consumer project).
6. Bump versions and publish in one batch once step 5 ships. Adding a `DIGEST.md` publishes
   nothing until the skill's version is bumped and pushed, and bumping all skills at once cuts
   one release per skill.

## Risks

- **Digest drifts from the body** — record the skill version in the `SKILLS.md` section header;
  regenerate on `rsk update`; the lint compares the digest's rules to the body for the
  pointer line only, so review in PRs still matters.
- **A drafting agent drops a constraint or invents one** — extract-only brief, hand review of
  the pilot, short enough to read in full.
- **Pins exceed the budget** — warn at generation; the budget is a number in config, not a
  rule.
- **Existing consumers** — until step 5 ships nothing changes for them; after it, a re-run of
  `rsk pin` rewrites `.rsk/CLAUDE.md` and `opencode.json`.
- **Claude Code later expands external imports** — the design still holds: the digest is what
  we want loaded, not the body.

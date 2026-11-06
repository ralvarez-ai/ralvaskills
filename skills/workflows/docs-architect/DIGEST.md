# Docs Architect digest

Purpose: one `docs/` tree where the folder is the doc type, frontmatter is the state, and moving a file is a state change.

## Always-on rules

- Pick the folder by the doc's job: `adr/` one decision; `ideas/` pre-plan; `plans/` scope and exit criteria; `designs/` how it is built; `reviews/` and `changes/` dated history; `guides/` and `methodology/` living `topic.md`; `backlog/`; `archive/` finished or obsolete.
- Scaffold: create `docs/README.md` and `docs/templates/`; create other folders when their first doc exists.
- Write-doc: copy the template, fill the frontmatter, link the doc to its origin.
- Frontmatter on every doc except `README.md`, `inbox.md`, `templates/`: `status`, `project`, `date`, `related`. ADRs use their own schema.
- Use relative markdown links; `related` paths are relative to `docs/`.
- Link ADRs from designs (`[ADR 0002](../adr/0002-slug.md)`) instead of restating reasoning.
- Write an ADR only for a real architecture decision (hard to reverse, surprising, real trade-off).
- ADR numbers: sequential; take the next by scanning `adr/`.
- Changing a decision: new ADR with `supersedes: N`; old one gets `status: superseded` and `superseded_by`.
- Backlog: queue layout (`inbox.md`, `<project>.md`, `done/YYYY-MM.md`) or item layout (`project.md`, `pending/`, `active/`, `done/`, `archive/`).
- Flow: inbox → ideas → plans → designs → adr → backlog. On close: item to `done/`, doc `status: done`, then `archive/`.
- Names: kebab-case, lowercase, no accents (`0004-use-nats.md`).
- Run the skill's `docs_lint.py [docs-dir]` (stdlib, read-only) before declaring a docs change done.

## Hard constraints

- Never edit an accepted ADR body; only `status` and `superseded_by`.
- Never rewrite reviews or changes; add a new dated doc.
- Never reuse an ADR number. Keep `adr/` flat; no per-project subfolder there.
- Create a per-project subfolder inside a type only after ~10 files of that type.
- Migrate: propose moves and renames, wait for approval, then move.
- Never duplicate these conventions in `CLAUDE.md`/`AGENTS.md`; point them at `docs/`.

Invoke the full skill when: setting up `docs/`, writing a plan/design/ADR/review, triaging a backlog, or reorganizing an existing docs folder.
Full guide: invoke the docs-architect skill.

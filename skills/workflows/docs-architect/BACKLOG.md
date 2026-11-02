# Backlog

Two layouts. Both keep one source of truth per item; pick one per repo by size and by how many branches touch the backlog, and do not mix them.

| | Queue layout (a list) | Item layout (a file per item) |
|---|---|---|
| Shape | `inbox.md` + `<project>.md` with sections + `done/YYYY-MM.md` | `project.md` + one file per item under `pending/`, `active/`, `done/`, `archive/` |
| Good when | a handful of items, one branch at a time, read top to bottom | dozens of items, parallel branches or agents, a portal or a tool that indexes it |
| State | the section header | the folder |
| Transition | move the line | rename the file (`git mv`); `git log --follow` keeps the thread |
| Order | the section order | priority, then git creation date |
| Two branches adding items | conflict line by line | no conflict: new file each |
| Checked by | the linter's size limits (`Doing` 2, `Next` 5) | `rdk check` (schema, ids, references) |

A repo that carries a list next to `project.md` has two sources of truth: convert one or delete the other.

## Queue layout

```
docs/backlog/
├── inbox.md          # raw capture, one line per thought, no format
├── <project>.md      # Doing / Next / Later / Someday
└── done/
    └── 2026-10.md    # closed items, one file per month
```

### Item format

```md
- [ ] Verify TON timers in the PLC agent  #plc-agent  P1
  plan: [plan](../plans/2026-10-12-plc-verification.md)
```

Every item has a project tag, a priority (P1/P2/P3) and, once it leaves the inbox, a link to the doc it came from (plan, design, review or change). An item without an origin link is allowed only in `Someday`.

### Sections per project file

| Section | Limit | Meaning |
|---|---|---|
| Doing | 1–2 | In progress now |
| Next | 3–5 | Ready to pick up; linked to an approved design |
| Later | none | Wanted, not scheduled |
| Someday | none | Maybe never |

The limits are what make the list a queue and not a pile. The linter warns when they are exceeded.

### Triage (weekly)

Take each inbox line and do exactly one of: delete it; append it to the project file with tag and priority; or open an `ideas/` doc when it needs thought. Leave the inbox empty.

### Pruning (monthly)

- Items in `Later` for 2–3 months without moving: delete, or demote to `Someday`. Do not archive them.
- Closed items move to `done/YYYY-MM.md` with the close date, keeping the origin link.

### Closing an item

1. Move the line to `done/YYYY-MM.md`, tick the box, add the date.
2. Set the originating doc to `status: done` once all of its items are closed.
3. Move that doc to `archive/` at the next monthly pass.

## Item layout

```
docs/backlog/
├── project.md        # slug, name, sections
├── pending/          # wanted, not scheduled (this includes blocked and parked work)
├── active/           # in progress now
├── done/             # closed; never rewritten
└── archive/          # dropped or obsolete
```

### `project.md`

```yaml
---
slug: plc-agent
name: PLC agent
sections:
    - corpus
    - verifier
---
```

- `slug` is the name other backlogs use to reference this one (`plc-agent#k7f3qz`).
- `sections` is a non-empty list of 1 to 4 themes. The section is a field of the item, not a folder: a new theme is a line here, not a new directory.
- The body is free; use it for what the queue layout kept in prose (where the project stands, known limits). Nothing that must be read as an item belongs there.

### Item file

Name: `<id>-<slug>.md`.

- `<id>` is six characters, random, lowercase, from the Crockford base32 alphabet without `i`, `l`, `o` and `u` (`[0-9a-hjkmnp-tv-z]{6}`). Random, not sequential: two branches that create items never collide, and no counter has to be serialised.
- `<slug>` is a kebab-case summary of the title, cut at a word boundary. It is for humans and `grep`; the id is the identity, so renaming the slug never breaks a reference.

```yaml
---
title: Verify TON timers in the PLC agent
priority: p1
section: verifier
topics:
    - timers
blocked_by:
    - k7f3qz
related:
    - plc-platform#a1b2c3
---
The body says what done means and the evidence that closes it: a commit, a
command, a document. Links to the origin go here:

plan: [plan](../plans/2026-10-12-plc-verification.md)
```

| Field | Required | Rule |
|---|---|---|
| `title` | yes | One line. Keep it short: the detail is the body, not the title. |
| `priority` | yes | `p0` (the project stops for it), `p1` (next), `p2` (this cycle), `p3` (when it fits) |
| `section` | yes | Must be one of the `sections` in `project.md` |
| `topics` | no | Free tags, for filtering |
| `blocked_by` | no | `<id>` in this backlog, or `<project>#<id>` across backlogs; each one must exist |
| `related` | no | Same syntax; a soft link, not a dependency |
| `resolution` | no | `dropped`, only in `archive/` |

Forbidden: `status` (the folder says it) and any date field (creation and closing come from git). A date written in the body is prose, not state.

### Lifecycle

State transitions are renames, in the same commit as the work they describe:

```
git mv docs/backlog/pending/k7f3qz-verify-ton-timers.md docs/backlog/active/
git mv docs/backlog/active/k7f3qz-verify-ton-timers.md docs/backlog/done/
git mv docs/backlog/pending/k7f3qz-verify-ton-timers.md docs/backlog/archive/   # add resolution: dropped
```

| Folder | Meaning | How it gets there |
|---|---|---|
| `pending/` | Wanted, not started. Blocked, parked and "someday" work lives here too, marked with `blocked_by` or a `p3` | Triage of an idea, a review or a plan |
| `active/` | Someone is on it now | `git mv` on the day work starts |
| `done/` | Closed, with the evidence in the body | `git mv` on the day it closes |
| `archive/` | Decided against, or obsolete | `git mv`, plus `resolution: dropped` when it was dropped |

`done/` never empties and is never rewritten: a mistake is a new item that links the old one. `archive/` is not a graveyard for finished work — that is `done/`.

### Triage, pruning and closing

- **Triage** — the same weekly pass: an idea becomes a `pending/` file with a title, a priority and a section; if it needs thought, it becomes an `ideas/` doc first and the item links it.
- **Pruning (monthly)** — a `pending/` item untouched for 2–3 months is deleted (it is in git) or moved to `archive/` with `resolution: dropped`. A `blocked_by` that no longer resolves is an error, not a warning: fix the reference or drop the item.
- **Closing** — `git mv` to `done/` and write in the body what closed it (commit, command, document). The same commit sets the originating plan or design to `status: done`; the doc moves to `archive/` at the monthly pass.
- **In progress limit** — keep `active/` to what a person is actually doing (1–2 each). Nothing enforces it; the folder is the promise.

### References across backlogs

One backlog per repo, and no copies. Work that belongs to another repo is an item there, referenced by `related: <project>#<id>` (a soft link) or `blocked_by` (a hard one). `project.md`'s `slug` is what makes the reference unambiguous. When an item moves to another repo, `git mv` it out and leave the reason in the commit.

### Converting from the queue layout

Mechanical, one commit, and the text is not rewritten:

1. State mapping: `Doing` → `active/`; `Next`, `Later`, `Someday`, blocked and parked → `pending/` (priority from the line: P1 → `p1`, and so on, P0 if the project uses it); `Done` or a ticked box → `done/` (a ticked box outside a `Done` heading is an error to fix first); `Dropped` → `archive/` with `resolution: dropped`.
2. A tracker in table form is not a backlog: move it to `docs/plans/` with `git mv` and frontmatter (`status`, `project`, `date`, `related`), and convert only the bullet items. A row that is a pointer to another repo's queue stays in the tracker.
3. With the portal's `rdk`: `rdk migrate docs/backlog/<project>.md --project <slug> --root docs/backlog`, then `git rm` the list and `rdk check --root docs/backlog` until it prints `ok`. Without it, write the files by hand in the same shape: the id can be any unused six characters from the alphabet.
4. Item files sit one level deeper than the list did: every relative link gains a `../`.
5. Report the text the conversion changed: a title split into title and body, a dropped table, a hashtag turned into a section.

The reverse move (items back to a list) is not supported: a list cannot hold ids, per-item history or parallel branches without conflicts.

### Validation

- With the portal's `rdk`: `rdk check --root docs/backlog` validates the schema, the ids, the section names, the `blocked_by` references and the root entries. It exits non-zero and names the file.
- Without it: the linter in this skill ([`scripts/docs_lint.py`](scripts/docs_lint.py)) checks the links inside the backlog and that `project.md` carries `slug`, `name` and `sections`; the rest of the schema is a review item, not a lint item.

## Both layouts

### Feedback loop

A review or a change that surfaces follow-up work creates new backlog items linking back to it. Never leave findings only inside the review.

### GitHub Issues

If the project tracks work in Issues, the backlog lives there. Keep `docs/` for the docs only and link the issue number from `related` or the item.

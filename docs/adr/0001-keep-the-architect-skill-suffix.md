---
adr: 1
status: superseded
date: 2026-11-11
project: ralvaskills
superseded_by: 2
---

# Skill names keep `-architect` as the conventions marker

## Context

An audit on 2026-11-11 counted **43 of 67 skills ending in `-architect`**, plus `hexagonal-arch` — 44 names, 66% of the catalogue. The suffix is doing four different jobs:

| Job | Examples |
|---|---|
| Structure and boundaries | `ddd-architect`, `hexagonal-arch`, `event-driven-architect` |
| Conventions for a stack | `go-architect`, `docker-architect`, `latex-architect`, `react-architect` |
| Produces an artifact | `narration-architect`, `mockup-architect`, `uru-thesis-architect` |
| Ordered decision procedure | `backend-performance-architect` |

The technical specification's [Naming Convention](../designs/2026-05-20-technical-specification.md#naming-convention) already defines `-architect` as "enforces structure/standards for a language, framework, or protocol". The family is documented, not accidental — but that definition was never stress-tested against a catalogue where it covers two thirds of the names, and it does not cover the artifact-producing members.

What a rename costs: **a skill name is the registry key.** `rsk` has no alias, redirect or deprecation mechanism (only CLI flags carry aliases). A rename deletes the old key, so `rsk install <old-name>` fails, existing `.rsk/rsk.mod` / `rsk.lock` entries stop resolving, global installs keep the stale name, and the published tarball for the old name is orphaned. The catalogue also cross-links skills by directory name, which puts the change at roughly **1,850 references across 206 files** (1,110 of them inside `skills/**`), 15 bundles in `internal/config/catalog.toml`, the registry index, the README, the landing page data, and a MAJOR version bump plus a release per renamed skill.

## Options considered

- **Keep the names, document what the suffix means** — zero breakage, one authoring rule. The cost is accepted: the suffix stops distinguishing architecture from conventions, so the description has to carry the subject.
- **Full sweep** — stack skills to `-standards`, artifact skills to `-designer` / `-author`, `backend-performance-architect` to `-playbook`, and `hexagonal-arch` to `hexagonal-architecture`. Coherent grammar, aligned with the existing `-reviewer` / `-builder` / `-planner` families. Costs 44 hard breaks and every consumer migrates by hand.
- **Phase 1 only** — the 12 names where `-architect` is factually wrong (artifact and prose skills, plus the truncated `hexagonal-arch`). Most of the clarity for a quarter of the churn, but 30 stack skills keep the ambiguous marker, which is the bulk of the confusion.
- **Bare subject names** (`go`, `docker`, `react`, `sql`) — shortest and unambiguous for stack skills; loses the role cue that the rest of the catalogue relies on, and collides conceptually with the official `docs` bundle.

## Decision

**Keep every current name.** `-architect` is this catalogue's marker for *conventions for a stack or a practice*, not for literal architecture. The subject is named by the description, not by the suffix. The suffix families and their meanings live in the specification's Naming Convention section, and that is the single place they are defined.

Renaming a published skill is a breaking change with no migration path. Any future rename is its own decision with its own ADR — never a side effect of a naming cleanup.

Scope: this covers the 43 `-architect` names and keeps them as they are. `hexagonal-arch` is the one genuine defect — the only truncated name in the catalogue, and it exists solely to justify a one-skill `-arch` row in the specification. It is tracked separately in [the backlog](../backlog/pending/n6v2wk-fix-the-truncated-hexagonal-arch-skill-name.md) and is not part of this decision.

## Consequences

- No consumer breaks, no versions bump, no releases. `go-architect`, `react-architect` and the rest keep their install keys.
- `-architect` no longer discriminates on its own: a reader cannot tell from the name whether a skill shapes structure or only lists conventions. Descriptions must lead with the subject, which the Description Quality rules already require.
- Authoring becomes one table lookup instead of a per-skill judgement call; new stack skills take `-architect` without debating `-standards` again.
- Revisit this ADR if `rsk` ships aliases or redirects (renames get cheap), or if the family passes roughly 50 of ~70 skills — at that point the marker is on so much of the catalogue that it cannot be used to filter anything.
- The audit that produced this decision is a snapshot, not an index. The mapping proposed for a full sweep is recorded below so a future reader does not have to re-derive it.

## Appendix — the rename map considered (snapshot, 2026-11-11)

Not binding. Recorded so the analysis is not repeated.

- **Keep `-architect`:** `ddd-architect`, `hexagonal-arch` (→ `hexagonal-architecture`), `event-driven-architect`, `docs-architect`.
- **→ `-standards`:** `go`, `python`, `latex`, `sql`, `protobuf`, `rest-api`, `grpc`, `connectrpc`, `mcp`, `docker`, `ci-cd`, `observability`, `secrets`, `grafana`, `react`, `nextjs`, `astro`, `hugo`, `fastapi`, `gin`, `nethttp` (→ `go-nethttp-standards`), `wails`, `remotion`, `reveal-js`, `manim`, `ui-ux`, `cli-tool`, `repo-tooling`, `ros2`.
- **→ `-designer`:** `mockup-designer`, `website-concept-designer`, `demo-script-designer`, `demo-deck-designer` (the current `demo-presentation-architect`), `explainer-designer`, `uru-thesis-defense-designer`.
- **→ `-author`:** `narration-author`, `uru-thesis-author`.
- **→ `-standards` (venue/institution):** `ml-paper-standards`, `uru-paper-standards`.
- **→ `-playbook`:** `backend-performance-playbook`.

---
status: active
project: ralvaskills
date: 2026-11-11
related:
    - adr/0001-keep-the-architect-skill-suffix.md
---

# Skill naming review — all 67 names

One row per skill, in the order of the `skills/` tree. Each was read against its own opening statement of purpose — "Design and refine slide-deck specifications", "Format a manuscript for Revista Tecnocientífica URU", "Its domain code knows nothing about the world" — plus its section headings and support-file inventory. Where the verdict turned on the body, the whole file was read.

The test is not "is the topic architecture" but **"does the name predict the artifact the skill delivers, in the role the catalogue assigns"** — `-reviewer` audits, `-builder` scaffolds, `-designer` designs an artifact, `-writer` produces prose, `-format` formats to external rules, `-architect` states the conventions and makes the structural calls for a stack.

**Act now** means it is worth breaking the name: the current one sends a reader to the wrong artifact, or the name is defective. Names that are merely less precise than ideal are marked *no* with the ideal recorded, because a skill name is the registry key and `rsk` has no alias — see [ADR 0001](../adr/0001-keep-the-architect-skill-suffix.md).

Counts: **10 act now** (MAJOR bump each), **15 ideal-but-keep**, **38 keep as-is**, plus 4 content changes that do not move a name.

## academic

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `ml-conference-paper-architect` | `ml-paper-format` | **yes** — 2.0.0 | Venue style files, anonymization, page limits, desk-reject checks. Formatting to external rules. |
| `study-artifact-builder` | keep | — | Builds the study artifact it names. |

## databases

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `sql-architect` | keep | — | Schema, keys, migrations and query shape are the architecture. |

## design

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `ddd-architect` | keep | — | Literal: bounded contexts, aggregates, boundaries. |
| `hexagonal-arch` | `hexagonal-architecture` | **yes** — 2.0.0 | The only truncated name in the catalogue, and the sole justification for the one-skill `-arch` row. |

## encoding

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `protobuf-architect` | keep | — | Package naming, field reservation, `$defs` sharing: schema design. |

## frameworks

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `astro-architect` | keep | — | Islands, rendering model, content layer. |
| `fastapi-architect` | `fastapi-standards` | no | Feature layout is the frame; Pydantic/DI/errors are standards. |
| `gin-architect` | `gin-standards` | no | Same shape as FastAPI. |
| `hugo-architect` | keep | — | Project layout and template hierarchy. |
| `manim-architect` | `manim-standards` | no | Scene, timing and asset conventions. |
| `nethttp-architect` | `go-nethttp-standards` | no | The only subject that omits its stack; the fix belongs in a standards pass, not a bundle-breaking rename. |
| `nextjs-architect` | keep | — | Server/client boundary and caching are structural. |
| `react-architect` | keep | — | Project structure and the three state layers are structural. |
| `remotion-architect` | `remotion-standards` | no | Composition and determinism conventions. |
| `reveal-js-architect` | `revealjs-standards` | no | Deck structure; also the only dotted-tool spelling in the catalogue. |
| `wails-architect` | `wails-standards` | no | Shell, proxy and cookie model — a frame with conventions inside. |

## frontend

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `demo-presentation-architect` | `demo-deck-designer` | **yes** — 2.0.0 | Outputs a slide-by-slide spec, and the `demo-` prefix currently collides with `demo-script`. |
| `demo-script-architect` | `demo-script-designer` | **yes** — 2.0.0 | Designs the presenter script; symmetric with the deck skill. |
| `explainer-architect` | `explainer-designer` | **yes** — 2.0.0 | Designs scenes, cues and a content file — an artifact, not a structure. |
| `landing-page-designer` | keep | — | The reference `-designer`: proposes page directions. |
| `mockup-architect` | `mockup-standards` | no | Opening line is "how to build mockups that can be compared, screenshotted and decided on", and it ships `pv-bar.css`/`pv-bar.js` — nearer a kit than a role. `mockup-kit` is the alternative if it ever moves. |
| `narration-architect` | `narration-writer` | **yes** — 2.0.0 | Writes voiceover and on-screen text. There is no "narration architect". |
| `ui-ux-architect` | `ui-ux-standards` | no | Design tokens, a11y floor and state discipline: conventions over structure. |
| `website-concept-architect` | `website-concept-designer` | **yes** — 2.0.0 | Runs a discovery interview and outputs 2–3 directions. |

## infra

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `ci-cd-architect` | keep | — | Pipeline taxonomy, composition and deployment strategy. |
| `docker-architect` | keep | — | Stage order and image composition are the architecture. |
| `grafana-architect` | `grafana-standards` | no | Dashboards and alert rules; no topology is designed. Weakest of the `-architect` set. |
| `observability-architect` | `observability-standards` | no | Signal naming, sampling and correlation: conventions. |
| `secrets-architect` | `secrets-standards` | no | One store plus a consumer matrix — a design, but "management" is the industry word. |

## languages

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `go-architect` | keep | — | Packages, interfaces, layout, `go.work`: the suffix in its literal sense. |
| `latex-architect` | `latex-standards` | no | Preamble order and build discipline. Expensive to rename: `ml-paper-format` and both `uru-thesis-*` skills depend on it. |
| `python-architect` | `python-standards` | no | Module layout leads, typing/testing/psycopg follow. |

## messaging

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `event-driven-architect` | keep | — | Event/command taxonomy, outbox, partitioning: structure. |

## meta

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `adhd-mode` | keep | — | A reply mode. |
| `caveman` | keep | — | A reply mode; the name is the mode. |
| `rsk-guide` | keep | — | Operator guide for one CLI. |
| `skill-builder` | keep | *content* | Name is right; it must carry the suffix table and the stack-prefix rule so new skills stop copying `-architect`. |
| `token-economy` | keep | — | Names the goal, not the tool. |

## personal

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `go-house-libs` | keep | — | Router over the shared modules; the `references/` pattern to copy when splitting. |
| `go-library-builder` | keep | — | Scaffolds library modules. |
| `pi-iteration-workflow` | keep | — | One codebase's deploy loop. |
| `uru-scientific-paper-architect` | `uru-paper-format` | **yes** — 2.0.0 | Manuscript formatting for one journal. |
| `uru-thesis-architect` | `uru-thesis-writer` | **yes** — 1.0.0 | Interview-first chapter drafter: it writes sections. |
| `uru-thesis-defense-architect` | `uru-thesis-defense-designer` | **yes** — 2.0.0 | Produces a fixed nine-section deck spec. |
| `uru-thesis-reviewer` | keep | — | Gains a clean pair: `uru-thesis-writer` / `uru-thesis-reviewer`. |

## protocols

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `connectrpc-architect` | keep | — | One `.proto` serving three protocols, interceptor chain. |
| `grpc-architect` | keep | — | Service/error/interceptor structure. |
| `mcp-architect` | keep | *content* | Server primitives, transport, stateless core: structure. Extract §9 auth and §13 testing. |
| `rest-api-architect` | keep | — | Resource design, versioning, pagination, errors. |

## quality

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `api-contract-reviewer` | keep | — | Audits contracts. |
| `backend-performance-architect` | keep | — | Its opening line is "the decision of what to change and the **design of the fix**", and it pairs with `performance-reviewer` (reviewer finds, architect designs). "Playbook" describes the ordering, not the artifact. |
| `codebase-hygiene-reviewer` | keep | — | Audits; proposes only. |
| `performance-reviewer` | keep | — | Finds and reports; the pair with the architect skill above is the catalogue's clearest role contrast. |
| `security-reviewer` | keep | — | Audits. |

## refactoring

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `code-design-refactor` | keep | — | Design-level rung of the ladder. |
| `design-patterns` | keep | — | A catalog, not a procedure. |
| `improve-codebase-architecture` | keep | — | Imperative and long, but it is the upstream skill's name; renaming would break the credit. |
| `logic-cleaner` | keep | — | Expression-level rung. |

## robotics

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `ros2-architect` | keep | — | Node topology, lifecycle, QoS. |

## tooling

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `cli-tool-architect` | keep | — | Subcommand structure and flag precedence are the architecture of a CLI. |
| `repo-tooling-architect` | keep | *content* | Repo-root layout. §8 CI restates `ci-cd-architect`; replace with a pointer. |

## workflows

| Skill | Recommended | Act now | Note |
|---|---|---|---|
| `commit-author` | keep | — | The only `-author`. `-writer` is reserved for prose skills; one word is not worth a break. |
| `docs-architect` | keep | — | The folder tree is the information architecture. |
| `feature-planner` | keep | — | Plans before building. |
| `grill-with-docs` | keep | — | Upstream name; renaming breaks the credit. |
| `tdd` | keep | — | The acronym is the term of art; the description carries the expansion. |

## Summary

| Action | Count |
|---|---|
| Rename (MAJOR bump) | 10 |
| Ideal name differs, not worth the break | 15 |
| Keep as-is | 38 |
| Content change only (name stays) | 4 — `go-architect`, `mcp-architect`, `repo-tooling-architect`, `skill-builder` |

After the 10: `-architect` 34 · `-designer` 6 · `-reviewer` 5 · `-builder` 3 · `-writer` 2 · `-format` 2 · `-author` 1 · role-less or one-off 15.

Two findings from reading the bodies rather than the descriptions:

- **The bodies call themselves "standards" or "conventions" almost universally** — `go`, `grpc`, `rest-api`, `protobuf`, `cli-tool`, `repo-tooling`, `docker`, `ui-ux` and `latex` all do. A strict self-description test would move around 26 names, not 15. The tier split is therefore a judgement about which skills make structural decisions, not a bright line — which is the same conclusion as the count: `-architect` is a category label.
- **The catalogue already uses "architect skills" as the collective noun** for this family: `security-reviewer` says it "catches the issues that architect skills already encode rules against". The family is load-bearing, not accidental.

The remaining 15 `-standards` candidates are that family fact recorded: the catalogue is mostly convention sets. Revisit them if `rsk` ships aliases, or if `latex-standards` is ever touched for another reason.

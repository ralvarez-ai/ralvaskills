---
adr: 2
status: accepted
date: 2026-11-11
project: ralvaskills
supersedes: 1
---

# Rename ten skills, and make a rename resolvable

## Context

[ADR 0001](0001-keep-the-architect-skill-suffix.md) kept every name. The review that followed ([skill naming review](../reviews/2026-11-11-skill-naming-review.md)) read all 67 skills' opening statements of purpose and separated two different problems that the count had merged:

- **`-architect` on conventions skills is defensible.** The catalogue's suffixes are roles the skill plays (`-reviewer` audits, `-builder` scaffolds), so `go-architect` reads as "the skill that makes the structural calls for Go". Ten skills fail that test — they produce a document, prose or a deck while calling themselves architects — and one (`hexagonal-arch`) is the catalogue's only truncated name.
- **The rest do not fail it.** Fourteen are less precise than `-standards` would be, but the tier split is a judgement with no bright line: the bodies of `grpc`, `rest-api`, `protobuf`, `cli-tool`, `repo-tooling`, `ui-ux` and `docker` all call themselves "conventions" or "standards" in their first paragraph.

The cost of a rename was measured, not estimated: **30 of 111 project manifests on this machine** reference at least one of the ten (26 for `hexagonal-arch` alone, four projects carrying two or more), plus 4 global symlinks and 2 bundles. About 190 references live inside this repository. `rsk` has no alias mechanism, so today a rename is a hard delete: the old key stops resolving and every consumer must be edited by hand.

## Options considered

- **Keep everything** ([ADR 0001](0001-keep-the-architect-skill-suffix.md)) — zero breakage, and the wrong roles stay in the always-loaded index. Fails on the ten names that mispredict their artifact.
- **Full sweep to `-standards`** — internally consistent, but 30+ renames across every bundle, and it trades one non-discriminating marker for another.
- **Ten renames, no mechanism** — cheapest semantics, but 30 manifests break and each is a hand edit with no upgrade path.
- **Ten renames plus a registry alias and a client migration** — the old name keeps resolving, the new name is written into `rsk.mod` / `rsk.lock` automatically, and the deprecation is invisible to an old binary. Chosen.

## Decision

Rename ten skills, and land the alias mechanism before the rename publishes.

| Now | Becomes |
|---|---|
| `narration-architect` | `narration-writer` |
| `explainer-architect` | `explainer-designer` |
| `website-concept-architect` | `website-concept-designer` |
| `demo-presentation-architect` | `demo-deck-designer` |
| `demo-script-architect` | `demo-script-designer` |
| `ml-conference-paper-architect` | `ml-paper-format` |
| `uru-scientific-paper-architect` | `uru-paper-format` |
| `uru-thesis-architect` | `uru-thesis-writer` |
| `uru-thesis-defense-architect` | `uru-thesis-defense-designer` |
| `hexagonal-arch` | `hexagonal-architecture` |

Each is a MAJOR bump; `uru-thesis-architect` goes 0.2.0 → 1.0.0.

**The mechanism.** `registry/index.json` gains a top-level `aliases` map and keeps the old entry as `deprecated`, with its versions intact, so `rsk install <old>@<pinned>` still resolves and an old binary sees "outdated" rather than "gone". `rsk` resolves the alias before lookup and rewrites `[skills]`, `pinned` and `rsk.lock` names in place, re-points the symlinks and re-renders `.rsk/SKILLS.md`; `rsk update --global` re-points the global links. Design and touchpoints: [skill rename aliases](../designs/2026-11-11-skill-rename-aliases.md).

**The other 57 names stay.** The 15 ideal-but-kept names are recorded with their ideal name in the review so they are not relitigated without new information. New skills follow one added authoring rule: **`-architect` requires a stack or protocol in the prefix** — no stack, no suffix.

## Consequences

- The rename is a breaking change for 30 projects, but it is repaired by `rsk update`, not by hand. Publication order is load-bearing: alias support ships in `rsk` one release before the ten publishes.
- Old `rsk` binaries lose access to the new versions but keep the old ones, because the deprecated entries stay in the index. No consumer is stranded.
- `hexagonal-arch` is a member of the `global` and `llm-app` bundles; those two lines, this repository's own global installs, and the spec's tree listing move with it.
- `-architect` remains the marker for 34 of 68 skills. That is a fact about a catalogue that is mostly convention sets, not a defect, and the stack-prefix rule stops it spreading to roles that do not fit.
- Revisit the 15 ideal-but-kept names if `rsk` ships aliases — they become a two-release change instead of a hand migration — or if `latex-architect` is touched for another reason, since three renamed skills depend on it.

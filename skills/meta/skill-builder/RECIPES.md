# Skill Builder — Templates & Checklists

Reference templates and the full validation checklist. Load on demand during scaffold.

## 1. Frontmatter template

```yaml
---
name: <skill-name>             # kebab-case; matches folder name
version: 1.0.0                 # always 1.0.0 unless tracking a planned feature (0.1.0)
description: <WHAT — short, comma-separated topics> — <version target if applicable>. Use when <trigger phrases>.
---
```

**Description quality (per [SPECS](../../../docs/designs/2026-05-20-technical-specification.md#description-quality)):**

- WHAT first, comma-separated, no marketing prose.
- WHEN last — trigger phrases / commands the user might invoke.
- Version-sensitive skills: include version target inline (e.g. "Go 1.27.3").
- Drop filler ("Enforces strict", "Comprehensive guide to").
- One line, YAML scalar form (no `>` or `|`).
- **Bad:** `description: Go best practices skill.` — no when.
- **Good:** `description: Go 1.27.3 architectural standards — memory-aligned structs, typed enums, interface design, goroutine safety, iterators, idiomatic errors. Use when writing, reviewing, or scaffolding Go code.`

## 2. Body skeleton

```markdown
# <Title>

<One-paragraph framing>. Targets **<version target if applicable>**.
Companion to [<related-skill>](../<path>/SKILL.md). Implementation skeletons in
[RECIPES.md](RECIPES.md). See [STACK.md](STACK.md) for pinned dependencies.

## 1. <Section name>

- **<Rule>** — <explanation>
- **<Rule>** — <explanation>

## 2. <Section name>

...

## N. Out of scope (if applicable)

- **<Concern>** → see [<other-skill>](../<path>/SKILL.md)

---

_Adapted from [<source>](<url>)._   <!-- only if attribution applies -->
```

Headings come from the body-shape preview confirmed during the interview.

## 3. STACK.md template

```markdown
# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| <dep> | <version> | <one-line> |
| ... | ... | ... |

## Notes

- <one or two architectural-opinion lines that the skill enforces>
- <inheritance: if this skill inherits a base stack, link it>

_Last reviewed: <YYYY-MM-DD>_
_Skill version at last review: <version>_
```

**Fetching live versions:**

- Go libraries → `https://proxy.golang.org/<module>/@latest`
- Python packages → `https://pypi.org/pypi/<package>/json`
- Tools / CLIs → `gh api repos/<owner>/<name>/releases/latest`

Always pin actual versions; never `latest`.

## 4. Standard cross-skill ties by category

| New skill category | Usually references |
|---|---|
| `frameworks/` (e.g. `fastapi-architect`) | matching language architect; `sql-architect`; `rest-api-architect` |
| `protocols/` (e.g. `rest-api-architect`) | usually framework-agnostic — minimal cross-references |
| `encoding/` (e.g. `protobuf-architect`) | language architect for code generation |
| `infra/` (e.g. `docker-architect`) | tooling defaults; sometimes language architects for per-language recipes |
| `quality/` (reviewer skills) | usually framework-agnostic; reference architect skills for "what they should already enforce" |
| `workflows/` | other workflow skills they chain with; pipeline references |

## 5. Naming suffix table

The suffix names the **role** the skill plays, not its topic — the subject belongs in the description. Canonical table: SPECS.md, Naming Convention.

| Suffix | Names the skill that… | Examples |
|---|---|---|
| `-architect` | states the conventions for a language, framework, protocol or tool | `go-architect`, `grpc-architect`, `docker-architect` |
| `-reviewer` | audits existing code or contracts, proposes only | `security-reviewer`, `api-contract-reviewer` |
| `-designer` | designs an artifact or spec a human reviews | `landing-page-designer`, `demo-deck-designer`, `explainer-designer` |
| `-writer` | produces prose a human reads or speaks | `narration-writer`, `uru-thesis-writer` |
| `-format` | formats an artifact to external rules | `ml-paper-format`, `uru-paper-format` |
| `-author` | produces prose or commits | `commit-author` |
| `-refactor` | restructures existing code | `code-design-refactor` |
| `-planner` | designs something before implementation | `feature-planner` |
| `-patterns` | catalogs patterns or principles for a domain | `design-patterns` |
| `-builder` | scaffolds new artifacts of a given type | `skill-builder` |
| `-guide` | documents one tool or workflow | `rsk-guide` |
| (none) | workflow tools and meta utilities that fit nothing above | `tdd`, `caveman` |

**`-architect` requires a stack or protocol in the prefix.** The prefix names a language, framework, protocol or tool — `go`, `react`, `grpc`, `docker`. A skill whose subject is a practice, a document, a paper or a deck is not an architect: **no stack, no suffix**; take the producing role (`-designer`, `-writer`, `-format`) instead. This is the rule that keeps the `-architect` count from growing on roles that do not fit — ten names moved off it in [ADR 0002](../../../docs/adr/0002-rename-ten-skills-with-registry-aliases.md).

**A published name is permanent**, because it is the registry key. A rename needs its own ADR and one `[aliases]` entry in `registry/aliases.toml`, in the same commit as the directory rename; the old entry stays in the index as `deprecated` so it keeps resolving.

**The name is checked, not trusted.** `internal/skill/names.go` enforces kebab-case, at most 64 characters, no doubled hyphen, `directory == frontmatter name`, uniqueness, and that the name is not retired — including as a whole hyphen-separated component, so `hexagonal-arch-v2` is rejected too. `go test ./internal/skill/` (`task skills:check`) runs it; `cmd/generate-registry` refuses to publish a violation.

If the topic doesn't fit any suffix, propose a new one and update SPECS — don't shoehorn.

## 6. SPECS.md atomic update list

Every new skill needs 4–6 edits to SPECS.md, applied in one pass:

1. **Folder Structure tree** — add `<category>/<name>/` entry with `# ✅ exists (v1.0.0 — short note)`.
2. **Roadmap row** in the matching category section — `| <name> | ✅ | v1.0.0 — one-line scope summary |`.
3. **Global bundle row** (if applicable) — `| <name> | <category> | <description> |`.
4. **Stack-specific bundle rows** (gin, nethttp, fastapi, go-grpc, etc.) — if applicable.
5. **Naming convention table** — only if introducing a new suffix.
6. **STACK.md needed table** — only if introducing a new category.

**Verify:** grep SPECS for the skill name; every mention should make sense in context. No `📋 planned` rows remain for a skill that now exists.

## 7. Validation checklist

Before declaring scaffold done:

- [ ] Folder created at `skills/<category>/<name>/`
- [ ] `SKILL.md` present with valid YAML frontmatter (`name`, `version`, `description`)
- [ ] **Description is one line** in YAML scalar form (no `>` or `|`); has both WHAT and WHEN; version target named if applicable
- [ ] Name in frontmatter matches folder name
- [ ] Naming suffix valid per SPECS convention table
- [ ] `-architect` used only with a stack or protocol prefix; a practice, document or deck takes `-designer`, `-writer` or `-format` instead
- [ ] Name is kebab-case, at most 64 characters, no doubled hyphen, and not a retired name or containing one as a whole hyphen-separated component
- [ ] `task skills:check` (`go test ./internal/skill/`) passes
- [ ] If this renames a published skill: an ADR and one `[aliases]` entry in `registry/aliases.toml`, in the same commit as the directory
- [ ] No nested `SKILL.md` inside (single discovery point)
- [ ] All `.md` files in the folder use `UPPER_SNAKE_CASE` (except `SKILL.md` and `STACK.md` themselves)
- [ ] `STACK.md` present if category requires (per SPECS table) and not exempt; or exemption noted in roadmap
- [ ] `RECIPES.md` created if the skill will house full reference implementations; each long code block in SKILL.md replaced with a one-line pointer
- [ ] No code block in `SKILL.md` exceeds ~15 lines unless truly load-bearing
- [ ] All cross-skill links resolve (relative paths point to real files)
- [ ] No emojis in `.md` files (per project convention)
- [ ] SPECS.md folder structure entry added with status marker
- [ ] SPECS.md roadmap row updated to ✅ with one-line summary
- [ ] SPECS.md bundle rows updated if skill is in any bundle
- [ ] Attribution footer added if skill is adapted from an external source
- [ ] No stale `📋 planned` rows remain for this skill in SPECS

If any item fails, fix before handing off.

## 8. STACK.md exemption rules

**Exempt from STACK.md:**

- `workflows/` — language-agnostic process skills.
- `refactoring/` — code-cleanup ladder is language-agnostic.
- `design/` — architectural patterns are version-agnostic.
- `quality/` — security/perf principles are stable (override only if framework-specific behavior).
- `meta/` — ralvaskills-ecosystem skills with no library stack (e.g. `skill-builder`, `rsk-guide`, `caveman`).

Mark the exemption explicitly in the roadmap notes when applicable.

## 9. Attribution footer

If the skill is adapted from an external source:

```markdown
---

_Adapted from [<source>](<url>)._
```

Verbatim copies require strict attribution. Evolved skills should still credit the origin — see `grill-with-docs`, `tdd`, `improve-codebase-architecture` for examples.

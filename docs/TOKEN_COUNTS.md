# SKILL.md Token Estimates

> Auto-generated. **Do not edit by hand.** Run `task tokens` to refresh.
>
> Estimate: ~4 bytes/token for bodies, ~3 bytes/token for descriptions (Claude tokenizer). Actual range ±15%.

_Last updated: 2026-11-02 · 57 skills · 15 bundles_

## Load model

Description tokens hit **every turn** for installed skills. Body tokens are paid only when a skill is actually invoked. Side files are never auto-loaded.

| What | When loaded | Estimated tokens |
|---|---|---:|
| All `SKILL.md` bodies | Only when invoked | ~142213 |
| All side files (`STACK` + `RECIPES` + topic files) | On-demand only | ~124734 |
| All `description:` fields (every skill) | Every turn, if all installed | ~6610 |

## Session profiles

Budgets: **global ≤ 1500** desc tokens (paid by every project), **session ≤ 2200** desc tokens (global + one project bundle). Status: `ok` < 90% · `warn` ≥ 90% · `OVER` ≥ 100%.

**Global baseline** (`rsk install global --global`): 15 skills · ~1510 / 1500 desc tokens every turn (`OVER`) · ~28343 body tokens if all invoked.

Per-project bundle additions on top of the global baseline:

| Bundle | Skills | + Desc tkns / turn | Session total | Budget | Status |
|---|---:|---:|---:|---:|---:|
| `docs` *(+5 external)* | 0 | +0 | ~1510 | 2200 | `ok` |
| `observability` | 2 | +151 | ~1661 | 2200 | `ok` |
| `event-driven` | 2 | +167 | ~1677 | 2200 | `ok` |
| `python-cli` | 3 | +261 | ~1771 | 2200 | `ok` |
| `code-review` | 3 | +269 | ~1779 | 2200 | `ok` |
| `ros2` | 3 | +273 | ~1783 | 2200 | `ok` |
| `design` *(+1 external)* | 3 | +313 | ~1823 | 2200 | `ok` |
| `go-cli` | 3 | +335 | ~1845 | 2200 | `ok` |
| `llm-app` *(+2 missing)* | 4 | +389 | ~1899 | 2200 | `ok` |
| `python-grpc` | 5 | +448 | ~1958 | 2200 | `ok` |
| `fastapi` | 5 | +480 | ~1990 | 2200 | `warn` |
| `go-grpc` | 5 | +522 | ~2032 | 2200 | `warn` |
| `gin` | 5 | +542 | ~2052 | 2200 | `warn` |
| `nethttp` | 5 | +554 | ~2064 | 2200 | `warn` |

## Personal / unbundled skills

These skills are not part of any bundle. They are installed individually with `rsk install <name> --include personal` and their desc tokens are only paid when explicitly installed.

| Skill | Category | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---|---:|---:|---:|
| `mcp-architect` | protocols | ~8119 | ~366 | ~10017 |
| `go-house-libs` | personal | ~2258 | ~265 | ~1523 |
| `go-library-builder` | personal | ~4040 | ~166 | ~3353 |
| `connectrpc-architect` | protocols | ~3465 | ~153 | ~2733 |
| `study-artifact-builder` | academic | ~2892 | ~149 | ~6369 |
| `ci-cd-architect` | infra | ~3247 | ~145 | ~3938 |
| `docs-architect` | workflows | ~1720 | ~144 | ~5214 |
| `uru-scientific-paper-architect` | personal | ~2233 | ~141 | ~2021 |
| `secrets-architect` | infra | ~2612 | ~132 | ~0 |
| `uru-thesis-defense-architect` | personal | ~1798 | ~132 | ~1525 |
| `astro-architect` | frameworks | ~3252 | ~132 | ~3941 |
| `website-concept-architect` | frontend | ~3083 | ~122 | ~0 |
| `ml-conference-paper-architect` | academic | ~2617 | ~121 | ~1762 |
| `uru-thesis-reviewer` | personal | ~3358 | ~119 | ~7587 |
| `latex-architect` | languages | ~3099 | ~116 | ~2180 |
| `uru-thesis-architect` | personal | ~1948 | ~114 | ~6629 |
| `pi-iteration-workflow` | personal | ~916 | ~108 | ~364 |
| `hugo-architect` | frameworks | ~2817 | ~107 | ~1407 |
| `demo-presentation-architect` | personal | ~2356 | ~101 | ~5548 |
| `demo-script-architect` | personal | ~1237 | ~96 | ~0 |
| `reveal-js-architect` | frameworks | ~2745 | ~93 | ~1203 |
| `codebase-hygiene-reviewer` | quality | ~1275 | ~90 | ~1032 |

## By category

| Category | Skills | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---:|---:|---:|---:|
| frameworks | 8 | ~20527 | ~890 | ~15621 |
| personal | 9 | ~20144 | ~1242 | ~28550 |
| protocols | 4 | ~16670 | ~727 | ~18751 |
| languages | 3 | ~13502 | ~348 | ~11222 |
| infra | 5 | ~12065 | ~517 | ~8889 |
| quality | 4 | ~7993 | ~359 | ~2183 |
| tooling | 2 | ~7639 | ~229 | ~4661 |
| workflows | 5 | ~6899 | ~512 | ~7957 |
| refactoring | 4 | ~6803 | ~408 | ~6101 |
| academic | 2 | ~5509 | ~270 | ~8131 |
| frontend | 2 | ~5313 | ~225 | ~1853 |
| meta | 3 | ~5121 | ~342 | ~1718 |
| design | 2 | ~4320 | ~179 | ~1519 |
| messaging | 1 | ~2878 | ~83 | ~2662 |
| robotics | 1 | ~2717 | ~105 | ~2321 |
| databases | 1 | ~2172 | ~90 | ~1359 |
| encoding | 1 | ~1941 | ~84 | ~1236 |

## By bundle

Catalog from `internal/config/catalog.toml`. `External` are `source = "official"` skills (e.g. `docx`, `xlsx`) that live outside this repo by design. `Missing` are `source = "local"` skills the bundle references that aren't in the local `skills/` tree.

| Bundle | Local | External | Missing | ~Body tkns (local) | ~Desc tkns (local) |
|---|---:|---:|---:|---:|---:|
| `global` | 15 | 0 | 0 | ~28343 | ~1510 |
| `docs` | 0 | 5 | 0 | ~0 | ~0 |
| `design` | 3 | 1 | 0 | ~8243 | ~313 |
| `go-grpc` | 5 | 0 | 0 | ~15527 | ~522 |
| `gin` | 5 | 0 | 0 | ~15539 | ~542 |
| `nethttp` | 5 | 0 | 0 | ~15702 | ~554 |
| `go-cli` | 3 | 0 | 0 | ~11769 | ~335 |
| `fastapi` | 5 | 0 | 0 | ~12060 | ~480 |
| `llm-app` | 4 | 0 | 2 | ~9704 | ~389 |
| `ros2` | 3 | 0 | 0 | ~8131 | ~273 |
| `python-grpc` | 5 | 0 | 0 | ~12010 | ~448 |
| `python-cli` | 3 | 0 | 0 | ~8252 | ~261 |
| `event-driven` | 2 | 0 | 0 | ~4819 | ~167 |
| `observability` | 2 | 0 | 0 | ~4235 | ~151 |
| `code-review` | 3 | 0 | 0 | ~6718 | ~269 |

## Per skill

| # | Skill | Category | Body bytes | ~Body tkns | ~Desc tkns | ~Stack tkns | ~Recipes tkns | ~Topic tkns |
|---|---|---|---:|---:|---:|---:|---:|---:|
| 1 | `mcp-architect` | protocols | 32371 | ~8119 | ~366 | ~1521 | ~8496 | ~0 |
| 2 | `go-architect` | languages | 27535 | ~6960 | ~153 | ~473 | ~2685 | ~5441 |
| 3 | `repo-tooling-architect` | tooling | 18156 | ~4801 | ~136 | ~685 | ~0 | ~1949 |
| 4 | `go-library-builder` | personal | 16158 | ~4040 | ~166 | ~628 | ~2725 | ~0 |
| 5 | `connectrpc-architect` | protocols | 13859 | ~3465 | ~153 | ~777 | ~1956 | ~0 |
| 6 | `python-architect` | languages | 13540 | ~3443 | ~79 | ~443 | ~0 | ~0 |
| 7 | `uru-thesis-reviewer` | personal | 13430 | ~3358 | ~119 | ~0 | ~0 | ~7587 |
| 8 | `astro-architect` | frameworks | 12582 | ~3252 | ~132 | ~604 | ~3337 | ~0 |
| 9 | `ci-cd-architect` | infra | 12985 | ~3247 | ~145 | ~624 | ~3314 | ~0 |
| 10 | `latex-architect` | languages | 12396 | ~3099 | ~116 | ~339 | ~1841 | ~0 |
| 11 | `website-concept-architect` | frontend | 11524 | ~3083 | ~122 | ~0 | ~0 | ~0 |
| 12 | `react-architect` | frameworks | 12104 | ~3026 | ~102 | ~465 | ~1426 | ~0 |
| 13 | `nextjs-architect` | frameworks | 11948 | ~2987 | ~108 | ~454 | ~1824 | ~0 |
| 14 | `rsk-guide` | meta | 11028 | ~2941 | ~199 | ~0 | ~0 | ~0 |
| 15 | `study-artifact-builder` | academic | 11566 | ~2892 | ~149 | ~0 | ~2272 | ~4097 |
| 16 | `event-driven-architect` | messaging | 11494 | ~2878 | ~83 | ~514 | ~2148 | ~0 |
| 17 | `cli-tool-architect` | tooling | 11334 | ~2838 | ~93 | ~591 | ~1436 | ~0 |
| 18 | `hugo-architect` | frameworks | 11266 | ~2817 | ~107 | ~434 | ~973 | ~0 |
| 19 | `reveal-js-architect` | frameworks | 10980 | ~2745 | ~93 | ~147 | ~1056 | ~0 |
| 20 | `ros2-architect` | robotics | 10867 | ~2717 | ~105 | ~607 | ~1714 | ~0 |
| 21 | `ml-conference-paper-architect` | academic | 10468 | ~2617 | ~121 | ~457 | ~1305 | ~0 |
| 22 | `secrets-architect` | infra | 10167 | ~2612 | ~132 | ~0 | ~0 | ~0 |
| 23 | `rest-api-architect` | protocols | 10412 | ~2603 | ~102 | ~506 | ~0 | ~4141 |
| 24 | `observability-architect` | infra | 10238 | ~2577 | ~74 | ~546 | ~874 | ~0 |
| 25 | `api-contract-reviewer` | quality | 9798 | ~2511 | ~108 | ~0 | ~0 | ~0 |
| 26 | `grpc-architect` | protocols | 9932 | ~2483 | ~106 | ~521 | ~833 | ~0 |
| 27 | `code-design-refactor` | refactoring | 9693 | ~2456 | ~96 | ~0 | ~0 | ~0 |
| 28 | `hexagonal-arch` | design | 9537 | ~2419 | ~101 | ~0 | ~768 | ~0 |
| 29 | `feature-planner` | workflows | 9202 | ~2369 | ~114 | ~0 | ~0 | ~0 |
| 30 | `demo-presentation-architect` | personal | 9421 | ~2356 | ~101 | ~0 | ~1776 | ~3772 |
| 31 | `go-house-libs` | personal | 9031 | ~2258 | ~265 | ~1523 | ~0 | ~0 |
| 32 | `uru-scientific-paper-architect` | personal | 8930 | ~2233 | ~141 | ~0 | ~2021 | ~0 |
| 33 | `ui-ux-architect` | frontend | 8920 | ~2230 | ~103 | ~466 | ~1387 | ~0 |
| 34 | `sql-architect` | databases | 8686 | ~2172 | ~90 | ~295 | ~0 | ~1064 |
| 35 | `security-reviewer` | quality | 8474 | ~2119 | ~77 | ~0 | ~539 | ~0 |
| 36 | `performance-reviewer` | quality | 8350 | ~2088 | ~84 | ~0 | ~612 | ~0 |
| 37 | `nethttp-architect` | frameworks | 7765 | ~1996 | ~120 | ~398 | ~1508 | ~0 |
| 38 | `docker-architect` | infra | 7854 | ~1971 | ~89 | ~333 | ~718 | ~924 |
| 39 | `uru-thesis-architect` | personal | 7789 | ~1948 | ~114 | ~0 | ~1679 | ~4950 |
| 40 | `protobuf-architect` | encoding | 7763 | ~1941 | ~84 | ~474 | ~762 | ~0 |
| 41 | `improve-codebase-architecture` | refactoring | 7656 | ~1914 | ~97 | ~0 | ~0 | ~3971 |
| 42 | `ddd-architect` | design | 7602 | ~1901 | ~78 | ~0 | ~0 | ~751 |
| 43 | `fastapi-architect` | frameworks | 7481 | ~1871 | ~120 | ~420 | ~1016 | ~0 |
| 44 | `gin-architect` | frameworks | 7136 | ~1833 | ~108 | ~413 | ~1146 | ~0 |
| 45 | `design-patterns` | refactoring | 7299 | ~1825 | ~122 | ~0 | ~0 | ~2130 |
| 46 | `uru-thesis-defense-architect` | personal | 7191 | ~1798 | ~132 | ~0 | ~1525 | ~0 |
| 47 | `skill-builder` | meta | 7030 | ~1758 | ~72 | ~0 | ~1718 | ~0 |
| 48 | `docs-architect` | workflows | 6878 | ~1720 | ~144 | ~0 | ~0 | ~5214 |
| 49 | `grafana-architect` | infra | 6632 | ~1658 | ~77 | ~407 | ~1149 | ~0 |
| 50 | `codebase-hygiene-reviewer` | quality | 5100 | ~1275 | ~90 | ~0 | ~1032 | ~0 |
| 51 | `demo-script-architect` | personal | 4828 | ~1237 | ~96 | ~0 | ~0 | ~0 |
| 52 | `tdd` | workflows | 4279 | ~1119 | ~69 | ~0 | ~0 | ~1350 |
| 53 | `pi-iteration-workflow` | personal | 3662 | ~916 | ~108 | ~0 | ~364 | ~0 |
| 54 | `grill-with-docs` | workflows | 3376 | ~889 | ~104 | ~0 | ~0 | ~1393 |
| 55 | `commit-author` | workflows | 3208 | ~802 | ~81 | ~0 | ~0 | ~0 |
| 56 | `logic-cleaner` | refactoring | 2432 | ~608 | ~93 | ~0 | ~0 | ~0 |
| 57 | `caveman` | meta | 1687 | ~422 | ~71 | ~0 | ~0 | ~0 |

**Totals:** 563030 body bytes · ~142213 body tokens · ~6610 desc tokens · ~124734 side tokens

## Topic files

Side files in a skill directory other than `STACK.md` and `RECIPES.md` — typically topic-specific reference docs referenced by name from `SKILL.md` (e.g. `PAYLOADS.md`, `AUTH_PATTERNS.md`). These contribute to the `~Topic tkns` column above and are loaded only when the parent skill chooses to read them.

- `uru-thesis-reviewer` (~7587 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
  - `TEMPLATES.md` — ~1830 tokens
  - `CHECKS.md` — ~807 tokens
- `go-architect` (~5441 tokens total):
  - `LAYOUTS.md` — ~5441 tokens
- `docs-architect` (~5214 tokens total):
  - `BACKLOG.md` — ~2339 tokens
  - `MIGRATION.md` — ~1076 tokens
  - `TEMPLATES.md` — ~1057 tokens
  - `FRONTMATTER.md` — ~742 tokens
- `uru-thesis-architect` (~4950 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
- `rest-api-architect` (~4141 tokens total):
  - `CONCURRENCY.md` — ~1056 tokens
  - `IDEMPOTENCY.md` — ~869 tokens
  - `PAYLOADS.md` — ~787 tokens
  - `AUTH_PATTERNS.md` — ~767 tokens
  - `STATUS_CODES.md` — ~662 tokens
- `study-artifact-builder` (~4097 tokens total):
  - `PIPELINES.md` — ~2239 tokens
  - `DATASETS.md` — ~1858 tokens
- `improve-codebase-architecture` (~3971 tokens total):
  - `HTML_REPORT.md` — ~1698 tokens
  - `LANGUAGE.md` — ~951 tokens
  - `INTERFACE_DESIGN.md` — ~681 tokens
  - `DEEPENING.md` — ~641 tokens
- `demo-presentation-architect` (~3772 tokens total):
  - `LAYOUTS.md` — ~3772 tokens
- `design-patterns` (~2130 tokens total):
  - `PATTERNS.md` — ~2130 tokens
- `repo-tooling-architect` (~1949 tokens total):
  - `TASKFILE.md` — ~1949 tokens
- `grill-with-docs` (~1393 tokens total):
  - `ADR_FORMAT.md` — ~703 tokens
  - `CONTEXT_FORMAT.md` — ~690 tokens
- `tdd` (~1350 tokens total):
  - `TESTS.md` — ~410 tokens
  - `MOCKING.md` — ~370 tokens
  - `DEEP_MODULES.md` — ~310 tokens
  - `INTERFACE_DESIGN.md` — ~163 tokens
  - `REFACTORING.md` — ~97 tokens
- `sql-architect` (~1064 tokens total):
  - `ENGINES.md` — ~563 tokens
  - `POSTGRES.md` — ~501 tokens
- `docker-architect` (~924 tokens total):
  - `CHECKLIST.md` — ~924 tokens
- `ddd-architect` (~751 tokens total):
  - `PATTERNS.md` — ~751 tokens

## Skills to consider trimming

Body > 2500 tokens — consider moving examples to `RECIPES.md` or topic files:

- `mcp-architect` (~8119 body tokens)
- `go-architect` (~6960 body tokens)
- `repo-tooling-architect` (~4801 body tokens)
- `go-library-builder` (~4040 body tokens)
- `connectrpc-architect` (~3465 body tokens)
- `python-architect` (~3443 body tokens)
- `uru-thesis-reviewer` (~3358 body tokens)
- `astro-architect` (~3252 body tokens)
- `ci-cd-architect` (~3247 body tokens)
- `latex-architect` (~3099 body tokens)
- `website-concept-architect` (~3083 body tokens)
- `react-architect` (~3026 body tokens)
- `nextjs-architect` (~2987 body tokens)
- `rsk-guide` (~2941 body tokens)
- `study-artifact-builder` (~2892 body tokens)
- `event-driven-architect` (~2878 body tokens)
- `cli-tool-architect` (~2838 body tokens)
- `hugo-architect` (~2817 body tokens)
- `reveal-js-architect` (~2745 body tokens)
- `ros2-architect` (~2717 body tokens)
- `ml-conference-paper-architect` (~2617 body tokens)
- `secrets-architect` (~2612 body tokens)
- `rest-api-architect` (~2603 body tokens)
- `observability-architect` (~2577 body tokens)
- `api-contract-reviewer` (~2511 body tokens)

Heaviest 10 descriptions — each desc token is paid every turn for any session that installs the skill:

- `mcp-architect` (~366 desc tokens)
- `go-house-libs` (~265 desc tokens)
- `rsk-guide` (~199 desc tokens)
- `go-library-builder` (~166 desc tokens)
- `go-architect` (~153 desc tokens)
- `connectrpc-architect` (~153 desc tokens)
- `study-artifact-builder` (~149 desc tokens)
- `ci-cd-architect` (~145 desc tokens)
- `docs-architect` (~144 desc tokens)
- `uru-scientific-paper-architect` (~141 desc tokens)

## Notes

- **Body tokens** cost only when the skill is invoked in a session.
- **Description tokens** cost every turn for any session that has the skill installed. The cost depends on which bundles are installed, not on the corpus total.
- Side files (`STACK.md`, `RECIPES.md`, topic files) are never auto-loaded; they cost 0 per turn.
- For exact counts: run each file through [`tiktoken`](https://github.com/openai/tiktoken) with `cl100k_base`.

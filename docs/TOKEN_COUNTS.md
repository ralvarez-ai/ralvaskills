# SKILL.md Token Estimates

> Auto-generated. **Do not edit by hand.** Run `task tokens` to refresh.
>
> Estimate: ~4 bytes/token for bodies, ~3 bytes/token for descriptions (Claude tokenizer). Actual range ±15%.

_Last updated: 2026-11-12 · 67 skills · 15 bundles_

## Load model

Description tokens hit **every turn** for installed skills. Body tokens are paid only when a skill is actually invoked. Side files are never auto-loaded.

| What | When loaded | Estimated tokens |
|---|---|---:|
| All `SKILL.md` bodies | Only when invoked | ~158591 |
| All side files (`STACK` + `RECIPES` + topic files) | On-demand only | ~191560 |
| All `description:` fields (every skill) | Every turn, if all installed | ~7637 |

## Session profiles

Budgets: **global ≤ 1500** desc tokens (paid by every project), **session ≤ 2200** desc tokens (global + one project bundle). Status: `ok` < 90% · `warn` ≥ 90% · `OVER` ≥ 100%.

**Global baseline** (`rsk install global --global`): 17 skills · ~1226 / 1500 desc tokens every turn (`ok`) · ~31114 body tokens if all invoked.

Per-project bundle additions on top of the global baseline:

| Bundle | Skills | + Desc tkns / turn | Session total | Budget | Status |
|---|---:|---:|---:|---:|---:|
| `docs` *(+5 external)* | 0 | +0 | ~1226 | 2200 | `ok` |
| `observability` | 2 | +151 | ~1377 | 2200 | `ok` |
| `event-driven` | 2 | +167 | ~1393 | 2200 | `ok` |
| `python-cli` | 3 | +261 | ~1487 | 2200 | `ok` |
| `code-review` | 3 | +269 | ~1495 | 2200 | `ok` |
| `ros2` | 3 | +273 | ~1499 | 2200 | `ok` |
| `design` *(+1 external)* | 3 | +313 | ~1539 | 2200 | `ok` |
| `llm-app` *(+2 missing)* | 4 | +359 | ~1585 | 2200 | `ok` |
| `go-cli` | 3 | +393 | ~1619 | 2200 | `ok` |
| `python-grpc` | 5 | +448 | ~1674 | 2200 | `ok` |
| `fastapi` | 5 | +480 | ~1706 | 2200 | `ok` |
| `go-grpc` | 5 | +580 | ~1806 | 2200 | `ok` |
| `gin` | 5 | +601 | ~1827 | 2200 | `ok` |
| `nethttp` | 5 | +613 | ~1839 | 2200 | `ok` |

## Personal / unbundled skills

These skills are not part of any bundle. They are installed individually with `rsk install <name> --include personal` and their desc tokens are only paid when explicitly installed.

| Skill | Category | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---|---:|---:|---:|
| `mcp-architect` | protocols | ~7795 | ~334 | ~11014 |
| `go-house-libs` | personal | ~2197 | ~265 | ~2071 |
| `backend-performance-architect` | quality | ~3413 | ~204 | ~3498 |
| `manim-architect` | frameworks | ~1776 | ~198 | ~4433 |
| `explainer-designer` | frontend | ~1704 | ~187 | ~1936 |
| `go-library-builder` | personal | ~4026 | ~167 | ~4058 |
| `landing-page-designer` | frontend | ~3057 | ~166 | ~1134 |
| `connectrpc-architect` | protocols | ~3046 | ~153 | ~3750 |
| `study-artifact-builder` | academic | ~2546 | ~149 | ~7193 |
| `ci-cd-architect` | infra | ~3085 | ~145 | ~4673 |
| `docs-architect` | workflows | ~1643 | ~144 | ~5829 |
| `uru-paper-format` | personal | ~2114 | ~141 | ~2567 |
| `remotion-architect` | frameworks | ~3218 | ~139 | ~4485 |
| `mockup-architect` | frontend | ~885 | ~136 | ~474 |
| `narration-writer` | frontend | ~2045 | ~134 | ~3027 |
| `uru-thesis-defense-designer` | personal | ~1777 | ~132 | ~2025 |
| `astro-architect` | frameworks | ~2889 | ~132 | ~4461 |
| `secrets-architect` | infra | ~2555 | ~132 | ~558 |
| `website-concept-designer` | frontend | ~2294 | ~122 | ~1198 |
| `ml-paper-format` | academic | ~2551 | ~121 | ~2207 |
| `uru-thesis-reviewer` | personal | ~3299 | ~119 | ~8215 |
| `wails-architect` | frameworks | ~2486 | ~118 | ~4593 |
| `latex-architect` | languages | ~3026 | ~116 | ~2719 |
| `uru-thesis-writer` | personal | ~1660 | ~114 | ~7162 |
| `pi-iteration-workflow` | personal | ~878 | ~108 | ~781 |
| `hugo-architect` | frameworks | ~2720 | ~107 | ~1836 |
| `demo-deck-designer` | frontend | ~2231 | ~101 | ~6050 |
| `demo-script-designer` | frontend | ~939 | ~96 | ~794 |
| `reveal-js-architect` | frameworks | ~2462 | ~93 | ~1729 |
| `codebase-hygiene-reviewer` | quality | ~1187 | ~90 | ~1460 |

## By category

| Category | Skills | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---:|---:|---:|---:|
| frameworks | 11 | ~26995 | ~1347 | ~33694 |
| personal | 7 | ~15951 | ~1046 | ~26879 |
| protocols | 4 | ~15896 | ~695 | ~22027 |
| frontend | 8 | ~15369 | ~1045 | ~16939 |
| languages | 3 | ~13482 | ~406 | ~16535 |
| infra | 5 | ~11698 | ~517 | ~11926 |
| quality | 5 | ~11193 | ~563 | ~7742 |
| meta | 5 | ~9752 | ~362 | ~8387 |
| refactoring | 4 | ~6628 | ~278 | ~8125 |
| workflows | 5 | ~6572 | ~423 | ~10578 |
| tooling | 2 | ~6155 | ~174 | ~7172 |
| academic | 2 | ~5097 | ~270 | ~9400 |
| design | 2 | ~4365 | ~149 | ~2471 |
| messaging | 1 | ~2770 | ~83 | ~3179 |
| robotics | 1 | ~2701 | ~105 | ~2816 |
| databases | 1 | ~2149 | ~90 | ~1823 |
| encoding | 1 | ~1818 | ~84 | ~1867 |

## By bundle

Catalog from `internal/config/catalog.toml`. `External` are `source = "official"` skills (e.g. `docx`, `xlsx`) that live outside this repo by design. `Missing` are `source = "local"` skills the bundle references that aren't in the local `skills/` tree.

| Bundle | Local | External | Missing | ~Body tkns (local) | ~Desc tkns (local) |
|---|---:|---:|---:|---:|---:|
| `global` | 17 | 0 | 0 | ~31114 | ~1226 |
| `docs` | 0 | 5 | 0 | ~0 | ~0 |
| `design` | 3 | 1 | 0 | ~8155 | ~313 |
| `go-grpc` | 5 | 0 | 0 | ~15526 | ~580 |
| `gin` | 5 | 0 | 0 | ~15538 | ~601 |
| `nethttp` | 5 | 0 | 0 | ~15760 | ~613 |
| `go-cli` | 3 | 0 | 0 | ~11878 | ~393 |
| `fastapi` | 5 | 0 | 0 | ~11672 | ~480 |
| `llm-app` | 4 | 0 | 2 | ~9350 | ~359 |
| `ros2` | 3 | 0 | 0 | ~7763 | ~273 |
| `python-grpc` | 5 | 0 | 0 | ~11494 | ~448 |
| `python-cli` | 3 | 0 | 0 | ~7846 | ~261 |
| `event-driven` | 2 | 0 | 0 | ~4588 | ~167 |
| `observability` | 2 | 0 | 0 | ~4208 | ~151 |
| `code-review` | 3 | 0 | 0 | ~6593 | ~269 |

## Per skill

| # | Skill | Category | Body bytes | ~Body tkns | ~Desc tkns | ~Stack tkns | ~Recipes tkns | ~Topic tkns |
|---|---|---|---:|---:|---:|---:|---:|---:|
| 1 | `mcp-architect` | protocols | 31177 | ~7795 | ~334 | ~1521 | ~8783 | ~710 |
| 2 | `go-architect` | languages | 28974 | ~7244 | ~211 | ~473 | ~4429 | ~7535 |
| 3 | `go-library-builder` | personal | 16102 | ~4026 | ~167 | ~628 | ~2725 | ~705 |
| 4 | `backend-performance-architect` | quality | 13650 | ~3413 | ~204 | ~0 | ~2747 | ~751 |
| 5 | `repo-tooling-architect` | tooling | 13230 | ~3371 | ~81 | ~685 | ~1334 | ~2679 |
| 6 | `uru-thesis-reviewer` | personal | 13195 | ~3299 | ~119 | ~0 | ~0 | ~8215 |
| 7 | `remotion-architect` | frameworks | 12870 | ~3218 | ~139 | ~316 | ~2588 | ~1581 |
| 8 | `python-architect` | languages | 12848 | ~3212 | ~79 | ~443 | ~203 | ~733 |
| 9 | `ci-cd-architect` | infra | 12339 | ~3085 | ~145 | ~624 | ~3520 | ~529 |
| 10 | `token-economy` | meta | 12249 | ~3063 | ~73 | ~0 | ~2033 | ~2171 |
| 11 | `landing-page-designer` | frontend | 12177 | ~3057 | ~166 | ~0 | ~0 | ~1134 |
| 12 | `connectrpc-architect` | protocols | 12182 | ~3046 | ~153 | ~777 | ~2421 | ~552 |
| 13 | `latex-architect` | languages | 12101 | ~3026 | ~116 | ~332 | ~1841 | ~546 |
| 14 | `nextjs-architect` | frameworks | 11890 | ~2973 | ~108 | ~367 | ~1824 | ~555 |
| 15 | `react-architect` | frameworks | 11872 | ~2968 | ~102 | ~415 | ~1426 | ~535 |
| 16 | `rsk-guide` | meta | 11136 | ~2967 | ~71 | ~0 | ~0 | ~739 |
| 17 | `astro-architect` | frameworks | 11496 | ~2889 | ~132 | ~604 | ~3298 | ~559 |
| 18 | `cli-tool-architect` | tooling | 11119 | ~2784 | ~93 | ~515 | ~1436 | ~523 |
| 19 | `event-driven-architect` | messaging | 11059 | ~2770 | ~83 | ~514 | ~2148 | ~517 |
| 20 | `hugo-architect` | frameworks | 10877 | ~2720 | ~107 | ~403 | ~973 | ~460 |
| 21 | `ros2-architect` | robotics | 10802 | ~2701 | ~105 | ~579 | ~1714 | ~523 |
| 22 | `rest-api-architect` | protocols | 10358 | ~2590 | ~102 | ~506 | ~0 | ~4720 |
| 23 | `secrets-architect` | infra | 9940 | ~2555 | ~132 | ~0 | ~0 | ~558 |
| 24 | `ml-paper-format` | academic | 10204 | ~2551 | ~121 | ~395 | ~1301 | ~511 |
| 25 | `study-artifact-builder` | academic | 10184 | ~2546 | ~149 | ~0 | ~2611 | ~4582 |
| 26 | `observability-architect` | infra | 10025 | ~2523 | ~74 | ~546 | ~874 | ~498 |
| 27 | `api-contract-reviewer` | quality | 9798 | ~2511 | ~108 | ~0 | ~0 | ~616 |
| 28 | `wails-architect` | frameworks | 9941 | ~2486 | ~118 | ~335 | ~3508 | ~750 |
| 29 | `grpc-architect` | protocols | 9860 | ~2465 | ~106 | ~521 | ~864 | ~652 |
| 30 | `reveal-js-architect` | frameworks | 9847 | ~2462 | ~93 | ~147 | ~1056 | ~526 |
| 31 | `hexagonal-architecture` | design | 9532 | ~2417 | ~71 | ~0 | ~768 | ~420 |
| 32 | `code-design-refactor` | refactoring | 9608 | ~2402 | ~60 | ~0 | ~0 | ~582 |
| 33 | `website-concept-designer` | frontend | 9175 | ~2294 | ~122 | ~0 | ~0 | ~1198 |
| 34 | `feature-planner` | workflows | 8678 | ~2238 | ~66 | ~0 | ~200 | ~394 |
| 35 | `demo-deck-designer` | frontend | 8921 | ~2231 | ~101 | ~0 | ~1776 | ~4274 |
| 36 | `ui-ux-architect` | frontend | 8855 | ~2214 | ~103 | ~412 | ~1387 | ~527 |
| 37 | `go-house-libs` | personal | 8785 | ~2197 | ~265 | ~1479 | ~0 | ~592 |
| 38 | `sql-architect` | databases | 8594 | ~2149 | ~90 | ~295 | ~0 | ~1528 |
| 39 | `uru-paper-format` | personal | 8456 | ~2114 | ~141 | ~0 | ~2017 | ~550 |
| 40 | `security-reviewer` | quality | 8273 | ~2069 | ~77 | ~0 | ~539 | ~544 |
| 41 | `narration-writer` | frontend | 8179 | ~2045 | ~134 | ~0 | ~1710 | ~1317 |
| 42 | `performance-reviewer` | quality | 8052 | ~2013 | ~84 | ~0 | ~612 | ~473 |
| 43 | `ddd-architect` | design | 7790 | ~1948 | ~78 | ~0 | ~0 | ~1283 |
| 44 | `nethttp-architect` | frameworks | 7581 | ~1927 | ~121 | ~403 | ~1609 | ~613 |
| 45 | `skill-builder` | meta | 7612 | ~1903 | ~72 | ~0 | ~2337 | ~492 |
| 46 | `fastapi-architect` | frameworks | 7481 | ~1871 | ~120 | ~420 | ~1016 | ~695 |
| 47 | `docker-architect` | infra | 7397 | ~1850 | ~89 | ~333 | ~786 | ~1480 |
| 48 | `design-patterns` | refactoring | 7310 | ~1828 | ~65 | ~0 | ~0 | ~2650 |
| 49 | `protobuf-architect` | encoding | 7271 | ~1818 | ~84 | ~474 | ~914 | ~479 |
| 50 | `improve-codebase-architecture` | refactoring | 7210 | ~1803 | ~60 | ~0 | ~0 | ~4488 |
| 51 | `uru-thesis-defense-designer` | personal | 7107 | ~1777 | ~132 | ~0 | ~1525 | ~500 |
| 52 | `manim-architect` | frameworks | 7103 | ~1776 | ~198 | ~284 | ~1471 | ~2678 |
| 53 | `gin-architect` | frameworks | 6769 | ~1705 | ~109 | ~419 | ~1266 | ~594 |
| 54 | `explainer-designer` | frontend | 6813 | ~1704 | ~187 | ~0 | ~790 | ~1146 |
| 55 | `grafana-architect` | infra | 6739 | ~1685 | ~77 | ~407 | ~1149 | ~622 |
| 56 | `uru-thesis-writer` | personal | 6637 | ~1660 | ~114 | ~0 | ~1697 | ~5465 |
| 57 | `docs-architect` | workflows | 6570 | ~1643 | ~144 | ~0 | ~0 | ~5829 |
| 58 | `adhd-mode` | meta | 5483 | ~1371 | ~75 | ~0 | ~0 | ~343 |
| 59 | `codebase-hygiene-reviewer` | quality | 4745 | ~1187 | ~90 | ~0 | ~1032 | ~428 |
| 60 | `tdd` | workflows | 4169 | ~1073 | ~69 | ~0 | ~0 | ~1781 |
| 61 | `demo-script-designer` | frontend | 3635 | ~939 | ~96 | ~0 | ~0 | ~794 |
| 62 | `mockup-architect` | frontend | 3540 | ~885 | ~136 | ~0 | ~0 | ~474 |
| 63 | `pi-iteration-workflow` | personal | 3511 | ~878 | ~108 | ~0 | ~387 | ~394 |
| 64 | `commit-author` | workflows | 3386 | ~847 | ~81 | ~0 | ~0 | ~409 |
| 65 | `grill-with-docs` | workflows | 3084 | ~771 | ~63 | ~0 | ~0 | ~1965 |
| 66 | `logic-cleaner` | refactoring | 2379 | ~595 | ~93 | ~0 | ~0 | ~405 |
| 67 | `caveman` | meta | 1789 | ~448 | ~71 | ~0 | ~0 | ~272 |

**Totals:** 631721 body bytes · ~158591 body tokens · ~7637 desc tokens · ~191560 side tokens

## Topic files

Side files in a skill directory other than `STACK.md` and `RECIPES.md` — typically topic-specific reference docs referenced by name from `SKILL.md` (e.g. `PAYLOADS.md`, `AUTH_PATTERNS.md`). These contribute to the `~Topic tkns` column above and are loaded only when the parent skill chooses to read them.

- `uru-thesis-reviewer` (~8215 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
  - `TEMPLATES.md` — ~1830 tokens
  - `CHECKS.md` — ~807 tokens
  - `DIGEST.md` — ~628 tokens
- `go-architect` (~7535 tokens total):
  - `LAYOUTS.md` — ~6658 tokens
  - `DIGEST.md` — ~877 tokens
- `docs-architect` (~5829 tokens total):
  - `BACKLOG.md` — ~2339 tokens
  - `MIGRATION.md` — ~1076 tokens
  - `TEMPLATES.md` — ~1057 tokens
  - `FRONTMATTER.md` — ~742 tokens
  - `DIGEST.md` — ~615 tokens
- `uru-thesis-writer` (~5465 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
  - `DIGEST.md` — ~515 tokens
- `rest-api-architect` (~4720 tokens total):
  - `CONCURRENCY.md` — ~1056 tokens
  - `IDEMPOTENCY.md` — ~869 tokens
  - `PAYLOADS.md` — ~787 tokens
  - `AUTH_PATTERNS.md` — ~767 tokens
  - `STATUS_CODES.md` — ~662 tokens
  - `DIGEST.md` — ~579 tokens
- `study-artifact-builder` (~4582 tokens total):
  - `PIPELINES.md` — ~2239 tokens
  - `DATASETS.md` — ~1858 tokens
  - `DIGEST.md` — ~485 tokens
- `improve-codebase-architecture` (~4488 tokens total):
  - `HTML_REPORT.md` — ~1686 tokens
  - `LANGUAGE.md` — ~951 tokens
  - `INTERFACE_DESIGN.md` — ~681 tokens
  - `DEEPENING.md` — ~641 tokens
  - `DIGEST.md` — ~529 tokens
- `demo-deck-designer` (~4274 tokens total):
  - `LAYOUTS.md` — ~3754 tokens
  - `DIGEST.md` — ~520 tokens
- `repo-tooling-architect` (~2679 tokens total):
  - `TASKFILE.md` — ~1949 tokens
  - `DIGEST.md` — ~730 tokens
- `manim-architect` (~2678 tokens total):
  - `TRAPS.md` — ~1078 tokens
  - `FRAMEWORK.md` — ~1061 tokens
  - `DIGEST.md` — ~539 tokens
- `design-patterns` (~2650 tokens total):
  - `PATTERNS.md` — ~2134 tokens
  - `DIGEST.md` — ~516 tokens
- `token-economy` (~2171 tokens total):
  - `GUARDRAILS.md` — ~961 tokens
  - `TOOLS.md` — ~657 tokens
  - `DIGEST.md` — ~553 tokens
- `grill-with-docs` (~1965 tokens total):
  - `CONTEXT_FORMAT.md` — ~854 tokens
  - `ADR_FORMAT.md` — ~703 tokens
  - `DIGEST.md` — ~408 tokens
- `tdd` (~1781 tokens total):
  - `DIGEST.md` — ~421 tokens
  - `TESTS.md` — ~410 tokens
  - `MOCKING.md` — ~370 tokens
  - `DEEP_MODULES.md` — ~320 tokens
  - `INTERFACE_DESIGN.md` — ~163 tokens
  - `REFACTORING.md` — ~97 tokens
- `remotion-architect` (~1581 tokens total):
  - `EXPLAINERS.md` — ~953 tokens
  - `DIGEST.md` — ~628 tokens
- `sql-architect` (~1528 tokens total):
  - `ENGINES.md` — ~563 tokens
  - `POSTGRES.md` — ~501 tokens
  - `DIGEST.md` — ~464 tokens
- `docker-architect` (~1480 tokens total):
  - `CHECKLIST.md` — ~924 tokens
  - `DIGEST.md` — ~556 tokens
- `narration-writer` (~1317 tokens total):
  - `CONVERSATIONAL.md` — ~732 tokens
  - `DIGEST.md` — ~585 tokens
- `ddd-architect` (~1283 tokens total):
  - `PATTERNS.md` — ~751 tokens
  - `DIGEST.md` — ~532 tokens
- `website-concept-designer` (~1198 tokens total):
  - `REFERENCE.md` — ~656 tokens
  - `DIGEST.md` — ~542 tokens
- `explainer-designer` (~1146 tokens total):
  - `LEARNERS.md` — ~688 tokens
  - `DIGEST.md` — ~458 tokens
- `landing-page-designer` (~1134 tokens total):
  - `DIGEST.md` — ~596 tokens
  - `REFERENCE.md` — ~538 tokens
- `demo-script-designer` (~794 tokens total):
  - `REFERENCE.md` — ~417 tokens
  - `DIGEST.md` — ~377 tokens
- `backend-performance-architect` (~751 tokens total):
  - `DIGEST.md` — ~751 tokens
- `wails-architect` (~750 tokens total):
  - `DIGEST.md` — ~750 tokens
- `rsk-guide` (~739 tokens total):
  - `DIGEST.md` — ~496 tokens
  - `REFERENCE.md` — ~243 tokens
- `python-architect` (~733 tokens total):
  - `DIGEST.md` — ~733 tokens
- `mcp-architect` (~710 tokens total):
  - `DIGEST.md` — ~710 tokens
- `go-library-builder` (~705 tokens total):
  - `DIGEST.md` — ~705 tokens
- `fastapi-architect` (~695 tokens total):
  - `DIGEST.md` — ~695 tokens
- `grpc-architect` (~652 tokens total):
  - `DIGEST.md` — ~652 tokens
- `grafana-architect` (~622 tokens total):
  - `DIGEST.md` — ~622 tokens
- `api-contract-reviewer` (~616 tokens total):
  - `DIGEST.md` — ~616 tokens
- `nethttp-architect` (~613 tokens total):
  - `DIGEST.md` — ~613 tokens
- `gin-architect` (~594 tokens total):
  - `DIGEST.md` — ~594 tokens
- `go-house-libs` (~592 tokens total):
  - `DIGEST.md` — ~592 tokens
- `code-design-refactor` (~582 tokens total):
  - `DIGEST.md` — ~448 tokens
  - `REFERENCE.md` — ~134 tokens
- `astro-architect` (~559 tokens total):
  - `DIGEST.md` — ~559 tokens
- `secrets-architect` (~558 tokens total):
  - `DIGEST.md` — ~558 tokens
- `nextjs-architect` (~555 tokens total):
  - `DIGEST.md` — ~555 tokens
- `connectrpc-architect` (~552 tokens total):
  - `DIGEST.md` — ~552 tokens
- `uru-paper-format` (~550 tokens total):
  - `DIGEST.md` — ~550 tokens
- `latex-architect` (~546 tokens total):
  - `DIGEST.md` — ~546 tokens
- `security-reviewer` (~544 tokens total):
  - `DIGEST.md` — ~544 tokens
- `react-architect` (~535 tokens total):
  - `DIGEST.md` — ~535 tokens
- `ci-cd-architect` (~529 tokens total):
  - `DIGEST.md` — ~529 tokens
- `ui-ux-architect` (~527 tokens total):
  - `DIGEST.md` — ~527 tokens
- `reveal-js-architect` (~526 tokens total):
  - `DIGEST.md` — ~526 tokens
- `cli-tool-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `ros2-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `event-driven-architect` (~517 tokens total):
  - `DIGEST.md` — ~517 tokens
- `ml-paper-format` (~511 tokens total):
  - `DIGEST.md` — ~511 tokens
- `uru-thesis-defense-designer` (~500 tokens total):
  - `DIGEST.md` — ~500 tokens
- `observability-architect` (~498 tokens total):
  - `DIGEST.md` — ~498 tokens
- `skill-builder` (~492 tokens total):
  - `DIGEST.md` — ~492 tokens
- `protobuf-architect` (~479 tokens total):
  - `DIGEST.md` — ~479 tokens
- `mockup-architect` (~474 tokens total):
  - `DIGEST.md` — ~474 tokens
- `performance-reviewer` (~473 tokens total):
  - `DIGEST.md` — ~473 tokens
- `hugo-architect` (~460 tokens total):
  - `DIGEST.md` — ~460 tokens
- `codebase-hygiene-reviewer` (~428 tokens total):
  - `DIGEST.md` — ~428 tokens
- `hexagonal-architecture` (~420 tokens total):
  - `DIGEST.md` — ~420 tokens
- `commit-author` (~409 tokens total):
  - `DIGEST.md` — ~409 tokens
- `logic-cleaner` (~405 tokens total):
  - `DIGEST.md` — ~405 tokens
- `pi-iteration-workflow` (~394 tokens total):
  - `DIGEST.md` — ~394 tokens
- `feature-planner` (~394 tokens total):
  - `DIGEST.md` — ~394 tokens
- `adhd-mode` (~343 tokens total):
  - `DIGEST.md` — ~343 tokens
- `caveman` (~272 tokens total):
  - `DIGEST.md` — ~272 tokens

## Skills to consider trimming

Body > 2500 tokens — consider moving examples to `RECIPES.md` or topic files:

- `mcp-architect` (~7795 body tokens)
- `go-architect` (~7244 body tokens)
- `go-library-builder` (~4026 body tokens)
- `backend-performance-architect` (~3413 body tokens)
- `repo-tooling-architect` (~3371 body tokens)
- `uru-thesis-reviewer` (~3299 body tokens)
- `remotion-architect` (~3218 body tokens)
- `python-architect` (~3212 body tokens)
- `ci-cd-architect` (~3085 body tokens)
- `token-economy` (~3063 body tokens)
- `landing-page-designer` (~3057 body tokens)
- `connectrpc-architect` (~3046 body tokens)
- `latex-architect` (~3026 body tokens)
- `nextjs-architect` (~2973 body tokens)
- `react-architect` (~2968 body tokens)
- `rsk-guide` (~2967 body tokens)
- `astro-architect` (~2889 body tokens)
- `cli-tool-architect` (~2784 body tokens)
- `event-driven-architect` (~2770 body tokens)
- `hugo-architect` (~2720 body tokens)
- `ros2-architect` (~2701 body tokens)
- `rest-api-architect` (~2590 body tokens)
- `secrets-architect` (~2555 body tokens)
- `ml-paper-format` (~2551 body tokens)
- `study-artifact-builder` (~2546 body tokens)
- `observability-architect` (~2523 body tokens)
- `api-contract-reviewer` (~2511 body tokens)

Heaviest 10 descriptions — each desc token is paid every turn for any session that installs the skill:

- `mcp-architect` (~334 desc tokens)
- `go-house-libs` (~265 desc tokens)
- `go-architect` (~211 desc tokens)
- `backend-performance-architect` (~204 desc tokens)
- `manim-architect` (~198 desc tokens)
- `explainer-designer` (~187 desc tokens)
- `go-library-builder` (~167 desc tokens)
- `landing-page-designer` (~166 desc tokens)
- `connectrpc-architect` (~153 desc tokens)
- `study-artifact-builder` (~149 desc tokens)

## Notes

- **Body tokens** cost only when the skill is invoked in a session.
- **Description tokens** cost every turn for any session that has the skill installed. The cost depends on which bundles are installed, not on the corpus total.
- Side files (`STACK.md`, `RECIPES.md`, topic files) are never auto-loaded; they cost 0 per turn.
- For exact counts: run each file through [`tiktoken`](https://github.com/openai/tiktoken) with `cl100k_base`.

# SKILL.md Token Estimates

> Auto-generated. **Do not edit by hand.** Run `task tokens` to refresh.
>
> Estimate: ~4 bytes/token for bodies, ~3 bytes/token for descriptions (Claude tokenizer). Actual range ±15%.

_Last updated: 2026-11-08 · 61 skills · 15 bundles_

## Load model

Description tokens hit **every turn** for installed skills. Body tokens are paid only when a skill is actually invoked. Side files are never auto-loaded.

| What | When loaded | Estimated tokens |
|---|---|---:|
| All `SKILL.md` bodies | Only when invoked | ~144660 |
| All side files (`STACK` + `RECIPES` + topic files) | On-demand only | ~170341 |
| All `description:` fields (every skill) | Every turn, if all installed | ~7317 |

## Session profiles

Budgets: **global ≤ 1500** desc tokens (paid by every project), **session ≤ 2200** desc tokens (global + one project bundle). Status: `ok` < 90% · `warn` ≥ 90% · `OVER` ≥ 100%.

**Global baseline** (`rsk install global --global`): 16 skills · ~1740 / 1500 desc tokens every turn (`OVER`) · ~29217 body tokens if all invoked.

Per-project bundle additions on top of the global baseline:

| Bundle | Skills | + Desc tkns / turn | Session total | Budget | Status |
|---|---:|---:|---:|---:|---:|
| `docs` *(+5 external)* | 0 | +0 | ~1740 | 2200 | `ok` |
| `observability` | 2 | +151 | ~1891 | 2200 | `ok` |
| `event-driven` | 2 | +167 | ~1907 | 2200 | `ok` |
| `python-cli` | 3 | +261 | ~2001 | 2200 | `warn` |
| `code-review` | 3 | +269 | ~2009 | 2200 | `warn` |
| `ros2` | 3 | +273 | ~2013 | 2200 | `warn` |
| `design` *(+1 external)* | 3 | +313 | ~2053 | 2200 | `warn` |
| `go-cli` | 3 | +335 | ~2075 | 2200 | `warn` |
| `llm-app` *(+2 missing)* | 4 | +389 | ~2129 | 2200 | `warn` |
| `python-grpc` | 5 | +448 | ~2188 | 2200 | `warn` |
| `fastapi` | 5 | +480 | ~2220 | 2200 | `OVER` |
| `go-grpc` | 5 | +522 | ~2262 | 2200 | `OVER` |
| `gin` | 5 | +543 | ~2283 | 2200 | `OVER` |
| `nethttp` | 5 | +555 | ~2295 | 2200 | `OVER` |

## Personal / unbundled skills

These skills are not part of any bundle. They are installed individually with `rsk install <name> --include personal` and their desc tokens are only paid when explicitly installed.

| Skill | Category | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---|---:|---:|---:|
| `mcp-architect` | protocols | ~7795 | ~334 | ~10985 |
| `go-house-libs` | personal | ~2197 | ~265 | ~2071 |
| `backend-performance-architect` | quality | ~3413 | ~204 | ~3526 |
| `go-library-builder` | personal | ~4008 | ~167 | ~4056 |
| `landing-page-designer` | frontend | ~3057 | ~166 | ~1133 |
| `connectrpc-architect` | protocols | ~3046 | ~153 | ~3750 |
| `study-artifact-builder` | academic | ~2546 | ~149 | ~7180 |
| `ci-cd-architect` | infra | ~3085 | ~145 | ~4673 |
| `docs-architect` | workflows | ~1643 | ~144 | ~5790 |
| `uru-scientific-paper-architect` | personal | ~2118 | ~141 | ~2572 |
| `mockup-architect` | frontend | ~885 | ~136 | ~439 |
| `uru-thesis-defense-architect` | personal | ~1782 | ~132 | ~2026 |
| `secrets-architect` | infra | ~2555 | ~132 | ~558 |
| `astro-architect` | frameworks | ~2889 | ~132 | ~4442 |
| `website-concept-architect` | frontend | ~2299 | ~122 | ~1199 |
| `ml-conference-paper-architect` | academic | ~2568 | ~121 | ~2219 |
| `uru-thesis-reviewer` | personal | ~3299 | ~119 | ~8215 |
| `latex-architect` | languages | ~3053 | ~116 | ~2720 |
| `uru-thesis-architect` | personal | ~1669 | ~114 | ~7164 |
| `pi-iteration-workflow` | personal | ~878 | ~108 | ~781 |
| `hugo-architect` | frameworks | ~2720 | ~107 | ~1817 |
| `demo-presentation-architect` | frontend | ~2233 | ~101 | ~6045 |
| `demo-script-architect` | frontend | ~939 | ~96 | ~795 |
| `reveal-js-architect` | frameworks | ~2476 | ~93 | ~1729 |
| `codebase-hygiene-reviewer` | quality | ~1187 | ~90 | ~1442 |

## By category

| Category | Skills | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---:|---:|---:|---:|
| frameworks | 8 | ~19529 | ~892 | ~20129 |
| personal | 7 | ~15951 | ~1046 | ~26885 |
| protocols | 4 | ~15896 | ~695 | ~21980 |
| languages | 3 | ~12639 | ~348 | ~15102 |
| infra | 5 | ~11698 | ~517 | ~11921 |
| frontend | 6 | ~11627 | ~724 | ~11937 |
| quality | 5 | ~11193 | ~563 | ~7752 |
| meta | 4 | ~8166 | ~572 | ~7287 |
| workflows | 5 | ~6510 | ~512 | ~10500 |
| refactoring | 4 | ~6476 | ~408 | ~8121 |
| tooling | 2 | ~6155 | ~229 | ~7186 |
| academic | 2 | ~5114 | ~270 | ~9399 |
| design | 2 | ~4268 | ~179 | ~2467 |
| messaging | 1 | ~2770 | ~83 | ~3169 |
| robotics | 1 | ~2701 | ~105 | ~2816 |
| databases | 1 | ~2149 | ~90 | ~1823 |
| encoding | 1 | ~1818 | ~84 | ~1867 |

## By bundle

Catalog from `internal/config/catalog.toml`. `External` are `source = "official"` skills (e.g. `docx`, `xlsx`) that live outside this repo by design. `Missing` are `source = "local"` skills the bundle references that aren't in the local `skills/` tree.

| Bundle | Local | External | Missing | ~Body tkns (local) | ~Desc tkns (local) |
|---|---:|---:|---:|---:|---:|
| `global` | 16 | 0 | 0 | ~29217 | ~1740 |
| `docs` | 0 | 5 | 0 | ~0 | ~0 |
| `design` | 3 | 1 | 0 | ~8155 | ~313 |
| `go-grpc` | 5 | 0 | 0 | ~14656 | ~522 |
| `gin` | 5 | 0 | 0 | ~14668 | ~543 |
| `nethttp` | 5 | 0 | 0 | ~14890 | ~555 |
| `go-cli` | 3 | 0 | 0 | ~11008 | ~335 |
| `fastapi` | 5 | 0 | 0 | ~11672 | ~480 |
| `llm-app` | 4 | 0 | 2 | ~9310 | ~389 |
| `ros2` | 3 | 0 | 0 | ~7763 | ~273 |
| `python-grpc` | 5 | 0 | 0 | ~11494 | ~448 |
| `python-cli` | 3 | 0 | 0 | ~7846 | ~261 |
| `event-driven` | 2 | 0 | 0 | ~4588 | ~167 |
| `observability` | 2 | 0 | 0 | ~4208 | ~151 |
| `code-review` | 3 | 0 | 0 | ~6593 | ~269 |

## Per skill

| # | Skill | Category | Body bytes | ~Body tkns | ~Desc tkns | ~Stack tkns | ~Recipes tkns | ~Topic tkns |
|---|---|---|---:|---:|---:|---:|---:|---:|
| 1 | `mcp-architect` | protocols | 31177 | ~7795 | ~334 | ~1521 | ~8783 | ~681 |
| 2 | `go-architect` | languages | 25495 | ~6374 | ~153 | ~473 | ~4425 | ~6105 |
| 3 | `go-library-builder` | personal | 16030 | ~4008 | ~167 | ~628 | ~2725 | ~703 |
| 4 | `backend-performance-architect` | quality | 13650 | ~3413 | ~204 | ~0 | ~2747 | ~779 |
| 5 | `repo-tooling-architect` | tooling | 13230 | ~3371 | ~136 | ~685 | ~1334 | ~2699 |
| 6 | `uru-thesis-reviewer` | personal | 13195 | ~3299 | ~119 | ~0 | ~0 | ~8215 |
| 7 | `python-architect` | languages | 12848 | ~3212 | ~79 | ~443 | ~203 | ~733 |
| 8 | `ci-cd-architect` | infra | 12339 | ~3085 | ~145 | ~624 | ~3520 | ~529 |
| 9 | `token-economy` | meta | 12249 | ~3063 | ~230 | ~0 | ~2033 | ~2150 |
| 10 | `landing-page-designer` | frontend | 12179 | ~3057 | ~166 | ~0 | ~0 | ~1133 |
| 11 | `latex-architect` | languages | 12209 | ~3053 | ~116 | ~339 | ~1841 | ~540 |
| 12 | `connectrpc-architect` | protocols | 12182 | ~3046 | ~153 | ~777 | ~2421 | ~552 |
| 13 | `nextjs-architect` | frameworks | 11890 | ~2973 | ~108 | ~367 | ~1824 | ~545 |
| 14 | `react-architect` | frameworks | 11872 | ~2968 | ~102 | ~415 | ~1426 | ~535 |
| 15 | `rsk-guide` | meta | 11136 | ~2967 | ~199 | ~0 | ~0 | ~713 |
| 16 | `astro-architect` | frameworks | 11496 | ~2889 | ~132 | ~604 | ~3298 | ~540 |
| 17 | `cli-tool-architect` | tooling | 11119 | ~2784 | ~93 | ~515 | ~1436 | ~517 |
| 18 | `event-driven-architect` | messaging | 11059 | ~2770 | ~83 | ~514 | ~2148 | ~507 |
| 19 | `hugo-architect` | frameworks | 10877 | ~2720 | ~107 | ~403 | ~973 | ~441 |
| 20 | `ros2-architect` | robotics | 10802 | ~2701 | ~105 | ~579 | ~1714 | ~523 |
| 21 | `rest-api-architect` | protocols | 10358 | ~2590 | ~102 | ~506 | ~0 | ~4702 |
| 22 | `ml-conference-paper-architect` | academic | 10270 | ~2568 | ~121 | ~395 | ~1305 | ~519 |
| 23 | `secrets-architect` | infra | 9940 | ~2555 | ~132 | ~0 | ~0 | ~558 |
| 24 | `study-artifact-builder` | academic | 10184 | ~2546 | ~149 | ~0 | ~2611 | ~4569 |
| 25 | `observability-architect` | infra | 10025 | ~2523 | ~74 | ~546 | ~874 | ~498 |
| 26 | `api-contract-reviewer` | quality | 9798 | ~2511 | ~108 | ~0 | ~0 | ~616 |
| 27 | `reveal-js-architect` | frameworks | 9901 | ~2476 | ~93 | ~147 | ~1056 | ~526 |
| 28 | `grpc-architect` | protocols | 9860 | ~2465 | ~106 | ~521 | ~864 | ~652 |
| 29 | `hexagonal-arch` | design | 9369 | ~2377 | ~101 | ~0 | ~768 | ~416 |
| 30 | `code-design-refactor` | refactoring | 9400 | ~2350 | ~96 | ~0 | ~0 | ~582 |
| 31 | `website-concept-architect` | frontend | 9196 | ~2299 | ~122 | ~0 | ~0 | ~1199 |
| 32 | `demo-presentation-architect` | frontend | 8932 | ~2233 | ~101 | ~0 | ~1776 | ~4269 |
| 33 | `feature-planner` | workflows | 8630 | ~2226 | ~114 | ~0 | ~200 | ~394 |
| 34 | `ui-ux-architect` | frontend | 8855 | ~2214 | ~103 | ~412 | ~1387 | ~527 |
| 35 | `go-house-libs` | personal | 8785 | ~2197 | ~265 | ~1479 | ~0 | ~592 |
| 36 | `sql-architect` | databases | 8594 | ~2149 | ~90 | ~295 | ~0 | ~1528 |
| 37 | `uru-scientific-paper-architect` | personal | 8472 | ~2118 | ~141 | ~0 | ~2021 | ~551 |
| 38 | `security-reviewer` | quality | 8273 | ~2069 | ~77 | ~0 | ~539 | ~544 |
| 39 | `performance-reviewer` | quality | 8052 | ~2013 | ~84 | ~0 | ~612 | ~473 |
| 40 | `nethttp-architect` | frameworks | 7581 | ~1927 | ~121 | ~403 | ~1609 | ~607 |
| 41 | `ddd-architect` | design | 7561 | ~1891 | ~78 | ~0 | ~0 | ~1283 |
| 42 | `fastapi-architect` | frameworks | 7481 | ~1871 | ~120 | ~420 | ~1016 | ~695 |
| 43 | `docker-architect` | infra | 7397 | ~1850 | ~89 | ~333 | ~786 | ~1480 |
| 44 | `protobuf-architect` | encoding | 7271 | ~1818 | ~84 | ~474 | ~914 | ~479 |
| 45 | `uru-thesis-defense-architect` | personal | 7126 | ~1782 | ~132 | ~0 | ~1525 | ~501 |
| 46 | `improve-codebase-architecture` | refactoring | 7070 | ~1768 | ~97 | ~0 | ~0 | ~4488 |
| 47 | `design-patterns` | refactoring | 7049 | ~1763 | ~122 | ~0 | ~0 | ~2646 |
| 48 | `skill-builder` | meta | 6854 | ~1714 | ~72 | ~0 | ~1728 | ~391 |
| 49 | `gin-architect` | frameworks | 6769 | ~1705 | ~109 | ~419 | ~1266 | ~594 |
| 50 | `grafana-architect` | infra | 6739 | ~1685 | ~77 | ~407 | ~1149 | ~617 |
| 51 | `uru-thesis-architect` | personal | 6676 | ~1669 | ~114 | ~0 | ~1698 | ~5466 |
| 52 | `docs-architect` | workflows | 6570 | ~1643 | ~144 | ~0 | ~0 | ~5790 |
| 53 | `codebase-hygiene-reviewer` | quality | 4745 | ~1187 | ~90 | ~0 | ~1032 | ~410 |
| 54 | `tdd` | workflows | 4080 | ~1051 | ~69 | ~0 | ~0 | ~1771 |
| 55 | `demo-script-architect` | frontend | 3636 | ~939 | ~96 | ~0 | ~0 | ~795 |
| 56 | `mockup-architect` | frontend | 3540 | ~885 | ~136 | ~0 | ~0 | ~439 |
| 57 | `pi-iteration-workflow` | personal | 3511 | ~878 | ~108 | ~0 | ~387 | ~394 |
| 58 | `commit-author` | workflows | 3386 | ~847 | ~81 | ~0 | ~0 | ~380 |
| 59 | `grill-with-docs` | workflows | 2972 | ~743 | ~104 | ~0 | ~0 | ~1965 |
| 60 | `logic-cleaner` | refactoring | 2379 | ~595 | ~93 | ~0 | ~0 | ~405 |
| 61 | `caveman` | meta | 1687 | ~422 | ~71 | ~0 | ~0 | ~272 |

**Totals:** 576008 body bytes · ~144660 body tokens · ~7317 desc tokens · ~170341 side tokens

## Topic files

Side files in a skill directory other than `STACK.md` and `RECIPES.md` — typically topic-specific reference docs referenced by name from `SKILL.md` (e.g. `PAYLOADS.md`, `AUTH_PATTERNS.md`). These contribute to the `~Topic tkns` column above and are loaded only when the parent skill chooses to read them.

- `uru-thesis-reviewer` (~8215 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
  - `TEMPLATES.md` — ~1830 tokens
  - `CHECKS.md` — ~807 tokens
  - `DIGEST.md` — ~628 tokens
- `go-architect` (~6105 tokens total):
  - `LAYOUTS.md` — ~5441 tokens
  - `DIGEST.md` — ~664 tokens
- `docs-architect` (~5790 tokens total):
  - `BACKLOG.md` — ~2339 tokens
  - `MIGRATION.md` — ~1076 tokens
  - `TEMPLATES.md` — ~1057 tokens
  - `FRONTMATTER.md` — ~742 tokens
  - `DIGEST.md` — ~576 tokens
- `uru-thesis-architect` (~5466 tokens total):
  - `NORMAS_URU_2020.md` — ~4950 tokens
  - `DIGEST.md` — ~516 tokens
- `rest-api-architect` (~4702 tokens total):
  - `CONCURRENCY.md` — ~1056 tokens
  - `IDEMPOTENCY.md` — ~869 tokens
  - `PAYLOADS.md` — ~787 tokens
  - `AUTH_PATTERNS.md` — ~767 tokens
  - `STATUS_CODES.md` — ~662 tokens
  - `DIGEST.md` — ~561 tokens
- `study-artifact-builder` (~4569 tokens total):
  - `PIPELINES.md` — ~2239 tokens
  - `DATASETS.md` — ~1858 tokens
  - `DIGEST.md` — ~472 tokens
- `improve-codebase-architecture` (~4488 tokens total):
  - `HTML_REPORT.md` — ~1686 tokens
  - `LANGUAGE.md` — ~951 tokens
  - `INTERFACE_DESIGN.md` — ~681 tokens
  - `DEEPENING.md` — ~641 tokens
  - `DIGEST.md` — ~529 tokens
- `demo-presentation-architect` (~4269 tokens total):
  - `LAYOUTS.md` — ~3754 tokens
  - `DIGEST.md` — ~515 tokens
- `repo-tooling-architect` (~2699 tokens total):
  - `TASKFILE.md` — ~1949 tokens
  - `DIGEST.md` — ~750 tokens
- `design-patterns` (~2646 tokens total):
  - `PATTERNS.md` — ~2130 tokens
  - `DIGEST.md` — ~516 tokens
- `token-economy` (~2150 tokens total):
  - `GUARDRAILS.md` — ~961 tokens
  - `TOOLS.md` — ~657 tokens
  - `DIGEST.md` — ~532 tokens
- `grill-with-docs` (~1965 tokens total):
  - `CONTEXT_FORMAT.md` — ~854 tokens
  - `ADR_FORMAT.md` — ~703 tokens
  - `DIGEST.md` — ~408 tokens
- `tdd` (~1771 tokens total):
  - `DIGEST.md` — ~421 tokens
  - `TESTS.md` — ~410 tokens
  - `MOCKING.md` — ~370 tokens
  - `DEEP_MODULES.md` — ~310 tokens
  - `INTERFACE_DESIGN.md` — ~163 tokens
  - `REFACTORING.md` — ~97 tokens
- `sql-architect` (~1528 tokens total):
  - `ENGINES.md` — ~563 tokens
  - `POSTGRES.md` — ~501 tokens
  - `DIGEST.md` — ~464 tokens
- `docker-architect` (~1480 tokens total):
  - `CHECKLIST.md` — ~924 tokens
  - `DIGEST.md` — ~556 tokens
- `ddd-architect` (~1283 tokens total):
  - `PATTERNS.md` — ~751 tokens
  - `DIGEST.md` — ~532 tokens
- `website-concept-architect` (~1199 tokens total):
  - `REFERENCE.md` — ~656 tokens
  - `DIGEST.md` — ~543 tokens
- `landing-page-designer` (~1133 tokens total):
  - `DIGEST.md` — ~595 tokens
  - `REFERENCE.md` — ~538 tokens
- `demo-script-architect` (~795 tokens total):
  - `REFERENCE.md` — ~417 tokens
  - `DIGEST.md` — ~378 tokens
- `backend-performance-architect` (~779 tokens total):
  - `DIGEST.md` — ~779 tokens
- `python-architect` (~733 tokens total):
  - `DIGEST.md` — ~733 tokens
- `rsk-guide` (~713 tokens total):
  - `DIGEST.md` — ~470 tokens
  - `REFERENCE.md` — ~243 tokens
- `go-library-builder` (~703 tokens total):
  - `DIGEST.md` — ~703 tokens
- `fastapi-architect` (~695 tokens total):
  - `DIGEST.md` — ~695 tokens
- `mcp-architect` (~681 tokens total):
  - `DIGEST.md` — ~681 tokens
- `grpc-architect` (~652 tokens total):
  - `DIGEST.md` — ~652 tokens
- `grafana-architect` (~617 tokens total):
  - `DIGEST.md` — ~617 tokens
- `api-contract-reviewer` (~616 tokens total):
  - `DIGEST.md` — ~616 tokens
- `nethttp-architect` (~607 tokens total):
  - `DIGEST.md` — ~607 tokens
- `gin-architect` (~594 tokens total):
  - `DIGEST.md` — ~594 tokens
- `go-house-libs` (~592 tokens total):
  - `DIGEST.md` — ~592 tokens
- `code-design-refactor` (~582 tokens total):
  - `DIGEST.md` — ~448 tokens
  - `REFERENCE.md` — ~134 tokens
- `secrets-architect` (~558 tokens total):
  - `DIGEST.md` — ~558 tokens
- `connectrpc-architect` (~552 tokens total):
  - `DIGEST.md` — ~552 tokens
- `uru-scientific-paper-architect` (~551 tokens total):
  - `DIGEST.md` — ~551 tokens
- `nextjs-architect` (~545 tokens total):
  - `DIGEST.md` — ~545 tokens
- `security-reviewer` (~544 tokens total):
  - `DIGEST.md` — ~544 tokens
- `latex-architect` (~540 tokens total):
  - `DIGEST.md` — ~540 tokens
- `astro-architect` (~540 tokens total):
  - `DIGEST.md` — ~540 tokens
- `react-architect` (~535 tokens total):
  - `DIGEST.md` — ~535 tokens
- `ci-cd-architect` (~529 tokens total):
  - `DIGEST.md` — ~529 tokens
- `ui-ux-architect` (~527 tokens total):
  - `DIGEST.md` — ~527 tokens
- `reveal-js-architect` (~526 tokens total):
  - `DIGEST.md` — ~526 tokens
- `ros2-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `ml-conference-paper-architect` (~519 tokens total):
  - `DIGEST.md` — ~519 tokens
- `cli-tool-architect` (~517 tokens total):
  - `DIGEST.md` — ~517 tokens
- `event-driven-architect` (~507 tokens total):
  - `DIGEST.md` — ~507 tokens
- `uru-thesis-defense-architect` (~501 tokens total):
  - `DIGEST.md` — ~501 tokens
- `observability-architect` (~498 tokens total):
  - `DIGEST.md` — ~498 tokens
- `protobuf-architect` (~479 tokens total):
  - `DIGEST.md` — ~479 tokens
- `performance-reviewer` (~473 tokens total):
  - `DIGEST.md` — ~473 tokens
- `hugo-architect` (~441 tokens total):
  - `DIGEST.md` — ~441 tokens
- `mockup-architect` (~439 tokens total):
  - `DIGEST.md` — ~439 tokens
- `hexagonal-arch` (~416 tokens total):
  - `DIGEST.md` — ~416 tokens
- `codebase-hygiene-reviewer` (~410 tokens total):
  - `DIGEST.md` — ~410 tokens
- `logic-cleaner` (~405 tokens total):
  - `DIGEST.md` — ~405 tokens
- `pi-iteration-workflow` (~394 tokens total):
  - `DIGEST.md` — ~394 tokens
- `feature-planner` (~394 tokens total):
  - `DIGEST.md` — ~394 tokens
- `skill-builder` (~391 tokens total):
  - `DIGEST.md` — ~391 tokens
- `commit-author` (~380 tokens total):
  - `DIGEST.md` — ~380 tokens
- `caveman` (~272 tokens total):
  - `DIGEST.md` — ~272 tokens

## Skills to consider trimming

Body > 2500 tokens — consider moving examples to `RECIPES.md` or topic files:

- `mcp-architect` (~7795 body tokens)
- `go-architect` (~6374 body tokens)
- `go-library-builder` (~4008 body tokens)
- `backend-performance-architect` (~3413 body tokens)
- `repo-tooling-architect` (~3371 body tokens)
- `uru-thesis-reviewer` (~3299 body tokens)
- `python-architect` (~3212 body tokens)
- `ci-cd-architect` (~3085 body tokens)
- `token-economy` (~3063 body tokens)
- `landing-page-designer` (~3057 body tokens)
- `latex-architect` (~3053 body tokens)
- `connectrpc-architect` (~3046 body tokens)
- `nextjs-architect` (~2973 body tokens)
- `react-architect` (~2968 body tokens)
- `rsk-guide` (~2967 body tokens)
- `astro-architect` (~2889 body tokens)
- `cli-tool-architect` (~2784 body tokens)
- `event-driven-architect` (~2770 body tokens)
- `hugo-architect` (~2720 body tokens)
- `ros2-architect` (~2701 body tokens)
- `rest-api-architect` (~2590 body tokens)
- `ml-conference-paper-architect` (~2568 body tokens)
- `secrets-architect` (~2555 body tokens)
- `study-artifact-builder` (~2546 body tokens)
- `observability-architect` (~2523 body tokens)
- `api-contract-reviewer` (~2511 body tokens)

Heaviest 10 descriptions — each desc token is paid every turn for any session that installs the skill:

- `mcp-architect` (~334 desc tokens)
- `go-house-libs` (~265 desc tokens)
- `token-economy` (~230 desc tokens)
- `backend-performance-architect` (~204 desc tokens)
- `rsk-guide` (~199 desc tokens)
- `go-library-builder` (~167 desc tokens)
- `landing-page-designer` (~166 desc tokens)
- `connectrpc-architect` (~153 desc tokens)
- `go-architect` (~153 desc tokens)
- `study-artifact-builder` (~149 desc tokens)

## Notes

- **Body tokens** cost only when the skill is invoked in a session.
- **Description tokens** cost every turn for any session that has the skill installed. The cost depends on which bundles are installed, not on the corpus total.
- Side files (`STACK.md`, `RECIPES.md`, topic files) are never auto-loaded; they cost 0 per turn.
- For exact counts: run each file through [`tiktoken`](https://github.com/openai/tiktoken) with `cl100k_base`.

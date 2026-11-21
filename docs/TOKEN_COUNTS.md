---
status: active
project: ralvaskills
date: 2026-11-21
related: []
---

# SKILL.md Token Estimates

> Auto-generated. **Do not edit by hand.** Run `task tokens` to refresh.
>
> Estimate: ~4 bytes/token for bodies, ~3 bytes/token for descriptions (Claude tokenizer). Actual range ±15%.

_Last updated: 2026-11-21 · 65 skills · 15 bundles_

## Load model

Description tokens hit **every turn** for installed skills. Body tokens are paid only when a skill is actually invoked. Side files are never auto-loaded.

| What | When loaded | Estimated tokens |
|---|---|---:|
| All `SKILL.md` bodies | Only when invoked | ~149869 |
| All side files (`STACK` + `RECIPES` + topic files) | On-demand only | ~170945 |
| All `description:` fields (every skill) | Every turn, if all installed | ~7083 |

## Session profiles

Budgets: **global ≤ 1500** desc tokens (paid by every project), **session ≤ 2200** desc tokens (global + one project bundle). Status: `ok` < 90% · `warn` ≥ 90% · `OVER` ≥ 100%.

**Global baseline** (`rsk install global --global`): 18 skills · ~1328 / 1500 desc tokens every turn (`ok`) · ~31970 body tokens if all invoked.

Per-project bundle additions on top of the global baseline:

| Bundle | Skills | + Desc tkns / turn | Session total | Budget | Status |
|---|---:|---:|---:|---:|---:|
| `docs` *(+5 external)* | 1 | +94 | ~1422 | 2200 | `ok` |
| `observability` | 2 | +150 | ~1478 | 2200 | `ok` |
| `event-driven` | 2 | +165 | ~1493 | 2200 | `ok` |
| `python-cli` | 3 | +287 | ~1615 | 2200 | `ok` |
| `ros2` | 3 | +300 | ~1628 | 2200 | `ok` |
| `design` *(+1 external)* | 3 | +310 | ~1638 | 2200 | `ok` |
| `code-review` | 4 | +359 | ~1687 | 2200 | `ok` |
| `llm-app` *(+2 missing)* | 4 | +386 | ~1714 | 2200 | `ok` |
| `go-cli` | 3 | +419 | ~1747 | 2200 | `ok` |
| `python-grpc` | 5 | +472 | ~1800 | 2200 | `ok` |
| `fastapi` | 5 | +505 | ~1833 | 2200 | `ok` |
| `go-grpc` | 5 | +604 | ~1932 | 2200 | `ok` |
| `gin` | 5 | +625 | ~1953 | 2200 | `ok` |
| `nethttp` | 5 | +637 | ~1965 | 2200 | `ok` |

## Personal / unbundled skills

These skills are not part of any bundle. They are installed individually with `rsk install <name> --include personal` and their desc tokens are only paid when explicitly installed.

| Skill | Category | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---|---:|---:|---:|
| `mcp-architect` | protocols | ~8509 | ~300 | ~12623 |
| `backend-performance-architect` | quality | ~3393 | ~202 | ~3492 |
| `manim-architect` | frameworks | ~1776 | ~197 | ~4433 |
| `explainer-designer` | frontend | ~1737 | ~186 | ~2401 |
| `landing-page-designer` | frontend | ~3057 | ~165 | ~1134 |
| `connectrpc-architect` | protocols | ~3028 | ~153 | ~3745 |
| `study-artifact-builder` | academic | ~2546 | ~149 | ~7193 |
| `ci-cd-architect` | infra | ~3061 | ~145 | ~4661 |
| `docs-architect` | workflows | ~1657 | ~143 | ~5815 |
| `remotion-architect` | frameworks | ~3207 | ~138 | ~4485 |
| `mockup-architect` | frontend | ~885 | ~136 | ~474 |
| `notify-on-finish` | tooling | ~1443 | ~135 | ~488 |
| `narration-writer` | frontend | ~2095 | ~133 | ~3027 |
| `astro-architect` | frameworks | ~2872 | ~131 | ~4450 |
| `secrets-architect` | infra | ~2550 | ~131 | ~558 |
| `ml-paper-format` | academic | ~2537 | ~121 | ~2203 |
| `website-concept-designer` | frontend | ~2278 | ~121 | ~1197 |
| `wails-architect` | frameworks | ~2481 | ~117 | ~4592 |
| `latex-architect` | languages | ~3009 | ~115 | ~2714 |
| `test-designer` | quality | ~2240 | ~107 | ~1132 |
| `hugo-architect` | frameworks | ~2695 | ~106 | ~1830 |
| `demo-deck-designer` | frontend | ~2221 | ~101 | ~6015 |
| `demo-script-designer` | frontend | ~937 | ~95 | ~793 |
| `reveal-js-architect` | frameworks | ~2443 | ~93 | ~1726 |
| `codebase-hygiene-reviewer` | quality | ~1184 | ~89 | ~1459 |

## By category

| Category | Skills | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---:|---:|---:|---:|
| frameworks | 11 | ~26841 | ~1338 | ~33648 |
| protocols | 4 | ~16552 | ~659 | ~23611 |
| frontend | 8 | ~15407 | ~1039 | ~17357 |
| quality | 7 | ~14611 | ~757 | ~10528 |
| languages | 3 | ~13401 | ~403 | ~16497 |
| infra | 5 | ~12026 | ~543 | ~12414 |
| meta | 6 | ~10665 | ~453 | ~8803 |
| tooling | 3 | ~7569 | ~308 | ~7631 |
| workflows | 6 | ~7392 | ~530 | ~10855 |
| refactoring | 4 | ~6600 | ~276 | ~8092 |
| academic | 2 | ~5083 | ~270 | ~9396 |
| design | 2 | ~4340 | ~148 | ~2461 |
| messaging | 1 | ~2749 | ~82 | ~3166 |
| robotics | 1 | ~2685 | ~105 | ~2806 |
| databases | 1 | ~2138 | ~89 | ~1815 |
| encoding | 1 | ~1810 | ~83 | ~1865 |

## By bundle

Catalog from `internal/config/catalog.toml`. `External` are `source = "official"` skills (e.g. `docx`, `xlsx`) that live outside this repo by design. `Missing` are `source = "local"` skills the bundle references that aren't in the local `skills/` tree.

| Bundle | Local | External | Missing | ~Body tkns (local) | ~Desc tkns (local) |
|---|---:|---:|---:|---:|---:|
| `global` | 18 | 0 | 0 | ~31970 | ~1328 |
| `docs` | 1 | 5 | 0 | ~780 | ~94 |
| `design` | 3 | 1 | 0 | ~8099 | ~310 |
| `go-grpc` | 5 | 0 | 0 | ~15826 | ~604 |
| `gin` | 5 | 0 | 0 | ~15834 | ~625 |
| `nethttp` | 5 | 0 | 0 | ~16053 | ~637 |
| `go-cli` | 3 | 0 | 0 | ~12206 | ~419 |
| `fastapi` | 5 | 0 | 0 | ~11985 | ~505 |
| `llm-app` | 4 | 0 | 2 | ~9682 | ~386 |
| `ros2` | 3 | 0 | 0 | ~8102 | ~300 |
| `python-grpc` | 5 | 0 | 0 | ~11810 | ~472 |
| `python-cli` | 3 | 0 | 0 | ~8190 | ~287 |
| `event-driven` | 2 | 0 | 0 | ~4559 | ~165 |
| `observability` | 2 | 0 | 0 | ~4186 | ~150 |
| `code-review` | 4 | 0 | 0 | ~7794 | ~359 |

## Per skill

| # | Skill | Category | Body bytes | ~Body tkns | ~Desc tkns | ~Stack tkns | ~Recipes tkns | ~Topic tkns |
|---|---|---|---:|---:|---:|---:|---:|---:|
| 1 | `mcp-architect` | protocols | 34034 | ~8509 | ~300 | ~1509 | ~10199 | ~915 |
| 2 | `go-architect` | languages | 28816 | ~7204 | ~210 | ~471 | ~4416 | ~7522 |
| 3 | `backend-performance-architect` | quality | 13570 | ~3393 | ~202 | ~0 | ~2742 | ~750 |
| 4 | `repo-tooling-architect` | tooling | 13164 | ~3353 | ~81 | ~682 | ~1328 | ~2672 |
| 5 | `remotion-architect` | frameworks | 12828 | ~3207 | ~138 | ~316 | ~2588 | ~1581 |
| 6 | `python-architect` | languages | 12752 | ~3188 | ~78 | ~438 | ~203 | ~733 |
| 7 | `token-economy` | meta | 12693 | ~3174 | ~72 | ~0 | ~2032 | ~2192 |
| 8 | `ci-cd-architect` | infra | 12241 | ~3061 | ~145 | ~621 | ~3511 | ~529 |
| 9 | `landing-page-designer` | frontend | 12177 | ~3057 | ~165 | ~0 | ~0 | ~1134 |
| 10 | `connectrpc-architect` | protocols | 12110 | ~3028 | ~153 | ~775 | ~2418 | ~552 |
| 11 | `latex-architect` | languages | 12035 | ~3009 | ~115 | ~331 | ~1837 | ~546 |
| 12 | `rsk-guide` | meta | 11100 | ~2957 | ~70 | ~0 | ~0 | ~738 |
| 13 | `nextjs-architect` | frameworks | 11808 | ~2952 | ~107 | ~365 | ~1820 | ~555 |
| 14 | `react-architect` | frameworks | 11798 | ~2950 | ~101 | ~412 | ~1421 | ~535 |
| 15 | `astro-architect` | frameworks | 11430 | ~2872 | ~131 | ~600 | ~3291 | ~559 |
| 16 | `cli-tool-architect` | tooling | 11073 | ~2773 | ~92 | ~509 | ~1429 | ~523 |
| 17 | `event-driven-architect` | messaging | 10975 | ~2749 | ~82 | ~513 | ~2136 | ~517 |
| 18 | `hugo-architect` | frameworks | 10777 | ~2695 | ~106 | ~399 | ~971 | ~460 |
| 19 | `ros2-architect` | robotics | 10738 | ~2685 | ~105 | ~576 | ~1707 | ~523 |
| 20 | `rest-api-architect` | protocols | 10278 | ~2570 | ~101 | ~505 | ~0 | ~4703 |
| 21 | `secrets-architect` | infra | 9916 | ~2550 | ~131 | ~0 | ~0 | ~558 |
| 22 | `study-artifact-builder` | academic | 10184 | ~2546 | ~149 | ~0 | ~2611 | ~4582 |
| 23 | `ml-paper-format` | academic | 10148 | ~2537 | ~121 | ~394 | ~1298 | ~511 |
| 24 | `observability-architect` | infra | 9975 | ~2511 | ~74 | ~544 | ~869 | ~498 |
| 25 | `api-contract-reviewer` | quality | 9734 | ~2494 | ~108 | ~0 | ~0 | ~616 |
| 26 | `wails-architect` | frameworks | 9923 | ~2481 | ~117 | ~335 | ~3508 | ~749 |
| 27 | `grpc-architect` | protocols | 9780 | ~2445 | ~105 | ~520 | ~863 | ~652 |
| 28 | `reveal-js-architect` | frameworks | 9771 | ~2443 | ~93 | ~146 | ~1054 | ~526 |
| 29 | `hexagonal-architecture` | design | 9484 | ~2405 | ~71 | ~0 | ~764 | ~420 |
| 30 | `code-design-refactor` | refactoring | 9574 | ~2394 | ~59 | ~0 | ~0 | ~582 |
| 31 | `website-concept-designer` | frontend | 9111 | ~2278 | ~121 | ~0 | ~0 | ~1197 |
| 32 | `test-designer` | quality | 8960 | ~2240 | ~107 | ~0 | ~692 | ~440 |
| 33 | `docker-architect` | infra | 8914 | ~2229 | ~117 | ~356 | ~1041 | ~1711 |
| 34 | `demo-deck-designer` | frontend | 8881 | ~2221 | ~101 | ~0 | ~1767 | ~4248 |
| 35 | `feature-planner` | workflows | 8612 | ~2221 | ~66 | ~0 | ~198 | ~394 |
| 36 | `ui-ux-architect` | frontend | 8787 | ~2197 | ~102 | ~411 | ~1378 | ~527 |
| 37 | `performance-reviewer` | quality | 8632 | ~2158 | ~83 | ~0 | ~608 | ~1815 |
| 38 | `sql-architect` | databases | 8552 | ~2138 | ~89 | ~293 | ~0 | ~1522 |
| 39 | `narration-writer` | frontend | 8378 | ~2095 | ~133 | ~0 | ~1710 | ~1317 |
| 40 | `security-reviewer` | quality | 8225 | ~2057 | ~77 | ~0 | ~535 | ~544 |
| 41 | `ddd-architect` | design | 7740 | ~1935 | ~77 | ~0 | ~0 | ~1277 |
| 42 | `nethttp-architect` | frameworks | 7521 | ~1912 | ~120 | ~401 | ~1606 | ~613 |
| 43 | `skill-builder` | meta | 7570 | ~1893 | ~71 | ~0 | ~2450 | ~492 |
| 44 | `fastapi-architect` | frameworks | 7437 | ~1860 | ~120 | ~419 | ~1014 | ~695 |
| 45 | `design-patterns` | refactoring | 7278 | ~1820 | ~64 | ~0 | ~0 | ~2638 |
| 46 | `protobuf-architect` | encoding | 7237 | ~1810 | ~83 | ~474 | ~912 | ~479 |
| 47 | `improve-codebase-architecture` | refactoring | 7204 | ~1801 | ~60 | ~0 | ~0 | ~4467 |
| 48 | `manim-architect` | frameworks | 7103 | ~1776 | ~197 | ~284 | ~1471 | ~2678 |
| 49 | `explainer-designer` | frontend | 6946 | ~1737 | ~186 | ~0 | ~1255 | ~1146 |
| 50 | `gin-architect` | frameworks | 6721 | ~1693 | ~108 | ~418 | ~1264 | ~594 |
| 51 | `grafana-architect` | infra | 6699 | ~1675 | ~76 | ~407 | ~1147 | ~622 |
| 52 | `docs-architect` | workflows | 6625 | ~1657 | ~143 | ~0 | ~0 | ~5815 |
| 53 | `notify-on-finish` | tooling | 5412 | ~1443 | ~135 | ~0 | ~0 | ~488 |
| 54 | `adhd-mode` | meta | 5613 | ~1404 | ~75 | ~0 | ~0 | ~343 |
| 55 | `codebase-hygiene-reviewer` | quality | 4733 | ~1184 | ~89 | ~0 | ~1031 | ~428 |
| 56 | `blast-radius-reviewer` | quality | 4337 | ~1085 | ~91 | ~0 | ~0 | ~327 |
| 57 | `tdd` | workflows | 4205 | ~1082 | ~69 | ~0 | ~0 | ~1781 |
| 58 | `demo-script-designer` | frontend | 3627 | ~937 | ~95 | ~0 | ~0 | ~793 |
| 59 | `commit-author` | workflows | 3566 | ~892 | ~81 | ~0 | ~0 | ~409 |
| 60 | `mockup-architect` | frontend | 3540 | ~885 | ~136 | ~0 | ~0 | ~474 |
| 61 | `unslop` | meta | 3119 | ~780 | ~94 | ~0 | ~0 | ~284 |
| 62 | `grill-with-docs` | workflows | 3104 | ~776 | ~63 | ~0 | ~0 | ~1947 |
| 63 | `guardrail-builder` | workflows | 3056 | ~764 | ~108 | ~0 | ~0 | ~311 |
| 64 | `logic-cleaner` | refactoring | 2340 | ~585 | ~93 | ~0 | ~0 | ~405 |
| 65 | `caveman` | meta | 1825 | ~457 | ~71 | ~0 | ~0 | ~272 |

**Totals:** 596496 body bytes · ~149869 body tokens · ~7083 desc tokens · ~170945 side tokens

## Topic files

Side files in a skill directory other than `STACK.md` and `RECIPES.md` — typically topic-specific reference docs referenced by name from `SKILL.md` (e.g. `PAYLOADS.md`, `AUTH_PATTERNS.md`). These contribute to the `~Topic tkns` column above and are loaded only when the parent skill chooses to read them.

- `go-architect` (~7522 tokens total):
  - `LAYOUTS.md` — ~6645 tokens
  - `DIGEST.md` — ~877 tokens
- `docs-architect` (~5815 tokens total):
  - `BACKLOG.md` — ~2334 tokens
  - `MIGRATION.md` — ~1071 tokens
  - `TEMPLATES.md` — ~1054 tokens
  - `FRONTMATTER.md` — ~741 tokens
  - `DIGEST.md` — ~615 tokens
- `rest-api-architect` (~4703 tokens total):
  - `CONCURRENCY.md` — ~1052 tokens
  - `IDEMPOTENCY.md` — ~864 tokens
  - `PAYLOADS.md` — ~785 tokens
  - `AUTH_PATTERNS.md` — ~764 tokens
  - `STATUS_CODES.md` — ~659 tokens
  - `DIGEST.md` — ~579 tokens
- `study-artifact-builder` (~4582 tokens total):
  - `PIPELINES.md` — ~2239 tokens
  - `DATASETS.md` — ~1858 tokens
  - `DIGEST.md` — ~485 tokens
- `improve-codebase-architecture` (~4467 tokens total):
  - `HTML_REPORT.md` — ~1676 tokens
  - `LANGUAGE.md` — ~947 tokens
  - `INTERFACE_DESIGN.md` — ~676 tokens
  - `DEEPENING.md` — ~639 tokens
  - `DIGEST.md` — ~529 tokens
- `demo-deck-designer` (~4248 tokens total):
  - `LAYOUTS.md` — ~3728 tokens
  - `DIGEST.md` — ~520 tokens
- `manim-architect` (~2678 tokens total):
  - `TRAPS.md` — ~1078 tokens
  - `FRAMEWORK.md` — ~1061 tokens
  - `DIGEST.md` — ~539 tokens
- `repo-tooling-architect` (~2672 tokens total):
  - `TASKFILE.md` — ~1943 tokens
  - `DIGEST.md` — ~729 tokens
- `design-patterns` (~2638 tokens total):
  - `PATTERNS.md` — ~2122 tokens
  - `DIGEST.md` — ~516 tokens
- `token-economy` (~2192 tokens total):
  - `GUARDRAILS.md` — ~961 tokens
  - `TOOLS.md` — ~656 tokens
  - `DIGEST.md` — ~575 tokens
- `grill-with-docs` (~1947 tokens total):
  - `CONTEXT_FORMAT.md` — ~853 tokens
  - `ADR_FORMAT.md` — ~686 tokens
  - `DIGEST.md` — ~408 tokens
- `performance-reviewer` (~1815 tokens total):
  - `BENCHMARKS.md` — ~1342 tokens
  - `DIGEST.md` — ~473 tokens
- `tdd` (~1781 tokens total):
  - `DIGEST.md` — ~421 tokens
  - `TESTS.md` — ~410 tokens
  - `MOCKING.md` — ~370 tokens
  - `DEEP_MODULES.md` — ~320 tokens
  - `INTERFACE_DESIGN.md` — ~163 tokens
  - `REFACTORING.md` — ~97 tokens
- `docker-architect` (~1711 tokens total):
  - `CHECKLIST.md` — ~1038 tokens
  - `DIGEST.md` — ~673 tokens
- `remotion-architect` (~1581 tokens total):
  - `EXPLAINERS.md` — ~953 tokens
  - `DIGEST.md` — ~628 tokens
- `sql-architect` (~1522 tokens total):
  - `ENGINES.md` — ~560 tokens
  - `POSTGRES.md` — ~498 tokens
  - `DIGEST.md` — ~464 tokens
- `narration-writer` (~1317 tokens total):
  - `CONVERSATIONAL.md` — ~732 tokens
  - `DIGEST.md` — ~585 tokens
- `ddd-architect` (~1277 tokens total):
  - `PATTERNS.md` — ~745 tokens
  - `DIGEST.md` — ~532 tokens
- `website-concept-designer` (~1197 tokens total):
  - `REFERENCE.md` — ~655 tokens
  - `DIGEST.md` — ~542 tokens
- `explainer-designer` (~1146 tokens total):
  - `LEARNERS.md` — ~688 tokens
  - `DIGEST.md` — ~458 tokens
- `landing-page-designer` (~1134 tokens total):
  - `DIGEST.md` — ~596 tokens
  - `REFERENCE.md` — ~538 tokens
- `mcp-architect` (~915 tokens total):
  - `DIGEST.md` — ~915 tokens
- `demo-script-designer` (~793 tokens total):
  - `REFERENCE.md` — ~417 tokens
  - `DIGEST.md` — ~376 tokens
- `backend-performance-architect` (~750 tokens total):
  - `DIGEST.md` — ~750 tokens
- `wails-architect` (~749 tokens total):
  - `DIGEST.md` — ~749 tokens
- `rsk-guide` (~738 tokens total):
  - `DIGEST.md` — ~496 tokens
  - `REFERENCE.md` — ~242 tokens
- `python-architect` (~733 tokens total):
  - `DIGEST.md` — ~733 tokens
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
- `ros2-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `cli-tool-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `event-driven-architect` (~517 tokens total):
  - `DIGEST.md` — ~517 tokens
- `ml-paper-format` (~511 tokens total):
  - `DIGEST.md` — ~511 tokens
- `observability-architect` (~498 tokens total):
  - `DIGEST.md` — ~498 tokens
- `skill-builder` (~492 tokens total):
  - `DIGEST.md` — ~492 tokens
- `notify-on-finish` (~488 tokens total):
  - `DIGEST.md` — ~488 tokens
- `protobuf-architect` (~479 tokens total):
  - `DIGEST.md` — ~479 tokens
- `mockup-architect` (~474 tokens total):
  - `DIGEST.md` — ~474 tokens
- `hugo-architect` (~460 tokens total):
  - `DIGEST.md` — ~460 tokens
- `test-designer` (~440 tokens total):
  - `DIGEST.md` — ~440 tokens
- `codebase-hygiene-reviewer` (~428 tokens total):
  - `DIGEST.md` — ~428 tokens
- `hexagonal-architecture` (~420 tokens total):
  - `DIGEST.md` — ~420 tokens
- `commit-author` (~409 tokens total):
  - `DIGEST.md` — ~409 tokens
- `logic-cleaner` (~405 tokens total):
  - `DIGEST.md` — ~405 tokens
- `feature-planner` (~394 tokens total):
  - `DIGEST.md` — ~394 tokens
- `adhd-mode` (~343 tokens total):
  - `DIGEST.md` — ~343 tokens
- `blast-radius-reviewer` (~327 tokens total):
  - `DIGEST.md` — ~327 tokens
- `guardrail-builder` (~311 tokens total):
  - `DIGEST.md` — ~311 tokens
- `unslop` (~284 tokens total):
  - `DIGEST.md` — ~284 tokens
- `caveman` (~272 tokens total):
  - `DIGEST.md` — ~272 tokens

## Skills to consider trimming

Body > 2500 tokens — consider moving examples to `RECIPES.md` or topic files:

- `mcp-architect` (~8509 body tokens)
- `go-architect` (~7204 body tokens)
- `backend-performance-architect` (~3393 body tokens)
- `repo-tooling-architect` (~3353 body tokens)
- `remotion-architect` (~3207 body tokens)
- `python-architect` (~3188 body tokens)
- `token-economy` (~3174 body tokens)
- `ci-cd-architect` (~3061 body tokens)
- `landing-page-designer` (~3057 body tokens)
- `connectrpc-architect` (~3028 body tokens)
- `latex-architect` (~3009 body tokens)
- `rsk-guide` (~2957 body tokens)
- `nextjs-architect` (~2952 body tokens)
- `react-architect` (~2950 body tokens)
- `astro-architect` (~2872 body tokens)
- `cli-tool-architect` (~2773 body tokens)
- `event-driven-architect` (~2749 body tokens)
- `hugo-architect` (~2695 body tokens)
- `ros2-architect` (~2685 body tokens)
- `rest-api-architect` (~2570 body tokens)
- `secrets-architect` (~2550 body tokens)
- `study-artifact-builder` (~2546 body tokens)
- `ml-paper-format` (~2537 body tokens)
- `observability-architect` (~2511 body tokens)

Heaviest 10 descriptions — each desc token is paid every turn for any session that installs the skill:

- `mcp-architect` (~300 desc tokens)
- `go-architect` (~210 desc tokens)
- `backend-performance-architect` (~202 desc tokens)
- `manim-architect` (~197 desc tokens)
- `explainer-designer` (~186 desc tokens)
- `landing-page-designer` (~165 desc tokens)
- `connectrpc-architect` (~153 desc tokens)
- `study-artifact-builder` (~149 desc tokens)
- `ci-cd-architect` (~145 desc tokens)
- `docs-architect` (~143 desc tokens)

## Notes

- **Body tokens** cost only when the skill is invoked in a session.
- **Description tokens** cost every turn for any session that has the skill installed. The cost depends on which bundles are installed, not on the corpus total.
- Side files (`STACK.md`, `RECIPES.md`, topic files) are never auto-loaded; they cost 0 per turn.
- For exact counts: run each file through [`tiktoken`](https://github.com/openai/tiktoken) with `cl100k_base`.

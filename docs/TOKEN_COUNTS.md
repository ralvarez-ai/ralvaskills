---
status: active
project: ralvaskills
date: 2026-11-26
related: []
---

# SKILL.md Token Estimates

> Auto-generated. **Do not edit by hand.** Run `task tokens` to refresh.
>
> Estimate: ~4 bytes/token for bodies, ~3 bytes/token for descriptions (Claude tokenizer). Actual range ±15%.

_Last updated: 2026-11-26 · 66 skills · 15 bundles_

## Load model

Description tokens hit **every turn** for installed skills. Body tokens are paid only when a skill is actually invoked. Side files are never auto-loaded.

| What | When loaded | Estimated tokens |
|---|---|---:|
| All `SKILL.md` bodies | Only when invoked | ~155068 |
| All side files (`STACK` + `RECIPES` + topic files) | On-demand only | ~178824 |
| All `description:` fields (every skill) | Every turn, if all installed | ~7409 |

## Session profiles

Budgets: **global ≤ 1500** desc tokens (paid by every project), **session ≤ 2200** desc tokens (global + one project bundle). Status: `ok` < 90% · `warn` ≥ 90% · `OVER` ≥ 100%.

**Global baseline** (`rsk install global --global`): 18 skills · ~1328 / 1500 desc tokens every turn (`ok`) · ~31969 body tokens if all invoked.

Per-project bundle additions on top of the global baseline:

| Bundle | Skills | + Desc tkns / turn | Session total | Budget | Status |
|---|---:|---:|---:|---:|---:|
| `docs` *(+5 external)* | 1 | +94 | ~1422 | 2200 | `ok` |
| `observability` | 2 | +150 | ~1478 | 2200 | `ok` |
| `event-driven` | 2 | +165 | ~1493 | 2200 | `ok` |
| `python-cli` | 3 | +305 | ~1633 | 2200 | `ok` |
| `design` *(+1 external)* | 3 | +310 | ~1638 | 2200 | `ok` |
| `ros2` | 3 | +318 | ~1646 | 2200 | `ok` |
| `code-review` | 4 | +360 | ~1688 | 2200 | `ok` |
| `llm-app` *(+2 missing)* | 4 | +404 | ~1732 | 2200 | `ok` |
| `go-cli` | 3 | +461 | ~1789 | 2200 | `ok` |
| `python-grpc` | 5 | +490 | ~1818 | 2200 | `ok` |
| `fastapi` | 5 | +523 | ~1851 | 2200 | `ok` |
| `go-grpc` | 5 | +646 | ~1974 | 2200 | `ok` |
| `gin` | 5 | +667 | ~1995 | 2200 | `warn` |
| `nethttp` | 5 | +679 | ~2007 | 2200 | `warn` |

## Personal / unbundled skills

These skills are not part of any bundle. They are installed individually with `rsk install <name> --include personal` and their desc tokens are only paid when explicitly installed.

| Skill | Category | ~Body tkns | ~Desc tkns | ~Side tkns |
|---|---|---:|---:|---:|
| `mcp-architect` | protocols | ~8753 | ~300 | ~13594 |
| `manim-architect` | frameworks | ~2251 | ~260 | ~5981 |
| `explainer-designer` | frontend | ~2638 | ~244 | ~4291 |
| `backend-performance-architect` | quality | ~3393 | ~202 | ~3492 |
| `landing-page-designer` | frontend | ~3057 | ~165 | ~1134 |
| `remotion-architect` | frameworks | ~3310 | ~154 | ~5358 |
| `connectrpc-architect` | protocols | ~3028 | ~153 | ~3745 |
| `study-artifact-builder` | academic | ~2546 | ~149 | ~7193 |
| `traefik-architect` | infra | ~2166 | ~146 | ~1702 |
| `ci-cd-architect` | infra | ~3061 | ~145 | ~4661 |
| `docs-architect` | workflows | ~1657 | ~143 | ~5815 |
| `mockup-architect` | frontend | ~885 | ~136 | ~474 |
| `notify-on-finish` | tooling | ~1443 | ~135 | ~488 |
| `narration-writer` | frontend | ~2388 | ~133 | ~3267 |
| `secrets-architect` | infra | ~2550 | ~131 | ~558 |
| `astro-architect` | frameworks | ~2872 | ~131 | ~4449 |
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
| frameworks | 11 | ~27419 | ~1417 | ~36068 |
| protocols | 4 | ~16872 | ~659 | ~24660 |
| frontend | 8 | ~16601 | ~1097 | ~19487 |
| infra | 6 | ~14709 | ~707 | ~14236 |
| quality | 7 | ~14692 | ~758 | ~10529 |
| languages | 3 | ~13745 | ~427 | ~16954 |
| meta | 6 | ~10664 | ~453 | ~8803 |
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
| `global` | 18 | 0 | 0 | ~31969 | ~1328 |
| `docs` | 1 | 5 | 0 | ~780 | ~94 |
| `design` | 3 | 1 | 0 | ~8099 | ~310 |
| `go-grpc` | 5 | 0 | 0 | ~16687 | ~646 |
| `gin` | 5 | 0 | 0 | ~16771 | ~667 |
| `nethttp` | 5 | 0 | 0 | ~16990 | ~679 |
| `go-cli` | 3 | 0 | 0 | ~13067 | ~461 |
| `fastapi` | 5 | 0 | 0 | ~12578 | ~523 |
| `llm-app` | 4 | 0 | 2 | ~10199 | ~404 |
| `ros2` | 3 | 0 | 0 | ~8619 | ~318 |
| `python-grpc` | 5 | 0 | 0 | ~12327 | ~490 |
| `python-cli` | 3 | 0 | 0 | ~8707 | ~305 |
| `event-driven` | 2 | 0 | 0 | ~4559 | ~165 |
| `observability` | 2 | 0 | 0 | ~4186 | ~150 |
| `code-review` | 4 | 0 | 0 | ~7875 | ~360 |

## Per skill

| # | Skill | Category | Body bytes | ~Body tkns | ~Desc tkns | ~Stack tkns | ~Recipes tkns | ~Topic tkns |
|---|---|---|---:|---:|---:|---:|---:|---:|
| 1 | `mcp-architect` | protocols | 35009 | ~8753 | ~300 | ~1509 | ~11108 | ~977 |
| 2 | `go-architect` | languages | 30191 | ~7548 | ~234 | ~471 | ~4455 | ~7940 |
| 3 | `backend-performance-architect` | quality | 13570 | ~3393 | ~202 | ~0 | ~2742 | ~750 |
| 4 | `repo-tooling-architect` | tooling | 13164 | ~3353 | ~81 | ~682 | ~1328 | ~2672 |
| 5 | `remotion-architect` | frameworks | 13239 | ~3310 | ~154 | ~316 | ~2588 | ~2454 |
| 6 | `python-architect` | languages | 12752 | ~3188 | ~78 | ~438 | ~203 | ~733 |
| 7 | `token-economy` | meta | 12693 | ~3174 | ~72 | ~0 | ~2032 | ~2192 |
| 8 | `ci-cd-architect` | infra | 12241 | ~3061 | ~145 | ~621 | ~3511 | ~529 |
| 9 | `landing-page-designer` | frontend | 12177 | ~3057 | ~165 | ~0 | ~0 | ~1134 |
| 10 | `connectrpc-architect` | protocols | 12110 | ~3028 | ~153 | ~775 | ~2418 | ~552 |
| 11 | `latex-architect` | languages | 12035 | ~3009 | ~115 | ~331 | ~1837 | ~546 |
| 12 | `rsk-guide` | meta | 11097 | ~2956 | ~70 | ~0 | ~0 | ~738 |
| 13 | `nextjs-architect` | frameworks | 11808 | ~2952 | ~107 | ~365 | ~1820 | ~555 |
| 14 | `react-architect` | frameworks | 11798 | ~2950 | ~101 | ~412 | ~1421 | ~535 |
| 15 | `astro-architect` | frameworks | 11430 | ~2872 | ~131 | ~600 | ~3290 | ~559 |
| 16 | `cli-tool-architect` | tooling | 11073 | ~2773 | ~92 | ~509 | ~1429 | ~523 |
| 17 | `event-driven-architect` | messaging | 10975 | ~2749 | ~82 | ~513 | ~2136 | ~517 |
| 18 | `docker-architect` | infra | 10984 | ~2746 | ~135 | ~356 | ~1041 | ~1831 |
| 19 | `hugo-architect` | frameworks | 10777 | ~2695 | ~106 | ~399 | ~971 | ~460 |
| 20 | `ros2-architect` | robotics | 10738 | ~2685 | ~105 | ~576 | ~1707 | ~523 |
| 21 | `rest-api-architect` | protocols | 10581 | ~2646 | ~101 | ~553 | ~0 | ~4733 |
| 22 | `explainer-designer` | frontend | 10552 | ~2638 | ~244 | ~0 | ~1767 | ~2524 |
| 23 | `api-contract-reviewer` | quality | 10059 | ~2575 | ~109 | ~0 | ~0 | ~617 |
| 24 | `secrets-architect` | infra | 9916 | ~2550 | ~131 | ~0 | ~0 | ~558 |
| 25 | `study-artifact-builder` | academic | 10184 | ~2546 | ~149 | ~0 | ~2611 | ~4582 |
| 26 | `ml-paper-format` | academic | 10148 | ~2537 | ~121 | ~394 | ~1298 | ~511 |
| 27 | `observability-architect` | infra | 9975 | ~2511 | ~74 | ~544 | ~869 | ~498 |
| 28 | `wails-architect` | frameworks | 9923 | ~2481 | ~117 | ~335 | ~3508 | ~749 |
| 29 | `grpc-architect` | protocols | 9780 | ~2445 | ~105 | ~520 | ~863 | ~652 |
| 30 | `reveal-js-architect` | frameworks | 9771 | ~2443 | ~93 | ~146 | ~1054 | ~526 |
| 31 | `hexagonal-architecture` | design | 9484 | ~2405 | ~71 | ~0 | ~764 | ~420 |
| 32 | `code-design-refactor` | refactoring | 9574 | ~2394 | ~59 | ~0 | ~0 | ~582 |
| 33 | `narration-writer` | frontend | 9551 | ~2388 | ~133 | ~0 | ~1881 | ~1386 |
| 34 | `website-concept-designer` | frontend | 9111 | ~2278 | ~121 | ~0 | ~0 | ~1197 |
| 35 | `manim-architect` | frameworks | 9003 | ~2251 | ~260 | ~284 | ~2768 | ~2929 |
| 36 | `test-designer` | quality | 8960 | ~2240 | ~107 | ~0 | ~692 | ~440 |
| 37 | `demo-deck-designer` | frontend | 8881 | ~2221 | ~101 | ~0 | ~1767 | ~4248 |
| 38 | `feature-planner` | workflows | 8612 | ~2221 | ~66 | ~0 | ~198 | ~394 |
| 39 | `ui-ux-architect` | frontend | 8787 | ~2197 | ~102 | ~411 | ~1378 | ~527 |
| 40 | `traefik-architect` | infra | 8664 | ~2166 | ~146 | ~210 | ~978 | ~514 |
| 41 | `performance-reviewer` | quality | 8632 | ~2158 | ~83 | ~0 | ~608 | ~1815 |
| 42 | `sql-architect` | databases | 8552 | ~2138 | ~89 | ~293 | ~0 | ~1522 |
| 43 | `security-reviewer` | quality | 8225 | ~2057 | ~77 | ~0 | ~535 | ~544 |
| 44 | `ddd-architect` | design | 7740 | ~1935 | ~77 | ~0 | ~0 | ~1277 |
| 45 | `nethttp-architect` | frameworks | 7521 | ~1912 | ~120 | ~401 | ~1606 | ~613 |
| 46 | `skill-builder` | meta | 7570 | ~1893 | ~71 | ~0 | ~2450 | ~492 |
| 47 | `fastapi-architect` | frameworks | 7437 | ~1860 | ~120 | ~419 | ~1014 | ~695 |
| 48 | `design-patterns` | refactoring | 7278 | ~1820 | ~64 | ~0 | ~0 | ~2638 |
| 49 | `protobuf-architect` | encoding | 7237 | ~1810 | ~83 | ~474 | ~912 | ~479 |
| 50 | `improve-codebase-architecture` | refactoring | 7204 | ~1801 | ~60 | ~0 | ~0 | ~4467 |
| 51 | `gin-architect` | frameworks | 6721 | ~1693 | ~108 | ~418 | ~1264 | ~594 |
| 52 | `grafana-architect` | infra | 6699 | ~1675 | ~76 | ~407 | ~1147 | ~622 |
| 53 | `docs-architect` | workflows | 6625 | ~1657 | ~143 | ~0 | ~0 | ~5815 |
| 54 | `notify-on-finish` | tooling | 5412 | ~1443 | ~135 | ~0 | ~0 | ~488 |
| 55 | `adhd-mode` | meta | 5613 | ~1404 | ~75 | ~0 | ~0 | ~343 |
| 56 | `codebase-hygiene-reviewer` | quality | 4733 | ~1184 | ~89 | ~0 | ~1031 | ~428 |
| 57 | `blast-radius-reviewer` | quality | 4337 | ~1085 | ~91 | ~0 | ~0 | ~327 |
| 58 | `tdd` | workflows | 4205 | ~1082 | ~69 | ~0 | ~0 | ~1781 |
| 59 | `demo-script-designer` | frontend | 3627 | ~937 | ~95 | ~0 | ~0 | ~793 |
| 60 | `commit-author` | workflows | 3566 | ~892 | ~81 | ~0 | ~0 | ~409 |
| 61 | `mockup-architect` | frontend | 3540 | ~885 | ~136 | ~0 | ~0 | ~474 |
| 62 | `unslop` | meta | 3119 | ~780 | ~94 | ~0 | ~0 | ~284 |
| 63 | `grill-with-docs` | workflows | 3104 | ~776 | ~63 | ~0 | ~0 | ~1947 |
| 64 | `guardrail-builder` | workflows | 3056 | ~764 | ~108 | ~0 | ~0 | ~311 |
| 65 | `logic-cleaner` | refactoring | 2340 | ~585 | ~93 | ~0 | ~0 | ~405 |
| 66 | `caveman` | meta | 1825 | ~457 | ~71 | ~0 | ~0 | ~272 |

**Totals:** 617295 body bytes · ~155068 body tokens · ~7409 desc tokens · ~178824 side tokens

## Topic files

Side files in a skill directory other than `STACK.md` and `RECIPES.md` — typically topic-specific reference docs referenced by name from `SKILL.md` (e.g. `PAYLOADS.md`, `AUTH_PATTERNS.md`). These contribute to the `~Topic tkns` column above and are loaded only when the parent skill chooses to read them.

- `go-architect` (~7940 tokens total):
  - `LAYOUTS.md` — ~7019 tokens
  - `DIGEST.md` — ~921 tokens
- `docs-architect` (~5815 tokens total):
  - `BACKLOG.md` — ~2334 tokens
  - `MIGRATION.md` — ~1071 tokens
  - `TEMPLATES.md` — ~1054 tokens
  - `FRONTMATTER.md` — ~741 tokens
  - `DIGEST.md` — ~615 tokens
- `rest-api-architect` (~4733 tokens total):
  - `CONCURRENCY.md` — ~1052 tokens
  - `IDEMPOTENCY.md` — ~864 tokens
  - `PAYLOADS.md` — ~785 tokens
  - `AUTH_PATTERNS.md` — ~764 tokens
  - `STATUS_CODES.md` — ~659 tokens
  - `DIGEST.md` — ~609 tokens
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
- `manim-architect` (~2929 tokens total):
  - `TRAPS.md` — ~1239 tokens
  - `FRAMEWORK.md` — ~1061 tokens
  - `DIGEST.md` — ~629 tokens
- `repo-tooling-architect` (~2672 tokens total):
  - `TASKFILE.md` — ~1943 tokens
  - `DIGEST.md` — ~729 tokens
- `design-patterns` (~2638 tokens total):
  - `PATTERNS.md` — ~2122 tokens
  - `DIGEST.md` — ~516 tokens
- `explainer-designer` (~2524 tokens total):
  - `VISUAL.md` — ~1212 tokens
  - `LEARNERS.md` — ~688 tokens
  - `DIGEST.md` — ~624 tokens
- `remotion-architect` (~2454 tokens total):
  - `EXPLAINERS.md` — ~1770 tokens
  - `DIGEST.md` — ~684 tokens
- `token-economy` (~2192 tokens total):
  - `GUARDRAILS.md` — ~961 tokens
  - `TOOLS.md` — ~656 tokens
  - `DIGEST.md` — ~575 tokens
- `grill-with-docs` (~1947 tokens total):
  - `CONTEXT_FORMAT.md` — ~853 tokens
  - `ADR_FORMAT.md` — ~686 tokens
  - `DIGEST.md` — ~408 tokens
- `docker-architect` (~1831 tokens total):
  - `CHECKLIST.md` — ~1083 tokens
  - `DIGEST.md` — ~748 tokens
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
- `sql-architect` (~1522 tokens total):
  - `ENGINES.md` — ~560 tokens
  - `POSTGRES.md` — ~498 tokens
  - `DIGEST.md` — ~464 tokens
- `narration-writer` (~1386 tokens total):
  - `CONVERSATIONAL.md` — ~732 tokens
  - `DIGEST.md` — ~654 tokens
- `ddd-architect` (~1277 tokens total):
  - `PATTERNS.md` — ~745 tokens
  - `DIGEST.md` — ~532 tokens
- `website-concept-designer` (~1197 tokens total):
  - `REFERENCE.md` — ~655 tokens
  - `DIGEST.md` — ~542 tokens
- `landing-page-designer` (~1134 tokens total):
  - `DIGEST.md` — ~596 tokens
  - `REFERENCE.md` — ~538 tokens
- `mcp-architect` (~977 tokens total):
  - `DIGEST.md` — ~977 tokens
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
- `api-contract-reviewer` (~617 tokens total):
  - `DIGEST.md` — ~617 tokens
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
- `cli-tool-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `ros2-architect` (~523 tokens total):
  - `DIGEST.md` — ~523 tokens
- `event-driven-architect` (~517 tokens total):
  - `DIGEST.md` — ~517 tokens
- `traefik-architect` (~514 tokens total):
  - `DIGEST.md` — ~514 tokens
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

- `mcp-architect` (~8753 body tokens)
- `go-architect` (~7548 body tokens)
- `backend-performance-architect` (~3393 body tokens)
- `repo-tooling-architect` (~3353 body tokens)
- `remotion-architect` (~3310 body tokens)
- `python-architect` (~3188 body tokens)
- `token-economy` (~3174 body tokens)
- `ci-cd-architect` (~3061 body tokens)
- `landing-page-designer` (~3057 body tokens)
- `connectrpc-architect` (~3028 body tokens)
- `latex-architect` (~3009 body tokens)
- `rsk-guide` (~2956 body tokens)
- `nextjs-architect` (~2952 body tokens)
- `react-architect` (~2950 body tokens)
- `astro-architect` (~2872 body tokens)
- `cli-tool-architect` (~2773 body tokens)
- `event-driven-architect` (~2749 body tokens)
- `docker-architect` (~2746 body tokens)
- `hugo-architect` (~2695 body tokens)
- `ros2-architect` (~2685 body tokens)
- `rest-api-architect` (~2646 body tokens)
- `explainer-designer` (~2638 body tokens)
- `api-contract-reviewer` (~2575 body tokens)
- `secrets-architect` (~2550 body tokens)
- `study-artifact-builder` (~2546 body tokens)
- `ml-paper-format` (~2537 body tokens)
- `observability-architect` (~2511 body tokens)

Heaviest 10 descriptions — each desc token is paid every turn for any session that installs the skill:

- `mcp-architect` (~300 desc tokens)
- `manim-architect` (~260 desc tokens)
- `explainer-designer` (~244 desc tokens)
- `go-architect` (~234 desc tokens)
- `backend-performance-architect` (~202 desc tokens)
- `landing-page-designer` (~165 desc tokens)
- `remotion-architect` (~154 desc tokens)
- `connectrpc-architect` (~153 desc tokens)
- `study-artifact-builder` (~149 desc tokens)
- `traefik-architect` (~146 desc tokens)

## Notes

- **Body tokens** cost only when the skill is invoked in a session.
- **Description tokens** cost every turn for any session that has the skill installed. The cost depends on which bundles are installed, not on the corpus total.
- Side files (`STACK.md`, `RECIPES.md`, topic files) are never auto-loaded; they cost 0 per turn.
- For exact counts: run each file through [`tiktoken`](https://github.com/openai/tiktoken) with `cl100k_base`.

# rsk landing page

The published page: React 19 + TypeScript + Vite, built as one prerendered
document. It replaced the hand-generated proposals in
`mockups/ui-proposals/2026-11-04/`, which stay as the design record.

## Commands

```bash
pnpm install          # needs NODE_AUTH_TOKEN (packages:read on GitHub Packages)
pnpm dev              # data, then vite on :5173
pnpm build            # tsc, vite build, then prerender dist/index.html
pnpm preview          # serve dist/
pnpm data             # regenerate src/data/catalog.json from the repository
pnpm check            # axe + overflow + reveals + JS-off against dist/
```

`pnpm dev` and `pnpm build` run `pnpm data` first (predev/prebuild), so the
numbers are never stale.

## Layers

```
src/data/        catalog.json (generated) + facts.ts (hand-checked constants)
src/lib/         useFilm (timeline) + useMotion (theme, reveals, counts, scroll)
src/components/  Header, Rail, Footer, Manifest        — page chrome and the device
src/components/ui/  Button, Tag, Panel, Section, Eyebrow, CopyBlock, FilmControl,
                    BrandMark, Wordmark               — the primitives
src/sections/    one component per page section
src/styles/      tokens, reset, layout, chrome, controls, panel, motion,
                 terminal, tradeoff + one sheet per page area
```

| Path | What |
|---|---|
| `src/data/catalog.json` | **Generated**, git-ignored. Skills, domains, bundles, releases, the repo's own `rsk.mod`/lock, the version. |
| `scripts/gen-data.mjs` | Reads `registry/index.json`, `skills/**`, `internal/config/catalog.toml`, `cmd/rsk/setup.go`, `internal/mcp/{tools,resource}.go`, `.rsk/rsk.mod`/`.lock` and the git tag. **Fails** when the tree drifts from 58 skills, 15 bundles, 12 `AddCommand`s or the three tool names. |
| `src/data/facts.ts` | Constants that cannot be derived: a real `rsk new` run, a real `project_profile` response, `rsk --help` groups, the quick start, the source list. |
| `scripts/prerender.mjs` | Renders the app to HTML into `dist/index.html` after the build, so the page reads and indexes without JavaScript. `main.tsx` hydrates it. |
| `scripts/verify.mjs` | The gate: axe at 8 combinations, overflow, unfired reveals, unsettled count-ups, JS-off, and a style-wiring check (mono stack, hero mask transition, theme button padding, panel radius) so a pruned or renamed selector cannot silently stop applying. Also runs each film (`#c-film`, `#a-film`): it must start on scroll, hide steps while playing and end fully revealed. |
| `scripts/compare.mjs` + `compare-film.mjs` + `compare-behaviour.mjs` | Parity against `mockups/ui-proposals/2026-11-04/direction-4-combined.html`: full-page captures and a per-element style dump (`WIDTHS=… THEMES=… WEB_URL=…` override the matrix), the films sampled at fixed times, and the filter, rail, theme, count-up, progress and copy behaviour. Static pixels alone missed two dead films, so run all three, and diff the PNGs at a strict threshold (about 2 levels): a loose one hid a tinted section. |

## Conventions

- **Tokens are semantic, styles read them.** `tokens.css` imports
  `@ralvarez-ts/tokens` and declares the page's roles on top: `--surface-page`,
  `--surface-raised`, `--surface-sunken`, `--text-strong`, `--text-muted`,
  `--text-accent`, `--border-subtle`, `--border-control`, `--status-ok|warn|danger|info`,
  `--radius-control|panel`, `--space-1..9`, `--duration-1..4`. A component never
  hard-codes a hex or a px radius.
- **Variants are data attributes**, not modifier class strings: `Button` takes
  `variant="primary|ghost"`, `Tag` takes `tone`, `Section` takes
  `tone="paper|sunken|ink"`. CSS reads `[data-variant]`, `[data-tone]`.
- **Class names name roles**: `manifest-file__line`, `call__message--request`,
  `agent-hero__line`, `filter__chip`. The old `k-`/`a-`/`b-`/`c-` prefixes came
  from the proposals' directions and said nothing about the page.
- **State classes are `is-*`** (`is-visible`, `is-filming`, `is-pending`,
  `is-match`), so a grep for `is-` finds everything JavaScript can toggle.
- **One sheet per component**, imported from `styles/index.css` in cascade order.
  Rules whose classes the prerendered page never uses are dropped; `verify.mjs`
  fails if a reveal never fires, which is how a missing rule surfaces.

## Decisions worth knowing

- **Tailwind is a build tool here, not a styling system.** The package's
  `@theme static` block only becomes real CSS variables when Tailwind's plugin
  compiles it. Without it the status roles (`--status-warn`, …) are undefined and
  `var()` falls back to nothing. Only `tailwindcss/theme.css` is imported; no
  utilities, no preflight.
- **`ui-core` is deliberately absent.** Its components are console-shaped
  (Card/Section with headers and control heights) and bringing it in means
  preflight over a bespoke art direction. Its conventions (tokens, variants,
  focus states) are followed by hand.
- **Sharp's geometry is overridden.** The package's sharp layer squares every
  corner; the page restores 4/6 px like the sibling landing pages.
- **The data step is the honesty mechanism.** Everything derivable is generated
  and asserted; everything hand-entered is marked `real run` / `real` / `sample`
  on the page.

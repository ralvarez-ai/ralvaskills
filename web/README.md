# rsk landing page

The published page: React 19 + TypeScript + Vite, built as one static document.
It replaced the hand-generated proposals in `mockups/ui-proposals/2026-11-04/`,
which stay as the design record and the screenshot set.

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

## What is where

| Path | What |
|---|---|
| `src/data/catalog.json` | **Generated**, git-ignored. Skills, domains, bundles, releases, the repo's own `rsk.mod`/lock, the version. |
| `scripts/gen-data.mjs` | Reads `registry/index.json`, `skills/**`, `internal/config/catalog.toml`, `cmd/rsk/setup.go`, `internal/mcp/{tools,resource}.go`, `.rsk/rsk.mod`/`.lock` and the git tag. **Fails** when the tree drifts from 58 skills, 15 bundles, 12 `AddCommand`s or the three tool names. |
| `src/data/facts.ts` | Hand-checked constants that cannot be derived: a real `rsk new` run, a real `project_profile` response, `rsk --help` groups, the quick start, the source list. |
| `src/styles/tokens.css` | `@ralvarez-ts/tokens` (`tokens.css` + `sharp.css`), the real package. No colour is copied. |
| `src/styles/base.css`, `a.css`, `b.css`, `c.css`, `d.css` | The page CSS, ported from the proposals. The palette block is aliases of the package (`--bg: var(--color-bg)`, `--line-2: var(--color-border-strong)`). |
| `src/lib/useFilm.ts` | The film engine: a pausable timeline, one play per view, final state = the markup. |
| `src/lib/useMotion.ts` | Theme, one-shot reveals, count-ups, reading progress, active-section rail. |
| `src/components/`, `src/sections/` | One component per section, chrome in `chrome.tsx`, the shared device in `Device.tsx`. |
| `scripts/prerender.mjs` | Renders the app to HTML into `dist/index.html` after the build, so the page reads and indexes without JavaScript. `main.tsx` hydrates it. |
| `scripts/verify.mjs` | The gate: axe at 8 combinations, overflow, unfired reveals, unsettled count-ups, JS-off. |

## Decisions worth knowing

- **Tailwind is a build tool here, not a styling system.** The package's
  `@theme static` block only becomes real CSS variables when Tailwind's plugin
  compiles it. Without it the status roles (`--warn`, `--ok`, `--danger`) are
  undefined and `var(--warn)` silently resolves to nothing. Only
  `tailwindcss/theme.css` is imported; utilities and preflight are not.
- **`ui-core` is deliberately absent.** Its components are console-shaped
  (Card/Section with headers and control heights) and bringing it in means
  preflight over a bespoke art direction. If a component ever fits, it is an
  import away.
- **Sharp's geometry is overridden.** The package's sharp layer squares every
  corner; the page restores 4/6 px like the sibling landing pages.
- **The data step is the honesty mechanism.** Everything derivable is generated
  and asserted; everything hand-entered is marked `real run` / `real` / `sample`
  on the page.

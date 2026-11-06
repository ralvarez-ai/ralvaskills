---
status: draft
project: ralvaskills
date: 2026-11-05
related: [designs/2026-11-04-ralvaskills-landing-page.md]
---

# Family chrome, the mockup preview bar, and the design-system refresh

Three linked pieces of work, in dependency order. Part 1 is a read-only comparison and is done;
parts 2 and 3 are proposals that need a decision before anything outside this repo changes.

## 1. Comparison: the landing pages today

Sites found: `ralvafi`, `ralvagym`, `ralvapass`, `ralvasend`, `ralvadocket` (all `site/`, React + Vite
+ Biome), plus this repo's `site/`. `ralvarez.dev` is Astro and not a product page.

| Piece | ralvafi (reference) | ralvagym / ralvasend | ralvaskills (here) |
|---|---|---|---|
| Family credit in the footer | `FamilyMark` (hexagon with the square-and-dot cut out) + "A ralvarez-apps product" linking to the org | same (`family-credit`) | **none**; footer shows the rsk `Wordmark` and a GitHub link |
| Header mark | app brand mark | not checked | `BrandMark` (two CSS squares), not the family hexagon |
| Side rail (section schema) | `Rail.tsx` | not checked | `Rail.tsx` (22 lines), present |
| Provenance block in the footer | "Where every claim comes from" `<details>` with a `SOURCES` list | not checked | `<details>` with `sources`; same idea, own copy |
| Tokens | own `tokens.css`, sharp/soft corners | own | own `tokens.css` |
| Brand kit asset | `ralvarez-brand-kit/ralvarez-apps/{sky,sky-lighter,rose}` | same | the kit has a `ralvarez-ai/` set (`rose`, `sky`, `sky-lighter`) that this page does not use |

Findings:

- This page should credit **"A ralvarez-ai product"** the way the apps credit ralvarez-apps. The
  hexagon/cut is shared; the brand kit already ships a `ralvarez-ai` variant to take the colour from.
- `FamilyMark.tsx`, `family-credit` CSS, the `<details>` source list and the `Rail` are copied per
  site. Five copies will drift. They are the first candidates for a shared package.
- ralvafi keeps the app's brand vectors in sync with `sync-shared.mjs` (copy at build, fail if the
  source is missing). That pattern is worth keeping for anything shared.

## 2. The mockup preview bar (ralvafi `mockups/landing/2026-11-04/a-palettes.html`)

What it is: a fixed bottom bar (`.pv`, `role="region"`), one `<fieldset>` per axis (palette, corners,
accent, logo mark), buttons with `aria-pressed`. State goes to `data-*` attributes on `<html>`,
persists in `localStorage` (`pv-<axis>`), and can be forced with a query string
(`?palette=b&accent=ink`), which is what lets `compare.mjs` render an exact variant. It reads the
page's own tokens, so it needs no styling of its own, and has a `forced-colors` rule. About 25 lines
of CSS and 10 of JS. It was in one ralvafi mockup (`a-palettes.html`); the first count of five across ralvafi,
ralvagym and ralvapass came from a loose grep and was wrong.

Proposal, smallest first:

1. **Extract it as a tiny asset, not an npm package.** `pv-bar.css` + `pv-bar.js`, where the JS
   builds the bar from a declarative config:
   `PV.init({ palette: ['a','b','c'], corners: ['soft','sharp'], accent: {teal:'#006a7f',…}, when: {accent: {palette:'b'}} })`.
   Mockups are self-contained HTML (they open from `file://`), so a build-time dependency does not
   fit; the file gets inlined or copied next to the mockup.
2. **Teach a skill to use it** (see below) instead of building a package registry for HTML.
3. Only if five or more repos need versioned updates, promote it to a package.

Which mockups get it: the ones that compare options (palette/corner/mark/layout variants). Pages
that show a single spec do not need it. Retrofit is a mechanical edit per file, so it should be
done per repo, in that repo, with its own screenshots, not from here.

**Skill:** a `mockup-architect` (name to confirm) next to `landing-page-designer`: the standards for
a mockup (self-contained, dated folder, `index.html` gallery, tokens as CSS variables, the preview
bar, `?param` forcing, reduced-motion, axe pass, README row) with `assets/pv-bar.{css,js}`. Alternative:
fold it into `landing-page-designer` as a reference. Separate is better: mockups are not only landing
pages (ralvafi has dozens of dashboards and flows).

## 3. Refresh `claude-design`

`ralvarezdev/claude-design/design-system/` is the ralvarez.dev personal brand extraction (one
`themes.css` of 161 lines, `theme-picker.js`, `colors_and_type.css`, 30 preview cards, one UI kit).
It does not yet know about:

- the product family (ralvarez-apps / ralvarez-ai): `FamilyMark`, the credit footer, rail, the
  ink-and-paper palette B, teal accent and the accent set (indigo, blue, plum, ochre, ink), sharp vs
  soft corners as a choice;
- the mockup preview bar as a component and the mockup standards;
- the brand-kit variants (`rose`, `sky`, `sky-lighter`) per org.

Proposal: add `ui_kits/ralvarez-products/` (family footer, header, rail, preview bar), add the new
palettes as `[data-palette]` blocks in `themes.css` rather than new files, add preview cards for
them, and say in the README that exploring new styles is expected: palettes and corner styles are
axes to vary, not decisions already taken. Run through the same `?param` forcing so each can be
screenshot.

## Risks

- Shared package vs copy: a shared UI package couples five sites' release cycles. Mitigation: share
  tokens and the three small components first, copy-at-build like `sync-shared.mjs`.
- `skills/**` publishes on push to `main`; a new skill ships as soon as it is merged. Mitigation:
  write it on a branch and confirm before merging.
- Editing other repos' mockups in bulk is wide; do it repo by repo, each with its own diff.

## Outcome of the retrofit (2026-11-05)

Converted: ralvafi 1 mockup, ralvapass 9 pages (theme axis), ralvasend `palettes.html` and `views.html`, ralvadocket 1
page; ralvagym none. Most option mockups are one file per option or JS-driven, which the bar does not fit without
reworking the generator (see the skill's retrofit checklist). Per-repo gaps are in each repo's
`docs/reviews/2026-11-05-family-chrome-gap.md` and the matching backlog items.

Checked: Playwright from `file://` on every converted page (query forcing, clicks, `localStorage`, 0 console errors),
ralvasend's existing axe audit of `palettes.html`, and `task lint` / Biome / `tsc`. Not checked: axe per variant, 390
and 1440 screenshots, a parser check on most pages, a browser look at this repo's footer or the skill's
`example.html`, and ralvafi's `a-palettes-check.mjs`. The `pv-bar.*` copies in ralvafi and ralvapass predate the
`swDark` fix. Items: `k3r9wd` (verification), `tzzc6d` (copy drift), `z7ckmp` (commit and publish).

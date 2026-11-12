# landing-page-designer digest

Designs landing-page proposals as self-contained HTML grounded in product code, real captures and brand: one memorable device per page, a choreographed hero film, a verification bar.

## Always-on rules

- Draw, do not refactor. Output is proposals under `<proposals-dir>/<date>/<product>/`; do not touch shipping app code.
- Do not commit or push unless the owner says so.
- Research first: facts table (claim, value, source file) from code, not marketing. Verify counts and names in code.
- Flag inconsistencies in NOTES.md; never fix silently. Label sample data on the page.
- Bar: a stranger gets what, who and why in five seconds; one product-only device; concrete claims; real pixels over illustration; state the trade-off.
- Keep each product's own mark, accent and tokens in both themes. Fonts: only the product's own, as woff2.
- Pick a personality before layout. Avoid the default skeleton (eyebrow, h1, lede, two buttons, trust list, glow).
- Hero film: 4-8 s, plays once, settles on a static final state that is plain HTML. Visible Pause/Replay. Plot only real behaviour.
- Motion: content visible by default. Opacity only in one-shot reveals. Scroll-linked effects animate transform, clip-path or stroke only.
- Captures: from the running product only, cropped with intent, both themes, separate phone crops. Crop out stale branding.
- Work in rounds: 3 small directions per product, owner feedback, one recommendation with its reason.

## Hard constraints

- Forbidden: stock photography, purple gradients, bento grids of grey cards with an icon each, "Trusted by 10,000 teams", invented testimonials, fake metrics.
- Each page: self-contained HTML, inline CSS and SVG, no CDN, under 200 KB excluding captures.
- Light and dark via `data-theme` with a toggle. One h1, skip link, focus-visible, 44 px targets.
- Before showing anyone: `verify.js` axe x8 (390x844 and 1440x800, light and dark, motion and reduced) zero violations; HTML parser check; smoke test; look at the screenshots.
- State what is not covered: screen readers, Safari/Firefox, real devices, user testing.

Invoke the full skill when: drawing, proposing or improving a product's public landing page, comparing directions, or asked to make it more attractive or dynamic. Not for porting a chosen design into the app.

Full guide: invoke the landing-page-designer skill.

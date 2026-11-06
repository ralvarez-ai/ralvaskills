# mockup-architect (digest)

Self-contained HTML UI mockups: dated folders, gallery index, token variables, and the fixed preview bar for pages that offer options.

## Always-on rules
- One `.html` opens from `file://`: inline or sibling CSS/JS, no build step, no CDN.
- Layout `mockups/<collection>/<YYYY-MM-DD>/` with an `index.html` gallery and a row in the collection README.
- Colours, radii, fonts are CSS variables on `:root`; variants are `:root[data-<axis>=<value>]`. Reuse the app's token names (`--bg --card --text --muted --line --line-strong --accent --btn-bg --btn-fg --shadow`).
- One question per page; a page offering options says what is decided and recommends one.
- Numbers are real or marked sample.
- Pages that offer options use the skill's pv-bar.css and pv-bar.js, copied next to the mockup. Axes go in the `pv-config` JSON before pv-bar.js. State: `<html data-*>`, `localStorage pv-*`, `?key=value` (query wins, not persisted). Add class `pv-clear` to the last section.
- Axes only for what the owner can choose (palette, corners, accent, logo mark, density, theme).
- Check contrast of every variant, not only the default.
- Retrofit: change the generator and regenerate, never only the output.
- Verify: parser check, axe per variant and theme, 390px and 1440px screenshots per variant, no console errors, works from `file://`.

## Hard constraints
- A mockup never imports the app's code; mockups never touch the app's source.
- Superseded rounds move to `archive/`; do not edit them.
- Never click in a script; render an exact variant by `?query` URL.

Invoke the full skill when: drawing, comparing or retrofitting a mockup, or when a mockup offers options to choose between.

Full guide: invoke the mockup-architect skill.

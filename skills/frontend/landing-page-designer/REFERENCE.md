# landing-page-designer reference

Lookup material moved out of SKILL.md.

## Traps (SKILL.md section 9; each cost time once)

- CSS specificity: a theme-switch rule for images must beat component rules (`img.sl`, not `.crop img`).
- Measuring a hidden element gives 0 (lens, zoom): show, then measure.
- Opacity-dimmed inactive steps fail contrast: emphasise with border/background instead.
- `role="list"` needs `listitem` children; section titles must not jump from `h1` to `h3`.
- Glow/mesh backgrounds need full-bleed, a mask, and `body{overflow-x:clip}`.
- An unbreakable `<code>` in a grid cell causes horizontal scroll at 390; a scrollable `pre` fails `scrollable-region-focusable` (give it `tabindex="0"`, `role="region"` and a unique label, or use `pre-wrap`).
- Generic class names collide across shared and product CSS (`.tl`, `.copy`): namespace family classes.
- `cqw` resolves against an ancestor container, never the element itself.
- An `aria-label` that differs from visible text breaks label-in-name; put extra context in an `.sr` span.
- A Python template string missing its `f` prefix leaks `{label}` into the accessible name: grep the output.
- `loading="lazy"` images are missing from full-page screenshots; not a defect.
- Pages assembled by extraction do not update when a source changes until the assembler is re-run.

## Files and tools (SKILL.md section 11)

```
<proposals-dir>/<date>/
  HANDOFF.md                  brief for the next session
  <product>/direction-N.html  pages; NOTES.md; fonts/; shots/
  _tools/verify.js            axe x8, overflow, images blocked, JS off
  _tools/family.py            shared shell, tokens, motion kit (film, count, type, onView)
  _tools/build_<product>.py   one generator per product
  _tools/settled.js           screenshot after the film settles
  _tools/standalone.py        viewing-only copy with inlined fonts and captures
  _tools/smoke_<product>.js   interaction smoke tests
```

Generators carry hard-coded paths from the session that made them; fix `sys.path` and output paths before running. `verify.js` expects `axe-core@4.10` installed in a sibling `verify/node_modules`.

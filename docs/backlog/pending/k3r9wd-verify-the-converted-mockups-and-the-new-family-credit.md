---
title: 'Verify the converted mockups and the new family credit'
priority: p2
section: ralvaskills
topics:
    - site
    - mockups
    - verification
---
The 2026-11-05 retrofit checked behaviour (query forcing, clicks, `localStorage`, 0 console errors, Playwright from
`file://`) but not most of the visual and accessibility gates, and nobody has looked at this repo's new footer in a
browser.

Not done: axe per variant on the converted pages; 390 and 1440 screenshots in light and dark; a parser check on
`skills/frontend/mockup-architect/assets/example.html`; a browser look at `site/` with the family credit; ralvafi's
`_tools/a-palettes-check.mjs` (it hung and was killed).

Done when each converted page has an axe run at 0 violations per variant and a light/dark screenshot pair, the
footer is checked in a browser at both widths, and the ralvafi check script finishes or is replaced.

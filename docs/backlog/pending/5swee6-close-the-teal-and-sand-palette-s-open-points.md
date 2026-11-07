---
title: 'Close the Teal & Sand palette''s open points'
priority: p2
section: landing
topics:
    - palette
    - accessibility
---
Teal & Sand replaced Ember on 2026-11-06 (`site/src/styles/tokens.css`, `terminal.css`, `favicon.svg`; chosen from `mockups/ui-proposals/2026-11-06/palettes-round-2.html`). The generator checked WCAG 2.2 AA for every pair and axe found 0 violations on the mockup, and `task site:check` passes on the site. Not validated: focus rings on the teal accent, `forced-colors`, tinted backgrounds beyond what axe measures, and the full page of each variant (only the first screen was reviewed in the contact sheets). The 4 to 10 level residue in the light `this-repo` band, found under Ember, is unchecked under the new palette. The brand question is the owner's: the rsk wordmark is now teal, no longer blue or vermilion.

Done when focus and forced-colors are checked in a browser (screenshots), the residue is explained or gone, and the wordmark decision and the palette change are recorded in `docs/designs/2026-11-04-ralvaskills-landing-page.md`.

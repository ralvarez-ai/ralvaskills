---
title: 'Decide whether the What the agent is told section comes back'
priority: p2
section: landing
topics:
    - design
    - owner-decision
---
It exists in `direction-3.html` and the web build had it, but the combined mockup has no such section, so it was removed for parity. If it should return, add it to the mockup's section list (`combined.py`) and to the rail together, regenerate, then port it to `site/` and extend `verify.mjs` and the comparison scripts. If not, record that decision in the design doc.

Done when the owner has decided and the design doc says which.

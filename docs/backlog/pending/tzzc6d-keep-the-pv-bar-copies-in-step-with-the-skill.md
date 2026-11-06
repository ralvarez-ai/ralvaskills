---
title: 'Keep the pv-bar copies in step with the skill'
priority: p3
section: ralvaskills
topics:
    - skills
    - mockups
---
`mockup-architect/assets/pv-bar.*` is the source. Copies sit in dated mockup folders in ralvafi, ralvapass, ralvasend and ralvadocket; the ralvafi and ralvapass ones predate the `swDark` fix. Add a small check (hash compare) or a documented refresh step, and decide whether the copies belong at the collection root instead of per dated folder.

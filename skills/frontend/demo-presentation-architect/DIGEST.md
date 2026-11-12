# Demo Presentation Architect (digest)

Designs a slide-deck spec as one `.md` file: slide budget, narrative arc, layouts from a fixed catalog, takeaway titles.

## Always-on rules
- Interview first, in order: language, then core info (topic, audience, goal, duration/slide budget, source material, constraints). Propose the outline; expand to full slide specs only after approval.
- Slide budget: ~90-120 s per content slide; subtract fixed slots.
- Fixed slots: `cover` opens; `section-divider` when >=3 sections (skipped for decks of 8 slides or fewer); `recap-numbered` (>=3 takeaways) and/or `closing-cta` close; a live-demo handoff ends with `handoff-band`.
- One idea per slide. Split on an "and"/comma joining unrelated phrases, word-budget overflow, two card-grids, or a transition that needs a new idea.
- A concept appears on at most two slides: introduction and recap.
- Title carries the takeaway, not the topic; exactly one `<em>` accent on the takeaway word.
- Eyebrow names the section, not the slide. Reading flow: eyebrow, title, body, callout, footer.
- Body word budgets: cover/quote <=40; title-bullets/card-grid <=120; hub-spoke/showcase-split <=140; recap/stat-row <=80. Over budget: move detail to notes or split.
- Speaker notes <=40 words: why and timing only, never slide text.
- Final copy, no placeholders. Headlines under ~10 words; parallel bullets; image fields carry `alt:`.
- Every slide has a Visual hierarchy block and a Transition line (except the last).
- Define a term before a slide uses it.

## Hard constraints
- Output is always one `.md`. Refuse HTML, PDF, PPTX, images, styling code; offer the `.md`.
- Never invent a layout name. If none fits, follow LAYOUTS §5 and wait for approval.
- Never mix languages unless the user asks for bilingual.
- Never switch narrative arc mid-deck.
- Presenter narration is a separate deliverable: hand it to demo-script-architect.

Invoke the full skill when: writing or revising slide-by-slide deck specs, picking layouts, or running the pre-handoff checklist.

Full guide: invoke the demo-presentation-architect skill.

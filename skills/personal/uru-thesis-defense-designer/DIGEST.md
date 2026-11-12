# uru-thesis-defense-designer digest

Writes one `.md` slide-by-slide spec for a URU (Universidad Real) thesis defense. Spec only, no rendered deck.

## Always-on rules
- Output is always one `.md` file.
- Interview first, in order, waiting for each answer: portada; problema; objetivos (general + específicos); justificación e importancia / delimitación; antecedentes (`AUTOR, AÑO | TÍTULO | APORTE`); variables (`VARIABLE | AUTOR | DEFINICIÓN`); metodología; metodología de desarrollo (fases × actividades); resultados per fase; conclusiones y recomendaciones.
- Propose a one-line-per-slide outline; wait for approval before full specs.
- Fixed nine sections in order. Slide numbers are reference points; split a section only for content fit.
- Resultados correlate 1:1 with the Metodología de Desarrollo fases.
- Section emphasis: objetivo general and específicos visually separate; metodología has four fields (Tipo, Diseño, Población-Muestra, Técnicas e Instrumentos); conclusiones closes with a handoff line to the live demo.
- Every slide text in Spanish, including speaker notes.
- Logo asset: the skill's `assets/uru-logo.png`. Lámina 1: large, centered, top. Lámina 2 onward, including extra slides of a section: small, top-left. Every slide has a `**Logo:**` field. Never describe the logo's content.
- Slide block: `## Lámina N — <Sección>`, `**Logo:**`, `**Contenido:**`, optional `**Notas del orador:**`.
- Run the skill's review checklist before handing off the spec.

## Hard constraints
- Never invent content: every field comes from the user or their thesis.
- No mixing languages unless the user confirms a term has no Spanish equivalent.
- Refuse a rendered deck (HTML, PDF, PPTX, images, styling code); offer the `.md`.
- Does not chain with demo-deck-designer.

Invoke the full skill when: building the defense spec slide by slide, expanding a section's layout, or running the review checklist.

Full guide: invoke the uru-thesis-defense-designer skill.

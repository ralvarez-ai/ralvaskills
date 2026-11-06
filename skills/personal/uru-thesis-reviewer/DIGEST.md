# URU Thesis Reviewer (digest)

Iterative feedback on URU thesis `.docx`/`.md`: numbered `.md` diff files per session; the author applies them by hand.

## Always-on rules
- Output goes in `<output-base>/YYYY-MM-DD/`; same-day extra sessions get `-a`, `-b`.
- Files kebab-case, numbered in reading order: `01-index.md`, `nn-section.md`, `11-referencias.md`.
- `01-index.md` is mandatory every session. Use `00-handoff.md` for open-work handoffs only.
- Diff shape: `-` exact old text, `+` proposed replacement. For substantive issues, `+` is a prompt, not rewritten text.
- Tag every block `[blocking]` (norm violation or factual error), `[suggest]` (real improvement), or `[nit]` (taste).
- Cite the rule when URU norms drive it (e.g. `NORMAS §V. CITAS`). RAE for language. APA 7th where URU is silent.
- Issue IDs `yyyy-mm-dd-<chapter>-<slug>`, assigned at first discovery, never restated as new.
- Status: `abierto`, `resuelto-en-borrador`, `aplicado`, `rechazado`, `needs-decision`. Severity and status are separate.
- Carry forward the previous `01-index.md` Seguimiento table at session start. Keep every issue not `aplicado`/`rechazado` unless its status moved.
- A reversed prior decision needs a one-line `**Por qué se revirtió:**` on its row.
- Pair any tutor comment number with its verbatim quoted text.
- Bibliography: propose a concrete replacement via search, or mark `[suggest] needs stronger source — none found in search`.
- Review order: substantive (structure, then content) before prose. Defer prose on a section with a substantive `[blocking]`.
- Feedback prose in Spanish; tags, headers and metadata in English.
- Ingest `.docx` via the docx skill, else pandoc into scratch `.review.md` (never in the session folder).
- Report to user: session path, severity totals, starting point. Never paste diffs into chat.

## Hard constraints
- Never edit, rewrite or replace the `.docx`.
- Only the author moves an issue to `aplicado`.
- No section rewrites: targeted diffs only. A total rewrite is `[blocking]` with reason.
- Never claim plagiarism; flag as `[suggest] verificar fuente`.
- Not grading; not defense/jury prep.
- Ask for an output path if only the `.docx` is given.
- Never translate Resumen to Abstract; flag mismatches only.

Invoke the full skill when: reviewing a URU thesis or chapter, writing session folders or diffs, checking citations or bibliography against the URU norms file, or carrying issues across sessions.

Full guide: invoke the uru-thesis-reviewer skill.

# explainer-architect digest

Designs teaching videos independent of the engine that draws them: scene rule, provenance, cues, versions, look check.

## Always-on rules

- The drawing teaches; the formula, code or table sits beside it as the summary. One idea per scene, drawn in the order the narration says it.
- Cue drawings to narration phrases (words before the phrase / measured rate), and fail the build when a phrase disappears. Scene length is `max(animation, words / rate + lead-in + tail)`; count spoken words.
- Every scene has an origin chip: SOURCE, TEXTBOOK, DERIVED, REAL RUN or ILLUSTRATION. Made-up values are flagged on screen and in the words. The last scene says what was not checked.
- Facts live in a content file with a source per fact; no number on screen is missing from it.
- Versions: old is v1, new is v2, comparisons are v2.a and v2.b. A new look is an option whose default is the old one.
- Before handoff: a contact sheet of end-state frames, the dense scenes at full size, the length gate with real audio, all at the final resolution.
- Interview first: audience, the one thing they must take away, sources and which are real, engine, who supplies the voice.

## Hard constraints

- Never overwrite a render someone has seen or published.
- Never show a made-up value without a label.
- Never time scenes from a guessed speech rate once the voice is measured.
- Never choose the voice provider for the owner.

Invoke the full skill when: planning, storyboarding or reviewing an explainer video, whatever engine draws it.

Full guide: invoke the explainer-architect skill.

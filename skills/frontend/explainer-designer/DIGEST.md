# explainer-designer digest

Designs teaching videos independent of the engine that draws them: scene rule, provenance, cues, versions, look check.

## Always-on rules

- The drawing teaches; the formula, code or table sits beside it as the summary. One idea per scene, drawn in the order the narration says it.
- Cue drawings to narration phrases (words before the phrase / measured rate), and fail the build when a phrase disappears. With audio, scene length is `voice + lead-in + tail + hold` from the measured file; a scene with no audio uses `max(animation, words / rate + lead-in + tail)`. Count spoken words.
- Every scene has an origin chip: SOURCE, TEXTBOOK, DERIVED, REAL RUN, ILLUSTRATION, UNCHECKED or PROPOSED. Made-up values are flagged on screen and in the words. Say what was not checked beside the claim and in the content file; a closing scene only if the owner asks.
- Facts live in a content file with a source per fact; no number on screen is missing from it.
- Versions: old is v1, new is v2, comparisons are v2.a and v2.b. A new look is an option whose default is the old one.
- Before handoff: a contact sheet of end-state frames (about 97%), the dense scenes at full size, the length gate with real audio, all at the final resolution; text readable where the video is shown.
- Interview first: audience, the one thing they must take away, sources and which are real, what exists and what is only designed, engine, who supplies the voice.
- Visual level (VISUAL.md): one focus at a time with the rest faded, a colour means one thing and a legend says it, solid is real and dashed is outside or not built, motion travels, labels stay readable on a phone.
- Define each domain word before using it, with a picture. A flow is drawn moving, a list is a table. Tag parts that are not built (BUILT, STILL TO BATTLE TEST, DESIGNED, NOT BUILT, NOT BUILT) and name any stand-in.

- For learners: one worked number per concept, predict then reveal, a pointer to a tool is an option off by default (LEARNERS.md).

## Hard constraints

- Never overwrite a render someone has seen or published.
- Never show a made-up value without a label.
- Never show a part as built when only a stand-in answers.
- Never time scenes from a guessed speech rate once the voice is measured.
- Never choose the voice provider for the owner.

Invoke the full skill when: planning, storyboarding or reviewing an explainer video, whatever engine draws it.

Full guide: invoke the explainer-designer skill.

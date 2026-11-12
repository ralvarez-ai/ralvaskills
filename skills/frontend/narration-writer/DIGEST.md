# narration-writer digest

Writes explainer-video narration and on-screen text in a controlled, STE-style English (about 80% of ASD-STE100), with a glossary, spoken forms, per-scene word budgets and a stdlib lint script.

## Always-on rules

- Narration is data: a script table of scenes (id, seconds, narration, caption, cue). The composition plays audio generated from it and never contains the words.
- One idea per sentence: 20 words or fewer (25 for a description), active voice with a named actor, simple tenses, plain verbs, no contractions, no idioms. Confirm the numbers in the official text.
- One word per meaning: glossary first, same term in every scene, a spoken form (`say`) for each acronym and mixed name.
- Plan at 2.5 words per second: `floor((seconds - lead-in - tail) * rate)`. Then measure the real voice on the slowest scene (not the mean, spoken words) and set that rate in the script; it times scenes and cues alike.
- The screen carries names, numbers and code; the voice carries the reason. No symbols, paths or identifiers in speech.
- Say what is not real: where the picture is an illustration, the narration says so.
- Run `scripts/narration_lint.py` (add `--audio-dir` for the length gate) before declaring a script done. Errors block; warnings are heuristics.
- Interview first: length, scene table, audience, language, who supplies the voice, what is real, sources.

- For a learner lesson the dial may be conversational (about 30 words a sentence, connectors allowed, glossary kept); see CONVERSATIONAL.md.

## Hard constraints

- Never claim ASD-STE100 compliance: no dictionary check is done, and the name is a trademark. Never copy the specification or its word list.
- Never say more than the sources behind the picture support: no timing, safety or score claims unless the source states them.
- Never speed up audio or trim it to hide an over-budget script; cut words, or lower the rate to the measured one when only the guess was wrong.
- Never pick the voice provider for the owner. Keys stay in env files that git ignores.
- Never type narration into a composition or call a speech engine at render time.

Invoke the full skill when: writing voiceover or captions for a video, simplifying dense technical prose, or fitting narration to scene timing.

Full guide: invoke the narration-writer skill.

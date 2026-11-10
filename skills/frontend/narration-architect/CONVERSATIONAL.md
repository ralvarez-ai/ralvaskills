# Conversational register and local voices

A third position of the dial in [SKILL.md](SKILL.md) section 2, and notes on offline voices. Use it when the viewer is a learner who will watch once, not a user who repeats steps.

## 1. When to leave the 80% dial

Short sentences of 20 words or fewer, one after the other, sound clipped when read by a synthetic voice: every sentence has the same length and the same fall. For a lesson that explains *why*, that rhythm reads as robotic. The controlled dial is right for steps; a lesson needs reasons, and a reason needs a connector.

## 2. The conversational dial

- **Limits:** about 30 words per sentence and 8 per paragraph. Set them in the script file (the lint script reads them from there), so the same lint runs on the looser text.
- **Keep:** one meaning per word and the glossary, active voice, no idioms, no jokes, plain punctuation, no symbols in speech, claims that trace to a source.
- **Allow:** "so", "because", "which is why", "now", "but", a question that the next sentence answers, and a sentence that tells the viewer what to look at ("with the formula on the screen").
- **Vary the length.** Put a short sentence after a long one. Two or three sentences of similar length in a row is the sign of the clipped rhythm.
- **Open with the problem, not the name.** "None of these points has a label" before "K-Means".
- **Say the reason once.** The voice gives why; the picture gives the number.
- **Name the lint level.** Write "conversational (30/8)" in the script header. Do not call it STE: it is no longer, and the name is a trademark.

## 3. A worked line

Controlled: "Measure the distance from each point to each centroid. Each point joins the nearest centroid."

Conversational: "Take one point and measure its distance to each centroid. The smallest distance wins, so the point joins the nearest one. Do this for every point, and the plot splits into regions."

The second states the rule, the consequence and the result, and it runs while the drawing shows a single point, then all of them.

## 4. Local (offline) voices

- A free local engine is good for drafts and for a language the paid service lacks. Keep its audio in a separate folder and write a separate output file (`.local.voiced`), so it never overwrites a paid take.
- **An echo or a metallic tail** usually comes from the model, not from the mux. Test by playing the raw clip before it is placed in the video. If the raw clip is clean, look at the mix (`amix` with `normalize=0`, no overlap of two clips of one scene); if it is not, change the voice or the speaking speed, and compare two or three before choosing.
- **Measure the rate on the slowest scene** with the real engine (section 4 of the main file) and set it in the script file; a local voice often differs from the paid one, so a rate measured on one does not carry over.
- Record the engine, voice and speed in the script header.

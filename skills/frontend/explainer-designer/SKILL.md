---
name: explainer-designer
version: 2.2.0
description: Engine-independent design of teaching videos (explainers for a mechanism, formula, protocol or data structure) - the drawing teaches and the formula sits beside it, one idea per scene, drawings cued to narration phrases, every scene marked as source, textbook, derived, real run or illustration, facts held in a content file with sources, words defined before they are used, a scene shaped to its idea (a flow is drawn moving, a list is a table), maturity tags for parts that are built, untested or only designed, v1/v2 and A/B versions that never overwrite a render, and a contact-sheet look check before handoff. Use when planning, storyboarding or reviewing an explainer video, whatever draws it (Remotion, Manim or another tool).
---

# Explainer Designer

Decides what a teaching video shows and in which order, and how to check it. It produces the scene plan, the content file and the review, not the code of an engine: for Remotion see [remotion-architect](../../frameworks/remotion-architect/SKILL.md), for Manim see [manim-architect](../../frameworks/manim-architect/SKILL.md). The words and the speech rate belong to [narration-writer](../narration-writer/SKILL.md); a live, presenter-led demo is [demo-script-designer](../demo-script-designer/SKILL.md). Templates and commands are in [RECIPES.md](RECIPES.md); videos for learners (a worked number per concept, predict then reveal, a pointer to a tool as a switch) are in [LEARNERS.md](LEARNERS.md); how a scene looks and moves (focus, colour meaning, line language, type, motion) is in [VISUAL.md](VISUAL.md).

## 1. Mental model

- **The drawing teaches; the formula, the code or the table sits beside it as the summary.** A scene that opens with a formula and then decorates it teaches nobody. Draw the object first (a graph being explored, a timeline, a table that fills in, a token that walks a path). The formula appears when the drawing has shown where each of its parts comes from.
- **One idea per scene.** A scene that needs "and also" is two scenes.
- **The picture and the voice split the work.** The picture carries names, numbers and code; the voice carries the reason (narration-writer, section 1).
- **A claim has a source, or it is labelled.** A viewer must be able to tell a fact from a worked example from an invention.

## 2. Interview first

Ask in order and wait for each answer: who watches and what they know already; the one thing they must be able to do or say afterwards; the sources (documents, runs, textbooks) and which are real; length and where the video is shown; which engine draws it (that picks the skill that follows); who supplies the voice (never choose it for the owner). Then ask what exists: which parts are built, which only designed, which never tested, and what a viewer would wrongly assume exists; and which domain words must be defined before the video uses them.

## 3. The scene rule

1. **Draw in the order the narration says it.** If the voice says A and then B, A appears first. Reorder the words before you reorder the animation.
2. **Cue drawings to phrases, not to seconds.** Keep the narration as data; an item (a row, a card, a step) carries the phrase at which it appears, and its time is the place of the phrase in the text times the measured length of the voice, minus a short lead so the item is on screen as it is named. The cue must fail the build when its phrase is no longer in the narration. A time typed by hand is only the fallback for a scene with no audio.
   - **A recorded clip (a rendered animation) is cued the same way:** cut the voice into one narration per part of the clip and start each audio when its part starts (the clip records where its parts start). One audio for a long clip drifts: in one video the voice ran 13 s ahead of the picture and said "this is a bounded model" over a different drawing.
3. **Scene length comes from the audio.** When the audio exists, a scene lasts `voice + lead-in + tail + a short hold`, read from the measured length of each file (record it in the audio manifest); a scene with no audio keeps a planned length from `max(animation time, words / rate + lead-in + tail)` with the measured rate (narration-writer, section 4). Hand-set seconds leave 30% to 60% of a video silent, and change with every voice. Count spoken words (after pronunciation forms) when you estimate.
4. **Leave silence.** A new drawing gets a second or two before the voice talks over it, and dense text is never read at the same time as it appears.
5. **End on a stable frame.** The last frame of a scene is the finished drawing, not half an animation.
6. **Do not draw twice what is said once.** If a number is on screen, the voice gives its meaning.

## 4. Provenance

Every scene carries a small chip with the origin of its content:

| Chip | Meaning |
|---|---|
| `SOURCE` | what the cited sources state |
| `TEXTBOOK` | general background; the reference is in the content file |
| `DERIVED` | worked out here from a source's data or formula (a calculation, not the source's claim) |
| `REAL RUN` | taken from a recorded run |
| `ILLUSTRATION` | made up to show an idea; it is not data |
| `UNCHECKED` | a claim nobody verified; the scene says so |
| `PROPOSED` | a design or a plan; the thing shown does not exist yet |

- Made-up values are flagged on screen **and** in the words ("this graph is made up").
- Facts live in a **content file** with a source per fact; the scene reads it. No number appears on screen that the content file does not hold. Shape in [RECIPES.md](RECIPES.md#1-the-content-file).
- Say what was not checked or not run beside the claim it qualifies, and keep the full list in the content file (`unchecked`). A closing scene that lists it exists only when the owner's rule allows one; ask in the interview.

**Maturity of a part.** The chip says where a scene's content comes from; a maturity tag says whether the thing shown exists. Tag the part, not the video:

Tags: `BUILT` (built and tested, even if little used), `STILL TO BATTLE TEST` (built, not yet in a scored run or against the real dependency), `DESIGNED, NOT BUILT`, and `NOT BUILT` (the diagram shows something that does not exist; name what stands in for it).

- Name the stand-in in the words and on screen: "a checker confirms it" is false when only a test double answers.
- A figure from one task is shown as one task, with its source, and says whether it was re-run here.
- When the subject is a system under construction, one scene may sort its parts into these columns. It is not the last scene.

## 5. Words first, then the shape of a scene

- **Define each domain word the first time it appears, with a picture.** Say what fills the role, not only what it is called: a "worker" may be an agent, a program that drives a model, or a person. A short scene near the start that defines the nouns, and shows what the system sees of each, saves every later scene from explaining them.
- **Give each scene the shape of its idea:** a list of facts is a table or cards; a sequence, a request and its answer, or a check that can fail is drawn moving; parts that sit together are layers; amounts are bars; what exists and what does not is columns ([RECIPES.md](RECIPES.md#7-the-shape-of-a-scene)).
- How the scene looks and moves (one focus at a time, a colour that means one thing, solid and dashed lines, motion that travels) is in [VISUAL.md](VISUAL.md).
- A scene made only of text and a table needs a reason. When the narration says "then", the drawing moves.
- A video of more than about five minutes keeps an orientation device: a small map of the system with the current part lit.

## 6. Versions of a video

- **Never overwrite a render that someone has seen or published.** The old one keeps its name (`v1`), the redesign is `v2`, a comparison is `v2.a` and `v2.b`. The same for the narration folder and the content file: copy, then change.
- Keep the old composition or scene id and add the new one beside it, so a reviewer can play both.
- A new look is an option whose default is the old look; existing renders do not move.
- Audio that was generated for the same words is reused by a new version; changed words need new audio.

## 7. Look before you hand off

1. Render, take one frame near the end of each scene (about 97%: at 85% a late fade or cue looks broken when it is not) and tile them in a contact sheet ([RECIPES.md](RECIPES.md#2-the-contact-sheet)). Read it for overlap, empty areas, near colours, dense text and a stamp or chip that covers a label. Before you change a tile that looks wrong, open the last second of its scene.
2. Open the dense scenes at full size; a sheet hides small type. The smallest text must be readable where the video is shown (a phone, a projector): fewer words and larger type beat a small label. An empty third of the frame means the drawing can be larger. In a flow scene something moves along the arrows; a fade of the whole arrow shows no direction.
3. Run the length gate with the real audio: none runs past its scene.
4. Listen once, and check each sentence against its source.
5. Judge at the final resolution, not at the preview size.

Videos of recorded runs: the steps to pick, verify and compare the data are in [RECIPES.md](RECIPES.md#5-a-video-of-recorded-runs).

## 8. Anti-patterns

- A formula first and a drawing as decoration.
- A made-up value with no label.
- Scene times from a guessed speech rate.
- A cue given in seconds that no one updates when the words change.
- A scene length typed by hand next to audio that has a measured length.
- One long voice over a clip with several parts.
- A new render saved over the old one.
- Two ideas in one scene, or a scene with no stable end frame.
- Reading the screen aloud.
- Judging the look from a thumbnail.
- A flow drawn as a table.
- A word used before it is defined.
- A part shown as if it existed when only a stand-in answers.
- A figure from one task shown as the result.
- Text smaller than the viewer can read.

## 9. Review checklist

- [ ] Each scene has one idea, drawn in the order it is said.
- [ ] Each scene has its origin chip; every illustration is flagged on screen and in the words.
- [ ] Every number on screen is in the content file with a source.
- [ ] Cues point at phrases, and the build fails if a phrase disappears.
- [ ] Scene lengths come from the audio of the take being rendered; silent time after the voice is a hold of a second or so, not a third of the video.
- [ ] A clip scene has its voice cut at the parts of the clip.
- [ ] The speech rate is measured, and no audio runs past its scene.
- [ ] What was not checked is said beside the claim and in the content file; a closing scene only when the owner asks.
- [ ] Every domain word is defined before it is used; each scene has the shape of its idea.
- [ ] Every part that is not built carries its maturity tag, and a stand-in is named.
- [ ] A figure says its source, and whether it was re-run.
- [ ] The squint, phone and grey checks of VISUAL.md pass: one focus, readable labels, no meaning carried by colour alone.
- [ ] Old versions are untouched; the new render has its own name.
- [ ] A contact sheet of end-state frames (about 97%) and the dense scenes were viewed at final size; the smallest text is readable where the video is shown.

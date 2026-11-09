# Narration Architect: recipes

Templates and mechanics for [SKILL.md](SKILL.md). Nothing here is tied to one project.

## 1. The script file

One JSON file per video. The lint script reads it, and the audio generator reads the `say` text. Keep the header, so a render can be repeated.

```json
{
  "video": "scan-cycle-explainer",
  "voice": {"provider": "chosen by the owner", "id": "voice id", "settings": "model, speed, stability"},
  "rate": 2.5, "leadIn": 0.5, "tail": 0.5,
  "limits": {"sentence": 20, "paragraph": 6},
  "glossary": [],
  "scenes": [
    {
      "id": "intro",
      "seconds": 10,
      "cue": "Title card at 0 s; diagram at 3 s",
      "caption": "Scan cycle",
      "narration": "A PLC runs its program in a scan cycle. The cycle repeats without stop."
    }
  ]
}
```

- `seconds` comes from the video's timeline table, not from the script. If the timeline changes, change it here too.
- `caption` is the on-screen text. `narration` is the spoken text. They may differ: the caption holds the exact identifier, the narration holds the spoken form.
- Never put a key or token in this file. Keep it in an env file that git ignores.

## 2. Glossary and pronunciation

One entry per term. The lint script reads `term` and `avoid`; the audio generator reads `say`.

```json
{"term": "scan cycle", "avoid": ["scan loop", "PLC cycle"], "say": "scan cycle"}
{"term": "PLC", "avoid": ["controller"], "say": "P L C"}
{"term": "STruC++", "avoid": [], "say": "S T r u C plus plus"}
```

- One meaning per term. If a word needs two meanings, it needs two terms.
- Apply `say` to the text just before it goes to the engine, and keep the narration itself readable.
- Listen to every acronym once with the real engine. Engines differ.

## 3. Budget math

`words allowed = floor((seconds - leadIn - tail) * rate)`

| Scene seconds | Lead-in and tail | Rate | Words allowed |
|---|---|---|---|
| 10 | 0.5 + 0.5 | 2.5 | 22 |
| 20 | 0.5 + 0.5 | 2.5 | 47 |
| 8 | 0.5 + 0.5 | 2.5 | 17 |

- 2.5 words per second is a planning rate. Measure the real voice on a sample scene and replace it.
- Count the words of the spoken text, after the `say` forms are applied.
- A cue-heavy scene (animation the viewer must watch) needs a lower rate. Use a smaller `rate` for that scene.

## 4. The lint script

`scripts/narration_lint.py` is stdlib only and read-only. It exits 1 on any error.

```
python3 scripts/narration_lint.py script.json
python3 scripts/narration_lint.py script.json --audio-dir public/narration
python3 scripts/narration_lint.py script.json --strict
```

| Check | Level |
|---|---|
| Sentence longer than `limits.sentence` | error |
| Contraction (`n't`, `'re`, `'ve`, `'ll`, `'d`, `'m`, `it's`) | error |
| Symbol that cannot be spoken (arrows, comparison signs, slashes, backticks, underscores) | error |
| A glossary `avoid` word | error |
| Scene words over the budget | error |
| Audio file longer than the scene minus the tail (`--audio-dir`, needs `ffprobe`) | error |
| Paragraph longer than `limits.paragraph` | warning |
| Possible passive voice | warning |
| No audio file for a scene | warning |

The passive check is a pattern, so it finds `is validated` and misses other forms. Warnings are heuristics. There is no approved-word dictionary, so a clean run does not make the text ASD-STE100 compliant.

## 5. Wiring audio in Remotion

One audio file per scene, inside the scene's own sequence, so it moves with the scene. See [remotion-architect](../../frameworks/remotion-architect/SKILL.md) for `Sequence`, `staticFile` and `<Audio>`.

```tsx
import {Audio} from '@remotion/media';
import {Sequence, staticFile} from 'remotion';

// `from` and `frames` come from the timeline table; LEAD_IN_FRAMES is Math.round(leadIn * fps).
<Sequence from={fromFrame} durationInFrames={frames}>
  <Sequence from={LEAD_IN_FRAMES}>
    <Audio src={staticFile(`narration/${scene.id}.mp3`)} />
  </Sequence>
</Sequence>
```

- Keep the generated audio out of the composition's source: write it to `public/narration/`, and git-ignore it if it can be rebuilt.
- Name each file by scene id, so the lint script finds it.
- Narration is generated from the script before the render. The render never calls a speech engine.
- If the audio is longer than its scene, cut words and generate again. Do not lengthen the scene to fit a script that was too long, unless the picture needs the time too.

## 6. Manim and other silent clips

- Render the clip with no sound and let the frame that embeds it carry the audio. The clip's own timing (`run_time`, `wait`) is the cue sheet: add up the seconds per part and put those in the script's `cue` field.
- For a clip that is published alone, join the audio with `ffmpeg -i clip.mp4 -i narration.mp3 -c:v copy -c:a aac -shortest out.mp4`.
- Check that the narration length does not pass the clip length, with the same `--audio-dir` gate.

## 7. Before and after

Dense:
"Once the build artifacts have been successfully generated, they're subsequently pushed to the registry by the pipeline, where they can then be consumed."

Controlled (80%):
"The pipeline builds the artifacts. Then it sends them to the registry. Other jobs can use them from there."

Rule shown: one idea per sentence, active voice with a named actor, no contraction, plain verbs.

Condition first:
"The fan must start if the temperature goes over 80." becomes "If the temperature is over 80, the fan must start."

Say what was not run:
"The result should hold for other programs." becomes "This was run on one program. It was not run on others."

## 8. Interview before writing

Ask in order, and wait for each answer:

1. Video length, and the scene table with the seconds of each scene.
2. Audience and their language level.
3. Language of the narration, and whether a translation is planned.
4. Where the audio is made: a hosted engine or a local one. This is the owner's decision (cost, privacy, rights to the voice).
5. What is real and what is an illustration. This sets the lines that say so.
6. The sources behind each scene, so no sentence says more than they do.
7. Terms that need a spoken form.

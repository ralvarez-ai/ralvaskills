# Explainer recipes

## 1. The content file

One JSON file per video. The scene code reads facts from it; narration stays in its own file ([narration-writer](../narration-writer/RECIPES.md)).

```json
{
  "id": "topic-v1",
  "number": 1,
  "sources": {
    "spec": {"label": "Vendor spec, section 4", "kind": "source"},
    "book": {"label": "Author, Book (2008), chapter 2", "kind": "textbook"},
    "run": {"label": "Recorded run of 2026-11-02", "kind": "run"}
  },
  "scenes": [
    {"id": "intro", "chip": "SOURCE", "sources": ["spec"]},
    {"id": "example", "chip": "ILLUSTRATION", "sources": []}
  ],
  "facts": {
    "limit": {"value": 5, "source": "spec"}
  },
  "unchecked": ["Only one input was run.", "The order of steps was read from the spec, not measured."]
}
```

- A lint step can fail the build when a scene's chip is missing, when a chip is `SOURCE` or `REAL RUN` and the scene lists no source, or when a source key does not exist.
- The last scene is generated from `unchecked`, so the list cannot be forgotten.

## 2. The contact sheet

Needs `ffmpeg`. `sections.json` is a list of `{id, start, end, seconds}` (the render writes it; see the engine skill).

```python
import json, subprocess, os, sys
video, sections, out = sys.argv[1:4]          # video.mp4 sections.json outdir
os.makedirs(out, exist_ok=True)
for i, s in enumerate(json.load(open(sections))):
    t = s["start"] + 0.85 * s["seconds"]
    subprocess.run(["ffmpeg", "-v", "error", "-y", "-ss", f"{t:.2f}", "-i", video,
                    "-frames:v", "1", f"{out}/{i:02d}-{s['id']}.png"], check=True)
```

Then tile: `ffmpeg -v error -y -pattern_type glob -i 'OUT/*.png' -vf 'scale=640:-1,tile=3xROWS' sheet.png`. Read the sheet first, then open single frames of the dense scenes.

## 3. The cue and the length

```python
def spoken_words(text, say):            # say: {"CTL": "C T L", ...}
    for term, form in say.items():
        text = text.replace(term, form)
    return len(text.split())

def cue_time(text, phrase, rate, lead):  # raises ValueError if phrase is gone
    at = text.index(phrase)
    return lead + spoken_words(text[:at], SAY) / rate

def scene_seconds(animation, text, rate, lead, tail):
    return max(animation, spoken_words(text, SAY) / rate + lead + tail)
```

Take `rate` from the narration script (measured), not from a constant in the scene code.

With audio, skip the estimate: `scene_seconds = audio_seconds + lead + tail + hold` (hold about 1.2 s), and `cue_time = max(0.3, lead + (index_of_phrase / len(text)) * audio_seconds - 0.35)`. Character position is a fair stand-in for time when the engine gives no word timestamps; use the timestamps when it does. Write `durations` in the manifest when the audio is made (`ffprobe -show_entries format=duration`), so the renderer needs no audio library.

## 4. A scene plan table

| Scene | Idea (one line) | Drawing | Formula or code beside it | Chip | Sources | Cue phrases |
|---|---|---|---|---|---|---|
| intro | what we compute | the object, empty | none yet | SOURCE | spec | "starts with" |

Fill the table before any animation code. A row with no drawing is a scene to cut or split.

## 5. A video of recorded runs

For a video whose scenes show the results of a batch of runs (a benchmark, a search, a training sweep). The steps come before any scene plan.

1. **Pick the run by a rule, not by a look.** Write the rule down (best score, widest effect, cleanest story) and name the runner-up. If the owner overrides the pick, record the swap and why in the content file.
2. **Rebuild the numbers from the primary record** (database, logs, result files), not from a review. Compare with each review that covers the same runs. If they disagree, keep the later one, say which, and do not copy a withdrawn claim. Record both totals when they differ (for example cost from the database against the review).
3. **Mark each fact as recorded or inferred.** An inference (such as "probably a timeout" from a gap in a log) keeps its word "probably" on screen and in the content file, with the evidence beside it.
4. **Compare two batches on the items both completed.** Show the all-items total and the matched-subset total side by side. Name the items that explain any gap, for instance a run that failed with no grade. One run per arm means a few items are inside the noise of a rerun: say so in the comparison scene.
5. **Put a caveat where it changes the reading.** A comparison scene carries its own caveats (cost stop rules, what the spend leaves out). If the repository's rule is no closing caveat scene, keep the rest in the content file's `unchecked` list.
6. **Name the version of the thing measured.** A run is a release only when its tag says so; otherwise the footer says "development" and the chip stays `REAL RUN`.
7. **Draft, then final.** Render a low-quality draft, read the contact sheet, then render once at full size and add the voice. Report what was not verified: audio not listened to, a source not re-read, a lint script not run.

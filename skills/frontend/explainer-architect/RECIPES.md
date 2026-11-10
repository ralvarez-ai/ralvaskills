# Explainer recipes

## 1. The content file

One JSON file per video. The scene code reads facts from it; narration stays in its own file ([narration-architect](../narration-architect/RECIPES.md)).

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

## 4. A scene plan table

| Scene | Idea (one line) | Drawing | Formula or code beside it | Chip | Sources | Cue phrases |
|---|---|---|---|---|---|---|
| intro | what we compute | the object, empty | none yet | SOURCE | spec | "starts with" |

Fill the table before any animation code. A row with no drawing is a scene to cut or split.

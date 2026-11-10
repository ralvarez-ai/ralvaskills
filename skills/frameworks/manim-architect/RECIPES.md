# Manim recipes

## 1. The base class

Reads the narration and content files, closes each scene to its time, writes the sections file.

```python
import json
from pathlib import Path
from manim import Scene, FadeOut

OUT = Path("out")

class Chapter(Scene):
    CHAPTER = ""          # id; names the narration, content and sections files
    RATE = 2.5            # fallback; the narration file's "rate" wins

    def construct(self):
        narration = json.loads(Path(f"narration/{self.CHAPTER}.json").read_text())
        self.words = {s["id"]: s["narration"] for s in narration["scenes"]}
        self.rate = narration.get("rate", self.RATE)
        self.lead = narration.get("leadIn", 0.5)
        self.tail = narration.get("tail", 0.5)
        self.timeline, self._current = [], None
        self.script()
        self.close_scene()
        OUT.mkdir(exist_ok=True)
        (OUT / f"{self.CHAPTER}.sections.json").write_text(json.dumps(self.timeline, indent=2))

    def open_scene(self, scene_id):
        self.close_scene()
        self._current = {"id": scene_id, "start": self.renderer.time}

    def close_scene(self):
        if self._current is None:
            return
        s = self._current
        needed = len(self.words[s["id"]].split()) / self.rate + self.lead + self.tail
        elapsed = self.renderer.time - s["start"]
        if elapsed < needed:
            self.wait(needed - elapsed)
        self.play(*[FadeOut(m) for m in self.mobjects], run_time=0.4)
        end = self.renderer.time
        self.timeline.append({"id": s["id"], "start": round(s["start"], 3),
                              "end": round(end, 3), "seconds": round(end - s["start"], 3)})
        self._current = None
```

Count spoken words (after pronunciation forms) when the narration has acronyms said as letters.

## 2. The cue helper

```python
def cue(scene, phrase, delay=0.0):
    text = scene.words[scene._current["id"]]
    at = text.index(phrase)                       # ValueError: the narration changed
    target = scene.lead + len(text[:at].split()) / scene.rate + delay
    behind = target - (scene.renderer.time - scene._current["start"])
    if behind > 0.03:
        scene.wait(behind)
```

Use it before the animation that should land on the word: `cue(self, "and the end mark")`.

## 3. Muxing the audio

```
ffmpeg -y -i video.mp4 -i s1.mp3 -i s2.mp3 \
  -filter_complex "[1:a]adelay=D1|D1[a1];[2:a]adelay=D2|D2[a2];[a1][a2]amix=inputs=2:normalize=0[aout]" \
  -map 0:v -map "[aout]" -c:v copy -c:a aac out.voiced.mp4
```

`Dn` is `(scene start + lead-in) * 1000` in milliseconds, from the sections file. Build the filter in a small script, one input per scene that has audio. Write a second output per voice (`.local.voiced.mp4` for a library voice) so one voice never overwrites another.

## 4. LaTeX helpers

```python
from manim import MathTex, TexTemplate, Text, VGroup, RIGHT, DOWN
TEX = TexTemplate(preamble=r"\usepackage{amsmath}\usepackage{amssymb}")

def tex(latex, size=32, color="#ffffff", max_width=12.6):
    f = MathTex(latex, tex_template=TEX, color=color, font_size=size * 1.2)
    if f.width > max_width:
        f.scale_to_fit_width(max_width)
    return f

def mix(size, color, *parts):                     # strings are text, ("m", latex) parts are formulas
    pieces = [tex(p[1], size, color, 9.0) if isinstance(p, tuple) else Text(p, font_size=size * 0.9, color=color)
              for p in parts]
    return VGroup(*pieces).arrange(RIGHT, buff=0.22, aligned_edge=DOWN)
```

## 5. Render commands

```
uv run manim --resolution 1920,1080 -qh --fps 30 --media_dir media chapters/foo.py Foo > out/foo.render.log 2>&1
uv run manim -ql --media_dir media-draft chapters/foo.py Foo > out/foo.draft.log 2>&1
```

Wrap them in one script that takes a registry id, passes `--quality` through, copies the result to `out/<id>.mp4` and prints only the path and the sections file.

## 6. Embedding a clip in Remotion

Render the clip silent, copy it to `public/clips/`, and play it in a scene with `<Video src={staticFile('clips/foo.mp4')} />` from `@remotion/media`; read its length from the sections file or a recorded metadata file and use it for the scene's `durationInFrames`. The Remotion frame owns the audio and the badges.

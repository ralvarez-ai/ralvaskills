# Manim framework: what to reach for

Manim Community (0.22 at the time of writing; check names against your installed version).
Tools of the library that matter for narrated explainers, and the ones that replace code you would write.

## Pick the scene class

| Need | Class |
|---|---|
| Normal 2D | `Scene` |
| Zoom or pan over a drawing | `MovingCameraScene` (`self.camera.frame.animate.scale(0.5).move_to(m)`) |
| Magnified inset | `ZoomedScene` |
| 3D objects, rotating camera | `ThreeDScene` (`set_camera_orientation`, `begin_ambient_camera_rotation`) |

The frame is 8 units high and ~14.2 wide (`config.frame_height`, `config.frame_width`); lay out in those units.

## Text and formulas

- `Text` (Pango, any font), `MarkupText` (inline colour and weight), `Tex` and `MathTex` (LaTeX).
- Isolate parts of a formula to colour or animate them: `MathTex("a", "^2", "+", "b", "^2", substrings_to_isolate=[...])`,
  then `f[0].set_color(...)`, `f.set_color_by_tex("a", BLUE)`. Index by piece, not by glyph.
- `TransformMatchingTex` / `TransformMatchingShapes` morph equal pieces and move them; use them for equation
  steps. For labels that differ, fade instead.
- `Variable`, `DecimalNumber` and `Integer` hold a number that an updater can change.
- Fonts: `with register_font("path.ttf"):` around the scene (user fonts), then `Text(..., font="Name")`.

## Building blocks

- Graphs: `Axes`, `NumberPlane`, `axes.plot(f)`, `axes.c2p(x, y)` (data to screen), `axes.get_area(...)`,
  `axes.get_riemann_rectangles(...)`, `BarChart`.
- Data: `Table`, `Matrix`, `DecimalNumber`, `Dot`, `Line`, `Arrow`, `Polygon`, `ImageMobject`.
- Layout: `VGroup.arrange(...)`, `next_to`, `align_to`, `to_edge`, `to_corner`, `get_center`, `scale_to_fit_width`.
- Annotation: `Brace` + `get_text`, `SurroundingRectangle`, `DoubleArrow`, `Underline`.

## Animation

- Appear: `Create`, `Write`, `FadeIn(shift=UP * 0.3)`, `GrowFromCenter`, `DrawBorderThenFill`.
- Change: `obj.animate.shift(...).set_color(...)` (one chained animation), `Transform` (changes the source
  object into the target), `ReplacementTransform` (the target takes over; the source is gone).
- Sequence: `LaggedStart(*anims, lag_ratio=0.15)`, `LaggedStartMap(FadeIn, group)`, `Succession`, `AnimationGroup`.
- Emphasis, one at a time: `Indicate`, `Flash`, `Circumscribe`, `FocusOn`, `ApplyWave`, `Wiggle`.
- `run_time` sets duration; `rate_func` the easing (`smooth` default, `linear` for steady motion, `there_and_back`).
- Values that move with the narration: a `ValueTracker` drives `always_redraw(lambda: ...)` objects and
  `add_updater`; animate `tracker.animate.set_value(x)` on the cue.
- Cue and hold times are `self.wait(seconds)`. Seconds do not depend on fps, so draft and final times match.

## Sections and sound built in

- `self.next_section("id")` marks a section; `manim --save_sections` writes a JSON with each section's
  name, duration and a split video file. It can replace a hand-written sections file when sections equal narration scenes.
- `self.add_sound("audio/hook.wav", time_offset=0.3)` puts audio in the render at the current scene time.
  Use it for a quick voiced draft. For a final, keep the separate mux (RECIPES.md, section 3): the video is not
  re-rendered when audio changes, and a paid voice and an offline voice stay as two outputs.

## CLI

| Flag | Use |
|---|---|
| `-ql` / `-qm` / `-qh` | 480p15 / 720p30 / 1080p60 (override with `-r 1920,1080 --fps 30`) |
| `-s` | Save the last frame as a PNG, no video: quick look at one scene |
| `-n 4,9` | Render only animations 4 to 9 of a long scene |
| `-a` | Every scene class in the file |
| `--media_dir DIR` | Cache and output location; separate per parallel render |
| `--disable_caching` | When a helper changed in a way the cache cannot see |
| `--save_sections` | Section JSON and split files |
| `--format png` / `--format gif` | Frames or a GIF |

## Frames for review without ffmpeg

At the end of a scene, `self.camera.get_image().save(f"frames/{scene_id}.png")` writes the exact end
state (Cairo renderer, works at `-ql`). `scripts/contact_sheet.py` (this skill) does the same from a finished video;
use whichever is closer to your render step. Tile either set into one sheet per version.

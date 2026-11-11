# manim-architect digest

Manim Community 0.22 scenes for formal objects that change shape, timed from the narration, with LaTeX formulas, cues and a separate audio mux.

## Always-on rules

- One chapter is one `Scene` subclass on a base class that reads the narration and content files, closes each scene to `max(animation, words / rate + lead-in + tail)` and writes a sections file (`id`, `start`, `end`, `seconds`).
- Cue drawings to narration phrases with a `cue(scene, phrase)` helper that raises when the phrase is gone.
- Morph shapes, not text: fade the new label in after the shape moves.
- Formulas are LaTeX through one helper with a minimal template, scaled to a maximum width. A failing formula fails the whole render: render after every edit that touches one.
- Colours and fonts come from the shared token file. Teaching helpers live in a package, one module per family; chapter files stay thin.
- Draft with `-ql`, review at 1920x1080 30 fps `-qh`. Always redirect output to a log file.
- Audio is muxed in a separate step from the sections file (`adelay`, `amix=normalize=0`, copy the video); scene code never calls a speech engine.
- Keep a drawing across scenes by family, set `z_index` on groups, use `Group` for images; teach a formula with one real worked case first (details in TRAPS.md).
- Review with one labelled contact sheet per version (`scripts/contact_sheet.py video.mp4 sections.json`) and full-size frames of dense scenes (explainer-architect). Library tools (scene classes, `next_section`, `add_sound`, CLI flags) are in FRAMEWORK.md.

## Hard constraints

- Never animate a `Transform` between two text labels.
- Never hard-code scene times; read the narration file and its measured rate.
- Never print render output into the session.
- Never claim a tool is missing without `command -v`.

Invoke the full skill when: writing, structuring, rendering or reviewing Manim scenes, alone or beside a Remotion video.

Full guide: invoke the manim-architect skill.

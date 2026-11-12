# Explainer videos in Remotion

The Remotion-specific half of a teaching video. What to teach, the scene rule, provenance chips, content files, cues and versions are engine-independent: read [explainer-designer](../../frontend/explainer-designer/SKILL.md) first. Words and speech rate: [narration-writer](../../frontend/narration-writer/SKILL.md). Scenes where a formal object changes shape can be drawn with [manim-architect](../manim-architect/SKILL.md) and embedded here.

## 1. Scene kinds as data

- Model scenes as a discriminated union (`kind`) in the props schema: `cards`, `timeline`, `stats`, `bars`, `unroll`, `stamp`, and so on. A new way to show something is a new kind with its own schema, not a special case inside an old one.
- Give a kind an optional `variant` (`a` | `b`) when two looks are being compared. A render prop (`--props='{"variant":"a"}'`) forces one for the whole video; without it each scene picks its own.
- The default render must not change when a kind or a variant is added: a new look is a prop (`style: "v1" | "v2"`) whose default is the old look. Existing compositions and Studio previews stay as they were.
- Read the facts of a scene from the content file, and refuse to build when a scene has no origin chip or a number the file does not hold (a `calculateMetadata` check or a prebuild script).

## 2. Timing from the audio

- The audio manifest of each take (`public/<take>/<video>/manifest.json`) lists the scenes that have audio and their `durations`. `calculateMetadata` reads it, so a take with another voice gets its own scene lengths and its own total length. A scene with no audio keeps the `seconds` of its content file.
- Items carry a `cue` phrase; a pure function resolves it to `atSec` from the narration text and the audio length before the scene renders. Check every cue at load (a registry that imports the narration files), not in a render.
- A scene that shows a clip can have `parts` (id, from, to): each part has its own narration and audio, placed at the part's start. The clip's script writes the start of each part to a JSON file next to its metadata. The scene-level audio is skipped when parts exist.

## 3. Charts and styling

- Charts are components that take data and a theme and draw SVG from `useCurrentFrame()`. Wire one into a scene through an optional field on a kind (for example `chart` on a `cards` scene), validated by the schema.
- Colours in SVG and inline styles are theme CSS variables (`var(--token)`); a JS object of colours that the stylesheet cannot see drifts.
- Scene timing: a scene lasts `max(animation, words / rate + lead-in + tail)` in frames, `Math.round(seconds * fps)`, with the measured rate from the narration script.

## 4. Embedding a Manim clip

Render the clip silent, copy it to `public/clips/`, play it with `<Video>` from `@remotion/media`, and size the scene from the clip's recorded length. The Remotion frame owns the badges, the sources and the audio (`<Audio>` per scene, narration-writer `RECIPES.md`).

## 5. Pitfalls

- Remotion bundles `src/` at the start of every render. Do not edit `src/` while a render loop runs: the next render in the loop picks up the half-edit.
- A render that fails at once with a webpack `TypeError` on `length` usually has a corrupt cache: delete `node_modules/.cache/webpack` and render again.
- Several tasks may write the same default output file (for example `out/video.mp4`). Render with an explicit output path per version: `remotion render src/index.ts <id> out/videos/<name>.mp4 --props=...`.
- Never kill a render with `pkill -f <pattern>` when the pattern also appears in your own command line; kill by PID or use a bracketed pattern (`[r]emotion`).
- Parallel renders of one project need a worktree and an install each; take stills one at a time per worktree.

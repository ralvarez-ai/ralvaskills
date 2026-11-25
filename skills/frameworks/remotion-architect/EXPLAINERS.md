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

## 6. Flow scenes

A request, a check that can fail or a message between two parties is drawn as something that travels ([VISUAL.md](../../frontend/explainer-designer/VISUAL.md)). In Remotion that is SVG driven by `interpolate`, with no extra package. Make each of these a scene kind of its own (`sequence`, `checklist`, `columns`) with its schema, not a special case of `cards`. Typechecked against Remotion 4.

```tsx
import {interpolate, useCurrentFrame} from 'remotion';

type Point = {x: number; y: number};
const clamp = {extrapolateLeft: 'clamp', extrapolateRight: 'clamp'} as const;

// An arrow that draws itself from `start` over `length` frames, with its label above it.
export const Message = ({from, to, start, length = 15, color, label}: {from: Point; to: Point; start: number; length?: number; color: string; label: string}) => {
  const frame = useCurrentFrame();
  const progress = interpolate(frame, [start, start + length], [0, 1], clamp);
  const end = {x: from.x + (to.x - from.x) * progress, y: from.y + (to.y - from.y) * progress};
  return (
    <g opacity={progress > 0 ? 1 : 0}>
      <line x1={from.x} y1={from.y} x2={end.x} y2={end.y} stroke={color} strokeWidth={4} />
      <text x={(from.x + to.x) / 2} y={from.y - 12} fill={color} fontSize={26} textAnchor="middle">{label}</text>
    </g>
  );
};

// A packet that travels and stays at the end. Two Packets in a row make a pass (down a list) or a bounce (back to the start).
export const Packet = ({from, to, start, length = 20, color}: {from: Point; to: Point; start: number; length?: number; color: string}) => {
  const frame = useCurrentFrame();
  const x = interpolate(frame, [start, start + length], [from.x, to.x], clamp);
  const y = interpolate(frame, [start, start + length], [from.y, to.y], clamp);
  return <circle cx={x} cy={y} r={12} fill={color} />;
};
```

- **A sequence diagram:** one vertical line per actor, one `Message` per sentence with `start` at the cue's frame, 0.45 to 0.55 of the height of a row between messages.
- **A check list with a bounce:** a dot per check turns from the neutral colour to the pass colour as the packet reaches its row; the failing check turns the fail colour and a second `Packet` takes it back to the start, with the reason beside it. Reset the dots between two attempts.
- **Columns by maturity:** a header, a rule and a stack of chips per column; delay each chip by subtracting frames (`frame - cueFrame - index * 8`) as the narration lists it.
- **Legibility floor:** pick a minimum font size for the final frame (for example 24 px at 1920x1080, a smaller one only for a short status tag) and check the theme, not each scene: a label under it cannot be read on a phone.
- **Maturity badge:** a status badge with a tone per state (`built`, `open`, `proposed`, `illustrative`) carries what the chip says about provenance and what a part's maturity is; one badge per part, not per video.
- **Review:** `remotion still <id> --frame=<n>` at about 97% of each scene's `durationInFrames`, tiled; open the last second of a scene that looks wrong before changing it.
- **Cues:** check them at load, as in section 2, and also that the cues of a scene follow the order of its narration: a cue earlier than the one before it is drawn late.


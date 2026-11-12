# remotion-architect digest

Builds Remotion 4.0.542 videos as real React source: compositions, frame-deterministic animation, timing, props schemas, and CLI / SSR / Lambda rendering.

## Always-on rules

- Interview first, in order, waiting for each answer: output (video/audio/still, size, fps, length); fixed or data-driven; assets; render target (CLI, Node, Lambda/Cloud Run, Player only); licensing; styling.
- A video is a pure function of `useCurrentFrame()`. Duration is an integer frame count: `Math.round(seconds * fps)`.
- Animate with `interpolate` (clamp both sides) and `spring`; delay by subtracting frames or wrapping in `<Sequence from>`.
- Chain scenes with `<Series>`; use `<TransitionSeries>` for transitions and subtract the overlap from total duration.
- Media: `<Img>`, `<Video>`/`<Audio>` from `@remotion/media`, `staticFile()` from `public/`. Fonts via `@remotion/google-fonts` or `@remotion/fonts`.
- Async work (fonts, fetches) uses `useDelayRender()` and always releases; data-driven length and props go in `calculateMetadata`.
- Props are typed with sensible `defaultProps`; Zod schema (`npx remotion add zod`) only for validation or specialized controls.
- All `remotion` and `@remotion/*` packages share one exact version, no `^`. Upgrade with `npx remotion upgrade`.
- Batch jobs bundle once and call `renderMedia` per variant; validate `inputProps` before server renders.
- Derive layout from `useVideoConfig()`, never hardcoded dimensions.
- Teaching videos: plan with explainer-designer, words from narration-writer; scenes as a union of kinds with optional `a`/`b` variants, a `style` prop whose default is the old look, charts wired through the schema, Manim clips embedded with `<Video>`. See EXPLAINERS.md.

## Hard constraints

- Never use CSS transitions/animations, Tailwind `animate-*`/`transition-*`, `Math.random()`, `Date.now()`, timers or `requestAnimationFrame` to drive visuals; use `random('seed')` and frame math.
- Never use native `<img>`, `<video>`, `<audio>` or absolute filesystem paths.
- Never leak a `delayRender` handle.
- Never overwrite a published render (old is v1, new is v2), and never change the default render when adding a look.
- Never assume the license is fine: raise the Company License question at https://www.remotion.dev/docs/license.

Invoke the full skill when: scaffolding a Remotion project, writing compositions or scenes, choosing a render target, or reviewing Remotion source.

Full guide: invoke the remotion-architect skill.

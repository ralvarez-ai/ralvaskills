# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| remotion | 4.0.542 | Core: compositions, `interpolate`, `spring`, `Sequence`, `Img` |
| @remotion/cli | 4.0.542 | `remotion studio` / `render` / `still` |
| @remotion/media | 4.0.542 | `<Video>` and `<Audio>` (recommended media components) |
| @remotion/transitions | 4.0.542 | `TransitionSeries` and transition presets |
| @remotion/renderer | 4.0.542 | `renderMedia`, `selectComposition` (server-side rendering) |
| @remotion/bundler | 4.0.542 | `bundle()` for SSR |
| @remotion/lambda | 4.0.542 | AWS Lambda rendering (opt-in) |
| @remotion/player | 4.0.542 | Embed a composition in a React app (opt-in) |
| react / react-dom | per `create-video` template | Declared by the Remotion template; follow it |

## Notes

- **Every `remotion` and `@remotion/*` package must be the exact same version, no `^`/`~`.** Upgrade with `npx remotion upgrade`.
- Add Remotion packages with `npx remotion add <pkg>` so the version matches automatically.
- Zod and `@remotion/zod-types` are opt-in for props schemas; install via `npx remotion add zod`.
- React-side conventions defer to [react-architect](../react-architect/SKILL.md).

_Last reviewed: 2026-11-09_
_Skill version at last review: 1.1.0_

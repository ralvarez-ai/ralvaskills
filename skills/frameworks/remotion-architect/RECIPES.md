# Remotion recipes

Full reference implementations for [SKILL.md](SKILL.md). Versions pinned in [STACK.md](STACK.md); every `remotion`/`@remotion/*` dependency uses the same exact version.

## Project skeleton

```
my-video/
├── package.json            # exact versions, no ^
├── remotion.config.ts
├── tsconfig.json
├── public/                 # staticFile() assets: fonts, music, footage
├── src/
│   ├── index.ts            # registerRoot
│   ├── Root.tsx            # <Composition>s, grouped in <Folder>
│   ├── lib/
│   │   ├── timing.ts       # seconds → frames, scene durations
│   │   └── animation.ts    # named helpers over interpolate/spring
│   └── promo/              # one folder per video
│       ├── Promo.tsx       # sequences the scenes
│       ├── schema.ts       # props type (+ Zod when needed)
│       └── scenes/
│           ├── Intro.tsx
│           └── Outro.tsx
├── scripts/render-batch.ts # optional SSR / batch
└── out/                    # render output (gitignored)
```

`src/index.ts`:

```ts
import { registerRoot } from "remotion";
import { RemotionRoot } from "./Root";

registerRoot(RemotionRoot);
```

`src/Root.tsx`:

```tsx
import { Composition, Folder } from "remotion";
import { Promo } from "./promo/Promo";
import { promoDefaults, promoDuration } from "./promo/schema";

export const RemotionRoot = () => (
  <Folder name="promo">
    <Composition
      id="promo"
      component={Promo}
      fps={30}
      width={1920}
      height={1080}
      durationInFrames={promoDuration(promoDefaults)}
      defaultProps={promoDefaults}
      calculateMetadata={({ props }) => ({
        durationInFrames: promoDuration(props),
      })}
    />
  </Folder>
);
```

`remotion.config.ts`:

```ts
import { Config } from "@remotion/cli/config";

Config.setVideoImageFormat("jpeg");
Config.setOverwriteOutput(true);
```

`package.json` scripts: `"studio": "remotion studio"`, `"render": "remotion render promo out/promo.mp4"`, `"still": "remotion still promo out/thumb.png --frame=0"`, `"upgrade": "remotion upgrade"`. Mirror them as Task targets (see [repo-tooling-architect](../../tooling/repo-tooling-architect/SKILL.md)).

## Scenes, timing and animation helpers

`src/lib/timing.ts`:

```ts
export const FPS = 30;
export const sec = (s: number, fps = FPS) => Math.round(s * fps);
```

`src/lib/animation.ts`:

```ts
import { interpolate, spring, useCurrentFrame, useVideoConfig } from "remotion";

export const useFadeIn = (start = 0, duration = 15) => {
  const frame = useCurrentFrame();
  return interpolate(frame, [start, start + duration], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });
};

export const usePop = (delay = 0) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  return spring({ frame: frame - delay, fps, config: { damping: 200 } });
};
```

A scene (`scenes/Intro.tsx`):

```tsx
import { AbsoluteFill } from "remotion";
import { useFadeIn, usePop } from "../../lib/animation";

export const Intro = ({ title }: { title: string }) => {
  const opacity = useFadeIn(0, 12);
  const scale = usePop(4);
  return (
    <AbsoluteFill
      style={{ justifyContent: "center", alignItems: "center", opacity }}
    >
      <h1 style={{ fontSize: 120, transform: `scale(${scale})` }}>{title}</h1>
    </AbsoluteFill>
  );
};
```

Sequencing with `Series` (`promo/Promo.tsx`):

```tsx
import { AbsoluteFill, Series } from "remotion";
import { Intro } from "./scenes/Intro";
import { Outro } from "./scenes/Outro";
import { sceneFrames, type PromoProps } from "./schema";

export const Promo = (props: PromoProps) => (
  <AbsoluteFill style={{ backgroundColor: props.background }}>
    <Series>
      <Series.Sequence durationInFrames={sceneFrames.intro}>
        <Intro title={props.title} />
      </Series.Sequence>
      <Series.Sequence durationInFrames={sceneFrames.outro}>
        <Outro cta={props.cta} />
      </Series.Sequence>
    </Series>
  </AbsoluteFill>
);
```

`promo/schema.ts`:

```ts
import { sec } from "../lib/timing";

export type PromoProps = { title: string; cta: string; background: string };

export const promoDefaults: PromoProps = {
  title: "Hello",
  cta: "Get started",
  background: "#0b0b10",
};

export const sceneFrames = { intro: sec(3), outro: sec(2) };

export const promoDuration = (_: PromoProps) =>
  sceneFrames.intro + sceneFrames.outro;
```

## Transitions

```tsx
import { linearTiming, TransitionSeries } from "@remotion/transitions";
import { fade } from "@remotion/transitions/fade";

const T = 15; // transition length in frames

<TransitionSeries>
  <TransitionSeries.Sequence durationInFrames={90}>
    <Intro title="Hi" />
  </TransitionSeries.Sequence>
  <TransitionSeries.Transition
    presentation={fade()}
    timing={linearTiming({ durationInFrames: T })}
  />
  <TransitionSeries.Sequence durationInFrames={60}>
    <Outro cta="Go" />
  </TransitionSeries.Sequence>
</TransitionSeries>;
// Total = 90 + 60 − T = 135. Put this arithmetic in schema.ts, not in your head.
```

## Fonts and media

```tsx
import { loadFont } from "@remotion/google-fonts/Inter";
import { Img, staticFile, useVideoConfig } from "remotion";
import { Audio, Video } from "@remotion/media";

const { fontFamily } = loadFont(); // module scope: blocks the render until loaded

export const Clip = () => {
  const { durationInFrames } = useVideoConfig();
  return (
    <>
      <Video src={staticFile("footage.mp4")} trimBefore={30} muted />
      <Audio
        src={staticFile("music.mp3")}
        volume={(f) => Math.min(1, f / 30, (durationInFrames - f) / 30)}
      />
      <Img src={staticFile("logo.png")} style={{ fontFamily }} />
    </>
  );
};
```

Local font file: `import { loadFont } from "@remotion/fonts"; loadFont({ family: "Brand", url: staticFile("Brand.woff2") });`.

## Data fetching and `useDelayRender`

Prefer `calculateMetadata` for anything that decides duration or props:

```tsx
<Composition
  id="report"
  component={Report}
  fps={30}
  width={1080}
  height={1920}
  durationInFrames={1}
  defaultProps={{ url: "https://example.com/data.json", rows: [] as Row[] }}
  calculateMetadata={async ({ props, abortSignal }) => {
    const rows: Row[] = await (await fetch(props.url, { signal: abortSignal })).json();
    return {
      props: { ...props, rows },
      durationInFrames: Math.max(1, rows.length * 45),
    };
  }}
/>
```

When a component itself must wait (a canvas, a third-party loader):

```tsx
import { useDelayRender } from "remotion";
import { useEffect, useState } from "react";

const { delayRender, continueRender, cancelRender } = useDelayRender();
const [handle] = useState(() => delayRender("load model"));
useEffect(() => {
  load().then(() => continueRender(handle)).catch(cancelRender);
}, [handle]);
```

## Zod schema (only when needed)

```ts
import { z } from "zod";
import { zColor } from "@remotion/zod-types";

export const promoSchema = z.object({
  title: z.string().min(1).max(60),
  cta: z.string(),
  background: zColor(),
});
export type PromoProps = z.infer<typeof promoSchema>;
// <Composition schema={promoSchema} defaultProps={promoDefaults} ... />
```

Validate server-side before rendering: `promoSchema.parse(untrustedInput)`.

## Server-side rendering

```ts
import path from "node:path";
import { bundle } from "@remotion/bundler";
import { renderMedia, selectComposition } from "@remotion/renderer";

const serveUrl = await bundle({ entryPoint: path.resolve("src/index.ts") });

export async function renderPromo(inputProps: unknown, out: string) {
  const props = promoSchema.parse(inputProps);
  const composition = await selectComposition({
    serveUrl,
    id: "promo",
    inputProps: props,
  });
  await renderMedia({
    composition,
    serveUrl,
    codec: "h264",
    outputLocation: out,
    inputProps: props,
  });
}
```

Bundle once per process (or per deploy), not per request. Put renders behind a bounded queue; each one uses a headless Chrome.

## Batch rendering from data

```ts
import { readFile } from "node:fs/promises";
// ...bundle() once as above
const rows: PromoProps[] = JSON.parse(await readFile("data/variants.json", "utf8"));
for (const [i, row] of rows.entries()) {
  await renderPromo(row, `out/promo-${String(i).padStart(3, "0")}.mp4`);
}
```

Run a few in parallel only after measuring memory; each `renderMedia` already spreads frames over `--concurrency` tabs.

## AWS Lambda

One-time setup (CLI, from the project root):

```sh
npx remotion lambda policies role
npx remotion lambda policies user          # then create the IAM user and set env vars
npx remotion lambda functions deploy
npx remotion lambda sites create src/index.ts --site-name=my-video
npx remotion lambda render <serve-url> promo --props='{"title":"Hi"}'
```

From code (`@remotion/lambda/client`):

```ts
import { getRenderProgress, renderMediaOnLambda } from "@remotion/lambda/client";

const { renderId, bucketName } = await renderMediaOnLambda({
  region: "us-east-1",
  functionName,
  serveUrl,
  composition: "promo",
  codec: "h264",
  inputProps: props,
});
const progress = await getRenderProgress({ renderId, bucketName, functionName, region: "us-east-1" });
```

Re-run `sites create` whenever composition code changes; the site is a deployed bundle. Keep the Lambda function version equal to the Remotion version (redeploy after `remotion upgrade`). Credentials come from the environment or a secret store, never the repo (see [secrets-architect](../../infra/secrets-architect/SKILL.md)).

## Player (in-browser preview)

```tsx
import { Player } from "@remotion/player";

<Player
  component={Promo}
  inputProps={promoDefaults}
  durationInFrames={promoDuration(promoDefaults)}
  compositionWidth={1920}
  compositionHeight={1080}
  fps={30}
  controls
  style={{ width: "100%" }}
/>;
```

The Player plays; it does not export. Export via CLI, SSR or Lambda.

## Docker (renderer)

Use Remotion's documented Dockerfile for the current version (Node base plus Chrome's system libraries, then `npx remotion browser ensure`). Keep it multi-stage and non-root per [docker-architect](../../infra/docker-architect/SKILL.md); skip distroless, which lacks Chrome's dependencies.

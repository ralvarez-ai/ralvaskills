---
title: 'Publish remotion-architect 1.0.0'
priority: p2
section: ralvaskills
topics:
    - registry
    - skills
---
New skill `remotion-architect` 1.0.0 (`skills/frameworks/`, standalone, not bundled): Remotion
4.0.542 programmatic video. Frame-deterministic compositions, `interpolate`/`spring`,
`Series`/`TransitionSeries`, `@remotion/media` `<Video>`, Zod props and `calculateMetadata`,
CLI/SSR/Lambda/Player rendering, and a Company License check. `RECIPES.md` carries the
skeleton and render recipes; `STACK.md` pins 4.0.542. Commit `a0cdfc0` also moved the skill
count to 63 (`AGENTS.md`, `gen-data.mjs`, `facts.ts`) and added the spec and changelog entries.

`registry/index.json` was deliberately left out of that commit: the publish run diffs the skills
against the committed index, so a pre-regenerated index would have hidden the new skill.

Resolved 2026-11-09: the `Publish Registry` run for `a0cdfc0` (42642621026) finished green,
released `remotion-architect@v1.0.0` (15:14:45Z) and committed the regenerated index
(`df390b8`), which now lists the skill.

Open points that are not part of this release:

- Not yet verified: that `task site:build` / `site:check` pass with the 63-skill index now
  that the index matches (needs `NODE_AUTH_TOKEN`).
- The skill was written from the Remotion docs and has not been exercised on a real project;
  try it on one video before pinning it anywhere.
- The free-tier terms of the Remotion license were not read; the skill only tells the agent to
  raise the question and link https://www.remotion.dev/docs/license.

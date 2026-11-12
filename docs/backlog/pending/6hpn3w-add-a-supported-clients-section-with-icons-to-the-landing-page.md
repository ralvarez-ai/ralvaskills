---
title: 'Add a supported-clients section with icons to the landing page'
priority: p2
section: landing
topics:
    - site
    - frontend
---
The page now names all ten clients in the FAQ and the files tree, but the breadth only shows up as prose. A short section — one row of marks with a one-line "what each one gets" — would make it visible on the scroll and give the ten-row MCP table in the README a visual counterpart.

Scope:

- **Placement**: between the hero film and the MCP tools section, so the client list arrives before the tool list it explains.
- **Icons**: inline SVG, no runtime dependency and no hotlinking. Check each mark's licence before shipping it (the `simple-icons` set is CC0 and covers Claude/Anthropic, OpenAI/Codex, Gemini and Cursor; Kimi, Qwen, Droid/Factory and Antigravity may need a monogram fallback rather than an invented logo). Brand marks are trademarks: fall back to a neutral monogram when there is no permitted asset, and never redraw one.
- **Light and dark**: every mark must hold contrast on both palettes; a single-colour variant per theme is safer than full-colour brand art. Reuse `site/src/styles/tokens` rather than adding hex values.
- **Motion**: the row is static content; no animation beyond the existing `reveal` opacity, and it must pass with `prefers-reduced-motion`.
- **Copy**: one line per client naming the thing that differs — `.agents/skills` for the shared eight, the config file for each, and `--read-only` for every client except Claude Code.

Data source: extend `site/scripts/gen-data.mjs` to emit the tool list from `internal/tool` (the IDs registered in `tool.go`, or the `ParseID` set) and assert the count, instead of hand-writing ten names into `facts.ts`. That keeps `task site:check` failing when an adapter is added without the page, the same way it already guards skills, bundles, commands and MCP tools.

Done when the section renders statically (prerendered text, no JS-only content), `task site:check` covers it (Biome, `tsc`, build, axe at 390 and 1440, js-off), and the icon set is recorded with its licence and source in the section file or `THIRD_PARTY_NOTICES.md`.

Plan: [Agent Skills clients: one adapter, ten targets](../../plans/2026-11-12-agent-skills-client-roadmap.md).

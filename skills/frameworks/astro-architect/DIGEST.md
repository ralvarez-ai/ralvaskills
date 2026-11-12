# Astro Architect (digest)

Astro 6+ standards: zero-JS-by-default output, typed content collections, scripts, routing, forms, deployment.

## Always-on rules
- Default to `.astro` components; use a framework island only when the interaction cannot be an event handler on server HTML.
- Interactivity ladder, cheapest first: CSS `:checked`/`:target`/`:hover`; `<details>`/`<dialog>`; processed `<script>` in the `.astro` file; `client:visible`/`client:idle` island; `client:load` last resort.
- Every `<script>` ships to the browser. Processed scripts are bundled, deduped, hoisted; a `<script>` in `src/components/` renders once per page.
- Collections live in `src/content.config.ts`; each needs a loader (required) and a schema (optional but load-bearing).
- Import Zod from `astro/zod`. Import order: `astro:content`, `astro/loaders`, `astro/zod`.
- Dates: `z.coerce.date()`.
- `getCollection` is unsorted: sort explicitly. Filter in the query. Use `getEntry` for lookup by id.
- Routes are file-based; static routes beat dynamic ones.
- Theme flash guard: `is:inline` script in `<head>`.
- Scoped `<style>` by default; theme switching via `[data-theme]` custom properties, no hydrated picker.
- Images: `<Image>` from `astro:assets` on imported files; set width/height; remote images need `image.domains` or `remotePatterns`.
- Fonts: self-host with `@font-face`.
- Tests: Vitest for units; Playwright against `astro preview`, not `astro dev`.
- Set `site:` in `astro.config.mjs` at scaffold time.
- Stay `static` unless a route needs a runtime; SSR adds a process, cold starts and secret storage.

## Hard constraints
- Never ship a framework component for an interaction a processed `<script>` can do.
- Never `@import` a font from a CDN in CSS; it blocks render.
- Never hydrate a theme picker.
- Re-verify auth inside every Action handler; the calling context is not trust.
- Actions, Server Islands and Sessions need an adapter. Don't use them on a static site.
- Sessions are cookies: keep them small.

Invoke when: scaffolding or reviewing an Astro site, wiring content collections, adding forms or Actions, choosing an adapter or rendering mode, or auditing client JS leaks.

Full guide: invoke the astro-architect skill.

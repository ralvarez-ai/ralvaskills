# nextjs-architect (digest)

Next.js 16 standards: App Router only, server components by default, server actions for writes, streaming, edge vs node, Image/Font/Metadata APIs.

## Always-on rules
- App Router only; no Pages Router code unless mid-migration.
- Every component is a server component; put `"use client"` on the smallest leaf that needs hooks, browser APIs or handlers.
- Props crossing server to client must be JSON-serializable: IDs, plain objects, ISO strings for dates.
- Reads: direct from server components through a server-only module. Writes: server actions.
- Server actions: `"use server"`, validate input with Zod, re-verify the session inside the action, `revalidatePath`/`revalidateTag` after mutation, return `{ok: true, ...} | {ok: false, ...}`.
- `notFound()` for missing data (closest `not-found.tsx`).
- Streaming: `loading.tsx` and `error.tsx` per route; per-region Suspense for slow regions only.
- Runtime defaults to node; use edge only where latency justifies it (auth checks, redirects, A/B).
- Protected layouts check the session server-side; middleware for redirects and flags.
- Auth cookies: `httpOnly`, `Secure`, `SameSite=Lax`.
- Don't cache user-specific data globally; trust cache defaults until measured slow.
- Metadata via `generateMetadata` and root `metadata`; never hand-write `<head>`.
- `priority` only on above-the-fold images; remote images need `remotePatterns`.
- Only `NEXT_PUBLIC_*` env vars reach the client; the rest are server-only.
- Tests: Vitest + RTL for components, action functions as units, Playwright against `next start`, or `next dev` in CI (build then start).

## Hard constraints
- Never `<img>` in App Router code; always `next/image` with width/height or `fill`.
- Never load fonts with `@import` or `<link rel="stylesheet">`; use `next/font`.
- Never trust the client: server actions and API routes re-verify the session.
- Keep DB code in server-only modules so it never bundles to the client.

Invoke the full skill when: scaffolding or reviewing a Next.js app, auditing server/client boundaries, choosing edge vs node, designing caching, auth or middleware, or picking a deployment target.

Full guide: invoke the nextjs-architect skill.

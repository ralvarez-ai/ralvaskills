# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| next | 16.2 | Framework (App Router only) |
| react | 19.2 | UI library (pinned by Next) — see [react-architect](../react-architect/STACK.md) |
| react-dom | 19.2 | DOM renderer |
| typescript | 6.0 | Static typing — strict mode mandatory |
| zod | 3.x latest | Server action validation + boundary types |
| @vercel/otel | latest | OpenTelemetry instrumentation when deploying on Vercel; generic OTel SDK otherwise — see [observability-architect](../../infra/observability-architect/STACK.md) |
| next-auth | 5.x (Auth.js) | Session-based auth for first-party apps |
| @biomejs/biome | 2.4 | Lint + format |
| vitest | 4.1 | Unit + component tests |
| @playwright/test | 1.60 | End-to-end tests against `next start` |

## Notes

- **Server actions are the canonical mutation mechanism.** API routes earn their place for webhooks, external clients, streaming.
- **Hybrid data access** (per the project's choice): direct DB from RSC for reads, server actions for UI-driven writes (route handlers for webhooks and external clients) — but `react-architect` and project conventions can override per project.
- **Runtime: node default; edge** only when latency justifies it (auth checks, redirects, A/B routing).
- **Inherits the `react-architect` stack** for client-side patterns; this file pins only what's Next-specific.

_Last reviewed: 2026-05-21_
_Skill version at last review: 1.0.0_

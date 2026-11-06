# react-architect (digest)

React 19 client standards: TypeScript strict, feature-based folders, hooks-first composition, TanStack Query for server state, zustand for cross-tree client state, Suspense and ErrorBoundary at every async boundary, Radix for a11y.

## Always-on rules
- `strict` tsconfig with `noUncheckedIndexedAccess` and `exactOptionalPropertyTypes`.
- Use `unknown` for untyped boundaries, then narrow; `as` only after a validated parse (Zod).
- Discriminated unions over enum plus cast.
- Feature folder owns its UI, hooks, types, schemas and API queries; `index.ts` is the public surface.
- Zod `schemas.ts` per feature; parse at the API boundary.
- Function components only; one component per file; `<Component>Props` colocated; named exports preferred.
- Loading, error and empty states are real components.
- Event handlers passed as props; `memo`, `useMemo`, `useCallback` only when profiling or identity justifies.
- `useEffect` synchronizes with non-React systems; it does not fetch.
- Server state in TanStack Query: hierarchical array keys, `staleTime` tuned per resource, `refetchOnWindowFocus: true` default.
- Context for static-ish cross-tree state (auth, theme, locale) with a stable provider value.
- zustand only when Context performance or a non-trivial reducer justifies it.
- Suspense and ErrorBoundary above every async boundary; one per logical region.
- Forms: controlled inputs, Zod on submit, disable submit while pending.
- Styling: Tailwind by default, CSS Modules for one-offs, design tokens centralized.
- Code-split at route boundaries; virtualize lists over 100 rows.
- Tests: Vitest, RTL queries by role/label/text, Playwright e2e, axe in tests; `tsc --noEmit` is the health gate.

## Hard constraints
- No `any`.
- No class components.
- Never store fetched server data in `useState` or zustand.
- Never fetch in `useEffect`.
- No CSS-in-JS with runtime cost (styled-components, emotion).
- No Jest in new code.

Invoke the full skill when: writing or reviewing React components, hooks, client-side state, forms, routing or client test setup.

Full guide: invoke the react-architect skill.

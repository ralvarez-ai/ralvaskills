# ui-ux-architect (digest)

Implementation standards for accessible, consistent, responsive frontends: WCAG 2.2 AA, Radix + Tailwind 4 + shadcn/ui, design tokens, four UI states, forms, motion, icons, testing.

## Always-on rules

- Every shipped component meets WCAG 2.2 AA. Semantic HTML first; a `<button>` is a button.
- Every interactive element is keyboard-operable with a visible focus indicator. Do not `outline: none` without a replacement.
- Contrast: 4.5:1 for normal text, 3:1 for large text and UI components.
- Every input has a `<label htmlFor>`. Placeholder text is not a label.
- ARIA only where semantic HTML cannot say it. `aria-label` on icon-only buttons; `aria-hidden` on decorative icons.
- Respect `prefers-reduced-motion`. Motion is functional only; 100–300 ms; 0 ms under reduce.
- Use Radix primitives, Tailwind 4, and shadcn/ui recipes. Components live in your repo.
- Centralize visual decisions as semantic design tokens (OKLCH colors, 4 px spacing scale).
- Mobile-first; container queries when layout depends on the container; touch targets 44×44 px.
- Every surface that touches data has loading, error, empty and success states.
- Error UI: plain language, retry, visible correlation ID, no stack traces.
- Forms: validate on blur; errors below the input via `aria-describedby`; disable submit with `aria-disabled`.
- One icon set per project (Lucide default); SVG, not icon fonts.
- Set `<html lang>` always. Route copy through a translation function if i18n is planned.
- Test: axe-core in Vitest and Playwright with zero WCAG 2.2 AA violations; Lighthouse a11y 95 or higher; manual keyboard pass.
- Start at maturity Level 1 (shared primitives in `components/ui/`).

## Hard constraints

- Never `<div onClick>` where a `<button>` belongs.
- Never use icon fonts.
- Never use the `disabled` attribute on a pending submit button; use `aria-disabled="true"`.

Invoke the full skill when: designing UI components, building a design system, or auditing accessibility; need the component recipes, token template or test setup.

Full guide: invoke the ui-ux-architect skill.

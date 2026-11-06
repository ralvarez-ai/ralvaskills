# reveal-js-architect digest

Builds reveal.js 6.0.1 decks as real source: `<section>` or Markdown-mode slides, vanilla by default, React extension via @revealjs/react 0.2.1.

## Always-on rules

- Interview first, in order, waiting for each answer: topic + slide count/duration; vanilla or React; HTML or Markdown-mode; theme; plugins; delivery (static file or npm + Vite).
- Output is reveal.js source (`index.html` or Markdown slides, plus `main.js` config), never a slide-spec document. Don't write slide copy.
- Horizontal `<section>` = new topic; vertical = sub-point; vertical nesting max one level.
- One idea per slide; split rather than cram. Fragments build within a slide, not instead of splitting.
- Speaker notes go in `<aside class="notes">`.
- Start from a built-in theme; custom CSS goes in a separate stylesheet loaded after it.
- Register only plugins the deck uses. `RevealMarkdown` when any slide uses `data-markdown`; `RevealNotes` when notes are used.
- Keep defaults `hash`, `controls`, `progress`, `center`; `transition: slide`; default 1280x720; `slideNumber` only for reference decks.
- Non-white theme: override `.reveal-viewport` background.
- Make content fit by tightening type or splitting; `overflow-y: auto` is a fallback only.
- Fixed canvas height: px on an inner wrapper, matching `Reveal.initialize({ height })`.

## Hard constraints

- Never set `position`, `overflow` or `display` on `<section>`; style a wrapper `<div>` inside it. Breakage is silent (blank slides on advance).
- Never set `height` on `<section>` or use `%` heights inside it.
- Nested `em` font-sizes multiply; size the inner region in `rem` or px, or reset the child to `1em`.
- Don't mix HTML and Markdown-mode in one deck unless one section genuinely needs HTML control.
- Don't reach for the React extension by habit; use it when the deck lives in a React app or needs component-driven slides.

Invoke the full skill when: structuring a deck, choosing theme, plugins or config, setting up vanilla or React projects, or reviewing deck source.

Full guide: invoke the reveal-js-architect skill.

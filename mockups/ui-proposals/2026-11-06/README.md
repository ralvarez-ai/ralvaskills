# Palette round 2 (2026-11-06)

**Decides:** which colour direction replaces Ember (warm paper, vermilion) on the landing page.

Open `palettes-round-2.html` from `file://`. The bar at the bottom switches the palette; the page's own
button switches light and dark. Force a variant with the query string, for example
`palettes-round-2.html?palette=teal-sand&theme=dark` (the query wins and is not stored).

It is the approved combined page from `../2026-11-04/direction-4-combined.html` with each palette as a
`html[data-palette=<name>]` override block, the preview bar from the `mockup-architect` skill inlined, and the
sample catalogue count updated to 60. Regenerate with `python3 -I _tools/palettes2.py` (it reuses the WCAG gate
and helpers of the previous round and refuses to write if any pair fails).

| Palette | Idea | Accent (light / dark) | Human line |
|---|---|---|---|
| Ember (today) | Warm paper, vermilion: the baseline | `#b3360c` / `#ff7a45` | amber |
| Blueprint | Cobalt on cool paper, navy night | `#1f4fd8` / `#7aa2ff` | amber |
| Teal & Sand | Deep teal on warm sand | `#0b6e6a` / `#4fd1c0` | coral |
| Claret | Raspberry on a blush ground | `#a4173f` / `#ff7d9c` | teal |
| Graphite & Gold | Neutral graphite, brass gold | `#8a5a00` / `#e6b93d` | teal |
| Mono & Magenta | Black and white, one magenta | `#b4147a` / `#ff6ec4` | amber |
| Emerald | Emerald on cool white, deep green night | `#0b7a4b` / `#3ddc97` | amber |

Avoided on purpose: the cyan-on-grey of the original, violet (Aurora, round 1), electric lime (Signal) and
vermilion/orange (Ember).

## Checks

- WCAG 2.2 AA by script: text, muted, accent and the four status colours at least 4.5:1 on page, band and card;
  the label on an accent fill at least 4.5:1; a strong border at least 3:1. All seven pass in both themes.
- axe-core (wcag2a, 2aa, 21a, 21aa, 22aa): 0 violations in 28 runs (7 palettes x 2 themes x 1440 and 390 px).
- No console errors; opens from `file://`. Screenshots of the first viewport: `shots/`.

## Recommendation

**Teal & Sand**, with **Graphite & Gold** as the quiet alternative. Teal & Sand keeps the page warm, so it
does not feel like a different product, but moves the hue out of the orange family you wanted to leave, and the
coral "You approve." keeps the human line distinct from the accent. Graphite & Gold is the most restrained.
Blueprint is the safest and the least distinctive; Mono & Magenta and Claret are the boldest.

Nothing here touches the site: porting a choice means editing the Ember layer at the end of
`site/src/styles/tokens.css` and `palettes.py`'s twin values.

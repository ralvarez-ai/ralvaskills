# Explainer visual language

How a scene looks and moves, whatever draws it. The scene rule (one idea, drawn in the order it is said) is in [SKILL.md](SKILL.md); the engines' code is in [manim-architect](../../frameworks/manim-architect/SKILL.md) and [remotion-architect](../../frameworks/remotion-architect/SKILL.md). The rules below are rules of thumb from reviewing finished videos: judge each against the frame, at the size the video is shown.

## 1. One focus at a time

- What is new is bright; what was already explained stays but fades back (about 30% opacity) when a new object takes the focus.
- Reveal in the order the narration names things. A frame where everything is on screen and the voice is still on the first object is a bug.
- About five to seven separate objects are what a viewer can follow in one scene. Group the rest (a box around a family, a column) or split the scene.
- Emphasis is a pulse of scale or a change of stroke, never a change that fills a shape with colour.

## 2. A colour means one thing

- Pick the meaning of each colour once for the whole video and keep it: for example green for passed or confirmed, red for refused or failed, amber for a warning or a stand-in, the brand colour for the subject. One entity (a kind of record, a role) keeps its colour in every scene.
- Show a legend the first time a colour carries a meaning, and never use red or green as decoration.
- Colour is never the only signal: the label, the shape or the position says the same thing, so a grey print or a viewer who does not see the difference loses nothing.

## 3. The line says what is real

- A solid outline is a thing that exists and sits inside; a dashed outline is something outside the system, planned, or not built.
- An arrow has a direction and a label above it; a thick line is the main path, a thin one a side path.
- A made-up value and a made-up name are flagged on screen and in the words; a chip or tag names the origin and the maturity ([SKILL.md](SKILL.md), section 4).

## 4. Type and space

- Three sizes are enough: a title, a label and a note, plus a smaller one for a short status tag. Set the floor for the smallest text from what is readable where the video is shown (a phone, a projector), and count the labels under it before a render.
- Identifiers, paths and code are in a monospaced face; sentences are in the text face. Align left edges inside a drawing.
- Keep bands: a title, the drawing, an optional one-line summary, a footer with the sources. The drawing fills its band: an empty third of the frame means it can be larger.
- An object keeps its place from scene to scene (the same box in the same spot), so the viewer finds it without searching.

## 5. Motion

- A thing that moves is something that travels (a request, a record, a token); a fade is an appearance; a pulse is emphasis. Do not animate everything.
- Stagger a list: each item a third of a second after the one before, as the voice reaches it. One to two seconds of stillness before a new dense drawing is spoken over.
- A check that can fail is drawn as an item that passes a list of checks, and a failure is drawn as a return with its reason. A message between two parties is an arrow that travels.
- Keep a shared object across two scenes, and let it change, instead of redrawing it: the redraw flashes and the second copy drifts from the first.
- A video of more than about five minutes carries an orientation device, a small map of the system with the current part lit.

## 6. Ways to show something

Pick the one that matches the idea ([RECIPES.md](RECIPES.md#7-the-shape-of-a-scene) has the table by what the narration says):

| Idea | Way to show it |
|---|---|
| a formula | built piece by piece, each piece coloured like the part of the drawing it names |
| the same example everywhere | a small panel in a corner that keeps its current values |
| a path through states | a token that walks the graph |
| a verdict from a table | the row that decides lights up |
| a fixpoint or a spread | a wave that fills in round by round |
| repetition over time | a timeline whose values change |
| a claim and its counterexample | two panels side by side |
| a measure that grows | a counter or a dial that counts up |
| a status on a statement | a stamp that lands on it |
| a question | the question first, a beat of silence, then the answer |
| the end of a chapter | a strip of one icon per idea |
| a request, a check, a message | something that travels along the path |
| what exists and what does not | columns by maturity |

## 7. Three checks

- **Squint:** blur the frame. Is the focus still obvious? If not, something competes with it.
- **Phone:** play it at the size of a phone. Is every label you need readable?
- **Grey:** drop the colour. Does every meaning still read from the label, the shape or the place?

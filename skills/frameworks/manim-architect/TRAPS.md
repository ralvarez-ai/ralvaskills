# Manim traps and worked patterns

Details behind [SKILL.md](SKILL.md). Read it when a scene keeps objects across scene changes, mixes images with vector shapes, or teaches a formula with a real number.

## 1. Objects that stay on screen across scenes

A video often keeps one drawing (a plot, a table) and changes only what sits around it. The base class fades out "everything that is not kept". Five traps:

- **Keep the whole family.** An animation may add parts one by one, or wrap them in an anonymous `AnimationGroup` group. The scene's `mobjects` list then holds a wrapper, not the object you marked. Compare the *family* (`get_family()`) of each scene mobject with the family of each kept one, and keep a mobject when any member matches. Comparing the top-level objects only fades out things you meant to keep.
- **Set `z_index` on the group.** A group is drawn by its own `z_index`, not by its parts'. An anonymous group made by an animation has `z_index = 0`, so a kept layer can end up above or below the wrong thing. Give every group the layer you want, and use fixed layers: plot frame lowest, region or heat images next, lines, points, markers on top.
- **Images do not go in a `VGroup`.** `ImageMobject` (a raster made from a grid, a heat map, class regions) needs `Group`. A `VGroup` raises at `add`.
- **Do not name an attribute `points` on a `VGroup` or `Mobject`.** Manim uses `.points` for its geometry; assigning your own list silently breaks drawing. Use `markers`, `dots`, `members`.
- **Never fade out what a scene change already removed.** If the base class fades everything not kept when a scene closes, a `FadeOut(old_things)` at the start of the next scene makes Manim put those things back at full opacity and fade them again: a flash of the previous scene. Drop a `FadeOut` whose target is no longer in the scene (override `play`, and keep its time as a wait), narrow a group to the parts still shown, and check every render for flashes: decode to small frames and look for a jump that returns within about half a second (a normal animation changes a little every frame; a flash is one big step in and one out).

## 2. A worked number before the general formula

A formula is easier to read once the viewer has used it on one concrete case. The pattern, in the order the narration says it:

1. Pick one real item from the data (a point, a row, a token). Choose it by a rule you can state (the one nearest a fixed spot, the largest, the first), and write the rule in the content file so the choice is not arbitrary.
2. Highlight it in the drawing (scale it up, dim the rest).
3. Draw what the formula measures (segments to the candidates, a bar, an area).
4. Show the formula with symbols, then, line by line, the same formula with the numbers put in, one line per candidate. Draw each line together with its segment, and colour both the same.
5. Mark the result (the smallest, the largest, the sum) with colour and one short phrase.
6. Remove the example and apply the rule to every item at once.

Numbers on screen are computed from the data the scene already holds, rounded the same way everywhere, and the scene's chip is `DERIVED` (they are a calculation, not a source's claim).

## 3. Showing a squared quantity

A square of a distance is easier to believe when it is a *square*.

- Build a square on the segment: with the segment from `a` to `b`, `edge = b - a`, `normal = [-edge[1], edge[0], 0]`, polygon `a, b, b + normal, a + normal`. Fill it at low opacity.
- Show the item with the largest distance first, with `d` and `d²` printed. A far item gives a big square and a near one a small square, which is the reason the squared sum is dominated by far items.
- Fade it out before the sum is drawn, so the sum is read as "all the squares added up" and not as a new object.

## 4. Draft sweep

Overlapping text, a missing set of dots and an empty half of the frame were real bugs in the first drafts and are cheap to find at draft quality.

- After each render, take a frame near the end of every scene (`ffmpeg -ss <t> -i video.mp4 -frames:v 1`) and look at it before the next edit.
- Check the right column against the left drawing: a text block taller than the column overlaps the footer band.
- Check that every object the narration names is on screen when it is named.

# Videos for learners

Additions to [SKILL.md](SKILL.md) for a video that teaches a concept to people who are meeting it for the first time, usually inside a course or a tool.

## 1. One worked number per concept

For each formula or rule, the video works one real case before it states the general form: choose an item from the data by a stated rule, show the quantities with their numbers, mark the result, then apply the rule to everything. Drawing pattern: [manim-architect TRAPS.md](../../frameworks/manim-architect/TRAPS.md#2-a-worked-number-before-the-general-formula). The scene chip is `DERIVED`, and the content file holds the choice rule and the inputs.

## 2. Show the reason for a term

When a formula squares, averages or normalises something, a scene shows what that does (a square built on a segment, a mean as the balance point, a bar scaled to one). If a term cannot be drawn, the voice says so and gives the reason in a sentence.

## 3. Predict, then reveal

Before a result, give the viewer a beat to guess: the voice asks the question, the picture holds, the answer follows. Use it where a wrong guess is common (which centroid is nearest, where the curve bends). One or two per video; more is a quiz.

## 4. Pointing to a tool is a switch

A video may end by sending the viewer to an app, a lab or a repository. Make that an option (a flag at render and narration time) that is off by default, and keep the words in their own field (`labRef`), so the same take can ship with or without it. The off version must not leave a gap: the last scene still ends on a stable frame.

## 5. Order of work for a new topic

1. Write the one thing the viewer must be able to say afterwards.
2. List the quantities the formula uses and decide how each is drawn.
3. Pick the worked case and its rule.
4. Write the narration at the conversational dial ([narration-architect CONVERSATIONAL.md](../narration-architect/CONVERSATIONAL.md)).
5. Draw, cued to phrases; render a draft; look at every scene's end frame.
6. Measure the voice, set the rate, render at final size.

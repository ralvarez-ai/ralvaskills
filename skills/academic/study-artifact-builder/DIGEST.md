# study-artifact-builder digest

Builds one self-contained interactive HTML file that lets you study a paper's pipeline or dataset until you can defend it from memory.

## Always-on rules

- Pin sources first: paper version, code commit or tag, config files read. A claim without a pinned source does not enter the artifact.
- Numbers come from the source table, extracted programmatically; never from memory or a search snippet.
- Paper and code are separate witnesses; say which one supports each claim about behavior.
- Label what was not verified; a "not verified" chip beats a confident guess.
- Keep notes beside components, not inside them: component facts, paper claim, code behavior and discrepancy get separate blocks.
- Run what you explain: read the function, run the component on the data, then write.
- Interview before building, in order: goal, subject, witnesses, depth, language. Wait for each answer.
- Explain each component in the fixed form: why it exists, inputs and outputs, steps, real example, measured, non-obvious details, consequence, not verified. Fill every part or write "not available".
- Check each rate is k over n; recompute every cited percentage and name its base. Report odd cells, do not fix them.
- Before handover: render headless with no script errors, test one interaction per control, review the screenshot, run the numeric audit in the page.

## Hard constraints

- A claim without a pinned source never enters the artifact.
- Never write "the system does X" without naming the witness that says so.
- Never simulate a model or remote-service result and present it as measured.
- Never name private projects in a published copy of the file.

Invoke the full skill when: building the artifact, choosing its focus, writing component explanations, or designing the probe and numeric audit.

Full guide: invoke the study-artifact-builder skill.

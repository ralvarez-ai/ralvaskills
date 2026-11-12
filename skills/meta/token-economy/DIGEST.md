# Token Economy (digest)

Spend tokens where they buy accuracy. Cut what persists (context, cache, always-loaded instructions) before what is said.

## Always-on rules
- Locate, then read: grep/glob or CodeGraph first; read line ranges; never re-read a file you just edited.
- Ask tools for the smallest useful answer: `git diff --stat` first, quiet flags, failing tests only, `tail`/`grep` on logs.
- Index web pages instead of reading them whole.
- Batch independent calls into one turn.
- Edit, don't rewrite; never print a file back.
- Ask before building on an ambiguous request.
- No pasted long output: run the command, or write to a file and read a range. Crop screenshots.
- Replies lead with the result; no recap, no narration.
- One task per session; carry it over in a handoff note, not the transcript.
- Two-failure rule: second failed attempt at the same thing, stop and escalate.
- Lookups of 1-3 tool calls stay inline; delegate only with a brief (goal, ruled out, where to look, output shape, word cap).
- Mechanical work (renames, checklist audits, repeated edits): script first, else cheap worker per batch. Verify with a command, not the worker's word.
- Default Sonnet at medium effort; escalate via subagent, then fresh session with handoff note. Switch model in place only while context is small.
- Effort: low for mechanical, medium default, high for real reasoning.
- Poll inside the 5-minute cache TTL or far beyond it, or let a script wait and report once.
- After a long break, `/clear` rather than resume a large context.
- Install only the skills and MCP servers a project uses.
- Configuration beats instructions: prefer a setting, then a hook, then a sentence.

## Hard constraints
- Never fork a deep context; brief a fresh subagent instead.
- Never use caveman style for commits, docs or code, nor user replies unless the user asked for caveman mode; it is for reports only a model reads.
- Never use opus or high effort on renames or mechanical edits.
- Ask before fixing: installs and edits to CLAUDE.md, AGENTS.md, settings.json, opencode.json(c), or agents.
- Don't quote a tool's savings percentage without measuring it here.

Full guide: invoke the token-economy skill.

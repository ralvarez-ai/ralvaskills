# Improve Codebase Architecture (digest)

Finds deepening opportunities (shallow modules made deep) for testability and AI-navigability, informed by `CONTEXT.md` and `docs/adr/`.

## Always-on rules
- Use the glossary exactly: module, interface, implementation, depth, seam, adapter, leverage, locality. Don't drift into "component", "service", "API" or "boundary".
- Deletion test: if deleting a module makes complexity reappear across callers, it earns its keep; if complexity vanishes, it was a pass-through.
- The interface is the test surface.
- One adapter is a hypothetical seam; two adapters is a real seam.
- Read `CONTEXT.md` and `docs/adr/` first. If they are missing, proceed silently.
- Hygiene-tool evidence is input, not a verdict.
- Present candidates as one self-contained HTML report in the OS temp directory, `architecture-review-<timestamp>.html`. Open it and give the user the absolute path.
- Each candidate card: files, problem, solution, benefits (locality, leverage, tests), before/after diagram, strength badge (`Strong`, `Worth exploring`, `Speculative`).
- End the report with a Top recommendation.
- Surface an ADR conflict only when the friction justifies reopening it; mark it clearly.
- Name new concepts in `CONTEXT.md` as they crystallize. Create the file lazily.
- Offer an ADR on rejection only when a future explorer would otherwise re-suggest the same thing.
- Deepen by dependency category: in-process, local-substitutable, remote-but-owned (ports and adapters), true-external (mock).
- Replace old shallow-module tests; don't layer new ones on top.

## Hard constraints
- Do NOT propose interfaces until the user picks a candidate. Ask "Which of these would you like to explore?"
- Never write the report into the repo.
- Introduce a port only when at least two adapters justify it.
- Don't list every theoretical refactor an ADR forbids.

Invoke when: the user wants to improve architecture, find refactoring opportunities, consolidate tightly-coupled modules, or make a codebase more testable and AI-navigable.

Full guide: invoke the improve-codebase-architecture skill.

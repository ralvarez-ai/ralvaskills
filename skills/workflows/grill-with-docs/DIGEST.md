# Grill With Docs (digest)

Grilling session that stress-tests an existing plan or design against the domain model, sharpens terms, and updates CONTEXT.md and ADRs inline.

## Always-on rules
- Interview relentlessly about every branch of the design. Ask one question at a time, each with your recommended answer. Wait for feedback before the next.
- If the codebase can answer a question, explore the code instead of asking.
- When the user uses a term that conflicts with CONTEXT.md, call it out immediately. Propose one canonical term; list the others as aliases to avoid.
- Stress-test domain relationships with concrete scenarios.
- When the user describes how something works, check the code. Surface any contradiction.
- Update CONTEXT.md inline as each term resolves. Do not batch.
- CONTEXT.md is a glossary only: no implementation details, specs, scratch notes or implementation decisions. Include only terms specific to this context.
- Offer an ADR only when all three hold: hard to reverse, surprising without context, the result of a real trade-off.
- Create files lazily: CONTEXT.md when the first term resolves, `docs/adr/` when the first ADR is needed.
- If CONTEXT-MAP.md exists, read it first. If the topic's context is unclear, ask.

## Hard constraints
- Never treat CONTEXT.md as a spec, scratch pad or repository for implementation decisions.
- Never offer an ADR unless all three criteria hold.

Invoke the full skill when: a plan or design already exists and needs stress-testing against the domain language. If no plan exists yet, use feature-planner first.

Full guide: invoke the grill-with-docs skill.

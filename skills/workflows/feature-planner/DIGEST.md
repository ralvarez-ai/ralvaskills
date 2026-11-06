# Feature Planner (digest)

First step of the new-feature pipeline: turns a vague request into a requirement, constraints, a design outline and ordered vertical slices. Hands off to grill-with-docs, then tdd.

## Always-on rules
- Each step (requirement, constraints, design outline, tasks) is a user checkpoint. Do not advance until the current step is confirmed.
- Requirement: who, the outcome as a present-tense fact, the success criterion, and what is out of scope. No design in it. Two paragraphs at most.
- Record constraints even when the answer is "no constraint". Check CONTEXT.md and `docs/adr/` before designing.
- Design outline: shape, new seams, risks (highest first), alternatives. Use pattern names only when they earn it.
- Tasks are vertical slices, each one red-green-refactor loop, ordered by risk: tracer bullet first, riskiest assumption second.
- Each slice ends in a test that fails before and passes after. Never slice horizontally (migration, then repository, then service, then handler, then tests).
- When the plan changes mid-implementation, re-plan the remaining slices.
- Keep the plan shorter than the code it describes.

## Hard constraints
- Skip planning for trivial changes, exploratory spikes, bug fixes with a reproduction, and pure refactors.
- Never write all tests first, then all implementation.

Invoke the full skill when: the user says "I want to build X" or starts non-trivial work with a vague shape and no plan yet. If a plan or design already exists, use grill-with-docs instead.

Full guide: invoke the feature-planner skill.

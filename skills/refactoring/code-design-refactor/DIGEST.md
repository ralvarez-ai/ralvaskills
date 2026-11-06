# Code Design Refactor (digest)

Design-level refactoring of functions, classes and modules within one package: extraction, decoupling, SRP, encapsulation, primitive obsession. Sits between logic-cleaner and improve-codebase-architecture.

## Always-on rules
- Tests first. Write characterization tests for current behavior before any structural change.
- One move at a time. Each move leaves the tree green. Commit each successful move as `refactor(scope): <move>`.
- Extract by conceptual seam, not by count. Extract a class when a function takes 6+ arguments that mean little without each other.
- Stop extracting when the piece exists only to be called from one place.
- Decouple two-way imports into one-way. Pass dependencies in; do not fetch them.
- Introduce a port only by the rule of three (one consumer: no; two: maybe; three: yes).
- SRP: one reason to change. `<Noun>Manager`, `<Noun>Service`, `<Noun>Helper` names hide responsibilities.
- Encapsulation: private fields by default. Intention methods (`activate()`) over setters. Tell, don't ask.
- Primitive obsession: wrap validated primitives (email, currency code) in immutable value objects.
- Stop when the next move's win is not obvious.

## Hard constraints
- Do not refactor code that has no tests. Write characterization tests first.
- Do not mix a refactor with a behavior change in the same diff.
- Do not polish code that is about to be deleted.
- Cross-module restructuring goes to improve-codebase-architecture. Do not expand scope.

Invoke the full skill when: asked to clean up, tidy, refactor or restructure a specific function, class or module; a review finds a too-big function or a class doing too much; a unit is hard to test because its shape is wrong.

Full guide: invoke the code-design-refactor skill.

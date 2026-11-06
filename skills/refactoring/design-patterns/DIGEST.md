# design-patterns digest

Skeptical catalog for Go and Python: which patterns to keep, avoid, or hand to another skill. Consult before adding a pattern; not a refactoring step.

## Always-on rules

- A pattern costs a name, at least one type or function, and a promised benefit. Weigh before adding.
- Rule of three: extract at the third place that wants the same behavior.
- Checks in order: is the problem present; would a language feature do; will there be a second implementation; does it make code easier to delete; does the name carry weight. If the first or third answer is no, skip.
- Keep when the need is real: Repository (aggregates needing persistence); Adapter (two or more real implementations); Strategy as function value; Decorator/middleware for cross-cutting concerns; Observer as domain events (several unknown consumers); Builder (many optional params); Factory method (type unknown at compile time). Iterator: use the language feature.
- Replace: Singleton with dependency injection (exceptions: runtime singletons you don't control, `os.Stdout`, `logging.getLogger()`); Abstract Factory with discriminated union + parse function; Visitor with `match` or type switch; Chain of Responsibility with middleware `func(handler) handler`; Mediator with domain events or direct calls; Template Method with composition; Prototype/Memento with native cloning and serialization.

## Hard constraints

- Never introduce a pattern speculatively.
- Interface-introducing patterns (Adapter, Strategy, Repository) only pay back with at least two implementations.

## Smells (judgment, not hard constraints)

- Go: `any` everywhere, empty struct receivers, returning interfaces from constructors.
- Python: `@staticmethod`-only classes, inheritance for reuse, metaclass magic, `*args/**kwargs` passthrough.

Invoke the full skill when: a pattern is proposed for a change, a cargo-cult pattern needs pushback, or the topic belongs to an owner skill (Repository and aggregates, ports and adapters, middleware, DI).

Full guide: invoke the design-patterns skill.

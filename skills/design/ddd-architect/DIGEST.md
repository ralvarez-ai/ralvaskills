# ddd-architect digest

Domain-Driven Design as practical structure: strategic decisions (bounded contexts, context mapping) first, tactical building blocks second. Event sourcing is out of scope.

## Always-on rules
- Strategic first: a bounded context is the scope where a term has one meaning; one context = one `CONTEXT.md`; multi-context repos get `CONTEXT-MAP.md` at the root.
- A context is a language boundary, not a microservice; a context owns its language.
- Name the integration pattern for every context pair; default to ACL for any external system you don't own.
- Record each relationship in `CONTEXT-MAP.md`, one line per pair.
- Entities: identity by `id`, equality by ID. Value objects: no identity, immutable, equal by attributes; reach for them aggressively (`EmailAddress`, not `string email`).
- Aggregate = transaction boundary; default small; invariants decide membership.
- One repository per aggregate root; collection semantics; named query methods; return aggregates, not DTOs.
- Repositories use sql-architect (psycopg + `.sql` or sqlx + `//go:embed`); no ORM by default.
- Domain events: past tense, immutable, minimal payload (IDs + facts subscribers need), published by the aggregate root after commit.
- Start logic as a method on the aggregate; promote to a domain service only when it spans aggregates.
- Application services orchestrate use cases and hold no domain logic.
- ACL: wrap the external model at the boundary, validate aggressively, one ACL per external system in the infrastructure layer.

## Hard constraints
- Don't share types across contexts to "save duplication"; the meaning is different.
- Reference other aggregates by ID, not by object.
- Don't publish events inside the transaction.
- Default to ACL for any external system you don't own; never let the external model leak inward.
- Skip DDD for CRUD-only services, throwaway prototypes, and pure data pipelines.

Invoke the full skill when: shaping a domain model, drawing context boundaries or a context map, choosing aggregate boundaries, or writing repositories and ACLs.

Full guide: invoke the ddd-architect skill.

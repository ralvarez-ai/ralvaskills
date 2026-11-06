# hexagonal-arch digest

Ports and adapters: the domain core knows nothing of the outside; dependencies point inward. The domain declares ports; adapters implement them. Pairs with ddd-architect.

## Always-on rules

- Domain imports nothing outward: no framework, DB driver or HTTP library in domain files.
- Ports are interfaces the domain declares, named in domain terms (`UserRepo`, not `PostgresClient`).
- Adapters are concrete, one per technology, and import the domain. Never the reverse.
- Go: the interface lives in the consuming (domain) package. Python: a `typing.Protocol` in the domain module.
- Name adapter folders by direction (`primary`/`secondary` or `driving`/`driven`), never "ports".
- Application services depend on ports; adapters are injected. HTTP and CLI handlers call application services.
- Behaviour lives on domain types (`User.ChangeEmail`). Adapters persist; they do not decide business rules.
- Port methods use domain types only. Translate infrastructure types inside the adapter.
- Introduce a port only with two adapters (a test fake counts). One adapter is a hypothetical seam: inline it.
- Tests: domain with no infrastructure; application services with fakes; each adapter with an integration test.

## Hard constraints

- Domain never imports framework, DB or HTTP packages.
- Dependency direction is inward. Check it by grepping domain imports.
- A flat or nested folder layout is not the pattern. Import direction is.

Invoke the full skill when: designing service structure, placing interfaces or ports, deciding whether a seam is worth it, or reviewing seam cleanliness.

Full guide: invoke the hexagonal-arch skill.

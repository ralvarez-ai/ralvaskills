# TDD (digest)

Test-driven development with a red-green-refactor loop. Tests verify behavior through public interfaces, not implementation. Language-agnostic.

## Always-on rules
- Plan first: confirm the interface changes and the behaviors to test, prioritized. Get user approval on the plan.
- Vertical slices via tracer bullets: one test, minimal code to pass, repeat. First, one end-to-end test that proves the path.
- Write only enough code to pass the current test. Do not anticipate future tests.
- Test names describe behavior (what), not implementation (how). Use the public interface only. A test should survive an internal refactor.
- Mock only at system boundaries: external APIs, time and randomness, sometimes the database or filesystem.
- Design for testability: inject dependencies, return results instead of side effects, keep the surface small.
- Refactor only after GREEN: extract duplication, deepen shallow modules, introduce value objects, apply SOLID where natural. Run tests after each step.
- Per-cycle checklist: behavior not implementation; public interface only; survives internal refactor; minimal code; no speculative features.

## Hard constraints
- Never write all tests first and then all implementation (horizontal slicing).
- Never refactor while RED. Get to GREEN first.
- Never verify through external means that bypass the interface (for example, querying the database directly).
- Do not mock your own classes, modules or internal collaborators.

Invoke the full skill when: the user wants to build a feature or fix a bug test-first, asks for red-green-refactor, integration tests, or test-first development.

Full guide: invoke the tdd skill.

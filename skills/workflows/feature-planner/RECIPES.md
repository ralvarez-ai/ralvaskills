# Feature Planner — Recipes

## Horizontal vs vertical slicing

**Wrong (horizontal slicing — common, painful):**

1. Add the database migration.
2. Add the repository code.
3. Add the service code.
4. Add the HTTP handler.
5. Add the tests.

**Right (vertical slicing):**

1. **Tracer bullet:** a stub endpoint returns a hardcoded response; one test passes end-to-end (proves the wire works).
2. **Riskiest assumption:** a single happy-path slice runs the real logic, hits the real DB, returns real data; one test.
3. **Edge case A** — one more test, the smallest amount of code to make it pass.
4. **Edge case B** — same.
5. **N. Refactor the now-duplicated bits** as the third occurrence appears ([logic-cleaner §6](../../refactoring/logic-cleaner/SKILL.md#6-duplication-rule-of-three)).

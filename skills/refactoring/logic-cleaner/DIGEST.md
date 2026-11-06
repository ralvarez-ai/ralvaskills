# Logic Cleaner (digest)

Expression-level cleanup, language-agnostic: control flow, conditionals, naming, magic values, duplication. Not architecture.

## Always-on rules
- Maximum 2 levels of indentation per function; extract deeper logic into helpers.
- Guard clauses with early returns at the top; keep the happy path left-aligned.
- Fail fast: halt on invalid states instead of passing nulls or error flags down the stack.
- Positive phrasing (`isPending`, not `isNotFinished`); avoid double negatives.
- Conditions with more than two logical operators -> named boolean variables or functions.
- Apply De Morgan's Laws to negated compound conditions.
- Declare variables close to first use; never shadow outer variables.
- Treat parameters as read-only; assign results to new variables.
- Avoid abbreviations unless universal (e.g. `ctx` in Go); use symmetric pairs (`start`/`stop`, `open`/`close`).
- Function names begin with strong, specific verbs (`CalculateTotal`, not `ProcessData` or `Handle`).
- Extract raw strings and numbers (except 0 and 1) into named constants that explain why they exist.
- Inject environment-dependent values (timeouts, retry limits, paths) via configuration.
- Rule of Three: duplicate once; refactor into a shared function on the third occurrence.

## Hard constraints
- Expression-level only. If a finding needs module boundaries or class structure changed, stop and surface it; don't expand scope.

Invoke the full skill when: asked to clean code, simplify expressions or flatten nesting, and you need the rule behind a specific case.

Full guide: invoke the logic-cleaner skill.

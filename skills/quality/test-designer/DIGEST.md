# Test Designer (digest)

Test suite design: logic tests separated from e2e tests, no time-based waiting, injectable clocks, readiness-gated e2e, redundancy pruning, flake triage.

## Always-on rules
- Two tiers, one per test file: logic (in-process, fakes for I/O, no network or real clock, whole tier in seconds) and e2e (few, real boundaries, journeys only). Separate directories or tags, commands and CI jobs.
- Prove rules and edge cases in logic tests; e2e never re-proves them.
- Never wait on time: no `sleep`, fixed timeouts or elapsed-time asserts. Await the signal, use a fake clock, or poll a condition with one generous deadline. Timeouts are ceilings, not expected durations.
- Inject clock, randomness and IDs into code that depends on them; test backoff as the sequence of requested delays.
- E2E: gate on readiness (health endpoint, open port, visible element), own and clean up its data, stay parallel-safe, stub third parties at the network edge.
- Mock only at system boundaries; prefer in-memory fakes over call-counting mocks.
- Prove each rule once, at the lowest tier. Delete tautological tests (mock-call asserts, expected values computed by the code under test), merge near-duplicates into table-driven cases, skip logic-free code.
- Triage a flake before touching it: hidden wait, shared state, ambient time, real race, environment.

## Hard constraints
- No retry or raised timeout to hide an undiagnosed flake; a retry needs a tracked owner and deadline.
- No `sleep` in a test as a synchronization mechanism.
- Load, stress and soak testing are out of scope.

Invoke the full skill when: designing or splitting a test suite, a test sleeps or flakes, or reviewing a PR that adds tests.

Full guide: invoke the test-designer skill.

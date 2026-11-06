# Event-Driven Architect (digest)

Async event services: event and command taxonomy, Protobuf schemas, topic naming, outbox, idempotency, DLQs, backpressure, schema evolution. Broker-agnostic.

## Always-on rules
- Events are past-tense facts (`OrderPlaced`). Commands are imperative with one handler (`PlaceOrder`). A command with several handlers is an event; rename it.
- Envelope on every event: `event_id` (UUID v7), `occurred_at` (RFC 3339), `aggregate_id`, `schema_version`.
- Payloads are minimal: IDs and the facts subscribers need, not full aggregate state.
- Topic: `<org>.<context>.<resource>.<version>.<event_name>`. Version lives in the topic name. One case convention per broker.
- Events are `.proto` messages, validated by protovalidate, guarded by `buf breaking` in CI.
- Partition by aggregate ID. Order is per key, not global.
- Consumers dedupe on `event_id` and make side effects idempotent.
- Every consumer has a DLQ named `<original>.dlq`. Retry first (typically 3 attempts, exponential backoff), then DLQ.
- Alert on any non-zero DLQ depth.
- Backpressure: bounded prefetch and concurrency; lag-based scaling; `503` with `Retry-After` when persistently overloaded.
- Sagas start as choreography; promote to orchestration at 5+ services, branching, retries or compensation.
- Propagate a correlation id through every event in a saga.
- Pick one broker per system.

## Hard constraints
- The outbox is mandatory for any write that emits an event. Publish-before-commit and publish-after-commit are broken under failure.
- Consumers must be idempotent; delivery is at-least-once.
- Additive schema changes stay in the version. Breaking changes go to a new `vN` topic, run side by side.
- Don't chase protocol-level exactly-once.
- Don't design for global ordering; it is a smell.

Invoke when: designing event flows, choosing event vs command vs RPC, naming topics or schemas, designing outbox, idempotency or DLQ handling, or auditing consistency.

Full guide: invoke the event-driven-architect skill.

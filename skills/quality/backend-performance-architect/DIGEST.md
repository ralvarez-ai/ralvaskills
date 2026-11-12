# backend-performance-architect digest

Impact-ordered backend latency and throughput playbook: remove work first, then query shape, then hops, then tail control, then transport.

## Always-on rules

- Order by impact: remove work (cache, precompute, off the hot path) → query and index shape → fewer or parallel hops → tail controls → serialization/transport/framework → micro-optimization.
- Measure first: p50/p95/p99 per endpoint, `EXPLAIN (ANALYZE, BUFFERS)`, `pg_stat_statements`, pprof, a trace of the slow request. No profile, no change.
- Stop a rung when the profile stops showing that layer. Never change the transport while a query dominates the endpoint.
- Cache in layers — in-process (bounded + TTL), shared (Valkey/Redis), CDN/HTTP — with an explicit invalidation story. Version keys; never cache an authorization decision without a revocation path.
- Coalesce concurrent misses (`singleflight`), jitter TTLs so a fleet does not expire together, and cache negatives briefly.
- Move email, push, audit, exports and webhooks behind a queue or outbox; return `202` and let a worker do the work.
- Index every FK not already covered by a composite index, and every `WHERE`/`JOIN` column; composite most-selective-first; partial index for soft delete; read the plan before adding one.
- Kill N+1 with batching/`IN`/`JOIN`; `COPY`/`UNNEST` for bulk writes; cursor pagination on hot lists, never `OFFSET`.
- Size connection pools explicitly and cap the total across replicas. Add a read replica only with explicit routing and bounded staleness; shard only after a single-node measurement.
- One pooled HTTP client or gRPC/Connect channel per target. Parallelize independent hops with bounded concurrency; prefer a batch endpoint over N calls.
- Every outbound call carries a deadline shorter than its caller's. Retry only idempotent calls, with a budget, backoff and jitter; break the circuit; hedge only after measurement.
- Bound queues and goroutines; shed low-priority load early (`429`/`503` + `Retry-After`); degrade with stale cache or a partial response instead of failing.
- Transport is a semantics choice, not a speed one: unary Connect/gRPC for internal RPC, REST/Connect for browsers, server-stream for progress, bidi only when both sides push — bidi does not reduce round trips.
- `fasthttp` only on an isolated listener for a measured, middleware-free hot endpoint; it loses `net/http` context cancellation, HTTP/2, streaming and tooling.
- Go: pprof on a private listener, `GOMEMLIMIT` before allocation work, explicit pool limits, request context propagated, benchmarks with `-benchmem`.

## Hard constraints

- No measurement, no optimization.
- Never shard, add a replica or add a queue before a single-node measurement says so.
- Never cache a per-user response under a shared key, nor an authorization decision without revocation.
- Never migrate the main API to `fasthttp`.

Invoke the full skill when: a service is slow, an endpoint misses its p99, a load test is coming, or someone proposes a faster framework, server or transport instead of a fix.

Full guide: invoke the backend-performance-architect skill.

# Performance Reviewer (digest)

Cross-language performance review: N+1, missing indexes, blocking async I/O, allocation hot paths, unbounded memory, slow algorithms. Measurement-grounded.

## Always-on rules
- Measure first. Every finding cites a measurement (`EXPLAIN ANALYZE`, pprof, py-spy, histogram, hyperfine). "Looks slow" is not a finding.
- Read the diff with the measurement in hand. Propose fixes that produce a new measurement.
- Watch for N+1: a loop that fires one query per item. A handler emitting more than 3 queries per request is suspect.
- Index every FK and every `WHERE` column. OFFSET pagination is O(N); use cursor pagination.
- Inside `async def`: no `time.sleep`, `requests.get`, or sync DB driver. Use `await asyncio.sleep`, `httpx.AsyncClient`, psycopg 3 async. Use `asyncio.to_thread` only when a sync API is unavoidable. Prefer `asyncio.TaskGroup` over bare `gather`.
- Go: bound goroutine spawns (worker pool or `errgroup` `SetLimit`). Propagate request context; no `context.Background()` in handlers.
- Every outbound HTTP/DB/queue/RPC call has a deadline. Retries use jittered exponential backoff. Reuse HTTP clients. Single-flight hot keys.
- Keep metric labels low-cardinality. No logging inside tight loops; trace request boundaries, not every inner call.
- Compile regexes once. Build strings with `join`/`strings.Builder`, not `+` in loops. Use set/map for lookups inside nested loops.

## Hard constraints
- No finding without a measurement. Gut-feel optimizations are rejected; if the code is not measurably slow, do not change it.
- Load testing and capacity planning are out of scope for code review.

Invoke the full skill when: a PR touches a hot path, an endpoint shows in p99 or error-rate dashboards, before a load test, after a metric regression, or new code loops over awaits or DB calls.

Full guide: invoke the performance-reviewer skill.

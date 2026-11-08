# Backend Performance — Reference Tables and Recipes

Loaded on demand. Measurement commands, SQL, Go snippets, and the decision tables behind §7 of [SKILL.md](SKILL.md).

## Contents

1. [Measurement recipes](#1-measurement-recipes)
2. [Postgres patterns](#2-postgres-patterns)
3. [Caching patterns](#3-caching-patterns)
4. [Tail-latency patterns](#4-tail-latency-patterns)
5. [Go transport choices](#5-go-transport-choices)
6. [Pre-flight checklist](#6-pre-flight-checklist)

## 1. Measurement recipes

**Postgres — the plan, not the guess.**

```sql
EXPLAIN (ANALYZE, BUFFERS, VERBOSE)
SELECT id, created_at, total
FROM movements
WHERE profile_id = $1 AND created_at >= $2
ORDER BY created_at DESC
LIMIT 50;
```

Read in this order: `Seq Scan` on a large table → missing index; `Sort Method: external merge` → the sort does not fit `work_mem`; `actual rows` off by orders of magnitude from `plan rows` → stale statistics (`ANALYZE`); `Buffers: shared read` high → cache misses at the storage layer.

**Postgres — where the time goes fleet-wide.**

```sql
SELECT queryid, calls,
       round(total_exec_time::numeric / calls, 3) AS mean_ms,
       round((total_exec_time::numeric / sum(total_exec_time) OVER ()) * 100, 1) AS pct
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 20;
```

If the top query is not the endpoint you are investigating, keep looking — the endpoint is probably calling it N times.

**Go — profiles.**

```bash
curl -s -o cpu.pprof 'http://127.0.0.1:6060/debug/pprof/profile?seconds=30'
go tool pprof -http=:8081 cpu.pprof

curl -s -o heap.pprof 'http://127.0.0.1:6060/debug/pprof/heap?gc=1'
go tool pprof -sample_index=alloc_space -http=:8081 heap.pprof
```

CPU profile under realistic load, not at idle. `alloc_space` finds allocation hot paths; `inuse_space` finds caches and leaks. Mutex/block profiles answer "where are requests waiting", which is usually the p99 question.

**Go — benchmarks that survive review.**

```bash
go test -run='^$' -bench=. -benchmem -count=6 ./... > new.txt
benchstat old.txt new.txt
```

`-count=6` + `benchstat` because a single run cannot distinguish a 5% change from noise. `-benchmem` because allocations are the cost that compounds.

**Load — shape over volume.**

```bash
k6 run --vus 64 --duration 3m script.js   # mixed reads/writes
hey -z 30s -c 64 http://127.0.0.1:8080/v1/movements
```

Report p50/p95/p99, error rate and saturation (pool in use, queue depth, CPU throttling), not just requests per second.

## 2. Postgres patterns

**Indexes.**

```sql
-- Every FK, unless an existing composite index already covers it.
CREATE INDEX CONCURRENTLY ix_movements_profile_id ON movements (profile_id);

-- The list endpoint's shape: filter + sort, most selective column first.
CREATE INDEX CONCURRENTLY ix_movements_profile_created
    ON movements (profile_id, created_at DESC);

-- Soft delete: index only the live rows.
CREATE INDEX CONCURRENTLY ix_movements_live
    ON movements (profile_id, created_at DESC)
    WHERE deleted_at IS NULL;
```

**Batch reads and writes.**

```go
rows := pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
    return []any{items[i].ID, items[i].Name}, nil
})
_, err := pool.CopyFrom(ctx, pgx.Identifier{"things"}, []string{"id", "name"}, rows)
```

`COPY` for large inserts, `UNNEST($1::uuid[], $2::text[])` for medium batches, multi-row `INSERT` for small ones. One round trip, one transaction.

**Pool sizing.**

```go
cfg, err := pgxpool.ParseConfig(dsn)
if err != nil {
    return nil, err
}
cfg.MaxConns = int32(4 * runtime.GOMAXPROCS(0)) // start here, then measure
cfg.MinConns = int32(runtime.GOMAXPROCS(0))
cfg.MaxConnLifetime = time.Hour
cfg.MaxConnIdleTime = 15 * time.Minute
```

The database has a connection ceiling; multiply the per-replica pool by the replica count and stay well under it. Queueing in the app (short waits) is better than queueing in the database (workers thrashing). A pooler earns its place when app replicas × pool size exceeds what Postgres can serve.

**Cursor pagination.**

```sql
SELECT id, created_at, total
FROM movements
WHERE profile_id = $1
  AND (created_at, id) < ($2, $3)   -- last row of the previous page
ORDER BY created_at DESC, id DESC
LIMIT $4;
```

The row comparison uses the composite index and stays O(log n) at any depth, unlike `OFFSET`.

## 3. Caching patterns

**Cache-aside with coalescing and jitter.**

```go
func (s *Store) Get(ctx context.Context, id string) (*Thing, error) {
    if v, ok := s.cache.get(id); ok {
        return v, nil
    }

    v, err, _ := s.group.Do(id, func() (any, error) {
        t, err := s.db.Load(ctx, id)
        if err != nil {
            return nil, err
        }
        s.cache.set(id, t, jitter(30*time.Second, 0.2))
        return t, nil
    })
    if err != nil {
        return nil, err
    }

    return v.(*Thing), nil
}

func jitter(base time.Duration, spread float64) time.Duration {
    return base + time.Duration(rand.Float64()*spread*float64(base))
}
```

`singleflight.Group.Do` makes concurrent misses share one load; the shared call uses the *first* caller's context, so keep the load short and do not derive it from a single request's deadline. Jitter stops a fleet from expiring the same key at the same instant.

**HTTP caching matrix.**

| Response | Header |
|---|---|
| Public, rarely changes | `Cache-Control: public, max-age=3600, stale-while-revalidate=60` |
| Public, must revalidate | `Cache-Control: no-cache` + `ETag` (answer `304` on match) |
| Per-user | `Cache-Control: private, no-store` — do not let a CDN hold it |
| Immutable asset | `Cache-Control: public, max-age=31536000, immutable` |

**What not to cache:** authorization decisions without revocation, anything keyed by a shared key that includes a user id, and writes that must be read-after-write consistent.

## 4. Tail-latency patterns

**Timeouts as a budget.** Give each hop less than the caller's remaining deadline; a child that outlives its parent is wasted work.

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()
```

**Retries with full jitter and a cap.**

```go
func backoff(attempt int, base, max time.Duration) time.Duration {
    d := min(base<<attempt, max)
    return time.Duration(rand.Int64N(int64(d))) // full jitter
}
```

Retry only idempotent calls, only on transient errors (`Unavailable`, `DeadlineExceeded`, `5xx`, connection reset), and only within a budget. Every retry multiplies load; without a budget it turns a brownout into an outage.

**Admission control that fails fast.**

```go
func Limit(max int64) func(http.Handler) http.Handler {
    var inFlight atomic.Int64
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if inFlight.Add(1) > max {
                inFlight.Add(-1)
                w.Header().Set("Retry-After", "1")
                http.Error(w, "overloaded", http.StatusServiceUnavailable)
                return
            }
            defer inFlight.Add(-1)
            next.ServeHTTP(w, r)
        })
    }
}
```

Pair it with a queue-depth or pool-saturation signal so the limit reflects pressure, not just a fixed number. Shed low-priority traffic first (analytics, prefetch, recomputation) and keep the user-facing path alive.

## 5. Go transport choices

**Reuse one channel per target.**

```go
conn, err := grpc.NewClient(addr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(8<<20)),
)
if err != nil {
    return nil, err
}
defer conn.Close() // process lifetime, not per request
```

Per-request `grpc.Dial` is the gRPC version of "open a connection per query". For Connect, one `http.Client` per target, with `Timeout` unset and deadlines carried by the context.

**HTTP/2 without TLS (`h2c`) on a private network.**

```go
// Server: Go 1.24+ native switches.
var p http.Protocols
p.SetHTTP1(true)
p.SetUnencryptedHTTP2(true)
srv.Protocols = &p

// Or with golang.org/x/net/http2/h2c, for older code:
srv.Handler = h2c.NewHandler(mux, &http2.Server{})
```

```go
// Client: h2c needs an HTTP/2 transport that dials plaintext.
tr := &http2.Transport{
    AllowHTTP: true,
    DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
        return (&net.Dialer{}).DialContext(ctx, network, addr)
    },
}
client := &http.Client{Transport: tr}
```

Without this, Connect clients on plaintext HTTP speak HTTP/1.1: keep-alive and pooling, but one in-flight request per connection. Enabling `h2c` is usually a bigger win than changing RPC libraries.

**Streaming (grpc-go).**

```go
func (s *server) Sync(stream pb.Sync_SyncServer) error {
    for {
        req, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            return nil
        }
        if err != nil {
            return err
        }
        if err := stream.Send(s.answer(req)); err != nil {
            return err
        }
    }
}
```

Bidi requires HTTP/2 end to end. Under a proxy or tunnel, confirm it speaks `h2` to the backend, or the stream will fail or silently degrade. Bidi is for protocols where both sides push; it does not batch unary calls.

**`fasthttp` as an isolated listener.**

```go
srv := &fasthttp.Server{
    Handler: func(ctx *fasthttp.RequestCtx) {
        ctx.SetContentType("application/json")
        ctx.SetBodyString(precomputed)
    },
}

go func() {
    _ = srv.ListenAndServe("127.0.0.1:8085") // private listener, not the public edge
}()
```

Use it only after the baseline endpoint is measurably bound by the HTTP layer, and only for a handler that needs no shared middleware. What you give up: `net/http` handlers and middleware, request-context cancellation on client disconnect, HTTP/2 and gRPC/Connect, `httptest`, and the tracing/metrics integrations. Measure both:

```bash
hey -z 30s -c 64 http://127.0.0.1:8080/hot > baseline.txt
hey -z 30s -c 64 http://127.0.0.1:8085/hot > fasthttp.txt
```

Compare p50/p99, throughput and CPU. If the baseline is already waiting on the database, the two numbers will match and the migration bought nothing.

## 6. Pre-flight checklist

Before changing code, answer these. Any "no" sends you back up a rung.

1. Which endpoint, measured how (p50/p99, load shape)?
2. Where does the trace say the time goes — database, a hop, CPU, queueing?
3. Is the offending query's plan known, and does an index explain it?
4. Would a cache remove a query, and what is its invalidation story?
5. Can the work leave the request (queue, outbox, precompute)?
6. Are independent hops already parallel, and are their clients pooled?
7. Does every outbound call have a deadline, and does every retry have a budget?
8. What happens at 3x traffic — bounded queues or OOM?
9. Is the proposed change at the rung the profile points to?
10. What measurement proves it worked, and what will you revert if it does not?

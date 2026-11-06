# ConnectRPC Architect (digest)

Go ConnectRPC: one `.proto` served as Connect, gRPC and gRPC-Web from one `http.Handler`, typed errors, interceptor chain, streaming rules, conformance tests.

## Always-on rules
- `.proto` is the RPC contract. Connect is the primary adapter for new Go services.
- Keep a native gRPC adapter from the same proto as insurance.
- REST is a separate public surface, not generated from proto. No `grpc-gateway`.
- Generate with buf via `task generate`, `paths=source_relative`. Never hand-edit `*connect.go`.
- Mount the path returned by `NewUserServiceHandler` on a `http.ServeMux`. The generated shim holds no business logic.
- Return `connect.NewError(code, err)`. Map domain errors to codes in one central place.
- Use `connect.CodeOf(err)` in logs and tests. Never invent an error-to-int map.
- Attach detail with `connect.NewErrorDetail`; include a request id.
- Interceptor order, outermost first: recovery, request-id, logging, auth, validation, metrics.
- Auth verifies tokens locally against cached JWKS; no per-request call to an auth service.
- Validation via `validate.NewInterceptor()` (protovalidate). No hand-written checks per handler.
- Every call and every stream has a deadline.
- Correlation id header: `X-Request-Id`.
- Health always on. Reflection dev only.
- Browsers: gRPC-Web is server-stream only. Make uploads unary.
- Server-streams take resume tokens (`start_after_id` or a cursor).
- Tests: `httptest` for Connect; one conformance table run across transports.

## Hard constraints
- Never leak stack traces or DB errors. Log server-side; return a generic message plus the correlation id.
- Never use `CodeInternal` as a default. Map every known domain error to a specific code.
- Never `Allow-Origin: *` on authenticated endpoints.
- Never enable reflection in prod; it leaks the service surface.
- Never leave a stream unbounded; an unbounded stream is a leak.
- Client- and bidi-streaming need HTTP/2.

Invoke when: designing or reviewing a Connect handler, choosing Connect vs vanilla gRPC vs REST, serving proto to browser callers, picking a streaming pattern, or writing conformance tests.

Full guide: invoke the connectrpc-architect skill.

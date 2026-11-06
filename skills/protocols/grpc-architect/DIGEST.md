# gRPC Architecture (digest)

Vanilla gRPC for backend-to-backend services: .proto services, status codes, interceptor chain, deadlines, bufconn tests.

## Always-on rules
- One service per file, named after the resource; methods verb-noun.
- Every method gets typed `<Verb><Noun>Request` / `<Verb><Noun>Response`; prefer `CreateUserResponse { User user = 1; }` over bare `User`.
- `google.protobuf.Empty` only as a response for fire-and-forget actions.
- Pagination by cursor: `ListUsersRequest { cursor, limit }` -> `ListUsersResponse { users, next_cursor }`.
- Return errors via `status.Error(code, msg)`; map domain errors to codes in one central place.
- Pick the specific standard code for every known domain error.
- `status.WithDetails` (ErrorInfo, BadRequest) for machine-readable context; include a correlation id.
- Interceptor chain, outermost first: recovery -> request-id -> log -> auth -> validation -> metrics. Same chain for streams.
- Validation centralized via protovalidate in an interceptor, returns INVALID_ARGUMENT.
- Generated handlers are thin shims; logic lives in service-layer code.
- Every call has a deadline. Server respects it, checks `ctx.Err()` in long work, propagates ctx downstream.
- Unary deadlines 5-30s; server-streams long but bounded.
- Handler wrapped in `context.WithTimeout` slightly under the client deadline.
- Correlation id in one header (e.g. `x-request-id`), read by interceptor, put in ctx, logged.
- Health service (`grpc.health.v1`) always on.
- Server-streams carry a resume cursor; bidi protocol sketched in .proto comments.
- Test with bufconn; real port only for TLS / HTTP/2 network-path checks.
- Table-driven tests: one row per (input, expected code, expected error type).
- grpcurl for manual testing, pinned in task runner or mise.
- gRPC for east-west backend traffic; REST for browser, public, cache-friendly, or small-scope APIs.

## Hard constraints
- Never send raw primitives or `google.protobuf.Empty` as request input.
- Never leak stack traces or DB errors to clients; log server-side, return generic INTERNAL with correlation id.
- Don't reach for INTERNAL as a default.
- Don't send large payloads in metadata; business data goes in message fields.
- Always set stream deadlines; an unbounded stream is a leak.
- Reflection on in dev, off in production.
- Don't use streaming to "save round-trips"; start with unary plus pagination.

Invoke the full skill when: designing or reviewing a gRPC service, choosing stream types, wiring interceptors, or mapping domain errors to status codes.

Full guide: invoke the grpc-architect skill.

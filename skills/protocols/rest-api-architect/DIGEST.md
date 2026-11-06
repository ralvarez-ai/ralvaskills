# REST API conventions (digest)

Framework-agnostic JSON-over-HTTP REST conventions: URLs, methods, status codes, versioning, pagination, errors, idempotency, concurrency, caching, auth headers, OpenAPI.

## Always-on rules
- Plural nouns, no verbs in URLs; resource ID in the path; nest one level at most; `snake_case` in URLs, query params and JSON.
- State changes are created sub-resources (`PUT /v1/orders/{id}/cancellation`), not `POST /v1/orders/{id}/cancel`.
- `PUT` takes the full representation; `PATCH` is JSON merge patch (RFC 7396).
- Use the status code table in the skill. `401` means missing or invalid credentials; `403` means valid credentials without permission.
- URL-prefix versioning only (`/v1/`). Bump on breaking changes only; send `Deprecation` and `Sunset` before removal.
- Cursor pagination, not offset. Filter `?status=`, sort `?sort=-created_at,name`, search `?q=`, sparse fields `?fields=`.
- Errors are RFC 9457 `application/problem+json`; `type` is a stable URL; `correlation_id` on every `5xx`.
- JSON: ISO 8601 timestamps with timezone as strings; money as decimal strings; UUID v7; `null` is intentional absence, a missing field is not provided.
- Auth via `Authorization: Bearer`. No tokens in query params; no PII in URLs.
- `Cache-Control` on every response; `Vary` when responses differ by user or encoding.
- Rate limits: `429` with `RateLimit-*` headers and `Retry-After`; limit per caller, not per IP.
- OpenAPI 3.1 is generated from code, snapshot-tested and linted in CI.

## Hard constraints
- Every `POST` and `PATCH` requires an `Idempotency-Key` header; without it the server returns `400 Bad Request`.
- Every editable resource exposes an `ETag`; every `PUT` and `PATCH` requires `If-Match` matching it, or returns `412 Precondition Failed`.
- Never `200 OK` for errors.
- Never change a published error `type` URL.
- Never `Access-Control-Allow-Origin: *` on authenticated endpoints.
- No PII in URLs. Never use Unix epoch numbers for timestamps.

Invoke the full skill when: designing or auditing REST endpoints; writing error contracts; handling idempotency or ETag/If-Match; choosing auth, pagination or caching behaviour; writing OpenAPI.

Full guide: invoke the rest-api-architect skill.

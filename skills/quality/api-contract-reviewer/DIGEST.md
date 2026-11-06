# api-contract-reviewer digest

Reviews REST OpenAPI specs and gRPC/protobuf contracts for stability, versioning and completeness. Output is a severity-keyed findings table.

## Always-on rules
- Passes in order: mechanical (`openapi-diff` for OpenAPI, `buf breaking` for proto), convention, completeness.
- Findings table columns: Severity | Rule | Location | Evidence | Fix.
- Severity: Critical = wire-breaking change in a stable version; High = expensive convention violation; Medium = completeness gap; Low = polish.
- Versioning: REST under `/vN/`; proto version in the package path, never in message names; additive changes stay in the current version; breaking changes bump the major version.
- Deprecate before remove: `deprecated` flag plus `Deprecation`/`Sunset` headers for REST.
- Proto: removed field numbers and names go to `reserved`; field types never change; enums have `*_UNSPECIFIED = 0`; no primitive wrappers (use `optional`).
- Errors: REST returns `application/problem+json` (RFC 9457); gRPC uses `status.Error` with standard codes mapped centrally; `type` URLs are stable; 5xx carry `correlation_id`; 422 lists fields.
- Idempotency: `Idempotency-Key` on POST/PATCH; `ETag` + `If-Match` on PUT/PATCH; cursor pagination, not offset.
- JSON: `snake_case`; ISO 8601 timestamps with timezone; money as strings; UUIDs canonical with dashes; `null` and missing are distinct in PATCH.
- OpenAPI: tag, summary and description on every operation; `examples:` on models; `requestBody.required: true` when mandatory; `info.version` matches the major version.
- gRPC: one service per file; `<Verb><Noun>Request` / `Response`; never `google.protobuf.Empty` as input.
- Specs are generated from code and snapshot-tested in CI; generated gRPC code is committed under `gen/`.

## Hard constraints
- Never reuse a field number; never change a field's type, even if values still fit.
- Breaking changes always bump the major version; run the new `vN` side by side.
- Never mix `{"detail": ...}` and `{"errors": [...]}` error shapes in one API.
- Never send `google.protobuf.Empty` as RPC input.
- Out of scope here: performance, security and architecture review; hand those off.

Invoke the full skill when: a PR changes `.proto`, OpenAPI or a wire-affecting handler; before publishing a new major version, deprecating a field or method, or a breaking-change release; or for a periodic spec-drift audit.

Full guide: invoke the api-contract-reviewer skill.

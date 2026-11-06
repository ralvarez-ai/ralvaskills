# observability-architect digest

Application-side signals: structured logs, Prometheus metrics, OTel traces, and how they correlate.

## Always-on rules

- Pillars: logs = what happened on this request; metrics = rate/aggregate; traces = the path.
- Correlate: `trace_id` in every log line; exemplars on histograms; `service.name` and `service.version` on traces, metrics and logs.
- Metric names: `<namespace>_<subsystem>_<name>_<unit>_<type>`. Unit always present (`_seconds`, `_bytes`, `_total`, `_ratio`); never `_ms`. Counters end `_total`.
- Labels bounded: `method`, templated `route`, `status_class`, `tenant`, `region`. Soft cap under 100 values per label.
- Logs: structured JSON to stderr; Go `log/slog`, Python `structlog`. Fixed `msg` string, values as fields.
- Mandatory log fields: `timestamp` (RFC 3339 UTC), `level`, `msg`, `trace_id`/`span_id` when traced, `service`.
- Errors: failed operation and relevant inputs in `msg`; stack trace in separate `error.stack`.
- Traces: auto-instrument first; span names verb + templated resource; OTel semconv attribute keys; set error status on failures; trace boundaries and sub-operations, not hot loops.
- Sampling: head 10% in prod, deterministic per `trace_id`; 100% of error traces; 100% in non-prod. Tail sampling deferred until needed.
- SLOs: RED for request-driven services, USE for resources; one SLO per critical user journey.

## Hard constraints

- PII and secrets never appear in any signal: no emails, names, phones, passwords, tokens, API keys, session cookies. Use `email_domain`.
- Never log full payloads; log size and shape.
- No unbounded labels: `user_id`, `email`, `request_id`, `correlation_id`.
- Raw `auth_token_id` is not safe to emit; UUID v7 IDs are fine in metric labels.
- Don't write more than 3-5 SLOs per service.

Invoke the full skill when: instrumenting code, naming metrics, wiring Go or Python stacks, or auditing what a service emits.

Full guide: invoke the observability-architect skill.

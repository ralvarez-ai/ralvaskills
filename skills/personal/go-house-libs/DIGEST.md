# Go House Libs (digest)

Maps a need to the ralvarezdev shared Go modules. Use them instead of re-writing their code.

## Always-on rules
- Pick the module from the skill's table and read `references/<NAME>.md` before wiring.
- Dependency direction only inward: `ginkit` -> restkit, ratelimit; `pgkit` -> restkit; `identity` -> email, pgkit, ratelimit; `identitygin` -> identity; `webpush` -> pgkit; `grpckit`, `connectrpckit` -> corekit. `sqlitekit`, `svckit`, `tick`, `resilience`, `corekit`, `contracts` import no other kit.
- Each adapter owns its third-party dependency in its own package; a consumer that skips it never pulls it in.
- A different HTTP framework gets its own `<framework>kit`, not a branch in restkit.
- A project keeps its one wire contract. New projects start on RFC 9457 problem+json.
- Promote code into a module only when a second project needs it; `mcpkit` is the exception (extracted up front). Adding a module = tag, migrate a consumer in one commit, update `references/` and `STACK.md`.
- Before calling a release additive, diff against the previous tag and check STACK.md pins.
- A change that turns a nil-returning call into an error is breaking; say so in the commit and docs.

## Hard constraints
- Do not hand-write lifecycle (use `svckit`), ticker loops (`tick.Every`), retry loops (`resilience/retry`) or SQLite open (`sqlitekit`).
- Do not re-declare a kit type locally: no second `Response[T]`, cursor codec, pgtype helper or mailer interface.
- Never pass `err.Error()` to `ginkit.Err` for a 5xx. Log it, answer generically; raw text only in debug mode.
- Agent-facing tool errors go through `mcpkit.ErrorMsg`.
- Idempotency fails closed: store outage refuses the request and never runs the handler; release the key on 5xx or panic.
- Generated code is types only. Hand-write calls on `restkit/client.Do`.
- Infrastructure is its own module: rate limiting is not login throttling; email and push are not identity concerns. Wire them, do not fold them in. Never a kit importing a consumer.

Invoke the full skill when: building or extending a Go API, MCP server, gRPC or Connect service, CLI, or Postgres/SQLite store, before hand-writing envelopes, cursors, contracts, Connect code, pgtype conversions, auth, rate limits, mail, push, lifecycle, tickers, retries or SQLite opens.

Full guide: invoke the go-house-libs skill.

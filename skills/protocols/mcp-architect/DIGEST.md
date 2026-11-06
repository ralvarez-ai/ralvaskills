# MCP Architect (digest)

Design and review MCP servers on spec 2026-07-28: tools, resources, prompts over JSON-RPC 2.0, stateless core.

## Always-on rules
- Primitive by who invokes: model -> tool; read-only data -> resource; user menu -> prompt.
- Tools: one verb, narrow scope, full JSON Schema with descriptions; description leads with the action, <300 tokens.
- Every result has `content[]`; with `outputSchema`, the server MUST populate `structuredContent` and SHOULD also emit a JSON text copy in `content[]`.
- Tool failures: `isError: true` inside the result. JSON-RPC errors are for protocol failures only.
- Set annotations explicitly: `title`, `readOnlyHint`, `destructiveHint` (err toward true), `idempotentHint`, `openWorldHint`. Omitted defaults are the most dangerous.
- Resources: `mimeType` on every resource; `ttlMs` + `cacheScope` on list and read results; deterministic list order.
- Stateless: no `initialize`, no `Mcp-Session-Id`. Client identity and capabilities ride in `_meta` on every request.
- Headers: `MCP-Protocol-Version`, `Mcp-Method`, `Mcp-Name` on Streamable HTTP POST.
- Mid-call input: MRTR (`InputRequiredResult`, `inputRequests`, `requestState`), not server push.
- Cross-call state: explicit-handle pattern (tool returns an ID the model passes back).
- OAuth 2.1 resource server: Authorization Code + PKCE, RFC 8707 `resource`, validate `aud` and RFC 9207 `iss`; CIMD over DCR.
- stdio: no auth; stderr for logs.
- MCP Apps panels only add to a complete text result; every capability stays reachable by a model-visible tool.
- Test: MCP Inspector, in-process transports, MRTR round trips, negative tests for URL, path and injection inputs.

## Hard constraints
- Never leak stack traces, DB errors or internal hostnames in either error layer.
- Never write debug prints to stdout on stdio.
- Don't put secrets or stack traces in `text` blocks.
- Don't break tool schemas in place; add a new tool and deprecate the old.
- Don't declare `outputSchema` for free-form prose tools.
- Validate on the server regardless of what the client claims.
- Never let tool descriptions be edited dynamically from untrusted input.
- The server MUST reject tokens whose audience claim doesn't match.
- Use a library for OAuth; don't hand-roll. Don't bolt OAuth onto stdio.
- Don't build new onboarding on DCR. Don't claim a spec version you haven't tested.
- Don't auto-chain destructive tools off untrusted text.
- SSRF: allowlist hosts when the use case permits; block RFC 1918, link-local, loopback and IMDS before the request.

Invoke the full skill when: designing, reviewing or scaffolding an MCP server, or authoring its Go JSON Schema contracts.

Full guide: invoke the mcp-architect skill.

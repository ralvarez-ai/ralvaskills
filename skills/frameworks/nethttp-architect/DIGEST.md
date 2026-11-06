# net/http Architecture (digest)

Go stdlib net/http on Go 1.26: enhanced ServeMux, no router framework, feature layout, RFC 9457 errors, JWT or IdP auth, graceful shutdown.

## Always-on rules
- Feature folder under `internal/` with handlers.go, service.go, repo.go, dto.go, queries/.
- handlers.go depends on service.go; never reaches into repo.go.
- service.go is pure Go, no net/http imports.
- dto.go holds wire types; never reuse domain structs as DTOs.
- Routes as `"GET /v1/users/{id}"`; read path values with `r.PathValue`.
- Be explicit on trailing slash: `"GET /v1/users/{$}"` matches exact path only.
- No router framework just for syntax sugar; write the prefix per route or use a small helper.
- Decode with `json.Decoder.DisallowUnknownFields()`; wrap body with `http.MaxBytesReader(w, r.Body, 1<<20)`.
- One shared `*validator.Validate` on the handler; validate with `validate.Struct`.
- Use one `writeJSON` helper; set headers before `WriteHeader`; `http.StatusNoContent` for empty bodies.
- Return DTOs, not domain types.
- Explicit constructors; wire deps in main.go; hand-wired for small services, `uber-go/fx` for larger graphs.
- Set `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` (mandatory in production); shutdown via `signal.NotifyContext`.
- Middleware is `func(http.Handler) http.Handler`; standard chain RequestID -> SLog -> Recover -> CORS -> GZip (outermost first): request id first so every response and log line carries it; logger outside Recover so panics are logged as 500s.
- Auth per route, never global; authorization via RequireScope composed with AuthRequired. With an external IdP, verify `aud` and `iss` explicitly.
- Errors as RFC 9457 problem+json; Recover emits 500 with correlation id, never a stack trace.
- Background work: `context.Background()` with timeout, log errors.
- Tests: `httptest.NewRecorder` with the handler directly; `testing/synctest` for timing-dependent goroutines.
- OpenAPI via kin-openapi, served at `/v1/openapi.json`, snapshot-tested in CI.

## Hard constraints
- Never leave the four server timeouts unset; defaults are unlimited.
- Never write your own logger middleware using `log.Printf`; use slog.
- Never reuse domain structs as DTOs.

Invoke the full skill when: scaffolding or reviewing a stdlib net/http service, wiring auth middleware, or setting server timeouts and shutdown.

Full guide: invoke the nethttp-architect skill.

# Gin Architecture (digest)

Gin 1.12 on Go 1.27.3: feature layout, struct-tag validation, RFC 9457 errors, JWT or external IdP auth, route-group versioning, OpenAPI.

## Always-on rules
- One folder per bounded context: handlers.go, service.go, repo.go, dto.go, SQL files.
- handlers.go depends on service.go; never reaches into repo.go.
- service.go is pure Go with no gin imports.
- dto.go holds wire types with `json:` and `binding:` tags.
- One `Register(rg, deps)` function per feature; keeps main.go thin.
- Version via `gin.RouterGroup` per version, sub-group per feature.
- Parse path params at the boundary; return 400 on parse failure.
- Bind with `c.ShouldBindJSON` / `ShouldBindUri` / `ShouldBindQuery`.
- `binding:"required"` on every non-optional field; `omitempty` on PATCH DTOs.
- Register custom validators once at startup.
- Return DTOs; `json:"-"` on fields that must never serialize.
- Explicit status codes; `c.Status(http.StatusNoContent)` for empty bodies.
- Explicit constructors; `Deps` struct wired in main.go; `uber-go/fx` for larger graphs.
- Open shared resources in main.go; close on `signal.NotifyContext` shutdown; `ReadHeaderTimeout` mandatory.
- Auth: in-house JWT (golang-jwt/jwt/v5 + argon2) or external IdP with JWKS cached in-process with TTL.
- Authorization per route via `RequireScope(...)` alongside `AuthRequired`.
- Errors as `application/problem+json` via one problem package; one mapping switch for sentinel errors.
- Middleware: `gin.New()` with RequestID, SLog, Recovery, CORS, gzip, in that order.
- Background work: derive context from `context.Background()`, set timeout, log errors.
- Tests: `httptest.NewRecorder` with the engine; rolled-back transaction per test.
- OpenAPI: swaggo/swag by default, kin-openapi for full 3.1 control; snapshot-test the spec in CI.

## Hard constraints
- Never `c.MustBindWith`; it panics.
- Never `gin.Default()` in production; use slog.
- Never global auth middleware; auth is per route.
- Never `c.JSON(204, nil)` for empty bodies.
- Never reuse domain structs as DTOs.
- Recovery before CORS.
- Serious work (retryable, distributed, scheduled) goes to a real task queue, not a goroutine.

Invoke the full skill when: scaffolding or reviewing a Gin service, wiring binding or validation, or setting up middleware order and problem responses.

Full guide: invoke the gin-architect skill.

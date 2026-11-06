# FastAPI Architecture (digest)

FastAPI 0.136 on Python 3.14: feature layout, Pydantic v2 request/response split, lifespan DI, URL-prefix versioning, RFC 9457 errors.

## Always-on rules
- One folder per bounded context: router.py, service.py, repo.py, schemas.py, SQL files.
- router.py depends on service.py; never reaches into repo.py.
- service.py is pure Python with no FastAPI imports.
- schemas.py holds Pydantic models, never reused as ORM models or DB rows.
- URL-prefix versioning: `/v1/...` via `APIRouter(prefix="/v1")`; deprecate by mounting `/v2` alongside.
- One APIRouter per feature; tags match feature folder names.
- Typed path params in the signature; `response_model` per route; explicit `status_code`.
- Three schemas per resource: Create, Update (PATCH partial), Response.
- `extra="forbid"` on every request model; `SecretStr` for passwords and tokens.
- Pydantic v2: `@field_validator` per field, `@model_validator(mode="after")` cross-field.
- Shared state via `Depends` (never imported as module-level globals): DB connections, HTTP clients, auth subjects. Async deps for I/O.
- Lifespan is the only place to open and close shared resources; settings validated once via pydantic-settings.
- Auth: OAuth2PasswordBearer + pyjwt + argon2-cffi in-house, or PyJWKClient for an external IdP; verify `aud` and `iss`.
- Authorization is route-level via `dependencies=[Depends(...)]`, not middleware.
- One exception handler per domain-exception family; validation errors mapped into Problem detail.
- Middleware order: CORS, GZip, Request ID, logging; auth stays a dependency.
- BackgroundTasks only for fire-and-forget work tied to one response.
- Tests: TestClient or httpx AsyncClient with ASGITransport; override deps via `app.dependency_overrides`.
- Snapshot the OpenAPI spec in CI; document non-default status codes via `responses={...}`.
- Routes are `async def` unless they call a sync library you decided not to wrap; run sync I/O with `asyncio.to_thread`.
- Open the DB pool once in lifespan; `response_model_exclude_unset=True` for large optional-heavy models.

## Hard constraints
- Never reuse one model for request and response.
- Never use `@app.on_event`; never open shared resources at module level.
- Never `time.sleep`, `requests`, or other blocking calls inside `async def`.
- Never `psycopg.connect()` per request.
- Never mutate `/v1`.
- Never let HTTPException and custom exceptions return different shapes.
- Never leak stack traces in `detail`.
- BackgroundTasks is not a queue; retryable or scheduled work goes to a real task queue.

Invoke the full skill when: scaffolding or reviewing a FastAPI service, designing schemas or error handlers, or wiring lifespan, auth, or middleware.

Full guide: invoke the fastapi-architect skill.

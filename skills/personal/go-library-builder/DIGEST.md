# Go Library Builder (digest)

Scaffolds reusable Go library modules in the ralvarezdev house style. Composition root over go-architect, hexagonal-architecture, ddd-architect.

## Always-on rules
- Check go-house-libs first. Scaffold a new module only when none fits and a second consumer exists.
- Shapes: flat single-port (domain root + one port, adapter subpackages); vertical-slice (package per bounded context, `app/`, `internal/acl/`); toolkit (no domain port, heavy deps in own subpackage); framework-adapter (wraps one framework around a framework-free core).
- Namespace is a parameter, default `github.com/ralvarezdev`; ask at bootstrap. Module path `<namespace>/<name>`, no `go-`/`ddd-`/`hex-` prefix.
- Layers are packages. Never `domain/`, `application/`, `ports/`, `adapters/` dirs.
- Adapter package = technology (`smtp`, `postgres`); adapter file named for its port; type `smtp.Mailer`, not `SmtpMailer`.
- Ports: declared by the consumer, ~3 methods max, grown by composition, `context.Context` first, no driver types in signatures.
- Accept interfaces, return `*T`. Constructor injection; `New(deps...)`.
- Config: plain validated struct with `setDefaults()` and `Config.Validate() error`. No functional options.
- Sentinels are package-prefixed (`user: not found`). Validation uses `errors.Join`; wrap `fmt.Errorf("<pkg>: <op>: %w")`.
- Business "no" is `Result{Allowed: false}` with nil error, not an error.
- Tests: in-memory adapter is the test double; `testcontainers-go` once per package in `TestMain`; `t.Parallel()`; reuse `pgkit/pgtest` and `ratelimit/valkeytest`.
- Persistence: sqlc (pgx/v5) default; goose via `pgkit/migrate` with library-owned `<name>_schema_version`.
- Bootstrap order: domain (stdlib + uuid) -> in-memory adapter -> production adapter -> persistence only if a DB adapter exists -> spec doc and README -> `task check` green -> `task tag VERSION=0.1.0` -> add the module to go-house-libs.

## Hard constraints
- Domain root imports stdlib + `google/uuid` only. One third-party import there is a bug.
- Only adapters import drivers, HTTP clients or crypto libraries. A driver error never escapes an adapter.
- The library never reads env, files or viper. No DI container (`fx`/`wire`) in the library root; an optional `fx.Module` may live in `<name>/fxmod`.
- No `init()`, no globals. No mock library, no generated mocks. ORMs are rejected.
- A port with one adapter is ceremony; the in-memory double counts as a second.
- A context never imports another's domain. Do not re-implement goose provider handling; forward to `pgkit/migrate`.
- Out of scope: deployable services and `cmd/` binaries.

Invoke the full skill when: bootstrapping a new Go library, adding an adapter, or auditing a module against the house style.

Full guide: invoke the go-library-builder skill.

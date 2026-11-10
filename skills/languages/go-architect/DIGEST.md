# go-architect (digest)

Go 1.27.3 standards: enums, struct layout, interfaces, concurrency, errors, DB access, packages, linting, go.work and app-repo layout.

## Always-on rules
- Enums: named type, `iota`, `String()`; `1 << iota` for bitmasks.
- Export sentinel errors for `errors.Is`; `errors.Join` to aggregate; wrap with `fmt.Errorf("op: %w")`.
- One block per kind for package-level `type`/`const`/`var`; order type → const → var → func.
- A literal used more than once becomes a named constant.
- Order struct fields largest → smallest; `New...` constructors; functional options for complex setup.
- Typed structs over `map[string]any`. Prefer stdlib: `cmp`, `slices`, `maps`, `iter`, `log/slog`, `net.JoinHostPort`.
- Define interfaces where used, 1–3 methods; accept interfaces, return concrete structs.
- Return `iter.Seq[T]` for large, short-circuitable or cursor-backed sequences; `[]T` for small bounded results.
- Embed `sync.RWMutex` directly above the fields it guards; `sync.WaitGroup.Go`; `time.Minute`, not `time.Second * 60`.
- `context.Context` first parameter for blocking or I/O; `signal.NotifyContext` for shutdown.
- Tests: table-driven, `package x_test`, `t.Parallel()`, `t.Helper()` first in helpers, `t.TempDir()`, `testing/synctest` for async code.
- DB: sqlc by default, generated code committed, CI fails on drift. sqlx + `//go:embed` only where sqlc cannot express the query. Don't mix the two in one package.
- DI through constructors; container only at the composition root; libraries expose `New(deps...)`.
- Lint: golangci-lint v2 (the skill's golangci.yml template); `gofmt` + `goimports`; `go fix` periodically.
- Dead code: `deadcode -test ./...`; confirm each hit before deleting.
- Packages: short, lowercase, named for what they provide; `shared`/`common`/`types`/`model` are catch-alls. One tiny top-level vocabulary package (`money`, `permission`, or the domain) for cross-context values and pure rules; `pkg/*` only when an external module imports it, otherwise `internal/`.
- Per binary: thin `cmd/` main; composition root in `internal/<service>`.
- No root layer bucket (`internal/app`, `application`, `services`, `core`, `domain`, `shared`): name each package for what it does; `internal/` is the privacy boundary, not a layer.
- Adapters: `<ctx>/postgres/` for the store, `<ctx>/client/<protocol>/` for outbound (`client/rest`, `client/grpc`), `<ctx>/<protocol>/` for inbound handlers (`mcp/`); no `remote/`, `transport/` or `handler/` folders, and servers stay app-level (`internal/api`, `internal/grpcserver`).

## Hard constraints
- Never embed `context.Context` in a struct.
- Never force a type assertion on an error; use `ok` or `errors.As`.
- Never use panic as control flow; panic only at init.
- Ports live in the context root package; a `port/` folder only survives a cycle or a cross-context contract, and is then named for the contract (`coreapi`). `adapter`/`adapters` are never folder names.
- golangci `exclusions.rules[].path` is a regex: `(a|b)` matches, `{a,b}` silently matches nothing, and a stale path after a rename disables the rule.
- Never name a package `util`, `common`, `helpers`, `shared`, `types` or `model`. No `init()`, no global state.
- Never name another context's exported type in a port signature; use the shared vocabulary or a consumer-local interface.
- Never `fmt.Sprintf("%s:%d", ...)` for host:port.
- Never hand-edit generated sqlc output.
- Never build SQL by string concatenation; always bind parameters.
- No anonymous interfaces: name every capability interface.
- Libraries never import fx or wire.
- Never use `use .` as the only `go.work` entry.

Invoke the full skill when: writing, reviewing or scaffolding Go code, or deciding where a Go module and its SPA live in a repo.

Full guide: invoke the go-architect skill.

# go-architect (digest)

Go 1.27.3: enums, struct layout, interfaces, concurrency, errors, DB access, packages, lint, go.work, app-repo layout.

## Always-on rules
- Enums: named type, `iota`, `String()`; `1 << iota` for bitmasks.
- Export sentinel errors for `errors.Is`; `errors.Join` to aggregate; wrap with `fmt.Errorf("op: %w")`.
- One block per kind at package level; order type → const → var → func.
- A literal used twice becomes a named constant.
- Order struct fields largest → smallest; `New...` constructors; functional options for complex setup.
- Typed structs over `map[string]any`; prefer stdlib (`cmp`, `slices`, `maps`, `iter`, `log/slog`); host:port via `net.JoinHostPort`, never `fmt.Sprintf`.
- Define interfaces where used, 1–3 methods; accept interfaces, return concrete types.
- `iter.Seq[T]` for large, short-circuitable or cursor-backed sequences; `[]T` for small bounded results.
- Embed `sync.RWMutex` directly above the fields it guards; `sync.WaitGroup.Go`; `time.Minute`, not `time.Second * 60`.
- `context.Context` first parameter for blocking or I/O; `signal.NotifyContext` for shutdown.
- Tests: table-driven, `package x_test`, `t.Parallel()`, `t.Helper()`, `t.TempDir()`, `testing/synctest` for async.
- DB: sqlc by default, generated code committed and drift-checked; never hand-edit it; sqlx + `//go:embed` only where sqlc cannot express the query; never mix the two in one package.
- DI through constructors; container only at the composition root; libraries expose `New(deps...)` and never import fx or wire.
- Lint: golangci-lint v2 (the skill's template); `gofmt` + `goimports`; `go fix` periodically.
- Dead code: `deadcode -test ./...`; confirm each hit before deleting.
- Packages: short, lowercase, named for what they provide; `shared`/`common`/`types`/`model` are catch-alls. One small top-level vocabulary package for cross-context values and pure rules; `pkg/*` only when an external module imports it, otherwise `internal/`.
- Per binary: thin `cmd/` main; composition root in `internal/<service>`.
- No root layer bucket (`internal/app`, `application`, `services`, `core`, `domain`, `shared`); name each package for what it does; `internal/` is a privacy boundary, not a layer.
- Adapters: `<ctx>/postgres/`, `<ctx>/client/<protocol>/` outbound, `<ctx>/<protocol>/` inbound; no `remote/`, `transport/` or `handler/` folders; servers stay app-level (`internal/api`, `internal/grpcserver`).

## Hard constraints
- Never embed `context.Context` in a struct.
- Never force a type assertion on an error; use `ok` or `errors.As`.
- Never use panic as control flow; panic only at init.
- Ports live in the context root; a `port/` folder only for a cycle or a cross-context contract, then named for the contract (`coreapi`). `adapter`/`adapters` are never folder names.
- golangci `exclusions.rules[].path` is a regex: `(a|b)` matches, `{a,b}` matches nothing, a stale path disables the rule.
- Never name a package `util`, `common`, `helpers`, `shared`, `types` or `model`. No `init()`, no global state.
- Never name another context's exported type in a port signature; use the shared vocabulary or a consumer-local interface.
- Never build SQL by string concatenation; always bind parameters.
- No anonymous interfaces: name every capability interface.
- Never use `use .` as the only `go.work` entry.

Invoke the full skill when: writing, reviewing or scaffolding Go code, or placing a Go module and its SPA.

Full guide: invoke the go-architect skill.

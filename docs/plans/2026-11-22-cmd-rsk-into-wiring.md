---
status: active
project: ralvaskills
date: 2026-11-22
related:
    - backlog/pending/a7zmcr-turn-cmd-rsk-into-wiring-move-the-logic-in-update-go-535-ins.md
---

# Turn `cmd/rsk` into wiring

`cmd/rsk` holds about 7000 lines outside tests, and six commands carry most of the logic: `update`, `install`, `status`, `claude` (the `tools` subcommands), `catalog` and `new`. The item `a7zmcr` asks to move that logic into `internal/<area>` and leave `cmd/rsk` with cobra registration and output only. This plan orders the move so each step ships on its own and no output changes.

## Where the files stand now (measured 2026-11-22)

| File | Lines | Cobra-coupled | Natural home |
|---|---|---|---|
| `update.go` + `update_git.go` + `update_lock.go` | 513 + 434 + 106 | `runUpdate*`, `runGitUpdate` take `*cobra.Command` | `internal/update` (344 lines today) |
| `install.go` | 518 | `runInstall*` take `*cobra.Command` | `internal/install` (1091 lines today: constraint solving, source selection) |
| `status.go` | 461 | `runStatus` | new `internal/status` |
| `claude.go` | 398 | five `runClaudeTools*` | `internal/tool` (`ClaudeTool` already lives there) |
| `catalog.go` | 374 | `runCatalogSkills`, `runCatalogBundles` | `internal/catalog` (165 lines today) |
| `new.go` | 307 | `runNew` | new `internal/project` |

The item's line counts (535, 426, 424, 398, 331, 304) have moved; the shape has not.

## The pattern to copy

`internal/sourceops` and `internal/resolve` already do this: plain structs in, plain structs out, no cobra import, and the CLI and the MCP tools both call them. Each moved command follows the same split:

1. an `Options` struct filled from flags, and a `Result` struct (rows, plan, warnings);
2. a function `Run(ctx, cfg, io.Writer, Options) (Result, error)` in `internal/<area>` that does the work and writes human text only through the writer it is given;
3. `cmd/rsk/<command>.go` keeps the cobra command, flag parsing, prompts, and the rendering of `Result`.

The writers matter because some logic prints progress as it goes (`relinkOutdated`, `reportInstall`). Pass the writer; do not return text and print it later, so interleaving with git output stays the same.

## Phases (smallest and least coupled first)

1. **`catalog`** (374 lines, package exists). Move `loadSkillsBySource`, `declaredSources`, `filterByBundle`, `skillsToEntries` and the bundle rows. Warm-up: it proves the pattern and the test move.
2. **`claude`** (398). Move the list, allow, deny and remove operations onto `internal/tool` next to `ClaudeTool`; the command keeps the table rendering.
3. **`new`** (307). New `internal/project` for manifest creation, MCP registration and the rsk-guide install; `cmd` keeps the flag and prompt resolution (`resolveForFlag`, `resolveMCPFlag`, `resolveGuideFlag`).
4. **`status`** (461). New `internal/status` for `scanLinked`, `buildStatusSections`, the git-lock drift and the rename tags; `cmd` keeps the `term.Table` rendering.
5. **`install`** (518). Move the global and project flows and `runInstallFromMod` into `internal/install`, which already owns the solver. The preview and confirm prompt stays in `cmd`.
6. **`update`** (1053 across three files, the largest and the one most tied to the registry, git and lock paths). Split in three steps, registry then git then lock, each its own commit. The two plans (`planGitUpdates`, `findOutdatedSkills`) are already pure enough to move first.

Order can change if a step uncovers coupling, but `update` stays last because it touches the most call paths and the roll-out used it daily.

## Rules for every step

- **No behaviour change.** Output text, exit codes, flags and the data-on-stdout contract stay identical. The step passes `cmd/rsk/output_test.go`, `parity_test.go` (CLI/MCP twins) and the existing per-command tests without edits to their expectations.
- **Tests move with the code.** A `*_test.go` that exercises moved logic moves to the new package; what stays in `cmd/rsk` tests flag parsing and rendering.
- **No cobra under `internal/`.** Add a lint rule (depguard or an import check in `task lint`) before phase 1 ends, so a later change cannot reintroduce it.
- **Gates per step:** `task lint`, `go test ./...`, `task deadcode`, `task lint:doccheck`. A step is not done while any fails.
- **Smoke on a real project copy** after phases 5 and 6: `rsk status`, `rsk update --dry-run` and `rsk install --dry-run` against a copy of a v2 project with the sandboxed config the roll-out used; compare output with the previous binary.
- **One commit per phase** (`refactor(<area>): move <thing> out of cmd/rsk`), no skill, registry or `internal/tool` adapter-config changes mixed in (adapter config-writing is an ask-first area).

## Decisions (owner, 2026-11-22)

- Rendering stays in `cmd/rsk`; only the logic moves.
- ADR 0004 records "`internal/` never imports cobra" (`cmdx` and `ui` excepted); `depguard` in `.golangci.yml` enforces it.
- One patch release (v3.10.1) after all six phases.

## Progress

- ADR 0004 and the depguard rule: `9485d8f`.
- Phase 1, `catalog`: done (`8a5784f`). `catalog.List` and `catalog.BundleRows`; warnings are returned, `cmd/rsk` prints them. Output of `catalog -o json`, `catalog bundles -o json`, `catalog bundle go-grpc -o json` and `catalog --source local` is byte-identical to the released 3.10.0. `task lint`, `go test ./...` and `lint:doccheck` pass; `deadcode` reports only `git.Clone`, which was already unreachable.

## Open questions (answered above, kept for the record)

- Is the target "cobra only" strict (rendering also moves out), or "wiring plus rendering" as above? The plan assumes rendering stays in `cmd`.
- ADR: the rule "`internal/` never imports cobra" is structural; record it as `docs/adr/0004` before phase 1?
- Release: one `v3.11.0` for the whole set, or ship per phase? No user-facing change, so a patch release at the end is enough.

## Exit criteria

- [ ] `cmd/rsk/{update,update_git,update_lock,install,status,claude,catalog,new}.go` each hold only cobra registration, flag parsing, prompts and rendering.
- [ ] `internal/**` imports no `spf13/cobra`, and a lint rule enforces it.
- [ ] `task lint`, `go test ./...`, `task deadcode` and `task lint:doccheck` pass.
- [ ] Output of `rsk status`, `update --dry-run`, `install --dry-run`, `catalog` and `new --help` is identical before and after on a project copy.
- [ ] `a7zmcr` moved to `done/`.

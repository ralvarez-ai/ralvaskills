# ralvaskills backlog

Queue for this repo. The 2026-10-22 declutter audit found 0 dead code and 0 unused
unexported symbols; the items below are size and command-structure work, not cleanup.

## Next

- [ ] Turn `cmd/rsk` into wiring: move the logic in `update.go` (535), `install.go` (426), `status.go` (424), `claude.go` (398), `catalog.go` (331) and `new.go` (304) into `internal/<area>`, leaving `cmd/rsk` with cobra registration only  #ralvaskills  P2
- [ ] Move `cmd/count-tokens/render.go` (451) into `internal/tokensreport`; leave `main.go` as the CLI  #ralvaskills  P2

## Later

- [ ] Split `internal/mcp/tools.go` (325) and `profile.go` (319) by domain if they keep growing  #ralvaskills  P3
- [ ] `setupCommands` in `cmd/rsk/setup.go` has 94 statements; split the registration groups when touching the command tree  #ralvaskills  P3
- [ ] One module host: `mcpkit`, `svckit` and `termkit` are required from
  `github.com/ralvarezdev/*` while the go.ralvarez.dev vanity exists. Platform queue R153 decides the
  host; afterwards this is one sweep of `go.mod` and imports, and the same sweep applies to
  `ralvaskills-executor-port`. Origin: [reuse audit 2026-11-01](../../../../../libraries/go/identity/docs/reviews/2026-11-01-reuse-audit.md)
  #ralvaskills  P3
- [ ] `internal/{config,fsx,cmdx}` (about 400 LOC) against the kits: `config` overlaps
  `viperkit/configparts` and `fsx`/`cmdx` are filesystem and command helpers with no kit today.
  Assess after the R153 sweep; extract only if the executor port needs the same code
  #ralvaskills  P3

- [ ] Remove or use the no-test dead internals: `joinKeys`, `catalog.ExcludePersonal`, `TargetScope.String`, `skill.ParseVersion`, `Version.String/Compare/Before/After`, `ui.Dim`. Check `Stringer`/formatting first; then add `deadcode:check` to CI (the `deadcode` task exists, but no gate)  #ralvaskills  P3

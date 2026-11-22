---
adr: 4
status: accepted
date: 2026-11-22
project: ralvaskills
---

# `cmd/rsk` is wiring, and `internal` never imports cobra

## Context

`cmd/rsk` holds about 7000 lines outside tests. Six commands (`update`, `install`, `status`, `claude`, `catalog`, `new`) carry logic that has nothing to do with cobra: planning updates, resolving constraints, scanning links, building the catalog. Their functions take a `*cobra.Command` only to reach `OutOrStdout`, `Context` and flags, which makes the logic hard to test without building a command and impossible for the MCP server to reuse. `internal/sourceops` and `internal/resolve` already show the other shape: plain structs in, plain structs out, and a CLI and an MCP tool that both call them.

Today only `internal/cmdx` and `internal/ui` import `spf13/cobra` under `internal/`; every other package is already clean.

## Options considered

- **Leave it.** No work, but the commands keep growing and the CLI and MCP surfaces can drift apart, as the CLI/MCP parity test already has to guard against.
- **Move logic and rendering out, leaving only cobra registration.** The strictest reading of the backlog item. It drags the `term` table rendering into `internal/` and changes more code per step for little gain.
- **Move the logic, keep flags, prompts and rendering in `cmd/rsk`.** The `sourceops` shape. `internal/<area>` returns a `Result`; `cmd/rsk` parses flags and renders it.

## Decision

The third option. Domain packages under `internal/` do not import `github.com/spf13/cobra`. `internal/cmdx` (cobra helpers) and `internal/ui` (the terminal UI) are the two exceptions, because they exist to support the CLI. A command's logic lives in `internal/<area>` behind an `Options` struct and a `Result`, writes progress only through an `io.Writer` it is given, and takes a `context.Context`. `cmd/rsk/<command>.go` keeps the cobra command, flag parsing, prompts and the rendering of the `Result`.

The rule is enforced by the `depguard` linter in `.golangci.yml`, so a later change cannot bring cobra back into a domain package.

## Consequences

- Moved logic is testable without cobra, and the MCP tools can reuse it.
- `cmd/rsk` tests shrink to flag parsing and rendering; the logic tests move with the code.
- Each move changes signatures, so it is done one command at a time with output unchanged (`plans/2026-11-22-cmd-rsk-into-wiring.md`).
- A new command that needs cobra types in domain code has to split them first, or earn an exception here.
- Revisit if `cmdx` or `ui` grow logic that is not about presentation.

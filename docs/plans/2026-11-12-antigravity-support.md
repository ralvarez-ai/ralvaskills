---
status: done
project: ralvaskills
date: 2026-11-12
related:
  - plans/2026-11-12-mcp-client-support.md
  - backlog/done/ag7k2m-add-google-antigravity-as-a-skill-and-mcp-target.md
---

# Google Antigravity as a skill and MCP target

## Scope

Add `antigravity` to the tool registry so `rsk` can link skills and register
`rsk mcp` for the Antigravity App, IDE and CLI.

In scope:

- one new adapter, `internal/tool/antigravity.go`, plus its tests
- `clientHasMCP` and `mcpConfigName` know the adapter, so preflight and
  `rsk sync` detect an existing registration
- tool lists in `rsk init`/`new`/`setup` help, the `install_skills` schema,
  README and CHANGELOG

Out of scope:

- writing `~/.gemini/config/mcp_config.json` (rsk registers project config
  only, as with every other client)
- `~/.gemini/antigravity-cli/skills` and `~/.gemini/skills`, which only the CLI
  reads and which outrank `~/.gemini/config/skills`
- Antigravity rules, custom agents, plugins, hooks and `skills.json` manifests

## What Antigravity reads

| Customization | Path |
| --- | --- |
| Global skills | `~/.gemini/config/skills/<name>/SKILL.md` (App, IDE and CLI) |
| Workspace skills | `<root>/.agents/skills/<name>/SKILL.md` (legacy: `.agent/`, `_agents/`, `_agent/`) |
| Workspace MCP | `<root>/.agents/mcp_config.json`, `mcpServers` map |
| Global MCP | `~/.gemini/config/mcp_config.json` |

Two facts drive the design:

1. **`.agents/skills` is already Codex's directory.** `install.ProjectSkillsDirs`
   deduplicates by path, so a project listing both tools links each skill once.
   `rsk.mod` listing both is therefore redundant on disk but harmless; the docs
   say so instead of refusing the combination.
2. **`.gemini` is home-relative, not XDG.** `internal.ConfigHome()` would send
   the global dir to `$XDG_CONFIG_HOME/.gemini/config/skills`, which Antigravity
   never reads, so the adapter uses `os.UserHomeDir()` like the Codex adapter.

Antigravity has no instructions-import file that rsk could generate: it scans
`skills/<name>/SKILL.md` one level deep and reads `AGENTS.md` itself. So
`SyncPinned` and `RemovePinned` are no-ops, matching Codex, and rsk never edits
`AGENTS.md`.

## Phases

1. **Adapter** — `antigravityTool` with ID `antigravity`, no-op pinned sync,
   `~/.gemini/config/skills`, `.agents/skills`, and
   `.agents/mcp_config.json` carrying
   `{"rsk": {"command": "rsk", "args": ["mcp", "--read-only"]}}` through the
   shared `registerMCPServer`/`removeMCPServer` JSON path. Refuse a symlinked
   `.agents` directory, as the Codex adapter does for `.codex`.
2. **Detection** — `clientHasMCP` and `mcpConfigName` gain an `AntigravityID`
   case, so `PreflightMCP`, `rsk new --mcp` and `rsk sync` work without a
   per-tool special case elsewhere.
3. **Registration in the registry** — `Register(&antigravityTool{})` in
   `tool.go`; `rsk init` picks it up from `tool.All()` automatically.
4. **Surfaces** — help text in `cmd/rsk/init.go`, `new.go` and `setup.go`, the
   `install_skills` `for` enum, README §MCP and §skills-directory, CHANGELOG.
5. **Tests** — directory paths, idempotent register/unregister with an
   unrelated server preserved and the file mode kept, conflicting `rsk` entry
   refused untouched, symlinked `.agents` refused, preflight and
   `ProjectHasMCP` recognise the adapter.

## Exit criteria

- [x] `rsk new --for antigravity` links skills into `.agents/skills/` and
      `--mcp` writes `.agents/mcp_config.json` with the read-only server.
- [x] `rsk new --for all --mcp` succeeds; `.agents/skills` holds one link per
      skill even with both `codex` and `antigravity` in `rsk.mod`.
- [x] Register → unregister restores a pre-existing `mcp_config.json`
      byte-for-byte apart from JSON formatting.
- [x] `task lint` and `go test ./...` pass.
- [x] README lists Antigravity's three paths; CHANGELOG has an `Added` entry.

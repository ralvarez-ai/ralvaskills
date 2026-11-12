---
title: 'Add Google Antigravity as a skill and MCP target'
priority: p2
section: ralvaskills
topics:
    - rsk
    - tool
---
Antigravity reads skills from `~/.gemini/config/skills/` (App, IDE and CLI) and `<root>/.agents/skills/`, and MCP servers from `<root>/.agents/mcp_config.json`; both are in scope for an `antigravity` adapter. Two wrinkles the plan records rather than hides: `.agents/skills` is the same directory Codex uses (`ProjectSkillsDirs` dedupes, so the double listing is harmless), and `~/.gemini` is home-relative, so `internal.ConfigHome()`/XDG would point at a directory Antigravity never reads.

Also update `clientHasMCP`/`mcpConfigName` so preflight and `rsk sync` detect the adapter, refresh the tool lists in `rsk init`/`new`/`setup` help, the `install_skills` `for` enum and the README.

Plan: [Google Antigravity as a skill and MCP target](../../plans/2026-11-12-antigravity-support.md).

Not in scope: the global `~/.gemini/config/mcp_config.json`, the CLI-only `~/.gemini/antigravity-cli/skills` and `~/.gemini/skills`, rules, custom agents, plugins and hooks.

Done when `rsk new --for antigravity --mcp` writes `.agents/skills` links plus a read-only `rsk` entry in `.agents/mcp_config.json`, unregister restores the file, `task lint` and `go test ./...` pass, and README/CHANGELOG name the new target.

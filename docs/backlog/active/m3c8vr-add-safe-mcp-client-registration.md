---
title: Add safe MCP registration for all supported clients
priority: p1
section: ralvaskills
topics:
    - rsk
    - mcp
related:
    - ralvaskills#t8k4mw
---

Implement [the MCP client support plan](../../plans/2026-11-12-mcp-client-support.md):
register the server for Claude Code, Codex, and OpenCode; preserve user-owned
configuration; and verify client-specific approval before enabling MCP installs.
Track any client behavior that cannot be verified separately rather than claim
that a client prompts on the basis of advisory MCP metadata.

Progress 2026-11-12: Codex registration uses a managed TOML block; OpenCode
supports v1 config and an existing v2 `mcp.servers` table. Both register a
read-only server until their install approval behavior is verified. Setup
preflights selected MCP configurations and refuses conflicts. OpenCode 2.1.4
normalized the generated v1 entry and connected in a scratch project.
The v3.8.0 CLI release includes this read-only integration; it does not close
the approval and configuration-compatibility work below.

Progress 2026-11-12 (later): `rsk sync` shipped in v3.9.0 and reconciles an
existing registration without touching a project that has none; v3.9.1 repairs
versionless official skills and upgrades OpenCode v2 entries written in the old
v1 shape. The machine-wide pinned version is now 3.9.1 in all 72 consumer
`mise.toml` files, and all 71 projects with a manifest have been synced, so the
read-only entries are live rather than deferred. Release and sweep detail:
[t8k4mw](../done/t8k4mw-add-rsk-sync-and-sweep-consumers.md).

Remaining: verify Codex project discovery and per-install approval in real
clients; enable the mutating tool only after confirmation is enforced. Decide
how to select native v2 for new OpenCode configs without breaking v1, support
or deliberately handle JSONC and higher-precedence config locations, improve
ownership tracking and recovery from failures after preflight, and test removal,
symlink and migration edge cases. The current adapter refuses alternate
project OpenCode config files and reports partial registration if a write fails.


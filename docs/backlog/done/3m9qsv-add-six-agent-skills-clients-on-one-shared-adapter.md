---
title: 'Add six Agent Skills clients on one shared JSON adapter'
priority: p2
section: ralvaskills
topics:
    - rsk
    - tool
---
Kimi Code, Factory Droid, Gemini CLI, Qwen Code, Cursor CLI and Sourcegraph Amp now join Claude Code, Codex, OpenCode and Antigravity. All six read `.agents/skills` in both scopes, so `rsk new --for all` links each skill once for eight clients and only the global directory and the MCP file are per-client.

`internal/tool/json.go` carries the shared shape: `jsonAdapter{id, skillsDir, projectDir, mcp}` and `jsonMCPConfig{file, key, server, alternate}`, reusing `registerMCPServer`/`removeMCPServer`. `mcpConfigName` and `clientHasMCP` go through `jsonMCPProvider`, so detection, preflight and `rsk sync` cover a new client without a case in `mcp.go`. Claude Code's `.mcp.json` moved onto the same path. `--for` help, `rsk new`/`rsk init` prose and the `install_skills` schema derive from `tool.Names()`.

Per-client details: Cursor needs `"type": "stdio"`; Amp uses the literal dotted key `amp.mcpServers` and refuses an existing `.amp/settings.jsonc`; Kimi's global directory moves with `KIMI_CODE_HOME`; Gemini and Qwen share their MCP config with hooks and permissions, so only a file that parsed is re-encoded. Every new registration is `--read-only` until its approval behavior is verified.

Plan: [Agent Skills clients: one adapter, ten targets](../../plans/2026-11-12-agent-skills-client-roadmap.md). Tier 2 (Crush, Goose, OpenHands) stays open in `9krt4v`; approval verification is `7wq4tn`.

Done when all ten targets register in one project, a second `rsk new --for all --mcp` is a no-op, `rsk destroy` removes every entry, and `task lint` plus `go test ./...` pass — all verified 2026-11-12 (312 tests).

Follow-up: the remaining plumbing in the bespoke clients — duplicated OpenCode JSON pair, three symlink guards, the OpenCode key check in the generic writer, and per-client MCP code left in `mcp.go` — is recorded in `5vq2hd`.

---
title: 'Fold the remaining MCP plumbing and fix the symlink messages'
priority: p2
section: ralvaskills
topics:
    - rsk
    - tool
---
The shared `jsonAdapter`/`jsonMCPConfig` removed the per-client boilerplate, but four near-copies and one tool-specific branch survived in the bespoke clients:

- `opencode.go` carried `readOCConfig`/`writeOCConfig`, byte-for-byte duplicates of the shared JSON pair.
- `refuseIfLink`, `refuseIfLinkDir` and an inline `Lstat` in `readCodexConfig` were three guards for one rule; the wording told a Kimi or Antigravity user to "point CLAUDE.md at a real file".
- `rejectAlternateOpenCodeConfigs` repeated `jsonMCPConfig.checkAlternates`.
- `registerMCPServer`/`removeMCPServer` compared the config key against OpenCode's to allow the pre-read-only entry, so the generic writer knew one client's history.

`linkTarget` also read the file's base name instead of its path, so every symlink refusal printed `→ unknown` — the one fact that tells you where the link goes.

Resolved 2026-11-12 in three commits: `7f0de5d` deleted the OpenCode JSON pair, unified the guards (`refuseLink` with `wantDir`), extracted `rejectAlternates`, fixed `linkTarget` to take the path, and replaced the key test with `legacy []map[string]any` passed by the caller. `d8422c7` moved each bespoke client's MCP code beside its adapter (`claude_mcp.go`, `opencode_mcp.go`, Codex's recognizer in `codex.go`) and introduced `mcpRecognizer` (`hasMCP`, `mcpConfigName`) to mirror `jsonMCPProvider`, so `mcp.go` is 243 lines of engine with no tool-ID switch. `791d6ec` recorded both in the changelog.

Deliberately unchanged: Codex's comment-delimited TOML block (a re-encode would eat user comments), OpenCode's acceptance of the legacy `mcp.rsk` entry and `.claude/skills/` instructions (consumers mid-migration depend on it), and Claude's per-tool `SyncPinned` (its `@AGENTS.md`/`@SKILLS.md` ordering has its own link tests).

Verified: 313 tests, `task lint`, `task lint:doccheck`, `deadcode`, and a live OpenCode v2 project where registration stayed idempotent. `TestSymlinkRefusalNamesTheTarget` fails on the old `linkTarget`.

Plan: [Agent Skills clients: one adapter, ten targets](../../plans/2026-11-12-agent-skills-client-roadmap.md).

---
status: done
project: ralvaskills
date: 2026-11-12
related:
  - plans/2026-11-12-antigravity-support.md
  - plans/2026-11-12-mcp-client-support.md
  - backlog/done/3m9qsv-add-six-agent-skills-clients-on-one-shared-adapter.md
  - backlog/done/5vq2hd-fold-the-remaining-mcp-plumbing-and-fix-symlink-messages.md
  - backlog/pending/9krt4v-add-crush-goose-and-openhands-mcp-targets.md
  - backlog/pending/7wq4tn-verify-antigravity-mcp-approval-and-decide-on-install-skills.md
---

# Agent Skills clients: one adapter, ten targets

## Scope

Support every CLI that reads the Agent Skills convention (`.agents/skills` plus
its own skills directory) and keeps MCP servers in a JSON object.

In scope:

- a shared `jsonAdapter` and `jsonMCPConfig` in `internal/tool`, so a client
  adapter is three fields instead of eight methods
- six new targets: Kimi Code, Factory Droid, Gemini CLI, Qwen Code, Cursor CLI,
  Amp
- moving Claude Code and Antigravity onto the same JSON MCP path, and listing
  the registry in CLI help instead of editing each command

Out of scope, documented in this plan but not built:

- **Tier 2** — Crush, Goose, OpenHands: no plain `mcpServers` JSON, or no
  `SKILL.md` discovery
- **Tier 3** — clients that need no adapter or cannot be supported safely

## How the code is shared

`internal/tool/json.go` holds the whole per-client surface:

```go
type jsonAdapter struct {
	id         ID
	skillsDir  func() string        // global skills directory
	projectDir func(string) string  // project skills directory
	mcp        jsonMCPConfig
}
```

`jsonMCPConfig` is `{file, key, server, alternate}`: the project-relative config
path, the key holding the server map (matched literally, so Amp's dotted
`amp.mcpServers` works), the entry rsk registers, and any alternate file the
client also reads. `Register`/`Unregister` reuse the existing
`registerMCPServer`/`removeMCPServer`, so idempotence, conflict refusal, atomic
writes and mode preservation were already solved.

Three consequences:

1. **A new JSON client is ~30 lines**: constants plus a constructor call. It
   needs no change in `mcp.go` — `clientHasMCP`, `mcpConfigName` and
   `ProjectMCPFile` go through the `jsonMCPProvider` interface.
2. **Skills cost nothing per client.** Every Tier 1 client reads
   `.agents/skills`, so `agentsSkillsDir` is their shared
   `ProjectSkillsDir`; `ProjectSkillsDirs` already deduplicates by path and
   `SyncPinned` is a no-op. Ten tools in `rsk.mod` produce one set of links.
3. **CLI lists read the registry.** `--for` help, `rsk new`/`rsk init` prose and
   the `install_skills` schema derive from `tool.Names()`; `toolsFromFlag`
   parses through `tool.ParseID` and builds `all` from `tool.All()`.

Claude Code keeps its own type because it writes the `CLAUDE.md` import, but its
MCP registration is now the same `jsonMCPConfig` path. Codex stays TOML and
OpenCode keeps its `mcp.servers`/legacy `mcp.rsk` handling.

## Tier 1 — built

| ID | Client | Global skills | Project skills | MCP file | Note |
| --- | --- | --- | --- | --- | --- |
| `antigravity` | Google Antigravity | `~/.gemini/config/skills` | `.agents/skills` | `.agents/mcp_config.json` | home-relative, ignores XDG |
| `kimi` | Kimi Code CLI | `$KIMI_CODE_HOME/skills`, else `~/.kimi-code/skills` | `.agents/skills` | `.kimi-code/mcp.json` | legacy `kimi-cli` is archived |
| `droid` | Factory Droid | `~/.factory/skills` | `.agents/skills` | `.factory/mcp.json` | also reads `.factory/skills` |
| `gemini` | Gemini CLI | `~/.gemini/skills` | `.agents/skills` | `.gemini/settings.json` | MCP shares settings with hooks and permissions |
| `qwen` | Qwen Code | `~/.qwen/skills` | `.agents/skills` | `.qwen/settings.json` | renames unparseable settings to `.corrupted` |
| `cursor` | Cursor CLI | `~/.cursor/skills` | `.agents/skills` | `.cursor/mcp.json` | needs `"type": "stdio"` |
| `amp` | Sourcegraph Amp | `~/.config/agents/skills` | `.agents/skills` | `.amp/settings.json` | dotted `amp.mcpServers`; `.jsonc` refused |

Every one registers `rsk mcp --read-only` until its approval behavior is
verified, matching Codex and OpenCode.

## Tier 2 — documented, not built

| Client | Skills | MCP | Why it is not a `jsonAdapter` |
| --- | --- | --- | --- |
| Crush | `.agents/skills`, `.crush/skills`, `~/.config/crush/skills` | `crushrc` (Bash, `mcp add …`); `crush.json` merges but is deprecated | Writing a Bash config is a different problem from merging JSON. `crush.json` still works, so a JSON adapter is possible and would be the smallest path |
| Goose | `~/.config/goose/` recipes | `~/.config/goose/config.yaml`, `extensions:` key | YAML, not JSON; a new reader/writer pair is needed |
| OpenHands | `.openhands/microagents/` | `~/.openhands/` config | "microagents" are not `SKILL.md` discovery, so `SyncPinned`/`ProjectSkillsDir` do not map |

Crush is the only one where the existing code is close: `crush.json` accepts the
same `{type, command, args}` server objects, and the only caveat is that Crush
prints a deprecation warning for JSON config.

## Tier 3 — excluded

| Client | Reason |
| --- | --- |
| GLM Coding Plan | Not a CLI. It is an Anthropic/OpenAI-compatible endpoint for Claude Code, Codex, OpenCode, Crush and Droid; `@z_ai/coding-helper` configures those clients and installs nothing of its own |
| Aider | No MCP support and no skill discovery; a `CONVENTIONS.md` link is not a target |
| Cline, Roo, Windsurf, Zed | IDE-first; no project-level agent config rsk can own |
| Claude Desktop | User-level config only, no project scope |
| `zcode-cli-x` | Third-party wrapper around the GLM API, not an agent consumers run |

## Risks the adapters carry

- **Shared settings files.** Gemini, Qwen and Amp keep MCP servers beside hooks,
  permissions and model settings. `writeJSONConfig` re-encodes the whole object,
  so formatting is normalised; keys are preserved. Qwen's habit of renaming an
  invalid `settings.json` to `.corrupted` is why rsk only re-encodes a file it
  just parsed.
- **Alternate configs.** Amp reads `.amp/settings.jsonc` as well; registration
  refuses when it exists rather than creating a second source of truth.
- **Cursor CLI.** Some releases have been reported not to read the project
  `.cursor/mcp.json` while the editor does. rsk writes the project file and the
  docs say so.
- **Gemini's personal tier** moved to Antigravity on 2026-06-18; the target
  remains for Gemini Code Assist enterprise and API-key users.

## Exit criteria

- [x] `rsk new --for all --mcp` registers all ten clients and writes one skills
      directory per target.
- [x] A registered adapter appears in `--for` help and `rsk new`/`rsk init`
      prose without editing those commands.
- [x] A JSON adapter is a constructor plus constants; `mcp.go` has no per-client
      case for them.
- [x] Table tests cover every JSON adapter for paths, idempotent register,
      unrelated-key preservation, conflict refusal and symlinked directories.
- [x] Tier 2 and Tier 3 are recorded here, and Crush/Goose/OpenHands have a
      backlog item.

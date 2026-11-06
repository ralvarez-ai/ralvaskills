# Token Economy — Tools

What each token tool cuts, how it is wired, and what it costs in turn. Referenced from
[SKILL.md](SKILL.md) section 8; install commands in [RECIPES §6](RECIPES.md#6-tool-setup).

| Tool | Cuts | Claude Code | OpenCode | Detect |
|---|---|---|---|---|
| context-mode | tool outputs (sandboxed, indexed, searched) | plugin | `opencode plugin add context-mode` | `ctx doctor`; `opencode plugin list` |
| RTK | shell command output | `rtk init -g --hook-only` (PreToolUse on Bash) | `rtk init -g --opencode --hook-only` | `rtk gain` |
| CodeGraph | grep + read for structure questions | `codegraph install -t claude`, `codegraph init` per repo | `codegraph install -t opencode` | `codegraph status`, `.codegraph/` |
| caveman | reply prose | skill | skill | — |

## What the tools themselves load

- **RTK** without `--hook-only` adds an `RTK.md` to the always-loaded instructions.
- **CodeGraph** writes a short marked block into the global `CLAUDE.md` / `AGENTS.md`. In Claude
  Code it also registers a `UserPromptSubmit` hook (`codegraph prompt-hook`) that, in indexed
  repos, pastes matching source into any prompt that names indexed symbols: from nothing to
  thousands of tokens per prompt, often unrelated to the question. Remove the hook and keep
  the MCP tool, which the agent calls when it needs it. `codegraph install` and `upgrade` put
  it back, so strip it after each (the setup repo's install task is the place).
- **context-mode** registers session-start and pre-Agent hooks; some versions inject tens of
  thousands of tokens on resume. Measure with `scripts/session-usage.py` before keeping it on
  long sessions.

## Overlap and pitfalls

- **RTK and context-mode both touch shell output.** RTK shrinks the command's output;
  context-mode keeps what is still large out of the context. After installing both, run a
  `git diff` and a test command once and confirm the output is neither lost nor doubled.
- **RTK only sees shell calls.** Read, Grep and Glob bypass it. A different `rtk` exists on
  crates.io: `rtk gain` failing means the wrong binary.
- **CodeGraph in Claude Code:** register it with `alwaysLoad: true` so its tool is not deferred
  behind tool search, or the agent greps instead. Indexing is per repo; ignore `.codegraph/`.
  Telemetry is on until `codegraph telemetry off`.
- **OpenCode 2.1.x:** `opencode mcp list` does not show servers under `mcp.servers`; confirm
  with `opencode debug config` or by asking the agent to list its tools.
- **Vendor savings figures are their own benchmarks.** Measure: `rtk gain`, `ctx stats`,
  `scripts/session-usage.py`.

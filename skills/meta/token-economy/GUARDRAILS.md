# Token Economy — Claude Code Guardrails

Settings, permission rules and hooks that enforce [SKILL.md](SKILL.md) section 7 instead of
asking for it. Order of preference: a setting (paid once, always applies), then a hook (runs
outside the model; only its output costs tokens), then an instruction (paid every turn,
enforces nothing). Check each key against your Claude Code version before applying.

## Settings

`~/.claude/settings.json` (user scope) or a project's `.claude/settings.json`:

```json
{
  "autoCompactWindow": 200000,
  "promptCacheTtl": "1h",
  "subagentPromptCacheTtl": "1h",
  "env": { "CLAUDE_CODE_FORK_SUBAGENT": "0" }
}
```

| Key | Effect | Trade-off |
|---|---|---|
| `autoCompactWindow` | Compacts at this context size instead of near the model's full window | A cap is the cheapest control; a project with legitimately long context overrides per session (`/autocompact`) |
| `promptCacheTtl: "1h"` | Main thread cache survives idle gaps up to an hour | 1h writes cost more than 5m writes; pays off once a session has gaps over 5 minutes |
| `subagentPromptCacheTtl: "1h"` | Same for subagents, forks and compaction | Same; matters most for subagents that poll |
| `CLAUDE_CODE_FORK_SUBAGENT=0` | The agent stops opening forks that inherit the full context | Forks remain available by hand |

**Cap or checkpoints, not both.** A cap holds the context at its value, so confirmations at
higher thresholds never fire. Pick a cap for the lowest bill, or the full window plus
confirmation checkpoints (hook below) for control.

## Permission rules

```json
{ "permissions": { "deny": ["Read(//tmp/claude-*/**)", "Read(//tmp/**/*.output)"] } }
```

Stops raw task outputs and transcripts being read into the context. Keep rules **scoped**: a
bare rule such as `"Read"` changes the tool definitions and invalidates the cache prefix.

## Hook patterns

Each is a small shell script; its output, when any, is the only token cost.

| Hook | Event | Guards against |
|---|---|---|
| Context checkpoint | `PreToolUse` | Context growing past thresholds (e.g. 200k, 400k) without the user agreeing; asks once per threshold |
| Context meter | `Stop` | Re-parsing the transcript in other hooks; records context size and time per turn |
| Cache watch | `UserPromptSubmit` | Paying an expired cache blind: `last turn + TTL < now` predicts a rebuild; warn, and block once above a cost threshold so the user can `/clear` or resend |
| Deep-fork guard | `PreToolUse` (Agent) | A fork opened when the context is already large |
| Transcript-read guard | `PreToolUse` (Read, Bash) | Raw transcripts and task outputs read into the context |
| Subagent report contract | `SubagentStop` | A subagent's final message dumping logs into the parent; block once, ask for a short report |
| Model-switch guard | `PreModelSwitch` | Switching model on a deep context, which discards the cache |
| Loop concision | `UserPromptSubmit` | Long replies to polls and automated reports |

**Avoid:**

- `prompt` or `agent` hook handlers: they spend model tokens on every fire.
- Blind output trimming (`PostToolUse` rewriting outputs to N characters): it hides the error
  you were looking for; only with a command allowlist.
- Blocking automatic compaction (`PreCompact`): the opposite of what is needed.
- Hooks that inject context on every prompt or every tool call; measure what any third-party
  hook adds with `scripts/session-usage.py`.

## Measuring

Transcripts live at `~/.claude/projects/<project>/<session>.jsonl`, subagents under
`<session>/subagents/`. Each `assistant` record carries the API `usage` (input, cache reads,
cache writes split by TTL, output, thinking); records repeat per content block, so deduplicate
by message id. Hook output appears as `attachment` records of type `hook_success`.
`scripts/session-usage.py` does both sums.

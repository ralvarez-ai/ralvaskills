---
name: token-economy
version: 1.0.0
description: Token-saving habits and setup for coding agents — context hygiene, batched tool calls, subagent exploration with model/effort routing, escalation and session handoff, cheap-model (DeepSeek Flash) guidance, and an audit of context-mode, RTK and CodeGraph for Claude Code and OpenCode. Use when the user wants to save tokens or cost, shrink context, choose models or subagents for a task, set up RTK/CodeGraph/context-mode, or invokes /token-economy.
---

# Token Economy

Spend tokens where they buy accuracy and nowhere else. Most waste is not long replies: it is
bytes that enter the context once and are re-sent on every later turn, retries on a model that
is out of its depth, and sessions that carry three unrelated tasks. This skill has two modes:
**working rules** (sections 2-6, 9) and an **audit** of the machine and project (section 8).
Templates in [RECIPES.md](RECIPES.md); the audit is [scripts/audit.sh](scripts/audit.sh).
Reply-style compression is [caveman](../caveman/SKILL.md); how skills get imported into a
project is [rsk-guide](../rsk-guide/SKILL.md).

## 1. When to invoke

- **Audit / setup** — new machine or project, "where are my tokens going", installing
  context-mode, RTK or CodeGraph, sessions hitting limits early.
- **Working** — before a broad exploration, when choosing a model or subagent, when a session
  has grown long or drifted across tasks, when a cheap model starts looping.
- **Not this skill** — "just be terse" is [caveman](../caveman/SKILL.md).

## 2. Where tokens go

Ranked by what they usually cost over a session:

1. **Resent context.** Every turn re-sends the whole conversation. A 5k-token file read on turn
   3 of a 40-turn session is paid ~37 more times. Caching makes repeats cheaper, not free, and
   cached tokens still count toward rate limits.
2. **Always-loaded instructions.** `CLAUDE.md` / `AGENTS.md` and everything they `@import`,
   OpenCode `instructions`, MCP tool schemas, skill descriptions. Paid on every turn of every
   session. A whole `SKILL.md` imported here defeats progressive disclosure.
3. **Turns.** Each extra round trip re-reads the prefix. Sequential calls that could have been
   one batched turn multiply this.
4. **Retries.** A failed edit or wrong guess costs its input, its output, and the cleanup.
5. **Thinking.** Billed as output. High effort on mechanical work is pure waste.
6. **Reply text.** Output is priced several times input, but it is the smallest volume.

**Cache busters:** changing the model, the instructions, or the MCP server set mid-session
discards the cached prefix. Decide those at session start.

**Rule: cut what persists before cutting what is said.**

## 3. Always-on rules

Habits only work if they are loaded every session, so they live in `AGENTS.md`, not here. Offer
to paste the block from [RECIPES §1](RECIPES.md#1-agentsmd-rules-block) (~25 lines). Its
substance:

- **Locate, then read.** Grep/glob first; read line ranges; never re-read a file you just
  edited (the edit tool already failed if it did not apply).
- **Ask tools for the smallest useful answer.** `git diff --stat` before the diff, `-q` flags,
  failing tests only, errors only from builds, `tail`/`grep` on logs. Large outputs go through
  context-mode or RTK.
- **Batch independent calls** into one turn.
- **Explore through a subagent** when the answer needs many files; keep the conclusion.
- **Replies lead with the result.** No recap of the diff, no narration of steps, no closing
  summary of what was just shown. Code, commits and docs stay in full prose.
- **One task per session.** Clear or start a new session when the task changes; carry state in
  a handoff note, not in the transcript.
- **Two-failure rule.** Second failed attempt at the same thing: stop, escalate (section 5).

## 4. Exploring with subagents

A subagent reads the files in its own context and returns only a summary. It is not free: it
starts with its own system prompt and tool schemas, so a lookup of 1-3 tool calls is cheaper
inline.

| Task | Model | Effort | Ask it to return |
|---|---|---|---|
| Locate (where is X, which files) | haiku | low | paths + one line each |
| Understand (how does X flow) | haiku or sonnet | medium | <=300 words with `file:line` |
| Sweep / review many files | sonnet | medium | findings list, ranked |
| Design call, hard bug | opus | high | a recommendation, not a survey |

- **Brief it like a peer:** goal, what is ruled out, where to look, the output shape and a
  word cap. A vague brief returns a vague summary and a second spawn.
- **Reports come back in caveman style.** Only a model reads them, so fragments with exact
  paths, symbols and `file:line` refs lose nothing. Put the request in the brief and in custom
  agent prompts ([caveman](../caveman/SKILL.md)); the same goes for handoff notes. Never for
  commits, docs, code or replies to the user.
- **Fan out only for independent questions**; send them in one turn.
- **Don't redo its work** inline; if the summary is thin, send a follow-up to the same agent.
- **Claude Code:** the built-in Explore agent runs on the *main* model. Define a project or user
  agent named `Explore` with `model: haiku` to make it cheap; `effort` is a frontmatter field
  ([RECIPES §2](RECIPES.md#2-claude-code-agents)).
- **OpenCode:** subagents are `agents` entries with `mode: subagent`; effort is a model
  variant (`provider/model#high`) ([RECIPES §3](RECIPES.md#3-opencode-config)). Its
  built-in `explore` already inherits the session model; overriding its `system` replaces
  the tuned prompt, so ask for terse reports in the brief instead.

## 5. Choosing the model and moving up

- **Default tier:** Sonnet at medium effort in Claude Code; the cheap model (e.g. DeepSeek
  V4.1 Flash) in OpenCode.
- **Escalate when:** the two-failure rule fires; a design decision will outlive the session;
  a bug's cause is not localised after one focused pass.
- **How to escalate**, cheapest first:
  - **Subagent at the higher tier** — the problem fits in a brief and the answer fits in a
    summary. Template in [RECIPES §5](RECIPES.md#5-subagent-brief).
  - **Fresh session at the higher tier** — the work needs many turns there, or the current
    context is polluted by dead ends. Write a handoff note
    ([RECIPES §4](RECIPES.md#4-handoff-note)) and tell the user to open the new session with it.
  - **Switch model in place** — only when the context is still small; it discards the cache.
- **Come back down** once the hard part is decided: mechanical edits go back to the cheap tier.
- **Effort:** low for mechanical edits, medium by default, high for real reasoning; `xhigh` /
  `max` only when the user asks.

## 6. Cheap models (DeepSeek V4.1 Flash and similar)

When tokens are already cheap, the waste is turns and retries, and quality falls off faster
with context length.

- **Don't compress instructions.** Small models need explicit steps and examples; terse
  prompts cost accuracy. Terse *replies* are fine.
- **Keep the prefix stable.** DeepSeek caches prefixes automatically and a hit costs about 2%
  of a miss; don't toggle MCP servers or instructions mid-session. Off-peak is half price.
- **Compact earlier.** Lower the compaction threshold so the model works on a short context
  ([RECIPES §3](RECIPES.md#3-opencode-config)).
- **Verify by running, not rereading.** A test run is cheaper than re-reading code to check it.
- **Retry narrowly.** After a failed edit, re-read only the target range, then retry once. A
  second failure escalates (section 5).
- **Hand the model smaller tasks.** One file, one change, explicit acceptance check.

## 7. Tools

| Tool | Cuts | Claude Code | OpenCode | Detect |
|---|---|---|---|---|
| context-mode | tool outputs (sandboxed, indexed, searched) | plugin | `opencode plugin add context-mode` | `ctx doctor`; `opencode mcp list` |
| RTK | shell command output | `rtk init -g` (PreToolUse hook) | `rtk init -g --opencode` (OpenCode >= 1.18.29) | `rtk gain` |
| CodeGraph | grep + read for structure questions | `codegraph install`, then `codegraph init` per repo | same `codegraph install` | `codegraph status`, `.codegraph/` |
| caveman | reply prose | skill | skill | — |

- **context-mode, RTK and CodeGraph are expected; their absence is a gap** the audit reports.
- **Watch what the tools themselves load.** Install RTK with `--hook-only` (otherwise it adds
  an `RTK.md` to the always-loaded instructions). CodeGraph writes a short marked block into
  the global `CLAUDE.md` / `AGENTS.md` and a prompt hook that stays silent outside indexed
  repos. context-mode's session-start and Agent hooks can inject tens of thousands of tokens
  on resume in some versions: measure them from the transcript before keeping it on long
  sessions.
- **RTK and context-mode both touch shell output.** RTK shrinks the command's output;
  context-mode keeps what is still large out of the context. After installing both, run a
  `git diff` and a test command once and confirm the output is neither lost nor doubled.
- **RTK only sees shell calls.** Read, Grep and Glob bypass it.
- **CodeGraph in Claude Code:** register it with `alwaysLoad: true` so its tool is not
  deferred behind tool search, or the agent will grep instead. Indexing is per repo
  (`codegraph init`); ignore `.codegraph/`. Telemetry is on until `codegraph telemetry off`.
- **Keep the setup in a config repo.** One repo per agent (or one for both) with an
  idempotent `task install` / `task check`, pinned tool versions, a backup before every
  edit, and marked blocks for instruction files so other tools' blocks survive. A new
  machine is then one command, and `check` doubles as the audit. Sketch in
  [RECIPES §7](RECIPES.md#7-setup-repo).
- **Vendor savings figures are their own benchmarks.** Measure yours: `rtk gain`, `ctx stats`.
- Install commands and versions in [RECIPES §6](RECIPES.md#6-tool-setup).

## 8. Audit mode

1. Run `bash scripts/audit.sh [project-dir]` (read-only). It reports `OK` / `GAP` / `WARN` for:
   tools installed and wired into each agent, the CodeGraph index, an `Explore` override,
   always-loaded instruction size (with `@imports` followed), whole `SKILL.md` files imported at
   start, and MCP server count.
2. Turn the output into a short table: finding, cost, fix.
3. **Ask before fixing.** Installing RTK, CodeGraph or context-mode, and editing `CLAUDE.md`,
   `AGENTS.md`, `opencode.json` or agent definitions all change global or shared config.
4. Offer the rules block (section 3) if `AGENTS.md` lacks one.

## 9. Anti-patterns

- Reading a whole file to find one function; re-reading after an edit.
- A subagent for a single grep; five sequential calls that could be one turn.
- Opus or high effort on renames and mechanical edits.
- Caveman style in commits, docs or code comments.
- One long session across unrelated tasks, relying on compaction to forget.
- Polling a log by re-reading it every minute.
- Importing whole `SKILL.md` bodies from `CLAUDE.md` or OpenCode `instructions`.
- Adding or removing MCP servers mid-session.
- Quoting a tool's savings percentage without measuring it here.

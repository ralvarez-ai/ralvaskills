# Token Economy — Recipes

Templates referenced from [SKILL.md](SKILL.md). Copy, then adjust names and caps.

## 1. AGENTS.md rules block

Paste near the top of the project's `AGENTS.md` (always loaded by Claude Code via `CLAUDE.md`
and by OpenCode natively). Keep it this short; every line is paid every turn.

```markdown
## Token discipline

- Locate before reading: grep/glob first, then read only the needed line ranges.
- Never re-read a file after editing it.
- Ask tools for the smallest useful output: `git diff --stat` first, quiet flags, failing
  tests only, errors only from builds, tail/grep for logs.
- Put independent tool calls in one turn. Prefer targeted edits to full-file rewrites.
- Ask one question (or show a short plan) before building on an ambiguous request.
- Never ask the user to paste long logs; run the command, or write output to a file and read
  a range.
- Delegate repetitive edits and checklist audits: a script first, else a cheap worker with one
  worked example and an acceptance command.
- When an answer needs many files, delegate the reading to a subagent and keep its summary.
- End every subagent brief with: "Report in caveman style: fragments, no filler; keep paths,
  symbols and code exact; `file:line` refs; under N words." Same for notes meant only for
  another model.
- Lead replies with the result. No recap of diffs, no step narration, no closing summary.
  Code, commits and docs stay in full sentences.
- One task per session. When the task changes, write a handoff note and start fresh.
- After two failed attempts at the same fix, stop and escalate to a stronger model
  (subagent with a brief, or a new session with a handoff note).
- Use low effort for mechanical edits, medium by default, high only for real reasoning.
```

## 2. Claude Code agents

`.claude/agents/Explore.md` (project) or `~/.claude/agents/Explore.md` (user). Overrides the
built-in Explore, which otherwise runs on the main model.

```markdown
---
name: Explore
description: Read-only search agent. Locates code and summarises how it fits; returns conclusions, not file dumps.
model: haiku
effort: low
tools: Read, Grep, Glob, Bash
---
Answer the question in the brief. In repos with a `.codegraph/` index, ask CodeGraph first,
then read only the ranges it points to; otherwise grep first and read line ranges.

Report in caveman style: fragments, no filler, keep paths, symbols and code exact, `file:line`
refs, under the brief's word cap (200 words without one). Never paste whole files.
```

Escalation agent, `.claude/agents/deep-reasoner.md`:

```markdown
---
name: deep-reasoner
description: Hard design decisions and bugs that resisted two attempts. Use with a full brief.
model: opus
effort: high
---
Work only from the brief and the files it names. Return one recommendation with its reason,
the risks, and the exact next steps. No survey of options you would not pick.
```

Worker for delegated mechanical work, `.claude/agents/worker.md`:

```markdown
---
name: worker
description: Applies a decided pattern across files, runs checklist audits and repetitive edits. Needs a spec, one worked example and an acceptance command.
model: haiku
effort: low
---
Apply the spec exactly as the worked example shows; never redesign. After each batch run the
acceptance command. Report in caveman style: files changed, items that failed with the error,
nothing else. An item that fails twice is reported, not retried.
```

Raise `effort` to `medium` on Explore for "how does X work" questions, or add a second agent.

## 3. OpenCode config

`opencode.json(c)` (project) or `~/.config/opencode/opencode.jsonc`. OpenCode 2 uses
`agents`, `mcp.servers` and `plugins`; the older `agent` key is deprecated. Effort is a model
variant, not an agent field. Check `opencode --version` and the schema at
`https://opencode.ai/config.json` before copying.

```json
{
  "$schema": "https://opencode.ai/config.json",
  "agents": {
    "escalate": {
      "description": "Hard design calls and bugs that resisted two attempts",
      "mode": "subagent",
      "model": "anthropic/claude-sonnet-5-5#high",
      "system": "Work from the brief. Return one recommendation, its risks and next steps."
    }
  },
  "compaction": { "auto": true, "buffer": "25%" }
}
```

- Leave the built-in `explore` alone: it inherits the session model, and setting `system`
  replaces its prompt. `permissions` is an array of `{action, resource, effect}` rules.
- `opencode mcp list` (2.1.x) does not show servers under `mcp.servers`; confirm with
  `opencode debug config` or by asking the agent to list its tools.

- OpenCode loads skills from `.opencode/skills/`, `.claude/skills/` and `.agents/skills/` with
  progressive disclosure, so listing `SKILL.md` files under `instructions` only adds cost.
- A larger compaction `buffer` compacts earlier; useful for cheap models (SKILL section 6).

## 4. Handoff note

Write to `docs/` (if the project keeps notes there) or a scratch file, then tell the user:
"open a new session on <model> and point it at <path>".

```markdown
# Handoff: <task in one line>

- **Goal:** <what done looks like, testable>
- **State:** <what is finished, with commit or file refs>
- **Decided:** <decisions and why, so they are not reopened>
- **Ruled out:** <approaches tried and why they failed>
- **Next:** <the first concrete step>
- **Files:** <paths that matter, with line ranges>
- **Checks:** <commands that prove it works>
```

## 5. Subagent brief

```text
Goal: <one sentence>.
Context: <what you already know; what is ruled out>.
Look at: <paths, symbols, docs>. Out of scope: <what not to touch or read>.
Return: <shape — list of file:line, or a summary> in at most <N> words.
Say "not found" rather than guessing.
```

## 6. Tool setup

Verified against each project's README on 2026-11-06; check the README again before
installing, since all three move fast. **Ask the user before running any of these** — they
change global agent config.

| Tool | Install | Wire into agents | Check |
|---|---|---|---|
| context-mode | `claude plugin marketplace add mksglu/context-mode`, `claude plugin install context-mode@context-mode` | OpenCode: `opencode plugin add context-mode`; update with `opencode plugin update` (no target) | `ctx doctor`, `ctx stats`, `opencode plugin list` |
| RTK (`rtk-ai/rtk`) | the repo's `install.sh` (checksum-verified, `RTK_VERSION=vX.Y.Z` pins) or `brew install rtk` | `rtk init -g --hook-only --auto-patch`; add `--opencode` for OpenCode (also re-applies the Claude hook) | `rtk --version`, `rtk gain` |
| CodeGraph (`colbymchenry/codegraph`) | `npm i -g @colbymchenry/codegraph@<version>` | `codegraph install -t claude,opencode -l global -y`; preview with `--print-config <id>`; `codegraph init` per repo | `codegraph status` |

- RTK: a different `rtk` exists on crates.io; `rtk gain` failing means the wrong binary.
- CodeGraph manual Claude Code entry (in `~/.claude.json` `mcpServers`):
  `{"type": "stdio", "command": "codegraph", "args": ["serve", "--mcp"], "alwaysLoad": true}`.
- Add `.codegraph/` to `.gitignore`.
- Restart the agents after wiring; context-mode upgrades apply on the next session.

## 7. Setup repo

One repo holds what gets applied to an agent's global config, so a machine is reproducible
and drift is visible.

```text
agents/            subagent definitions (e.g. Explore on a small model)
instructions/      the rules block, applied between markers
scripts/           apply-block (idempotent marked-block writer, backs up first)
Taskfile.yml       install, install:tools, check, verify, uninstall, lint
```

- **Pin** tool versions as Taskfile vars; give each tool task a `status:` that checks the
  pinned version is installed and wired, so `install` is a no-op when nothing changed.
- **Use the tools' own installers** for their config (`rtk init`, `codegraph install`,
  `opencode plugin add`); they back up and merge. Write only inside your own markers.
- **`check`** reports `ok` / `MISS` per piece and is read-only; `verify` exits 1 on gaps.
- One repo per agent keeps each small; share the pinned versions between them.

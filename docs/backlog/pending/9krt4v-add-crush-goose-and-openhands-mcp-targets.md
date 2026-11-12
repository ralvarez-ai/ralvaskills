---
title: 'Add Crush, Goose and OpenHands MCP targets'
priority: p3
section: ralvaskills
topics:
    - rsk
    - tool
---
The Agent Skills clients that read `.agents/skills` and keep MCP servers in a JSON object are covered by `internal/tool/jsonAdapter` (ten targets after the 2026-11-12 work). Three are left out because their MCP surface is not that shape:

- **Crush** — MCP moved to `crushrc`, a Bash script using `mcp add`; `crush.json` still merges but is deprecated and prints a warning. Closest to done: a `jsonAdapter` writing `crush.json` with `{type, command, args}` would work today, at the cost of writing a deprecated format.
- **Goose** — `~/.config/goose/config.yaml` under `extensions:`, so it needs a YAML reader/writer, not `jsonMCPConfig`.
- **OpenHands** — skills are "microagents" in `.openhands/microagents/`, not `SKILL.md` discovery, so the skills half of the Tool interface does not map.

Decide per client in this order: Crush (smallest, decide deprecated-JSON or Bash), Goose (only if used, YAML cost), OpenHands (needs a skills model, lowest value).

Plan: [Agent Skills clients: one adapter, ten targets](../../plans/2026-11-12-agent-skills-client-roadmap.md).

Done when each client is either implemented with a test that registers and unregisters against a realistic config, or recorded as won't-do with the reason.

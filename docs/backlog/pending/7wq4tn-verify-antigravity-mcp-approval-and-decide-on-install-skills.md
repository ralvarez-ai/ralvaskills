---
title: 'Verify Antigravity MCP approval and decide on install_skills'
priority: p3
section: ralvaskills
topics:
    - rsk
    - tool
---
Antigravity registers `rsk mcp --read-only` because its per-tool approval behavior is unverified, the same way Codex and OpenCode are registered. The question is whether Antigravity forces a confirmation on a destructive MCP tool: if it does, `install_skills` can join the Antigravity entry and agents can install skills instead of only inspecting them.

Method: start Antigravity in a project whose `.agents/mcp_config.json` carries a non-read-only `rsk` entry, call `install_skills`, and record what the client shows before the call runs. Compare with the Codex and OpenCode findings in the MCP client support plan, and update the README sentence that groups all three clients under `--read-only`.

Plan: [Google Antigravity as a skill and MCP target](../../plans/2026-11-12-antigravity-support.md).

Done when the approval behavior is recorded for the Antigravity App, IDE and CLI (or named as untested), the README and the plan reflect the outcome, and the adapter's `--read-only` flag is either kept with a reason or dropped with a test.

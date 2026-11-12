---
title: 'Migrate the 30 consumer projects to the renamed skills'
priority: p2
section: ralvaskills
topics:
    - rsk
    - skill-naming
blocked_by: []
related:
  - r4n8kw
  - n6v2wk
---
Step 3 of [ADR 0002](../../adr/0002-rename-ten-skills-with-registry-aliases.md). Measured 2026-11-11: 30 of 111 project manifests under `~/Dev` reference a renamed skill, and five carry more than one.

| Old name | Projects |
|---|---|
| `hexagonal-arch` | 26 — `ralvarez-go/{email,identity,notifier,outbox,ratelimit,webpush}`, `ralvarez-apps/{autobacklog,autoreports,ralvadocket,ralvafi}`, `ralvarez-edu/uns-platform`, `ralvarezdev/{desaparecidos-venezuela,m2p,strucgo,strucgo-runtime}`, `ralvarez-studio/automation-labs-{iec-core,openhands-plc,plc-platform,plc-platform-mcp-apps,plc-platform-mcp-apps-panels,plc-platform-env-overlay,plc-platform-plugins,plc-platform-trace,coding-agents-benchmarks/third-party/openhands-plc}`, `repuestos-edge/repuestos-edge-insiders`, `teamvoltimor/vtitan` |
| `website-concept-architect` | `ralvafi`, `ralvarez.dev`, `repuestos-edge-insiders` |
| `narration-architect`, `explainer-architect` | `ralvarez-edu/ml-lab`, `automation-labs-videos-plc` |
| `demo-script-architect` | `automation-labs-vtscada-testing` |
| `uru-scientific-paper-architect`, `uru-thesis-architect`, `uru-thesis-defense-architect` | `ralvarez-edu/uns-platform` |
| `demo-presentation-architect`, `ml-conference-paper-architect` | none |

Run `rsk update --dry-run` per repo first (it prints one `would rename` line per name and writes nothing), then `rsk update`, then commit the two-line `rsk.mod`/`rsk.lock` diff in that repo. Then `rsk update --global --dry-run` and `rsk update --global` for the four machine-wide links (`demo-script-architect`, `hexagonal-arch` in both `~/.claude/skills` and `~/.config/opencode/skills`).

Done when no `rsk.mod` or `rsk.lock` under `~/Dev` contains a retired name, both global targets link the new names only, and `.rsk/SKILLS.md` in each repo names the new skills.

---
title: 'Land the SKILLS.md imports the newer rsk wrote'
priority: p2
section: ralvaskills
topics:
    - rsk
    - git
related:
  - m4pnr7
---
A newer `rsk` writes the generated block, so `.rsk/CLAUDE.md` becomes the single import `@SKILLS.md` (it was empty, or `@../.claude/skills/docs-architect/SKILL.md`) and `.rsk/SKILLS.md` is re-rendered from the pinned digests. Nothing commits those files, so 22 repos and worktrees under `~/Dev` carry them uncommitted — found 2026-11-12 during the rename migration, and none of them is rename-related (no retired name is in any manifest).

Mostly `.rsk/CLAUDE.md` (18): `ralvadocket.wt/{split-backlog,generic-decisions,generic-ingest}`, `ralvapass.rename`, `ralvasend.rename`, `homelab-nats-apps`, `homelab.rename`, `automation-labs-autoplc-harbor/third-party/harbor-kit`, `automation-labs-agents4plc-harbor/third-party/harbor-kit`, `automation-labs-coding-agents-benchmarks-env-overlay`, `automation-labs-coding-agents-benchmarks-lean-default`, `automation-labs-plc-prompts-exp5`, and five `third-party/` trees under `automation-labs-coding-agents-benchmarks` and `automation-labs-openhands-plc`. `.rsk/SKILLS.md` only: `automation-labs-coding-agents-benchmarks/{third-party/openevolve-plc,release}`. Plus two tracked by other items: `automation-labs-iec-core`'s `rsk.mod` pin (m4pnr7) and `automation-labs-titanic`'s `rsk.lock` (it ignores `.rsk` in git).

Committing an empty-to-`@SKILLS.md` import is safe: the file only points at a generated sibling. Decide per worktree whether `.rsk/SKILLS.md` is tracked there at all.
Done when no repo or worktree under `~/Dev` has an uncommitted `.rsk/CLAUDE.md` or `.rsk/SKILLS.md`, or the ones that stay uncommitted are named and justified here.

Resolved 2026-11-12: eleven trees committed, one commit each, `chore(rsk): render SKILLS.md and import it from CLAUDE.md` — `ralvadocket.wt/{split-backlog,generic-decisions,generic-ingest}` (e3c6473, c08190d, 7830cf4), `ralvapass.rename` (4259de2), `ralvasend.rename` (cc8c365), `homelab.rename` (22a02f6), `automation-labs-coding-agents-benchmarks-{env-overlay,lean-default}` (8d08832e, 5707fe14) and `corekit`'s SKILLS.md-only (14f492c); `homelab-nats-apps` and `automation-labs-plc-prompts-exp5` carried their expanded pin set in the same commit (7449344, 9def054). Staged by path, never `git add -A`; `.rsk/CLAUDE.md`, `.rsk/SKILLS.md`, `.rsk/rsk.lock` (and `rsk.mod` where the pin set grew) only. Nothing pushed.

Ruled out, left uncommitted: the six `third-party/` checkouts are submodules parked on a release tag with a detached HEAD — `harbor-kit` under `automation-labs-autoplc-harbor` and `automation-labs-agents4plc-harbor`, plus `agents4plc-harbor`, `autoplc-harbor`, `openevolve-plc` and `plc-prompts`. A commit there is unreachable and disappears at the next checkout of the pinned tag, and recording it in the superproject would move a vendored pin. Their `.rsk` drift is local tooling state only. `automation-labs-titanic` and `automation-labs-iec-core` are `m4pnr7`'s.

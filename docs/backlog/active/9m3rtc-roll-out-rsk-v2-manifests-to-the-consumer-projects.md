---
title: 'Plan the roll-out of v2 manifests and the git transport to consumers'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#9d3mhw
---
Once the first CLI release is tagged: a project becomes v2 the first time it installs a skill from a non-default source, which includes an official skill, and an older `rsk` then refuses its `rsk.mod` and `rsk.lock`. 26 of 164 lock files on the owner's machine carry an official skill. The consumers pin `rsk` in `mise.toml`, so the pin must move before any of them updates.

**Deliverable.** A release note with the migration rule and the `rsk update` behaviour for official links; the pin bump across the consumers (see the earlier sweep items); a check that no consumer is left on an older pin with a v2 file.

**Done when.** every consumer with an official skill runs a release that reads v2.

Update 2026-11-21 (after step 1f): the release note also needs the new surface: `rsk resolve`, `install --from`, `--trust` and `--save` on install, `rsk sync --trust`, `rsk pin --yes`, the digest confirmation for non-default pins, and that a lock restores through the machine's own transport. A consumer script that pins a remote skill without a terminal now needs `--yes` (the digest is shown, not asked).

Update 2026-11-22: the release is `v3.10.0` (changelog section "v3.10.0", with the migration rule). It is tagged locally and waiting on the owner's push (`c3qwq1`). Start this item once the release assets exist.

Update 2026-11-22 (later): the plan is [plans/2026-11-22-consumer-rollout-rsk-3-10.md](../../plans/2026-11-22-consumer-rollout-rsk-3-10.md). The survey: 149 projects have a lock, 21 use an official skill, none is v2 yet, 85 `mise.toml` files pin rsk (79 at 3.9.1, 6 at 3.7.0) and 62 lock projects have no pin. The release assets do not exist yet, so step 0 (`c3qwq1`) is the gate.

Update 2026-11-22 (release out): `v3.10.0` is released and verified (`c3qwq1` done). The pilot dry run on a copy of `ralvarez-ts/ui` passed (see the plan). Next: the owner approves the real pilot, then the 21 official-skill projects in batches. Use `rsk update`, not `rsk sync`, in consumers: sync also registers the MCP server.

Update 2026-11-22 (end of day): the pilot dry run on a copy passed; the real pilot (`ralvarez-ts/ui`: pin 3.9.1, clean tree, 12 skills, one official; it has uncommitted backlog work of its own that must be left alone) waits on the owner's go. Steps, in `plans/2026-11-22-consumer-rollout-rsk-3-10.md`: bump the pin, `rsk update --yes`, one `chore(rsk)` commit with explicit pathspecs.


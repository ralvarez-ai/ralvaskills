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

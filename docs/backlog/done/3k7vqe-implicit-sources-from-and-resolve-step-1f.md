---
title: 'Implicit sources, --from, rsk resolve and path sources (step 1f)'
priority: p1
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#9d3mhw
---
Step of [the phase 1 plan](../../plans/2026-11-21-source-chain-phase-1.md); design: [private and remote sources](../../designs/2026-11-21-private-and-remote-skill-sources.md).

**Deliverable.** `rsk install @owner/repo/name` for an unregistered source shows the URL and transport and asks; `--yes` never covers it; CI passes `--trust @owner/repo`; nothing is saved without `--save`; SSH first, HTTPS as fallback, and the lock keeps identity and commit so each machine picks its own transport. Pinning a skill from a non-default source shows its `DIGEST.md` text and asks (trust rule 7). `rsk install --from ./path/to/skill` as a one-off, not persisted. `rsk resolve <ref>` explains source, path, version, commit and why a bare name is ambiguous. Declared path sources join the chain with their own walker (the local-clone walker reads `<path>/skills`, a declared path is walked as given, honoring roots).

**Tests.** Implicit source refused without a TTY and with `--yes` alone; `--trust` grants for one run only; nothing written without `--save`; digest confirmation on pin; `--from` leaves no lock record; `rsk resolve` output for bare, qualified, ambiguous and unknown refs; path source discovery, symlink refusal and root confinement.

**Done when.** its tests are green, `task lint`, `go test ./...` and `task registry:check` pass, and registry and local-clone behaviour is byte-identical. This is the first releasable point: with it the trust rules ship with the git transport, and the CLI can be tagged.

Resolved 2026-11-21 (`4513836` declared path sources in the chain, `198a35d` `rsk resolve`, `b03eaea` implicit sources in `internal/sources`, `61a746f` install from an unregistered source and `rsk sync --trust`, `5dd1413` digest confirmation, `cc2fc81` `--from`). The decisions made while building are in the plan, section 1f. Verified with a real binary against `vercel-labs/agent-skills`: refusal without a terminal and with `--yes` alone, the prompt on a pty (no and yes), `--trust` with `--save` (entry written, no `trust.json`), the HTTPS fallback with SSH stubbed out, a `rsk sync --trust` restore on a second sandbox through its own transport, `rsk resolve` for declared-untrusted, unregistered and bare names, and `--from` in a project. The run found one bug the tests missed (the `--from` flag was never registered) and one trap (`--save` leaves bare-name installs failing until the source is trusted), both fixed. Old and new binaries print identical output for the read-only commands on a copy of the real config.


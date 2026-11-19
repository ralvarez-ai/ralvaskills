---
title: 'Add rsk sync and reconcile every consumer project'
priority: p1
section: ralvaskills
topics:
    - rsk
    - cli
related:
    - ralvaskills#m3c8vr
    - ralvaskills#k3vx7p
---
A project's skill links live in machine-local directories and are not committed, so a clone, a cache move, or an aborted install left every consumer with a manifest and a lockfile that no command could turn back into working links. `rsk install` re-resolved versions instead of honoring `rsk.lock`, and `rsk update` is the wrong tool when the goal is to reproduce the locked state.

Done when a command restores the configured clients' links from the lock, refreshes the pinned instructions, leaves `rsk.mod` and `rsk.lock` untouched, refuses a link that points elsewhere or a user-owned path, supports registry, local and official sources, and previews with `--dry-run`. MCP reconciliation must be guarded: register only for a project that already has an rsk server unless asked, preflight every selected client, and require an `rsk` on PATH that supports `mcp --read-only`.

Resolved 2026-11-12 (`7856a5f`, released as v3.9.0): `rsk sync` plans every missing link from the locked cache, applies them, then calls `install.SyncPinnedAllTools`; a second run is a no-op and the manifests stay byte-identical. `--mcp`/`--mcp=false` control registration, and the default only acts when `tool.ProjectHasMCP` finds an existing server. `checkSyncMCPExecutable` runs before any write.

The first consumer sweep exposed two compatibility gaps, fixed in v3.9.1 (`e3cc191`): official skills have no frontmatter version, so their lock entry is versionless and sync rejected 11 projects outright; and OpenCode v2 configs whose `mcp.servers.rsk` was written in the older v1 shape were reported as a conflict. Sync now resolves official entries from the official cache and upgrades both known legacy v1 shapes in place. Tests cover legacy manifests without a `version` key, v1 MCP entries, alternate OpenCode configs, stale links, broken locks, local and official sources, the PATH guard, and idempotence.

Sweep result: all 72 `mise.toml` files under `~/Dev/active` pin 3.9.1, and **71 of 71** projects with `.rsk/rsk.mod` synced (822 links restored across two passes; a final dry run reports 0 pending links and 0 failures). `automation-labs-openevolve-plc-adapter` has no manifest and nothing to sync. `corekit` had a lock behind its links and was repaired with `rsk update --yes`. Three projects had legacy OpenCode MCP entries upgraded to `["rsk","mcp","--read-only"]`.

Left over: consumer working trees keep their `mise.toml` and generated config changes uncommitted; the MCP approval, JSONC and ownership questions remain in [m3c8vr](../active/m3c8vr-add-safe-mcp-client-registration.md).

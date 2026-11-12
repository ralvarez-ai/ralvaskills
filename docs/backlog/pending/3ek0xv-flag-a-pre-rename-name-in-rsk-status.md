---
title: 'Flag a pre-rename name in `rsk status`'
priority: p3
section: ralvaskills
topics:
    - rsk
    - skill-naming
related:
    - 96x2s8
---
The design's second open question: "Should `rsk status` flag an installed old name before any update runs? Cheap and useful; decide when step 1 is written" ([design](../../designs/2026-11-11-skill-rename-aliases.md) §Open questions). Step 1 shipped in `rsk` v3.7.0, the ten renames published 2026-11-12, and all 30 consumer projects were migrated, so no manifest under `~/Dev` carries a retired name today. This is the safety net for the case that remains: a checkout that has not run `rsk update` since the rename, or one restored from an older commit, looks healthy until something resolves a name and prints the redirect.

Why it is not free: `rsk status` makes no registry call — it renders the manifest, the links and the staleness of what is installed — so it cannot see the `aliases` map in `registry/index.json`. Two offline sources already exist: `localAliases` (`cmd/rsk/rename.go:156`) reads `registry/aliases.toml` from a local ralvaskills clone, and `IndexCache` (`internal/mcp/cache.go:59`) serves a cached index for its TTL. Decide which one status may use and what it does when neither is available — silence is fine; a network call from `status` is not.

Done when `rsk status` marks a skill whose name is an alias key, in every place it lists skills (project and global), the mark travels through `--output json` as a field rather than prose, and the command still works with no network and no cache.

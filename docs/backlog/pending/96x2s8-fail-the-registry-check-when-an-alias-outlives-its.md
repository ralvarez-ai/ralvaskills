---
title: 'Fail the registry check when an alias outlives its window'
priority: p2
section: ralvaskills
topics:
    - registry
    - skill-naming
    - rsk
related: []
---
The rename design records the risk and the mitigation in one line: "**Alias rot.** The map grows silently. `task registry:check` should fail when an alias has been in the index for more than two minor releases, so removals are deliberate" ([design](../../designs/2026-11-11-skill-rename-aliases.md) §Risks). Nothing implements it. `registry/aliases.toml` is a plain `[aliases] old = "new"` table, and `applyAliases` (`cmd/generate-registry/main.go:362`) only checks that the target exists in the index, so an alias resolves forever and the map only grows.

Two decisions come before the check:

- **The clock.** Alias support shipped in `rsk` **v3.7.0** (2026-11-12) and the ten aliases landed with it, so "two minor releases" counted from the feature's first release already ran out at 3.9.0, which is released. Either record when each alias was introduced (`since = "3.7.0"` per entry, or a parallel table) and count rsk minors from there, or measure age by the index `generated_at` days. Pick one; do not leave the window implicit.
- **The verdict.** A hard `--check` failure blocks every unrelated publish until the aliases are dealt with, which is the point but also a foot-gun mid-sweep. Warn first, fail one minor later, is the shape that stays usable.

Retiring an alias is not just deleting a key: `registry/aliases.toml` is both the published map and the guard's permanent record, and `CheckRetirements` (`internal/skill/names.go`) must keep rejecting the retired name even after the redirect is gone. So retirement needs the table split — an active alias set that feeds `registry/index.json` and a retired-name set that only the guard reads — or an `expired`/`enabled` marker per entry. That split is the real work; the check is the easy half.

Done when `task registry:check` names an alias past its window, the ten are either expired (redirect dropped from the index, name still refused by `TestCatalogueNames`) or refreshed with the decision written down here, and a test covers an alias that has outlived its window.

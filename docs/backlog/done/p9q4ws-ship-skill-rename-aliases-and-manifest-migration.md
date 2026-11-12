---
title: 'Ship skill-rename aliases and manifest migration in rsk'
priority: p2
section: ralvaskills
topics:
    - rsk
    - registry
related:
  - n6v2wk
---
Renames decided in [ADR 0002](../../adr/0002-rename-ten-skills-with-registry-aliases.md) cannot land until the redirect exists, so this ships first, as a minor release with no skill changed. Design: [skill rename aliases](../../designs/2026-11-11-skill-rename-aliases.md).

Registry: `registry/aliases.toml` as the source of truth, `cmd/generate-registry` emitting the top-level `aliases` map plus `deprecated: true` on the old entry with its versions intact, and `index.json` moved to version 2 with version 1 still accepted.

Client: `aliases` carried through `internal/source` (`IndexEntry`/`registryIndex`, `FindVersion` resolving to the canonical name) and `internal/mcp` (`cachedIndex`, so the MCP tools accept the old name); `skillNamesFromArgs`/`findSkillByName` resolving before lookup; new `internal/manifest/rename.go` with `RenameMod`/`RenameLock`; `skill.Relink`; the migration wired into `bumpLockForProject` and a global relink over `config.GlobalTargets`; `.rsk/SKILLS.md` re-rendered through the existing `refreshPinnedDigest`.

The guard lands with this item: `registry/aliases.toml`, `internal/skill/names.go`, `TestCatalogueNames`, and the `walkSkills` gate in `cmd/generate-registry` (already added — the table is empty until the rename commit).

Landed: the alias map, `FindVersion` resolution, the MCP cache, `RenameMod`/`RenameLock`, the project and global migration in `rsk update`, the generator emitting aliases and deprecation, `registry-remove` preserving them, and the guard. Verified by `internal/manifest/rename_test.go`, `internal/source/registry_test.go`, `cmd/rsk/rename_test.go`, `go test ./...` and a clean `golangci-lint` run. Original done-when: the four test suites in the design pass (`manifest`, `source`, `mcp`, CLI), a temp project with an old name in `rsk.mod` and `rsk.lock` migrates on `rsk update` with `--dry-run` writing nothing, a stale global symlink is re-pointed, and `task lint` plus `task registry:check` are clean.

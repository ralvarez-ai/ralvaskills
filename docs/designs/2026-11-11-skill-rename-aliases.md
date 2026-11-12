---
status: active
project: ralvaskills
date: 2026-11-11
related:
    - adr/0002-rename-ten-skills-with-registry-aliases.md
---

# Skill rename aliases and manifest migration

How a published skill is renamed so that the old name keeps resolving and every consumer's `rsk.mod` / `rsk.lock` ends up on the new name without a hand edit. Decision: [ADR 0002](../adr/0002-rename-ten-skills-with-registry-aliases.md).

## Approach

A skill name is the registry key, so the rename has two halves:

1. **Registry (server):** an `aliases` map plus a `deprecated` old entry, so `old` never 404s and `old@<pinned version>` still installs. This is the GitHub-style redirect.
2. **Client:** `rsk` resolves the alias before lookup, then rewrites the project's manifest, lock, symlinks and rendered digest. The old name is gone from the project on the first `rsk update` after the rename.

Neither half works alone. A map without migration leaves 30 manifests on the old name forever; migration without a map cannot find the old name in the index.

## Published index shape

```json
{
  "version": 2,
  "generated_at": "2026-11-11T00:00:00Z",
  "aliases": {
    "hexagonal-arch": "hexagonal-architecture",
    "narration-architect": "narration-writer"
  },
  "skills": {
    "hexagonal-arch":  { "latest": "1.2.1", "deprecated": true, "versions": { "1.2.1": { "…": "…" } } },
    "hexagonal-architecture": { "latest": "2.0.0", "versions": { "2.0.0": { "…": "…" } } }
  }
}
```

Rules the generator enforces:

- An alias is a **single hop**: the target must exist in `skills` and must not itself be an alias key. Chains are rejected at generation time, not resolved at runtime.
- The old entry is `deprecated: true` with its published versions untouched and `latest` left at the last old release, so an old binary reports "outdated" instead of "not found".
- The map is generated from a new `registry/aliases.toml`, which is the reviewable source of truth; the skill's own frontmatter is not involved because `internal/skill` only parses `version:` and `description:`.

`version` moves to 2 for readers that want to assert the feature; `rsk` accepts 1 and 2, and an older `rsk` simply ignores the two new keys.

## Client resolution

`internal/source/registry.go`:

- `registryIndex` gains `Aliases map[string]string`.
- `Index(ctx)` returns `(skills, aliases, error)` — or a small `Snapshot{Skills, Aliases}` so the MCP cache can reuse it.
- `FindVersion(ctx, name, version)` resolves `name` through the map first and returns the skill under its **canonical** name, so every caller sees the new one.

`internal/mcp/cache.go`: `cachedIndex` gains `Aliases`; the MCP server must resolve too, or `rsk mcp` keeps proposing the old name to an agent while the CLI has moved.

`cmd/rsk/resolve.go`: `skillNamesFromArgs` and `findSkillByName` resolve before lookup, so `install`, `update`, `pin` and `uninstall` all accept either name.

Printed once per resolution, on stderr as a log line:

```
↻ hexagonal-arch is now hexagonal-architecture
```

## Manifest migration

New file `internal/manifest/rename.go`, two pure functions:

```go
type Renamed struct{ Old, New string }

// RenameMod rewrites [skills] keys and the pinned list.
func RenameMod(m Mod, aliases map[string]string) (Mod, []Renamed, bool)

// RenameLock rewrites each [[skills]] name.
func RenameLock(l Lock, aliases map[string]string) (Lock, []Renamed, bool)
```

Call order in `cmd/rsk/update.go` (`bumpLockForProject`), which is already the one place that has the project, the lock and the link targets:

1. `RenameMod` → `WriteMod` when changed.
2. `RenameLock` → `WriteLock` when changed.
3. For each `Renamed`: `skill.Relink(target, old, new)` on every target from `projectSkillsDirs` — write the new symlink, remove the old one, keep the atomic-write discipline of `internal/fsx`.
4. Continue the existing resolution and version bump under the new names.
5. `refreshPinnedDigest` re-renders `.rsk/SKILLS.md` from the new `pinned` list.

`internal/skill/linker.go` gains `Relink(target, old, new string) error`; `Link` stays as it is.

**Global installs** have no manifest: the symlink is the state. `rsk update --global` scans `config.GlobalTargets`, `os.Readlink`s each entry, and re-points any link whose cache path begins with `directory(base)/<old>/`. The cache itself is keyed `<name>/<version>`, so the new version needs a download anyway; the relink happens after the fetch.

Idempotence: after the first run no key matches an alias, both functions return `changed=false`, and nothing is written.

## Guard against resurrection

The retired-name table is the permanent record, so the repository can enforce it from the day the first rename lands:

- `registry/aliases.toml` holds `[aliases] old = "new"`. It is the source of truth for both the published map and the check.
- `internal/skill/names.go` holds the rules: `ValidateName` (kebab-case, ≤64, no doubled hyphen), `CheckRetirements` (exact match **and** any whole hyphen-separated component, so `hexagonal-arch-v2` is rejected too), `ValidateRetirements` (no empty target, no self-reference, no chain), and `CatalogueIssues` (directory equals the frontmatter `name`, name is legal, not retired, not duplicated).
- `TestCatalogueNames` walks `skills/` against the real table, so `go test ./...` — `task test`, and every CI job that runs it — fails the moment a skill takes a retired name. Fixture tests cover the rules themselves, so the guard is tested even while the table is empty.
- `cmd/generate-registry` loads the table first, validates it, and calls `ValidateName` + `CheckRetirements` inside `walkSkills`, so publication fails before a tarball is cut. `--aliases` defaults to `registry/aliases.toml`.

The check reads **skill names only**. A prose mention of a retired name in `docs/`, `CHANGELOG.md` or a skill body is expected and never flagged.

## Rollout

| Step | Ships | Why this order |
|---|---|---|
| 1 | `rsk` alias resolution + migration, one minor release | Without it, the first post-rename `rsk update` cannot resolve the old key and the manifest is stuck |
| 2 | Ten renamed skills, ten MAJOR releases, index regenerated with `aliases` and the deprecated entries | The redirect exists before anything needs it |
| 3 | `rsk update --dry-run` then `rsk update` across the 30 projects; `rsk update --global` | Each project gets a two-line `rsk.mod`/`rsk.lock` diff to commit |

`hexagonal-arch` goes first inside step 2: it owns 26 of the 30 projects and both bundles, so it is the case that exercises the mechanism.

## Risks

- **A missed resolver.** `internal/mcp` is the second lookup path and the easiest to forget; the alias test must cover the MCP tool, not only the CLI.
- **A project with WIP.** The migration writes `rsk.mod` and `rsk.lock`; a repo with uncommitted work gets a mixed diff. `--dry-run` plus `git status` per repo before the real run.
- **A pinned old version.** `hexagonal-arch@1.2.1` keeps installing from the deprecated entry, but `rsk update` will move it to `2.0.0` — a MAJOR jump into a renamed skill. The release notes must say so, and the deprecated entry's `latest` must stay at the old version so an *old* binary never jumps.
- **Alias rot.** The map grows silently. `task registry:check` should fail when an alias has been in the index for more than two minor releases, so removals are deliberate.

## Test plan

- `internal/manifest/rename_test.go`: mod with a skills key + pinned entry; lock with two entries, one aliased; no-op when nothing matches; chain rejected upstream.
- `internal/source/registry_test.go`: fixture index with `aliases` + `deprecated`; `Find` on the old name returns the new name and the new version; `FindVersion(old, pinned)` returns the old archive.
- `internal/mcp/server_test.go`: `search_skills` and `install_skills` accept the old name.
- CLI golden test: a temp project with an old name in both files, `rsk update --dry-run` shows the rename and writes nothing, then the real run leaves both files on the new name and one symlink.
- Global test: a temp dir with a stale symlink, assert the relink and that the old link is gone.

## Open questions

- Does the alias map belong in `registry/index.json` only, or also in a small embedded fallback so `rsk` can migrate offline? Only needed if a rename ever lands without a network round-trip; not now.
- Should `rsk status` flag an installed old name before any update runs? Cheap and useful; decide when step 1 is written.

## Status

Step 1 landed first, with no rename published behind it:

- `registry/index.json` carries `aliases`, and `cmd/generate-registry` emits them, marks the old entry `deprecated`, refuses an alias whose target is not in the index, and bumps the index to version 2 only when the table is non-empty. `cmd/registry-remove` preserves both fields when it rewrites the index.
- `internal/source` resolves the alias in `FindVersion` (the returned skill carries the canonical name); `internal/mcp` caches the alias map and `catalog.FromIndex` skips deprecated entries.
- `internal/manifest.RenameMod` / `RenameLock` rewrite the two files; `cmd/rsk/rename.go` wires project and global migration into `rsk update`, honouring `--dry-run`.
- The retired-name guard ([`registry/aliases.toml`](../../registry/aliases.toml) plus `internal/skill/names.go`) is in place and empty.
- Tests: `internal/manifest/rename_test.go`, `internal/source/registry_test.go` (alias resolves to the canonical name), `cmd/rsk/rename_test.go` (dry run writes nothing; a real run rewrites both files and keeps the version).

Step 2 landed with it in one commit: the ten directories renamed, `registry/aliases.toml` carrying the ten aliases, `registry/index.json` regenerated at version 2 (ten new entries at MAJOR, ten old entries `deprecated` at their last version), the `hexagonal-arch` lines moved in the `global` and `llm-app` bundles, and the repository's ~190 internal references, README and specification updated. `task skills:check` passes because the table and the directories moved together.

Step 3 — the 30 consumer projects and the four global links — is [`m3w7qa`](../backlog/pending/m3w7qa-migrate-the-30-consumer-projects-to-the-renamed-skills.md), and it runs after the ten releases are cut.

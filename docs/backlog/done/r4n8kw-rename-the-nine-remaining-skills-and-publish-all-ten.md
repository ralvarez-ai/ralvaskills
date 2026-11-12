---
title: 'Rename the nine remaining skills and publish all ten'
priority: p2
section: ralvaskills
topics:
    - skill-naming
    - registry
related:
  - n6v2wk
---
Step 2 of [ADR 0002](../../adr/0002-rename-ten-skills-with-registry-aliases.md). Step 1 has landed in the working tree: alias resolution in `internal/source`, `manifest.RenameMod`/`RenameLock`, the migration in `rsk update`, the `registry/aliases.toml` guard, and `cmd/generate-registry` emitting `aliases` + `deprecated`. `hexagonal-arch` is tracked separately in `n6v2wk` and closes in the same commit.

The names and their MAJOR bumps: `narration-architect` → `narration-writer` 2.0.0 (1.1.1), `explainer-architect` → `explainer-designer` 2.0.0 (1.1.0), `website-concept-architect` → `website-concept-designer` 2.0.0 (1.1.0), `demo-presentation-architect` → `demo-deck-designer` 2.0.0 (1.2.0), `demo-script-architect` → `demo-script-designer` 2.0.0 (1.1.0), `ml-conference-paper-architect` → `ml-paper-format` 2.0.0 (1.1.0), `uru-scientific-paper-architect` → `uru-paper-format` 2.0.0 (1.2.0), `uru-thesis-architect` → `uru-thesis-writer` 1.0.0 (0.2.0), `uru-thesis-defense-architect` → `uru-thesis-defense-designer` 2.0.0 (1.1.0), plus `hexagonal-architecture` 2.0.0 (1.2.1).

Per name: `git mv` the directory, bump `version`, update the frontmatter `name`, the H1, `DIGEST.md`, and every cross-skill link (about 190 references across the repo; `grep -rl <old-name> skills docs README.md`). Then: the ten `[aliases]` lines in `registry/aliases.toml` (same commit — the guard fails otherwise, by design), the two `hexagonal-arch` lines in `internal/config/catalog.toml` (`global`, `llm-app`), the stale example at `cmd/rsk/install.go:45` (it names `demo-script-architect`, which was never personal — use `uru-thesis-writer --include personal`), `README.md`, the tree listing and the `-arch` row in `docs/designs/2026-05-20-technical-specification.md`, and a CHANGELOG entry.

Done when `task skills:check`, `go test ./...`, `golangci-lint run ./...`, `task registry:check`, `task registry:generate` + `registry:sync-index`, `task site:check` and `docs_lint.py docs` are clean, the index carries `aliases` with `deprecated: true` on all ten old entries and version 2, and the ten releases are tagged `<new-name>@v<version>`.

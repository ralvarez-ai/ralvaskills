---
title: 'Fix the truncated hexagonal-arch skill name'
priority: p3
section: ralvaskills
topics:
    - skill-naming
    - registry
---
`hexagonal-arch` is the only truncated skill name in the catalogue, and the specification's `-arch` row exists for it alone. Rename it to `hexagonal-architecture` so the name matches the pattern and the spec can drop the one-skill suffix. Decided in [ADR 0002](../../adr/0002-rename-ten-skills-with-registry-aliases.md), which renames ten skills and ships the alias mechanism; this item tracks this name because it owns 26 of the 30 affected projects and both bundles.

Touches: the directory, the `name` frontmatter, the `# Hexagonal Arch` heading, `DIGEST.md`, roughly 20 cross-skill links (`go-architect`, `go-library-builder`, `feature-planner`, `code-design-refactor`, `design-patterns`, `ddd-architect`), `internal/config/catalog.toml` (bundle `global` and `llm-app`), the registry index, `README.md` and the spec tree. Historical records — `CHANGELOG.md`, `docs/reviews/`, `docs/changes/`, `docs/TOKEN_COUNTS.*`, `mockups/` — keep the old name. No alias exists, so `rsk install hexagonal-arch` breaks: release it as `2.0.0` with a CHANGELOG note, and update this repo's own global installs.

Done when `task registry:generate`, `task lint`, `task site:check` and `docs_lint.py` pass, the old name is gone from `registry/index.json`, and the `-arch` row is removed from the spec's Naming Convention table.

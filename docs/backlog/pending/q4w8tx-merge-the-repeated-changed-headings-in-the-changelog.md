---
title: 'Merge the repeated Changed headings in the changelog'
priority: p3
section: ralvaskills
topics:
    - docs
---
`## Unreleased` in `CHANGELOG.md` carries four separate `### Changed` blocks (`### Fixed` and `### Added` sit between them), because each session appended a new heading instead of adding to the existing one. The content is fine; the shape makes the uncommitted change look larger than it is and breaks the Keep a Changelog order the file declares at the top.

Merge them into one `### Added` → `### Changed` → `### Fixed` → `### Removed` sequence for Unreleased only, keeping every entry and its order within each type. `## v3.7.0` and the older versions are history: leave them alone.
Done when `## Unreleased` holds one heading per type and no entry is lost.

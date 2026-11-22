---
title: 'Follow-ups found while building step 1f'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#3k7vqe
---
Found while building and running step 1f; none blocks the first release.

- ~~A versionless skill digests as "no digest, and not readable here".~~ Fixed 2026-11-22: `digest.Load` reads the frontmatter leniently, so the description is kept.
- **The MCP `install_skills` pin shows no digest text.** The CLI shows it and asks (trust rule 7); the tool writes `.rsk/SKILLS.md` for a non-default source after the client's approval of the input only. Decide whether the tool returns the text and refuses to pin until the user confirms it.
- **`rsk update` rewrites `.rsk/SKILLS.md` when a pinned remote skill's digest changes, with no confirmation** (trust rule 7 is only enforced at pin time). Decide whether update should show the new text, and whether the MCP `install_skills` pin must (1h).
- **`rsk status` labels a `--from` skill as `reg`.** `detectSkillSource` falls back to the registry for a link outside every known cache; teach it declared path sources or print `path`.
- **`rsk status --global` orders its tool sections differently between runs** (map iteration; the previous binary does it too), so its output is not byte-stable.
- **`deadcode` reports `git.Clone` as unreachable** since the official cache moved onto the git transport (`997870e`).
- **The landing page counts commands from `setup.go`'s `AddCommand` list,** which leaves out `rsk source` and `rsk resolve`.
- ~~`AGENTS.md` still calls the `rsk mcp` tools "planned" and lists three.~~ Fixed 2026-11-22 (`2a2e35d`).
- **The `rsk-guide` skill does not know `resolve`, `--from`, `--trust`, `--save`, `source scan` and the six source tools yet.** Approved by the owner 2026-11-22; blocked by the permission classifier, tracked as `7bxgw3`.
- **An implicit source cannot carry `roots`, `excludes` or a `rev`.** A repository that needs them (`cursor/plugins` has `pr-review-canvas` twice) is registered with `rsk source add --roots`. `rsk source scan` (1h) should show the clash before a user hits it.

---
status: draft
project: ralvaskills
date: 2026-11-04
related: [designs/2026-05-20-technical-specification.md, guides/registry-flow.md]
---

# Landing page for ralvaskills — work brief

## Scope

**In scope.** Three landing-page directions for `ralvaskills` (the `rsk` CLI + the skill
catalogue), as self-contained HTML proposals, each built around one product-only device, then one
long page that combines the winning ideas and passes the full verification bar.

**Out of scope.** Porting the page into the live site, changing the Pages deploy (it currently
serves `registry/` at `https://skills.ralvarez.dev`), and any edit under `skills/**` or
`registry/index.json` — those auto-publish on push.

## Read first

1. `.claude/skills/landing-page-designer/SKILL.md` — the method. Every section below is that
   skill applied to this repo; where the two disagree, the skill wins.
2. `README.md` and `AGENTS.md` — the product's own words.
3. `docs/designs/2026-05-20-technical-specification.md`, `docs/guides/registry-flow.md`.
4. `internal/` — the source of truth for every claim.

## Output layout

| What | Where |
|---|---|
| Direction HTML (3 in round 1, then the combined page) | `mockups/ui-proposals/2026-11-04/direction-<n>.html` |
| Screenshots, both themes × 1440/390 | `mockups/ui-proposals/2026-11-04/shots/` |
| Fonts (only if used, woff2) | `mockups/ui-proposals/2026-11-04/fonts/` |
| Notes: concept, real-vs-sample, flags, verification table, handoff | `docs/designs/2026-11-04-ralvaskills-landing-page.md` (new, frontmatter per `docs/templates/design.md`) |

`mockups/` does not exist yet — create it. Markdown never goes under `mockups/`; HTML never goes
under `docs/`. Add one row to the `docs/plans/` and `docs/designs/` tables in `docs/README.md`.

## Verified facts (anchors — extend, don't trust)

Checked on 2026-11-04 against this tree. The page may use these without further reading; anything
new goes in the facts table in the design doc, claim → value → source file.

| Claim | Value | Source |
|---|---|---|
| Skills published | **58**, one `SKILL.md` each | `registry/index.json`, `find skills -name SKILL.md` |
| Bundles | **15** (`global`, `docs`, `design`, `go-grpc`, `gin`, `nethttp`, `go-cli`, `fastapi`, `llm-app`, `ros2`, `python-grpc`, `python-cli`, `event-driven`, `observability`, `code-review`) | `internal/config/catalog.toml` |
| CLI subcommands | `init`, `new`, `destroy`, `install`, `uninstall`, `update`, `pin`, `tools`, `catalog`, `list`, `status`, `completion`, `mcp` | `cmd/rsk/setup.go` |
| MCP tools | `project_profile`, `search_skills`, `install_skills` (`install_skills` is the only destructive one) | `internal/mcp/tools.go` |
| MCP resource | `rsk://catalog` (name `rsk-catalog`) | `internal/mcp/resource.go` |
| Machine config | `~/.config/rsk/config.json`; cache `~/.config/rsk/cache/<name>/<version>` | `internal/config/config.go` |
| Manifest | `.rsk/rsk.mod` + `.rsk/rsk.lock`, project-scope only | `internal/manifest` |
| Publishing | push to `main` → release `<skill>@v<version>` → R2 backup → index commit → Pages | `.github/workflows/publish-registry.yml` |
| Page host today | `https://skills.ralvarez.dev` serves `registry/` as the site root (`index.json` only) | same workflow, `CNAME` |

**Drift already spotted — flag in the design doc, do not fix:** `README.md` says "Latest release
`v3.3.0`" while this tree builds `rsk 3.4.1`; the `rsk://catalog` resource description in
`internal/mcp/resource.go` is Spanish while the rest of the CLI surface is English.

## Round 1 — three directions (mandatory)

One device per direction, each different in mechanism, not in colour. Suggested axes; replace one
if a better device appears in the code:

- **A. The manifest as the argument.** Do what `rsk new` does: a repo's `.rsk/rsk.mod` before and
  after, the diff being the product. Device: a real `rsk.mod` file that assembles itself, `.lock`
  resolving beneath it.
- **B. The catalogue as a real number.** 58 skills / 15 bundles as the hero figure, with the
  family tree behind it (domains × skills), filterable by the same keywords `rsk catalog` uses.
  Device: a count that is the page's spine, not a stat strip.
- **C. The agent loop.** `project_profile` detects `go.mod` → proposes `go-architect` → user
  approves → `install_skills`. Device: the three-call exchange shown as the three calls it really
  is, with the approve step visibly a human one.

Each direction is a full page, not a hero section. Five-second test, real assets only, the
trade-off stated: `rsk` is one binary and one config file across every project, and in exchange you
depend on the registry being reachable (or on a local clone).

## Hard requirements

- Self-contained HTML: inline CSS and SVG, **no CDN**, under **200 KB** excluding captures.
- Light and dark through `data-theme` + a toggle (stored preference, system default).
- Landmarks, one `h1`, no heading gaps, skip link, `focus-visible`, 44 px targets.
- A **hero film** per direction: 4–8 s, plays once on view, settles on the plain HTML final state,
  visible Pause/Replay, respects `prefers-reduced-motion`, readable with JS off.
- 390 px designed on purpose (own crops, reflowed diagrams), not just stacked.
- Every number and skill name traceable to a file; sample data labelled on the page.
- No stock photography, no purple gradient, no grey bento grid, no invented testimonials or
  metrics.

## Verification (all must pass before showing anyone)

1. `verify.js` with axe-core at **390×844 and 1440×800 × light/dark × motion/reduced = 8 runs,
   zero violations**; plus no horizontal scroll, images-blocked still readable, JS-off still
   readable. Playwright and chromium are already installed on this machine; put the harness in
   `mockups/ui-proposals/2026-11-04/_tools/`.
2. HTML parser check for swallowed text (a raw `<tag>` or `<a@b.c>` in copy).
3. Smoke test: every control clicked at both widths, fail on `pageerror`, axe re-run in end states.
4. Screenshots: settled film + reduced motion, both themes, both widths, looked at, and at least
   two rounds of changes made because of what the screenshots showed.

State what is **not** covered (screen reader, Firefox/Safari, real devices, user testing of the
five-second claim) in the design doc.

## Boundaries

- Never edit `skills/**` or `registry/index.json`; never commit or push unless the owner asks.
- `landing-page-designer` is already installed in this project (`.claude/skills/`, tracked in
  `.rsk/rsk.mod`, not pinned). Do not install, pin, update or uninstall any skill without an
  explicit yes.
- Do not touch `internal/` or `cmd/` — this is a drawing job.

## Exit criteria

- [ ] Three directions in `mockups/ui-proposals/2026-11-04/`, each a full page with its own device.
- [ ] A contact sheet: all three first viewports side by side, both themes, 1440 and 390.
- [ ] One recommendation, with the reason it beats the other two.
- [ ] Winner combined into a long page that passes verification (8 axe runs clean, parser check
      clean, smoke clean, under 200 KB).
- [ ] `docs/designs/2026-11-04-ralvaskills-landing-page.md` filled: concept, real vs sample, flags
      (including the two drifts above), facts table, verification results, what is not covered,
      and the next ideas for whoever continues.
- [ ] `docs/README.md` index updated.

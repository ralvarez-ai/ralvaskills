---
status: draft
project: ralvaskills
date: 2026-11-04
related: [plans/2026-11-04-landing-page-brief.md, designs/2026-05-20-technical-specification.md, guides/registry-flow.md]
---

# Landing page for ralvaskills — round 1 and the combined page

Proposals only. Nothing here changes `cmd/`, `internal/`, `skills/**`, `registry/index.json` or the
Pages deploy. The method is `.claude/skills/landing-page-designer/SKILL.md`; the brief is
[plans/2026-11-04-landing-page-brief.md](../plans/2026-11-04-landing-page-brief.md).

## Approach

### Files

| What | Where |
|---|---|
| Direction A, the manifest | `mockups/ui-proposals/2026-11-04/direction-1.html` (39 KB) |
| Direction B, the count | `mockups/ui-proposals/2026-11-04/direction-2.html` (53 KB) |
| Direction C, the agent loop | `mockups/ui-proposals/2026-11-04/direction-3.html` (39 KB) |
| Combined long page (C + A + B) | `mockups/ui-proposals/2026-11-04/direction-4-combined.html` (82 KB) |
| Contact sheets | `shots/contact-sheet.png` (A/B/C), `shots/contact-sheet-combined.png` (C vs combined) |
| Screenshots | `shots/<page>-<1440\|390>-<light\|dark>[-reduced\|-full].png` |
| Generator and harness | `_tools/` (`build.py`, `family.py`, `data.py`, `dir_a.py`, `dir_b.py`, `dir_c.py`, `combined.py`, `verify.js`, `smoke.js`, `shots.js`, `contact.js`, `sec.js`, `parsecheck.py`) |

Rebuild: `cd mockups/ui-proposals/2026-11-04/_tools && python3 build.py`. The generator reads
`registry/index.json`, `skills/`, `internal/config/catalog.toml`, `cmd/rsk/setup.go`,
`internal/mcp/tools.go`, `internal/mcp/resource.go`, `.rsk/rsk.mod`/`.lock` and the latest
`v*` git tag at build time, and **asserts** the anchors (58 skills, 15 bundles, 12 `AddCommand` +
cobra `completion` = 13, the three tool names, one mutating tool, `rsk://catalog`). If the tree
drifts, the build fails instead of the page lying. Harness setup: `npm install` in `_tools/`
(Playwright 1.64 + axe-core 4.10; `node_modules/` is git-ignored). It uses the Chromium already in
`~/.cache/ms-playwright` (override with `CHROME=`).

### Concept per direction

All three share one family layer (`family.py`): termkit TokyoNight tokens, the header, kicker,
buttons, film engine, reveals, copy button, footer with the real block wordmark, and a source list
("Where every number on this page comes from"). They differ in mechanism:

- **A — the manifest as the argument.** *"`go.mod`, but for your agent’s skills."* A dark terminal
  runs `rsk new` → `rsk install go-grpc` → `rsk pin go-architect`, and a real-format
  `.rsk/rsk.mod` assembles itself with a diff gutter while `rsk.lock` resolves beneath it. Then: the
  six files `rsk new` really writes, this repo's own `rsk.mod` (night band), the 13 commands grouped
  as `rsk --help` groups them, the trade-off, quick start. Personality: exact, editorial.
- **B — the catalogue as a real number.** A giant **58** over a unit chart: one square per skill,
  grouped by `skills/<domain>/`, the 9 personal ones dashed. The count is the spine: 58 → 15 bundles
  (a live filter that behaves like `rsk catalog`, personal hidden unless `--include personal`) →
  139 releases (weekly histogram from `published_at`) → 13 commands + 3 tools → "the trade-off, in
  ones". Personality: loud, numeric.
- **C — the agent loop.** *"Your agent asks. You approve. rsk installs."* Swimlanes You · Agent ·
  rsk mcp with the three calls as they really are: `project_profile` (read-only, a real response),
  the approval (the one human step, in the warning colour, in the You lane), `install_skills`
  (writes). On phones the lanes reflow into a numbered thread. Then: the tools with their real
  annotations, the 13 detection rules, the real `CLAUDE.md` rule *"Never install without asking."*,
  the trade-off, quick start. Personality: protocol, calm.

### Recommendation: C, and why it beats A and B

**C wins** because it is the only device that shows what nobody else does. A's "go.mod for skills"
is the clearest sentence on any of the pages, but it would describe any package manager for prompts.
B's 58 proves scale, but a count says nothing about why rsk differs. C shows the two things that
are particular to rsk, and both come from the code: the agent proposes skills from file evidence
(`signals.go`, with the proof file in the payload), and the only write is a human decision forced on
every call (`requiresUserInteraction` in `tools.go`, plus the `CLAUDE.md` rule `rsk new` writes).
That is also the trade-off made visible: the agent proposes, you approve. Five-second test: the h1
says what happens and who decides; the lede says the clients (Claude Code, OpenCode) and the
catalogue size; the lanes show the mechanism.

C's weakness is that it assumes the reader knows what an MCP tool is, and says little to someone
who just wants the CLI. The combined page fixes that with A and B.

### The combined page (`direction-4-combined.html`)

One story, one sample project, end to end:

1. **The loop** (C's hero film): profile → you approve → install, for `go-architect`,
   `docker-architect`, `cli-tool-architect`.
2. **What it writes** (A's device, second film): the same three skills landing in `rsk.mod` and
   `rsk.lock`, framed as "after call 3". The ◆ line in the terminal stands for the agent's call.
3. **What `rsk new` writes** (A): the real six-file tree and the real `CLAUDE.md`.
4. **Catalogue** (B): **58** as the heading, the bundle filter with a live count, personal hidden.
5. **Tools** and **Detection** (C): annotations and the 13 signal rules, with the full real JSON.
6. **The real file** (A, night band): this repository's own `rsk.mod` and lock.
7. **Commands** (A), **trade-off** (weighted give/take), **quick start** (copy button), **FAQ** (7
   questions, each answer traced to a file).

Pacing alternates paper / tinted / paper / night; a section rail (01–11, right edge, ≥1400 px)
tracks position. 81.7 KB, two films, both with Pause/Replay.

### Real vs sample

| On the page | Status | How it is labelled |
|---|---|---|
| 58 / 17 domains / 15 bundles / 139 releases / 13 commands / 3 tools | real, computed at build | footer source list |
| `rsk new` terminal lines and the six-file tree, `CLAUDE.md` content | **real run**, 2026-11-04, scratch folder `/tmp/rsk-demo-go` (go.mod, Dockerfile, cmd/demo/main.go) with `RSK_CONFIG_HOME` pointed at an empty temp dir; paths shortened, one sandbox-only warning omitted | "real run" badge |
| `project_profile` response | **real call** to the rsk MCP for the same scratch folder | "real" badge |
| `install_skills` reply, the post-install `rsk.mod`/`rsk.lock`, the pin | **sample**: the shape from `internal/mcp/install_tool.go` / `internal/install/apply.go` with real latest versions; not run, because nothing may be installed in this session | "sample" badge + caption |
| The approval dialog | **drawn**, not a capture of the Claude Code prompt | caption says so |
| This repo's `.rsk/rsk.mod` / `rsk.lock` | real, working tree on 2026-11-04 (includes the uncommitted `landing-page-designer` entry) | "real · working tree" badge |
| Release histogram | real, from `published_at` | caption + table in `<details>` |

### Facts table

| Claim | Value | Source |
|---|---|---|
| Skills published | 58, one `SKILL.md` each | `registry/index.json`; `skills/*/*/SKILL.md` |
| Domains | 17 (academic 2, databases 1, design 2, encoding 1, frameworks 8, frontend 3, infra 5, languages 3, messaging 1, meta 3, personal 9, protocols 4, quality 4, refactoring 4, robotics 1, tooling 2, workflows 5) | `skills/<domain>/` |
| Personal skills | 9, hidden by `rsk catalog` unless `--include personal` | `registry/index.json` `personal`; `cmd/rsk/catalog.go:138-140` |
| Bundles | 15 | `internal/config/catalog.toml` |
| `go-grpc` bundle | go-architect, grpc-architect, protobuf-architect, docker-architect, sql-architect | `internal/config/catalog.toml` |
| Custom bundles | `~/.config/rsk/catalog.toml` replaces same-name bundles, appends new ones | header comment of `internal/config/catalog.toml` |
| Skill releases | 139 since 2026-05-23; peak week 43 (first publish) | `registry/index.json` `versions[].published_at` |
| CLI subcommands | 13: init, new, destroy, install, uninstall, update, pin, tools, catalog, list, status, mcp + cobra `completion` | `cmd/rsk/setup.go`; `bin/rsk --help` |
| Help groups | Machine setup / Project lifecycle / Install… / Project pinning / Tool configuration / Read-only views / Other | `bin/rsk --help` (3.4.1) |
| MCP tools | `project_profile`, `search_skills` (readOnly, idempotent), `install_skills` (destructive, `anthropic/requiresUserInteraction`) | `internal/mcp/tools.go:75-128` |
| MCP resource | `rsk://catalog`, name `rsk-catalog`, `application/json`, descriptions cut at 200 runes | `internal/mcp/resource.go` |
| Detection rules | 13 signal kinds, each with a proof file; refactor skills never proposed by a signal | `internal/mcp/signals.go` |
| What install writes | `m.Skills[name] = "*"` (or constraint), lock entry name/version/source, pin list, then sync pinned to every tool | `internal/install/apply.go:93-121` |
| MCP registration | `.mcp.json` (Claude Code) and `opencode.json` (OpenCode) | `internal/tool/mcp.go:16`; real run |
| Machine config | `~/.config/rsk/config.json` | `internal/config/config.go:20` |
| Version | 3.4.1 | `git tag v3.4.1`; `bin/rsk --version` |
| License | GPL-3.0 | `LICENSE` |
| Quick start | `go install github.com/ralvarezdev/ralvaskills/v3/cmd/rsk@latest`, `rsk init`, `rsk new`, `rsk install <name>`, `rsk pin <name>` | `README.md:139-160` |
| Brand mark | the six-row block wordmark, `◆` glyph, row gradient = accent darkened 0 / .18 / .34 / .48 in Lab, cycled over 6 rows | `internal/ui/theme.go`; termkit `banner.go` |
| Colours | termkit TokyoNight accent `#2E7DE9`/`#7AA2F7`, danger, warning, success, info, muted | termkit v0.65.4 `theme.go` via `internal/ui/theme.go` |
| Trade-off | needs the registry reachable or a local clone (`rsk init` asks for "Path to your local ralvaskills repo clone"); `rsk` must be on PATH for the MCP server | `rsk init --help`; `README.md:182` |

## Components and interfaces

- **Family layer** (`family.py`): tokens for light/dark from termkit plus the Tokyo Night
  backgrounds; contrast of every text token against every background is **asserted ≥ 4.5:1 at
  build**. `data-theme` + toggle (`aria-pressed`, stored in `localStorage['rsk-theme']`, system
  default, circular View Transition wipe when motion is allowed). A head script adds `html.anim` only
  when `prefers-reduced-motion` allows; every hide/animate rule sits under `.anim`.
- **Film engine** (`kFilm`): steps `{at, dur, fn(p)}` on a pausable timeline, `reset()` for the start
  state, `end()` cleanup; the final state **is** the plain HTML. Plays once when 30 % visible, pauses
  when the tab is hidden, Pause/Play/Replay button (hidden without JS). Typed text is driven by the
  timeline (Pause freezes it); the typed layer is `aria-hidden`, the full text sits once in `.sr`.
  Hidden steps keep their box (no layout shift).
- **Reveals**: one-shot IntersectionObserver, opacity + transform, ending at 1. No scroll-linked
  opacity.
- **Catalogue filter**: chips with `aria-pressed`, a live count, an `aria-live` note; non-matching
  tags are de-emphasised with a dashed border and muted text (contrast-safe), never opacity.
- **390 px**: A puts the device between h1 and lede; B shrinks the squares; C reflows lanes into a
  numbered thread with "Agent → rsk mcp" labels; chips scroll horizontally; tables become stacks.

### Flags (found while reading the code — flagged, not fixed)

1. **README version drift** (from the brief): `README.md:136` says "Latest release `v3.3.0`"; the
   latest tag is `v3.4.1` and `bin/rsk --version` prints `rsk 3.4.1`. The pages read the version
   from the tag.
2. **Spanish on the agent surface** (from the brief, extended): the `rsk://catalog` description in
   `internal/mcp/resource.go` is Spanish ("Catálogo de skills disponibles…"), and so are three
   JSON-schema field descriptions the agent sees: `ProjectProfileIn.Path` "Ruta del proyecto; por
   defecto el cwd" and `SearchSkillsIn.Limit` "Por defecto 10, máximo 25" (`internal/mcp/tools.go`),
   `InstallSkillsIn.Scope` "project (por defecto) | global" (`internal/mcp/install_tool.go`). The
   live tool schemas in this session show the Spanish strings.
3. **`llm-app` bundle names two skills that do not exist**: `llm-app-architect` and
   `agent-architect` (`internal/config/catalog.toml:107-108`) are in neither `skills/` nor the
   registry index. On B and the combined page, picking `llm-app` shows a count of 4 and a note naming the two missing skills.
   Not tested: what `rsk install llm-app` does with the missing two.
4. **AGENTS.md vs code on MCP registration**: AGENTS.md says the server is registered in
   `.claude/settings.json` under `mcp`; the code writes `.mcp.json` (`internal/tool/mcp.go:16`) and
   `opencode.json`, as the real run confirms.
5. **`rsk init --force` help text**: "Overwrite existing rsk.mod if present" (`cmd/rsk/setup.go:56`),
   but `init` writes `~/.config/rsk/config.json`, not `rsk.mod`.
6. **`rsk new` without a terminal says yes**: with stdin closed and no `--mcp`/`--guide`, the real run
   registered the MCP server and tried to install `rsk-guide` (skipped only because the sandbox had
   no config). `cmd/rsk/new.go:203-209,259` asks through `termkit.Confirm(cmd, false, …)`; whether
   that is meant to default to yes outside a TTY is for the owner to decide. The README phrase
   "registered by `rsk new`" matches the observed behaviour.
7. **Personal skill names are public**: the 9 personal skills are
   listed in the public `registry/index.json`. The pages show their count and hide the names unless
   you press `--include personal` (mirroring the CLI). Owner call whether the public page should
   name them at all.
8. **Brand gaps**: termkit carries only semantic colours, so backgrounds and body text are upstream
   Tokyo Night (`#1a1b26`, `#16161e`, `#c0caf5`; light `#f3f4f7`, `#e1e2e7`, `#343b58`), not rsk
   code. termkit's light accent `#2E7DE9` is 3.1:1 on paper, so light-theme text uses it darkened 0.34
   in Lab (the wordmark's third step, `#2a5193`); success/warning/danger/info are darkened the same
   way. termkit's `Brand` colour (`#9854F1`, purple) is deliberately unused. rsk has no web font: the
   pages use system sans and mono stacks, so they will look different on macOS and Windows than in
   these Fedora screenshots.

### Verification

Run on 2026-11-04, Chromium (Playwright build 1248), axe-core 4.10.3, tags wcag2a/2aa/21a/21aa/22aa
+ best-practice.

| Check | A | B | C | Combined |
|---|---|---|---|---|
| axe, 390×844 + 1440×800 × light/dark × motion/reduced (8 runs) | 0 violations | 0 | 0 | 0 |
| No horizontal scroll (all 8 runs) | pass | pass | pass | pass |
| Every reveal fired before axe (motion runs) | pass | pass | pass | pass |
| JS off: h1 present, nothing left hidden, no overflow (390 + 1440) | pass | pass | pass | pass |
| Images blocked (no `<img>` is load-bearing; text present) | pass | pass | pass | pass |
| HTML parser check (unknown tags, leaked `{templates}`) | pass | pass | pass | pass |
| Smoke at 1440 + 390: skip link by keyboard, film autoplay → pause → resume → settle → replay, theme ×3 (attribute, `aria-pressed`, storage), copy (label + clipboard), every in-page link, every `<details>`, no uncovered buttons, no `pageerror`/console error, axe + overflow again in the end state | pass | pass (plus filter counts: go-grpc 5, llm-app 4, docs 0, global 15, all 49, all + personal 58) | pass | pass (two films) |
| Size < 200 KB | 39.1 KB | 52.7 KB | 39.0 KB | 81.7 KB |

The harness was checked against a planted page (missing alt, heading skip, low contrast): it
reported all three. It caught one real gap of its own: the first runs scrolled with
`scroll-behavior: smooth` still on, so the below-fold reveals never fired and axe (which skips
`opacity: 0`) never saw that copy. `settle()` now scrolls instantly, and an unfired reveal counts as
a failure.

**Screenshot rounds** (changes made because of what the shots showed):

1. A: h1 at 5.6 rem took four lines and pushed the CTA off the 1440 fold → 4.5 rem; the lock table
   wrapped names and clipped "registry" → source column dropped (all are registry), nowrap; kicker
   wrapped at 390 → shortened.
2. Family: in dark theme the night band was indistinguishable from the page → darker band
   (`#0f1017`) with hairline borders; quick-start comments wrapped mid-line → comment above each
   command.
3. B: at 390 the 58 was timid and the frameworks row wrapped → number min 8 rem, 14 px squares;
   the filter readout competed with the section number → smaller.
4. C: the human step fell below the 1440 fold → requests became labels on the arrows, h1 smaller;
   "You approve." in light was muddy → raw termkit warning `#8C6C3E` for the display line.
5. Combined: A's `grid-area: dev` leaked into the new grid and crushed the "What it writes" column →
   reset; the section rail only appeared ≥1500 px (never tested) → ≥1400 px, then it sat next to the
   film's ①②③ and read as a second numbering → moved to the right edge.

### Not covered

- **Screen readers** (NVDA, VoiceOver, Orca): not run. axe checks structure, not the listening
  experience of the lanes or the live count.
- **Firefox and Safari**: not run. View Transitions (theme wipe) and `color-mix()` fall back or
  differ; `text-wrap` is progressive.
- **Real devices**: only emulated viewports at DPR 1; no touch, no real 390 phone, no zoom/400 %
  reflow test, no Windows High Contrast run (a `forced-colors` block exists but was not seen).
- **The five-second claim**: judged by me, not tested with people.
- **Copy accuracy beyond the facts table**: phrasings such as "makes Claude Code ask you on every
  call" rest on the code comment for `requiresUserInteraction`; OpenCode's behaviour with that meta
  key was not checked.
- **`install_skills` and `rsk install` output**: not run (no installs allowed), so call 3 and the
  post-install files are sample.
- **Fonts**: system stacks only, as rendered on this Fedora machine.

## Risks

- Generated from the tree at build time: re-run `python3 build.py` after the catalogue changes or the
  numbers go stale; a broken anchor fails the build on purpose.
- The scratch-folder outputs (`rsk new`, `project_profile`) are pasted into `_tools/data.py`; they
  will not follow CLI changes until someone re-runs them.
- Hosting: `https://skills.ralvarez.dev` serves `registry/` as the site root today. Putting a page
  there means changing the Pages artifact without breaking `index.json` consumers. That is out of
  scope here and touches the publish workflow.

### Next ideas (handoff)

- Re-run `rsk install` / `install_skills` in a sandbox once the owner allows it, and replace the
  sample call 3 and post-install files with a real capture.
- Capture the real Claude Code approval prompt for the drawn dialog, crop both themes, and label it.
- Decide flags 1–7. Fixing 1, 2, 3 and 5 needs no design input; 6 and 7 are product decisions.
- Try an `rsk status` drift moment (a lock older than the registry) as a second film for the
  updates story.
- Port the combined page into whatever serves `skills.ralvarez.dev` (separate job; keep
  `index.json` at its URL).

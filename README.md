# ralvaskills

**Note:** This is my personal toolkit for daily development, made public for anyone who wants to use, fork, or learn from it.

A growing collection of opinionated, Staff-level skills for Claude Code, Codex, OpenCode, Google Antigravity, Kimi Code, Factory Droid, Gemini CLI, Qwen Code, Cursor CLI, and Sourcegraph Amp.

By default, AI coding agents produce functional but generic code. They nest logic too deeply, ignore architectural patterns, and write verbose commit messages. These skills enforce strict standards — functioning as an automated Senior Engineer that demands clean, maintainable, and highly optimized code.

See [`docs/designs/2026-05-20-technical-specification.md`](docs/designs/2026-05-20-technical-specification.md) for the full technical specification, authoring guide, and roadmap.

---

## Skill Library

Skills are grouped by what they shape. Every folder with a `SKILL.md` is a skill; folders containing other folders are categories.

### `academic/` — papers and study

- **`ml-paper-format`** — ML conference papers (NeurIPS, ICML, ICLR, AAAI): official style-file usage, anonymization and camera-ready toggles, page-limit accounting, natbib citations.
- **`study-artifact-builder`** — Single-file interactive HTML artifact to study a paper's pipeline or dataset until it can be defended from memory: source pinning, paper-versus-code separation, per-component explanations.

### `languages/` — language standards

- **`go-architect`** — Go 1.27 standards: memory-aligned structs, typed enums, interface design, goroutine safety, iterators, idiomatic errors, sqlx + `//go:embed` SQL pattern, `go.work` multi-module layout, and the two app-repo shapes for a Go module with a SPA.
- **`python-architect`** — Python 3.14 enterprise standards: modern typing (PEP 649), immutable dataclasses, Protocol-based DI, asyncio discipline, pytest 9, psycopg + `.sql` files via `importlib.resources`.
- **`latex-architect`** — LaTeX on TeX Live 2026: latexmk-only builds, preamble and package-order discipline, biblatex + biber bibliographies, semantic cross-references (cleveref, siunitx).

### `frameworks/` — framework-specific patterns

- **`fastapi-architect`** — FastAPI 0.136 on Python 3.14: feature layout, Pydantic v2 request/response separation, async DI with lifespan, URL-prefix versioning, RFC 9457 errors, OAuth2+JWT.
- **`gin-architect`** — Gin 1.12 on Go 1.27.3: feature layout, struct-tag validation, RFC 9457 errors, JWT/IdP auth, route groups for versioning, OpenAPI.
- **`nethttp-architect`** — Go stdlib `net/http` (1.22+ ServeMux, no router): feature layout, struct-tag validation, RFC 9457, JWT, graceful shutdown, OpenAPI via kin-openapi.
- **`nextjs-architect`** — Next.js 16 with React 19, App Router only: server components by default with explicit `"use client"` boundaries, server actions, streaming Suspense, edge vs node runtime.
- **`react-architect`** — React 19 + TypeScript strict: feature-based components, hooks-first composition, TanStack Query for server state, zustand for cross-tree client state, Suspense + ErrorBoundary, Radix for a11y.
- **`hugo-architect`** — Hugo 0.161 Extended static sites: project layout, front matter conventions, template hierarchy, Hugo Modules over submodules, Page Bundles, Hugo Pipes asset pipeline, i18n, render hooks.
- **`reveal-js-architect`** — reveal.js 6.0.1 presentations: section/fragment structure, HTML vs Markdown-mode authoring, theming, plugin config, vanilla setup or the `@revealjs/react` extension.
- **`astro-architect`** — Astro 6+: islands architecture with zero JS by default, Content Layer collections with loader + Zod validation, `<script>` processing, `getStaticPaths`, Actions.
- **`wails-architect`** — Wails v3.1 desktop/mobile shell around an existing SPA: embedded build, same-origin API proxy with server-held session cookies (WebKit drops `Set-Cookie` on custom schemes), per-OS window chrome, Android fetch shim and system-bar insets, packaging and signed Play AAB.

### `protocols/` — wire-level conventions

- **`rest-api-architect`** — Cross-language REST conventions: resource URLs, method semantics, URL-prefix versioning, cursor pagination, snake_case JSON, ISO 8601 timestamps, RFC 9457, Idempotency-Key, ETag/If-Match, OpenAPI as source of truth.
- **`grpc-architect`** — Vanilla gRPC: `.proto` services, `status.Error` with standard codes, domain→code mapping, interceptor chain (auth/log/recovery/validation/metrics), client deadlines, context propagation, bufconn testing.
- **`mcp-architect`** — MCP 2026-07-28 servers: tool/resource/prompt primitives, stateless core (no `initialize` handshake), header-based routing, Multi Round-Trip Requests, cacheable list results, OAuth 2.1 + RFC 8707/9207, tool annotations, structured output, prompt-injection/SSRF defenses. Python (FastMCP) and Go SDK recipes. Go schema-driven contracts (JSON Schema as source of truth, generated Go types) and MCP Apps panels as an addition to a complete text result.
- **`connectrpc-architect`** — ConnectRPC: one `.proto` serving Connect, gRPC and gRPC-Web from a single `http.Handler` via `protoc-gen-connect-go`, `connect.Code` error mapping, mandatory interceptor chain.

### `encoding/` — schemas and wire formats

- **`protobuf-architect`** — Protocol Buffers (proto3): Buf-style package naming, protovalidate (CEL), buf toolchain for lint/breaking/generation, field-number reservation, well-known types.

### `databases/` — persistence

- **`sql-architect`** — SQL standards: UUID v7 PKs, snake_case, soft delete, forward-only migrations, parameter binding, N+1 prevention, EXPLAIN-driven indexing. PostgreSQL 18 primary; MySQL 9 / SQLite 3.53 noted.

### `messaging/` — event-driven systems

- **`event-driven-architect`** — Broker-agnostic event taxonomy (NATS/Kafka/RabbitMQ): event/command schemas in Protobuf, topic naming, mandatory outbox, partitioning, idempotency, DLQs, schema evolution.

### `frontend/` — UI and accessibility

- **`ui-ux-architect`** — WCAG 2.2 AA, Radix + Tailwind 4 + shadcn/ui, design-token theming, mandatory loading/error/empty/success states, mobile-first responsive, keyboard parity, contrast checked in CI.
- **`website-concept-designer`** — Pre-code website concepting: purpose, audience, content inventory, tone, visual archetype, structure; outputs 2-3 concept directions plus an ASCII homepage wireframe. Hands off to `ui-ux-architect` / `frontend-design`.
- **`landing-page-designer`** — Landing-page proposals as self-contained HTML grounded in the product's code, real captures and brand: one memorable device per page, a choreographed hero film, accessibility.
- **`mockup-architect`** — Mockups as self-contained HTML in dated folders, with the fixed preview bar (palette, corners, accent, mark) driven by data attributes and `?query` so variants are reproducible.
- **`demo-deck-designer`** — Slide-deck specifications in markdown: content distribution across slides, within-slide organization, layout from a fixed catalog. Outputs `.md` only.
- **`demo-script-designer`** — Presenter-centric demo scripts with narrative flow, visual-cue mapping, progressive capability reveals.

### `design/` — architectural patterns

- **`ddd-architect`** — Domain-Driven Design: strategic-first (bounded contexts, context mapping) with aggregates, value objects, domain events, repositories. Event sourcing out of scope.
- **`hexagonal-architecture`** — Ports & adapters: dependencies point inward; domain declares ports, adapters implement them; domain never imports framework/DB/HTTP.

### `infra/` — runtime, packaging, observability

- **`docker-architect`** — Multi-stage builds, per-language base defaults (distroless Go, slim Python/Node), BuildKit cache mounts, non-root, multi-arch amd64+arm64, digest-pinned bases, Trivy scanning, Compose v2.
- **`ci-cd-architect`** — CI/CD principles: pipeline taxonomy, trigger design, SHA-pinned supply-chain hygiene, OIDC cloud auth, caching, matrix strategy, test gates, release automation. Suggestion-mode; GitHub Actions recipes.
- **`observability-architect`** — Application-side signal production: structured logs, Prometheus metrics, OTel traces, signal correlation, head sampling, PII discipline, RED+USE.
- **`grafana-architect`** — Signal consumption: dashboards-as-code (Grizzly), per-service folders, one-question-per-panel, unified alerting with runbooks, low-cardinality discipline.
- **`secrets-architect`** — Reference-based secret handling for one external store (`$SECRETS_DIR`): consumer decision matrix (Compose secrets, `*_FILE`, systemd `LoadCredential`, writable bind mount, CI), `700`/`600` + uid rules, per-secret migration workflow, live verification, pitfalls.

### `robotics/` — robotics platforms

- **`ros2-architect`** — ROS2 across Jazzy/Kilted/Lyrical: colcon + ament_cmake/ament_python, lifecycle nodes, explicit QoS, services/actions, parameters, Python launch DSL. C++ (`rclcpp`) and Python (`rclpy`) equal first-class. Cross-platform dev env via Pixi + RoboStack.

### `quality/` — reviewers

- **`api-contract-reviewer`** — Reviews REST + gRPC contracts for stability, versioning, completeness, backwards compatibility. Runs `buf breaking` / `openapi-diff`. Severity-keyed findings.
- **`performance-reviewer`** — Cross-language perf review: N+1, missing indexes, blocking I/O in async, allocation hot paths, unbounded memory. Findings grounded in EXPLAIN / pprof / py-spy / metrics.
- **`security-reviewer`** — Cross-language security review: injection, auth/authz, secrets, insecure defaults, deserialization, CSRF/SSRF/IDOR, dep vulns. Critical/High/Medium/Low with file:line + fixes.
- **`codebase-hygiene-reviewer`** — Evidence-based declutter audit: dead code, unused exports, duplication, oversized files/functions, modularization candidates. Proposes, never edits.

### `tooling/` — developer tooling

- **`cli-tool-architect`** — Cross-language CLI standards: subcommand structure, flag/env/config/default precedence, TOML in XDG, stdout-data/stderr-logs split, `--output json|yaml`, exit codes, NO_COLOR, completions.
- **`repo-tooling-architect`** — Repo-root productivity files: `.editorconfig`, `.gitignore`, version pinning (mise), task runner (Task), minimal pre-commit, dotenv + external secret manager, Renovate.

### `workflows/` — process

- **`feature-planner`** — Upstream feature planning: requirement clarification, constraint discovery, design outline, vertical-slice task breakdown ordered by risk. Hands off to `grill-with-docs` and `tdd`.
- **`grill-with-docs`** — Stress-tests plans against the project's domain model, sharpens terminology, updates CONTEXT.md and ADRs inline as decisions crystallise.
- **`tdd`** — Red-green-refactor loop emphasizing behavior verification through public interfaces.
- **`commit-author`** — Generates concise Conventional Commits messages from a staged diff. Enforces full type set, imperative subject lines, no AI co-author attribution.
- **`docs-architect`** — Project documentation system: a folder-per-doc-type `docs/` tree (ADRs, backlog, ideas, plans, designs, reviews, changes, guides), mandatory frontmatter, lifecycle by moving files, backlog as a queue or one file per item, a stdlib linter. Scaffold, write-doc and migrate modes.

### `refactoring/` — code cleanup ladder

Three rungs from expression-level polish to system-level architecture, plus a skeptical pattern catalog. Each rung explicitly defers up or down when scope changes.

- **`logic-cleaner`** — Expression-level cleanup: guard clauses, naming, magic values, conditional simplification, duplication. Language-agnostic.
- **`code-design-refactor`** — Design-level refactoring: extraction, decoupling, SRP, encapsulation, primitive obsession. Sits between `logic-cleaner` and `improve-codebase-architecture`.
- **`improve-codebase-architecture`** — System-level architectural friction analysis with deepening opportunities informed by CONTEXT.md and ADRs.
- **`design-patterns`** — Skeptical reference catalog for modern Go and Python: keeps repository, adapter, strategy, decorator, observer, builder; flags the rest as anti-patterns. Consult before applying a pattern in the other refactoring skills.

### `meta/` — ralvaskills ecosystem

Skills about authoring/operating the toolkit itself, plus communication modes that aren't really workflows.

- **`skill-builder`** — Meta-skill that scaffolds new ralvaskills skills per SPECS.md (SKILL.md skeleton, optional STACK/RECIPES, folder placement, SPECS updates). Interview-first.
- **`rsk-guide`** — Operator's guide for the `rsk` CLI: the `.rsk/` project manifest (`rsk new`, `rsk install`, `rsk pin`), bundle/global installs, the official Anthropic cache, and the `rsk mcp` server the agent queries for project standards.
- **`caveman`** — Ultra-compressed communication mode that reduces token usage by ~75% while preserving technical accuracy.
- **`adhd-mode`** — Reply shape for a reader with ADHD: next action first, numbered steps, restated state, one next step.
- **`token-economy`** — Token-saving habits and setup for Claude Code and OpenCode: always-on rules, subagent exploration, model escalation and handoff, cheap-model guidance, and an audit of context-mode, RTK and CodeGraph.

### `personal/` — author-specific, not bundle-installable

These are personal to me and excluded from bundle installs. Listed here for transparency; not intended for reuse.

- **`pi-iteration-workflow`** — Windows-to-Pi SSH edit-test loop for the `voldemorbot` robot codebase: pixi/ROS2 rebuild triggers, `screen` for flaky-WiFi survival, PowerShell SSH quoting pitfalls, Pi Zero sudo caching.
- **`uru-thesis-writer`** — Interview-first, section-by-section drafter for URU thesis (TEG) chapters, Spanish output, URU chapter order, never fabricates data. Complements `uru-thesis-reviewer`.
- **`uru-thesis-reviewer`** — Continuous-feedback reviewer for URU (Universidad Rafael Urdaneta) Trabajos Especiales de Grado. Produces ordered `.md` diff files; never edits the source `.docx`.
- **`uru-thesis-defense-designer`** — URU-style thesis defense deck specs: fixed institutional slide structure, Spanish-only content, logo placement rules. Outputs `.md` only.
- **`uru-paper-format`** — Manuscript formatting for Revista Tecnocientífica URU: typography rules, IMRyD-derived structure by article type, citation/reference formats. Outputs one `.md` manuscript spec.
- **`go-library-builder`** — Scaffolds reusable Go 1.27.3 library modules in the ralvarezdev house style (flat, vertical-slice, toolkit or framework-adapter shapes): domain-pure root package, narrow ports, technology-named adapter subpackages (prod/dev/test), config structs over options, testcontainers over mocks, mise + Task + golangci tooling. Module namespace is a parameter.
- **`go-house-libs`** — Maps a need to the ralvarezdev shared Go modules (restkit, ginkit, pgkit, mcpkit, contracts, connectrpckit, identity, identitygin, ratelimit, email, webpush, termkit, svckit, tick, sqlitekit, resilience): a router plus one `references/` file per module with its package layout, a recipe and the gotchas; the consumer-side companion of `go-library-builder`.

---

## `rsk` CLI

A Go CLI that manages skill installation across projects without copy-pasting folders. Tracked in [`docs/designs/2026-05-20-technical-specification.md`](docs/designs/2026-05-20-technical-specification.md). Published builds are on the [releases page](https://github.com/ralvarezdev/ralvaskills/releases).

```
go install go.ralva.dev/ralvaskills/v3/cmd/rsk@latest
# or grab a prebuilt binary from the Releases page

# Machine setup (once)
rsk init                                  # configure source, tools, global dirs

# Per-project workflow
rsk new                                   # initialize .rsk/ in this project
rsk install <name>                        # install a bundle or skill (auto-resolved)
rsk install <name> --pin                  # install + auto-load in this project
rsk install                               # re-install from rsk.mod (after git clone)
rsk sync --dry-run                        # preview reconciliation using rsk.lock
rsk sync                                  # restore links and sync configured clients

# Global / machine-wide
rsk install <name> --global               # link into configured tool directories
rsk install <name> --global --for claude-code

# Browse & inspect
rsk catalog                               # what's available (skills + bundles)
rsk list                                  # what's installed in this project
rsk list --global                         # what's installed system-wide
rsk status                                # combined view + drift

# Updates and removal
rsk update [name]                         # pull latest and re-link
rsk uninstall <name>                      # remove (cleans manifest in project scope)
rsk pin <name> [--remove]                # toggle auto-load in CLAUDE.md / opencode.json (`rsk unpin` still works)
rsk destroy                               # remove .rsk/ and tool config entries (--yes skips the prompt)

# MCP server (registered by `rsk new`)
rsk mcp                                   # stdio MCP server the agent queries (launched by the client)

# Tool permissions (Claude Code by default; --for selects the tool)
rsk tools list                            # current permissions for this project
rsk tools allow <rule>                    # allow a tool
rsk tools deny <rule>                     # deny a tool
rsk tools remove <rule>                   # drop a rule
# (the older `rsk claude tools ...` path still works as a hidden alias)
```

### MCP server

`rsk new --mcp` registers `rsk mcp` in every selected client's project config, preserving unrelated servers and keys:

| Client (`--for`) | Project MCP config |
| --- | --- |
| `claude-code` | `.mcp.json` (plus a short pointer in `./CLAUDE.md`) |
| `codex` | `.codex/config.toml` |
| `opencode` | `opencode.json` |
| `antigravity` | `.agents/mcp_config.json` |
| `kimi` | `.kimi-code/mcp.json` |
| `droid` | `.factory/mcp.json` |
| `gemini` | `.gemini/settings.json` |
| `qwen` | `.qwen/settings.json` |
| `cursor` | `.cursor/mcp.json` |
| `amp` | `.amp/settings.json` |

The server provides `project_profile`, `search_skills`, and the `rsk://catalog` resource over stdio. Claude Code additionally receives `install_skills` with Claude-specific interaction metadata; every other client is registered read-only (`--read-only`), which omits that tool until its per-install approval behavior is verified. Codex loads project configuration only in trusted projects.

The client launches the server, so **`rsk` must be on PATH** on every machine that opens the project — `rsk new` warns and skips registration when it is not. Re-run `rsk new` to (re)register; it never overwrites other MCP servers. `rsk destroy` removes the registration and the pointer.

### Interactive mode

Running `rsk` with no arguments and a real terminal opens the interactive session instead of printing help: a flat picker over the command tree (type `/` to filter, `↑`/`↓` to move, `?` for the help overlay), where any command that takes flags or positional arguments shows a typed form and the run's table, chart, or raw output opens as a paged result. Result rows offer their own actions (`i` install, `u` uninstall, `a`/`d`/`x` allow/deny/remove), so you can act on a row without leaving it. Recently used commands float to the top, and the theme, history, and last-used form values are remembered between runs. Non-interactive contexts (pipes, CI) fall back to the plain `--help` output.

---

## Recommended Official Claude Code Skills

These official Claude Code skills complement the ones above and are highly recommended:

- **`frontend-design`** — Distinctive, production-grade frontend interfaces with high design quality.
- **`xlsx`** — Spreadsheet files (`.xlsx`, `.xlsm`, `.csv`, `.tsv`) for data manipulation, analysis, formatting, conversion.

---

## Attribution

Some skills are adapted from other people's work and changed for this toolkit. Each one carries a credit footer; license notices are in [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

| Skill | Source | License |
| --- | --- | --- |
| `tdd`, `grill-with-docs`, `improve-codebase-architecture` | [Matt Pocock's skills](https://github.com/mattpocock/skills) | MIT |
| `caveman` | [JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman) by Julius Brussee | Apache-2.0 |
| `adhd-mode` | [ayghri/i-have-adhd](https://github.com/ayghri/i-have-adhd) by Ayoub Ghriss | MIT |

---

## Installation

These skills are drop-in modules for ten agent clients: Claude Code, Codex, OpenCode, Google Antigravity, Kimi Code, Factory Droid, Gemini CLI, Qwen Code, Cursor CLI, and Sourcegraph Amp.

### Global installation

Clone this repository and copy the skill folders to your global skills directory.

**Windows (PowerShell):**

```powershell
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.config\opencode\skills\"
Copy-Item -Path .\skills\* -Destination "$env:USERPROFILE\.config\opencode\skills\" -Recurse -Force

New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.claude\skills\"
Copy-Item -Path .\skills\* -Destination "$env:USERPROFILE\.claude\skills\" -Recurse -Force
```

**macOS / Linux:**

```bash
mkdir -p ~/.config/opencode/skills/
cp -R ./skills/* ~/.config/opencode/skills/

mkdir -p ~/.claude/skills/
cp -R ./skills/* ~/.claude/skills/
```

### Per-project installation

Drop individual skill folders into your project's `.claude/skills/` (Claude Code), `.agents/skills/` (Codex, Antigravity, Kimi, Droid, Gemini, Qwen, Cursor, Amp), or `.opencode/skills/` (OpenCode) directory. Each skill is self-contained — copy only what you need.

`.agents/skills/` is the Agent Skills convention every client but Claude Code and OpenCode reads, so `rsk new --for all` links each skill once even with eight tools selected. To enable all of them in an existing rsk project, run `rsk new --for all`, then `rsk install` to link its existing skills into each native directory. Review old `.claude/skills/` links before removing any: Claude may still use them, and OpenCode also discovers them for compatibility. Run `rsk init` to configure each client's global target — `~/.agents/skills`, `~/.gemini/config/skills`, `~/.kimi-code/skills`, `~/.factory/skills`, `~/.qwen/skills`, `~/.cursor/skills` or `~/.config/agents/skills`. Every one of them discovers installed skills on demand; rsk's pinned digests are not automatically loaded.

`rsk sync` uses the exact cached registry versions in `rsk.lock` and preserves
`rsk.mod` and `rsk.lock`. It refuses conflicting existing links. It refreshes
pinned instructions and reconciles MCP registration when an rsk server is
already configured; use `--mcp` to register one, or `--mcp=false` to leave MCP
configuration alone. The `rsk` binary on PATH must support read-only MCP mode
before Codex, OpenCode or Antigravity MCP registrations can be synchronized.

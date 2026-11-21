---
status: draft
project: ralvaskills
date: 2026-11-21
related:
    - designs/2026-05-20-technical-specification.md
    - designs/2026-11-11-skill-rename-aliases.md
    - plans/2026-10-02-mcp-server-plan.md
    - plans/2026-11-17-per-skill-permissions.md
    - guides/registry-flow.md
---

# Private and remote skill sources

A design for installing skills from more than the two sources `rsk` supports today, including
private repositories, arbitrary git refs, direct paths and URLs, with reproducible locks and no
secrets stored by `rsk`. Written to be evaluated outside this repository: it states the problem,
the constraints the current code imposes, the options, one proposed design, and the open questions.

## Context and problem

`rsk` resolves a skill name through exactly two sources: a local clone **or** the hosted registry,
plus the `anthropics/skills` cache.

- `install.LocalSource` returns either `source.NewLocal(cfg.RepoPath)` or
  `source.NewRegistry(...)`, never both (`internal/install/resolve.go:26-31`, gated by
  `config.LocalMode`, `internal/config/config.go:50`).
- The second source is always the official cache (`internal/source/official.go:16`, URL hardcoded).
- `skill.Source` is a closed enum: `local`, `official`, `registry`
  (`internal/skill/source.go:8-28`). `ParseSource` rejects anything else, including on lock-file
  reads.
- Resolution order is written out in five places: `install.FindSkill`
  (`internal/install/resolve.go:33-56`), `LocalSource` (`:26`), `newLocalSource`
  (`cmd/rsk/resolve.go:19`), `loadSkillsBySource` (`cmd/rsk/catalog.go:167`) and the MCP server.
- The registry client has no authentication: `NewRegistry` builds a bare `http.Client`
  (`internal/source/registry.go:73`). There is no header, token or credential concept anywhere.
- The lock records name, version and source kind only (`internal/manifest/lock.go:22-34`), keyed by
  bare name (`UpsertLockEntry`, `:78`). A remote skill pinned to a git ref is therefore not
  reproducible, and two sources with the same skill name collide silently.
- `personal/` is an install filter, not a privacy boundary: those skills live in the public repo,
  appear in the public `registry/index.json` with `"personal": true`, and ship as public release
  tarballs. `--include personal` (`cmd/rsk/install.go`) and the MCP refusal
  (`internal/mcp/install_tool.go:145`) only discourage use.

Consequences:

1. A project cannot use the public registry **and** a private skill repo at the same time.
2. A skill that exists only in a git repository cannot be installed at all.
3. A private registry cannot be read by `rsk`.
4. Nothing records where a remote skill came from, so `rsk update` and `rsk status` cannot verify it.
5. Publishing is welded to this repository: `cmd/generate-registry` hardcodes GitHub release asset
   URLs (`cmd/generate-registry/main.go:285`) and this repo's `skills/` layout.

## Goals and non-goals

Goals:

- **G1** A project or machine can declare N ordered sources of mixed kinds: hosted registry, private
  registry, git repository, local path, official cache.
- **G2** Remote install by git ref, by tarball URL, and by path, without publishing anything.
- **G3** Private sources authenticate without `rsk` storing a secret: it stores a *reference* to an
  environment variable, keychain entry, or the user's git credential helper.
- **G4** Every install is reproducible: the lock records the transport identity (git commit or
  tarball SHA-256) and the source id.
- **G5** Repo-wide access is sufficient to onboard a private skills repository; no per-skill ACL.
- **G6** Arbitrary nesting inside a source works deterministically: same name means same skill,
  duplicates are an error.
- **G7** Anyone can produce a compatible index from their own repo without forking this repo's
  workflow.

Non-goals (deliberately deferred):

- Per-skill entitlements, metering or selling skills.
- OCI artifacts as a transport.
- Encryption of skills at rest (plaintext ends up in the agent's context anyway).
- Signature verification of skill content (see Open questions).
- Windows symlink strategy (unchanged from today).

## Constraints imposed by the current code

| # | Constraint | Where | Effect on this design |
|---|---|---|---|
| 1 | One local **XOR** one registry | `internal/install/resolve.go:26-31`, `internal/config/config.go:50` | Source chain is mandatory; this is the primary change |
| 2 | Closed `Source` enum | `internal/skill/source.go:8-28`, `Label()` `:37-48` | Replace kind-based enum with source ids; keep wire compatibility |
| 3 | Hardcoded order in five call sites | `install.FindSkill`, `LocalSource`, `newLocalSource`, `loadSkillsBySource`, MCP | One `Chain` and one factory, consumed everywhere |
| 4 | Lock has no provenance | `internal/manifest/lock.go:22-43` | Lock schema bump to v2 |
| 5 | Lock keyed by bare name | `internal/manifest/lock.go:78` | Qualified `source:name` keys |
| 6 | No HTTP auth | `internal/source/registry.go:73` | Auth transport layer |
| 7 | No integrity metadata | `internal/source`, `cmd/generate-registry` | SHA-256 in index and lock |
| 8 | Publishing tied to this repo | `cmd/generate-registry/main.go:285` | Extract a buildable library + CLI |
| 9 | `personal/` is path-based and public | `internal/skill/util.go:9` | Split `personal` (published) from `private` (never published) |
| 10 | Filesystem lookup re-walks the root | `internal/source/fs.go` | Cache per process; matters for monorepos |
| 11 | `rsk.mod` is `map[string]string` | `internal/manifest/mod.go:19-26` | Schema bump to v2 for qualified names + `[[source]]` |

## Options

| Option | What changes | Pros | Cons | Effort |
|---|---|---|---|---|
| **A. Private registry only** | Auth on `Registry`, second URL | One new transport; discovery and versions stay first-class; fits multi-team | Requires a bucket, an auth front, a publisher, token distribution; still no git-ref installs | M |
| **B. Source chain, git-first** | Source spec + factory + chain; git transport; lock v2 | Fixes XOR blocker; private access is git's problem; no new infra; highest reuse | Git-only versioning (commit ids, frontmatter version is display); needs drift handling | M |
| **C. Vendor / submodules only** | None in `rsk` | Zero code | No resolution, no updates, no separation of public/private | S |
| **D. OCI artifacts** | New transport | Nice provenance, registries already exist | Tooling and auth cost, no skill-specific benefit yet | L |

**Proposed decision:** B. Treat A as one more transport in the same chain, so choosing git first does
not close the door on a private registry later. C remains possible as a degenerate path source.

Accepting this design should produce an ADR (next free number: 0003) covering the source-chain
decision and the lock v2 bump.

## Design

### 1. Source model

A source is a named, typed descriptor. It is declared in the machine config and may be overridden or
added per project.

```go
// internal/source/spec.go (new)
type Kind string // "registry" | "git" | "path" | "official"

type Spec struct {
	ID       string
	Type     Kind
	URL      string   // registry base URL, git URL, or filesystem path
	Ref      string   // git ref; empty means the source's default branch
	Roots    []string // directories to walk inside the source; empty means "."
	Excludes []string // glob patterns pruned before walking
	Priority int      // lower wins; declaration order breaks ties
	Private  bool     // never listed outside this machine/project
	Auth     string   // "none" | "git" | "env:VAR" | "keychain:service"
}
```

```go
// internal/source/chain.go (new)
type Resolver interface {
	ID() string
	All(ctx context.Context) ([]skill.Skill, error)
	Find(ctx context.Context, name string) (skill.Skill, error)
	Spec() Spec
}

type Chain []Resolver

func New(spec Spec, opts Options) (Resolver, error) // the only factory
func Build(specs []Spec, opts Options) (Chain, error) // sorted by priority
func (c Chain) Find(ctx context.Context, ref string) (skill.Skill, error)
func (c Chain) All(ctx context.Context) ([]skill.Skill, error)
```

Resolution rules:

1. `source:name` resolves in exactly that source; an unknown source is an error.
2. A bare `name` is resolved in priority order; the first hit wins only if it is unique.
3. If two sources provide the same bare name, that is an **error** that names both sources.
   `rsk install --from <source>` disambiguates.
4. Ordering is `priority` ascending, then declaration order, then id.
5. `private = true` sources are excluded from `catalog` listings, MCP resources and
   `search_skills` unless the command is run with `--include private` (mirroring today's
   `--include personal`).

Existing commands keep their flags; `newLocalSource`/`LocalSource` become thin wrappers over
`Build(cfg.Sources, ...)` so the five call sites collapse to one.

### 2. Authentication and authorization

Rule: **`rsk` never persists a secret.** `Auth` is a reference; the value is read at use time and
never logged. Tokens go in headers, never URLs (upstream `http.Client` errors include the URL).

| Transport | Mechanism | Granularity | Secret lives in | Revocation | Fit |
|---|---|---|---|---|---|
| git | SSH key / deploy key (read-only) | per repo | `~/.ssh`, agent | remove key | default for a developer |
| git | SSH cert from an org CA | repos the CA allows | agent | CA policy | orgs with many repos |
| git | HTTPS fine-grained PAT | selected repos | git credential helper / keychain | expire or rotate | headless machines |
| git | GitHub App installation token | installation repos | minted, 1 h TTL | uninstall app | fleets; best CI story |
| git | OIDC → cloud STS | one workflow run | none | n/a | publish/build only |
| registry | `Authorization: Bearer` | one source | `env:` or keychain | rotate token | simplest private registry |
| registry | Basic auth | one source | keychain | rotate | legacy hosts only |
| registry | mTLS | one source or org | OS store | remove cert | internal infrastructure |
| registry | Cloudflare Access service token / JWT | one application | env / keychain | rotate | fits the current CF-fronted host |
| registry | Presigned object URLs (S3/R2 SigV4) | one object, short TTL | none on disk | expiry | tarball downloads |

Authorization granularity ladder — pick one per deployment:

| Level | How | Supported by |
|---|---|---|
| Source-wide (recommended) | one repo or one registry = one credential = one policy | git and registry |
| Team | one source per team, same mechanics | git and registry |
| Subtree | separate repos; a client that can clone can read every directory | requires repo split |
| Per-skill entitlement | server filters the index per principal | registry only |

Read and publish credentials are separate. Developers hold read only; publishing runs in CI with an
App token or OIDC and never on a laptop.

### 3. Repo-wide setup

Repo-wide is the supported and recommended boundary: one skills repository, one credential, every
skill in it shares the policy. Privacy class — not skill domain — decides the repo split.

```toml
# .rsk/rsk.mod (project) or ~/.config/rsk/config.json (machine)
[[source]]
id       = "acme-private"
type     = "git"
url      = "git@github.com:acme/skills-private.git"
ref      = "main"
roots    = ["skills"]
priority = 10
private  = true
```

Onboarding a second machine is: clone access plus the same two lines. No secret is copied, because
git supplies it.

### 4. Schemas

#### 4.1 Source descriptor

As above, in `config.json` (`sources`) and `rsk.mod` (`[[source]]`). Project sources override machine
sources with the same `id`; `RSK_SOURCE_<ID>_TOKEN` overrides an `env:` auth reference.

#### 4.2 `rsk.mod` v2

```toml
version = 2
tools   = ["claude-code", "opencode"]
pinned  = ["acme-private:go-architect"]

[skills]
"acme-private:go-architect" = "*"
"go-grpc"                   = "^1.2.0"   # bare = unambiguous across sources
```

`schema.Mod` bumps to 2. Version 1 files are read as before (all names bare, implicit sources from
`repo_path`/`registry_url`) and rewritten as v2 on the next write.

#### 4.3 `rsk.lock` v2

```toml
version = 2

[[skills]]
name      = "acme-private:go-architect"
source    = "acme-private"
transport = "git"
url       = "git@github.com:acme/skills-private.git"
ref       = "main"
commit    = "9f3c1ab7e2d4c1e0f8b6a5d4c3b2a1908f7e6d5c"
version   = "1.2.0"            # frontmatter, display only
permissions = ["Bash(make lint)"]
granted     = []

[[skills]]
name      = "go-grpc"
source    = "ralvaskills"
transport = "registry"
url       = "https://skills.ralvarez.dev"
version   = "1.4.0"
sha256    = "3b1f…"
```

`LockEntry` gains `Source`, `Transport`, `URL`, `Ref`, `Commit`, `SHA256`; `Name` becomes qualified.
For git, the commit id is the integrity record. For registry and URL transports, `sha256` is the
verification record and is required from index v2 onward.

#### 4.4 Hosting manifest: `rsk-skills.toml`

Every repository that carries skills can declare what `rsk` should see. Optional for a path source,
recommended for git sources, required for monorepos.

```toml
# rsk-skills.toml at the repository root
private  = true
roots    = ["skills", "packages/*/skills"]
excludes = ["**/node_modules/**", "**/testdata/**", "**/.venv/**", "**/fixtures/**"]

[[bundle]]
name        = "acme-backend"
description = "Backend conventions for the acme platform"
skills      = ["python-architect", "sql-architect", "hexagonal-architecture"]
```

This is not a replacement for `SKILL.md`; it only constrains discovery. Without it, a git source
walks `roots` from `SKILL.md` in the source spec, and a monorepo will discover skills inside
`node_modules`, vendored tool directories and test fixtures.

#### 4.5 Index v2

`IndexEntry` gains `source`, `private`, `homepage`, `withdrawn`; `IndexVersion` gains a required
`sha256`. `deprecated` keeps its current meaning (a redirect target exists). `withdrawn` means the
skill is no longer published and must not be installed fresh; existing installs remain resolvable
only while their pinned tarball exists.

#### 4.6 Naming

- Directory name = skill name (already enforced at publish: `cmd/generate-registry/main.go:111`).
- Names are unique within a source. A duplicate is a hard error, not a first-match win.
- `rsk.mod` and `rsk.lock` use qualified `source:name` once more than one source is configured.
- Retired names keep working through `registry/aliases.toml`; `withdrawn` is not a retirement
  (`internal/skill/names.go:20-34`).

### 5. Nested layouts

Already supported for discovery: `skill.Walk` recurses and stops at the first directory containing
`SKILL.md` (`internal/skill/registry.go:53`), and the publisher does the same
(`cmd/generate-registry/main.go:88-127`). A nested layout works today; what is missing is
determinism.

Required guardrails:

1. **Unique names per source.** `skill.FindByName` returns the first lexical match
   (`internal/skill/util.go:17`); two `api/` directories silently shadow each other. Make it an
   error listing both paths.
2. **Roots and excludes.** Honour `rsk-skills.toml` / the source spec instead of walking the whole
   checkout.
3. **Validation parity.** Publishing validates names and retirements (`ValidateName`,
   `CheckRetirements`); filesystem sources do not. Run the same checks in `Walk` so a repo that
   works locally cannot fail to publish.
4. **Replace path-based `personal`.** `IsPersonalPath` scans for a `personal` segment
   (`internal/skill/util.go:9`); under arbitrary nesting it means nothing. Use the source's
   `private` flag, plus an explicit `personal = true` per skill in `rsk-skills.toml` for the
   published-but-author-specific case.
5. **Cache the walk.** `fsSource.find` re-walks the root on every lookup
   (`internal/source/fs.go`). Cache per process, keyed by root mtime.
6. **Symlinks are not followed.** `filepath.WalkDir` does not descend into symlinked directories, so
   a skill symlinked into a source is invisible. Document it; do not follow links (escape risk).

### 6. Resolution, update and cache semantics

| Operation | Behaviour |
|---|---|
| `install name` | chain resolves; ambiguity is an error; preview shows source id next to each skill |
| `install source:name` | exact source; no fallback |
| `install --from ./path` | one-off path source, not persisted |
| `update` | per transport: registry compares version, git fetches the ref and compares commit to the lock, path reports drift only |
| `status` | reads the lock; verifies the cached path still matches the recorded commit/digest; no network unless `--remote` |
| `catalog` | lists public sources by default, adds private sources with `--include private` |
| `mcp` | same policy; `install_skills` refuses private sources unless the project opted in |

Cache layout: `<cache>/<source-id>/<name>/<version|commit>/`. Separate roots per source prevent a
private extraction from being mistaken for a public one. All extraction keeps the existing zip-slip
and entry-type guards (`internal/source/registry.go:256-324`).

Index fetching gets ETag/`If-None-Match` and a TTL inside `source.Registry`, lifting the idea that
already exists only in MCP (`internal/mcp/cache.go:61`) so the CLI stops refetching on every command.

### 7. Publishing, decoupled

Extract the generator into `internal/registrybuild` and expose it as `rsk registry build`:

```
rsk registry build \
  --skills-dir . \
  --output-dir dist \
  --archive-url-template 'https://skills.acme.dev/{name}/{version}.tar.gz' \
  --existing-index dist/index.json
```

Flags already exist for `--skills-dir`, `--output-dir`, `--existing-index` and `--github-repo`
(`cmd/generate-registry/main.go:230-250`); the change is (a) an archive URL *template* instead of a
GitHub-only construction at `:285`, (b) SHA-256 per version, and (c) living behind `rsk` so another
org does not fork the workflow.

### 8. Migrating existing skills from public to private

1. Choose the set, then split by privacy class: anything unpublishable goes to
   `acme/skills-private`, everything else stays public.
2. Mirror the current layout (`skills/<domain>/<name>/SKILL.md`) so no skill content changes.
3. Add the `[[source]]` block, install from it, confirm `rsk status` reports the new source id.
4. In the public repo: remove the directories and mark the index entries `deprecated: true`
   (`IndexEntry.Deprecated` already exists and is skipped by `All` and `FromIndex`). Do **not** add
   them to `registry/aliases.toml`; retirements map old to new names and reject reuse.
5. Keep the published release tarballs. Deleting them breaks every pinned lock.
6. Accept that existing users keep a cached copy silently: `update` skips names missing from the
   index (`cmd/rsk/update.go:474`). Ship a migration note (uninstall, add the private source,
   reinstall) rather than relying on `update`.
7. If any moved skill contains internal names or credentials, scrub them: they are already in public
   git history, the published index and public tarballs.

### 9. Phasing

| Phase | Deliverable | Tests | Risk |
|---|---|---|---|
| 0 | `Spec`, `Chain`, factory; qualified names; `mod`/`lock` v2 | resolution-order, ambiguity, v1→v2 migration | schema migration touches every command |
| 1 | git transport (fetch, checkout, commit drift); path transport; `--from` | bare-repo fixture, private repo over SSH, offline cache hit | cache layout change |
| 2 | URL/tarball transport; SHA-256 in index and lock; auth transport + redaction | local auth server, digest mismatch aborts, no secret in logs | secret leakage in errors |
| 3 | `rsk registry build`; `rsk-skills.toml`; docs; migration guide | golden index, nested-layout fixture, duplicate-name failure | publishing workflow rewrite |

Each phase leaves `rsk` usable; the chain in phase 0 with only the existing two sources is a no-op
behaviourally.

## Risks

- **Schema migration** — `rsk.mod`/`rsk.lock` v2 must read v1 forever or migration stops adoption.
  Mitigation: read v1, write v2, add a `rsk migrate` dry-run.
- **Executable content** — a remote `SKILL.md` is instruction content executed by the agent. Pair
  remote sources with the per-skill permissions work (`plans/2026-11-17-per-skill-permissions.md`)
  and never auto-install from a newly added source without confirmation.
- **Secret leakage** — tokens in URLs, error strings or logs. Mitigation: header-only auth, a
  redacting transport, a test that greps captured logs for the token.
- **Name collisions** — bare names across sources. Mitigation: ambiguity is an error; qualified
  references in the lock; `install --from` to disambiguate.
- **Monorepo cost** — walking a large checkout on every lookup. Mitigation: `roots`/`excludes`,
  per-process cache, `rsk-skills.toml`.
- **False privacy** — `personal: true` in a public index stays a filter. Mitigation: rename the
  concept, document that only `private` sources and authenticated transports enforce anything.
- **Windows symlinks** — unchanged; installing from a remote source does not fix it.

## Open questions

1. **`personal` naming.** Keep `personal` for published-but-author-specific and add `private`, or
   rename `personal` to something explicit and break the flag? (Proposal: keep, document, add
   `private`; `--include personal` and `--include private` stay distinct.)
2. **Is a private registry required, or is git enough?** (Proposal: git first; build the authenticated
   registry only when a team needs discovery without clone access.)
3. **Repo granularity.** One private repo, or one per team? (Proposal: one per privacy class; split
   further only if teams need different credentials.)
4. **Signature verification.** Add minisign/sigstore over tarballs now, or rely on TLS + SHA-256 +
   git commit identity? (Proposal: defer; the digest is the integrity record, signing is a separate
   ADR.)
5. **Subtree privacy.** Is sparse checkout plus a server-side filter worth supporting, or is "one
   repo per privacy class" the documented answer? (Proposal: the latter.)
6. **`--include private` versus per-source policy.** Should privacy be a flag, a source field, or
   both? (Proposal: source field controls listing; the flag is required to install.)

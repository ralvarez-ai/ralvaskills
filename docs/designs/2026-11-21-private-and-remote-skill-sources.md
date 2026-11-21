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
private repositories, arbitrary git refs, direct paths and URLs, with reproducible locks, globally
unique skill references (`@owner/repo/name`) and no secrets stored by `rsk`. Written to be evaluated
outside this repository: it states the problem, the constraints the current code imposes, the
options, one proposed design, the migration already executed for the author's personal skills, and
the open questions.

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
- Resolution order and source construction are spread over many sites, not five:
  `install.FindSkill` (`internal/install/resolve.go:33-56`), `LocalSource` (`:26`),
  `newLocalSource` (`cmd/rsk/resolve.go:19`), `loadSkillsBySource` (`cmd/rsk/catalog.go:167`), and
  direct `source.NewOfficial` / `NewRegistry` / `NewLocal` construction in `cmd/rsk/install.go:102,349`,
  `new.go:298`, `update_lock.go:47`, `sync.go:229,234`, `update.go:313`, `catalog.go:182`,
  `internal/update/check.go:132`, `internal/mcp/install_tool.go:100` and
  `internal/mcp/server.go:53,56`. Three of them (`update.go`, `update/check.go`, `mcp/server.go`) need the
  concrete `*Registry` for `Index`, not a `Resolver`, so the factory must hand that back too.
- The registry client has no authentication: `NewRegistry` builds a bare `http.Client`
  (`internal/source/registry.go:73`). There is no header, token or credential concept anywhere.
- The lock records name, version and source kind only (`internal/manifest/lock.go:22-34`), keyed by
  bare name (`UpsertLockEntry`, `:78`). A remote skill pinned to a git ref is therefore not
  reproducible, and two sources with the same skill name collide silently.
- `personal/` was an install filter, not a privacy boundary: those skills lived in the public repo,
  appeared in the public `registry/index.json` with `"personal": true`, and shipped as public
  release tarballs. `--include personal` (`cmd/rsk/install.go`) and the MCP refusal
  (`internal/mcp/install_tool.go:145`) only discouraged use. The seven skills have since moved to a
  private repository (see Migration); the path-based concept is still in the code.

Consequences:

1. A project cannot use the public registry **and** a private skill repo at the same time.
2. A skill that exists only in a git repository cannot be installed at all.
3. Registries cannot be private by design (see Decisions): private skills belong in git, so this is
   a boundary to state, not a gap to fill.
4. Nothing records where a remote skill came from, so `rsk update` and `rsk status` cannot verify it.
5. Publishing is welded to this repository: `cmd/generate-registry` hardcodes GitHub release asset
   URLs (`cmd/generate-registry/main.go:328`) and this repo's `skills/` layout.
6. Two publishers cannot ship a skill with the same name and be told apart.

## Goals and non-goals

Goals:

- **G1** A project or machine can declare N sources of mixed kinds: hosted registry, private
  registry, git repository, local path, official cache.
- **G2** Remote install by git ref, by tarball URL, and by path, without publishing anything.
- **G3** Private sources are git repositories or local paths, authenticated by git itself (SSH key,
  credential helper) or by the filesystem. `rsk` stores no secret and has no credential concept at
  all. Registries are public and unauthenticated. Sources are managed from the CLI, globally by
  default.
- **G4** Every install from a git or registry source is reproducible: the lock records the
  transport identity (git commit or tarball SHA-256) and the source identity. Path sources are the
  stated exception.
- **G5** Repo-wide access is sufficient to onboard a private skills repository; no per-skill ACL.
- **G6** Arbitrary nesting inside a source works deterministically: same name means same skill,
  duplicates are an error.
- **G7** Anyone can produce a compatible public index from their own repo without forking this
  repo's workflow.
- **G8** A skill has a globally unique reference, `@owner/repo/name`, so same-named skills from
  different publishers can be told apart.

Non-goals (deliberately deferred):

- Per-skill entitlements, metering or selling skills.
- OCI artifacts as a transport.
- Encryption of skills at rest (plaintext ends up in the agent's context anyway).
- Signature verification of skill content (deferred, see Decisions).
- Windows symlink strategy (unchanged from today).
- Installing two same-named skills into one scope (see 5.6 and 2.3).

## Constraints imposed by the current code

| # | Constraint | Where | Effect on this design |
|---|---|---|---|
| 1 | One local **XOR** one registry | `internal/install/resolve.go:26-31`, `internal/config/config.go:50` | Source chain is mandatory; this is the primary change |
| 2 | Closed `Source` enum (40 references) | `internal/skill/source.go:8-28`, `Label()` `:37-48` | Replace kind-based enum with source identities; keep wire compatibility |
| 3 | Order and construction hardcoded at about a dozen sites | listed in Context, plus `internal/mcp` | One `Chain` and one factory, consumed everywhere; a lint or test forbids `source.New*` outside the factory |
| 4 | Lock has no provenance | `internal/manifest/lock.go:22-43` | Lock schema bump to v2 |
| 5 | Lock keyed by bare name | `internal/manifest/lock.go:78,91` (`UpsertLockEntry`, `RemoveLockEntry`) | Qualified keys and a lock migration |
| 6 | No HTTP auth | `internal/source/registry.go:73` | Stays that way: registries are public by design, so no auth layer is added |
| 7 | No integrity metadata | `internal/source`, `cmd/generate-registry` | SHA-256 in index and lock |
| 8 | Publishing tied to this repo | `cmd/generate-registry/main.go:328` (URL), `:32` (`releaseURLPath`) | Extract a buildable library + CLI |
| 9 | `personal/` is path-based and public | `internal/skill/util.go:9`, `internal/skill/registry.go:102` | Split `personal` (published) from `private` (never published) |
| 10 | Filesystem lookup re-walks the root | `internal/source/fs.go` | Cache per process; matters for monorepos |
| 11 | `rsk.mod` is `map[string]string` | `internal/manifest/mod.go:19-26` | Schema bump to v2 for qualified names + `[[source]]` |
| 12 | Installed skills are symlinks named by bare name | `internal/update/check.go:177` | Two same-named skills cannot share one skills directory |
| 13 | Installed version is parsed from the cache path `<registryCache>/<name>/<version>/` | `internal/update/check.go:182` | Do not change the default registry layout; new sources get a new root |
| 14 | `schema.Check` accepts only `0` or exactly the current version | `internal/schema/schema.go` | Must become a supported range; older `rsk` hard-fails on a v2 file |
| 15 | `internal/git` is `Pull`, `Clone`, `Behind` only; stderr goes straight to the terminal | `internal/git/git.go` | No fetch-by-SHA, export, hardening or redaction exists; extend it, do not assume it |

## Options

| Option | What changes | Pros | Cons | Effort |
|---|---|---|---|---|
| **A. Private registry** | Auth on `Registry`, second URL | One new transport; discovery and versions stay first-class; fits multi-team | Requires a bucket, an auth front, a publisher, token distribution; still no git-ref installs; puts secrets inside `rsk` | M |
| **B. Source chain, git-first** | Source spec + factory + chain; git transport; lock v2 | Fixes XOR blocker; private access is git's problem; no new infra; highest reuse | Git-only versioning (commit ids, frontmatter version is display); needs drift handling | M |
| **C. Vendor / submodules only** | None in `rsk` | Zero code | No resolution, no updates, no separation of public/private | S |
| **D. OCI artifacts** | New transport | Nice provenance, registries already exist | Tooling and auth cost, no skill-specific benefit yet | L |

**Decision:** B. A is rejected for private skills: a private source is a git repository or a path, and
registries stay public. That keeps secrets out of `rsk` entirely. C remains possible as a degenerate
path source.

Accepting this design should produce an ADR (next free number: 0003) covering the source-chain
decision, the reference grammar and the lock v2 bump. The decisions it rests on were settled on
2026-11-21 (see Decisions): private skills in git only with no private registry, privacy as a source field
only, schema v2 written only when needed, version selection through tags, and any number of private
sources. The grammar in 2.1 and the trust rules in section 3 are hard requirements of phase 1, not
follow-ups. Subtree privacy (question 5) is the one question still open.

## Design

### 1. Source model

A source is a named, typed descriptor with a canonical identity. It is declared in the machine
config and may be added per project, with the restrictions in section 3.

```go
// internal/source/spec.go (new)
type Kind string // "registry" | "git" | "path" | "official"

type Spec struct {
	ID       string   // canonical identity "@owner/repo"; derived for git, declared otherwise
	Type     Kind
	URL      string   // registry base URL, git URL, or filesystem path
	Rev      string   // git ref; empty means the source's default branch
	Roots    []string // directories to walk inside the source; empty means "."
	Excludes []string // glob patterns pruned before walking
	Private  bool     // git and path sources only; kept out of surfaces that can leave the machine unprompted
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

func New(spec Spec, opts Options) (Resolver, error)   // the only factory
func Build(specs []Spec, opts Options) (Chain, error) // declaration order
func (c Chain) Find(ctx context.Context, ref Ref) (skill.Skill, error)
func (c Chain) All(ctx context.Context) ([]skill.Skill, error)
```

There is no `priority` field. Duplicates across sources are an error, so an ordering could never
decide a resolution; the chain keeps declaration order only for stable listings.

Resolution rules:

1. A qualified reference (`@owner/repo/name`) resolves in exactly that source. An unknown source is
   an error, not a fallback (see 2.2 for implicit sources).
2. A bare `name` queries **every** source. It resolves only if exactly one source has it.
3. Two sources providing the same bare name is an **error** naming both. The user disambiguates with
   the qualified form.
4. **Fail closed.** If any source cannot be reached while resolving a bare name, the command fails.
   Otherwise an expired credential on the private source would silently resolve the name from the
   public one (dependency confusion). "Reachable" means a name index for that source exists and is
   younger than a TTL (default 24 h), or was refreshed now; a stale-but-present index is accepted
   inside the TTL, so the rule does not force a network round trip per lookup. A source that has
   never been fetched, or whose refresh fails past the TTL, fails the command. `--offline` accepts
   any existing index and says so.
5. **Lock first.** Installing from `rsk.lock` never consults the chain: the lock names the source,
   the URL and the commit or digest.
6. **Adding a source never changes how an existing bare name resolves.** `rsk source add` runs the
   scan first and warns about every name that now exists in two sources; resolution of a bare name
   already in `rsk.mod` stays pinned to its recorded source until the user re-resolves it.
7. `private = true` sources are excluded from `catalog` listings, MCP resources, `rsk://catalog` and
   `search_skills` unless the command is run with `--include private` (mirroring today's
   `--include personal`). Registering the source is the opt-in to install from it: no flag is
   needed per install.

Existing commands keep their flags; `newLocalSource`/`LocalSource` become thin wrappers over
`Build(cfg.Sources, ...)`, and the direct `source.NewOfficial`/`NewRegistry` calls move behind it, so the
dozen sites collapse to one.

Threat model for `private`: it keeps names, descriptions and paths of private skills out of output
that is shared or captured without a deliberate act (MCP responses, `--output json` pasted into a
public issue, screen shares). It does not protect content from a user who already holds the
credential.

### 2. References and implicit sources

#### 2.1 Grammar

```
@<owner>/<repo>/<name>[@<version>]       # qualified: owner/repo identify the source
@<owner>/<repo>//<path/to/name>          # qualified with a repo-relative path (disambiguation)
<name>[@<version>]                       # bare: resolved by the chain (rules 2-4)
```

- Skill names are kebab-case (`ValidateName`), so `/` and `@` never occur inside a name.
- A leading `@` means qualified. The version is whatever follows the last `@`, which keeps
  `ParseNameVersion` and the `<skill>@v<version>` release tags working.
- The folder a skill sits in is not part of its identity: the name is the directory basename
  (`internal/skill/registry.go:97`), found under the source's `roots`. Use the `//` form only when
  two directories in one source share a name; it selects one copy and the installed directory is
  still the basename. The lock records the repo-relative `path` so a later reorganisation of the
  repository cannot change what a locked commit resolves to.
- For a git source, `@owner/repo` is the GitHub owner and repository. Other hosts use an explicit
  `rsk source add` and a declared identity (default `@local/<dirname>` for a path).
- One source can have several identities: a **canonical** one and **aliases**. The default registry
  is canonically `@ralvarez-ai/ralvaskills`, the organisation that owns the repository, with
  `@ralvarezdev/ralvaskills` as an alias (the owner the generator's release URLs still use,
  `cmd/generate-registry/main.go:35`; GitHub redirects it today, verified). Bare names keep resolving there, so
  nothing existing breaks. A third-party public registry declares its identity in its index
  (`source`, see 5.5); the identity is a label, and the lock pins the URL.
- Alias rules: the lock and `rsk.mod` always store the **canonical** form, so two spellings never
  produce two entries; an alias is accepted on input and normalised; an alias must map to the same
  source and transport as its canonical identity; the built-in and declared identities are checked
  **before** the implicit GitHub derivation in 2.2, otherwise `@ralvarezdev/ralvaskills` would be
  cloned as a git repository and bypass the registry's SHA-256 record.
- The official cache is an ordinary git source, `@anthropics/skills`, not a special kind.
- Lock and `rsk.mod` keys use the qualified form **iff the skill comes from a source other than the
  default registry**. That is a property of the skill, not of the machine, so the same committed file
  reads the same everywhere (a rule based on "how many sources are configured" would churn the file
  between machines).
- Identities are case-folded (GitHub owners and repositories are case-insensitive). The identity of
  a GitHub-hosted git source must equal the one derived from its URL, so `@anthropics/skills`
  cannot be registered against another URL. Declared identities are only for non-GitHub hosts and
  paths; the scope `@local` is reserved for paths.
- **Alias claims are never taken from fetched data.** Aliases come only from the built-in table
  shipped with `rsk` or from the user's own machine config; a fetched index or `rsk-skills.toml`
  that claims `@ralvarez-ai/ralvaskills` is ignored and reported.
- **Versions in a git source.** The skill's version is the frontmatter at the resolved commit.
  `name@1.2.0` resolves through a tag `<name>@v1.2.0` in that repository (the convention this repo
  already uses for releases) and is an error when the tag is missing; version ranges in `rsk.mod`
  (`^1.2.0`) are not supported for git sources in the first cut, only `*` or an exact tag. Because
  git has no version index, `update` compares content, not versions (see 5.3).

How the kinds differ:

| Kind | Reference | Lock `transport` | Integrity record |
|---|---|---|---|
| Default registry | `go-grpc` or `@ralvarez-ai/ralvaskills/go-grpc` | `registry` | `sha256` |
| Remote git | `@owner/repo/name` | `git` | `commit` and `path` |
| Official | `@anthropics/skills/name` | `git` | `commit` and `path` |
| Local path | `@local/<dir>/name` | `path` | none; `status` reports drift only |

Local clone mode (`repo_path`) behaves like the path row. The installed symlink is always the bare
name, so the source shows only in the lock, in `status` and in the symlink target; `catalog` and
`list` replace today's `ralva/anthr/reg` labels (`Source.Label()`) with the identity.

Examples: `@ralvarezdev/skills/go-architect`, `@acme/platform-skills/python-architect@1.2.0`,
`@acme/monorepo//packages/api/skills/api`.

#### 2.2 Implicit sources

`rsk install @owner/repo/name` works without registering the source: `rsk` derives
`git@github.com:owner/repo.git` (SSH, so a private repo works with the user's key; HTTPS as a
fallback), clones, and records the transport identity in the lock. Because the reference clones
code and instructions an agent will follow:

- First use of an unregistered `@owner/repo` shows the URL and asks for confirmation. `--yes` does
  not cover it. CI passes `--trust @owner/repo` explicitly; the MCP server never can, and
  `install_skills` refuses unknown sources.
- Nothing is written to machine config unless `--save` is given.
- SSH is tried first and HTTPS second. The lock stores the identity and commit; on restore the
  machine's own source configuration decides the transport URL, so one machine's SSH choice does not
  break another's.

#### 2.3 What qualification does not solve

A qualified reference disambiguates **resolution**, not **installation**. Installed skills are
symlinks named by the bare name (`internal/update/check.go:177`) and the agent keys on the
frontmatter `name`; `@a/b/x` cannot be a directory name. Two same-named skills in one scope is
therefore a hard error. They can coexist across scopes (one project, one global), where the agent's
own precedence applies, or via an explicit `--as <alias>` rename at install time if the skill
tolerates it.

### 3. Authentication, authorization and trust

Rule: **`rsk` holds no credential at all.** Private skills live in git repositories or on the
filesystem, so authentication is git's and the operating system's: `rsk` never reads, stores,
forwards or logs a secret, and no `rsk` file, flag or MCP parameter carries one. Registries are
public and unauthenticated; a registry source marked `private` is rejected.

| Transport | Mechanism | Granularity | Secret lives in | Revocation | Fit |
|---|---|---|---|---|---|
| git | SSH key / deploy key (read-only) | per repo | `~/.ssh`, agent | remove key | default for a developer |
| git | SSH cert from an org CA | repos the CA allows | agent | CA policy | orgs with many repos |
| git | HTTPS fine-grained PAT | selected repos | git credential helper | expire or rotate | headless machines |
| git | GitHub App installation token | installation repos | minted, 1 h TTL | uninstall app | fleets; best CI story |
| git | OIDC to cloud STS | one workflow run | none | n/a | publish/build only |
| path | filesystem permissions | one directory | OS | chmod | local checkouts |

Authorization is source-wide: one repository = one credential = one policy, and there is no limit
on how many such sources a user registers. A client that can clone a repository can read every
directory in it, so finer granularity means splitting repositories (decision 5). Read and publish
credentials are separate: developers hold read only; publishing runs in CI with an App token or
OIDC and never on a laptop.

#### Trust rules (a source spec is untrusted input)

A project's committed `rsk.mod` can name any source, so cloning a repository must not be enough to
make `rsk` send a secret or run code.

1. **No credential field exists.** A `[[source]]` has no `auth`, `token` or `env:` reference, so a
   committed `rsk.mod` cannot make `rsk` send a secret anywhere. What git sends to a host (SSH key
   offer, credential-helper entry for that host) is the user's own git configuration, which is why
   rules 2 and 3 gate which hosts are contacted at all.
2. **First use of a project-declared or implicit source needs confirmation** showing the URL
   and transport. The grant is stored per machine (`~/.config/rsk/trust.json`) keyed by identity
   **and** URL, and is made only by `rsk source trust <id>`, which requires a TTY and fails without one. When an
   agent runs `rsk` through a shell tool there is no TTY, so the command refuses and names the
   exact `rsk source trust` line for the user to run themselves, the same hand-back the rsk
   workflows already use for `sudo`.
3. **Git URL allowlist.** Accept `ssh://`, `https://` and scp-like `git@host:path` only. Reject
   `ext::`, `file://` and any value starting with `-`; pass `--` before the URL; set
   `protocol.ext.allow=never`, `protocol.file.allow=never`; never recurse submodules.
4. **Non-interactive git.** `GIT_TERMINAL_PROMPT=0`, ssh `-o BatchMode=yes`, stdin closed. The MCP
   server is a stdio process and must never block on a prompt or write to stdout from git.
5. **Redaction.** Git stderr is captured and scrubbed, not forwarded (`internal/git` forwards it
   to `os.Stderr` today).
6. **Roots are untrusted.** `roots` and `excludes` from a remote `rsk-skills.toml` must be relative,
   contain no `..`, and resolve inside the checkout; symlinks are not followed.
7. **Remote text enters always-on context.** A pinned skill's `DIGEST.md` is written into
   `.rsk/SKILLS.md`, which the agent loads every session. Pinning a skill from a non-default source
   shows the digest text and asks for confirmation.
8. **`rsk.lock` is untrusted too.** Restoring from a lock names a URL; an unknown or untrusted
   source in a lock goes through rule 2 like any other, so the lock-first rule (resolution rule 5)
   skips the chain but never the trust gate.
9. **Identity binds to URL.** A project `[[source]]` whose id matches a machine source but whose URL
   differs is rejected, and a trust grant is keyed by identity and URL, so changing the URL
   revokes it. The same holds for `RSK_REGISTRY_URL` (`internal/config/config.go`): overriding
   the default registry URL keeps the identity, but the lock pins the URL, and a mismatch is
   reported instead of silently followed.
10. **No credential storage of any kind**, including keychains (a libsecret/dbus dependency on
    Linux that nothing here needs).

### 4. Repo-wide setup and the CLI

Repo-wide is the supported and recommended boundary: one skills repository, one credential, every
skill in it shares the policy. Privacy class, not skill domain, decides the repo split. There is no limit on how
many private sources a user registers: one per team, per client or per project is fine, each with
its own credential.

Sources are managed from the CLI. Machine scope is the default; project scope writes the committed
`rsk.mod` and is subject to the trust gate (rule 2). `--private` applies to git and path sources only:

```
rsk source add @acme/skills-private --roots skills --private        # machine config
rsk source add @acme/skills-private --project                       # [[source]] in .rsk/rsk.mod
rsk source add ./my-skills --id @me/scratch                         # path source
rsk source list | remove <id> | test <id>                           # test = reachability + transport
rsk install @acme/skills-private/go-architect [--as alias] [--trust @acme/skills-private]
rsk install --from ./path/to/skill                                  # one-off, not persisted
```

`source add` derives the SSH URL from `@owner/repo`; pass a full URL for other hosts. Each
state-changing command states the action and waits for confirmation, as `rsk-guide` requires for
`install`/`pin`. Equivalent TOML in `~/.config/rsk/config.json` (`sources`) or `.rsk/rsk.mod`:

```toml
[[source]]
id      = "@acme/skills-private"
type    = "git"
url     = "git@github.com:acme/skills-private.git"
rev     = "main"
roots   = ["skills"]
private = true
```

Onboarding a second machine is: clone access plus `rsk source add`. No secret is copied, because
git supplies it.

#### Command and MCP surface

Every operation is exposed on both surfaces so an agent can do the tedious setup and a human can do
the same from a terminal. The existing server already follows this shape: `rsk install` and
`install_skills` both call `internal/install`, and `install_skills` carries
`_meta["anthropic/requiresUserInteraction"]` so the client forces approval.

| CLI command | MCP tool | Kind |
|---|---|---|
| `rsk source list` / `show <id>` | `list_sources` | read-only |
| `rsk source test <id>` (reachability, transport, skill count, duplicate names) | `test_source` | read-only |
| `rsk source scan <path-or-@owner/repo>` (dry-run discovery: skills found, name clashes, naming and frontmatter violations) | `scan_source` | read-only |
| `rsk resolve <ref>` (source, path, version and commit a reference gives; why a bare name is ambiguous) | `resolve_reference` | read-only |
| `rsk source add <@owner/repo\|url\|path> [--id --rev --roots --private --project --save]` | `add_source` | mutating, forced approval |
| `rsk source remove <id>` | `remove_source` | mutating, forced approval |
| `rsk install <ref>... [--as --from]` (exists; gains qualified references) | `install_skills` (exists) | mutating, forced approval |
| `rsk repo init <dir>` (scaffold a skills repo: layout, `rsk-skills.toml`, `AGENTS.md`, `.gitignore`, `.editorconfig`, then `rsk new`) | `init_skills_repo` | mutating, forced approval |
| `rsk repo import --from <src> <names>... --map name=languages/go` (copy skills, keep layout and subcategories) | `import_skills` | mutating, forced approval |
| `rsk source trust <id>` / `untrust <id>` (grant or revoke the per-machine trust record; requires a TTY) | none | CLI only, user-only |
| `rsk registry deprecate <names>...` (the index edit done by hand in the migration) | none | CLI only, publishing |
| `rsk migrate --dry-run` (`rsk.mod` / `rsk.lock` v1 to v2) | none | CLI only |

`project_profile` and `search_skills` stay; `project_profile` gains a `sources` section
(configured, proposed, untrusted). Creating the GitHub repository stays out of `rsk`: the agent
runs `gh` itself under the user's normal permission prompts.

Parity rules:

1. **Core first.** Each operation is a function in `internal/` taking and returning plain structs;
   the CLI and MCP wrappers only parse input and render output.
2. **Preview parity.** Every mutating command has `--dry-run --output json`; every mutating MCP tool
   takes `dry_run: true` and returns the same JSON. The agent shows the preview, then applies.
3. **A parity test** enumerates the registered MCP tools and fails when one has no CLI counterpart.

What an agent must never do (it can be steered by a remote skill's text):

- **No trust parameters over MCP.** `--trust` and `--yes` are CLI-only. No tool schema or `rsk`
  file accepts a secret, so none can land in a transcript.
- **Forced approval on every mutating tool.** The client renders the tool's input, so inputs stay
  explicit (full URL, transport, `private`), never an indirection the user cannot read.
- **`add_source` cannot grant trust.** It registers the source as untrusted. The first clone needs
  `rsk source trust <id>` in the user's own terminal (trust rule 2); until then `install_skills` and
  `test_source` return an error that names that command.
- **File-writing tools are confined.** `init_skills_repo` and `import_skills` write only inside the
  directory they are given, refuse a non-empty target and any `..` in `--map` paths.
- **Tool output from a source is data**, never instructions: descriptions and digests from a remote
  repository are not acted on.

Typical agent flow for "set up my private skills repo": `project_profile`, `scan_source`, then
`init_skills_repo` and `import_skills` (each previewed with `dry_run`), `add_source`, `test_source`,
`install_skills`, with the user approving each mutating step. The `rsk-guide` skill documents the
order; updating it is a separate change because it publishes on push.

### 5. Schemas

#### 5.1 Source descriptor

As above, in `config.json` (`sources`) and `rsk.mod` (`[[source]]`). A project source with the same
`id` as a machine source must have the same URL (trust rule 9); otherwise it is rejected.

#### 5.2 `rsk.mod` v2

```toml
version = 2
tools   = ["claude-code", "opencode"]
pinned  = ["@acme/skills-private/go-architect"]

[skills]
"@acme/skills-private/go-architect" = "*"
"go-grpc"                           = "^1.2.0"   # bare = unambiguous across sources
```

`schema.Mod` bumps to 2. Version 1 files are read as before (all names bare, implicit sources from
`repo_path`/`registry_url`). Compatibility:

- `schema.Check` must accept a range, not `stored == current` (`internal/schema/schema.go`).
- Write v2 **only when a non-default source or a qualified key exists**; otherwise keep writing v1,
  so teammates on an older `rsk` are unaffected until they adopt a new source.
- An older `rsk` hard-fails on a v2 file with the existing "upgrade rsk" message. That is accepted.

#### 5.3 `rsk.lock` v2

```toml
version = 2

[[skills]]
name      = "@acme/skills-private/go-architect"
source    = "@acme/skills-private"
transport = "git"
url       = "git@github.com:acme/skills-private.git"
rev       = "main"
commit    = "9f3c1ab7e2d4c1e0f8b6a5d4c3b2a1908f7e6d5c"
path      = "skills/backend/go-architect"
version   = "1.2.0"            # frontmatter, display only
permissions = ["Bash(make lint)"]
granted     = []

[[skills]]
name      = "go-grpc"
source    = "@ralvarez-ai/ralvaskills"
transport = "registry"
url       = "https://skills.ralvarez.dev"
version   = "1.4.0"
sha256    = "3b1f..."
```

`LockEntry` gains `Source`, `Transport`, `URL`, `Rev`, `Commit`, `Tree`, `Path`, `SHA256`; `Name` becomes qualified.
For git, the commit id is the integrity record and `tree` (`git rev-parse <commit>:<path>`, the object
id of the skill's own directory) is the content identity. `update` compares `tree`, not `commit`:
in a monorepo every unrelated commit moves `commit`, and comparing it would report every skill as
outdated. For registry and URL transports, `sha256` is the
verification record and is required from index v2 onward. A **path** source records neither, and **no absolute path**: the lock keeps the identity only (the
machine's source config supplies the path), so a committed lock never carries a machine-specific
location. It is not reproducible and `status` reports drift only.

v1 `source` values map onto identities when read: `registry` to the default registry's canonical
identity, `official` to `@anthropics/skills`, `local` to the local clone's `@local/<dir>`.
`ParseSource` (`internal/skill/source.go`) rejects unknown values today, so it must accept the
identity form before any v2 file is read.

#### 5.4 Hosting manifest: `rsk-skills.toml`

Every repository that carries skills can declare what `rsk` should see. Optional for a path source,
recommended for git sources, required for monorepos.

```toml
# rsk-skills.toml at the repository root
roots    = ["skills", "packages/*/skills"]
excludes = ["**/node_modules/**", "**/testdata/**", "**/.venv/**", "**/fixtures/**"]

[[bundle]]
name        = "acme-backend"
description = "Backend conventions for the acme platform"
skills      = ["python-architect", "sql-architect", "hexagonal-architecture"]
```

This is not a replacement for `SKILL.md`; it only constrains discovery. Without it, a git source
walks `roots` from the source spec, and a monorepo will discover skills inside `node_modules`,
vendored tool directories and test fixtures. Its content is untrusted (trust rule 6). **Privacy is
a property of the source entry the user registered, never of the file**: a `private` key in
`rsk-skills.toml` is ignored, so a remote repository cannot declare itself public or private.

#### 5.5 Index v2

`IndexEntry` gains `source`, `homepage`, `withdrawn`; `IndexVersion` gains a required
`sha256`.

- `deprecated` keeps its **current** meaning: the entry is skipped by `All` and `FromIndex`
  (`internal/source/registry.go:39,125`). There is no redirect; renames redirect through the
  separate `aliases` map.
- `withdrawn` is new: no longer published, must not be installed fresh, existing locks still
  resolve while the pinned tarball exists.
- **Backfill.** Existing versions have no digest. Compute SHA-256 from the published release assets
  once, and make tarball construction deterministic (sorted entries, zero mtime) before any rebuild,
  or a rebuild changes every digest.

#### 5.6 Naming

- Directory name = skill name (already enforced at publish: `cmd/generate-registry/main.go:111`).
- Names are unique within a source. A duplicate is a hard error, not a first-match win.
- Names are unique within one install scope (2.3).
- Retired names keep working through `registry/aliases.toml`; `withdrawn` is not a retirement
  (`internal/skill/names.go:20-34`).

#### 5.7 Local overlay for private sources

Pins from a `private` source stay on the machine. `rsk` writes them to gitignored siblings and
merges both layers on read:

| Committed (shared) | Local (gitignored) |
|---|---|
| `.rsk/rsk.mod` | `.rsk/rsk.local.mod` |
| `.rsk/rsk.lock` | `.rsk/rsk.local.lock` |
| `.rsk/SKILLS.md` (generated pins) | `.rsk/SKILLS.local.md` |

- **Why the digest too.** `.rsk/SKILLS.md` is imported into the agent's context and committed, so a
  pinned private skill's `DIGEST.md` would otherwise leak through it, not just through the lock.
- **Exclusion.** `rsk` appends the local files to `.git/info/exclude`, not the project `.gitignore`
  (repo-tooling rule: local-only exclusions stay out of the shared file).
- **Adapters.** Claude Code gets the local digest through `CLAUDE.local.md` (`@.rsk/SKILLS.local.md`);
  the OpenCode and Codex adapters need an equivalent that is not committed. Whether each client
  honours an uncommitted import, and what a teammate without the file sees, is **unverified** and a
  phase 1 test.
- **Trade-off.** Private pins are reproducible on the machine that made them, not across a team by
  git alone. A teammate with the same source registered re-resolves through `rsk install`. A project
  that is itself private can opt in to sharing with `--shared`, which writes the committed files.
- Default-registry and public-source pins keep going to the committed files.

### 6. Nested layouts

Already supported for discovery: `skill.Walk` recurses and stops at the first directory containing
`SKILL.md` (`internal/skill/registry.go:53`), and the publisher does the same
(`cmd/generate-registry/main.go:88-127`). The skill name is the directory's basename
(`internal/skill/registry.go:97`), so subcategories such as `languages/go/` or `academic/uru/` are
plain directories and work at any depth. What is missing is determinism.

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

### 7. Resolution, update and cache semantics

| Operation | Behaviour |
|---|---|
| `install name` | chain queries every source; ambiguity is an error; preview shows the qualified reference next to each skill |
| `install @owner/repo/name` | exact source; no fallback; unknown source follows 2.2 |
| `install --from ./path` | one-off path source, not persisted, not reproducible |
| `update` | per transport: registry compares version, git fetches the ref and compares the skill directory's tree id to the lock, path reports drift only. A qualified skill needs only its own source |
| `status` | reads the lock; verifies the cached path still matches the recorded commit/digest; no network unless `--remote` |
| `catalog` | lists public sources by default, adds private sources with `--include private`; an unreachable source warns and the rest are listed, with a partial marker |
| `mcp` | same policy; `install_skills` refuses private or unknown sources unless the project manifest already lists them |

Failure semantics split by command: **listing** (`catalog`, `search_skills`) degrades per source;
**resolving a bare name** fails closed (rule 4); **installing from the lock** needs only the sources
the lock names. An `--offline` mode uses the lock and the cache and never touches the network.

#### Git materialisation

One bare clone per source at `<cache>/git/<owner>/<repo>.git` (nested, because a flattened
`<owner>__<repo>` collides: `a__b/c` and `a/b__c` are both valid GitHub names), updated by
`git fetch origin <rev>` (or `fetch origin <sha>` for a locked commit). A skill is exported with
`git archive <commit> <path>` against the **local** bare clone (GitHub does not serve `git archive`
remotely) into `<cache>/sources/<owner>/<repo>/<name>/<commit>/`. The export applies the registry
extraction guards: only regular files and directories, every path confined to the skill directory,
a size cap, no symlinks or submodule entries. A repository's `.gitattributes` (`export-ignore`,
`export-subst`) is honoured by `git archive` and can alter content; the lock's `tree` is computed
before export. If a locked commit is no longer
reachable (force-push), the command fails naming the commit; the mitigation is to pin a tag or
protected branch. `internal/git` gains `Fetch`, `Archive` and the hardening in section 3.

#### Cache layout and existing installs

The default registry keeps `<registryCache>/<name>/<version>/`, because
`InstalledVersionFromTargets` parses that shape (`internal/update/check.go:182`) and every existing
symlink points into it. Only the new sources use `<cache>/sources/<owner>/<repo>/<name>/<commit|version>/`,
so no migration of installed links is needed. A private extraction can never be mistaken for a
public one because the roots differ. All extraction keeps the existing zip-slip and entry-type
guards (`internal/source/registry.go:256-324`).

Index fetching gets ETag/`If-None-Match` and a TTL inside `source.Registry`, lifting the idea that
already exists only in MCP (`internal/mcp/cache.go:61`) so the CLI stops refetching on every command.

### 8. Publishing, decoupled

Extract the generator into `internal/registrybuild` and expose it as `rsk registry build`:

```
rsk registry build \
  --skills-dir . \
  --output-dir dist \
  --archive-url-template 'https://skills.acme.dev/{name}/{version}.tar.gz' \
  --existing-index dist/index.json
```

Flags already exist for `--skills-dir`, `--output-dir`, `--existing-index` and `--github-repo`
(`cmd/generate-registry/main.go:230-260`); `gitHubRepoDefault` (`:35`) should become the canonical
`ralvarez-ai/ralvaskills`; the change is (a) an archive URL *template* instead of
a GitHub-only construction at `:328`, (b) SHA-256 per version, and (c) living behind `rsk` so
another org does not fork the workflow.

### 9. Migration: personal skills to a private repository

#### Plan

1. Choose the set, then split by privacy class: anything unpublishable goes to a private repo,
   everything else stays public.
2. Mirror the current layout (`skills/<domain>/<name>/SKILL.md`) so no skill content changes;
   subcategories are plain directories.
3. Register the source, install from it, confirm `rsk status` reports the new source identity.
4. In the public repo: remove the directories and hide the index entries.
5. Keep the published release tarballs. Deleting them breaks every pinned lock; `registry-remove`
   is therefore not run.
6. Existing users keep a cached copy silently: `update` skips names missing from the index
   (`cmd/rsk/update.go:474`). Ship a migration note (uninstall, add the private source, reinstall).
7. If a moved skill contains internal names, key paths or hosts, scrub them: they are already in
   public git history, the published index and public tarballs.

#### Executed (2026-11-21)

- **Private repo** `ralvarezdev/skills` created (private) with the seven skills, layout
  `skills/languages/go/` (`go-house-libs`, `go-library-builder`), `skills/workflows/`
  (`pi-iteration-workflow`), `skills/academic/uru/` (the four `uru-*` skills). `rsk` is initialised
  there for Claude Code and OpenCode with the MCP server and five pins; the pinned digests total
  about 2740 tokens against the 2000-token budget, an open decision.
- **ralvaskills**: `skills/personal/` removed, README section removed, the seven index entries set
  `deprecated: true` (the only hiding mechanism that exists today; `withdrawn` is not built), the
  site gate and `AGENTS.md` moved from 72 to 65 skills, token counts regenerated, changelog entry.
- **Aliases:** the three `uru-*` retirement aliases were deleted from `registry/aliases.toml` and
  the index by decision (2026-11-21). Those old names no longer redirect and may be reused.
- **Still open after the move**: the code still has the path-based `--include personal` concept, the
  skills' content remains in public history, and both repositories still need pushing by the owner.
- **Not usable yet through `rsk`**: the private repo is installable only by hand until the source
  chain exists. This design is what makes `rsk source add @ralvarezdev/skills` work.

### 10. Phasing

| Phase | Deliverable | Tests | Risk |
|---|---|---|---|
| 0a | `Spec`, `Chain`, factory collapsing the five call sites; no schema change | resolution order, ambiguity, fail-closed, behaviour identical to today | refactor touches 40 `Source` references |
| 0b | `@owner/repo/name` grammar; qualified keys; `mod`/`lock` v2 written only when needed; `schema.Check` range | v1 to v2 migration, old-binary error message, `ParseNameVersion` with `@` | schema change touches every command |
| 1 | Extend `internal/git` (fetch, archive, hardening, redaction); git and path transports; `--from`; `rsk source`, `rsk resolve`; implicit sources with trust prompt; the local overlay (5.7) and its adapter checks; read-only MCP tools (`list_sources`, `test_source`, `scan_source`, `resolve_reference`) and `add_source`/`remove_source` | bare-repo fixture, private repo over SSH, offline cache hit, URL allowlist, `roots` traversal | new cache root |
| 2 | Public URL/tarball transport; SHA-256 in index and lock plus backfill | digest mismatch aborts, backfilled digests match the release assets, tarballs rebuild deterministically | backfill touches every published version |
| 3 | `rsk registry build`/`deprecate`; `rsk repo init`/`import` and their MCP tools; MCP/CLI parity test; `rsk-skills.toml`; docs; migration guide; replace path-based `personal` | golden index, nested-layout fixture, duplicate-name failure | publishing workflow rewrite |

Ordering and exit criteria:

- Two items are independent of the chain and cheap, so they need not wait: SHA-256 in the index and
  lock with the digest backfill (phase 2's integrity half; it fixes a gap in the existing registry
  path), and `rsk registry deprecate` (the personal-skills move needed it and it was done by hand).
- The trust rules (section 3) ship **in the same release as** the git transport and implicit
  sources. Cloning arbitrary repositories without them is not an acceptable intermediate state.
- Until phase 1 lands, a private skills repository is used by hand (a symlink into the tool's
  skills directory) or through the existing local mode; this is stated in that repository's README.
- A phase is done when its tests are green, `task lint` and `registry:check` pass, v1 files still
  round-trip, and the behaviour of the default registry is byte-identical for existing projects.

Each phase leaves `rsk` usable; phase 0a with only the existing two sources is a no-op
behaviourally, and phase 0b writes v1 files until a new source appears.

## Risks

- **Credential exfiltration through a project file.** A committed `[[source]]` pointing at an
  attacker URL could make git offer the user's keys to a hostile host or pull hostile instructions.
  Mitigation: trust rules 1 and 2
  (no credential field exists, and a project source needs a per-machine trust grant).
- **Git as an execution surface.** `ext::` URLs, option injection, `file://`, submodules.
  Mitigation: trust rule 3.
- **Dependency confusion.** Private source unreachable, public one answers. Mitigation: rules 2 to 5.
- **Schema migration.** v2 must read v1 forever, and an older `rsk` hard-fails on v2. Mitigation:
  write v2 only when needed, `schema.Check` as a range, an `rsk migrate` dry-run.
- **Agent-driven setup as an attack path.** An agent that can add sources can be steered by a
  remote skill's text. Mitigation: the MCP surface carries no secret or trust parameters, every
  mutating tool forces approval, `add_source` registers as untrusted, source text is treated as data.
- **Alias or identity hijack.** A fetched index claiming a canonical identity would redirect
  trusted names. Mitigation: aliases only from the built-in table or machine config; GitHub identities
  must match their URL.
- **Committed manifests disclose private sources.** Mitigation: private-source pins and their
  digests go to gitignored local files (5.7), with `--shared` as the explicit opt-out. The skill
  names still appear in any command output the user pastes.
- **Executable content.** A remote `SKILL.md` is instruction content executed by the agent, and its
  digest lands in always-on context. Pair remote sources with the per-skill permissions work
  (`plans/2026-11-17-per-skill-permissions.md`), trust rules 2 and 7, and never auto-install from a
  newly added source.
- **Credential leakage through git.** `rsk` holds no secret, but git's stderr and URLs can carry a
  user's token if they embed one in a remote URL. Mitigation: captured and scrubbed git stderr, and
  rejecting URLs with embedded credentials (`https://user:token@host`).
- **Name collisions.** Qualification fixes resolution, not the on-disk name. Mitigation: 2.3, one
  name per scope, `--as`.
- **Monorepo cost.** Walking a large checkout on every lookup. Mitigation: `roots`/`excludes`,
  per-process cache, `rsk-skills.toml`.
- **False privacy.** `personal: true` in a public index stays a filter. Mitigation: rename the
  concept, document that only `private` sources and authenticated transports enforce anything.
- **Locked commit disappears.** Force-push or branch deletion. Mitigation: pin tags or protected
  branches, fail loudly naming the commit.
- **Redirect fragility.** Release URLs and published tarball URLs use the old owner and depend on
  GitHub's redirect, which breaks if a new repository ever takes the old name. Mitigation: canonical
  owner in the generator, and the lock pins the resolved URL.
- **Windows symlinks.** Unchanged; installing from a remote source does not fix it.

## Decisions

Settled 2026-11-21, each with the answer given.

| # | Question | Decision |
|---|---|---|
| 1 | `personal` naming | Keep `personal` for published-but-author-specific, add `private`; the two flags stay distinct. |
| 2 | Private registry or git | Git only. There is no private registry, by design: private sources are git repositories or paths, registries stay public and unauthenticated, and `rsk` holds no credential. |
| 3 | Repo granularity | No limit: a user can register as many private repositories as they want. One per privacy class is guidance, not a rule. |
| 4 | Signature verification | Deferred to its own ADR. The integrity records are stored as below. |
| 5 | Subtree privacy | Not supported. A client that can clone a repository reads all of it; splitting repositories by privacy class is the answer, because the only enforcing alternative is a per-principal registry, which decision 2 rules out. |
| 6 | Privacy flag or source field | A source field only; registering the source is the opt-in to install. |
| 7 | Dangling aliases | Deleted. The three `uru-*` entries are gone from `registry/aliases.toml` and the index. |
| 8 | Non-GitHub hosts | No `@host/...` form. Use `rsk source add` with a declared identity; `@owner/repo` always means GitHub. |
| 10 | Committed manifests disclose private sources | Warn, and keep private-source pins and digests in gitignored local files (5.7). |
| 11 | Git versions | Tags `<name>@v<version>` when present, `*` otherwise; ranges deferred. |

Where the integrity records live (decision 4): the **SHA-256** is published in the registry index
and copied into the project's `rsk.lock`; the **commit and tree ids** of a git skill are in
`rsk.lock` (or `rsk.local.lock` for a private source); extracted content and clones sit in the
local cache under `~/.config/rsk/cache`; and the per-machine trust grants sit in
`~/.config/rsk/trust.json`. Nothing about verification is held on a server other than the registry
index itself.

The pin-budget question from the first draft was dropped: it is a choice in the private repository,
not part of this design.

## Open questions

None. Subtree privacy (question 5) was settled by ruling out private registries (decisions 2 and 5).

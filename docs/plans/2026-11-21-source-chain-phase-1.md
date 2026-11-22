---
status: active
project: ralvaskills
date: 2026-11-21
related:
    - designs/2026-11-21-private-and-remote-skill-sources.md
    - adr/0003-source-chain-and-reference-grammar.md
    - backlog/active/9d3mhw-git-and-path-transports-rsk-source-and-trust-rules.md
---

# Phase 1: git and path sources, `rsk source`, and the trust rules

Phases 0a and 0b are done: one factory and chain, the `@owner/repo/name` grammar, and v2 manifests written only when needed. This plan breaks phase 1 of [the design](../designs/2026-11-21-private-and-remote-skill-sources.md) into steps that each leave `rsk` usable and safe to ship.

## Scope

In: hardened git access, source declarations and the per-machine trust grant, the `git` and `path` transports, the official cache on git, `rsk source`, `rsk resolve`, `--from`, implicit sources, the `private` flag with the local overlay, and the MCP tools for sources.

Out: SHA-256 digests and URL transports (phase 2), `rsk registry build`, `repo init/import`, `rsk-skills.toml` and the MCP/CLI parity test (phase 3), signature verification (own ADR), version ranges for git sources.

## Rules for the order

1. **The trust gate sits in the fetch path, not in a command.** An untrusted git source cannot be cloned whatever the entry point (CLI, lock restore, MCP), so any intermediate commit is safe. The design requires the trust rules to ship with the git transport; this makes that true commit by commit.
2. **`private` is rejected until the overlay exists.** `source.New` already refuses `Private`, `Rev`, `Roots` and `Excludes`. Each is accepted only in the step that honors it, so a source is never believed private while its pins still land in committed files.
3. **No CLI release is tagged before step 1f.** Steps 1a to 1e may land on `main` (the CLI publishes on tag, not on push), but a release carrying git sources without implicit-source confirmation is not acceptable. The first release is 1a to 1f; 1g and 1h follow (decision 4).
4. **The registry and the local clone are untouched.** Their code paths stay as they are, and `rskgit.Pull/Clone/Behind` keep forwarding output for the local clone. The official cache is the one existing source that moves (step 1d, decision 2).

## Steps

### 1a. Hardened git runner (`internal/git`, no consumer yet)

- `RemoteURL` type with a parser that accepts `ssh://`, `https://` and scp-like `git@host:path`, and rejects `ext::`, `file://`, a leading `-` and embedded credentials (`https://user:token@host`). The runner takes only this type.
- One runner for every new call: `GIT_TERMINAL_PROMPT=0`, `GIT_SSH_COMMAND` with `-o BatchMode=yes`, stdin closed, `-c protocol.ext.allow=never -c protocol.file.allow=never`, `--` before the URL, no submodule recursion, stderr captured and scrubbed instead of forwarded.
- `Fetch(ref or sha)` into a bare clone, `RevParse` (commit, and tree of a path), `ListTree(commit, roots)` returning mode and path, and `Archive(commit, paths)` into a directory with the extraction guards (regular files and directories only, paths confined, size cap, no symlink or submodule entries).
- Tests: table tests for the URL parser and the argument lists (no process spawned). Integration tests use a local bare repository through an unexported, in-package test hook that lifts the `file` restriction, since the public API cannot reach a local path. The hook cannot be set from outside the package.

### 1b. Source declarations and trust (`internal/sources`, `internal/trust`, `cmd/rsk/source*.go`)

- Sources are declared in a new machine-wide file, `~/.config/rsk/sources.json`, with its own schema version (decision 3). `config.json` is not touched, so an older `rsk` pinned in some project neither fails on it nor rewrites it, and no source can be dropped silently. Project-scoped `[[source]]` in `rsk.mod` is not part of phase 1: sources are global, which also keeps trust rule 2's project-declared case out of this release.
- `trust.json` records grants keyed by identity **and** URL. `rsk source trust <id>` writes one and refuses without a TTY, naming the exact line to run. `untrust` removes it.
- `rsk source add | list | remove | test`, always machine scope. `add` derives the SSH URL from `@owner/repo`, requires the identity to match the URL for GitHub hosts, reserves `@local`, and registers the source as **untrusted**. `--dry-run --output json` on every mutating command.
- Identity rules: a declared id must match the URL for GitHub hosts; changing a URL revokes its grant. The built-in official source (step 1d) carries a baked-in grant for its exact URL.
- Still no git transport: `source.New` has no `git` kind yet, so this step only manages declarations.

### 1c. Git transport for public sources (`internal/source/git.go`)

- Kind `git`, honoring `Rev`, `Roots` and `Excludes`. One bare clone per source at `<cache>/git/<owner>/<repo>.git`, updated with `Fetch`.
- Discovery uses `ListTree`, not a checkout: apply `roots` and `excludes`, find `SKILL.md` directories, reject symlink entries, and fail on two skills with one name (the first-match behaviour of `FindByName` stays only for the old sources). Then `Archive` only the chosen skill directories into `<cache>/sources/<owner>/<repo>/<name>/<commit>/`.
- The fetch path checks the trust grant first and fails closed with the `rsk source trust` hint.
- Fail-closed name index (resolution rule 4): a per-source name index with a fetch time, 24 h TTL, refreshed now when stale, `--offline` accepting any existing index and saying so. A source that has never been fetched, or whose refresh fails past the TTL, fails a bare-name lookup.
- `catalog` lists a git source and degrades per source: an unreachable or untrusted source is a warning and the rest are listed. (`Chain.All` is still not needed: the catalog walks the chain's sources one by one to keep its per-source warnings.)

Built as planned, with these settled while building: declared **path** sources are left out of the chain until step 1f, since a path source needs its own walker; a **project** install of a git skill is refused until step 1e can lock it, while `--global` works (it only links); `install` from `rsk.mod`, `update` and `new` keep using the built-in sources, so they never depend on a declared source; a skill's `Source` holds the repository identity for a declared source (`skill.DeclaredSource`), which keeps `Skill` under the linter's copy threshold. First real run, against `cursor/plugins`: it holds `pr-review-canvas` in two plugin directories, so the duplicate-name rule refused the whole source until `--roots pr-review-canvas` narrowed it. That is the design working, and it argues for `rsk source scan` (step 1h) to show the clashes before a user hits them.

> **Order changed while building.** Step 1e was built before 1d. Folding official onto per-commit extractions means an installed link must be re-pointed on update, which needs the lock to record a commit and tree and `update` to compare them: that is 1e. Done in that order, official updates never regress. Step 1e is complete; 1d follows and rests on it.

### 1d. The official cache on the git transport - done

Decision 2: `@anthropics/skills` becomes an ordinary git source instead of its own kind, as the design already describes (2.1).

- A built-in spec (type `git`, the existing URL, trusted for that exact URL) replaces `KindOfficial`; `Official`, `FindCached` and the lazy-clone logic from 0a go away, and the chain's cache-first pass is deleted with them.
- Built with these refinements: the official source keeps `skill.SourceOfficial` as its kind instead of becoming a declared identity, which kept bundles, labels and lock keys untouched; a legacy global link into the old clone is recognised by its target; tests reach the official source through an in-memory clone (`internal/source/sourcetest`) because they were silently reaching GitHub otherwise.
- Its lock entries gain `commit`, `tree` and `path` like any git skill, so `update` finally detects drift for official skills instead of treating them as versionless.
- **The old working clone is left in place.** Installed official symlinks point into `<official_cache>/skills/<name>`; leaving that directory alone means the 26 projects that use it keep working with no relink. New installs use the new layout, and the first `rsk update` finds no recorded `tree`, treats the skill as changed and re-points its links. Nothing deletes the old clone automatically.
- `rsk update --include official` fetches through the new path; the `official_cache` config key stays for the old clone and is read only to find it.
- Bare-name resolution still queries every source, and an unreachable official source now fails closed like any other (today it is only consulted when nothing else has the name). That is a behaviour change to call out in the changelog.

### 1e. Lock, mod and update for git sources (`internal/manifest`, `cmd/rsk`) - done

- Lock v2 entries gain `url`, `rev`, `commit`, `tree` and `path`; keys use the qualified form; restoring from a lock names a URL and still goes through the trust gate (trust rule 8) but never the chain.
- `update` compares the skill directory's `tree` id, not `commit`, so an unrelated commit in a monorepo does not report every skill outdated. `status` verifies the cached path against the recorded commit offline.
- `name@1.2.0` resolves through a tag `<name>@v1.2.0` and is an error when the tag is missing; `*` or an exact tag only.
- Unreachable locked commit (force-push) fails naming the commit.

### 1f. Implicit sources, `--from`, `rsk resolve` - done

- `rsk install @owner/repo/name` for an unregistered source: shows the URL and transport and asks; `--yes` does not cover it; CI passes `--trust @owner/repo`; nothing is saved without `--save`. SSH first, HTTPS as fallback, and the lock stores identity and commit so each machine picks its own transport.
- Pinning a skill from a non-default source shows its `DIGEST.md` text and asks (trust rule 7), since that text enters always-on context.
- `rsk install --from ./path/to/skill` as a one-off path source, not persisted, not reproducible; `rsk resolve <ref>` explains what a reference resolves to and why a bare name is ambiguous.
- **First releasable point.** Steps 1a to 1f together satisfy the "trust rules ship with the git transport" requirement.

Built as planned, with these settled while building:

- **One-run approval goes through the gate, not around it.** `sources.ChainOptions.Implicit` holds the entries approved for this run (by `--trust`, or by the person at the terminal). `ChainAt` validates each one like a declared source (GitHub only, identity equals URL, not built in, not already declared) and grants it in memory on top of `trust.json`, so the git source's own gate is the one that lets it through and nothing is written. The MCP chain never sets it, and a test pins that `install_skills` refuses an unregistered source without reaching for it.
- **Only a person can approve.** `askHuman` never proceeds on its own: `term.Confirm` auto-approves with `--yes`, without a terminal and inside the session, which is right for "Proceed?" and wrong here, so the implicit and digest prompts use their own helper. The terminal check is a seam (`isInteractive`) so the tests can drive both sides.
- **`--trust` covers an implicit source only.** For a declared source it is refused with the `rsk source trust` line, and an id nothing uses is an error. It is still a flag a shell can pass, so an agent that runs `rsk` through a shell tool can approve a source for one run; the control there is the harness's permission prompt and the digest rule, not rsk. Recorded as residual risk.
- **`--save` writes the entry only, never `trust.json`, even with `--trust`.** `--trust` needs no terminal, so letting it persist a grant would let a script make the grant that only `rsk source trust` at a terminal may. A registered but untrusted source fails every bare-name lookup (fail closed), which the `--save` output says, with the command that undoes it. This is stricter than the suggestion to write the grant when `--trust` is also given.
- **SSH first, HTTPS second, decided in `internal/sources` (`Candidate`, `Choose`).** After approval a probe (`git ls-remote`, injectable) tries SSH and then HTTPS; every URL passes the gate before it is contacted, and the approval shows both. Any SSH failure counts as unavailable, since BatchMode makes a missing key, an unknown host key and a refused connection look the same. The chosen URL is what the lock and `--save` record. `--offline` refuses an implicit source.
- **Restore: the machine decides the transport.** A lock URL that differs from this machine's is followed when both name the same GitHub repository (so SSH on one machine and HTTPS on another both work) and refused otherwise, which keeps the earlier refusal of a lock naming another host. A lock naming an undeclared source goes through the same approval as an install (`rsk sync --trust`); it never opens a source on its own.
- **Digest confirmation: `--yes` covers it, a missing terminal does not.** The text is always shown first. The non-default test is "from a git source", which includes the official one, matching the qualified-key rule. `rsk pin` gained `--yes` for it. Not covered: `rsk update` rewrites `.rsk/SKILLS.md` through `digest.Sync` when a pinned skill's digest changes, with no confirmation, and the MCP `install_skills` pin shows nothing (1h). `digest.Load` also gives a versionless skill the line "no digest, and not readable here" instead of its description (`skill.ReadFrontmatter` requires a version), so for a versionless skill, official ones included, the confirmation shows that line; fixing it changes generated files, so it is left for a decision.
- **`--from` takes one skill directory** (a real directory, not a link), links it and records nothing. It combines with `--global`; with no `--global` it links into the project's tool directories, with no `rsk.mod` or `rsk.lock` entry. A directory of skills is pointed at `rsk source add`. `install.ScopeOneOff` is the link-only scope.
- **A declared path source is `Spec.Declared`.** The local clone and a declared path are both `KindPath`; the flag tells the walker apart (`<path>/skills` versus as given), and `Chain.Primary` is the registry or the clone, which is what a bundle member resolves from, so a path source holding the same name cannot make a bundle ambiguous. Its skills have kind `path` (`skill.SourcePath`), never reach a lock, and `install.Apply` refuses them in a project.
- **`rsk resolve` uses `Chain.Survey`**, which asks every source and keeps each answer where `Find` stops at the first it cannot consult, and applies the same verdict rules; a test checks that it never says "resolved" where `Find` fails. It exits non-zero unless the reference resolves or is a bundle. It is registered outside `setup.go`'s `AddCommand` list like `rsk source`, so the landing page's command count, which `gen-data.mjs` reads from that list, does not move and does not include either command; the page should say so or count them when it is next touched.
- **A real run found what the tests could not:** `--from` was read but never registered, so the command panicked while every unit test, which calls `runInstall`, passed. A test now checks that each flag the commands read is registered and runs the whole command once.

### 1g. `private` and the local overlay - done

- Honor `Private`: excluded from `catalog`, MCP resources, `rsk://catalog` and `search_skills` unless `--include private`.
- Pins from a private source go to gitignored `.rsk/rsk.local.mod`, `rsk.local.lock` and `SKILLS.local.md`, merged with the committed files on read; the local files are appended to `.git/info/exclude`; `--shared` opts in to the committed files.
- Adapters: Claude Code reads the local digest through `CLAUDE.local.md`; OpenCode and Codex need an equivalent that is not committed. Whether each client honors an uncommitted import is **unverified** and is the first thing this step tests. I will not change `internal/tool/*` before you approve that part.

Built as planned (the owner approved the `internal/tool/*` change on 2026-11-22), with these settled while building:

- **The clients were tested first, for real.** `claude -p` and `opencode run` in a sandbox project, with a canary string: Claude Code loads a `CLAUDE.local.md` that imports `.rsk/SKILLS.local.md`; a missing imported file, and a missing `CLAUDE.local.md`, are ignored without an error. OpenCode loads a file named in the committed `instructions` array and ignores one that is missing. In the end-to-end run the owner's checkout showed the private skill and a clone with only the committed files showed `NONE`, in both clients. Codex has no `@` import in `AGENTS.md` (`AGENTS.override.md` replaces the file, it does not add to it), and rsk writes no digest for Codex or the eight JSON clients, so they need no refusal: there is nothing to leak. Codex itself was not run.
- **Privacy is a property of the source, not of `skill.Skill`.** Adding a field would have taken `Skill` past the 128-byte copy gate. `Spec.Private` and `sources.Entry.Private` carry it; the install path asks `PrivateSources` (identities) and the manifest carries a per-skill layer mark (`Mod.Local`, `LockEntry.Local`).
- **The reader merges, the writer splits.** `ReadMod` and `ReadLock` merge the local files and mark what came from them; `WriteMod` and `WriteLock` send each skill back to its layer. That is what makes the invariant hold for `uninstall`, `update`, `pin`, `sync` and `rename` without touching them: none can write a private skill to a committed file. A name in both layers is an error.
- **What is excluded is more than the local files.** The first real run showed the install links (`.claude/skills/<name>`, one per client) are named after the skill and would have been committed, so they are in `.git/info/exclude` with the local files. The exclude file is found by reading `.git` (a worktree's `commondir` included), with no git binary.
- **A client that cannot import gets nothing, which is the refusal.** `tool.LocalImporter` is implemented by Claude Code and OpenCode only.
- **The MCP server cannot see a private source at all** (`ChainOptions.ExcludePrivate`, `Env.HidePrivate`), which is stricter than the design's "unless the manifest already lists it": a private source is used from the CLI. `project_profile` and its install counts leave the local layer out.
- **Not covered:** an implicit source cannot be private (nothing says it is), so its pins are committed; the install path warns in the changelog only. A project whose `rsk.mod` is committed with private skills from before the flag existed keeps them there until reinstalled. Windows path handling of the exclude file is untested.


### 1h. MCP tools for sources - done

- Read-only: `list_sources`, `test_source`, `scan_source`, `resolve_reference`. Mutating with forced approval: `add_source`, `remove_source`. `project_profile` gains a `sources` section. `install_skills` refuses unknown or private sources unless the manifest already lists them.
- No trust or `--yes` parameter over MCP; `add_source` registers as untrusted; tool output from a source is data. Every mutating tool takes `dry_run`.
- The CLI/MCP parity test (a phase 3 item in the design) lands here, since these tools ship here.

Built as planned, with these settled while building:

- **`internal/sourceops` is the core** (`List`, `PlanAdd`, `PlanRemove`, `Test`, `Scan` over plain structs and an `Env` that holds the file paths and the test seams). A mutating operation is a plan plus the function that applies it, so the CLI previews and asks, and the MCP tool applies unless `dry_run`. `rsk source list|add|remove|test` were moved onto it with their tests unchanged.
- **`scan` needed a lenient discovery.** Resolving a source still fails on a name two directories claim; `Path.Scan` and `Git.Scan` report the clash instead, leave the clashing names out of the skills, and add frontmatter issues (a name that differs from the directory, no version, no description). Scanning a repository fetches it, so it needs a registered, trusted source; an unregistered `@owner/repo` is refused with the `rsk source add` line, since the alternative is cloning code nobody approved.
- **`install_skills` resolves through the declared sources behind the trust gate.** The design says it refuses unknown or private sources unless the manifest lists them; the server now builds its chain with `sources.ChainAt` (declared and path sources, never `Implicit`), so an unregistered source is still refused and a trusted one installs. A side effect is the CLI's: a registered but untrusted source fails a bare-name install over MCP too (fail closed). `Deps.SourcesPath` empty means no declared sources, which keeps the tests that do not exercise them off the real files; the command sets both paths.
- **The schema test is the guard against a trust parameter**: it fails when any tool's input has a property whose name contains trust, yes, token, secret, password, credential, auth, key, private, confirm or force.
- **Parity.** `mcp.ToolNames()` lists every tool (a test checks it against what the server registers), and `TestEveryMCPToolHasACLICounterpart` maps each to a command in the tree and requires `--dry-run` on the twins of the mutating tools. `project_profile` is the one agent-only tool, listed with that reason.
- **The MCP pin still shows no digest.** `install_skills` with `pin` writes `.rsk/SKILLS.md` for a non-default source without the text confirmation the CLI makes (trust rule 7). The client's forced approval shows the input, not the digest. Left open in `g1zk6n`.
- **Not done here:** `AGENTS.md` still calls the `rsk mcp` tools "planned" (editing it needs your approval) and `rsk-guide` does not know the new tools or commands (it publishes on push).

## Decisions (settled 2026-11-21)

1. **Fetch strategy: list the tree, archive only the chosen skills.** No checkout; symlinks and submodules are rejected from the listing; `tree` is a direct `rev-parse`.
2. **Official cache: folded into the git transport now** (step 1d), with the old clone left in place so installed links keep working.
3. **Declared sources live in `~/.config/rsk/sources.json`**, machine-wide, with its own schema version. `config.json` stays v1 and untouched.
4. **Release cut: one release after step 1f**; `private` with the overlay (1g) and the MCP tools (1h) follow. Until 1g, `--private` is refused, so a private skills repository can be registered only as a normal source and its pins land in the project's committed `.rsk/` files: acceptable in a private project, not in a public one.

## Risks

- **First network-facing code in `rsk`.** Mitigated by the allowlist, the non-interactive runner and the fail-closed trust gate being the first two steps, before any transport exists.
- **Tests that need a remote.** Private-repo-over-SSH cannot run in CI without an sshd. The automated suite uses local bare repositories; the SSH path is a manual acceptance run against your real private skills repository, which needs your credentials and your go-ahead.
- **`git` version.** Fetching by SHA and `git archive` against a bare clone assume a recent `git`; the runner checks the version once and reports it.
- **Cache growth.** One bare clone per source plus per-commit extractions; `rsk source remove` deletes both, and a later step can prune old commits.
- **Moving the official cache.** The only existing source whose layout changes. Leaving the old clone in place and re-pointing links on the next `update` avoids a forced relink; the risk is a long tail of projects whose official links point at a clone nobody updates, which `rsk status` should report.
- **A large step 1g.** It is the only step touching adapter config-writing, so it is isolated, approved separately and can slip without blocking the first release.

## Exit criteria

- [x] An untrusted git source cannot be cloned from the CLI, a lock restore or MCP.
- [x] `rsk source add @owner/repo`, `trust`, `install @owner/repo/name` installs from a local bare-repo fixture and records `commit`, `tree` and `path` in a v2 lock.
- [x] A second machine restores from that lock without the chain and without a secret.
- [x] A source that is down fails a bare-name lookup instead of falling through to another source.
- [ ] URL allowlist, option injection, embedded credentials and `roots` traversal tests pass.
- [ ] `update` reports a skill outdated only when its own directory changed.
- [ ] Registry and local clone behaviour is byte-identical for existing projects; an existing official install keeps working untouched and is re-pointed by the next `update`; `task lint`, `go test ./...`, `task registry:check` and `task site:check` pass.
- [x] The manual SSH acceptance run against the private skills repository passes (2026-11-21, see `6w4kzn` in done; one finding, the runner honors `url.*.insteadOf`).

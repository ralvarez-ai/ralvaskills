---
status: draft
project: ralvaskills
date: 2026-11-21
related:
    - designs/2026-11-21-private-and-remote-skill-sources.md
    - adr/0003-source-chain-and-reference-grammar.md
    - backlog/pending/9d3mhw-git-and-path-transports-rsk-source-and-trust-rules.md
---

# Phase 1: git and path sources, `rsk source`, and the trust rules

Phases 0a and 0b are done: one factory and chain, the `@owner/repo/name` grammar, and v2 manifests written only when needed. This plan breaks phase 1 of [the design](../designs/2026-11-21-private-and-remote-skill-sources.md) into steps that each leave `rsk` usable and safe to ship.

## Scope

In: hardened git access, source declarations and the per-machine trust grant, the `git` and `path` transports, `rsk source`, `rsk resolve`, `--from`, implicit sources, the `private` flag with the local overlay, and the MCP tools for sources.

Out: SHA-256 digests and URL transports (phase 2), `rsk registry build`, `repo init/import`, `rsk-skills.toml` and the MCP/CLI parity test (phase 3), signature verification (own ADR), version ranges for git sources.

## Rules for the order

1. **The trust gate sits in the fetch path, not in a command.** An untrusted git source cannot be cloned whatever the entry point (CLI, lock restore, MCP), so any intermediate commit is safe. The design requires the trust rules to ship with the git transport; this makes that true commit by commit.
2. **`private` is rejected until the overlay exists.** `source.New` already refuses `Private`, `Rev`, `Roots` and `Excludes`. Each is accepted only in the step that honors it, so a source is never believed private while its pins still land in committed files.
3. **No CLI release is tagged before step 1e.** Steps 1a to 1d may land on `main` (the CLI publishes on tag, not on push), but a release carrying git sources without implicit-source confirmation is not acceptable.
4. **Old behaviour is untouched.** The registry, local clone and official cache keep their code paths; `rskgit.Pull/Clone/Behind` keep forwarding output as today. Only new code uses the hardened runner.

## Steps

### 1a. Hardened git runner (`internal/git`, no consumer yet)

- `RemoteURL` type with a parser that accepts `ssh://`, `https://` and scp-like `git@host:path`, and rejects `ext::`, `file://`, a leading `-` and embedded credentials (`https://user:token@host`). The runner takes only this type.
- One runner for every new call: `GIT_TERMINAL_PROMPT=0`, `GIT_SSH_COMMAND` with `-o BatchMode=yes`, stdin closed, `-c protocol.ext.allow=never -c protocol.file.allow=never`, `--` before the URL, no submodule recursion, stderr captured and scrubbed instead of forwarded.
- `Fetch(ref or sha)` into a bare clone, `RevParse` (commit, and tree of a path), `ListTree(commit, roots)` returning mode and path, and `Archive(commit, paths)` into a directory with the extraction guards (regular files and directories only, paths confined, size cap, no symlink or submodule entries).
- Tests: table tests for the URL parser and the argument lists (no process spawned). Integration tests use a local bare repository through an unexported, in-package test hook that lifts the `file` restriction, since the public API cannot reach a local path. The hook cannot be set from outside the package.

### 1b. Source declarations and trust (`internal/config`, `internal/trust`, `cmd/rsk/source*.go`)

- `config.json` gains `sources`. A file with sources is written as config v2 and an older `rsk` refuses it, so it can never save the file and silently drop them. A machine with no sources stays v1.
- `trust.json` records grants keyed by identity **and** URL. `rsk source trust <id>` writes one and refuses without a TTY, naming the exact line to run. `untrust` removes it.
- `rsk source add | list | remove | test`, machine scope first. `add` derives the SSH URL from `@owner/repo`, requires the identity to match the URL for GitHub hosts, reserves `@local`, and registers the source as **untrusted**. `--dry-run --output json` on every mutating command.
- Identity rules: a project source with the same id as a machine source must have the same URL; changing a URL revokes its grant.
- Still no git transport: `source.New` has no `git` kind yet, so this step only manages declarations.

### 1c. Git transport for public sources (`internal/source/git.go`)

- Kind `git`, honoring `Rev`, `Roots` and `Excludes`. One bare clone per source at `<cache>/git/<owner>/<repo>.git`, updated with `Fetch`.
- Discovery uses `ListTree`, not a checkout: apply `roots` and `excludes`, find `SKILL.md` directories, reject symlink entries, and fail on two skills with one name (the first-match behaviour of `FindByName` stays only for the old sources). Then `Archive` only the chosen skill directories into `<cache>/sources/<owner>/<repo>/<name>/<commit>/`.
- The fetch path checks the trust grant first and fails closed with the `rsk source trust` hint.
- Fail-closed name index (resolution rule 4): a per-source name index with a fetch time, 24 h TTL, refreshed now when stale, `--offline` accepting any existing index and saying so. A source that has never been fetched, or whose refresh fails past the TTL, fails a bare-name lookup.
- `catalog` lists a git source and degrades per source with a partial marker; `Chain.All` appears here because it now has a caller.

### 1d. Lock, mod and update for git sources (`internal/manifest`, `cmd/rsk`)

- Lock v2 entries gain `url`, `rev`, `commit`, `tree` and `path`; keys use the qualified form; restoring from a lock names a URL and still goes through the trust gate (trust rule 8) but never the chain.
- `update` compares the skill directory's `tree` id, not `commit`, so an unrelated commit in a monorepo does not report every skill outdated. `status` verifies the cached path against the recorded commit offline.
- `name@1.2.0` resolves through a tag `<name>@v1.2.0` and is an error when the tag is missing; `*` or an exact tag only.
- Unreachable locked commit (force-push) fails naming the commit.

### 1e. Implicit sources, `--from`, `rsk resolve`

- `rsk install @owner/repo/name` for an unregistered source: shows the URL and transport and asks; `--yes` does not cover it; CI passes `--trust @owner/repo`; nothing is saved without `--save`. SSH first, HTTPS as fallback, and the lock stores identity and commit so each machine picks its own transport.
- Pinning a skill from a non-default source shows its `DIGEST.md` text and asks (trust rule 7), since that text enters always-on context.
- `rsk install --from ./path/to/skill` as a one-off path source, not persisted, not reproducible; `rsk resolve <ref>` explains what a reference resolves to and why a bare name is ambiguous.
- **First releasable point.** Steps 1a to 1e together satisfy the "trust rules ship with the git transport" requirement.

### 1f. `private` and the local overlay (touches `internal/tool/*`: needs your approval)

- Honor `Private`: excluded from `catalog`, MCP resources, `rsk://catalog` and `search_skills` unless `--include private`.
- Pins from a private source go to gitignored `.rsk/rsk.local.mod`, `rsk.local.lock` and `SKILLS.local.md`, merged with the committed files on read; the local files are appended to `.git/info/exclude`; `--shared` opts in to the committed files.
- Adapters: Claude Code reads the local digest through `CLAUDE.local.md`; OpenCode and Codex need an equivalent that is not committed. Whether each client honors an uncommitted import is **unverified** and is the first thing this step tests. I will not change `internal/tool/*` before you approve that part.

### 1g. MCP tools for sources

- Read-only: `list_sources`, `test_source`, `scan_source`, `resolve_reference`. Mutating with forced approval: `add_source`, `remove_source`. `project_profile` gains a `sources` section. `install_skills` refuses unknown or private sources unless the manifest already lists them.
- No trust or `--yes` parameter over MCP; `add_source` registers as untrusted; tool output from a source is data. Every mutating tool takes `dry_run`.
- The CLI/MCP parity test (a phase 3 item in the design) lands here, since these tools ship here.

## Decisions I need from you

1. **Fetch strategy.** `ListTree` plus per-skill `Archive` (recommended: no checkout of a large monorepo, symlinks visible before extraction, one code path) versus a full shallow checkout and `skill.Walk` (less code, but fetches and walks everything).
2. **Official cache.** Keep it as its own kind in phase 1 (recommended: zero risk to the 26 projects that use it) or fold it into the git transport now, which the design lists as a later revisit.
3. **Config bump.** A `config.json` with `sources` becomes v2 and an older `rsk` refuses it (recommended, same pattern as the manifests) versus keeping v1 and accepting that an older `rsk` may drop sources when it rewrites the file.
4. **Release cut.** One release after 1e (recommended), with 1f and 1g as follow-ups, or hold everything for one release.

## Risks

- **First network-facing code in `rsk`.** Mitigated by the allowlist, the non-interactive runner and the fail-closed trust gate being the first two steps, before any transport exists.
- **Tests that need a remote.** Private-repo-over-SSH cannot run in CI without an sshd. The automated suite uses local bare repositories; the SSH path is a manual acceptance run against your real private skills repository, which needs your credentials and your go-ahead.
- **`git` version.** Fetching by SHA and `git archive` against a bare clone assume a recent `git`; the runner checks the version once and reports it.
- **Cache growth.** One bare clone per source plus per-commit extractions; `rsk source remove` deletes both, and a later step can prune old commits.
- **A large step 1f.** It is the only step touching adapter config-writing, so it is isolated, approved separately and can slip without blocking the first release.

## Exit criteria

- [ ] An untrusted git source cannot be cloned from the CLI, a lock restore or MCP.
- [ ] `rsk source add @owner/repo`, `trust`, `install @owner/repo/name` installs from a local bare-repo fixture and records `commit`, `tree` and `path` in a v2 lock.
- [ ] A second machine restores from that lock without the chain and without a secret.
- [ ] A source that is down fails a bare-name lookup instead of falling through to another source.
- [ ] URL allowlist, option injection, embedded credentials and `roots` traversal tests pass.
- [ ] `update` reports a skill outdated only when its own directory changed.
- [ ] Registry, local clone and official behaviour is byte-identical for existing projects, and `task lint`, `go test ./...`, `task registry:check` and `task site:check` pass.
- [ ] The manual SSH acceptance run against the private skills repository passes (needs your approval to run).

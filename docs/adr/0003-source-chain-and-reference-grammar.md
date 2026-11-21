---
adr: 3
status: accepted
date: 2026-11-21
project: ralvaskills
---

# One source chain, one reference grammar, and no credential in rsk

## Context

`rsk` resolves a skill through exactly two sources: a local clone **or** the hosted registry, then the `anthropics/skills` cache. Which source answers is decided at about a dozen call sites that each construct sources directly (`source.NewOfficial`, `NewRegistry`, `NewLocal`), and `skill.Source` is a closed enum of three kinds. Three consequences follow, set out in [the design](../designs/2026-11-21-private-and-remote-skill-sources.md):

- A project cannot use the public registry and a private skill repository at the same time, and a skill that lives only in a git repository cannot be installed.
- The lock records name, version and a source *kind*, keyed by bare name. Two publishers shipping the same skill name collide silently, and nothing records where a remote skill came from.
- Nothing in `rsk` carries a credential, so any private source has to be authenticated by something outside it.

The seven personal skills have already moved to a private git repository that `rsk` cannot yet install from. The decisions below were settled on 2026-11-21 (design, "Decisions").

## Options considered

- **Private registry** — authentication on the registry client, a second URL. Needs a bucket, an auth front, a publisher and token distribution; puts a secret inside `rsk`; still offers no git-ref installs.
- **Source chain, git-first** — a typed source descriptor, one factory, an ordered chain, a git transport, a lock that records provenance. Private access becomes git's problem. Chosen.
- **Vendor or submodules only** — no code, but no resolution, no updates and no public/private separation. Remains possible as a degenerate path source.
- **OCI artifacts** — good provenance, but tooling and auth cost with no skill-specific benefit yet.

## Decision

**1. A source chain.** A source is a `Spec` (`ID`, `Type`, `URL`, `Rev`, `Roots`, `Excludes`, `Private`) built by one factory, `source.New`, and held in an ordered `Chain`. No other code constructs a source, and the concrete constructors stop being exported so the compiler enforces it. There is no `priority` field: duplicates are an error, so an ordering could never decide a resolution; declaration order only keeps listings stable. Resolution:

1. A qualified reference resolves in exactly that source; an unknown source is an error, not a fallback.
2. A bare name queries every source and resolves only if exactly one has it.
3. Two sources with the same bare name is an error naming both.
4. **Fail closed.** If a source cannot be consulted while resolving a bare name, the command fails, so an expired credential on a private source can never fall through to the public one.
5. **Lock first.** Installing from `rsk.lock` never consults the chain.
6. Adding a source never changes how an already-recorded bare name resolves.
7. A `private` source stays out of listings, MCP resources and search unless asked for.

**2. A reference grammar with one canonical identity.** `@<owner>/<repo>/<name>[@<version>]` is qualified; `<name>[@<version>]` is bare. A source has one **canonical** identity and may have **aliases**: the default registry is canonically `@ralvarez-ai/ralvaskills` with `@ralvarezdev/ralvaskills` as an alias. The lock and `rsk.mod` store only the canonical form; an alias is accepted on input and normalised; aliases come only from the table built into `rsk` or the user's machine config, never from fetched data; built-in identities are checked before the implicit GitHub derivation. Keys are qualified if and only if the skill comes from a source other than the default registry, so the same committed file reads the same on every machine.

**3. Lock and `rsk.mod` v2, written only when needed.** The lock gains source, transport, URL, ref, commit, tree, path and digest. A file is written as v1 until a non-default source or a qualified key exists, so teammates on an older `rsk` are unaffected until they adopt one; `schema.Check` becomes a supported range, and an older `rsk` hard-fails on a v2 file with the existing "upgrade rsk" message.

**4. Private skills are git or path sources only, and `rsk` holds no credential.** Authentication is git's (SSH key, credential helper) or the filesystem's. Registries are public and unauthenticated, and a `private` registry is rejected. There is no `auth`, `token` or `env:` field, no keychain, and no secret in any flag or MCP parameter. Trust over a project-declared source is a per-machine grant that only a human at a TTY can make.

The work lands in phases (0a refactor with no behaviour change, 0b grammar and schema, then git transport, digests, publishing). Signature verification is deferred to its own ADR.

## Consequences

- Easier: a project can mix the public registry, private git repositories, paths and the official cache; the source of an install is recorded and verifiable; adding a source kind is one factory case, not a dozen edits.
- Easier: no secret handling, rotation or storage to build or audit in `rsk`.
- Harder: an old `rsk` cannot read a v2 file, and every command that reads the lock changes with 0b. Bare-name resolution now errors where a name exists in two sources; today the registry silently wins over the official cache, so that case changes from a silent pick to an error.
- Harder: a private registry is permanently out of scope. Per-skill access control means splitting repositories by privacy class.
- Harder: remote skill text reaches always-on agent context, so the trust rules in the design's section 3 ship in the same release as the git transport, not after it.
- Revisit: signature verification (own ADR), version ranges for git sources, and whether `official` should become an ordinary `git` source once the git transport exists.

---
status: draft
project: ralvaskills
date: 2026-11-22
related:
    - plans/2026-11-21-source-chain-phase-1.md
    - backlog/pending/9m3rtc-roll-out-rsk-v2-manifests-to-the-consumer-projects.md
    - backlog/pending/c3qwq1-push-main-and-v3-10-0-and-run-the-release-workflow.md
---

# Roll-out of `rsk` 3.10.0 to the consumers

`rsk` 3.10.0 writes `rsk.mod` and `rsk.lock` as version 2 the first time a project holds a skill from a non-default source, and that includes an official skill. An older `rsk` refuses a v2 file ("upgrade rsk"). The consumers pin `rsk` in `mise.toml`, so the pin has to move before any project can become v2. This plan orders that, measures what is affected, and says what each step may and may not touch. It edits no consumer: every step that does is the owner's to approve.

## What the survey found (read-only, 2026-11-22, under `~/Dev`)

| Fact | Number |
|---|---|
| Projects with an `.rsk/rsk.lock` | 149 |
| Of those, with an official skill (they become v2 on their first `rsk update` or `rsk install` after the bump) | 21 |
| Already v2 | 0 |
| `mise.toml` files pinning `github:ralvarez-ai/ralvaskills` | 85: 79 at 3.9.1, 6 at 3.7.0 |
| Lock projects whose directory tree has no `mise.toml` pin (they run whatever `rsk` is on `PATH`) | 62 |
| Official-skill projects pinned at 3.9.1 / 3.7.0 / no pin | 12 / 1 / 8 |
| Release `v3.10.0` on GitHub | not yet (the tag is pushed, the Release CLI workflow has not run) |

The numbers move: re-run the survey in step 0 rather than trusting them.

## Why order matters

1. A v2 file written by 3.10.0 is unreadable by 3.9.1, 3.9.0 and everything before. Within one project the failure is loud and local, but across a shared repository a teammate on the old pin is locked out until they bump.
2. Nothing writes v2 on its own: only `rsk update`, `rsk install` or `rsk pin` of a non-default skill does. So the bump can be gradual and a project not touched stays readable by anything.
3. `mise exec -- rsk` resolves the pin of the directory it runs in, so inside a consumer it runs the old binary until the pin moves (the trap recorded in `r9tk2v`).
4. The 62 projects with no pin depend on the owner's global `rsk`; they are safe once that one binary is 3.10.0 and every other machine that opens them is too.

## Steps

0. **Gate: the release exists.** `c3qwq1` done: assets on the GitHub Release, `rsk --version` from the released binary prints 3.10.0. Re-run the survey (script below). Nothing else starts before this.
1. **Move the owner's own binary** to 3.10.0 (global install, and this repository's `mise.toml`). Run `rsk status` and `rsk list` in one v1 project to confirm a v1 file still reads and is not rewritten.
2. **Pilot on one official-skill project** the owner picks (a small one). With the owner's go: bump its `mise.toml` pin, run `rsk update --dry-run`, then `rsk update`, and check that `rsk.lock` is v2, that `rsk sync` restores it, and that the `.rsk` diff is only the lock, the mod and `SKILLS.md`. Commit it with explicit pathspecs, one `chore(rsk)` commit. Keep this project as the control for step 5.
3. **Bump the pins in batches**, official-skill projects first (21), then the rest, in the order the earlier sweeps used. One `mise.toml` edit per repo, committed on its own (`chore(mise): bump rsk to 3.10.0`), pushed or recorded in `w8nq2s` when a push is refused. A repository on 3.7.0 moves straight to 3.10.0.
4. **Let `rsk update` run per project only after its pin moved**, and only the projects the owner wants refreshed: the pin bump alone changes nothing about `.rsk`.
5. **Verify with the check below** after each batch; it must print nothing.
6. **The 62 unpinned projects** need no edit. Record which machines have a global `rsk` older than 3.10.0; a v2 project opened on one of them fails with the upgrade message.

## The check (read-only)

For every project with a v2 `.rsk/rsk.lock` or `rsk.mod`, the `mise.toml` pin in or above it must be 3.10.0 or newer. The survey script in this repository's history (the Python snippet in the session that wrote this plan) does the walk; the rule is the only part that matters, and a project failing it is the one to bump first. A project with no pin is reported separately as "unpinned", not as a failure.

## What may change in a consumer

Only `mise.toml` (one line) and, after an explicit `rsk update`, `.rsk/rsk.mod`, `.rsk/rsk.lock` and `.rsk/SKILLS.md`. A skill file, `AGENTS.md`, `CLAUDE.md` or `opencode.json` is never edited by this plan. If an `rsk update` wants to change anything else, stop and look.

## Risks

- **A teammate on an older pin pulls a v2 file.** Mitigation: bump every pin in a repository in the same commit that could introduce v2, and tell teammates to run `mise install`.
- **An unpinned project opened on an old global `rsk`.** Mitigation: step 6, and the release note's migration rule.
- **`rsk update` re-points official links** (they move from the old working clone to per-commit extractions). Expected, and the point of v2 for official skills; the pilot confirms it before the batch.
- **The pinned digests change text** (a versionless skill now digests as its description instead of "no digest, and not readable here"). Harmless, but it shows up in `SKILLS.md` diffs; say so in the commit body.
- **Pushes that cannot be made** (eleven repositories earlier): record them in `w8nq2s`, do not retry blindly.

## Rollback

A v2 file written by mistake is repaired by running the old `rsk` after restoring the two files from git (`git checkout -- .rsk/rsk.mod .rsk/rsk.lock`); nothing else is rewritten. The pin bump is one line to revert.

## Exit criteria

- [ ] `v3.10.0` assets exist and the owner's `rsk --version` prints it.
- [ ] The pilot project is v2, restores with `rsk sync`, and its diff is only the three `.rsk` files.
- [ ] Every repository whose `mise.toml` pins `ralvarez-ai/ralvaskills` is at 3.10.0 or newer, or is recorded in `w8nq2s`.
- [ ] The check prints nothing: no v2 project is under an older pin.
- [ ] The 21 official-skill projects have run `rsk update` and are v2, or the owner has said which ones are left on purpose.

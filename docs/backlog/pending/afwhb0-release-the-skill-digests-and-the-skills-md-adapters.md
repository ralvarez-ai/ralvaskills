---
title: 'Release the skill digests and the SKILLS.md adapters'
priority: p1
section: ralvaskills
topics:
    - release
    - registry
    - cli
---
Everything is committed on `main` and not pushed. The order matters:

1. Revert the local `registry/index.json` (`git checkout registry/index.json`); it was synced only to unblock `task site:dev` and its release URLs do not exist yet.
2. `task lint`, `task registry:check` and the Go tests on the final tree.
3. Push `main`. The publish workflow releases the 59 bumped skills, `mockup-architect` 1.0.0 and the two moved demo skills, then commits the index. `git pull` afterwards.
4. Release the CLI (`task cli:tag`): `internal/digest` and the adapter change only reach consumers with a new `rsk`.
5. In this repo run `rsk update` (or `pin`) so `.rsk/CLAUDE.md`, `.rsk/SKILLS.md` and `opencode.json` use the new format, and confirm with `claude -p` that a digest line reaches the model.

Done when the index lists 60 skills, an `rsk` release carries the adapter, and this repo's `.rsk/SKILLS.md` is generated from published digests. Supersedes the push step of `sykdrq` and the ralvaskills part of `z7ckmp`. Design: `docs/designs/2026-11-06-skill-digests-and-skills-md.md`.

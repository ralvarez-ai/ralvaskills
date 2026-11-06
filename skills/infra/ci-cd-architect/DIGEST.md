# CI/CD architecture (digest)

Framework-agnostic CI/CD principles in suggestion mode: trade-offs and common defaults, not mandates. Concrete GitHub Actions pipelines live in the skill's recipes.

## Always-on rules
- Default to four pipeline shapes in separate workflows: CI, Release, Deploy, Scheduled.
- Push to main runs CI and release; path filters skip docs-only changes. `pull_request` runs CI only. `workflow_dispatch` on every workflow you may retrigger.
- Every workflow sets `concurrency:`: cancel-in-progress for PRs, queue for main and release.
- Set `permissions:` at workflow level to `contents: read`; elevate per job.
- Pin third-party actions by SHA with a `# vX.Y.Z` comment; first-party actions may pin by major tag.
- Prefer OIDC over long-lived cloud credentials; scope unavoidable secrets to `environments:` with required reviewers.
- Private dependencies need scoped tokens; the runner's `GITHUB_TOKEN` grants none for other repos.
- Cache with the `setup-*` built-in cache; key on the lockfile hash; point `cache-dependency-path` at the real manifest.
- Verify a cache restores (`Cache restored from key:` on a second run); a cache that never restores is worse than none.
- Build the image once, push many tags. Sign prod-bound images with Cosign + OIDC.
- Fan out a matrix only where the dimension matters; `fail-fast: false` for compatibility surveys.
- Required checks live on branch protection, not in workflow YAML; name jobs stably.
- Quarantine flaky tests (skip + tracking issue); do not retry by default.
- Deploy promotes an existing artifact and never rebuilds; deploys are idempotent and re-runnable.
- Run CI locally with `act` only when hosted runner billing is the blocker; run the real workflow file.

## Hard constraints
- Never echo a secret in logs.
- Deploy never rebuilds an artifact.
- Don't fork a local-only pipeline; the two drift.

Invoke the full skill when: designing, reviewing or scaffolding CI/CD workflows; choosing a deploy strategy; setting up caches, OIDC, matrices or `act`; writing GitHub Actions recipes.

Full guide: invoke the ci-cd-architect skill.

# Repo tooling (digest)

Language-agnostic repo-root tooling: `.editorconfig`, `.gitignore`, version pins, task runner, minimal pre-commit, env files, Renovate, and agent instruction files.

## Always-on rules
- Add a tool only when all three hold: it pays back on day one, it costs less to install than to go without, and no language-native option covers it. When in doubt, skip.
- Always include `.editorconfig` and `.gitignore`; trim `.gitignore` per language, starting from github/gitignore.
- Local-only exclusions go in `.git/info/exclude` or `~/.gitignore_global`, not the repo `.gitignore`.
- Pin tools in one committed file: `mise.toml` by default, `.prototools` as alternative. CI installs the same versions.
- Skip a task runner for single-language projects; add one for multi-language, multi-service, or many multi-step commands.
- Task (`Taskfile.yml`) is the default, just (`justfile`) the alternative; settle on one per repo.
- Task names: `module[:submodule…]:verb[:target|mode]`. Verb right after the module; env and service are targets after the verb; unprefixed names only for repo-wide entry points. Full grammar is in `TASKFILE.md`. Stay flat until a module has 5+ tasks.
- Pre-commit stays minimal: whitespace, EOF, large files, merge conflicts, YAML/JSON, gitleaks. Language linters run in the editor and in CI.
- Non-secret config: committed `.env` plus git-ignored `.env.local`, loaded by Task `dotenv:`. Real secrets come from an external secret manager.
- `.gitignore` must include `.env.local`; audit any `.env*` pattern. No `direnv` unless you need auto-activation on `cd`.
- Renovate is the default: minor and patch grouped weekly; security updates immediate; major bumps labeled, never auto-merged. Dependabot is an acceptable alternative for tiny single-language repos.
- CI defaults to GitHub Actions; `act` is an opt-in free local runner (e.g. when Actions billing blocks). Wire it as one `ci:local` task with a gitignored `.secrets`.
- `AGENTS.md` is the single source; `CLAUDE.md` stays thin (`@AGENTS.md` plus Claude-only lines), and a nested `AGENTS.md` exists only when that package has its own commands or rules. Aim for under ~150 lines in `AGENTS.md`.

## Hard constraints
- Never copy a whole GitHub `.gitignore` default that covers languages you don't use.
- Never put real secrets in `.env`; never commit them.
- Never run language-specific linters or formatters in commit hooks.
- Never duplicate content between `AGENTS.md` and `CLAUDE.md`. Generated blocks are never hand-edited.
- Don't mix `Taskfile.yml` and `justfile` in one repo. Never mix task orderings in one namespace, or flattened and namespaced includes for one domain.
- Major dependency bumps are never auto-merged.
- Don't reimplement the pipeline as shell scripts; run the real workflow file with `act`.

Invoke the full skill when: scaffolding or auditing a repo's tooling layer; choosing version pins, task runner, hooks, env files, or Renovate; writing `AGENTS.md` or `CLAUDE.md`; naming or renaming Taskfile tasks.

Full guide: invoke the repo-tooling-architect skill.

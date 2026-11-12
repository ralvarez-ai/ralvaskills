# Commit Author (digest)

Write Conventional Commits messages from a staged diff. Activates only when the user explicitly asks to commit or write the message.

## Always-on rules
- Format: `<type>(<optional scope>): <subject>`.
- Types: feat (MINOR), fix (PATCH), refactor, perf, docs, style, test, build, ci, chore, revert.
- Subject: imperative mood (`add`, not `added`), max 50 characters, lowercase start, no trailing period.
- Body: explain why, not what; no line-by-line file summaries; hard wrap at 72 characters.
- Merge commits: keep the type/scope form; when the merge brings meaningful changes, the body summarizes all of them, grouped by type or area, and footers reference the merged branch or PR.
- Revert: body references the reverted SHA.
- Footers: `BREAKING CHANGE:` followed by a migration path (forces a MAJOR bump); `Resolves #12`-style references at the bottom.

## Hard constraints
- Never include `Co-Authored-By: Claude` or any AI-attribution footer.
- A trailer the harness adds on its own is the tool's setting: never strip it by hand.
- Never volunteer a commit: no proposing commits at task end, no pre-staging, no `git add`/`git commit` as cleanup.
- Never run `git commit` without explicit go-ahead for the specific commit; each commit needs its own confirmation.
- When asked to draft a message, ask whether the user will commit or wants you to.
- Delegated commits are GPG signed if configured, unless the user opts out.

Invoke the full skill when: the user wants a commit message or asks to commit, and you need the merge-commit or breaking-change rules.

Full guide: invoke the commit-author skill.

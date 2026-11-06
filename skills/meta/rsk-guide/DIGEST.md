# rsk-guide digest

Operator's summary of the `rsk` CLI: manage ralvaskills skills via the `.rsk/` manifest, bundles, global installs, and the `rsk mcp` server.

## Always-on rules

- Treat install, pin, update and uninstall as architectural decisions, not housekeeping.
- Before any state-changing action, state the action plainly: which skills or bundles, target scope (`--global` / `--for` / project manifest), and what changes on disk (including `.rsk/CLAUDE.md` and `./CLAUDE.md` when pinning).
- Quote the exact command, then wait for the user to confirm, redirect or decline.
- Use `--dry-run` when available if the user wants to preview.
- Confirm each bundle and each pin separately, even after a blanket "yes".
- Pin only skills that truly need auto-loading; installed-but-unpinned skills are not auto-loaded.
- Read-only commands (`rsk list`, `rsk catalog`, `rsk status` without `--stack`) need no confirmation.
- Verify flags against `rsk --help`; the CLI is still in progress (draft).

## Hard constraints

- Never run `rsk install`, `rsk uninstall`, `rsk update`, `rsk new`, `rsk destroy`, `rsk pin` or `rsk unpin` without explicit user confirmation.
- Ask first before `rsk status --stack` (network call). It is not yet implemented.
- Never promise drift checks until `--stack` ships.
- `rsk` must be on PATH for the `rsk mcp` server; `rsk new` skips MCP registration otherwise.

## MCP server (`rsk mcp`)

- Flow: `project_profile` (read-only) → propose to user → `install_skills` only after approval.
- Client forces approval on every `install_skills` call; this does not replace asking the user first.

Invoke the full rsk-guide skill for: setting up a project manifest or bundle, choosing pin vs install, global or official-cache installs, new-machine setup, the bundle list, `.rsk/` layout, or MCP server details.

Full guide: invoke the rsk-guide skill.

---
status: draft
project: ralvaskills
date: 2026-11-12
related:
  - plans/2026-10-02-mcp-server-plan.md
---

# MCP registration for Claude Code, Codex, and OpenCode

## Scope

Register the existing `rsk mcp` stdio server in each selected client's project
configuration, preserve unrelated configuration, and document and verify each
client's approval behavior for `install_skills`. Keep the existing server and
tool contracts; do not add a new skill or change the registry.

OpenCode v1 support remains a compatibility requirement. OpenCode v2's native
MCP layout differs, but its migration guide says supported v1 config is still
accepted; the adapter must not silently convert a v1-only user's config to a
v2-only format. Codex project configuration applies only to trusted projects.

## Current findings

1. `internal/tool/mcp.go` writes Claude's `.mcp.json` and OpenCode's v1
   `opencode.json` layout. The Codex adapter returns an unsupported error and
   `rsk new --for all --mcp` is rejected before configuration.
2. The generic JSON writer replaces an existing non-object server container
   and overwrites an existing `rsk` entry. Removal deletes that entry by name,
   without checking whether rsk owns its current contents. Atomic writes can
   replace a symlink; MCP configuration does not use `refuseIfLink`.
3. OpenCode's reader handles only root `opencode.json` as strict JSON. It does
   not inspect `.jsonc` or `.opencode/` config locations, which can take
   precedence over the file rsk writes. A second `rsk` definition there can
   override or conflict with registration.
4. `registerProjectMCP` writes clients sequentially after `rsk.mod` and pinned
   configuration. A later failure can leave earlier clients registered. Destroy
   also removes configuration sequentially; neither path records ownership.
5. `install_skills` mutates links and manifests. Its `destructiveHint` is
   advisory and `anthropic/requiresUserInteraction` is Claude-specific, yet the
   tool description says the client asks for approval. Codex and OpenCode
   approval behavior has not been verified. The tool input's `for` description
   also omits Codex even though the parser accepts registered IDs.
6. Tests cover Claude idempotence and one OpenCode v1 registration, but not
   Codex, v2, collisions, invalid container types, JSONC, symlinks, partial
   failure, or actual client approval and stdio startup.

## Phases

1. **Define and test configuration ownership.** Inventory supported project
   locations and precedence for each client. Preview target paths and detect
   malformed files, non-object containers, symlinks, duplicate server names,
   and conflicting user-owned `rsk` entries before writing. Matching existing
   entries are no-ops. Remove or migrate only entries proven to be managed by
   rsk, including a conservative rule for legacy entries. Never claim ownership
   of an arbitrary same-name server. Preserve comments, unrelated settings,
   and file permissions; use format-aware edits or fail with an actionable
   message when safe preservation cannot be guaranteed. Preflight every target
   before mutation, and design rollback or a clear recoverable partial-failure
   report for failures after preflight.
2. **Add the Codex adapter.** Write only the managed `mcp_servers.rsk` entry in
   project `.codex/config.toml`, using `command = "rsk"` and initially
   `args = ["mcp", "--read-only"]` until install approval is verified.
   Preserve TOML comments and other tables. Support idempotent registration and
   ownership-aware removal; do not edit the user's global Codex configuration
   or `AGENTS.md`. Explain the trusted-project and PATH requirements. Remove the
   `--mcp` guard only when the adapter and its failure paths pass tests.
3. **Support OpenCode configuration deliberately.** Test v1 compatibility in
   the installed v2 client first. Provide native v2 `mcp.servers.rsk` support
   while retaining a documented v1 path; choose v1 or v2 from an explicit or
   unambiguous configuration/version signal, and fail clearly if ambiguous.
   Do not silently migrate other servers. Detect legacy `mcp.rsk` and native
   `mcp.servers.rsk`, migrate only an identifiable managed entry, and avoid
   duplicate definitions. Respect project config precedence and JSONC comments,
   or refuse unsupported configurations rather than write an ineffective or
   destructive parallel file. Removal must handle managed v1 and v2 entries.
4. **Verify approvals before enabling writes.** Keep Claude's interaction
   metadata and destructive annotation, but treat annotations as hints for other
   clients. Configure and test Codex per-tool `approval_mode = "prompt"` and
   OpenCode's `rsk_install_skills` ask permission where supported, including
   nested Code Mode calls, noninteractive mode, policy precedence, and bypass or
   auto-approval modes. Ensure the user sees names, scope, and targets before
   approval. If a client or mode cannot enforce the required per-install
   confirmation, expose only read-only tools in that mode until an enforceable
   approval mechanism is available; do not claim the metadata guarantees it.
   Correct the tool description and schema, including the Codex target.
5. **Complete lifecycle and documentation.** Wire `rsk new`, reruns, PATH
   warnings, and `rsk destroy` to the safe adapters. Preserve unrelated entries
   and avoid registering some clients before discovering another is invalid.
   Update CLI help, README, and changelog with supported versions, migration,
   trust, permissions, and recovery. The bundled `rsk-guide` may require a
   separately authorized edit and version bump because it publishes on push.

## Exit criteria

- [ ] All three clients discover and launch `rsk mcp` from a scratch project;
      protocol output stays on stdout and diagnostics on stderr.
- [ ] Registration and removal are idempotent, preserve unrelated configuration
      and comments, and refuse conflicts, invalid shapes, and symlinked files.
- [ ] Tests cover existing and missing files, v1/v2 OpenCode, mixed/duplicate
      entries, TOML and JSONC, config precedence, collision and failure paths,
      repeated setup, destroy, and interrupted writes.
- [ ] Real client tests record exact versions and demonstrate read-only calls
      and the install approval or safe read-only fallback, including relevant
      auto-approval and noninteractive modes. Automated tests verify generated
      policies and MCP annotations but do not substitute for client behavior.
- [ ] `go test ./...`, `task lint`, and `task registry:check` pass. Affected
      documentation and CLI help describe limitations without promising a
      client-side guarantee that has not been verified.

## References

- [Codex MCP configuration](https://learn.chatgpt.com/docs/extend/mcp?surface=cli)
- [Codex configuration reference](https://learn.chatgpt.com/docs/config-file/config-reference)
- [OpenCode v1 migration](https://opencode.ai/docs/migrate-v1)
- [OpenCode MCP servers](https://opencode.ai/docs/mcp-servers)
- [OpenCode permissions](https://opencode.ai/docs/permissions)
- [Claude Code MCP](https://code.claude.com/docs/en/mcp)

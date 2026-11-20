---
title: 'Publish docker-architect 1.5.0 with Docker Sandboxes'
priority: p2
section: ralvaskills
topics:
    - registry
    - skills
---
`docker-architect` 1.5.0 (`skills/infra/docker-architect/`) is committed: a new section 12 on Docker Sandboxes (`sbx`) for isolated coding agents, with the direct versus `--clone` workspace choice, credentials through `sbx secret set`, deny-by-default network policy, no host `docker.sock`, and an `sbx` command recipe. `CHECKLIST.md`, `DIGEST.md` and `STACK.md` carry the matching rules. It is published only when the `Publish Registry` run for the push finishes green and commits the regenerated index.

- The facts come from the Docker Sandboxes docs. The install method and the supported-agents list were left out because they were not confirmed, so `RECIPES.md` points to the docs for the agent names.
- `task tokens` has not been run for the new size, so `docs/TOKEN_COUNTS.md` still shows the 1.4.1 counts.
- `CHANGELOG.md` has no entry for the release yet.

Done when the index lists `docker-architect` at `v1.5.0`, the `docker-architect@v1.5.0` release exists, the token counts and changelog are updated, and the unconfirmed install and agents details are checked against the docs.

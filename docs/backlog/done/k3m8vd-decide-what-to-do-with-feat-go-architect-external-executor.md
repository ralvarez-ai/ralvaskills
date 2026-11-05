---
title: 'Decide what to do with feat/go-architect-external-executor and its worktree'
priority: p2
section: ralvaskills
topics:
    - git
    - owner-decision
---
The branch has 10 commits that are not in `main` and lives in the worktree `../ralvaskills-executor-port`, which has one uncommitted change. It exists on `origin` too. It was left alone when the merged landing branches were cleaned up, because deleting it would lose that work. It edits a skill, so merging it publishes on push (bump the version and check the diff first).

Done when the branch is merged (or deliberately dropped), the worktree is removed with `git worktree remove` and the remote branch is deleted.

Resolved 2026-11-04: rebased over main, merged and cleaned up. Only the external-program recipe survived (go-architect 1.15.0). The branch's repo-config commits (skeleton AGENTS.md, rsk pins, MCP registration) described the worktree itself and were dropped; its TASKFILE.md work was already in main.

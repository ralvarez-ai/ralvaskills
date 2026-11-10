---
title: '`go-architect` 1.20.0: ports at the context root, no `adapter`/`adapters` folder, and `path` in golangci exclusions is a regex'
priority: p2
section: ralvaskills
topics:
    - go-architect
    - golangci
---

Three rules keep being re-derived in review because the skill does not state them:

- **Ports live in the context root package.** `ledger/sale` declares `Sales` and `Repository`; a `port/` folder only
  survives when a cycle or a cross-context contract forces it, and then it is named for the contract (`coreapi`), never
  for the role. `plc/openplc/port` should dissolve into `package openplc`; `ralvafi/internal/financemcp/port` cannot
  (cycle) and becomes `coreapi`.
- **`adapter` and `adapters` are never folder names.** Every entry inside is already named for what it adapts; the
  bucket hides the dependency direction the archtest ranks on.
- **In golangci-lint 2.x, `exclusions.rules[].path` is a regex, not a glob.** `identity/{a,b}/` matches nothing, so a
  brace list silently disables the rule; `identity/(a|b)/` works. The template must use the second form and say why.

Publish as 1.20.0; the estate re-pin is tracked in ralvafi.

**2026-11-10: done, pending the push.** `go-architect` is 1.20.0: SKILL.md gained the ports-at-the-root rule (a `port/`
folder only survives a cycle or a cross-context contract, and is named for the contract — `coreapi`), the ban on
`adapter`/`adapters` as folder names, and the golangci note that `exclusions.rules[].path` is a **regex** (`{a,b}`
matches nothing; a stale path after a rename silently disables the rule). DIGEST.md carries both as always-on bullets so
`rsk pin` materialises them.

`task registry:generate` + `registry:sync-index` rebuilt the index; only the `go-architect` entry was spliced in, because
the regenerated index also carried three skills from another session's uncommitted edits (remotion-architect,
explainer-architect, narration-architect). Publication happens on push, so the estate re-pin (ralvafi `fcv886`) waits for
that.


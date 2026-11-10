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

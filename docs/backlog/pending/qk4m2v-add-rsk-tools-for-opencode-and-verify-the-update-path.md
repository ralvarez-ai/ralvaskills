---
title: 'Add `rsk tools --for opencode` and verify the update path'
priority: p2
section: ralvaskills
topics:
    - rsk
    - permissions
related:
    - ralvaskills#v7pd1a
---
`rsk tools list|allow|deny|remove` still answers "other tools are not supported yet" (`cmd/rsk/claude.go`), while the machinery for OpenCode already exists: `openCodeTool` implements `GrantPermissions`/`RevokePermissions` over the `permissions` array (`internal/tool/opencode.go`) and `rsk install`/`rsk update` use it. What is missing is the hand-operated side, so a user who does not want to install a skill still cannot add a rule with rsk.

Why it is not a copy of the Claude command: the shapes differ. Claude stores two string lists (`allow`, `deny`); OpenCode stores `{action, resource, effect}` objects where the **last matching rule wins**. A `deny` inserted before an allow is dead weight, and a `list` that drops `effect` hides why a command still runs. The `--output json` contract has to say which shape it is emitting.

Second half: `rsk update`'s reconcile is covered by unit tests (`internal/permission`) and shares `reconcileLock` with install, which was exercised end to end, but no real `rsk update` has run against a published version — the local path pulls a repo, so it needs either a network pull or a tag pushed to a scratch remote.

Done when `rsk tools list|allow|deny|remove --for opencode` read and write `opencode.json`, `deny` lands after the allow so it wins, `--output json` reports `{action, resource, effect}` per rule, the Claude path is unchanged, and an `rsk update` on a skill bumped to a version with a new `allowed-tools` entry writes the new rule and reports the withdrawn one.

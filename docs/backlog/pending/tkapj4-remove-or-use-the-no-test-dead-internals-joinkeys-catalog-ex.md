---
title: 'Remove or use the no-test dead internals'
priority: p3
section: ralvaskills
topics:
    - ralvaskills
---

`joinKeys`, `catalog.ExcludePersonal`, `TargetScope.String`, `skill.ParseVersion`, `Version.String/Compare/Before/After`, `ui.Dim`. Check `Stringer`/formatting first; then add `deadcode:check` to CI (the `deadcode` task exists, but no gate)

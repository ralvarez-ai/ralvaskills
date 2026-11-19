---
title: Stop rsk from leaving .opencode skill symlinks untracked in consumers
priority: p3
section: ralvaskills
topics:
    - rsk
    - opencode
---
`.claude/skills/` is ignored in every consumer, but seven of the repos touched
in this pass had no `.opencode/skills/` entry, so `rsk install` left a dozen
skill symlinks into `~/.config/rsk/cache` showing as untracked forever. They
were ignored by hand in that pass.

Pick one:

- have `rsk install`/`new` ensure the ignore entry next to the symlinks it
  writes, as part of the same generated block idea that owns the tool configs;
- or document the entry in the adapter docs and drop it from the consumer
  setup step.

Symlinks into a machine-local cache should never be committable; a consumer
that commits them breaks on the next machine.

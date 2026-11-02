---
title: '`internal/{config,fsx,cmdx}` (about 400 LOC) against the kits'
priority: p3
section: ralvaskills
topics:
    - ralvaskills
---

`config` overlaps `viperkit/configparts` and `fsx`/`cmdx` are filesystem and command helpers with no kit today. Assess after the R153 sweep; extract only if the executor port needs the same code

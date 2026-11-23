---
title: Fix the Publish Registry Pages deploy when the artifact is missing
priority: p1
section: ralvaskills
topics:
    - ci
    - registry
---
Both `Publish Registry` runs on 2026-11-23 failed at `Deploy to Pages` with `No artifacts named "github-pages" were found for this workflow run`, although `Upload Pages artifact` (`actions/upload-pages-artifact@v5`) reported success. Run `45275217475` still created `go-architect@v1.21.0` and committed the index (`64c007d`), so the skill is published; only the Pages refresh is broken. The second run (`45276049512`) reproduced it from `main`.

Check the `registry` upload path, the Pages environment and the `actions/upload-pages-artifact`/`deploy-pages` majors against GitHub's current guidance, and re-run from `main`. Done when `https://skills.ralva.dev/index.json` serves the new index after a workflow run and the deploy step is green.

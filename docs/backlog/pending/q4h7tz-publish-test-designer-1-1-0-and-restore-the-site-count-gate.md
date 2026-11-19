---
title: 'Publish test-designer 1.1.0 and restore the site count gate'
priority: p2
section: ralvaskills
topics:
    - registry
    - skills
    - site
---
`test-designer` 1.1.0 (`skills/quality/test-designer/`) is committed: logic tests kept apart from e2e tests, no sleep-based waits (await a signal, fake clock, or deadline-bounded polling), injectable clocks, readiness-gated e2e, redundancy pruning and flake triage. `task registry:check` reports it as a new skill at `v1.1.0`; it is published only when the `Publish Registry` run for the push finishes green and commits the regenerated index.

- The skill has no bundle yet. Decide whether it joins `code-review` (reviewers only, so probably not) or a testing grouping next to `tdd`.
- `site/scripts/gen-data.mjs` asserts 67 skills, while the tree holds 69 (`notify-on-finish` is also unpublished), so `task site:check` fails until the counts in the script, `site/src/data/facts.ts` and `AGENTS.md` are brought to the published number.
- Its first release carries 1.1.0 directly, so no `test-designer@v1.0.0` tag will exist.

Done when the index lists the skill, `test-designer@v1.1.0` is released, the site count gate passes on the updated `main`, and the bundle decision is recorded.

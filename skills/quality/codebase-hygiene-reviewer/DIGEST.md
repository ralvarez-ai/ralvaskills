# Codebase Hygiene Reviewer (digest)

Evidence-based declutter audit: dead code, unused exports, duplication, oversized files, split candidates. Proposes only; never edits.

## Always-on rules
- Audit tracked files only (`git ls-files`); exclude generated code and note the exclusion in the report.
- For a library, enumerate real consumers before calling anything dead. "This repo does not use it" is not "nobody uses it".
- Tests count as roots by default.
- Keep categories apart: dead code (delete), unused export (unexport or deprecate), duplication (extract or share), size (candidate only), modularization (hand off).
- Verify every hit before it reaches the report. Static tools miss reflection, cobra `RunE`, handlers in maps, framework conventions.
- Migration `Down`/`Status`/`Version` and interface-only implementations are not findings.
- Size is a signal. A split needs a named conceptual seam, not a file count. No seam means "monitor".
- Re-run the tools after each deletion batch and report the delta.
- Write the report to the OS temp dir as `hygiene-report-<timestamp>.md`, not into the repo unless asked.

## Hard constraints
- Propose only. Deleting is a separate, user-approved step, one batch at a time, tests green between batches.
- Never recommend hand-editing generated code. Fix the generator input or exclude it.
- Library API removals are semver-breaking: route them through deprecation, not deletion.

Invoke the full skill when: cleaning up a repo or library, hunting dead code, sizing a split, or gathering evidence before a shared-library release.

Full guide: invoke the codebase-hygiene-reviewer skill.

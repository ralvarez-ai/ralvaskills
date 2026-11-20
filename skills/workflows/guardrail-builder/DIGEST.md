# Guardrail Builder (digest)

Turn a repeated agent mistake into enforcement, highest level first, and keep a rule-to-enforcer table.

## Always-on rules
- A mistake class counts once it has happened twice; find them in commits, reverts, review comments and workaround comments.
- Fix at the highest level that works: 1 architecture (one owner, one way, hide internals, delete copyable old code), 2 types, then a lint or CI check whose error names the fix, 3 a behavior test, 4 docs last and only for judgment calls.
- Prove each new check fails on a real past mistake; run the same command locally and in CI.
- Exceptions sit on the offending line with a reason, an expiry date and a human's approval.
- Keep a rule-to-enforcer table in `AGENTS.md`. A rule with no enforcer that was broken again is a repeat: fix it at a higher level in the same change.
- Configuration beats instructions: setting, then hook, then a sentence.

## Hard constraints
- Do not add a docs-only rule when a lint or type can enforce it.
- Do not keep a rule whose mistake can no longer happen.

Invoke the full skill when: the operator corrects the same mistake twice, or asks to harden a repo against agent mistakes.

Full guide: invoke the guardrail-builder skill.

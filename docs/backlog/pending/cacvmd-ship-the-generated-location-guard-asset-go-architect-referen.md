---
title: Ship the generated-location guard asset go-architect references
priority: p2
section: ralvaskills
topics:
    - codegen
    - tooling
---
`go-architect` 1.21.0 tells readers to add a script that fails when `*.gen.go`, `*.pb.go`, `*.gen.ts` or `schema.d.ts` sits outside `**/gen/**`, with allowlists for sqlc (`**/postgres/db/`), router codegen and tool scaffolding. The skill has no such asset yet, so the recommendation is prose only.

Add `assets/check-generated-location.sh`, document its usage next to the other style guards, and bump the skill with the asset. The script should accept a root, use the repo's `gen/` convention, and print the offending paths with the allowlist rule that let them through.

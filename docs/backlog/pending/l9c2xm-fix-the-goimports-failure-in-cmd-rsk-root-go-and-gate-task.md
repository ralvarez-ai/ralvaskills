---
title: 'Fix the goimports failure in cmd/rsk/root.go and gate task lint'
priority: p2
section: ralvaskills
topics:
    - code
    - tooling
---
`task lint` fails on `main` with `cmd/rsk/root.go:13:1: File is not properly formatted (goimports)`, and `gofmt -l cmd/rsk/root.go` agrees, so it is not a local toolchain quirk. No workflow runs `task lint` — `check-site.yml`, `publish-registry.yml`, `release.yml` and `deploy-site.yml` do not — which is how the drift survived.

Found while committing `backend-performance-architect` (`e3e812d`); that commit touches no Go files, so the failure predates it.

Fix the formatting, then decide whether `task lint`, `task registry:check` and `task tokens:lint` belong in `check-site.yml` or a separate workflow, so a red gate stops depending on someone running it by hand.

Done when `task lint` is green on `main` and CI enforces it.

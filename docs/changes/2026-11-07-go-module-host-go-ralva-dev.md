---
status: done
project: ralvaskills
date: 2026-11-07
related: [backlog/done/x3jr7m-move-the-go-libraries-and-eight-apps-to-go-ralva-dev.md]
---

# Go modules move to the `go.ralva.dev` host

## Why
Module paths depended on `go.ralvarez.dev` and on the GitHub organization of
each repository. One vanity host for libraries and apps keeps import paths
stable when repositories change organization.

## Impact and migration
- Libraries: every module declares `go.ralva.dev/<name>` from its next tag:
  corekit v0.6.0, email v0.7.0, ratelimit v0.8.0, restkit v0.19.0,
  ginkit v0.19.0, connectrpckit v0.4.0, resilience v0.4.0, grpckit v0.10.0,
  pgkit v0.13.0, tick v0.3.0, outbox v0.5.0, repokit v0.3.0, identity v3.4.0,
  contracts v0.34.0, identitygin v0.10.0, identityclient v0.18.0,
  identityotel v0.2.0, svckit v0.5.0, mcpkit v0.12.0, mcpkit/mcpconfig v0.2.0,
  natskit v0.5.0, webpush v0.8.0, notifierclient v0.5.0, sqlitekit v0.3.0,
  stylekit v0.3.0, termkit v0.67.0, viperkit v0.4.0, wailskit v0.2.0.
  Earlier tags keep `go.ralvarez.dev`, which the Worker serves for ever.
- Apps: autobacklog, autoreports, ralvadocket, ralvafi, ralvagym, ralvapass,
  ralvasend and ralvaskills are `go.ralva.dev/<name>` in the working tree; no
  app tag declares it yet.
- `go-vanity` serves every module on every host (ADR 0002); the apps are in
  its `MODULES` with their organization as owner.
- Consumers: bump each `require` to the tag above, rewrite the imports, run
  `go mod tidy` with `GOPRIVATE=go.ralva.dev/*`.

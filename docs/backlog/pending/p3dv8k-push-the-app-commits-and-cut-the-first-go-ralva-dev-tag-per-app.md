---
title: 'Push the app commits and cut the first `go.ralva.dev` tag per app'
priority: p2
section: ralvaskills
topics:
    - go-module-host
    - ralva-dev
---
`ralvaskills` and `autoreports` are pushed (2026-11-07). Six apps still have local commits (`go.mod`, imports, module path, changelog) that are not pushed: autobacklog, ralvadocket, ralvafi, ralvagym, ralvapass, ralvasend. None has a tag that declares the new path, so `go install go.ralva.dev/ralvaskills/v3/cmd/rsk@latest` fails and released `v3.x` keep the old path. Push, then tag each app (minor bump, as the libraries) and move its `CHANGELOG` entry from `Unreleased` to the version. The vanity Worker already serves the eight apps. Pushing and tagging are outward-facing: confirm per app.

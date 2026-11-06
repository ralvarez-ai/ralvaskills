# Repo Tooling Recipes

Config samples moved out of [SKILL.md](SKILL.md); section numbers match it.

## §4 Sample `Taskfile.yml`

```yaml
# Taskfile.yml
version: '3'

dotenv: ['.env', '.env.local']

tasks:
  default:
    desc: List all tasks
    cmds:
      - task --list

  proto:generate:
    desc: Regenerate Go code from .proto files
    sources:
      - 'proto/**/*.proto'
    generates:
      - 'internal/**/*.pb.go'
    cmds:
      - buf generate

  test:unit:
    desc: Run unit tests
    cmds:
      - go test ./...

  test:race:
    desc: Run all tests with race detector
    cmds:
      - go test -race ./...

  docker:up:
    desc: Start dev stack (DB + app)
    cmds:
      - docker compose up -d

  docker:down:
    desc: Stop dev stack
    cmds:
      - docker compose down
```

## §4 `includes:` example

```yaml
# Taskfile.yml
includes:
  docker: ./Taskfile.docker.yml
  proto: ./Taskfile.proto.yml
```

```yaml
# Taskfile.docker.yml
version: '3'
tasks:
  up:
    desc: Start dev stack (DB + app)
    cmds:
      - docker compose up -d
  down:
    desc: Stop dev stack
    cmds:
      - docker compose down
```

## §5 `.pre-commit-config.yaml` — minimal set

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
        exclude: '\.md$'             # markdown uses trailing space for hard breaks
      - id: end-of-file-fixer
      - id: check-added-large-files
        args: ['--maxkb=500']
      - id: check-merge-conflict
      - id: check-yaml
      - id: check-json

  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.21.0
    hooks:
      - id: gitleaks                 # secret detection
```

## §7 `renovate.json`

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": ["config:recommended"],
  "schedule": ["before 6am on Monday"],
  "packageRules": [
    {
      "matchUpdateTypes": ["minor", "patch"],
      "groupName": "minor + patch updates"
    },
    {
      "matchPackageNames": ["go", "python", "node"],
      "matchUpdateTypes": ["major"],
      "automerge": false,
      "labels": ["language-bump"]
    }
  ],
  "vulnerabilityAlerts": {
    "enabled": true,
    "labels": ["security"]
  }
}
```

## §8 act: config, task wiring and gotchas

```
# .actrc — tracked, no secrets
-P ubuntu-latest=catthehacker/ubuntu:act-latest
--container-architecture=linux/amd64
```

```
# .secrets.example — tracked template; .secrets itself is gitignored
MY_TOKEN=
```

Wire it into the task runner so it's one command, not a memorized `act`
invocation:

```yaml
  ci:local:
    desc: Run CI locally via act
    cmds:
      - test -f .secrets || (echo "Missing .secrets — cp .secrets.example .secrets" && exit 1)
      - act push --secret-file .secrets --artifact-server-path .artifacts
```

- **Switch the workflow's trigger to `workflow_dispatch`-only** if the goal
  is "never risk a billed run" rather than "have a free local option in
  addition to GitHub's" — pushing commits/tags then can't accidentally
  spin up a hosted runner; `workflow_dispatch` stays available from the
  Actions tab as a manual, opt-in escape hatch.
- **`--artifact-server-path`** is required for `actions/upload-artifact` /
  `actions/download-artifact` to work between jobs locally — without it
  those calls silently no-op instead of erroring.
- **Secrets and repo variables** come from `--secret-file` / `--var-file`,
  not GitHub's secret store — gitignore the real file, commit a
  `.example` template (same shape as the `.env`/`.env.local` pattern in
  §6).
- **A tag-triggered job needs a synthesized event file**, not just
  `act push` — `{"ref": "refs/tags/vX.Y.Z"}` via `-e event.json` sets
  `github.ref`/`github.ref_name` the same way a real tag push would, which
  matters for any `if: startsWith(github.ref, 'refs/tags/v')` gate or
  version string baked into an image tag.
- **A job that builds and pushes a Docker image to a registry running on
  the same machine hits two more gotchas, in sequence:**
  1. If that registry sits behind a tunnel or reverse proxy whose
     domain is only meaningful for *remote* access, pushing through it
     from the same host that runs the registry can hit the tunnel's own
     edge timeout on a multi-MB layer upload (often ~100s) —
     traffic that never needed to leave the machine. Make the registry
     host a repo variable (`REGISTRY_HOST`, default the real domain) and
     override it to the registry's local address (e.g. `localhost:5000`)
     for local `act` runs specifically.
  2. Bypassing to `localhost` then fails with `connection refused` from
     inside the builder, even though the *job's own* container is
     host-networked under `act` — `docker/setup-buildx-action`'s default
     `docker-container` driver builds in a separate, isolated builder
     container with its own network namespace. Fix: `driver-opts:
     network=host` on that action.
  3. Whatever you use to observe the result, use the task runner's actual
     exit code, not a wrapper's — `task ci:local:publish | tee log | tail`
     reports `tail`'s exit status, which is `0` even when the piped
     command failed. Check the real job output (or the exit code without
     piping through a second command) before calling a run green.

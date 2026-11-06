# docker-architect digest

Docker Engine 29, Compose v2 and BuildKit standards for writing or reviewing Dockerfiles and Compose files.

## Always-on rules

- Always multi-stage; the runtime stage copies only artifacts.
- Layer order least to most volatile. `.dockerignore` is mandatory.
- Final stage runs as non-root, UID >= 10000. `HEALTHCHECK` on long-running services.
- apt: `--no-install-recommends` and `rm -rf /var/lib/apt/lists/*` in the same `RUN`.
- Bases: Go `gcr.io/distroless/static-debian12:nonroot`; Python `python:3.14-slim` (uv in builder); Node `node:22-slim`, LTS only.
- Pin bases by digest in production; Renovate refreshes them.
- BuildKit: cache, secret and bind mounts; here-docs for multi-line scripts. Build secrets come from the environment.
- Build `linux/amd64` and `linux/arm64`.
- Compose: `docker-compose.yaml` plus `docker-compose.override.yaml` for dev; `docker compose` (v2); no `version:` key; `depends_on` with `condition: service_healthy`, which needs a working healthcheck.
- Named volumes for state; explicit networks; resource limits on every service.
- Restart `unless-stopped` for long-running, `no` for batch. `init: true` when the app spawns children. `TZ=Etc/UTC` in the image.
- Logs: json-file with size and count rotation, or journald on Linux hosts.
- Deploys reference image digests. Tags: semver, plus `latest` and `1.2` aliases for local convenience; CI also pushes `:<short-sha>`.
- Trivy scan every image in CI; `.trivyignore` for documented exceptions; SBOM on releases. Sign images with Cosign for prod-bound registries.

## Hard constraints

- Never ship the toolchain in the runtime image.
- Never bake secrets into layers. Build-time: `--mount=type=secret`. Runtime: orchestrator.
- Never run as root in the final stage.
- Never bake dev conveniences into `docker-compose.yaml`.
- Don't use the legacy `docker-compose` binary.
- Avoid alpine for Python (musl breaks scientific wheels).
- Never silent allowlists in Trivy.
- Never rely on the host timezone.

Invoke the full skill when: writing or reviewing a Dockerfile or Compose file, adding Trivy or multi-arch steps, or scaffolding a new project (checklist).

Full guide: invoke the docker-architect skill.

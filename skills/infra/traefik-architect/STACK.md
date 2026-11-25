# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| traefik | v3.8 | Reverse proxy; Docker provider, label routing, ACME |
| docker-compose | 2 (v2 plugin) | Service orchestration - `docker compose` subcommand |
| docker-socket-proxy | current release | Read-only, filtered access to the Docker socket for Traefik |
| ACME DNS challenge | lego provider of your DNS host | Wildcard certificates without inbound access |

## Notes

- **Pin Traefik by digest in production** and let Renovate refresh it; a `v3` tag moves.
- **Traefik v3 rule syntax** (``Host(`x`)``) differs from v2 in matchers such as `HostRegexp`; do not paste v2 rules unchecked.
- **Inherits** the Docker and Compose baseline from [docker-architect](../docker-architect/STACK.md).

_Last reviewed: 2026-11-25_
_Skill version at last review: 1.0.0_

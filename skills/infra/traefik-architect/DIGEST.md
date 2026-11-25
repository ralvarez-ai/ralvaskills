# traefik-architect digest

Traefik v3 and Compose v2 standards for routing web services by name and keeping host ports from colliding.

## Always-on rules

- One Traefik per machine on a shared external `proxy` network; it publishes only 80 (and 443 with TLS).
- `providers.docker.exposedbydefault=false`; only `traefik.enable=true` containers are routed.
- Routed services have no `ports:`. Router, service and middleware names equal the container name.
- Rule ``Host(`<service>.${BASE_DOMAIN}`)``; set `loadbalancer.server.port` explicitly. `BASE_DOMAIN` lives in the untracked `.env`, never in a tracked compose file.
- Names: `*.localhost` locally (SSH `-L` to reach a remote one), a per-machine wildcard `*.<machine>.<domain>` for stable names. Avoid `.local`, `.dev`, `.app`. A wildcard covers one label; never wildcard the apex.
- Wildcard HTTPS through the DNS challenge; the API token is passed by file, never inline.
- Bind the proxy to a private interface; use a socket proxy, not the raw Docker socket, on shared hosts.
- Proxy HTTP(S) only. TCP/UDP services keep fixed ports.
- Fixed ports: one flat list per machine, unique per `(machine, port, protocol)`, outside the ephemeral range (`net.ipv4.ip_local_port_range`). Claimed but not listening is dormant, not free.
- Labs publish nothing, use `-p lab-<name>`, and let Docker or the OS pick ports.

## Hard constraints

- Never publish `ports:` on a service that is also routed by label.
- Never run two proxies on 80/443.
- Never mount the Docker socket into an internet-facing proxy.
- Never add an internal hostname to a public tunnel.
- Never treat "not listening" as "free".

Invoke the full skill when: adding a proxy, routing a service by name, choosing domain names or DNS records, setting up HTTPS, or fixing a port collision.

Full guide: invoke the traefik-architect skill.

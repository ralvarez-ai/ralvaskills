# Traefik Recipes

Reference skeletons for [SKILL.md](SKILL.md). Replace `<base>` and `<ip>` with the machine's values; nothing here is specific to one host.

## 1. Proxy compose (one per machine)

```yaml
name: proxy
services:
  traefik:
    image: traefik:v3.8
    container_name: traefik
    restart: unless-stopped
    command:
      - --providers.docker.endpoint=tcp://socket-proxy:2375
      - --providers.docker.exposedbydefault=false
      - --providers.docker.network=proxy
      - --entrypoints.web.address=:80
      - --api.dashboard=false
      - --metrics.prometheus=true
    ports: ["127.0.0.1:80:80"]
    networks: [proxy, socket]
    depends_on: [socket-proxy]
  socket-proxy:
    image: tecnativa/docker-socket-proxy
    restart: unless-stopped
    environment: { CONTAINERS: 1 }
    volumes: ["/var/run/docker.sock:/var/run/docker.sock:ro"]
    networks: [socket]
networks:
  proxy: { name: proxy }
  socket: { internal: true }
```

- Bind `ports` to the machine's private IP instead of `127.0.0.1` to serve other hosts on the LAN or VPN.
- The `socket` network is `internal`, so the socket proxy has no route out.
- Pin the image by digest in production and refresh it with Renovate.

## 2. Routed service (no `ports:`)

```yaml
services:
  grafana:
    image: grafana/grafana
    container_name: grafana
    networks: [proxy, default]
    labels:
      - traefik.enable=true
      - traefik.http.routers.grafana.rule=Host(`grafana.${BASE_DOMAIN}`)
      - traefik.http.routers.grafana.entrypoints=web
      - traefik.http.services.grafana.loadbalancer.server.port=3000
networks:
  proxy: { external: true }
```

`.env` per machine: `BASE_DOMAIN=localhost` on a laptop, `BASE_DOMAIN=<machine>.<domain>` on a server.

## 3. HTTPS with a wildcard certificate (DNS challenge)

Add to the Traefik `command`:

```yaml
- --entrypoints.websecure.address=:443
- --entrypoints.web.http.redirections.entrypoint.to=websecure
- --entrypoints.web.http.redirections.entrypoint.scheme=https
- --entrypoints.websecure.http.tls.certresolver=le
- --entrypoints.websecure.http.tls.domains[0].main=<machine>.<domain>
- --entrypoints.websecure.http.tls.domains[0].sans=*.<machine>.<domain>
- --certificatesresolvers.le.acme.dnschallenge.provider=cloudflare
- --certificatesresolvers.le.acme.storage=/acme/acme.json
```

Pass the provider token by file and persist the store:

```yaml
environment:
  CF_DNS_API_TOKEN_FILE: /run/secrets/dns_token
secrets: [dns_token]
volumes: [acme:/acme]
```

Routed services then use `entrypoints=websecure` and need no per-service certificate label. Other DNS providers use their own variable names; the `_FILE` suffix works for any of them.

## 4. Lab and test patterns

```sh
docker compose -p lab-foo up -d        # isolated project name
docker compose -p lab-foo down -v      # removes only this lab
docker compose -p lab-foo port db 5432 # the host port Docker chose
```

For a host process that must reach a service, publish to loopback with a Docker-chosen port: `ports: ["127.0.0.1::5432"]`.

## 5. Reaching a remote machine with no DNS

```sh
ssh -L 8080:localhost:80 <host>
```

Then open `http://grafana.localhost:8080`. The laptop resolves `*.localhost` to itself, the tunnel carries the request, and the remote Traefik routes by name. Use a local port above 1024.

## 6. Auditing fixed ports

```sh
ss -ltn                                       # live listeners
docker ps -a --format '{{.Names}} {{.Ports}}' # published by existing containers
docker compose config --format json --profile '*' # claims in a compose file
sysctl net.ipv4.ip_local_port_range           # the ephemeral range to stay out of
```

A stack taken down with `docker compose down` leaves no container, so only its compose file still records its ports. Keep the allocation list in version control and treat it as the source of claims; use the live listing only to find listeners nobody claimed.

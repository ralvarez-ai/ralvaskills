---
name: traefik-architect
version: 1.1.0
description: Traefik v3 standards - host-based routing from Docker labels, one proxy per machine on a shared network, BASE_DOMAIN per machine, *.localhost / sslip.io / wildcard DNS naming, DNS-challenge HTTPS, private-interface binding, socket proxy, and the fixed-port convention for the TCP/UDP services a proxy cannot cover. Use when publishing web services behind a proxy, fixing host port collisions, or setting up local, lab and homelab stacks.
---

# Traefik Architecture & Port Discipline

One Traefik per machine replaces per-service published ports for everything that speaks HTTP(S). Services keep no `ports:` entry and are reached by name. Targets **Traefik v3** and **Compose v2**. Companion to [docker-architect](../docker-architect/SKILL.md). Full compose and label examples in [RECIPES.md](RECIPES.md); pinned versions in [STACK.md](STACK.md).

## 1. When a proxy fits, and when it does not

- **HTTP(S) services go behind the proxy.** Web UIs, APIs, dashboards. The proxy routes by the `Host` header, so any number of services share ports 80 and 443.
- **TCP and UDP services keep a published port.** Databases, brokers, game servers carry no hostname. A proxy can forward raw TCP on an extra entrypoint, but only one service per port, so do not rely on it.
- **Container to container never goes through the proxy.** Use the service name on the shared network (`http://grafana:3000`). The proxy is for clients outside the Compose network.
- **Count before building.** If most published ports are dev or debug conveniences, deleting them shrinks the problem more than any tool.

## 2. Layout: one proxy, one shared network

- **One Traefik container per machine**, its own compose project, publishing only 80 (and 443 once TLS is on).
- **A named external network** (`proxy`) that Traefik and every routed service join. Services also keep their own default network for private traffic such as a database.
- **`providers.docker.exposedbydefault=false`.** Only containers labelled `traefik.enable=true` become routes. An unlabelled container is never exposed by accident.
- **`providers.docker.network=proxy`**, so Traefik always uses the shared network even when a container is on several.

## 3. Label conventions

- **No `ports:` on a routed service.** The label is the only publication.
- **Router, service and middleware names equal the container name.** They share one namespace across all projects, so a generic name like `web` collides. Service = container = router.
- **One rule per service:** ``Host(`<service>.${BASE_DOMAIN}`)``. Path routes only when two apps must share a hostname, with a `StripPrefix` middleware and a check that the app's links still work.
- **State the backend port explicitly** with `loadbalancer.server.port`. Traefik guesses from `EXPOSE` otherwise.
- **Add a healthcheck** to the service; Traefik skips unhealthy containers only when one exists.

## 4. `BASE_DOMAIN` per machine

- **The only thing that differs per machine is the base.** Keep it in one variable, `BASE_DOMAIN`, in the untracked `.env`, and write every rule as ``Host(`app.${BASE_DOMAIN}`)``. The same compose file then runs on a laptop and a server.
- **Service names stay identical everywhere**: `grafana.localhost` on the laptop, `grafana.<server-base>` on the server.
- **Never hard-code a domain** in a tracked compose file.
- **Environments on one machine share the base and differ in the name**: `<service>.<base>` for prod, `<service>-dev.<base>` for dev, with router names `<container>-dev`. Keep the suffix in the same label (`-dev`, not `.dev.`), so one wildcard record and one wildcard certificate cover every environment.
- **Dev labels live in an opt-in override** (`compose.proxy.yml` with `ports: !reset []`), so the default dev flow keeps working without a proxy.
- **Different machines get different bases**, never different service names.

## 5. Name schemes (no DNS work to start)

| Scheme | Setup | Use |
|---|---|---|
| `*.localhost` | None. Resolves to 127.0.0.1 | Same machine. Default for local work |
| `*.localhost` over an SSH tunnel | `ssh -L 8080:localhost:80 <host>`, then `app.localhost:8080` | Reach a remote machine with no DNS, no exposed port |
| `<service>.<ip>.sslip.io` | None. The name encodes the IP | Private LAN IP, zero accounts. Third-party dependency |
| Own wildcard `*.<machine>.<domain>` | One `A` record per machine | Stable names, real HTTPS |
| `*.internal`, `*.home.arpa` | Local DNS server | Private, reserved suffixes |

- **Avoid `.local`** (mDNS clashes) and **`.dev` / `.app`** (HSTS forces HTTPS in browsers).
- **A wildcard covers one label.** `*.vm.example.com` matches `grafana.vm.example.com`, not `a.b.vm.example.com` or `vm.example.com` itself.
- **Put a machine's wildcard under its own subtree**, never on the apex. An apex wildcard hides typos in production hostnames.
- **Private IP in a public zone is harmless** but needs the DNS provider's proxy off (DNS only). Some routers drop such answers (DNS-rebinding protection); allow-list the suffix.
- **Start with `*.localhost`, add DNS later.** Moving over means changing `BASE_DOMAIN`, adding the record and recreating the routes. Labels do not change.

## 6. HTTPS

- **Use the DNS challenge for wildcard certificates.** It proves ownership through the DNS provider's API, so internal services get valid certificates with no inbound connection.
- **One wildcard cert per machine** declared on the `websecure` entrypoint (`main` = the machine base, `sans` = its wildcard), not one certificate per service.
- **The API token is a secret.** Pass it by file (`<VAR>_FILE`), scoped to DNS edit on one zone. See [secrets-architect](../secrets-architect/SKILL.md). Never in the compose file or `.env`.
- **Persist `acme.json`** on a volume, mode 600.
- **Redirect `web` to `websecure`** once the certificate works. Until then plain HTTP is acceptable for internal-only services.

## 7. Binding and exposure

- **Bind the proxy to the private interface**, for example `10.0.0.5:80:80` or `127.0.0.1:80:80`. Never `0.0.0.0` on a machine with a public interface unless the service is meant to be public.
- **Public exposure is opt-in per hostname** (a tunnel route or a firewall rule). Internal names are never added to a tunnel.
- **Do not mount the Docker socket directly into Traefik on a shared or exposed host.** Put a socket proxy in front that allows read-only container listing. Skeleton in [RECIPES.md](RECIPES.md).
- **The dashboard is off by default** (`api.insecure=false`). If enabled, route it like any service, on a private name.

## 8. Fixed ports for what a proxy cannot cover

General publishing rules (loopback, dynamic ports and ranges, the ephemeral range, no published databases) live in [docker-architect](../docker-architect/SKILL.md) section 13. With a proxy in place they reduce to:

| Service | Host port |
|---|---|
| HTTP(S), including gRPC once TLS is on | none, a hostname |
| Reached only by other containers | none |
| Dependency of a backend run on the host | fixed in a reserved range, or dynamic |
| External non-HTTP client | fixed, listed |

- **Allocate the fixed ones flat, from one list per machine.** One record of `machine, port, protocol, service`. A port is claimed when the record is committed. Routed hostnames are not recorded there.
- **Uniqueness is per `(machine, port, protocol)`.** Ignore the bind address; a `127.0.0.1` versus `0.0.0.0` overlap is a safe false positive.
- **Check claims against reality.** A listing of listening sockets (`ss -ltn`) misses stopped stacks. Compare the record with the compose files (`docker compose config --profile '*'`) and with the live listeners.
- **Claimed but not listening is dormant, not free.**

## 9. Labs and tests

- **Labs publish nothing.** Services talk over the Compose network by name; the runner joins that network.
- **One compose project name per lab** (`-p lab-<name>`), so networks and containers never collide and `down` removes only that lab.
- **Test code binds port 0** or uses Testcontainers, so the OS or the library assigns a free port.
- **Two branches of one project** run side by side by changing only the project name and the hostname.

## 10. Anti-patterns

- Publishing `ports:` on a service that is also routed by label.
- Two proxies fighting over 80 and 443.
- A wildcard on the apex domain.
- A hard-coded domain in a tracked compose file.
- A generic router name (`web`, `app`) shared by two projects.
- The Docker socket mounted into an internet-facing proxy.
- An internal hostname added to a public tunnel.
- Treating "not listening" as "free".

## 11. Out of scope

- **Compose, networks, healthchecks, restart policy** -> [docker-architect](../docker-architect/SKILL.md).
- **Where the API token lives and how it is mounted** -> [secrets-architect](../secrets-architect/SKILL.md).
- **Proxy metrics and access logs** -> [observability-architect](../observability-architect/SKILL.md). Enable the Prometheus metrics entrypoint and scrape it.
- **`.env` and task-runner conventions** -> [repo-tooling-architect](../../tooling/repo-tooling-architect/SKILL.md).

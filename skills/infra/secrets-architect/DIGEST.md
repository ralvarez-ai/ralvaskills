# secrets-architect digest

One external secrets store (`$SECRETS_DIR`, default `~/Secrets`) consumed by reference. Repos, `.env` files and container environments hold paths, never values.

## Always-on rules
- Match the mechanism to the consumer: Compose `secrets:` for container files read through `<VAR>_FILE`; `LoadCredential=` for systemd units; `$SECRETS_DIR` path for shell and go-task; plain writable bind mount for writable state; host repo secret for CI.
- Build a path default as `${SECRETS_DIR:-$HOME/Secrets}`; keep `SECRETS_DIR` defined once in the stack `.env`.
- Build secrets come from `environment: VAR` with the caller exporting `VAR`.
- Layout is `<kind>/<owner>/<name>`; dirs `700`, files `600`, with `644` only for direct-open readers.
- Check the reader's uid (`docker inspect <ctr> --format '{{.Config.User}}'`) before choosing a mode.
- Path variables end in `_FILE`; names are specific; secret file names are lowercase, hyphenated, no extension.
- Verify secret content, not the filename.
- Go modules (`contents:read`) and npm `@scope/*` (`read:packages`) need separate tokens.
- Migration per secret: inventory consumer, place file, pick mechanism, wire, remove value from `.env`/repo, verify, document, commit.
- Verify live: config renders; recreate and observe; prove the credential path, not a fallback; compare `sha256sum` of store vs live without printing values; run `docker compose build`, not just `config`.

## Hard constraints
- No secret value in git, in a `.env`, or in a container environment.
- Never `file:` a build secret outside the build context; source it from the environment.
- Don't blanket `chmod 600`; one non-root reader breaks.
- Compose `secrets:` on non-swarm is a bind mount; `mode:`, `uid:`, `gid:` are ignored, only the host file mode matters.
- systemd `EnvironmentFile` does not expand variables; use `%h/Secrets/…` or `LoadCredential`.
- Moving a bind-mounted file does not affect a running container; recreate it before trusting "it still works".

Invoke the full skill when: wiring a credential into a repo, container, systemd unit or CI job; migrating secrets out of `.env` or git; or auditing the secrets store.

Full guide: invoke the secrets-architect skill.

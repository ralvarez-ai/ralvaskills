# Security Reviewer (digest)

Cross-language code and configuration security review: injection, authz, secrets, insecure defaults, deserialization, CSRF/SSRF/IDOR, dependency vulns. Not a penetration test.

## Always-on rules
- Two passes: tool pass (gitleaks, semgrep, gosec, bandit, trivy, npm audit / pip-audit), then a read pass of the diff. Tools miss intent.
- Every SQL query uses parameter binding. Subprocess calls use list-form args.
- Passwords: Argon2id. Pin the JWT algorithm. Verify `aud` and `iss` for external IdP tokens.
- Tenant or customer id comes from the authenticated context, never from the request body.
- Check authorization per endpoint; do not assume middleware did it.
- Secrets: none in code, committed env files, Dockerfile layers, or logs. Rotate any secret ever committed, even if later removed.
- Parse untrusted input strictly: `yaml.safe_load`, strict JSON decoding, XML external entities disabled.
- Every resource lookup checks ownership. Return 404, not 403, for another caller's resource.
- Validate user-supplied outbound URLs (block private ranges, require HTTPS, allow-list).
- Login endpoints: rate-limit aggressively; errors say "Invalid credentials", never reveal whether a user exists.
- Rate-limit per caller, not per IP.
- Report each finding as Critical, High, Medium, Low or Info, with file:line, evidence and fix. End with a one-sentence summary.

## Hard constraints
- Never put tokens in URLs. Use `Authorization: Bearer` only.
- Never use `Access-Control-Allow-Origin: *` on authenticated endpoints.
- Never `pickle.loads` untrusted data. Never `yaml.load`; use `yaml.safe_load`.
- Never accept the algorithm a token claims; never allow `none`.
- Never use `:latest` tags in Dockerfiles. Never hand-edit lockfiles.
- Flag deeper work (threat model, compliance audit, pentest) explicitly; do not claim it is covered.

Invoke the full skill when: the user asks for a security review or audit, before release of a service touching auth, payments, PII or external integration, after dependency updates with security advisories, or for a new endpoint, auth flow or SQL query.

Full guide: invoke the security-reviewer skill.

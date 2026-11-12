# Wails Architecture (digest)

Wails v3.1 shell around an existing web app: embedded SPA, same-origin proxy, server-held cookies, per-OS chrome, Android shim, packaging.

## Always-on rules
- A shell is a thin client: it embeds the SPA build, proxies API calls and holds no data.
- Never embed the server; a desktop binary cannot run Postgres or the identity service.
- Give the shell its own `go.mod`, lint and archtest, outside the app module.
- One asset handler: API routes to a proxy, everything else to the SPA with an `index.html` fallback.
- A missing file with an extension is a 404, never `index.html`.
- Proxy the API on the webview's origin: no CORS and no base-URL config in the SPA.
- Rewrite the request host; delete `Origin` (or rewrite it for an origin-checking API).
- Answer 502 when the API is unreachable.
- WebKit ignores `Set-Cookie` on custom-scheme responses: hold cookies in a Go jar.
- Attach cookies outgoing, strip `Set-Cookie` from responses, never let the page hold a token.
- In-memory jar means logout on restart; a keychain store is the upgrade.
- Prefer OIDC with a loopback redirect when the API supports it.
- Return `application.WebviewWindowOptions` from one `runtime.GOOS` switch.
- Linux: frameless + injected controls via `JS`; the `CSS` option does not apply under GTK 4.
- Windows: Mica + system theme. macOS: transparent or hidden-inset title bar. Phones: plain options.
- Android has no environment: set the API URL in code; emulator host is `10.0.2.2`.
- Android webview sends every request as a bodiless `GET`; non-200 reads as "asset not found" — wrap `fetch`, bridge to a Go service.
- Fetch shim limits: no XHR, EventSource, streaming bodies, redirects or abort propagation.
- Pad the Android root view by `WindowInsetsCompat`; scaffold regeneration drops the patch.
- Use a full JDK 21; Gradle needs `javac` and `jlink` fails on newer JDKs.
- Rebuild the SPA in every build task; the scaffold task is status-checked and ships a stale `dist`.
- Android release: signed AAB, upload key from secrets, `bundle` not `assemble:aab`, version from the tag, signature via `keytool`.
- The webview has no service worker: no offline cache, no Web Push, no passkeys, no keychain.
- Unit-test the jar, the proxy and the SPA fallback; hand-check the webview per OS.

## Hard constraints
- Never rewrite `Secure`/`Domain`/`SameSite` to fix cookies; the cookie is never stored.
- Never claim a platform works because it compiles.
- Never put the upload key or passwords in the repo or binary; a private-repo GitHub provider needs a token in the binary — prefer an `endpoint` provider with a signed manifest.
- Never design around service workers, Web Push, passkeys or a keychain in the webview.
- Add native features only after the shell works, in value order.

Invoke the full skill when: scaffolding or reviewing a Wails shell, debugging sign-in inside the webview, or wiring Android builds and Play releases.

Full guide: invoke the wails-architect skill.

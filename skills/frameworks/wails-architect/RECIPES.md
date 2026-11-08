# Wails Recipes

Reference implementations for the patterns in [SKILL.md](SKILL.md). All names are placeholders.

## Project layout

```
shell/                              # its own module, outside the app module
├── go.mod
├── main.go                         # wires handler, window, app
├── handler.go                      # SPA + API proxy on one ServeMux
├── proxy.go                        # cookie jar + reverse proxy
├── spa.go                          # index.html fallback
├── window/
│   ├── window.go                   # per-OS application.WebviewWindowOptions
│   ├── controls.js                 # injected minimise/maximise/close cluster
│   └── controls.css
├── api_default.go                  # //go:build !android
├── api_android.go                  # //go:build android
├── frontend/
│   ├── package.json                # build = copy ../web/dist here
│   └── dist/                       # gitignored
├── build/                          # scaffold: Taskfile.yml + per-OS assets
├── scripts/android-version.sh
├── Taskfile.yml
└── bin/                            # gitignored
```

`//go:embed all:frontend/dist` fails if the directory is absent, so create it (with a `.gitkeep`) or build the SPA first. The build task always copies the fresh SPA; do not rely on the scaffold's status-checked frontend task.

## main.go

```go
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/url"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"example.com/shell/window"
)

var defaultAPI = "https://api.example.com" // overridden by -ldflags on Android dev builds

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	apiURL := os.Getenv("APP_API_URL")
	if apiURL == "" {
		apiURL = defaultAPI
	}
	target, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("APP_API_URL: %v", err)
	}
	dist, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	handler, _ := newHandler(dist, target)

	app := application.New(application.Options{
		Name:   "App",
		Assets: application.AssetOptions{Handler: handler},
		Mac:    application.MacOptions{ApplicationShouldTerminateAfterLastWindowClosed: true},
	})
	app.Window.NewWithOptions(window.Options(window.Config{
		Title: "App",
		URL:   "/login", // the shell's first screen, not the public landing page
	}))
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## Handler: one origin for assets and API

```go
func newHandler(dist fs.FS, target *url.URL) (http.Handler, *Jar) {
	jar := NewJar()
	proxy := &Proxy{Target: target, Jar: jar, Client: &http.Client{Timeout: 30 * time.Second}}

	mux := http.NewServeMux()
	mux.Handle("/v1/", proxy)      // everything the SPA calls
	mux.Handle("GET /health", proxy)
	mux.Handle("/", SPA(dist))     // everything else
	return mux, jar
}
```

## Cookie jar (RFC 6265 path match)

```go
type stored struct {
	Name, Value, Path string
	Expires           time.Time
}

type Jar struct {
	mu      sync.Mutex
	cookies []stored
}

func (j *Jar) SetFrom(resp *http.Response) {
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, c := range resp.Cookies() {
		if c.MaxAge < 0 { // deletion
			j.cookies = slices.DeleteFunc(j.cookies, func(s stored) bool {
				return s.Name == c.Name && s.Path == c.Path
			})
			continue
		}
		s := stored{Name: c.Name, Value: c.Value, Path: orSlash(c.Path), Expires: c.Expires}
		j.cookies = slices.DeleteFunc(j.cookies, func(old stored) bool {
			return old.Name == s.Name && old.Path == s.Path
		})
		j.cookies = append(j.cookies, s)
	}
}

func (j *Jar) AddTo(req *http.Request) {
	j.mu.Lock()
	defer j.mu.Unlock()
	now := time.Now()
	var pairs []string
	for _, c := range j.cookies {
		if !c.Expires.IsZero() && now.After(c.Expires) {
			continue
		}
		if !c.pathMatch(req.URL.Path) {
			continue
		}
		pairs = append(pairs, c.Name+"="+c.Value)
	}
	if len(pairs) > 0 {
		req.Header.Set("Cookie", strings.Join(pairs, "; "))
	}
}

func (c stored) pathMatch(p string) bool {
	return c.Path == "" || p == c.Path || strings.HasPrefix(p, strings.TrimSuffix(c.Path, "/")+"/")
}

func orSlash(p string) string {
	if p == "" {
		return "/"
	}
	return p
}
```

The jar is in memory: a restart signs the user out. Replace it with a keychain-backed store, keeping these three methods.

## Reverse proxy

```go
type Proxy struct {
	Target *url.URL
	Jar    *Jar
	Client *http.Client
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out := r.Clone(r.Context())
	out.URL.Scheme, out.URL.Host = p.Target.Scheme, p.Target.Host
	out.URL.Path = singleJoiningSlash(p.Target.Path, r.URL.Path)
	out.Host = p.Target.Host
	out.Header.Del("Cookie")
	out.Header.Del("Origin") // an origin-checking API needs a configured rewrite instead
	p.Jar.AddTo(out)

	resp, err := p.Client.Do(out)
	if err != nil {
		http.Error(w, "API unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	p.Jar.SetFrom(resp) // keep them in Go…
	for k, vs := range resp.Header {
		if strings.EqualFold(k, "Set-Cookie") {
			continue // …and never hand them to the webview
		}
		w.Header()[k] = vs
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
```

## SPA fallback

```go
func SPA(dist fs.FS) http.Handler {
	files := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}
		if _, err := fs.Stat(dist, name); err != nil {
			if path.Ext(name) != "" { // a missing asset is a 404, never index.html
				http.NotFound(w, r)
				return
			}
			http.ServeFileFS(w, r, dist, "index.html") // Go 1.22+
			return
		}
		files.ServeHTTP(w, r)
	})
}
```

## Per-OS window options

```go
func Options(cfg Config) application.WebviewWindowOptions {
	opts := application.WebviewWindowOptions{
		Title: cfg.Title, URL: cfg.URL, Width: 1200, Height: 800,
	}
	switch runtime.GOOS {
	case "darwin":
		opts.Mac.TitleBar = application.MacTitleBar{AppearsTransparent: true}
		if cfg.MacInset {
			opts.Mac.TitleBar = application.MacTitleBarHiddenInset
		}
	case "windows":
		opts.Windows = application.WindowsWindow{
			BackdropType: application.Mica,
			Theme:        application.SystemDefault,
		}
	case "android", "ios":
		// the system draws the bars; nothing to inject
	default: // linux
		if !cfg.NativeChrome {
			opts.Frameless = true
			opts.JS = controlsScript(cfg.DragSelector, cfg.ControlsHeight)
		}
	}
	return opts
}
```

## Frameless Linux controls

The drag mechanism is Wails' `--wails-draggable` custom property. Inject the styles from JS: the window's `CSS` option did not apply under GTK 4.

```js
(async () => {
  const { Window } = await import('/wails/runtime.js');
  const style = document.createElement('style');
  style.textContent = `
    #drag { position: fixed; inset: 0 130px auto 0; height: 10px;
            --wails-draggable: drag; z-index: 2147483646; }
    .app-header { --wails-draggable: drag; }
    .app-header :is(a, button, input, [tabindex]) { --wails-draggable: no-drag; }
    #controls { position: fixed; top: 0; right: 0; display: flex;
                --wails-draggable: no-drag; z-index: 2147483647; }
  `;
  document.head.append(style);

  const strip = document.createElement('div');
  strip.id = 'drag';
  strip.addEventListener('dblclick', () => Window.ToggleMaximise());

  const box = document.createElement('div');
  box.id = 'controls';
  box.innerHTML =
    '<button data-a="min" aria-label="Minimise">&#8211;</button>' +
    '<button data-a="max" aria-label="Maximise">&#9633;</button>' +
    '<button data-a="close" aria-label="Close">&#10005;</button>';
  box.addEventListener('click', (e) => {
    const a = e.target.closest('button')?.dataset.a;
    if (a === 'min') Window.Minimise();
    if (a === 'max') Window.ToggleMaximise();
    if (a === 'close') Window.Close();
  });
  document.body.append(strip, box);
})();
```

## Android fetch shim

Inject this **before** the app bundle (`<head>`), and enable it only on Android.

```js
(() => {
  const routes = ['/v1/', '/health'];
  const native = window.fetch.bind(window);
  window.fetch = async (input, init = {}) => {
    const req = new Request(input, init);
    const url = new URL(req.url, location.href);
    if (!window.wails || !routes.some((r) => url.pathname.startsWith(r))) {
      return native(input, init);
    }
    const bytes = new Uint8Array(await req.arrayBuffer());
    const res = await window.wails.Call('Fetcher.Fetch', {
      method: req.method,
      url: url.pathname + url.search,
      headers: Object.fromEntries(req.headers),
      body: btoa(String.fromCharCode(...bytes)),
    });
    return new Response(Uint8Array.from(atob(res.body), (c) => c.charCodeAt(0)), res);
  };
})();
```

```go
type FetchRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type FetchResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Fetcher struct{ Handler http.Handler }

func (f *Fetcher) Fetch(in FetchRequest) (FetchResponse, error) {
	body, err := base64.StdEncoding.DecodeString(in.Body)
	if err != nil {
		return FetchResponse{}, err
	}
	req := httptest.NewRequest(in.Method, in.URL, bytes.NewReader(body))
	for k, v := range in.Headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	f.Handler.ServeHTTP(rec, req)

	headers := map[string]string{}
	for k := range rec.Header() {
		headers[k] = rec.Header().Get(k)
	}
	return FetchResponse{
		Status:  rec.Code,
		Headers: headers,
		Body:    base64.StdEncoding.EncodeToString(rec.Body.Bytes()),
	}, nil
}
```

Register it as a service and turn the shim on for Android:

```go
app := application.New(application.Options{
	Assets:   application.AssetOptions{Handler: handler},
	Services: []application.Service{application.NewService(&Fetcher{Handler: handler})},
})
```

Limits: `fetch` only; no `XMLHttpRequest`, `EventSource`, streaming request bodies, redirect following or abort propagation.

## Platform API address

```go
//go:build !android

package main

var defaultAPI = "https://api.example.com"
```

```go
//go:build android

package main

// The emulator reaches the host at 10.0.2.2. A release build is compiled with
// -ldflags "-X main.defaultAPI=https://api.example.com".
var defaultAPI = "http://10.0.2.2:8080"
```

## Android system-bar insets

```java
View root = findViewById(android.R.id.content);
ViewCompat.setOnApplyWindowInsetsListener(root, (v, insets) -> {
    Insets bars = insets.getInsets(WindowInsetsCompat.Type.systemBars());
    v.setPadding(bars.left, bars.top, bars.right, bars.bottom);
    return insets;
});
```

`wails3 task common:update:build-assets` regenerates `MainActivity` and drops this patch. Reapply it, or keep the patch in a small overlay the build applies.

## Version from the git tag

```bash
#!/usr/bin/env bash
# v1.13.1-rc.4 -> VERSION_NAME=1.13.1-rc.4 VERSION_CODE=1130104
# v1.13.1      -> VERSION_NAME=1.13.1      VERSION_CODE=1130199 (a final outranks its RCs)
set -euo pipefail
tag="${RELEASE_TAG:-$(git describe --tags --abbrev=0 --match 'v[0-9]*')}"
re='^v([0-9]+)\.([0-9]+)\.([0-9]+)(-rc\.([0-9]+))?$'
[[ $tag =~ $re ]] || { echo "tag '$tag' is not vMAJOR.MINOR.PATCH[-rc.N]" >&2; exit 1; }
maj=${BASH_REMATCH[1]} min=${BASH_REMATCH[2]} pat=${BASH_REMATCH[3]} rc=${BASH_REMATCH[5]:-99}
(( min < 100 && pat < 100 && rc < 100 )) || { echo "minor, patch and rc must be < 100" >&2; exit 1; }
echo "VERSION_NAME=${tag#v}"
echo "VERSION_CODE=$((maj * 1000000 + min * 10000 + pat * 100 + rc))"
```

## Taskfile: build and signed AAB

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  linux: ./build/linux/Taskfile.yml
  android: ./build/android/Taskfile.yml

tasks:
  build:
    cmds:
      - pnpm -C ../web build
      - npm --prefix frontend run build # always refresh the embedded SPA
      - task: "{{.GOOS | default OS}}:build"

  release:android:aab:
    vars:
      KEYDIR:
        sh: echo "${SECRETS_DIR:-$HOME/Secrets}/android"
    preconditions:
      - sh: test -f "{{.KEYDIR}}/upload-keystore.p12"
        msg: "no upload keystore at {{.KEYDIR}}"
    cmds:
      - pnpm -C ../web build
      - npm --prefix frontend run build
      - task: android:bundle # not assemble:aab:release (skips the SPA copy)
    env:
      VERSION_NAME:
        sh: scripts/android-version.sh | sed -n 's/^VERSION_NAME=//p'
      VERSION_CODE:
        sh: scripts/android-version.sh | sed -n 's/^VERSION_CODE=//p'
      ANDROID_KEYSTORE_FILE: '{{.KEYDIR}}/upload-keystore.p12'
      ANDROID_KEYSTORE_PASSWORD:
        sh: cat "{{.KEYDIR}}/keystore-password"
      ANDROID_KEY_ALIAS:
        sh: cat "{{.KEYDIR}}/key-alias"
      ANDROID_KEY_PASSWORD:
        sh: cat "{{.KEYDIR}}/keystore-password"
```

## Tests worth writing

```go
func TestProxyStripsCookiesFromThePage(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "abc", Path: "/"})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer api.Close()
	target, _ := url.Parse(api.URL)

	proxy := &Proxy{Target: target, Jar: NewJar(), Client: api.Client()}
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/login", nil))

	if rec.Header().Get("Set-Cookie") != "" {
		t.Fatal("Set-Cookie leaked to the webview")
	}
	req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	proxy.Jar.AddTo(req)
	if got := req.Header.Get("Cookie"); got != "session=abc" {
		t.Fatalf("Cookie = %q, want session=abc", got)
	}
}
```

Also table-test `stored.pathMatch` (`/` matches everything, `/v1` does not match `/v10`, `/v1/` matches `/v1/me`) and the SPA fallback (`/settings` → 200 HTML, `/missing.js` → 404).

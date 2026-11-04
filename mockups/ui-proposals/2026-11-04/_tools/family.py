"""Shared family layer for the rsk landing proposals.

Tokens, base CSS, the shell (header, footer, theme toggle), and the motion kit
(film engine, scroll-linked reveals, count-ups, section rail, copy button).
Pages import this and add their own composition.

Two palettes, on purpose:

- **Page chrome** is ``@ralvarez-ts/tokens`` (the ``sharp`` layer) — the same
  package the ralvagym, finance, identity and notifier apps consume. The page is
  single-file HTML with no CDN, so it reads a pinned snapshot
  (``tokens.json``, written by ``tokens_extract.py``) instead of the package at
  runtime. ``python3 tokens_extract.py <pkg> --check`` fails when the package
  moves, so the page cannot drift into a hand-copied palette.
- **Terminal surfaces** (``.tn``) keep rsk's real termkit TokyoNight theme
  (github.com/ralvarezdev/termkit theme.go via internal/ui/theme.go), so the
  product inside the frame still looks like the product.

Text contrast of every page token against every page background is asserted
at build; the TokyoNight surface is fixed-dark and checked separately.
"""

from __future__ import annotations

import html
import json
from pathlib import Path

# --- colour maths (port of go-colorful BlendLab toward black, as rsk's darken) ---


def _lin(c: float) -> float:
    c /= 255
    return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4


def _rgb(h: str) -> tuple[int, int, int]:
    h = h.lstrip("#")
    return int(h[0:2], 16), int(h[2:4], 16), int(h[4:6], 16)


def lum(h: str) -> float:
    r, g, b = (_lin(c) for c in _rgb(h))
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def contrast(a: str, b: str) -> float:
    hi, lo = sorted((lum(a), lum(b)), reverse=True)
    return (hi + 0.05) / (lo + 0.05)


def _to_lab(h: str) -> tuple[float, float, float]:
    r, g, b = (_lin(c) for c in _rgb(h))
    x = (0.4124 * r + 0.3576 * g + 0.1805 * b) / 0.95047
    y = 0.2126 * r + 0.7152 * g + 0.0722 * b
    z = (0.0193 * r + 0.1192 * g + 0.9505 * b) / 1.08883

    def f(t: float) -> float:
        return t ** (1 / 3) if t > 216 / 24389 else (24389 / 27 * t + 16) / 116

    fx, fy, fz = f(x), f(y), f(z)
    return 116 * fy - 16, 500 * (fx - fy), 200 * (fy - fz)


def _from_lab(L: float, a: float, b: float) -> str:
    fy = (L + 16) / 116
    fx, fz = fy + a / 500, fy - b / 200

    def inv(t: float) -> float:
        return t**3 if t**3 > 216 / 24389 else (116 * t - 16) / (24389 / 27)

    x, y, z = inv(fx) * 0.95047, inv(fy), inv(fz) * 1.08883
    rgb = (
        3.2406 * x - 1.5372 * y - 0.4986 * z,
        -0.9689 * x + 1.8758 * y + 0.0415 * z,
        0.0557 * x - 0.2040 * y + 1.0570 * z,
    )

    def g(c: float) -> int:
        c = max(0.0, min(1.0, c))
        c = 12.92 * c if c <= 0.0031308 else 1.055 * c ** (1 / 2.4) - 0.055
        return round(c * 255)

    return "#%02x%02x%02x" % tuple(g(c) for c in rgb)


def darken(h: str, t: float) -> str:
    L, a, b = _to_lab(h)
    return _from_lab(L * (1 - t), a * (1 - t), b * (1 - t))


# --- tokens ---

# termkit TokyoNight (Light, Dark) — internal/ui/theme.go: `Theme = termkit.TokyoNight`
TK = {
    "accent": ("#2E7DE9", "#7AA2F7"),
    "danger": ("#F52A65", "#F7768E"),
    "warning": ("#8C6C3E", "#E0AF68"),
    "success": ("#587539", "#9ECE6A"),
    "info": ("#007197", "#7DCFFF"),
    "muted": ("#848CB5", "#565F89"),
}
# rsk wordmark gradient steps — internal/ui/theme.go wordmarkGradientSteps,
# cycled over six rows by termkit banner.go (i % len).
WM_STEPS = [0, 0.18, 0.34, 0.48]

# Ralvarez family brand kit — the real package, pinned. `tokens.json` holds the
# effective sharp value of every role for both themes, extracted from
# @ralvarez-ts/tokens by tokens_extract.py. Regenerate it when the package moves;
# never hand-edit it.
_TOKENS = json.loads((Path(__file__).with_name("tokens.json")).read_text())
_META = _TOKENS["_meta"]
assert _META["package"] == "@ralvarez-ts/tokens", _META

# Roles the package does not own, layered on top: it has no `info` tone, no
# terminal palette, and no ink label for a cyan fill (its `on-accent` is the
# label for `accent-solid`). The device shadow is a page composition choice;
# sharp itself carries no elevation.
APP = {
    "light": {
        "on-accent": "#ffffff",
        "info": "#115c7a",
        "shadow": "0 1px 0 rgba(20,23,26,.05), 0 18px 40px -30px rgba(20,23,26,.45)",
    },
    "dark": {
        "on-accent": "#06161a",
        "info": "#7fdeed",
        "shadow": "0 1px 0 rgba(0,0,0,.35), 0 18px 44px -30px rgba(0,0,0,.9)",
    },
}
LIGHT = {**_TOKENS["light"], **APP["light"]}
DARK = {**_TOKENS["dark"], **APP["dark"]}
# The always-dark band (family `.ink`): darker than the dark theme page.
NIGHT = dict(DARK, bg="#0a0c0f", surface="#161a1f")
# rsk's real terminal surface: termkit TokyoNight (dark), fixed in both themes.
TERM = {
    "bg": "#16161e",
    "bg2": "#1a1b26",
    "surface": "#1f2335",
    "line": "#2a2e45",
    "line-2": "#3b4261",
    "text": "#c0caf5",
    "muted": "#a9b1d6",
    "accent": "#7aa2f7",
    "accent-t": "#7aa2f7",
    "on-accent": "#16161e",
    "ok": "#9ece6a",
    "warn": "#e0af68",
    "danger": "#f7768e",
    "info": "#7dcfff",
    "sel": "#283457",
    "shadow": "0 18px 44px -30px rgba(0,0,0,.9)",
}
for _i in range(6):
    _t = WM_STEPS[_i % len(WM_STEPS)]
    # The wordmark still uses rsk's own gradient, inside the terminal / ink bands.
    LIGHT[f"wm{_i}"] = darken(TK["accent"][0], _t)
    DARK[f"wm{_i}"] = darken(TK["accent"][1], _t)


def _check() -> None:
    for name, th in (("light", LIGHT), ("dark", DARK), ("night", NIGHT)):
        for bg in ("bg", "bg2", "surface"):
            for fg in ("text", "muted", "accent-t", "ok", "warn", "danger", "info"):
                c = contrast(th[fg], th[bg])
                assert c >= 4.5, f"{name}: {fg} on {bg} = {c:.2f}"
        assert contrast(th["on-accent"], th["accent"]) >= 4.5, name
    # The package's promise for --color-border-strong is a >=3:1 control
    # boundary (WCAG 1.4.11). The old hand-copied hairline was 1.5:1, so this
    # assertion is the reason the page takes the role from the package at all.
    for name, th in (("light", LIGHT), ("dark", DARK)):
        c = contrast(th["line-2"], th["bg"])
        assert c >= 3.0, f"{name}: control boundary {c:.2f}"
    for bg in ("bg", "bg2", "surface"):
        for fg in ("text", "muted", "accent-t", "ok", "warn", "danger", "info"):
            c = contrast(TERM[fg], TERM[bg])
            assert c >= 4.5, f"term: {fg} on {bg} = {c:.2f}"


_check()


def _vars(th: dict) -> str:
    return ";".join(f"--{k}:{v}" for k, v in th.items())


FONTS = """
--mono:ui-monospace,"JetBrains Mono","Cascadia Mono","SF Mono",Menlo,Consolas,"DejaVu Sans Mono","Liberation Mono",monospace;
--sans:system-ui,-apple-system,"Segoe UI",Roboto,"Noto Sans","Helvetica Neue",Arial,sans-serif;
"""

BASE_CSS = (
    ":root{color-scheme:light;--s1:4px;--s2:8px;--s3:12px;--s4:16px;--s5:24px;--s6:32px;--s7:48px;--s8:72px;--s9:112px;"
    "--wrap:1200px;--gut:clamp(20px,4vw,56px);"
    f"--rk:{_TOKENS['radii']['control']};--rc:{_TOKENS['radii']['card']};"
    "--e-out:cubic-bezier(.16,1,.3,1);--e-io:cubic-bezier(.65,0,.35,1);--d1:160ms;--d2:320ms;--d3:640ms;--d4:1100ms;"
    + FONTS + _vars(LIGHT) + "}"
    "@media (prefers-color-scheme:dark){:root:not([data-theme]){color-scheme:dark;" + _vars(DARK) + "}}"
    "[data-theme=light]{color-scheme:light;" + _vars(LIGHT) + "}"
    "[data-theme=dark]{color-scheme:dark;" + _vars(DARK) + "}"
    ".night{color-scheme:dark;" + _vars(NIGHT) + ";background:var(--bg);color:var(--text)}"
    ".tn{color-scheme:dark;" + _vars(TERM) + ";background:var(--bg);color:var(--text)}"
    + r"""
*,*::before,*::after{box-sizing:border-box}
html{-webkit-text-size-adjust:100%;scroll-behavior:smooth;scroll-padding-top:88px}
@media (prefers-reduced-motion:reduce){html{scroll-behavior:auto}}
body{margin:0;background:var(--bg);color:var(--text);font:400 1.0625rem/1.6 var(--sans);overflow-x:clip;-webkit-font-smoothing:antialiased;text-rendering:optimizeLegibility}
img,svg{max-width:100%}
h1,h2,h3{line-height:.99;margin:0;text-wrap:balance;letter-spacing:-.04em;font-weight:800}
p{margin:0 0 1em;text-wrap:pretty}
a{color:var(--accent-t);text-underline-offset:3px;text-decoration-thickness:1px}
a:hover{text-decoration-thickness:2px}
code,pre,kbd,.mono{font-family:var(--mono)}
code{font-size:.9em;overflow-wrap:anywhere}
h1 code,h2 code,h3 code{white-space:nowrap;overflow-wrap:normal}
pre{margin:0;white-space:pre-wrap;overflow-wrap:anywhere;font-size:.86rem;line-height:1.6}
.num{font-variant-numeric:tabular-nums}
::selection{background:var(--accent);color:var(--on-accent)}
:focus-visible{outline:2px solid var(--accent-t);outline-offset:3px;border-radius:2px}
.night :focus-visible,.tn :focus-visible{outline-color:var(--accent-t)}
.wrap{width:100%;max-width:var(--wrap);margin:0 auto;padding:0 var(--gut)}
.sr{position:absolute!important;width:1px;height:1px;overflow:hidden;clip:rect(0 0 0 0);clip-path:inset(50%);white-space:nowrap;margin:-1px;padding:0;border:0}
.skip{position:absolute;left:12px;top:-80px;z-index:50;background:var(--accent-t);color:var(--on-accent);padding:12px 16px;border-radius:var(--rk);font-weight:700;min-height:44px;text-decoration:none}
.skip:focus{top:12px}
/* family micro-detail: the brand-kit square and dot (12px, radius 2.6 at 64-unit scale) */
.fd{display:inline-grid;grid-template-columns:repeat(2,5px);gap:2px;vertical-align:middle}
.fd i{width:5px;height:5px;border-radius:1px;background:currentColor}
.fd i:last-child{border-radius:50%;background:var(--accent)}
/* header + scroll progress */
.k-top{position:sticky;top:0;z-index:20;background:color-mix(in srgb,var(--bg) 88%,transparent);backdrop-filter:saturate(1.4) blur(12px);-webkit-backdrop-filter:saturate(1.4) blur(12px);border-bottom:1px solid transparent;transition:border-color var(--d2)}
.k-top.stuck{border-bottom-color:var(--line)}
.k-top .wrap{display:flex;align-items:center;gap:8px;min-height:64px}
.k-brand{display:inline-flex;align-items:center;gap:10px;min-height:44px;color:var(--text);text-decoration:none;font:700 1.05rem var(--mono);letter-spacing:-.01em}
.k-brand .fd{color:var(--text)}
.k-brand small{font:500 .78rem var(--mono);color:var(--muted);letter-spacing:0}
.k-nav{margin-left:auto;display:flex;align-items:center;gap:2px}
.k-nav a{display:inline-flex;align-items:center;min-height:44px;padding:0 12px;color:var(--muted);text-decoration:none;font:500 .9rem var(--sans);border-radius:var(--rk);border-bottom:2px solid transparent}
.k-nav a:hover{color:var(--text);background:var(--sel)}
.k-nav a[aria-current]{color:var(--text);border-bottom-color:var(--accent)}
@media (max-width:760px){.k-nav a[data-opt]{display:none}}
.k-prog{position:absolute;left:0;right:0;bottom:-1px;height:2px;pointer-events:none;background:color-mix(in srgb,var(--line) 45%,transparent)}
.k-prog i{display:block;height:100%;background:var(--accent);transform:scaleX(var(--p,0));transform-origin:0 50%;transition:transform 90ms linear}
html:not(.js) .k-prog{display:none}
/* buttons */
.k-btn{display:inline-flex;align-items:center;justify-content:center;gap:8px;min-height:44px;min-width:44px;padding:0 18px;border-radius:var(--rk);border:1px solid var(--line-2);background:transparent;color:var(--text);font:600 .95rem/1 var(--sans);cursor:pointer;text-decoration:none;white-space:nowrap;transition:background var(--d1),border-color var(--d1),color var(--d1)}
.k-btn:hover{border-color:var(--text)}
.k-btn.primary{background:var(--accent-t);border-color:var(--accent-t);color:var(--on-accent)}
.k-btn.primary:hover{background:color-mix(in srgb,var(--accent-t) 86%,var(--text));border-color:color-mix(in srgb,var(--accent-t) 86%,var(--text))}
.k-btn.ghost{background:transparent}
.anim .k-btn:active{transform:translateY(1px)}
.k-theme{margin-left:6px;padding:0 12px}
.k-theme .ico{font-size:1.1rem;line-height:1}
@media (max-width:420px){.k-theme .t{position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%)}}
/* kicker + section rhythm */
.k-kick{display:inline-flex;align-items:center;gap:10px;font:600 .74rem/1.5 var(--mono);letter-spacing:.1em;text-transform:uppercase;color:var(--muted);margin-bottom:14px}
.k-kick .fd{margin:-2px 2px 0 0}
.k-kick .gl{color:var(--accent-t)}
.k-sec{position:relative;padding:clamp(72px,9vw,128px) 0}
.k-sec h2{font-size:clamp(2rem,4.4vw,3.4rem);margin:.12em 0 .5em;max-width:22ch}
.k-lede{font-size:clamp(1.05rem,1.5vw,1.22rem);color:var(--muted);max-width:62ch}
.k-rule{height:2px;background:var(--text);border:0;margin:0 0 clamp(28px,4vw,48px)}
.anim .k-rule{transform:scaleX(0);transform-origin:0 50%;transition:transform .95s var(--e-out)}
.anim .k-rule.in{transform:none}
/* hero grid texture (decorative) */
.k-grid{position:absolute;inset:0;z-index:-1;pointer-events:none;overflow:hidden;
 background-image:linear-gradient(to right,color-mix(in srgb,var(--line) 58%,transparent) 1px,transparent 1px),linear-gradient(to bottom,color-mix(in srgb,var(--line) 58%,transparent) 1px,transparent 1px);
 background-size:72px 72px;
 mask-image:radial-gradient(120% 90% at 78% 0%,#000 0%,transparent 68%);
 -webkit-mask-image:radial-gradient(120% 90% at 78% 0%,#000 0%,transparent 68%)}
/* badges */
.k-badge{display:inline-flex;align-items:center;gap:6px;padding:1px 7px;border:1px solid currentColor;border-radius:3px;font:600 .68rem/1.7 var(--mono);letter-spacing:.06em;text-transform:uppercase;white-space:nowrap;vertical-align:1px}
.k-badge.ro{color:var(--ok)}.k-badge.de{color:var(--danger)}.k-badge.hu{color:var(--warn)}.k-badge.sa{color:var(--muted);border-style:dashed}
/* code window */
.k-win{background:var(--surface);border:1px solid var(--line);border-radius:var(--rc);box-shadow:var(--shadow);overflow:hidden}
.k-win-h{display:flex;align-items:center;gap:10px;padding:9px 14px;border-bottom:1px solid var(--line);font:600 .72rem var(--mono);letter-spacing:.06em;text-transform:uppercase;color:var(--muted);background:color-mix(in srgb,var(--surface) 80%,var(--bg))}
.k-win-h .dots{display:inline-flex;gap:4px}.k-win-h .dots i{width:6px;height:6px;border-radius:1px;background:currentColor;opacity:.65;display:block}
.k-win-b{padding:14px 16px}
/* copy block */
.k-copy{position:relative}
.k-copy pre{padding:18px 20px;padding-right:96px;background:var(--surface);border:1px solid var(--line);border-radius:var(--rc)}
.k-copy .k-btn{position:absolute;top:8px;right:8px;min-height:40px;padding:0 14px;background:var(--surface)}
.k-cm{color:var(--muted)}.k-c{display:block;margin-bottom:.55em}
/* film control: label + a real progress bar */
.k-film-ctl{display:inline-flex;align-items:center;gap:10px;min-height:40px;min-width:44px;padding:0 12px;border-radius:var(--rk);border:1px solid var(--line-2);background:transparent;color:var(--muted);font:600 .72rem var(--mono);letter-spacing:.06em;text-transform:uppercase;cursor:pointer}
.k-film-ctl:hover{color:var(--text);border-color:var(--text)}
.k-film-ctl .ic{display:inline-block;width:10px;text-align:center}
.k-film-ctl .bar{width:42px;height:2px;background:var(--line);position:relative;overflow:hidden}
.k-film-ctl .bar i{position:absolute;inset:0;background:var(--accent);transform:scaleX(var(--p,1));transform-origin:0 50%}
/* section rail */
.k-rail{display:none}
@media (min-width:1400px){
 .k-rail{display:block;position:fixed;left:18px;top:50%;transform:translateY(-50%);z-index:15}
 .k-rail ol{list-style:none;margin:0;padding:0;display:grid;gap:2px}
 .k-rail a{display:flex;align-items:center;gap:10px;min-height:30px;padding:2px 6px;color:var(--muted);text-decoration:none;font:500 .7rem var(--mono);border-radius:3px}
 .k-rail a::before{content:"";width:12px;height:2px;background:var(--line-2);transition:width var(--d2) var(--e-out),background var(--d2)}
 .k-rail a .n{color:var(--muted);opacity:0;transform:translateX(-4px);transition:opacity var(--d2),transform var(--d2) var(--e-out)}
 .k-rail a .lb{opacity:0;transform:translateX(-4px);transition:opacity var(--d2),transform var(--d2) var(--e-out);white-space:nowrap;background:var(--bg);padding:3px 6px;border:1px solid var(--line);border-radius:3px;color:var(--text)}
 .k-rail:hover a .lb,.k-rail a:focus-visible .lb,.k-rail:hover a .n,.k-rail a:focus-visible .n,.k-rail a[aria-current] .n{opacity:1;transform:none}
 .k-rail a[aria-current]::before{width:28px;background:var(--accent)}
 .k-rail a[aria-current] .n{color:var(--accent-t)}
}
/* reveals: content visible unless motion is allowed; one-shot, ends at 1 */
.anim .rv{opacity:0;transform:translateY(20px);transition:opacity .55s var(--e-out),transform .8s var(--e-out);transition-delay:calc(var(--i,0) * 70ms)}
.anim .rv.in{opacity:1;transform:none}
.anim .rv--l{transform:translateX(-26px)}
.anim .rv--s{transform:scale(.975)}
.anim .rv--d{transform:translateY(32px)}
/* film: hidden steps keep their box (no layout shift) */
.f-run .f-h{visibility:hidden}
.anim .f-run .f-h{visibility:visible;opacity:0;transform:translateY(6px)}
.anim .f-run [data-f]{transition:opacity .35s ease,transform .45s var(--e-out),background-color .5s ease}
.f-caret::after{content:"";display:inline-block;width:.55em;height:1.05em;margin-left:1px;vertical-align:-.15em;background:currentColor;animation:k-blink 1s steps(1) infinite}
@keyframes k-blink{50%{opacity:0}}
/* scroll-linked helpers (Chromium 115+; progressive elsewhere) */
@supports (animation-timeline:view()){
 .anim .par{animation:k-par linear both;animation-timeline:view();animation-range:entry 0% cover 70%}
 @keyframes k-par{from{transform:translateY(26px)}to{transform:translateY(-26px)}}
 .anim .drift-x{animation:k-drift-x linear both;animation-timeline:view();animation-range:entry 10% exit 60%}
 @keyframes k-drift-x{from{transform:translateX(-4%)}to{transform:translateX(3%)}}
}
/* terminal surface: termkit TokyoNight, fixed in both themes */
.tn{--wm0:#2e7de9;--wm1:#2d65bb;--wm2:#2a5193;--wm3:#264072;--wm4:#2e7de9;--wm5:#2d65bb}
.tn .k-win-h{background:color-mix(in srgb,var(--surface) 80%,#000);border-color:var(--line)}
/* trade-off */
.k-trade{display:grid;gap:clamp(24px,4vw,64px);grid-template-columns:1fr}
@media (min-width:900px){.k-trade{grid-template-columns:1.1fr 1fr;align-items:start}}
.k-trade .give,.k-trade .take{border-top:2px solid var(--line);padding-top:20px}
.k-trade .give{border-color:var(--ok)}.k-trade .take{border-color:var(--warn)}
.k-trade h3{font-size:clamp(1.4rem,2.6vw,2rem);margin-bottom:.5em}
.k-trade ul{margin:0;padding-left:1.1em}.k-trade li{margin:.35em 0}
/* footer */
.k-foot{border-top:1px solid var(--line);padding:56px 0 64px;font-size:.92rem;color:var(--muted)}
.k-foot .wrap{display:grid;gap:28px}
@media (min-width:900px){.k-foot .wrap{grid-template-columns:auto 1fr;align-items:start}}
.k-foot details{max-width:72ch}
.k-foot summary{cursor:pointer;min-height:44px;display:flex;align-items:center;font-weight:600;color:var(--text)}
.k-src{margin:8px 0 0;padding:0;list-style:none;display:grid;gap:6px}
.k-src li{display:grid;grid-template-columns:minmax(0,1fr);gap:0}
@media (min-width:640px){.k-src li{grid-template-columns:minmax(0,17ch) minmax(0,1fr);gap:16px}}
.k-wm{display:block;width:180px;height:auto}
.wm0{fill:var(--wm0)}.wm1{fill:var(--wm1)}.wm2{fill:var(--wm2)}.wm3{fill:var(--wm3)}.wm4{fill:var(--wm4)}.wm5{fill:var(--wm5)}
/* view-transition theme wipe */
::view-transition-old(root),::view-transition-new(root){animation:none;mix-blend-mode:normal}
@media (forced-colors:active){.k-badge,.k-btn,.k-film-ctl,.k-win{border:1px solid CanvasText}.k-btn.primary{forced-color-adjust:none;background:Highlight;color:HighlightText}.fd i,.k-rule,.k-prog i,.k-rail a::before{background:CanvasText;forced-color-adjust:none}}
"""
)

HEAD_JS = r"""(function(){var d=document.documentElement,s=null;try{s=localStorage.getItem('rsk-theme')}catch(e){}
d.setAttribute('data-theme',s||(matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'));
if(!matchMedia('(prefers-reduced-motion: reduce)').matches)d.classList.add('anim');d.classList.add('js')})();"""

KIT_JS = r"""
(function(){
var d=document.documentElement, motion=d.classList.contains('anim');
/* theme toggle: stored preference, system default, circular wipe when motion allowed */
var tb=document.querySelector('.k-theme');
function cur(){return d.getAttribute('data-theme')}
function sync(){if(tb)tb.setAttribute('aria-pressed',cur()==='dark'?'true':'false')}
if(tb){tb.hidden=false;sync();tb.addEventListener('click',function(){
  var next=cur()==='dark'?'light':'dark';
  var apply=function(){d.setAttribute('data-theme',next);try{localStorage.setItem('rsk-theme',next)}catch(e){};sync()};
  if(motion&&document.startViewTransition){var r=tb.getBoundingClientRect(),x=r.left+r.width/2,y=r.top+r.height/2,
    R=Math.hypot(Math.max(x,innerWidth-x),Math.max(y,innerHeight-y));
    var vt=document.startViewTransition(apply);
    vt.ready.then(function(){d.animate({clipPath:['circle(0px at '+x+'px '+y+'px)','circle('+R+'px at '+x+'px '+y+'px)']},{duration:520,easing:'cubic-bezier(.16,1,.3,1)',pseudoElement:'::view-transition-new(root)'})}).catch(function(){});
  }else apply();
})}
/* sticky header hairline + reading progress */
var top=document.querySelector('.k-top'),bar=document.querySelector('.k-prog i'),queued=false;
function measure(){
  queued=false;
  if(top)top.classList.toggle('stuck',scrollY>8);
  if(bar){var h=d, max=h.scrollHeight-innerHeight, p=max>0?Math.min(1,Math.max(0,scrollY/max)):0;bar.style.setProperty('--p',p.toFixed(4))}
}
function onScroll(){if(!queued){queued=true;requestAnimationFrame(measure)}}
addEventListener('scroll',onScroll,{passive:true});addEventListener('resize',onScroll,{passive:true});measure();
/* one-shot reveals + rule draws + section current */
var targets=[].slice.call(document.querySelectorAll('.rv,.k-rule'));
if(motion&&'IntersectionObserver' in window){var io=new IntersectionObserver(function(es){es.forEach(function(e){if(e.isIntersecting){e.target.classList.add('in');io.unobserve(e.target)}})},{rootMargin:'0px 0px -8% 0px'});targets.forEach(function(el){io.observe(el)})}
else targets.forEach(function(el){el.classList.add('in')});
var links=[].slice.call(document.querySelectorAll('.k-nav a[href^="#"],.k-rail a[href^="#"]'));
if(links.length&&'IntersectionObserver' in window){
  var byId={};links.forEach(function(a){var id=a.getAttribute('href').slice(1);(byId[id]=byId[id]||[]).push(a)});
  var secs=Object.keys(byId).map(function(id){return document.getElementById(id)}).filter(Boolean);
  var navio=new IntersectionObserver(function(es){es.forEach(function(e){if(!e.isIntersecting)return;
    links.forEach(function(a){a.removeAttribute('aria-current')});
    (byId[e.target.id]||[]).forEach(function(a){a.setAttribute('aria-current','true')})})},{rootMargin:'-45% 0px -50% 0px'});
  secs.forEach(function(s){navio.observe(s)});
}
/* count-ups: the plain HTML already holds the real number, so JS off and
   reduced motion both show it; only the motion path starts from zero. */
var counts=[].slice.call(document.querySelectorAll('[data-count]'));
function finishCount(el){el.textContent=el.getAttribute('data-final');el.setAttribute('data-done','1')}
counts.forEach(function(el){
  el.setAttribute('data-final',el.textContent.trim());
  if(!motion||!('IntersectionObserver' in window)){el.setAttribute('data-done','1');return}
  var o=new IntersectionObserver(function(es){if(!es[0].isIntersecting)return;o.disconnect();
    var to=+el.getAttribute('data-count')||0,from=+(el.getAttribute('data-from')||0),t0=null,dur=1200;
    function frame(now){if(t0===null)t0=now;var p=Math.min(1,(now-t0)/dur),e=1-Math.pow(1-p,4);
      el.textContent=Math.round(from+(to-from)*e).toLocaleString('en-US');
      if(p<1)requestAnimationFrame(frame);else finishCount(el)}
    requestAnimationFrame(frame);
  },{threshold:.55});
  o.observe(el);
});
window.kCountsDone=function(){counts.forEach(function(el){if(el.getAttribute('data-done')!=='1')finishCount(el)})};
/* copy buttons */
document.querySelectorAll('[data-copy]').forEach(function(b){b.addEventListener('click',function(){
  var src=document.getElementById(b.getAttribute('data-copy')),txt=src?src.getAttribute('data-raw')||src.textContent:'';
  var lbl=b.querySelector('.l'),live=document.getElementById('k-live');
  function done(ok){lbl.textContent=ok?'Copied':'Copy failed';if(live)live.textContent=ok?'Commands copied to the clipboard':'Copy failed; select the text instead';setTimeout(function(){lbl.textContent='Copy'},1800)}
  try{navigator.clipboard.writeText(txt).then(function(){done(true)},function(){done(false)})}catch(e){done(false)}
})});
/* film engine: steps [{at,dur,fn(p)}], reset() = start state, end() = cleanup; the plain HTML is the final state */
window.kFilm=function(root,cfg){
  var btn=root.querySelector('.k-film-ctl'); if(!btn)return; root.setAttribute('data-film','');
  var ic=btn.querySelector('.ic'),lb=btn.querySelector('.l');
  var t=0,last=0,raf=0,playing=false,started=false,done=false;
  function set(state){btn.setAttribute('data-state',state);
    if(state==='playing'){ic.textContent='❚❚';lb.textContent='Pause'}
    else if(state==='paused'){ic.textContent='▶';lb.textContent='Play'}
    else{ic.textContent='↻';lb.textContent='Replay'}}
  function apply(){cfg.steps.forEach(function(s){if(t>=s.at){var p=s.dur?Math.min(1,(t-s.at)/s.dur):1;if(s._p!==p){s._p=p;s.fn(p)}}});btn.style.setProperty('--p',Math.min(1,t/cfg.total).toFixed(4))}
  function frame(now){if(!playing)return;t+=Math.min(64,now-last);last=now;apply();if(t>=cfg.total){finish();return}raf=requestAnimationFrame(frame)}
  function reset(){t=0;done=false;started=true;btn.style.setProperty('--p','0');cfg.steps.forEach(function(s){s._p=undefined});root.classList.add('f-run');cfg.reset()}
  function play(){if(!started||done)reset();playing=true;last=performance.now();set('playing');raf=requestAnimationFrame(frame)}
  function pause(){playing=false;cancelAnimationFrame(raf);set('paused')}
  function finish(){playing=false;done=true;cancelAnimationFrame(raf);cfg.steps.forEach(function(s){if(s._p!==1){s._p=1;s.fn(1)}});cfg.end();root.classList.remove('f-run');btn.style.setProperty('--p','1');set('done')}
  btn.hidden=false;set('done');
  btn.addEventListener('click',function(){playing?pause():play()});
  document.addEventListener('visibilitychange',function(){if(document.hidden&&playing)pause()});
  root.kFilm={finish:finish,play:play,pause:pause,state:function(){return btn.getAttribute('data-state')}};
  if(motion&&'IntersectionObserver' in window){var o=new IntersectionObserver(function(es){if(es[0].isIntersecting&&!started){o.disconnect();play()}},{threshold:.3});o.observe(root)}
};
/* film helpers */
window.kF={
  hide:function(root){root.querySelectorAll('[data-f]').forEach(function(e){e.classList.add('f-h')})},
  show:function(root,k){root.querySelectorAll('[data-f="'+k+'"]').forEach(function(e){e.classList.remove('f-h')})},
  type:function(el,p){var full=el.getAttribute('data-t');el.textContent=full.slice(0,Math.round(full.length*p));el.classList.toggle('f-caret',p<1)},
  end:function(root){root.querySelectorAll('.f-h').forEach(function(e){e.classList.remove('f-h')});root.querySelectorAll('[data-t]').forEach(function(e){e.textContent=e.getAttribute('data-t');e.classList.remove('f-caret')})}
};
})();
"""

WORDMARK = """██████╗  ███████╗ ██╗  ██╗
██╔══██╗ ██╔════╝ ██║ ██╔╝
██████╔╝ ███████╗ █████╔╝
██╔══██╗ ╚════██║ ██╔═██╗
██║  ██║ ███████║ ██║  ██╗
╚═╝  ╚═╝ ╚══════╝ ╚═╝  ╚═╝"""


def wordmark_svg(cls: str = "k-wm", label: str = "rsk") -> str:
    """The real rsk banner wordmark (internal/ui/theme.go), one SVG text row per
    line, shaded with the CLI's own row gradient."""
    rows = WORDMARK.split("\n")
    w, lh = 26 * 12, 22
    out = [
        f'<svg class="{cls}" viewBox="0 -16 {w} {lh * len(rows) + 2}" role="img" aria-label="{esc(label)}">'
    ]
    for i, r in enumerate(rows):
        out.append(
            f'<text class="wm{i}" x="0" y="{i * lh}" font-family="DejaVu Sans Mono,Menlo,Consolas,monospace" '
            f'font-size="20" textLength="{len(r.rstrip()) * 12}" lengthAdjust="spacingAndGlyphs" xml:space="preserve">{esc(r.rstrip())}</text>'
        )
    out.append("</svg>")
    return "".join(out)


def esc(s: str) -> str:
    return html.escape(s, quote=True)


def header(nav: list[tuple[str, str, bool]], version: str) -> str:
    links = "".join(
        f'<a href="{h}"{" data-opt" if opt else ""}>{esc(t)}</a>' for t, h, opt in nav
    )
    return f"""<a class="skip" href="#main">Skip to content</a>
<header class="k-top"><div class="wrap">
<a class="k-brand" href="#top"><span class="fd" aria-hidden="true"><i></i><i></i></span>rsk <small>{esc(version)}</small></a>
<nav class="k-nav" aria-label="Page">{links}
<button class="k-btn ghost k-theme" type="button" aria-pressed="false" hidden><span class="ico" aria-hidden="true">◐</span><span class="t">Dark theme</span></button></nav>
<div class="k-prog" aria-hidden="true"><i></i></div>
</div></header>"""


def rail(items: list[tuple[str, str]]) -> str:
    """Persistent section rail (family style): a tick per section, active by
    IntersectionObserver, label on hover/focus. Hidden below 1400 px."""
    links = "".join(
        f'<li><a href="#{i}"><span class="n">{n:02d}</span><span class="lb">{esc(t)}</span></a></li>'
        for n, (i, t) in enumerate(items, 1)
    )
    return f'<nav class="k-rail" aria-label="Sections"><ol>{links}</ol></nav>'


def copy_block(id_: str, lines: list[tuple[str, str]], label: str) -> str:
    """lines: (command, comment). Raw text (commands only) goes to the clipboard."""
    raw = "\n".join(c for c, _ in lines if c)
    body = []
    for c, cm in lines:
        if cm:
            body.append(f"<span class=k-cm># {esc(cm)}</span>")
        if c:
            body.append(f"<span class=k-c>{esc(c)}</span>")
    return f"""<div class="k-copy"><pre id="{id_}" data-raw="{esc(raw)}" aria-label="{esc(label)}" role="region" tabindex="0">{chr(10).join(body)}</pre>
<button class="k-btn" type="button" data-copy="{id_}"><span class="l">Copy</span><span class="sr"> {esc(label)}</span></button></div>"""


def film_ctl() -> str:
    return '<button class="k-film-ctl" type="button" hidden><span class="ic" aria-hidden="true">↻</span><span class="l">Replay</span><span class="bar" aria-hidden="true"><i></i></span><span class="sr"> the hero animation</span></button>'


def footer(sources: list[tuple[str, str]], version: str) -> str:
    items = "".join(
        f"<li><span>{esc(c)}</span><code>{esc(s)}</code></li>" for c, s in sources
    )
    return f"""<footer class="k-foot"><div class="wrap">
<div>{wordmark_svg()}<p style="margin-top:14px">rsk {esc(version)} · GPL-3.0 · <a href="https://github.com/ralvarezdev/ralvaskills">github.com/ralvarezdev/ralvaskills</a></p></div>
<details><summary>Where every number on this page comes from</summary>
<p>Counted in the repository on 2026-11-04. Anything marked <span class="k-badge sa">sample</span> is illustrative.</p>
<ul class="k-src">{items}</ul></details>
</div></footer>"""


def page(title: str, desc: str, css: str, body: str, js: str, extra_head: str = "") -> str:
    return f"""<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{esc(title)}</title>
<meta name="description" content="{esc(desc)}">
{extra_head}<meta name="color-scheme" content="light dark">
<script>{HEAD_JS}</script>
<style>{BASE_CSS}{css}</style>
</head>
<body>
{body}
<p id="k-live" class="sr" aria-live="polite"></p>
<script>{KIT_JS}</script>
<script>{js}</script>
</body>
</html>
"""

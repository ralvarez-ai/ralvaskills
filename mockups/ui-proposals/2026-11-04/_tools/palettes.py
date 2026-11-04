#!/usr/bin/env python3
"""Colour-palette proposals for the rsk landing page.

Reads the approved combined page (direction-4-combined.html) and writes one
copy per palette, plus palette-viewer.html to compare them. The mockup's colour
is entirely variables (the only literals are three #000), so a palette is an
override block, not a rewrite. The source file is never modified.

    python3 -I _tools/palettes.py

Every palette is checked against WCAG 2.2 before it is written: text, muted,
accent and the four status colours must reach 4.5:1 on the page, sunken and card
surfaces; the label on an accent fill 4.5:1; a strong border 3:1 (1.4.11).
The script exits 1 if any palette fails, so a proposal cannot ship broken.
"""
import html
import re
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent.parent
SRC = HERE / "direction-4-combined.html"

# bg=page, sunken=alternate bands, card=raised surface, strong=control border.
# human: optional colour for the "You approve." headline line only.
PALETTES = {
    "ember": {
        "title": "Ember",
        "idea": "Warm paper and a vermilion accent. Editorial rather than console: it steps away from the cyan-on-grey look of most developer tools.",
        "light": dict(bg="#f7f2ea", sunken="#efe8dc", card="#fffdf9", border="#e2d9c8", strong="#7d7260", text="#1c1710", muted="#5a5040", accent="#b3360c", on="#ffffff", tint="rgba(179,54,12,.10)", ok="#1d6a3b", warn="#7a4f00", danger="#a3261f", info="#8a3a00"),
        "dark": dict(bg="#14100c", sunken="#1c1711", card="#221c15", border="#372e23", strong="#76684f", text="#f4ede0", muted="#b8ab93", accent="#ff7a45", on="#1a0b04", tint="rgba(255,122,69,.15)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#ffb48a"),
    },
    "aurora": {
        "title": "Aurora",
        "idea": "Violet over a cool, slightly tinted ground, with mint for success. Fresh and clearly a product, but the most common look among AI tools.",
        "light": dict(bg="#f6f5fb", sunken="#edebf7", card="#ffffff", border="#dedaf0", strong="#706b9c", text="#17152b", muted="#524f77", accent="#5a3ed8", on="#ffffff", tint="rgba(90,62,216,.10)", ok="#0c6a4c", warn="#7a4d00", danger="#b0242c", info="#2d4fc0"),
        "dark": dict(bg="#0d0b1a", sunken="#141128", card="#1a1633", border="#2d2851", strong="#6a6598", text="#f0eeff", muted="#aeaad9", accent="#a994ff", on="#130a33", tint="rgba(169,148,255,.16)", ok="#5eead4", warn="#fbbf24", danger="#ff90a4", info="#8bd3ff"),
    },
    "signal": {
        "title": "Signal",
        "idea": "Near-black ink and an electric lime. The loudest in dark mode; in light the lime cannot carry text, so it becomes a deep forest green and loses most of the effect.",
        "light": dict(bg="#f3f4ee", sunken="#e8eadf", card="#ffffff", border="#d8dccb", strong="#6a705c", text="#11140d", muted="#4b5140", accent="#2c6604", on="#ffffff", tint="rgba(44,102,4,.10)", ok="#0e6b5c", warn="#7a4d00", danger="#a8281f", info="#1f5f8a"),
        "dark": dict(bg="#0a0c08", sunken="#11140d", card="#171b12", border="#292f20", strong="#6b765d", text="#eff3e6", muted="#a9b29a", accent="#b8f23a", on="#10200a", tint="rgba(184,242,58,.14)", ok="#5eead4", warn="#ffb454", danger="#ff8a80", info="#7dd3fc"),
    },
}
# Ember again, with a cool counterpoint for the human: "You approve." leaves the
# amber of "writes" and takes a deep teal, so the two warm headline lines stop
# competing. Only that line changes; the warning amber keeps its meaning.
PALETTES["ember-duo"] = {
    "title": "Ember + counterpoint",
    "idea": "Ember with the human given its own cool colour: “You approve.” is teal, “rsk installs.” vermilion, and the amber stays reserved for warnings.",
    "light": dict(PALETTES["ember"]["light"], human="#0b6560"),
    "dark": dict(PALETTES["ember"]["dark"], human="#5fd6c8"),
}


def _lum(h):
    h = h.lstrip("#")
    r, g, b = (int(h[i:i + 2], 16) / 255 for i in (0, 2, 4))
    f = lambda c: c / 12.92 if c <= 0.03928 else ((c + 0.055) / 1.055) ** 2.4
    return 0.2126 * f(r) + 0.7152 * f(g) + 0.0722 * f(b)


def contrast(a, b):
    la, lb = sorted((_lum(a), _lum(b)), reverse=True)
    return (la + 0.05) / (lb + 0.05)


def mix(h, to, t):
    h, to = h.lstrip("#"), to.lstrip("#")
    c = [round(int(h[i:i + 2], 16) * (1 - t) + int(to[i:i + 2], 16) * t) for i in (0, 2, 4)]
    return "#%02x%02x%02x" % tuple(c)


def check(p):
    """Return a list of failure strings; empty means the palette is sound."""
    bad = []
    for scheme in ("light", "dark"):
        v = p[scheme]
        for surf in ("bg", "sunken", "card"):
            for k in ("text", "muted", "accent", "ok", "warn", "danger", "info") + (("human",) if "human" in v else ()):
                if (c := contrast(v[k], v[surf])) < 4.5:
                    bad.append(f"{scheme} {k} on {surf}: {c:.2f} < 4.5")
            if (c := contrast(v["strong"], v[surf])) < 3:
                bad.append(f"{scheme} strong border on {surf}: {c:.2f} < 3")
        if (c := contrast(v["on"], v["accent"])) < 4.5:
            bad.append(f"{scheme} label on accent: {c:.2f} < 4.5")
        if "human" in v and (c := contrast(v["bg"], v["human"])) < 4.5:
            bad.append(f"{scheme} page-coloured label on human fill: {c:.2f} < 4.5")
    return bad


def wordmark(accent, ground, dark):
    """The six wordmark rows: the accent, then steps toward the ground."""
    toward = "#000000" if not dark else ground
    steps = (0, 0.19, 0.36, 0.5, 0, 0.19)
    return [mix(accent, toward, t) for t in steps]


def scheme_block(sel, v, dark):
    wm = wordmark(v["accent"], v["bg"], dark)
    out = (f"--bg:{v['bg']};--bg2:{v['sunken']};--surface:{v['card']};--line:{v['border']};--line-2:{v['strong']};"
           f"--text:{v['text']};--muted:{v['muted']};--accent:{v['accent']};--accent-t:{v['accent']};--on-accent:{v['on']};"
           f"--sel:{v['tint']};--ok:{v['ok']};--warn:{v['warn']};--danger:{v['danger']};--info:{v['info']};")
    out += "".join(f"--wm{i}:{c};" for i, c in enumerate(wm))
    return f"{sel}{{{out}}}"


def override_css(name, p):
    L, D = p["light"], p["dark"]
    # Terminal and the night band keep the dark palette in both themes, as the
    # original does; the terminal sits one step off the page so it reads as a window.
    term = dict(D, bg=D["sunken"], sunken=D["bg"])
    css = [
        scheme_block("html[data-theme=light]", L, False),
        scheme_block("html[data-theme=dark]", D, True),
        scheme_block("html .night", D, True),
        scheme_block("html .tn", term, True),
    ]
    for scheme, sel in (("light", "html[data-theme=light]"), ("dark", "html[data-theme=dark]")):
        if "human" in p[scheme]:
            h = p[scheme]["human"]
            # Everything in the hero that means "the human": the headline line, the
            # YOU lane, the request arrow, the approval dialog and its button, and
            # the human badge. Other uses of --warn (WRITES labels, ...) stay amber.
            css.append(
                f"{sel} .c-h1 .you,{sel} .c-lane.y,{sel} .c-arrow.ask,{sel} .c-you .lbl,{sel} .k-badge.hu{{color:{h}}}\n"
                f"{sel} .c-lane.y,{sel} .c-you{{border-color:{h}}}\n"
                f"{sel} .c-you .b.al{{border-color:{h};background:{h};color:var(--bg)}}"
            )
    return f"/* palette: {name} (generated by _tools/palettes.py) */\n" + "\n".join(css)


THEME_PARAM = ("<script>/* ?theme=light|dark picks the theme for this visit (used by the viewer) */"
               "(function(){try{var t=new URLSearchParams(location.search).get('theme');"
               "if(t==='light'||t==='dark')localStorage.setItem('rsk-theme',t)}catch(e){}})()</script>")


def build_variant(src, name, title, css):
    out = src.replace("<head>", "<head>" + THEME_PARAM, 1)
    out = re.sub(r"<title>(.*?)</title>", lambda m: f"<title>{m.group(1)} — palette: {title}</title>", out, count=1)
    if css:
        out = out.replace("</head>", f"<style>\n{css}\n</style>\n</head>", 1)
    return out


def swatch_rows(p):
    keys = [("bg", "page"), ("sunken", "band"), ("card", "card"), ("border", "border"), ("text", "text"), ("muted", "muted"),
            ("accent", "accent"), ("ok", "ok"), ("warn", "warn"), ("danger", "danger"), ("info", "info")]
    cells = []
    for scheme in ("light", "dark"):
        v = p[scheme]
        row = [f'<th scope="row">{scheme}</th>']
        for k, label in keys:
            fg = "#000" if _lum(v[k]) > 0.4 else "#fff"
            row.append(f'<td style="background:{v[k]};color:{fg}" title="{label} {v[k]}"><span>{label}</span><code>{v[k]}</code></td>')
        if "human" in v:
            fg = "#000" if _lum(v["human"]) > 0.4 else "#fff"
            row.append(f'<td style="background:{v["human"]};color:{fg}" title="human"><span>human</span><code>{v["human"]}</code></td>')
        cells.append("<tr>" + "".join(row) + "</tr>")
    return "<table class=\"sw\"><tbody>" + "".join(cells) + "</tbody></table>"


def ratios(p):
    out = []
    for scheme in ("light", "dark"):
        v = p[scheme]
        out.append(f"{scheme}: text {contrast(v['text'], v['bg']):.1f}:1, accent {contrast(v['accent'], v['bg']):.1f}:1, "
                   f"label on accent {contrast(v['on'], v['accent']):.1f}:1, border {contrast(v['strong'], v['bg']):.1f}:1")
    return " · ".join(out)


def viewer(order):
    variants = [("current", "Current", "The page as it is today: cool grey and one cyan, amber for the human.")] + [
        (n, PALETTES[n]["title"], PALETTES[n]["idea"]) for n in order]
    buttons = "".join(f'<button type="button" data-v="{n}" aria-pressed="{"true" if i == 0 else "false"}">{html.escape(t)}</button>' for i, (n, t, _) in enumerate(variants))
    cards = []
    for n in order:
        p = PALETTES[n]
        cards.append(f'<article><h2>{html.escape(p["title"])}</h2><p>{html.escape(p["idea"])}</p>{swatch_rows(p)}'
                     f'<p class="r">{ratios(p)}</p><p><a href="palette-{n}.html">Open {html.escape(p["title"])}</a> · '
                     f'<a href="palette-{n}.html?theme=light">light</a> · <a href="palette-{n}.html?theme=dark">dark</a></p></article>')
    notes = "const NOTES=" + str({n: d for n, _, d in variants}).replace("'", '"') + ";"
    return f"""<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>rsk landing page — palette proposals</title>
<style>
:root{{color-scheme:light dark;font:15px/1.5 system-ui,sans-serif;--bg:#f4f5f7;--fg:#14171a;--mu:#565d66;--ln:#c9ced4;--ac:#005a6e}}
@media (prefers-color-scheme:dark){{:root{{--bg:#0c0e11;--fg:#eceef1;--mu:#a0a7b1;--ln:#33383e;--ac:#5fd0e6}}}}
body{{margin:0;background:var(--bg);color:var(--fg)}}
header.bar{{position:sticky;top:0;z-index:5;background:var(--bg);border-bottom:1px solid var(--ln);padding:12px 20px;display:flex;flex-wrap:wrap;gap:8px 24px;align-items:center}}
h1{{font-size:1rem;margin:0}} .grp{{display:flex;gap:6px;align-items:center}} .grp b{{color:var(--mu);font-weight:600;font-size:.8rem;margin-right:2px}}
button{{font:inherit;padding:6px 12px;border:1px solid var(--ln);background:transparent;color:var(--fg);border-radius:4px;cursor:pointer}}
button[aria-pressed=true]{{background:var(--fg);color:var(--bg);border-color:var(--fg)}} button:focus-visible,a:focus-visible{{outline:2px solid var(--ac);outline-offset:2px}}
#note{{margin:0;padding:10px 20px;color:var(--mu);border-bottom:1px solid var(--ln)}}
#stage{{display:flex;justify-content:center;padding:16px;background:color-mix(in srgb,var(--fg) 6%,var(--bg))}}
iframe{{border:1px solid var(--ln);background:#fff;height:78vh;width:100%;max-width:1440px}} #stage.m iframe{{max-width:390px}}
main{{padding:20px;max-width:1500px;margin:0 auto;display:grid;gap:20px;grid-template-columns:repeat(auto-fit,minmax(520px,1fr))}}
article{{border:1px solid var(--ln);padding:16px;border-radius:6px}} article h2{{margin:0 0 4px;font-size:1.1rem}} article p{{margin:6px 0}}
.r{{font-size:.8rem;color:var(--mu)}} a{{color:var(--ac)}}
.sw{{border-collapse:collapse;width:100%;table-layout:fixed;margin:10px 0}} .sw th{{width:48px;text-align:left;font-size:.75rem;color:var(--mu);font-weight:600}}
.sw td{{padding:8px 4px;font-size:.62rem;text-align:center;vertical-align:top;border:1px solid var(--ln)}} .sw td span{{display:block;font-weight:700}} .sw td code{{font-size:.6rem}}
</style></head><body>
<header class="bar"><h1>Palette proposals</h1>
<div class="grp" role="group" aria-label="Palette"><b>Palette</b>{buttons}</div>
<div class="grp" role="group" aria-label="Theme"><b>Theme</b><button type="button" data-t="light" aria-pressed="true">Light</button><button type="button" data-t="dark" aria-pressed="false">Dark</button></div>
<div class="grp" role="group" aria-label="Width"><b>Width</b><button type="button" data-w="d" aria-pressed="true">Desktop</button><button type="button" data-w="m" aria-pressed="false">Mobile 390</button></div></header>
<p id="note" aria-live="polite"></p>
<div id="stage"><iframe id="f" title="Landing page preview"></iframe></div>
<main>{"".join(cards)}</main>
<script>
{notes}
var st={{v:'current',t:'light',w:'d'}},f=document.getElementById('f'),stage=document.getElementById('stage');
function go(){{f.src='palette-'+st.v+'.html?theme='+st.t;stage.className=st.w==='m'?'m':'';document.getElementById('note').textContent=NOTES[st.v]}}
document.querySelectorAll('[data-v],[data-t],[data-w]').forEach(function(b){{b.addEventListener('click',function(){{
 var k=b.dataset.v?'v':b.dataset.t?'t':'w';st[k]=b.dataset.v||b.dataset.t||b.dataset.w;
 document.querySelectorAll('['+(k==='v'?'data-v':k==='t'?'data-t':'data-w')+']').forEach(function(x){{x.setAttribute('aria-pressed',x===b?'true':'false')}});go()}})}});
go();
</script></body></html>
"""


def main():
    src = SRC.read_text()
    failed = False
    for name, p in PALETTES.items():
        bad = check(p)
        if bad:
            failed = True
            print(f"FAIL {name}")
            for b in bad:
                print("   ", b)
    if failed:
        return 1
    (HERE / "palette-current.html").write_text(build_variant(src, "current", "current", ""))
    for name, p in PALETTES.items():
        (HERE / f"palette-{name}.html").write_text(build_variant(src, name, p["title"], override_css(name, p)))
        print(f"ok   palette-{name}.html  ({ratios(p)})")
    (HERE / "palette-viewer.html").write_text(viewer(list(PALETTES)))
    print("ok   palette-viewer.html")
    return 0


if __name__ == "__main__":
    sys.exit(main())

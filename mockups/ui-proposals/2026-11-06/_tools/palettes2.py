#!/usr/bin/env python3
"""Round 2 of the rsk landing-page palettes: one self-contained page.

Reads the approved combined page from the 2026-11-04 round and writes
../palettes-round-2.html: the same page with every palette as a
`html[data-palette=<name>]` override block and the mockup preview bar inlined
(palette axis; the page's own button flips light/dark, `?theme=` forces it).
The source file is never modified.

    python3 -I _tools/palettes2.py

The WCAG 2.2 gate is the one from round 1 (palettes.check): text, muted, accent
and the four status colours 4.5:1 on page, band and card; the label on an accent
fill 4.5:1; a strong border 3:1. The script exits 1 if any palette fails.
Ember is included unchanged so the new ones are judged against what ships today.
"""
import html
import importlib.util
import json
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent.parent
ROUND1 = HERE.parent / "2026-11-04"
SPEC = importlib.util.spec_from_file_location("round1", ROUND1 / "_tools" / "palettes.py")
r1 = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(r1)
ASSETS = HERE.parent.parent.parent / "skills" / "frontend" / "mockup-architect" / "assets"

PALETTES = {"ember": dict(r1.PALETTES["ember"], title="Ember (today)",
    idea="What the site ships now: warm paper, vermilion. The baseline the rest are judged against.")}

PALETTES["blueprint"] = {
    "title": "Blueprint",
    "idea": "Cobalt on a cool paper, navy at night. Technical and calm; the amber stays for the human, so \"You approve.\" is the warm line on a cold page.",
    "light": dict(bg="#f4f6fb", sunken="#e9edf6", card="#ffffff", border="#d5dcec", strong="#66718f", text="#0f1a33", muted="#4a5573",
                  accent="#1f4fd8", on="#ffffff", tint="rgba(31,79,216,.10)", ok="#17693a", warn="#7a4f00", danger="#a3261f", info="#0b6a7a"),
    "dark": dict(bg="#0a1020", sunken="#0f172c", card="#141d36", border="#26324f", strong="#66739a", text="#e9eefc", muted="#a6b0cf",
                 accent="#7aa2ff", on="#0a1633", tint="rgba(122,162,255,.16)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#6fd3e6"),
}
PALETTES["teal-sand"] = {
    "title": "Teal & Sand",
    "idea": "Deep teal on warm sand, a coral for the human. Earthy without being orange, and the least like a developer-tool default.",
    "light": dict(bg="#f6f1e7", sunken="#ede6d6", card="#fffcf5", border="#ded5bf", strong="#756b55", text="#1a1d1c", muted="#4f564e",
                  accent="#0b6e6a", on="#ffffff", tint="rgba(11,110,106,.10)", ok="#1d6a3b", warn="#7a4f00", danger="#a3261f", info="#1f5e8a", human="#a83e16"),
    "dark": dict(bg="#0d1514", sunken="#121d1b", card="#172523", border="#243735", strong="#5f7d78", text="#e8f1ee", muted="#9fb8b2",
                 accent="#4fd1c0", on="#062320", tint="rgba(79,209,192,.15)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#8cc4ee", human="#ff9b73"),
}
PALETTES["claret"] = {
    "title": "Claret",
    "idea": "Raspberry on a blush-tinted ground, with a teal human. Warm and a little editorial, but clearly not the orange family.",
    "light": dict(bg="#faf5f4", sunken="#f2e8e6", card="#fffdfc", border="#e6d4d1", strong="#826c68", text="#231416", muted="#5b4547",
                  accent="#a4173f", on="#ffffff", tint="rgba(164,23,63,.10)", ok="#1d6a3b", warn="#7a4f00", danger="#a3261f", info="#1f5e8a", human="#0b6560"),
    "dark": dict(bg="#150c0f", sunken="#1d1216", card="#25171c", border="#3a252b", strong="#85666d", text="#f6ebed", muted="#c0a7ac",
                 accent="#ff7d9c", on="#2b0612", tint="rgba(255,125,156,.15)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#8cc4ee", human="#5fd6c8"),
}
PALETTES["graphite-gold"] = {
    "title": "Graphite & Gold",
    "idea": "Neutral graphite with a brass-gold accent. Quiet, the most restrained; teal takes the human so the gold and the warning amber do not blur.",
    "light": dict(bg="#f5f4f1", sunken="#eceae5", card="#ffffff", border="#dcd9d1", strong="#736f66", text="#181714", muted="#54514a",
                  accent="#8a5a00", on="#ffffff", tint="rgba(138,90,0,.10)", ok="#1d6a3b", warn="#7a4f00", danger="#a3261f", info="#1f5e8a", human="#0b6560"),
    "dark": dict(bg="#121211", sunken="#181817", card="#1e1e1c", border="#31302c", strong="#6e6b62", text="#f1efe9", muted="#b0ada3",
                 accent="#e6b93d", on="#1f1600", tint="rgba(230,185,61,.14)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#8cc4ee", human="#5fd6c8"),
}
PALETTES["mono-magenta"] = {
    "title": "Mono & Magenta",
    "idea": "Black and white with one hot magenta. The boldest and the most brand-like; the colour is only the accent, so everything else reads as type.",
    "light": dict(bg="#fafafa", sunken="#f1f1f2", card="#ffffff", border="#e0e0e3", strong="#6e6e75", text="#0e0e10", muted="#505056",
                  accent="#b4147a", on="#ffffff", tint="rgba(180,20,122,.09)", ok="#1d6a3b", warn="#7a4f00", danger="#a3261f", info="#1f5e8a"),
    "dark": dict(bg="#0b0b0d", sunken="#111114", card="#17171b", border="#2a2a30", strong="#6c6c76", text="#f4f4f6", muted="#a8a8b2",
                 accent="#ff6ec4", on="#2a0a1d", tint="rgba(255,110,196,.14)", ok="#7fd49a", warn="#f2c94c", danger="#ff918c", info="#8cc4ee"),
}
PALETTES["emerald"] = {
    "title": "Emerald",
    "idea": "Emerald on a cool white, a deep green night. Fresh and legible; it reads as \"go\", which suits an installer, but the green-for-success meaning is shared with the status colour.",
    "light": dict(bg="#f4f7f5", sunken="#e8eeea", card="#ffffff", border="#d3ddd7", strong="#62746b", text="#0f1a14", muted="#44544b",
                  accent="#0b7a4b", on="#ffffff", tint="rgba(11,122,75,.10)", ok="#17693a", warn="#7a4f00", danger="#a3261f", info="#1f5e8a"),
    "dark": dict(bg="#08130e", sunken="#0d1b14", card="#12251b", border="#1f3a2c", strong="#5b8a72", text="#e6f3ec", muted="#9fc2b0",
                 accent="#3ddc97", on="#052618", tint="rgba(61,220,151,.14)", ok="#8fe0a8", warn="#f2c94c", danger="#ff918c", info="#8cc4ee"),
}


def override_css(name, p):
    pre = f"html[data-palette={name}]"
    L, D = p["light"], p["dark"]
    term = dict(D, bg=D["sunken"], sunken=D["bg"])
    css = [
        r1.scheme_block(f"{pre}[data-theme=light]", L, False),
        r1.scheme_block(f"{pre}[data-theme=dark]", D, True),
        r1.scheme_block(f"{pre} .night", D, True),
        r1.scheme_block(f"{pre} .tn", term, True),
    ]
    for scheme in ("light", "dark"):
        if "human" in p[scheme]:
            s, h = f"{pre}[data-theme={scheme}]", p[scheme]["human"]
            css.append(f"{s} .c-h1 .you,{s} .c-lane.y,{s} .c-arrow.ask,{s} .c-you .lbl,{s} .k-badge.hu{{color:{h}}}\n"
                       f"{s} .c-lane.y,{s} .c-you{{border-color:{h}}}\n"
                       f"{s} .c-you .b.al{{border-color:{h};background:{h};color:var(--bg)}}")
    return f"/* palette: {name} */\n" + "\n".join(css)


def build():
    src = (ROUND1 / "direction-4-combined.html").read_text()
    css = "\n".join(override_css(n, p) for n, p in PALETTES.items())
    bridge = (".pv{--card:var(--surface);--line-strong:var(--line-2);--btn-bg:var(--text);--btn-fg:var(--bg);"
              "--shadow:0 8px 30px rgba(0,0,0,.25)}")
    axes = [{"key": "palette", "label": "Palette", "values": [
        {"v": n, "label": p["title"], "sw": p["light"]["accent"], "swDark": p["dark"]["accent"]} for n, p in PALETTES.items()]}]
    cfg = f'<script type="application/json" id="pv-config">{json.dumps({"axes": axes})}</script>'
    bar_css = (ASSETS / "pv-bar.css").read_text()
    bar_js = (ASSETS / "pv-bar.js").read_text()
    theme_param = ("<script>(function(){try{var t=new URLSearchParams(location.search).get('theme');"
                   "if(t==='light'||t==='dark')localStorage.setItem('rsk-theme',t)}catch(e){}})()</script>")
    head = (f"{theme_param}<style>\n{css}\n{bridge}\n{bar_css}\n</style>\n{cfg}\n<script>\n{bar_js}\n</script>\n")
    # the source page says 58; the catalogue is 60 now (honest data)
    src = src.replace("58 skills", "60 skills").replace('data-count="58">58', 'data-count="60">60')
    out = src.replace("<head>", "<head>" + head, 1)
    out = out.replace("</title>", " — palette round 2</title>", 1) if "</title>" in out else out
    # keep the last section clear of the bar
    out = out.replace("</body>", "<style>body{padding-bottom:132px}@media(max-width:900px){body{padding-bottom:260px}}</style>\n</body>", 1)
    return out


def main():
    failed = False
    for name, p in PALETTES.items():
        bad = r1.check(p)
        if bad:
            failed = True
            print(f"FAIL {name}")
            for b in bad:
                print("   ", b)
        else:
            print(f"ok   {name:14} {r1.ratios(p)}")
    if failed:
        return 1
    (HERE / "palettes-round-2.html").write_text(build())
    print("wrote palettes-round-2.html")
    return 0


if __name__ == "__main__":
    sys.exit(main())

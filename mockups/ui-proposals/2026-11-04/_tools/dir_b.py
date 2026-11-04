"""Direction B — the catalogue as a real number.

Device: 58, one square per published skill, grouped by its skills/<domain>/
folder; the same count drives the page (58 → 15 bundles → 139 releases →
13 commands + 3 tools → 1 dependency). The explorer filters the way
`rsk catalog` does: by bundle, personal hidden unless --include personal.
"""

from __future__ import annotations

import data as D
from family import copy_block, esc, film_ctl, footer, header, page

CSS = r"""
.b-num{font:800 clamp(8rem,14vw,12.5rem)/.82 var(--mono);letter-spacing:-.07em;color:var(--accent-t);font-variant-numeric:tabular-nums;display:block}
.b-sec-n{font:800 clamp(3.6rem,9vw,7.5rem)/.85 var(--mono);letter-spacing:-.06em;color:var(--accent-t);display:block;margin-bottom:.12em;font-variant-numeric:tabular-nums}
.b-hero{padding:clamp(28px,5vw,72px) 0 clamp(56px,8vw,112px)}
.b-grid{display:grid;gap:32px clamp(32px,5vw,80px);grid-template-columns:minmax(0,1fr)}
@media (min-width:1080px){.b-grid{grid-template-columns:minmax(0,5fr) minmax(0,6fr);align-items:start}}
.b-h1{font-size:clamp(1.8rem,3.3vw,2.75rem);font-weight:750;letter-spacing:-.03em;line-height:1.08}
.b-h1 .b-num{margin-bottom:.1em}
.b-h1 .w{display:block;max-width:19ch}
.b-lede{color:var(--muted);font-size:clamp(1.05rem,1.4vw,1.2rem);max-width:46ch;margin:22px 0 26px}
.b-cta{display:flex;flex-wrap:wrap;gap:12px}
/* unit chart */
.b-units{margin:0}
.b-rows{list-style:none;margin:0;padding:0;border-top:1px solid var(--line)}
.b-rows li{display:grid;grid-template-columns:minmax(0,8.5em) minmax(0,1fr) 2.2em;align-items:center;gap:12px;min-height:30px;border-bottom:1px solid color-mix(in srgb,var(--line) 55%,transparent);padding:3px 0}
.b-rows .dn{font:600 .85rem var(--mono);color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.b-rows .ct{font:700 .9rem var(--mono);text-align:right;color:var(--muted);font-variant-numeric:tabular-nums}
.b-sq{display:flex;flex-wrap:wrap;gap:4px}
.b-sq i{display:block;width:18px;height:18px;border-radius:3px;background:var(--accent)}
@media (max-width:600px){.b-sq{gap:3px}.b-sq i{width:14px;height:14px}.b-rows li{grid-template-columns:minmax(0,7.2em) minmax(0,1fr) 1.6em;gap:8px}}
.b-rows li.p .b-sq i{background:transparent;border:2px dashed var(--muted)}
.anim .f-run .b-rows li.f-h{opacity:1;visibility:visible;transform:none}
.anim .f-run .b-rows li.f-h .b-sq i{transform:scale(0)}
.anim .b-sq i{transition:transform .38s cubic-bezier(.3,1.5,.5,1)}
.b-cap{display:flex;flex-wrap:wrap;gap:10px 18px;justify-content:space-between;align-items:center;margin-top:14px;font-size:.86rem;color:var(--muted)}
.b-cap p{margin:0;max-width:60ch}
/* explorer */
.b-ex-top{display:grid;gap:24px;grid-template-columns:minmax(0,1fr);align-items:end}
@media (min-width:980px){.b-ex-top{grid-template-columns:minmax(0,1fr) auto}}
.b-count{font:800 clamp(3rem,6vw,5.2rem)/.85 var(--mono);letter-spacing:-.06em;color:var(--accent-t);font-variant-numeric:tabular-nums;min-width:2ch;display:inline-block;text-align:right}
.b-count-l{font:600 .95rem var(--mono);color:var(--muted)}
.b-chips{display:flex;flex-wrap:wrap;gap:8px;margin:28px 0 12px;padding:0;list-style:none}
.b-chip{min-height:44px;min-width:44px;padding:0 14px;border-radius:999px;border:1px solid var(--line);background:var(--surface);color:var(--text);font:600 .85rem var(--mono);cursor:pointer}
.b-chip:hover{border-color:var(--accent-t)}
.b-chip[aria-pressed=true]{background:var(--accent-t);border-color:var(--accent-t);color:var(--on-accent)}
.b-chip.tg[aria-pressed=true]{background:var(--warn);border-color:var(--warn);color:var(--bg)}
@media (max-width:700px){.b-chips{flex-wrap:nowrap;overflow-x:auto;padding-bottom:6px;margin-inline:calc(-1*clamp(20px,4vw,56px));padding-inline:clamp(20px,4vw,56px);scrollbar-width:thin}.b-chip{flex:none}}
.b-note{min-height:3.2em;color:var(--muted);font-size:.95rem;margin:0 0 22px;max-width:80ch}
.b-note code{color:var(--text)}
.b-field{columns:4 15rem;column-gap:28px}
.b-dom{break-inside:avoid;margin:0 0 22px;padding-top:10px;border-top:2px solid var(--text)}
.b-dom h3{font:700 .8rem var(--mono);letter-spacing:.06em;text-transform:uppercase;display:flex;justify-content:space-between;margin-bottom:8px}
.b-dom ul{list-style:none;margin:0;padding:0;display:flex;flex-wrap:wrap;gap:6px}
.b-tag{font:500 .82rem/1 var(--mono);padding:7px 9px;border:1px solid var(--line);border-radius:6px;background:var(--surface)}
.b-field.f .b-tag{border-style:dashed;color:var(--muted);background:transparent}
.b-field.f .b-tag.on{border:1px solid var(--accent-t);border-left-width:4px;color:var(--text);background:var(--sel)}
.b-tag .v{color:var(--muted);margin-left:6px;font-size:.92em}
.b-field .ph{font-size:.85rem;color:var(--muted);font-style:italic}
.b-field:not(.inc) .b-tag.p{display:none}
.b-field.inc .ph{display:none}
/* bundles list */
.b-bl{list-style:none;margin:0;padding:0;display:grid;grid-template-columns:minmax(0,1fr);border-top:2px solid var(--text)}
@media (min-width:900px){.b-bl{grid-template-columns:repeat(3,minmax(0,1fr));column-gap:32px}}
.b-bl li{padding:14px 0;border-bottom:1px solid var(--line);display:grid;grid-template-columns:auto minmax(0,1fr);gap:2px 14px}
.b-bl .n{font:800 1.6rem/1 var(--mono);color:var(--accent-t);grid-row:span 2;min-width:1.4ch;text-align:right}
.b-bl b{font-family:var(--mono)}
.b-bl span{color:var(--muted);font-size:.92rem}
/* releases */
.b-rel{display:grid;gap:32px;grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.b-rel{grid-template-columns:minmax(0,.8fr) minmax(0,1.4fr);align-items:end}}
.b-chart svg{display:block;width:100%;height:auto}
.b-chart .bar{fill:var(--accent)}
.b-chart .ax{stroke:var(--text);stroke-width:1.5}
.b-chart text{fill:var(--muted);font:500 12px var(--mono)}
.b-chart text.pk{fill:var(--text);font-weight:700}
.b-chart details{margin-top:10px;font-size:.9rem}
.b-chart summary{min-height:44px;display:flex;align-items:center;cursor:pointer;color:var(--accent-t);font-weight:600}
.b-chart table{border-collapse:collapse;font:500 .82rem var(--mono)}
.b-chart td,.b-chart th{padding:3px 14px 3px 0;text-align:left;border-bottom:1px solid var(--line)}
/* cli */
.b-cli{display:grid;gap:48px;grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.b-cli{grid-template-columns:1fr 1fr}}
.b-cmds{display:flex;flex-wrap:wrap;gap:8px;list-style:none;margin:0;padding:0}
.b-cmds code{display:inline-block;padding:6px 10px;border:1px solid var(--line);border-radius:6px;background:var(--surface);font-size:.92rem}
.b-tools{list-style:none;margin:0;padding:0;border-top:2px solid var(--text)}
.b-tools li{padding:14px 0;border-bottom:1px solid var(--line)}
.b-tools code{font-weight:700;font-size:1rem;margin-right:10px}
.b-tools p{margin:6px 0 0;color:var(--muted);font-size:.95rem}
/* ones */
.b-ones{list-style:none;margin:0;padding:0;display:grid;grid-template-columns:minmax(0,1fr);border-top:2px solid var(--text)}
@media (min-width:900px){.b-ones{grid-template-columns:repeat(3,minmax(0,1fr));column-gap:40px}}
.b-ones li{padding:22px 0;border-bottom:1px solid var(--line)}
.b-ones .b-sec-n{font-size:clamp(3rem,6vw,5rem)}
.b-ones li.take .b-sec-n{color:var(--warn)}
.b-ones h3{font-size:1.25rem;margin:.2em 0 .4em}
.b-ones p{color:var(--muted);margin:0}
.b-start{display:grid;gap:32px;grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.b-start{grid-template-columns:minmax(0,.8fr) minmax(0,1.2fr)}}
"""


def _bundles_of(d: dict) -> dict[str, list[str]]:
    m: dict[str, list[str]] = {}
    for b in d["bundles"]:
        for n in b["in_catalogue"]:
            m.setdefault(n, []).append(b["name"])
    return m


def hero(d: dict) -> str:
    doms = sorted(d["domains"].items(), key=lambda kv: (kv[0] == "personal", -len(kv[1]), kv[0]))
    rows = "".join(
        f'<li data-f="r{i}"{" class=p" if dom == "personal" else ""}><span class="dn">{esc(dom)}</span>'
        f'<span class="b-sq" aria-hidden="true">{"<i></i>" * len(sk)}</span><span class="ct">{len(sk)}</span></li>'
        for i, (dom, sk) in enumerate(doms)
    )
    n = len(d["skills"])
    return f"""
<section class="b-hero" id="top" aria-labelledby="b-h1"><div class="wrap b-grid">
<div>
<p class="k-kick"><span class="gl" aria-hidden="true">◆</span> rsk · the catalogue</p>
<h1 class="b-h1" id="b-h1"><span class="b-num" aria-hidden="true" id="b-count0">{n}</span><span class="sr">{n} </span><span class="w">skills your coding agent can load, installed and versioned by one CLI.</span></h1>
<p class="b-lede">Each one is a <code>SKILL.md</code> with the standards for one job: <code>go-architect</code>, <code>sql-architect</code>, <code>mcp-architect</code>. <code>rsk</code> installs them, alone or in bundles, into Claude Code and OpenCode.</p>
<div class="b-cta"><a class="k-btn primary" href="#start">Install rsk</a><a class="k-btn" href="#catalogue">Filter the {n}</a></div>
</div>
<figure class="b-units" id="b-film" aria-labelledby="b-units-cap">
<ul class="b-rows" aria-label="Skills per domain">{rows}</ul>
<figcaption class="b-cap" id="b-units-cap"><p><span class="k-badge ro">real</span> One square per published skill, by the <code>skills/&lt;domain&gt;/</code> folder it lives in. Dashed: the {d["personal"]} marked personal, which <code>rsk catalog</code> hides unless you pass <code>--include personal</code>.</p>{film_ctl()}</figcaption>
</figure>
</div></section>"""


def explorer(d: dict) -> str:
    bo = _bundles_of(d)
    doms = sorted(d["domains"].items(), key=lambda kv: (kv[0] == "personal", -len(kv[1]), kv[0]))
    blocks = []
    for dom, sk in doms:
        tags = "".join(
            f'<li class="b-tag{" p" if s["personal"] else ""}" data-b="{esc(" ".join(bo.get(s["name"], [])))}">{esc(s["name"])}<span class="v">{esc(s["latest"])}</span></li>'
            for s in sk
        )
        ph = f'<li class="ph">{len(sk)} personal skills, hidden by default</li>' if dom == "personal" else ""
        blocks.append(f'<section class="b-dom" aria-label="{esc(dom)}"><h3><span>{esc(dom)}</span><span>{len(sk)}</span></h3><ul>{tags}{ph}</ul></section>')
    chips = '<li><button class="b-chip" type="button" aria-pressed="true" data-bundle="">all</button></li>' + "".join(
        f'<li><button class="b-chip" type="button" aria-pressed="false" data-bundle="{esc(b["name"])}">{esc(b["name"])}</button></li>'
        for b in d["bundles"]
    ) + '<li><button class="b-chip tg" type="button" aria-pressed="false" id="b-inc">--include personal</button></li>'
    notes = {}
    for b in d["bundles"]:
        t = f'{b["description"]}. <code>rsk install {esc(b["name"])}</code>'
        if b["official"]:
            t += f' — includes {len(b["official"])} official Anthropic skill{"s" if len(b["official"]) > 1 else ""} ({", ".join(b["official"])}), fetched from Anthropic and not part of the {len(d["skills"])}.'
        if b["missing"]:
            t += f' — also names {" and ".join(b["missing"])}, which are not in the catalogue (flagged in the notes).'
        notes[b["name"]] = t
    import json

    pub = len(d["skills"]) - d["personal"]
    return f"""
<section class="k-sec" id="catalogue" aria-labelledby="b-ex-h" style="background:var(--bg2)"><div class="wrap">
<div class="b-ex-top">
<div><span class="b-sec-n" aria-hidden="true">{len(d["bundles"])}</span><h2 id="b-ex-h" style="margin-top:0"><span class="sr">{len(d["bundles"])} </span>bundles. Filter the catalogue the way <code>rsk catalog</code> does.</h2></div>
<p style="margin:0;text-align:right"><span class="b-count" id="b-count" aria-hidden="true">{pub}</span><br><span class="b-count-l" id="b-count-l" aria-hidden="true">shown of {len(d["skills"])}</span></p>
</div>
<ul class="b-chips" aria-label="Filter by bundle" hidden id="b-chips">{chips}</ul>
<p class="b-note" id="b-note" aria-live="polite">Showing the {pub} skills <code>rsk catalog</code> lists by default; the {d["personal"]} personal ones need <code>--include personal</code>.</p>
<div class="b-field" id="b-field">{"".join(blocks)}</div>
<script type="application/json" id="b-notes">{json.dumps(notes).replace("</", "<\\/")}</script>
</div></section>"""


def bundles(d: dict) -> str:
    items = "".join(
        f'<li><span class="n">{len(b["skills"])}</span><b>{esc(b["name"])}</b><span>{esc(b["description"])}</span></li>'
        for b in d["bundles"]
    )
    return f"""
<section class="k-sec" id="bundles" aria-labelledby="b-bl-h"><div class="wrap">
<p class="k-kick rv">From internal/config/catalog.toml</p>
<h2 id="b-bl-h" class="rv">One name, a whole stack: <code>rsk install go-grpc</code> brings five.</h2>
<ul class="b-bl rv" aria-label="Bundles and how many skills each names">{items}</ul>
<p style="margin-top:18px;color:var(--muted);font-size:.95rem">Your own <code>~/.config/rsk/catalog.toml</code> replaces a bundle with the same name or adds new ones.</p>
</div></section>"""


def releases(d: dict) -> str:
    hist = d["hist"]
    W, H, pad, base = 760, 240, 30, 200
    bw = (W - pad * 2) / len(hist)
    mx = max(c for _, c in hist)
    bars, labels = [], []
    last_m = None
    for i, (wk, c) in enumerate(hist):
        h = (c / mx) * (base - 30)
        x = pad + i * bw
        bars.append(f'<rect class="bar" x="{x + 2:.1f}" y="{base - h:.1f}" width="{bw - 4:.1f}" height="{h:.1f}" rx="2"/>')
        if c == mx:
            labels.append(f'<text class="pk" x="{x + bw / 2:.1f}" y="{base - h - 8:.1f}" text-anchor="middle">{c}</text>')
        if wk.month != last_m:
            labels.append(f'<text x="{x + 2:.1f}" y="{base + 22}">{wk.strftime("%b")}</text>')
            last_m = wk.month
    svg = (
        f'<svg viewBox="0 0 {W} {H}" aria-hidden="true" focusable="false">{"".join(bars)}'
        f'<line class="ax" x1="{pad}" x2="{W - pad}" y1="{base}" y2="{base}"/>{"".join(labels)}</svg>'
    )
    rows = "".join(f"<tr><td>{wk.isoformat()}</td><td>{c}</td></tr>" for wk, c in hist)
    peak_wk = next(wk for wk, c in hist if c == mx)
    return f"""
<section class="k-sec night" id="releases" aria-labelledby="b-rel-h"><div class="wrap b-rel">
<div class="rv"><span class="b-sec-n" aria-hidden="true">{d["releases"]}</span><h2 id="b-rel-h" style="margin-top:0"><span class="sr">{d["releases"]} </span>skill releases since {d["first"].strftime("%-d %B %Y")}.</h2>
<p class="k-lede">Each skill is versioned on its own and published as <code>&lt;skill&gt;@v&lt;version&gt;</code> when its folder changes on <code>main</code>. The lock is what keeps your project on the version you chose.</p></div>
<figure class="b-chart rv" style="margin:0" aria-labelledby="b-chart-cap">
{svg}
<figcaption id="b-chart-cap" style="font-size:.9rem;color:var(--muted)"><span class="k-badge ro">real</span> Releases per week from <code>published_at</code> in <code>registry/index.json</code>. Peak: {mx} in the week of {peak_wk.strftime("%-d %b")}, the first publish.</figcaption>
<details><summary>Releases per week, as a table</summary><table><thead><tr><th scope="col">Week of</th><th scope="col">Releases</th></tr></thead><tbody>{rows}</tbody></table></details>
</figure>
</div></section>"""


def cli() -> str:
    cmds = "".join(f"<li><code>{esc(c)}</code></li>" for _, cs in D.HELP_GROUPS for c, _ in cs)
    tools = "".join(
        f'<li><code>{esc(n)}</code><span class="k-badge {k}">{"read-only" if k == "ro" else "writes · asks you"}</span><p>{esc(t)}</p></li>'
        for n, t, k in D.TOOLS
    )
    return f"""
<section class="k-sec" id="cli" aria-labelledby="b-cli-h"><div class="wrap">
<h2 id="b-cli-h" class="sr">Commands and MCP tools</h2>
<div class="b-cli">
<div class="rv"><span class="b-sec-n" aria-hidden="true">13</span><h3 style="font-size:clamp(1.5rem,2.6vw,2.1rem);margin-bottom:.6em"><span class="sr">13 </span>commands for people.</h3><ul class="b-cmds" aria-label="rsk subcommands">{cmds}</ul></div>
<div class="rv"><span class="b-sec-n" aria-hidden="true">3</span><h3 style="font-size:clamp(1.5rem,2.6vw,2.1rem);margin-bottom:.6em"><span class="sr">3 </span>MCP tools for agents, one that writes.</h3><ul class="b-tools">{tools}</ul>
<p style="margin-top:14px;color:var(--muted);font-size:.95rem">Plus the <code>rsk://catalog</code> resource: every skill, its latest version, descriptions cut to 200 characters.</p></div>
</div></div></section>"""


def ones() -> str:
    return """
<section class="k-sec" id="trade" aria-labelledby="b-one-h" style="background:var(--bg2)"><div class="wrap">
<h2 id="b-one-h" class="rv">The trade-off, in ones.</h2>
<ul class="b-ones rv">
<li><span class="b-sec-n" aria-hidden="true">1</span><h3><span class="sr">1 </span>binary</h3><p>The same <code>rsk</code> for every project on the machine.</p></li>
<li><span class="b-sec-n" aria-hidden="true">1</span><h3><span class="sr">1 </span>config file</h3><p><code>~/.config/rsk/config.json</code>, written once by <code>rsk init</code>.</p></li>
<li class="take"><span class="b-sec-n" aria-hidden="true">1</span><h3><span class="sr">1 </span>dependency you accept</h3><p>The registry at skills.ralvarez.dev must be reachable to install or update, or you point rsk at a local clone of the repository.</p></li>
</ul></div></section>"""


def start() -> str:
    return f"""
<section class="k-sec" id="start" aria-labelledby="b-start-h"><div class="wrap b-start">
<div class="rv"><span class="b-sec-n" aria-hidden="true">5</span><h2 id="b-start-h" style="margin-top:0"><span class="sr">5 </span>commands from nothing to a pinned skill.</h2></div>
<div class="rv">{copy_block("b-qs", D.QUICK_START, "quick-start commands")}</div>
</div></section>"""


JS = r"""
(function(){
/* hero film: squares fill domain by domain while the count climbs to 58 */
var root=document.getElementById('b-film'),num=document.getElementById('b-count0');
if(root&&window.kFilm){var rows=[].slice.call(root.querySelectorAll('[data-f]')),total=+num.textContent,steps=[],acc=0,t0=300,per=3600/total;
 rows.forEach(function(r,i){var n=r.querySelectorAll('i').length,at=t0+acc*per;steps.push({at:at,fn:function(){r.classList.remove('f-h')}});
  [].forEach.call(r.querySelectorAll('i'),function(sq,j){sq.style.transitionDelay=(j*per/1000).toFixed(3)+'s'});acc+=n});
 steps.push({at:t0,dur:3600+300,fn:function(p){num.textContent=Math.round(total*p)}});
 kFilm(root,{total:t0+3600+900,steps:steps,reset:function(){kF.hide(root);num.textContent='0'},end:function(){kF.end(root);num.textContent=total;root.querySelectorAll('i').forEach(function(i){i.style.transitionDelay=''})}});}
/* explorer */
var chips=document.getElementById('b-chips');if(!chips)return;chips.hidden=false;
var field=document.getElementById('b-field'),note=document.getElementById('b-note'),cnt=document.getElementById('b-count'),cntl=document.getElementById('b-count-l'),inc=document.getElementById('b-inc'),notes=JSON.parse(document.getElementById('b-notes').textContent);
var tags=[].slice.call(field.querySelectorAll('.b-tag')),cur='',all=tags.length;
function render(){var incl=inc.getAttribute('aria-pressed')==='true';field.classList.toggle('inc',incl);field.classList.toggle('f',!!cur);var n=0;
 tags.forEach(function(t){var vis=incl||!t.classList.contains('p'),on=!cur||(' '+t.dataset.b+' ').indexOf(' '+cur+' ')>=0;t.classList.toggle('on',!!cur&&on);if(vis&&on)n++});
 cnt.textContent=n;cntl.textContent=(cur?'in '+cur+', ':'shown ')+'of '+all;
 note.innerHTML=cur?notes[cur]:(incl?'Showing all '+all+' skills, personal ones included.':'Showing the '+n+' skills <code>rsk catalog</code> lists by default; the personal ones need <code>--include personal</code>.');
 var l=document.getElementById('k-live');if(l)l.textContent=n+' of '+all+' skills match'+(cur?' bundle '+cur:'');}
chips.addEventListener('click',function(e){var b=e.target.closest('button');if(!b)return;
 if(b===inc){inc.setAttribute('aria-pressed',inc.getAttribute('aria-pressed')==='true'?'false':'true')}
 else{cur=b.dataset.bundle;chips.querySelectorAll('[data-bundle]').forEach(function(c){c.setAttribute('aria-pressed',c===b?'true':'false')})}
 render()});
})();
"""


def build(d: dict) -> str:
    nav = [("Catalogue", "#catalogue", True), ("Releases", "#releases", True), ("Trade-off", "#trade", True), ("Install", "#start", False)]
    body = (
        header(nav, d["version"])
        + '<main id="main">'
        + hero(d)
        + explorer(d)
        + bundles(d)
        + releases(d)
        + cli()
        + ones()
        + start()
        + "</main>"
        + footer(D.SOURCES, d["version"])
    )
    return page(
        "rsk — 58 skills for your coding agent (direction B)",
        "58 skills for Claude Code and OpenCode in 17 domains and 15 bundles, installed and versioned by the rsk CLI.",
        CSS,
        body,
        JS,
    )

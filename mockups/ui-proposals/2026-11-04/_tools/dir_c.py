"""Direction C — the agent loop.

Device: the exchange as the three calls it really is — project_profile
(read-only) → the client asks you (the one human step) → install_skills
(writes). Swimlanes You · Agent · rsk mcp at desktop; a thread on phones.
"""

from __future__ import annotations

import json

import data as D
from family import copy_block, esc, film_ctl, footer, header, page

NAMES = [c["skill"] for c in D.PROFILE["candidates"]]

CSS = r"""
.c-hero{padding:clamp(24px,3.5vw,48px) 0 clamp(56px,8vw,112px)}
.c-top{display:grid;gap:18px 56px;grid-template-columns:minmax(0,1fr);align-items:end;margin-bottom:clamp(24px,3vw,36px)}
@media (min-width:1000px){.c-top{grid-template-columns:minmax(0,7fr) minmax(0,5fr)}}
.c-h1{font-size:clamp(2.5rem,4.6vw,4rem);font-weight:800;letter-spacing:-.045em;line-height:.98}
.c-h1 span{display:block}
.c-h1 .you{color:var(--warn)}
[data-theme=light] .c-h1 .you{color:#8c6c3e}
.c-h1 .rsk{color:var(--accent-t)}
.c-lede{color:var(--muted);font-size:clamp(1.05rem,1.4vw,1.2rem);margin:0 0 20px}
.c-lede code{color:var(--text)}
.c-cta{display:flex;gap:12px;flex-wrap:wrap}
/* sequence */
.c-fig{margin:0}
.c-seq{list-style:none;margin:0;padding:0;display:grid;gap:10px;position:relative}
.c-lanes{display:none}
.c-ex{position:relative;display:grid;gap:8px;grid-template-columns:minmax(0,1fr)}
.c-no{font:800 1rem/1 var(--mono);display:inline-flex;align-items:center;justify-content:center;width:1.9em;height:1.9em;border-radius:50%;border:2px solid currentColor;color:var(--accent-t);margin-right:10px;flex:none}
.c-ex.h .c-no{color:var(--warn)}
.c-nod{display:none!important}
.c-route{display:flex;align-items:center;font:600 .82rem var(--mono);color:var(--muted);margin:0}
.c-route b{color:var(--text);margin:0 .4em}
.c-card{background:var(--surface);border:1px solid var(--line);border-radius:10px;padding:10px 14px;font:500 .84rem/1.6 var(--mono);box-shadow:var(--shadow)}
.c-card .m{font-weight:700;color:var(--accent-t);margin-right:8px}
.c-card dl{margin:4px 0 0;display:grid;grid-template-columns:auto minmax(0,1fr);gap:2px 12px}
.c-card dt{color:var(--muted)}.c-card dd{margin:0;overflow-wrap:anywhere}
.c-card .pf{color:var(--muted)}
.c-card.back{border-style:dashed;box-shadow:none}
.c-arrow{display:none}
.c-you{border:2px solid var(--warn);border-radius:12px;background:var(--surface);padding:12px 14px;font-size:.92rem;box-shadow:var(--shadow)}
.c-you .q{font-weight:700;margin:0 0 6px}
.c-you .q code{color:var(--accent-t)}
.c-you .nm{font:500 .82rem var(--mono);color:var(--muted);margin:0 0 10px}
.c-you .acts{display:flex;gap:8px;align-items:center;flex-wrap:wrap}
.c-you .b{display:inline-flex;align-items:center;min-height:34px;padding:0 12px;border-radius:8px;border:1px solid var(--line);font:600 .84rem var(--sans);color:var(--muted)}
.c-you .b.al{border-color:var(--warn);background:var(--warn);color:var(--bg)}
.c-you .b.al .p{display:none}
.c-you.pend .b.al{background:transparent;color:var(--text)}
.c-you.pend .b.al .p{display:inline}.c-you.pend .b.al .d{display:none}
.c-you .lbl{font:600 .72rem var(--mono);color:var(--warn);text-transform:uppercase;letter-spacing:.06em;margin-left:auto}
@media (min-width:1000px){
 .c-seq{grid-template-columns:repeat(3,minmax(0,1fr));column-gap:0;row-gap:18px;
  background-image:repeating-linear-gradient(to bottom,var(--line) 0 6px,transparent 6px 12px),repeating-linear-gradient(to bottom,var(--line) 0 6px,transparent 6px 12px),repeating-linear-gradient(to bottom,var(--line) 0 6px,transparent 6px 12px);
  background-size:2px 100%;background-repeat:no-repeat;background-position:16.6667% 0,50% 0,83.3333% 0}
 .c-lanes{display:contents}
 .c-lane{text-align:center;font:700 .8rem var(--mono);letter-spacing:.06em;text-transform:uppercase;padding:8px 0;background:var(--bg);position:relative;z-index:1;justify-self:center;min-width:12em;border:1px solid var(--line);border-radius:999px}
 .c-lane.y{color:var(--warn);border-color:var(--warn)}
 .c-ex{grid-column:1/-1;grid-template-columns:repeat(3,minmax(0,1fr));column-gap:0;row-gap:6px}
 .c-ex>.c-route{position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%);white-space:nowrap}
 .c-nod{display:inline-flex!important;position:absolute;left:-3.1em;top:4px;margin:0;background:var(--bg)}
 .c-m{grid-column:2/4;display:grid;justify-items:center;position:relative}
 .c-m.l{grid-column:1/3}
 .c-arrow{display:block;width:50%;height:22px;position:relative;font:700 .8rem var(--mono);color:var(--accent-t)}
 .c-arrow::before{content:"";position:absolute;left:0;right:0;top:50%;border-top:2px solid currentColor}
 .c-arrow.r::after{content:"";position:absolute;right:-1px;top:calc(50% - 6px);border-style:solid;border-width:6px 0 6px 10px;border-color:transparent transparent transparent currentColor}
 .c-arrow.lt::after{content:"";position:absolute;left:-1px;top:calc(50% - 6px);border-style:solid;border-width:6px 10px 6px 0;border-color:transparent currentColor transparent transparent}
 .c-arrow.back{color:var(--muted)}.c-arrow.back::before{border-top-style:dashed}
 .c-arrow.ask{color:var(--warn)}
 .c-m .c-card{width:min(100%,58ch);margin-top:2px}
 .c-m.req .c-card{order:-1;background:none;border:0;box-shadow:none;padding:0 0 2px;text-align:center;margin:0}
 .c-seq{row-gap:12px}
 .c-ex{row-gap:2px}
 .c-card dl.sg{grid-template-columns:auto auto;column-gap:10px}
 .c-youcell{grid-column:1/2;padding:0 14px}
}
.anim .f-run .c-ex .f-h .c-arrow,.anim .f-run .c-m.f-h .c-arrow{clip-path:inset(0 100% 0 0)}
.anim .f-run .c-m.l.f-h .c-arrow{clip-path:inset(0 0 0 100%)}
.anim .c-arrow{transition:clip-path .55s cubic-bezier(.3,.7,.2,1);clip-path:inset(0 0 0 0)}
.c-cap{display:flex;flex-wrap:wrap;gap:10px 18px;align-items:center;justify-content:space-between;margin-top:22px;font-size:.86rem;color:var(--muted)}
.c-cap p{margin:0;max-width:78ch}
/* tools */
.c-tools{list-style:none;margin:0;padding:0;border-top:2px solid var(--text)}
.c-tools>li{display:grid;grid-template-columns:minmax(0,1fr);gap:8px 32px;padding:22px 0;border-bottom:1px solid var(--line)}
@media (min-width:900px){.c-tools>li{grid-template-columns:minmax(0,15em) minmax(0,1fr) minmax(0,17em)}}
.c-tools h3{font:700 1.15rem var(--mono)}
.c-tools p{margin:0;color:var(--muted)}
.c-tools pre{font-size:.78rem;color:var(--muted)}
.c-tools li.w h3{color:var(--danger)}
/* signals */
.c-two{display:grid;gap:40px;grid-template-columns:minmax(0,1fr)}
@media (min-width:1000px){.c-two{grid-template-columns:minmax(0,1.15fr) minmax(0,1fr);align-items:start}}
.c-sig{width:100%;border-collapse:collapse;font-size:.92rem}
.c-sig th{text-align:left;font:600 .76rem var(--mono);color:var(--muted);text-transform:uppercase;letter-spacing:.05em;padding:8px 10px 8px 0;border-bottom:2px solid var(--text)}
.c-sig td{padding:9px 10px 9px 0;border-bottom:1px solid var(--line);vertical-align:top}
.c-sig td:first-child{font:600 .85rem var(--mono);white-space:nowrap}
.c-sig td:last-child{font:500 .84rem var(--mono);color:var(--accent-t)}
.c-sig tr.hit td:first-child::before{content:"● ";color:var(--ok)}
@media (max-width:640px){.c-sig thead{position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%)}.c-sig tr{display:grid;grid-template-columns:1fr;padding:10px 0;border-bottom:1px solid var(--line)}.c-sig td{border:0;padding:1px 0}}
.c-json pre{background:var(--surface);border:1px solid var(--line);border-radius:12px;padding:18px 20px;font-size:.8rem}
.c-told{display:grid;gap:40px;grid-template-columns:minmax(0,1fr)}
@media (min-width:1000px){.c-told{grid-template-columns:minmax(0,1fr) minmax(0,1.1fr);align-items:center}}
.c-told blockquote{margin:0;font-size:clamp(1.5rem,3vw,2.4rem);line-height:1.2;font-weight:700;letter-spacing:-.02em;border-left:4px solid var(--warn);padding-left:24px}
.c-told pre{background:var(--surface);border:1px solid var(--line);border-radius:12px;padding:18px 20px}
.c-trade{display:grid;gap:36px;grid-template-columns:minmax(0,1fr)}
@media (min-width:1000px){.c-trade{grid-template-columns:minmax(0,1.2fr) minmax(0,1fr)}}
.c-trade .big{font-size:clamp(1.6rem,3.2vw,2.6rem);font-weight:750;line-height:1.15;letter-spacing:-.02em;margin:0}
.c-trade .big em{font-style:normal;color:var(--warn)}
.c-trade ul{margin:0;padding-left:1.1em;color:var(--muted)}.c-trade li{margin:.5em 0}
.c-start{display:grid;gap:32px;grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.c-start{grid-template-columns:minmax(0,.8fr) minmax(0,1.2fr)}}
"""


def seq(d: dict) -> str:
    p = D.PROFILE
    sig = "".join(f'<dt>{esc(s["kind"])}</dt><dd><span class="pf">←</span> {esc(s["proof"])}</dd>' for s in p["signals"])
    cand = " · ".join(f'{esc(c["skill"])} {esc(c["latest"])}' for c in p["candidates"])
    names = ", ".join(NAMES)
    inst = " · ".join(f'{esc(n)} {esc(d["latest"][n])}' for n in NAMES)
    return f"""
<figure class="c-fig" id="c-film" aria-describedby="c-cap">
<ol class="c-seq" aria-label="The three calls, in order">
<li class="c-lanes" aria-hidden="true"><span class="c-lane y">You</span><span class="c-lane">Agent</span><span class="c-lane">rsk mcp</span></li>
<li class="c-ex">
 <span class="c-no c-nod" aria-hidden="true">1</span><p class="c-route"><span class="c-no">1</span>Agent<b>→</b>rsk mcp</p>
 <div class="c-m req" data-f="a1"><div class="c-arrow r" aria-hidden="true"></div>
  <div class="c-card"><span class="m">project_profile</span><span class="k-badge ro">read-only</span><span class="pf"> {{}}</span></div></div>
 <div class="c-m" data-f="a2"><div class="c-arrow lt back" aria-hidden="true"></div>
  <div class="c-card back"><span class="sr">rsk mcp replies: </span><dl><dt>signals</dt><dd><dl class="sg" style="margin:0">{sig}</dl></dd><dt>proposes</dt><dd>{cand}</dd></dl></div></div>
</li>
<li class="c-ex h">
 <span class="c-no c-nod" aria-hidden="true">2</span><p class="c-route"><span class="c-no">2</span>Agent<b>→</b>You</p>
 <div class="c-m l req" data-f="b1"><div class="c-arrow lt ask" aria-hidden="true"></div>
  <div class="c-card"><span class="m">install_skills</span><span class="k-badge de">writes</span><span class="pf"> needs your approval</span></div></div>
 <div class="c-youcell" data-f="b2"><div class="c-you" id="c-you"><p class="q">Allow <code>install_skills</code> from rsk?</p><p class="nm">{esc(names)}</p>
  <p class="acts" style="margin:0"><span class="b">Deny</span><span class="b al"><span class="p">Allow</span><span class="d">✓ Allowed</span></span><span class="lbl">you, every call</span></p></div></div>
</li>
<li class="c-ex">
 <span class="c-no c-nod" aria-hidden="true">3</span><p class="c-route"><span class="c-no">3</span>Agent<b>→</b>rsk mcp</p>
 <div class="c-m req" data-f="c1"><div class="c-arrow r" aria-hidden="true"></div>
  <div class="c-card"><span class="m">install_skills</span><span class="k-badge de">writes</span> <span class="pf">{{"names": [{esc(", ".join(chr(34)+n+chr(34) for n in NAMES))}]}}</span></div></div>
 <div class="c-m" data-f="c2"><div class="c-arrow lt back" aria-hidden="true"></div>
  <div class="c-card back"><span class="sr">rsk mcp replies: </span><dl><dt>installed</dt><dd>{inst}</dd><dt>failed</dt><dd>none</dd></dl></div></div>
</li>
</ol>
<figcaption class="c-cap" id="c-cap"><p><span class="k-badge ro">real</span> Call 1 is a real <code>project_profile</code> response for a scratch Go folder (go.mod, Dockerfile, cmd/demo/main.go). <span class="k-badge sa">sample</span> The approval is drawn, not a capture of the client dialog, and call 3’s reply is the shape from <code>internal/mcp/install_tool.go</code> with real latest versions.</p>{film_ctl()}</figcaption>
</figure>"""


def hero(d: dict) -> str:
    return f"""
<section class="c-hero" id="top" aria-labelledby="c-h1"><div class="wrap">
<div class="c-top">
<div><p class="k-kick"><span class="gl" aria-hidden="true">◆</span> rsk mcp · for Claude Code and OpenCode</p>
<h1 class="c-h1" id="c-h1"><span>Your agent asks.</span> <span class="you">You approve.</span> <span class="rsk">rsk installs.</span></h1></div>
<div><p class="c-lede"><code>rsk</code> gives your coding agent three tools. Two only read: your project’s files and a catalogue of {len(d["skills"])} skills. The one that writes, <code>install_skills</code>, makes Claude Code ask you on every call.</p>
<div class="c-cta"><a class="k-btn primary" href="#start">Install rsk</a><a class="k-btn" href="#tools">The three tools</a></div></div>
</div>
{seq(d)}
</div></section>"""


def tools() -> str:
    ann = {
        "ro": "readOnlyHint: true\ndestructiveHint: false\nidempotentHint: true",
        "de": "readOnlyHint: false\ndestructiveHint: true\n_meta: anthropic/\n  requiresUserInteraction",
    }
    items = "".join(
        f'<li{" class=w" if k == "de" else ""}><h3>{esc(n)}</h3><p>{esc(t)}</p><pre aria-label="{esc(n)} annotations">{esc(ann[k])}</pre></li>'
        for n, t, k in D.TOOLS
    )
    items += '<li><h3>rsk://catalog</h3><p>A resource, not a tool: every skill with its latest version, descriptions cut to 200 characters. Name <code>rsk-catalog</code>.</p><pre aria-label="rsk://catalog details">mimeType:\n  application/json</pre></li>'
    return f"""
<section class="k-sec" id="tools" aria-labelledby="c-tools-h" style="background:var(--bg2)"><div class="wrap">
<p class="k-kick rv">From internal/mcp/tools.go</p>
<h2 id="c-tools-h" class="rv">Three tools. One writes, and it is marked to say so.</h2>
<p class="k-lede rv">The descriptions are the ones the agent reads. <code>destructiveHint</code> is only advice to the client; the <code>requiresUserInteraction</code> flag is what makes Claude Code prompt you even in auto-approve modes.</p>
<ul class="c-tools rv">{items}</ul>
</div></section>"""


def signals() -> str:
    hits = {s["kind"] for s in D.PROFILE["signals"]}
    rows = "".join(
        f'<tr{" class=hit" if k in hits else ""}><td>{esc(k)}</td><td>{esc(f)}</td><td>{esc(s)}</td></tr>' for k, f, s in D.SIGNALS
    )
    js = json.dumps(D.PROFILE, indent=1)
    return f"""
<section class="k-sec" id="profile" aria-labelledby="c-sig-h"><div class="wrap">
<p class="k-kick rv">From internal/mcp/signals.go</p>
<h2 id="c-sig-h" class="rv">What <code>project_profile</code> looks for.</h2>
<div class="c-two">
<div class="rv"><table class="c-sig"><caption class="sr">Signals, the files that prove them, and the skills proposed</caption><thead><tr><th scope="col">Signal</th><th scope="col">Proved by</th><th scope="col">Proposes</th></tr></thead><tbody>{rows}</tbody></table>
<p style="margin-top:14px;color:var(--muted);font-size:.92rem">It proposes, it does not rank. Refactoring skills are never proposed by a file; only when the agent is already refactoring. ● matched in the scratch folder.</p></div>
<div class="c-json rv"><p class="k-kick" style="margin-bottom:10px">The full response <span class="k-badge ro">real</span></p><pre tabindex="0" role="region" aria-label="project_profile JSON response">{esc(js)}</pre></div>
</div></div></section>"""


def told() -> str:
    return f"""
<section class="k-sec night" id="told" aria-labelledby="c-told-h"><div class="wrap c-told">
<div class="rv"><p class="k-kick">What the agent is told</p>
<h2 id="c-told-h">rsk new writes the rule into CLAUDE.md.</h2>
<blockquote><p style="margin:0">“Never install without asking.”</p></blockquote></div>
<div class="rv"><pre tabindex="0" role="region" aria-label="CLAUDE.md written by rsk new">{esc(D.RSK_NEW_CLAUDE_MD)}</pre>
<p style="margin-top:14px;color:var(--muted);font-size:.92rem"><span class="k-badge ro">real</span> From a real <code>rsk new</code> run. It also registers the server in <code>.mcp.json</code> and <code>opencode.json</code>; <code>rsk destroy</code> removes both and the note.</p></div>
</div></section>"""


def trade() -> str:
    return """
<section class="k-sec" id="trade" aria-labelledby="c-trade-h"><div class="wrap c-trade">
<div class="rv"><p class="k-kick">The trade-off</p><h2 id="c-trade-h" class="sr">The trade-off</h2>
<p class="big">One binary, one config file, every project. In exchange, <em>the registry has to be reachable</em> when you install, or you point rsk at a local clone.</p></div>
<ul class="rv">
<li><code>rsk</code> must be on PATH on every machine that opens the project; the client launches <code>rsk mcp</code> per session and stops it after.</li>
<li>The agent only proposes from files it can see. A stack with no marker file gets no proposal; ask it to use <code>search_skills</code>.</li>
<li>The catalogue is one author’s standards. Fork it and point rsk at your clone to change one.</li>
</ul></div></section>"""


def start() -> str:
    return f"""
<section class="k-sec" id="start" aria-labelledby="c-start-h" style="background:var(--bg2)"><div class="wrap c-start">
<div class="rv"><p class="k-kick">Quick start</p><h2 id="c-start-h">Set it up once. Then let the agent ask.</h2>
<p class="k-lede"><code>rsk new</code> registers the MCP server for both clients. Skip the last two lines if you want the agent to propose them.</p></div>
<div class="rv">{copy_block("c-qs", D.QUICK_START, "quick-start commands")}</div>
</div></section>"""


JS = r"""
(function(){var root=document.getElementById('c-film');if(!root||!window.kFilm)return;var you=document.getElementById('c-you');
var s=[];function at(ms,fn){s.push({at:ms,fn:fn})}function sh(k){return function(){kF.show(root,k)}}
at(250,sh('a1'));at(1150,sh('a2'));at(2450,sh('b1'));at(3050,sh('b2'));at(4250,function(){you.classList.remove('pend')});
at(4850,sh('c1'));at(5750,sh('c2'));
kFilm(root,{total:6800,steps:s,reset:function(){kF.hide(root);you.classList.add('pend')},end:function(){kF.end(root);you.classList.remove('pend')}});
})();
"""


def build(d: dict) -> str:
    nav = [("Tools", "#tools", True), ("Detection", "#profile", True), ("Trade-off", "#trade", True), ("Install", "#start", False)]
    body = (
        header(nav, d["version"])
        + '<main id="main">'
        + hero(d)
        + tools()
        + signals()
        + told()
        + trade()
        + start()
        + "</main>"
        + footer(D.SOURCES, d["version"])
    )
    return page(
        "rsk — your agent asks, you approve, rsk installs (direction C)",
        "rsk mcp gives Claude Code and OpenCode three tools to profile a project and install the skills that apply, with your approval.",
        CSS,
        body,
        JS,
    )

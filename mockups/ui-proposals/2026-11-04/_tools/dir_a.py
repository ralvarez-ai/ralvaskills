"""Direction A — the manifest as the argument.

Device: a real-format .rsk/rsk.mod that assembles itself as `rsk new`,
`rsk install go-grpc` and `rsk pin go-architect` run, with rsk.lock resolving
beneath it. Personality: exact, editorial; paper page, one dark terminal.
"""

from __future__ import annotations

import data as D
from family import copy_block, esc, film_ctl, footer, header, page, rail, wordmark_svg

GRPC = ["go-architect", "grpc-architect", "protobuf-architect", "docker-architect", "sql-architect"]

CSS = r"""
.a-hero{position:relative;padding:clamp(32px,6vw,88px) 0 clamp(64px,9vw,128px)}
.a-hero .wrap{position:relative;z-index:1}
.a-grid{display:grid;gap:22px 64px;grid-template-columns:minmax(0,1fr);grid-template-areas:"kick" "h" "dev" "lede" "cta"}
@media (min-width:1000px){.a-grid{grid-template-columns:minmax(0,5fr) minmax(0,6.4fr);grid-template-rows:auto auto auto 1fr;grid-template-areas:"kick dev" "h dev" "lede dev" "cta dev"}}
.a-kick{grid-area:kick}.a-h{grid-area:h}.a-lede{grid-area:lede}.a-cta{grid-area:cta}.a-dev{grid-area:dev}
.a-h{font-size:clamp(2.6rem,5.2vw,4.5rem);font-weight:800;letter-spacing:-.05em;line-height:.96}
.a-h code{font:inherit;font-family:var(--mono);font-weight:700;letter-spacing:-.04em;color:var(--accent-t);font-size:.94em}
.a-lede{font-size:clamp(1.08rem,1.5vw,1.28rem);color:var(--muted);max-width:44ch;margin:0}
.a-lede code{color:var(--text)}
.a-cta{display:flex;flex-wrap:wrap;gap:12px;align-items:center}
.a-dev{margin:0;display:grid;gap:14px;align-self:start}
.a-term{border-radius:var(--rc);overflow:hidden;border:1px solid var(--line);box-shadow:var(--shadow)}
.a-cmds{list-style:none;margin:0;padding:14px 18px 16px;font:500 .9rem/1.75 var(--mono)}
.a-cmds .ps{color:var(--accent-t);margin-right:.6em}
.a-cmds .out{color:var(--ok);padding-left:1.3em}
.a-cmds .ag{color:var(--warn)}.a-cmds .ag .ps{color:var(--warn)}.a-cmds .ag code{color:var(--warn)}
.a-cmds code{font-size:1em;color:var(--text)}
.a-files{display:grid;gap:14px;grid-template-columns:minmax(0,1fr)}
@media (min-width:600px){.a-files{grid-template-columns:minmax(0,1.35fr) minmax(0,1fr)}}
.a-mod{font:500 .9rem/1.75 var(--mono);margin:0;padding:12px 0;list-style:none;counter-reset:ln}
.a-mod li{display:grid;grid-template-columns:2.4em 1.4em minmax(0,1fr);padding-right:12px;counter-increment:ln;white-space:pre-wrap;overflow-wrap:anywhere}
.a-mod li::before{content:counter(ln);color:var(--muted);text-align:right;padding-right:.7em;font-size:.82em}
.a-mod .g{color:var(--ok);font-weight:700}
.a-mod li.add{background:color-mix(in srgb,var(--ok) 9%,transparent)}
.a-mod .k{color:var(--text)}.a-mod .s{color:var(--ok)}.a-mod .n{color:var(--warn)}.a-mod .t{color:var(--accent-t);font-weight:700}
.a-mod li.chg{background:color-mix(in srgb,var(--accent) 16%,transparent)}
.a-lock{width:100%;border-collapse:collapse;font:500 .84rem/1.4 var(--mono)}
.a-lock caption{text-align:left;padding:0 0 6px;color:var(--muted);font:500 .78rem var(--sans)}
.a-lock th{text-align:left;font-weight:600;color:var(--muted);padding:6px 12px;border-bottom:1px solid var(--line);font-size:.72rem;text-transform:uppercase;letter-spacing:.06em}
.a-lock td{white-space:nowrap;padding:8px 12px;border-bottom:1px solid color-mix(in srgb,var(--line) 55%,transparent)}
.a-lock td.v{color:var(--accent-t);font-weight:700;text-align:right}
.a-lock th.v{text-align:right}
.a-cap{font-size:.84rem;color:var(--muted);margin:0;display:flex;flex-wrap:wrap;gap:8px 14px;align-items:center;justify-content:space-between}
.a-cap p{margin:0;max-width:62ch}
/* files section */
.a-two{display:grid;gap:clamp(28px,4vw,64px);grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.a-two{grid-template-columns:minmax(0,1fr) minmax(0,1fr);align-items:start}}
.a-tree{list-style:none;margin:0;padding:0;font:500 .92rem/1.5 var(--mono);border-top:2px solid var(--text)}
.a-tree li{display:grid;grid-template-columns:1.6em minmax(0,15em) minmax(0,1fr);gap:8px;padding:10px 0;border-bottom:1px solid var(--line);align-items:baseline}
.a-tree .g{color:var(--ok);font-weight:700}
.a-tree .yours{color:var(--muted)}
.a-tree .d{font:400 .9rem/1.45 var(--sans);color:var(--muted)}
@media (max-width:600px){.a-tree li{grid-template-columns:1.6em minmax(0,1fr)}.a-tree .d{grid-column:2}}
.a-step{font:600 .72rem var(--mono);letter-spacing:.08em;text-transform:uppercase;color:var(--muted);margin:36px 0 10px}
/* this repo, night band */
.a-night h2{color:var(--text)}
.a-real{display:grid;gap:28px;grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.a-real{grid-template-columns:minmax(0,1.1fr) minmax(0,1fr)}}
.a-real pre{font-size:clamp(.86rem,1.2vw,1rem);line-height:1.75;padding:22px 24px;background:var(--surface);border:1px solid var(--line);border-radius:var(--rc)}
.a-notes{list-style:none;margin:0;padding:0;display:grid;gap:18px}
.a-notes li{border-left:2px solid var(--accent-t);padding:2px 0 2px 16px;color:var(--muted)}
.a-notes b{color:var(--text)}
/* commands */
.a-help{list-style:none;margin:0;padding:0;border-top:2px solid var(--text)}
.a-help>li{display:grid;grid-template-columns:minmax(0,1fr);gap:6px 32px;padding:18px 0;border-bottom:1px solid var(--line)}
@media (min-width:800px){.a-help>li{grid-template-columns:minmax(0,16em) minmax(0,1fr)}}
.a-help h3{font:600 .72rem var(--mono);letter-spacing:.08em;text-transform:uppercase;color:var(--muted);padding-top:4px}
.a-help dl{margin:0;display:grid;gap:6px}
.a-help dl div{display:grid;grid-template-columns:7.5em minmax(0,1fr);gap:12px}
.a-help dt{font:700 1rem var(--mono);color:var(--accent-t)}.a-help dd{margin:0}
.a-start{display:grid;gap:clamp(24px,4vw,56px);grid-template-columns:minmax(0,1fr)}
@media (min-width:980px){.a-start{grid-template-columns:minmax(0,.8fr) minmax(0,1.2fr);align-items:start}}
"""


def _mod_lines(names: list[str] = GRPC) -> list[tuple[str, str]]:
    """(data-f key, html) per rsk.mod line, in file order (skills sorted)."""
    s = lambda x: f'<span class="s">{esc(x)}</span>'  # noqa: E731
    lines = [
        ("new", '<span class="k">version</span> = <span class="n">1</span>'),
        ("new", f'<span class="k">tools</span> = [{s(chr(34) + "claude-code" + chr(34))}, {s(chr(34) + "opencode" + chr(34))}]'),
        ("new", f'<span class="k">pinned</span> = <span class="pin">[{s(chr(34) + "go-architect" + chr(34))}]</span>'),
        ("new", ""),
        ("new", '<span class="t">[skills]</span>'),
    ]
    for n in sorted(names):
        lines.append((f"s{names.index(n)}", f'  <span class="k">{esc(n)}</span> = {s(chr(34) + "*" + chr(34))}'))
    return lines


def device(d: dict, names: list[str] = GRPC, c2: str = "rsk install go-grpc", agent: bool = False, cap: str = "") -> str:
    mod = "".join(
        f'<li class="add" data-f="{k}"{" id=a-pin" if "pinned" in h else ""}><span class="g" aria-hidden="true">+</span><span>{h}</span></li>'
        for k, h in _mod_lines(names)
    )
    lock = "".join(
        f'<tr data-f="s{i}"><td>{esc(n)}</td><td class="v">{esc(d["latest"][n])}</td></tr>'
        for i, n in enumerate(names)
    )
    cmd = lambda k, c: (  # noqa: E731
        f'<li data-f="{k}"><span class="ps" aria-hidden="true">$</span><code><span data-t="{esc(c)}" aria-hidden="true">{esc(c)}</span><span class="sr">{esc(c)}</span></code></li>'
    )
    out = "".join(f'<li class="out" data-f="o1">{esc(o)}</li>' for o in D.RSK_NEW_OUT)
    if agent:
        line2 = f'<li data-f="c2" class="ag"><span class="ps" aria-hidden="true">◆</span><code><span data-t="{esc(c2)}" aria-hidden="true">{esc(c2)}</span><span class="sr">{esc(c2)}</span></code></li>'
    else:
        line2 = cmd("c2", c2)
    cap = cap or f"""<span class="k-badge sa">sample project</span> The <code>rsk new</code> lines are a real run (paths shortened). The install and pin results are what <code>internal/install/apply.go</code> writes for the real <code>go-grpc</code> bundle, with versions from the registry on 2026-11-04 — not a capture."""
    return f"""<figure class="a-dev" id="a-film" aria-describedby="a-cap">
<div class="a-term tn"><div class="k-win-h"><span class="dots" aria-hidden="true"><i></i><i></i><i></i></span>~/demo — a Go project</div>
<ol class="a-cmds" aria-label="Commands run">{cmd("c1", "rsk new")}{out}{line2}{cmd("c3", "rsk pin go-architect")}</ol></div>
<div class="a-files">
<div class="k-win"><div class="k-win-h">.rsk/rsk.mod <span class="k-badge ro" style="margin-left:auto">commit</span></div><ol class="a-mod" aria-label=".rsk/rsk.mod contents">{mod}</ol></div>
<div class="k-win"><div class="k-win-h">.rsk/rsk.lock <span class="k-badge ro" style="margin-left:auto">commit</span></div><div class="k-win-b" style="padding:6px 0 4px">
<table class="a-lock"><caption class="sr">rsk.lock resolved versions, all from the registry, shown as a table</caption><thead><tr><th scope="col">skill</th><th scope="col" class="v">version</th></tr></thead><tbody>{lock}</tbody></table></div></div>
</div>
<figcaption class="a-cap" id="a-cap"><p>{cap}</p>{film_ctl()}</figcaption>
</figure>"""


def hero(d: dict) -> str:
    return f"""
<section class="a-hero" id="top" aria-labelledby="a-h1"><div class="k-grid" aria-hidden="true"></div><div class="wrap a-grid">
<p class="k-kick a-kick rv"><span class="fd" aria-hidden="true"><i></i><i></i></span> Skills for Claude Code + OpenCode</p>
<h1 class="a-h rv" id="a-h1" style="--i:1"><code>go.mod</code>, but for your agent’s skills.</h1>
<p class="a-lede rv" style="--i:2">rsk installs versioned skills from a catalogue of {len(d["skills"])} into Claude Code and OpenCode, and writes what it did into two files you commit. A fresh clone runs <code>rsk install</code> and gets the same skills at the same versions.</p>
<div class="a-cta rv" style="--i:3"><a class="k-btn primary" href="#start">Install rsk</a><a class="k-btn" href="#files">See what it writes</a></div>
{device(d)}
</div></section>"""


def files(d: dict) -> str:
    rows_new = [
        ("+", ".rsk/rsk.mod", "The manifest: which tools, which skills, which are pinned. Commit it."),
        ("+", ".rsk/CLAUDE.md", "Imports for pinned skills; empty until you pin one. Commit it."),
        ("+", "CLAUDE.md", "Starts with @AGENTS.md, then @.rsk/CLAUDE.md, then a short note for the agent."),
        ("+", "AGENTS.md", "A skeleton, only when the project has none. An existing one is never touched."),
        ("+", ".mcp.json", "Registers rsk mcp for Claude Code."),
        ("+", "opencode.json", "Registers rsk mcp for OpenCode."),
    ]
    rows_inst = [
        ("+", ".rsk/rsk.lock", "Exact versions and sources. Commit it, like go.sum."),
        ("+", ".claude/skills/<name>", "A symlink into rsk’s cache under ~/.config/rsk/. OpenCode reads the same folder."),
    ]
    tree = lambda rows: "".join(  # noqa: E731
        f'<li><span class="g" aria-hidden="true">{g}</span><span>{esc(p)}<span class="sr"> (new)</span></span><span class="d">{esc(t)}</span></li>'
        for g, p, t in rows
    )
    yours = "".join(
        f'<li><span aria-hidden="true"></span><span class="yours">{esc(p)}</span><span class="d">yours, untouched</span></li>'
        for p in ("go.mod", "Dockerfile", "cmd/demo/main.go")
    )
    return f"""
<section class="k-sec" id="files" aria-labelledby="a-files-h"><div class="wrap a-two">
<div class="rv">
<div class="k-rule" aria-hidden="true"></div>
<p class="k-kick">Before and after</p>
<h2 id="a-files-h">What <code>rsk new</code> actually writes.</h2>
<p class="k-lede">Six files in the project root and <code>.rsk/</code>, from a real run on a three-file Go project. Installing adds the lock and the links. Nothing lands outside the project except the cache.</p>
<p class="a-step">rsk new <span class="k-badge ro">real run</span></p>
<ul class="a-tree" aria-label="Files after rsk new">{yours}{tree(rows_new)}</ul>
<p class="a-step">rsk install go-grpc <span class="k-badge sa">from the code</span></p>
<ul class="a-tree" aria-label="Files added by rsk install">{tree(rows_inst)}</ul>
</div>
<div class="rv">
<div class="k-win"><div class="k-win-h">CLAUDE.md, as rsk new wrote it <span class="k-badge ro" style="margin-left:auto">real</span></div>
<pre class="k-win-b" tabindex="0" role="region" aria-label="CLAUDE.md written by rsk new">{esc(D.RSK_NEW_CLAUDE_MD)}</pre></div>
<p style="margin-top:16px;color:var(--muted);font-size:.95rem"><code>@AGENTS.md</code> stays the first line so the guide you wrote loads before anything rsk generates. <code>rsk destroy</code> removes the note and the MCP registration and leaves <code>AGENTS.md</code> alone.</p>
</div>
</div></section>"""


def this_repo(d: dict) -> str:
    lock = "".join(
        f'<tr><td>{esc(e["name"])}</td><td class="v">{esc(e["version"])}</td><td>{esc(e["source"])}</td></tr>'
        for e in sorted(d["repo_lock"], key=lambda e: e["name"])
    )
    n = len(d["repo_lock"])
    return f"""
<section class="k-sec night a-night" id="this-repo" aria-labelledby="a-real-h"><div class="wrap">
<div class="k-rule" aria-hidden="true"></div>
<p class="k-kick rv">The real file</p>
<h2 id="a-real-h" class="rv">The repository that ships rsk runs on rsk.</h2>
<div class="a-real">
<div class="rv"><div class="k-win"><div class="k-win-h">ralvaskills/.rsk/rsk.mod <span class="k-badge ro" style="margin-left:auto">real · working tree 2026-11-04</span></div>
<pre class="k-win-b" tabindex="0" role="region" aria-label="ralvaskills .rsk/rsk.mod">{esc(d["repo_mod"].rstrip())}</pre></div></div>
<div class="rv">
<ul class="a-notes">
<li><b>{n} skills, {len(d["repo_pinned"])} pinned.</b> Pinned skills are imported through <code>.rsk/CLAUDE.md</code> and <code>opencode.json</code>, so they load in every session.</li>
<li><b><code>"*"</code> means latest;</b> the lock records the exact version that resolved, so a clone gets the same ones.</li>
<li><b>Both tools from one list.</b> <code>tools = ["claude-code", "opencode"]</code> drives which adapter files rsk writes.</li>
</ul>
<table class="a-lock" style="margin-top:28px"><caption>ralvaskills/.rsk/rsk.lock, shown as a table</caption><thead><tr><th scope="col">skill</th><th scope="col" class="v">version</th><th scope="col">source</th></tr></thead><tbody>{lock}</tbody></table>
</div></div>
</div></section>"""


def commands() -> str:
    groups = "".join(
        f'<li><h3>{esc(g)}</h3><dl>'
        + "".join(f"<div><dt>{esc(c)}</dt><dd>{esc(t)}</dd></div>" for c, t in cmds)
        + "</dl></li>"
        for g, cmds in D.HELP_GROUPS
    )
    return f"""
<section class="k-sec" id="commands" aria-labelledby="a-cmd-h"><div class="wrap">
<div class="k-rule" aria-hidden="true"></div>
<p class="k-kick rv">The whole surface</p>
<h2 id="a-cmd-h" class="rv">Thirteen commands, grouped the way <code>rsk --help</code> groups them.</h2>
<ul class="a-help rv">{groups}</ul>
<p style="margin-top:20px;color:var(--muted)">Data on stdout, logs on stderr, <code>--output json</code> on the read-only views, <code>NO_COLOR</code> honoured.</p>
</div></section>"""


def trade() -> str:
    return """
<section class="k-sec" id="trade" aria-labelledby="a-trade-h" style="background:var(--bg2)"><div class="wrap">
<div class="k-rule" aria-hidden="true"></div>
<p class="k-kick rv">The trade-off</p>
<h2 id="a-trade-h" class="rv">One binary and one config file for every project. One dependency in exchange.</h2>
<div class="k-trade">
<div class="give rv"><h3>What you get</h3><ul>
<li>One binary and one machine config, <code>~/.config/rsk/config.json</code>, for every project.</li>
<li>Versions in a lock you commit. A clone runs <code>rsk install</code> and matches.</li>
<li>Claude Code and OpenCode from one install: skills are linked once into <code>.claude/skills/</code>.</li>
</ul></div>
<div class="take rv"><h3>What you accept</h3><ul>
<li><b>The registry at skills.ralvarez.dev has to be reachable</b> when you install or update, or you point rsk at a local clone of the repository.</li>
<li>The MCP server needs <code>rsk</code> on PATH on every machine that opens the project.</li>
<li>The catalogue is one author’s standards. To change one, fork the repo and use your clone.</li>
</ul></div>
</div></div></section>"""


def start() -> str:
    return f"""
<section class="k-sec" id="start" aria-labelledby="a-start-h"><div class="wrap a-start">
<div class="k-rule" aria-hidden="true" style="grid-column:1/-1"></div>
<div class="rv"><p class="k-kick">Quick start</p><h2 id="a-start-h">Five commands from nothing to a pinned skill.</h2>
<p class="k-lede">Needs Go to build from source, or take a binary from the Releases page.</p></div>
<div class="rv">{copy_block("a-qs", D.QUICK_START, "quick-start commands")}</div>
</div></section>"""


JS = r"""
(function(){var root=document.getElementById('a-film');if(!root||!window.kFilm)return;
var pin=root.querySelector('#a-pin .pin'),pinLine=root.querySelector('#a-pin'),full=pin.innerHTML;
var steps=[];function at(ms,fn,dur){steps.push({at:ms,dur:dur||0,fn:fn})}
var show=function(k){return function(){kF.show(root,k)}};
var typer=function(k){var el=root.querySelector('[data-f="'+k+'"] [data-t]');return function(p){kF.type(el,p)}};
at(150,show('c1'));at(200,typer('c1'),550);
at(950,show('o1'));at(1150,show('new'));
at(1700,show('c2'));at(1750,typer('c2'),900);
for(var i=0,n=root.querySelectorAll('tr[data-f]').length;i<n;i++)at(2850+i*230,show('s'+i));
at(4200,show('c3'));at(4250,typer('c3'),800);
at(5250,function(){pin.innerHTML=full;pinLine.classList.add('chg')});
kFilm(root,{total:6400,steps:steps,
 reset:function(){kF.hide(root);root.querySelectorAll('[data-t]').forEach(function(e){e.textContent=''});pin.innerHTML='[]';pinLine.classList.remove('chg')},
 end:function(){kF.end(root);pin.innerHTML=full;setTimeout(function(){pinLine.classList.remove('chg')},900)}});
})();
"""


def build(d: dict) -> str:
    nav = [("What it writes", "#files", True), ("Commands", "#commands", True), ("Trade-off", "#trade", True), ("Install", "#start", False)]
    body = (
        header(nav, d["version"])
        + rail([("top", "The manifest"), ("files", "rsk new"), ("this-repo", "Real file"), ("commands", "Commands"), ("trade", "Trade-off"), ("start", "Install")])
        + '<main id="main">'
        + hero(d)
        + files(d)
        + this_repo(d)
        + commands()
        + trade()
        + start()
        + "</main>"
        + footer(D.SOURCES, d["version"])
    )
    return page(
        "rsk — go.mod for your agent’s skills (direction A)",
        "rsk installs versioned skills into Claude Code and OpenCode and records them in .rsk/rsk.mod and rsk.lock.",
        CSS,
        body,
        JS,
    )

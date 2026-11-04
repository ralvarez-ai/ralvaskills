"""Combined page — direction C's agent loop as the spine, with A's manifest as
the result of call 3 and B's count as the catalogue the agent chooses from.

One story, one scratch project: project_profile finds go.mod, Dockerfile and
cmd/demo/main.go → you approve → install_skills writes those three skills into
rsk.mod / rsk.lock → the 58 they came from → the tools and rules behind it.
"""

from __future__ import annotations

import data as D
import dir_a as A
import dir_b as B
import dir_c as C
from family import esc, footer, header, page

SECTIONS = [
    ("top", "The loop"),
    ("writes", "What it writes"),
    ("files", "rsk new"),
    ("catalogue", "Catalogue"),
    ("tools", "Tools"),
    ("profile", "Detection"),
    ("this-repo", "Real file"),
    ("commands", "Commands"),
    ("trade", "Trade-off"),
    ("start", "Install"),
    ("faq", "FAQ"),
]

CSS = r"""
.d-writes{display:grid;gap:32px clamp(32px,5vw,72px);grid-template-columns:minmax(0,1fr)}
@media (min-width:1080px){.d-writes{grid-template-columns:minmax(0,4fr) minmax(0,7fr);align-items:start}}
.d-writes .k-lede{max-width:40ch}
.d-writes .a-dev{grid-area:auto}
#catalogue .b-count{font-size:clamp(2.4rem,4vw,3.4rem)}
.d-call{display:inline-flex;align-items:center;gap:10px;font:700 .85rem var(--mono);color:var(--warn);margin-bottom:6px}
.d-call .c-no{color:var(--accent-t);margin:0}
.d-rail{display:none}
@media (min-width:1400px){
 .d-rail{display:block;position:fixed;right:14px;top:50%;transform:translateY(-50%);z-index:15}
 .d-rail ol{list-style:none;margin:0;padding:0;display:grid;gap:2px}
 .d-rail a{display:flex;flex-direction:row-reverse;align-items:center;gap:10px;min-height:44px;min-width:44px;padding:0 10px;border-radius:8px;color:var(--muted);text-decoration:none;font:600 .78rem var(--mono)}
 .d-rail a .lb{opacity:0;transform:translateX(4px);transition:opacity .2s,transform .2s;white-space:nowrap;background:var(--surface);border:1px solid var(--line);padding:4px 8px;border-radius:6px;color:var(--text)}
 .d-rail a:hover .lb,.d-rail a:focus-visible .lb{opacity:1;transform:none}
 .d-rail a[aria-current=true]{color:var(--accent-t)}
 .d-rail a[aria-current=true] .n{border-bottom:2px solid currentColor}
}
.d-faq{display:grid;gap:0;grid-template-columns:minmax(0,1fr);border-top:2px solid var(--text);max-width:900px}
.d-faq details{border-bottom:1px solid var(--line)}
.d-faq summary{min-height:56px;display:flex;align-items:center;justify-content:space-between;gap:16px;cursor:pointer;font-weight:700;font-size:1.1rem;list-style:none}
.d-faq summary::-webkit-details-marker{display:none}
.d-faq summary::after{content:"+";font:700 1.4rem var(--mono);color:var(--accent-t)}
.d-faq details[open] summary::after{content:"−"}
.d-faq details p{color:var(--muted);margin:0 0 18px;max-width:70ch}
"""


def rail() -> str:
    items = "".join(
        f'<li><a href="#{i}"><span class="n">{n:02d}</span><span class="lb">{esc(t)}</span></a></li>' for n, (i, t) in enumerate(SECTIONS, 1)
    )
    return f'<nav class="d-rail" aria-label="Sections"><ol>{items}</ol></nav>'


def writes(d: dict) -> str:
    names = C.NAMES
    cap = (
        '<span class="k-badge sa">sample project</span> Same scratch folder as above. The <code>rsk new</code> lines are a real run (paths shortened); '
        "the ◆ line stands for call 3. The file contents are what <code>internal/install/apply.go</code> writes for those three skills, "
        "with versions from the registry on 2026-11-04 — not a capture."
    )
    dev = A.device(d, names=names, c2="install_skills · approved by you", agent=True, cap=cap)
    return f"""
<section class="k-sec" id="writes" aria-labelledby="d-wr-h" style="background:var(--bg2)"><div class="wrap d-writes">
<div class="rv"><p class="d-call"><span class="c-no">3</span> after call 3</p>
<h2 id="d-wr-h">It lands in two files you commit.</h2>
<p class="k-lede"><code>install_skills</code> writes the skills into <code>.rsk/rsk.mod</code> and their exact versions into <code>.rsk/rsk.lock</code>. Whoever clones the repo runs <code>rsk install</code> and gets the same three, at the same versions. Think <code>go.mod</code> and <code>go.sum</code>, for your agent’s skills.</p></div>
{dev}
</div></section>"""


def catalogue(d: dict) -> str:
    html = B.explorer(d)
    nb = len(d["bundles"])
    a = f'<div><span class="b-sec-n" aria-hidden="true">{nb}</span><h2 id="b-ex-h" style="margin-top:0"><span class="sr">{nb} </span>bundles. Filter the catalogue the way <code>rsk catalog</code> does.</h2></div>'
    n = len(d["skills"])
    b = (
        f'<div><span class="b-sec-n" aria-hidden="true">{n}</span><h2 id="b-ex-h" style="margin-top:0"><span class="sr">{n} </span>skills the agent chooses from, in {nb} bundles.</h2>'
        f'<p class="k-lede">Grouped by the folder they live in. Pick a bundle to see what <code>rsk install &lt;bundle&gt;</code> brings; personal skills stay hidden unless you pass <code>--include personal</code>, as in the CLI.</p></div>'
    )
    assert a in html
    return html.replace(a, b)


def faq(d: dict) -> str:
    qa = [
        ("Do I need the MCP server?", "No. Every step works from the CLI: <code>rsk install</code>, <code>rsk pin</code>, <code>rsk catalog</code>. The server is there so the agent can propose skills from your files instead of you naming them."),
        ("Which clients does it support?", "Claude Code and OpenCode. <code>rsk new --for claude-code|opencode|all</code> picks which config files rsk writes; installed skills are linked once into <code>.claude/skills/</code>."),
        ("Can I install a skill for every project?", "Yes: <code>rsk install &lt;name&gt; --global</code> links it into <code>~/.claude/skills/</code> and <code>~/.config/opencode/skills/</code>. The <code>global</code> bundle is meant for that."),
        ("Can I define my own bundles?", "Write <code>~/.config/rsk/catalog.toml</code>. A bundle with the same name replaces the built-in one; new names are added."),
        ("What are personal skills?", f"{d['personal']} of the {len(d['skills'])} are marked personal in the registry index. <code>rsk catalog</code> hides them unless you pass <code>--include personal</code>."),
        ("How do updates work?", "<code>rsk status</code> shows installed skills and version drift; <code>rsk update [name]</code> pulls the latest and re-links. Each skill is versioned and released on its own."),
        ("What does it cost?", "Nothing. rsk and the catalogue are GPL-3.0, on GitHub."),
    ]
    items = "".join(f"<details><summary>{q}</summary><p>{a}</p></details>" for q, a in qa)
    return f"""
<section class="k-sec" id="faq" aria-labelledby="d-faq-h"><div class="wrap">
<p class="k-kick rv">Questions</p><h2 id="d-faq-h" class="rv">Before you install.</h2>
<div class="d-faq rv">{items}</div>
</div></section>"""


RAIL_JS = r"""
(function(){var links=[].slice.call(document.querySelectorAll('.d-rail a'));if(!links.length||!('IntersectionObserver' in window))return;
var map={};links.forEach(function(a){map[a.getAttribute('href').slice(1)]=a});
var io=new IntersectionObserver(function(es){es.forEach(function(e){if(e.isIntersecting){links.forEach(function(a){a.removeAttribute('aria-current')});var a=map[e.target.id];if(a)a.setAttribute('aria-current','true')}})},{rootMargin:'-45% 0px -50% 0px'});
Object.keys(map).forEach(function(id){var s=document.getElementById(id);if(s)io.observe(s)});})();
"""


def build(d: dict) -> str:
    nav = [("How it works", "#top", True), ("Catalogue", "#catalogue", True), ("Trade-off", "#trade", True), ("Install", "#start", False)]
    tools = C.tools().replace(' style="background:var(--bg2)"', "", 1)
    start = C.start().replace(' style="background:var(--bg2)"', "", 1)
    body = (
        header(nav, d["version"])
        + rail()
        + '<main id="main">'
        + C.hero(d)
        + writes(d)
        + A.files(d)
        + catalogue(d)
        + tools
        + C.signals()
        + A.this_repo(d)
        + A.commands()
        + A.trade()
        + start
        + faq(d)
        + "</main>"
        + footer(D.SOURCES, d["version"])
    )
    return page(
        "rsk — your agent asks, you approve, rsk installs",
        "rsk installs versioned skills into Claude Code and OpenCode. Its MCP server lets the agent propose them from your project's files; the one tool that writes asks you first.",
        A.CSS + B.CSS + C.CSS + CSS,
        body,
        C.JS + A.JS + B.JS + RAIL_JS,
    )

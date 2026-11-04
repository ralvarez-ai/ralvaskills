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
from family import esc, footer, header, page, rail

TITLE = "rsk — your agent asks, you approve, rsk installs"
DESC = (
    "rsk installs versioned skills into Claude Code and OpenCode. Its MCP server "
    "lets the agent propose them from your project's files; the one tool that "
    "writes asks you first."
)

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
@media (min-width:1080px){.d-writes>.rv:first-child{position:sticky;top:96px;align-self:start}}
#catalogue .b-count{font-size:clamp(2.4rem,4vw,3.4rem)}
.d-call{display:inline-flex;align-items:center;gap:10px;font:600 .72rem var(--mono);letter-spacing:.08em;text-transform:uppercase;color:var(--warn);margin-bottom:6px}
.d-call .c-no{color:var(--accent-t);margin:0;width:1.7em;height:1.7em}
.d-faq{display:grid;gap:0;grid-template-columns:minmax(0,1fr);border-top:2px solid var(--text);max-width:900px}
.d-faq details{border-bottom:1px solid var(--line)}
.d-faq summary{min-height:56px;display:flex;align-items:center;justify-content:space-between;gap:16px;cursor:pointer;font-weight:700;font-size:1.1rem;list-style:none}
.d-faq summary::-webkit-details-marker{display:none}
.d-faq summary::after{content:"+";font:700 1.4rem var(--mono);color:var(--accent-t)}
.d-faq details[open] summary::after{content:"−"}
.d-faq details p{color:var(--muted);margin:0 0 18px;max-width:70ch}
.d-writes .rv:first-child .k-rule{margin-bottom:22px}
"""


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
<div class="rv"><div class="k-rule" aria-hidden="true"></div><p class="d-call"><span class="c-no">3</span> after call 3</p>
<h2 id="d-wr-h">It lands in two files you commit.</h2>
<p class="k-lede"><code>install_skills</code> writes the skills into <code>.rsk/rsk.mod</code> and their exact versions into <code>.rsk/rsk.lock</code>. Whoever clones the repo runs <code>rsk install</code> and gets the same three, at the same versions. Think <code>go.mod</code> and <code>go.sum</code>, for your agent’s skills.</p></div>
{dev}
</div></section>"""


def catalogue(d: dict) -> str:
    html = B.explorer(d)
    nb = len(d["bundles"])
    a = f'<div><span class="b-sec-n" aria-hidden="true" data-count="{nb}">{nb}</span><h2 id="b-ex-h" style="margin-top:0"><span class="sr">{nb} </span>bundles. Filter the catalogue the way <code>rsk catalog</code> does.</h2></div>'
    n = len(d["skills"])
    b = (
        f'<div><span class="b-sec-n" aria-hidden="true" data-count="{n}">{n}</span><h2 id="b-ex-h" style="margin-top:0"><span class="sr">{n} </span>skills the agent chooses from, in {nb} bundles.</h2>'
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
<div class="k-rule" aria-hidden="true"></div>
<p class="k-kick rv">Questions</p><h2 id="d-faq-h" class="rv">Before you install.</h2>
<div class="d-faq rv">{items}</div>
</div></section>"""


def build(d: dict, extra_head: str = "") -> str:
    nav = [("How it works", "#top", True), ("Catalogue", "#catalogue", True), ("Trade-off", "#trade", True), ("Install", "#start", False)]
    tools = C.tools().replace(' style="background:var(--bg2)"', "", 1)
    start = C.start().replace(' style="background:var(--bg2)"', "", 1)
    body = (
        header(nav, d["version"])
        + rail(SECTIONS)
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
        TITLE,
        DESC,
        A.CSS + B.CSS + C.CSS + CSS,
        body,
        C.JS + A.JS + B.JS,
        extra_head=extra_head,
    )

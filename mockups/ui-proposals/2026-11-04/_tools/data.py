"""Facts for the rsk landing proposals, read from the repository at build time.

Every number on the pages comes from here. The asserts are the anchors from
docs/plans/2026-11-04-landing-page-brief.md: if the tree drifts, the build
fails instead of the page lying.
"""

from __future__ import annotations

import datetime as dt
import json
import re
import subprocess
import tomllib
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
OUT = Path(__file__).resolve().parents[1]


def _first_sentence(s: str) -> str:
    m = re.match(r"(.+?[.!?])(\s|$)", s.strip())
    return (m.group(1) if m else s).strip()


def load() -> dict:
    idx = json.loads((ROOT / "registry/index.json").read_text())["skills"]
    skills = []
    for p in sorted((ROOT / "skills").glob("*/*/SKILL.md")):
        dom, name = p.parent.parent.name, p.parent.name
        e = idx[name]
        skills.append(
            {
                "name": name,
                "domain": dom,
                "latest": e["latest"],
                "releases": len(e["versions"]),
                "personal": bool(e.get("personal")),
                "desc": _first_sentence(e["description"]),
            }
        )
    assert len(skills) == 58 == len(idx), len(skills)

    domains: dict[str, list[dict]] = {}
    for s in skills:
        domains.setdefault(s["domain"], []).append(s)

    bundles = tomllib.loads((ROOT / "internal/config/catalog.toml").read_text())["bundle"]
    assert len(bundles) == 15
    names = {s["name"] for s in skills}
    for b in bundles:
        b["in_catalogue"] = [x["name"] for x in b["skills"] if x["name"] in names]
        b["official"] = [x["name"] for x in b["skills"] if x["source"] == "official"]
        b["missing"] = [
            x["name"]
            for x in b["skills"]
            if x["source"] != "official" and x["name"] not in names
        ]

    setup = (ROOT / "cmd/rsk/setup.go").read_text()
    assert setup.count("rootCmd.AddCommand(") == 12  # + cobra's default `completion` = 13

    tools_go = (ROOT / "internal/mcp/tools.go").read_text()
    tools = re.findall(r'Name:\s+"([a-z_]+)"', tools_go)
    assert tools == ["project_profile", "search_skills", "install_skills"], tools
    assert tools_go.count("mutatingAnnotations(\"") == 1
    res = (ROOT / "internal/mcp/resource.go").read_text()
    assert '"rsk://catalog"' in res and '"rsk-catalog"' in res

    version = subprocess.run(
        ["git", "describe", "--tags", "--abbrev=0", "--match", "v*"],
        cwd=ROOT,
        capture_output=True,
        text=True,
        check=True,
    ).stdout.strip()

    stamps = [
        dt.datetime.fromisoformat(v["published_at"].replace("Z", "+00:00"))
        for e in idx.values()
        for v in e["versions"].values()
    ]
    first = min(stamps)
    week0 = (first - dt.timedelta(days=first.weekday())).date()
    weeks = Counter(((s.date() - week0).days // 7) for s in stamps)
    nweeks = max(weeks) + 1
    hist = [(week0 + dt.timedelta(weeks=i), weeks.get(i, 0)) for i in range(nweeks)]

    mod = (ROOT / ".rsk/rsk.mod").read_text()
    lock = tomllib.loads((ROOT / ".rsk/rsk.lock").read_text())["skills"]

    return {
        "skills": skills,
        "domains": domains,
        "bundles": bundles,
        "version": version.lstrip("v"),
        "releases": len(stamps),
        "first": first.date(),
        "hist": hist,
        "latest": {s["name"]: s["latest"] for s in skills},
        "repo_mod": mod,
        "repo_lock": lock,
        "repo_pinned": tomllib.loads(mod)["pinned"],
        "personal": sum(s["personal"] for s in skills),
    }


# Real `rsk new --for all --no-color` run on 2026-11-04 in a scratch folder
# (/tmp/rsk-demo-go: go.mod, Dockerfile, cmd/demo/main.go) with RSK_CONFIG_HOME
# pointed at an empty temp dir. Paths shortened; one sandbox-only warning
# (no config.json, so the rsk-guide install was skipped) omitted.
RSK_NEW_OUT = [
    "✓  registered the rsk MCP server (.mcp.json, opencode.json)",
    "✓  initialized .rsk/ for all",
]
RSK_NEW_CLAUDE_MD = """@AGENTS.md
@.rsk/CLAUDE.md

<!-- rsk-mcp:start -->
## Project standards

Before scaffolding or a non-trivial change, call the rsk MCP `project_profile` tool to learn which skills apply here. Installed skills appear automatically in your skill list — invoke them with the Skill tool. Never install without asking.
<!-- rsk-mcp:end -->"""

# Real project_profile response for the same scratch folder (rsk MCP, 2026-11-04).
PROFILE = {
    "candidates": [
        {"because": ["language:go"], "installed": False, "latest": "1.14.0", "skill": "go-architect"},
        {"because": ["infra:docker"], "installed": False, "latest": "1.3.0", "skill": "docker-architect"},
        {"because": ["repo:cli"], "installed": False, "latest": "1.0.0", "skill": "cli-tool-architect"},
    ],
    "installed": [],
    "manifest": {"exists": False, "pinned": 0, "skills": 0},
    "project_root": "/tmp/rsk-demo-go",
    "signals": [
        {"kind": "language:go", "proof": "go.mod"},
        {"kind": "infra:docker", "proof": "Dockerfile"},
        {"kind": "repo:cli", "proof": "cmd/demo/main.go"},
    ],
}

# internal/mcp/signals.go — signal kind, the files that prove it, skills proposed.
SIGNALS = [
    ("language:go", "go.mod", "go-architect"),
    ("language:python", "pyproject.toml, requirements.txt, uv.lock", "python-architect"),
    ("language:typescript", "tsconfig.json, package.json", "—"),
    ("framework:nextjs", "next.config.*, or next in package.json", "nextjs-architect, react-architect"),
    ("framework:astro", "astro.config.*", "astro-architect"),
    ("protocol:grpc", "any .proto, or google.golang.org/grpc in go.mod", "grpc-architect, protobuf-architect"),
    ("protocol:rest", "openapi.yaml/.yml, swagger.json, any .http", "rest-api-architect"),
    ("infra:docker", "Dockerfile, docker-compose.yml/.yaml", "docker-architect"),
    ("infra:k8s", "Chart.yaml, kustomization.yaml/.yml", "—"),
    ("infra:ci", "files in .github/workflows", "ci-cd-architect"),
    ("repo:tooling", ".golangci.yml, mise.toml, Taskfile.yml", "repo-tooling-architect"),
    ("repo:cli", "main.go, or cmd/*/main.go", "cli-tool-architect"),
    ("data:postgres", ".sql files with sqlx or psycopg", "—"),
]

# rsk --help groups (cmd/rsk/setup.go, bin/rsk --help on 3.4.1)
HELP_GROUPS = [
    ("Machine setup", [("init", "Set up rsk for this machine.")]),
    ("Project lifecycle", [("new", "Initialize an rsk project in the current directory."), ("destroy", "Remove the rsk project from the current directory.")]),
    ("Install / uninstall / update", [("install", "Install bundles or skills."), ("uninstall", "Remove installed bundles or skills."), ("update", "Pull the latest skills and re-link.")]),
    ("Project pinning", [("pin", "Pin a skill so it is imported into each tool's project config.")]),
    ("Tool configuration", [("tools", "Manage tool permissions (allow/deny/list/remove).")]),
    ("Read-only views", [("catalog", "Browse the catalog of available skills and bundles."), ("list", "Show installed skills (project / --global)."), ("status", "Show installed skills and version drift.")]),
    ("Other", [("mcp", "Run the MCP server over stdio."), ("completion", "Generate the autocompletion script for the specified shell.")]),
]
assert sum(len(c) for _, c in HELP_GROUPS) == 13

TOOLS = [
    ("project_profile", "Detect this project's stack from its files and propose the ralvaskills skills that apply. Call it before scaffolding or a non-trivial change.", "ro"),
    ("search_skills", "Search the ralvaskills catalog by keyword. Use it to find a skill by name or intent when project_profile does not name one.", "ro"),
    ("install_skills", "Install one or more skills or bundles into this project by name, resolving them through the ralvaskills catalog. Writes symlinks and the project manifest. The client asks the user to approve.", "de"),
]

QUICK_START = [
    ("go install github.com/ralvarezdev/ralvaskills/v3/cmd/rsk@latest", "or a binary from Releases"),
    ("rsk init", "once per machine"),
    ("rsk new", "in a project: writes .rsk/ and registers rsk mcp"),
    ("rsk install go-grpc", "a bundle or a skill"),
    ("rsk pin go-architect", "load it in every session"),
]

SOURCES = [
    ("58 skills", "registry/index.json · skills/*/*/SKILL.md"),
    ("17 domains", "skills/<domain>/"),
    ("15 bundles", "internal/config/catalog.toml"),
    ("13 commands", "cmd/rsk/setup.go (+ cobra completion)"),
    ("3 MCP tools, 1 writes", "internal/mcp/tools.go"),
    ("rsk://catalog", "internal/mcp/resource.go"),
    ("Detection rules", "internal/mcp/signals.go"),
    ("Releases and dates", "registry/index.json versions[].published_at"),
    ("Version", "git tag v3.4.1 · bin/rsk --version"),
    ("Config and cache", "internal/config/config.go"),
    ("What install writes", "internal/install/apply.go"),
    ("rsk new output", "real run, scratch folder, 2026-11-04"),
    ("project_profile JSON", "real call, scratch folder, 2026-11-04"),
    ("Colours", "termkit TokyoNight via internal/ui/theme.go"),
    ("Wordmark", "internal/ui/theme.go"),
]

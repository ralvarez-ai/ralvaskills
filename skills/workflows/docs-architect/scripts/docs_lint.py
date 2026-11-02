#!/usr/bin/env python3
"""Validate a docs/ tree against the docs-architect convention.

Read-only. Stdlib only. Exit code 1 when any error is found.

  docs_lint.py [docs-dir] [--max-draft-days N] [--index]
"""

import argparse
import re
import sys
from datetime import date, datetime
from pathlib import Path

DATED_DIRS = {"ideas", "plans", "designs", "reviews", "changes"}
UNDATED_DIRS = {"guides", "methodology"}
STANDARD_STATUS = {"draft", "active", "done", "obsolete"}
ADR_STATUS = {"proposed", "accepted", "superseded", "deprecated"}
EXEMPT_DIRS = {"backlog", "templates"}
ADR_NAME = re.compile(r"^(\d{4})-[a-z0-9]+(-[a-z0-9]+)*\.md$")
DATED_NAME = re.compile(r"^\d{4}-\d{2}-\d{2}-")
LINK = re.compile(r"\[[^\]]*\]\(([^)\s]+)\)")
FENCE = re.compile(r"^```.*?^```", re.S | re.M)
BACKLOG_LIMITS = {"doing": 2, "next": 5}
BACKLOG_STATES = {"pending", "active", "done", "archive"}
PRIORITIES = {"p0", "p1", "p2", "p3"}
ITEM_NAME = re.compile(r"^([0-9a-hjkmnp-tv-z]{6})-[a-z0-9]+(-[a-z0-9]+)*\.md$")
ITEM_REF = re.compile(r"^([a-z0-9][a-z0-9-]*#)?[0-9a-hjkmnp-tv-z]{6}$")
ITEM_FORBIDDEN = {"status", "date", "created", "closed"}
ITEM_REQUIRED = ("title", "priority", "section")


def parse_frontmatter(text):
    """Return (dict, body) for a leading ---/--- block, or (None, text)."""
    if not text.startswith("---\n"):
        return None, text
    end = text.find("\n---", 4)
    if end == -1:
        return None, text
    meta = {}
    last = None
    for line in text[4:end].splitlines():
        if not line.strip() or line.lstrip().startswith("#"):
            last = None
            continue
        item = re.match(r"^\s*-\s+(.*)$", line)
        if item and last is not None:
            if not isinstance(meta[last], list):
                meta[last] = []
            meta[last].append(item.group(1).strip().strip("'\""))
            continue
        if ":" not in line:
            last = None
            continue
        key, _, value = line.partition(":")
        value = re.sub(r"\s+#.*$", "", value).strip()
        if value.startswith("[") and value.endswith("]"):
            inner = value[1:-1].strip()
            value = [v.strip().strip("'\"") for v in inner.split(",")] if inner else []
        elif value:
            value = value.strip("'\"")
        else:
            value = []
        last = key.strip()
        meta[last] = value
    return meta, text[end + 4 :]


def parse_date(value):
    try:
        return datetime.strptime(str(value), "%Y-%m-%d").date()
    except ValueError:
        return None


def title_of(body, fallback):
    for line in body.splitlines():
        if line.startswith("# "):
            return line[2:].strip()
    return fallback


class Linter:
    def __init__(self, root, max_draft_days):
        self.root = root
        self.max_draft_days = max_draft_days
        self.findings = []
        self.docs = {}
        self.inbound = set()

    def add(self, level, path, msg):
        self.findings.append((level, path.relative_to(self.root).as_posix(), msg))

    def top_dir(self, path):
        parts = path.relative_to(self.root).parts
        return parts[0] if len(parts) > 1 else ""

    def load(self):
        for path in sorted(self.root.rglob("*.md")):
            text = path.read_text(encoding="utf-8")
            meta, body = parse_frontmatter(text)
            self.docs[path] = (meta, body)

    def check_links(self, path, meta, body):
        targets = [m for m in LINK.findall(FENCE.sub("", body))]
        related = (meta or {}).get("related", [])
        if isinstance(related, str):
            related = [related] if related else []
        for rel in related:
            if ITEM_REF.match(str(rel)):
                continue  # <id> or <project>#<id>: a backlog reference, not a path
            target = (self.root / rel).resolve()
            if not target.exists():
                self.add("error", path, f"related path not found: {rel}")
            else:
                self.inbound.add(target)
        for link in targets:
            if re.match(r"^([a-z]+:|#)", link):
                continue
            target = (path.parent / link.split("#")[0]).resolve()
            if not target.exists():
                self.add("error", path, f"broken link: {link}")
            else:
                self.inbound.add(target)

    def check_standard(self, path, meta, folder):
        for key in ("status", "project", "date"):
            if not meta.get(key):
                self.add("error", path, f"frontmatter missing '{key}'")
        status = meta.get("status")
        if status and status not in STANDARD_STATUS:
            self.add("error", path, f"status '{status}' not in {sorted(STANDARD_STATUS)}")
        if "related" not in meta:
            self.add("warn", path, "frontmatter missing 'related'")
        created = parse_date(meta.get("date", ""))
        if meta.get("date") and not created:
            self.add("error", path, f"date '{meta['date']}' is not YYYY-MM-DD")
        if status == "draft" and created and folder != "archive":
            age = (date.today() - created).days
            if age > self.max_draft_days:
                self.add("warn", path, f"draft for {age} days: promote or delete")
        if folder in DATED_DIRS and not DATED_NAME.match(path.name):
            self.add("warn", path, "filename should start with YYYY-MM-DD-")
        if folder in UNDATED_DIRS and DATED_NAME.match(path.name):
            self.add("warn", path, "living doc: filename should not carry a date")

    def check_adrs(self):
        by_number = {}
        for path, (meta, _) in self.docs.items():
            if self.top_dir(path) != "adr" or path.name.lower() == "readme.md":
                continue
            match = ADR_NAME.match(path.name)
            if not match:
                self.add("error", path, "ADR filename must be NNNN-slug.md (kebab-case)")
                continue
            if meta is None:
                continue
            number = int(match.group(1))
            if str(meta.get("adr", "")) != str(number):
                self.add("error", path, f"frontmatter adr '{meta.get('adr')}' != filename {number}")
            if meta.get("status") not in ADR_STATUS:
                self.add("error", path, f"status must be one of {sorted(ADR_STATUS)}")
            for key in ("date", "project"):
                if not meta.get(key):
                    self.add("error", path, f"frontmatter missing '{key}'")
            if number in by_number:
                self.add("error", path, f"duplicate ADR number {number}")
            by_number[number] = (path, meta)
        numbers = sorted(by_number)
        for missing in sorted(set(range(1, numbers[-1] + 1)) - set(numbers)) if numbers else []:
            self.add("warn", self.root / "adr" / ".", f"gap in ADR numbering: {missing} is missing")
        for number, (path, meta) in by_number.items():
            old = meta.get("supersedes")
            if old:
                other = by_number.get(int(old)) if str(old).isdigit() else None
                if not other:
                    self.add("error", path, f"supersedes unknown ADR {old}")
                elif str(other[1].get("superseded_by", "")) != str(number) or other[1].get("status") != "superseded":
                    self.add("error", other[0], f"must be superseded_by {number} with status superseded")
            new = meta.get("superseded_by")
            if new and meta.get("status") != "superseded":
                self.add("error", path, "superseded_by set but status is not superseded")
            if meta.get("status") == "superseded" and not new:
                self.add("error", path, "status superseded but superseded_by is empty")

    def check_backlog_project(self, path, meta):
        """project.md of the item layout: slug, name and a non-empty sections list."""
        if meta is None:
            self.add("error", path, "project.md is missing frontmatter")
            return
        for key in ("slug", "name"):
            if not meta.get(key):
                self.add("error", path, f"frontmatter missing '{key}'")
        sections = meta.get("sections")
        if not isinstance(sections, list) or not sections:
            self.add("error", path, "frontmatter needs a non-empty 'sections' list")

    def check_backlog_item(self, path, meta, state, sections):
        """One file per item: id in the name, schema in the frontmatter."""
        if not ITEM_NAME.match(path.name):
            self.add("error", path, "item filename must be <id>-<slug>.md (id: 6 chars of [0-9a-hjkmnp-tv-z])")
        if meta is None:
            self.add("error", path, "missing frontmatter")
            return
        for key in ITEM_REQUIRED:
            if not meta.get(key):
                self.add("error", path, f"frontmatter missing '{key}'")
        for key in ITEM_FORBIDDEN:
            if key in meta:
                self.add("error", path, f"frontmatter must not carry '{key}'")
        priority = str(meta.get("priority", ""))
        if priority and priority not in PRIORITIES:
            self.add("error", path, f"priority '{priority}' not in {sorted(PRIORITIES)}")
        section = str(meta.get("section", ""))
        if section and sections is not None and section not in sections:
            self.add("error", path, f"section '{section}' not in project.md sections {sections}")
        for key in ("blocked_by", "related"):
            value = meta.get(key)
            if value is None:
                continue
            refs = value if isinstance(value, list) else [value]
            for ref in refs:
                if ref and not ITEM_REF.match(str(ref)):
                    self.add("error", path, f"{key}: '{ref}' must be <id> or <project>#<id>")
        resolution = meta.get("resolution")
        if resolution and (resolution != "dropped" or state != "archive"):
            self.add("error", path, "'resolution: dropped' is the only allowed value, and only in archive/")

    def check_backlog(self, path, body):
        section = None
        counts = {}
        for line in body.splitlines():
            heading = re.match(r"^#{1,6}\s+(\w+)", line)
            if heading:
                section = heading.group(1).lower()
            elif section in BACKLOG_LIMITS and re.match(r"^- \[ \]", line):
                counts[section] = counts.get(section, 0) + 1
        for section, count in counts.items():
            if count > BACKLOG_LIMITS[section]:
                self.add("warn", path, f"{section.title()} has {count} items (limit {BACKLOG_LIMITS[section]})")

    def run(self):
        self.load()
        projects = {}
        for path, (meta, _) in self.docs.items():
            if path.parent.name == "backlog" and path.name == "project.md":
                projects[path.parent] = (meta or {}).get("sections")
        for path, (meta, body) in self.docs.items():
            folder = self.top_dir(path)
            if folder == "backlog":
                parts = path.relative_to(self.root).parts
                backlog_dir = path.parents[len(parts) - 2]
                sub = parts[1] if len(parts) > 2 else ""
                if backlog_dir in projects:
                    if not sub and path.name == "project.md":
                        self.check_backlog_project(path, meta)
                    elif sub in BACKLOG_STATES:
                        self.check_backlog_item(path, meta, sub, projects[backlog_dir])
                    elif not sub and path.name == "inbox.md":
                        pass
                    else:
                        where = sub or path.name
                        self.add("warn", path, f"unexpected entry in an item-layout backlog: {where}")
                elif path.parent.name != "done" and path.name != "inbox.md":
                    self.check_backlog(path, body)
                self.check_links(path, meta, body)
                continue
            if folder == "templates":
                continue
            if path.name.lower() == "readme.md":
                self.check_links(path, meta, body)
                continue
            if meta is None:
                self.add("error", path, "missing frontmatter")
                continue
            if folder != "adr":
                self.check_standard(path, meta, folder)
            self.check_links(path, meta, body)
        self.check_adrs()
        for path, (meta, _) in self.docs.items():
            folder = self.top_dir(path)
            if folder in EXEMPT_DIRS | {"archive", "adr", ""}:
                continue
            if path.name.lower() == "readme.md":
                continue
            if path.resolve() not in self.inbound:
                self.add("warn", path, "orphan: no doc links to it")

    def index(self):
        rows = {}
        for path, (meta, body) in self.docs.items():
            folder = self.top_dir(path)
            if folder in EXEMPT_DIRS or not meta or path.name.lower() == "readme.md":
                continue
            rel = path.relative_to(self.root).as_posix()
            rows.setdefault(folder, []).append(
                (str(meta.get("date", "")), rel, title_of(body, path.stem), str(meta.get("status", "")))
            )
        out = []
        for folder in sorted(rows):
            out += [f"## {folder}/", "", "| Date | Doc | Status |", "|---|---|---|"]
            for when, rel, title, status in sorted(rows[folder]):
                out.append(f"| {when} | [{title}]({rel}) | {status} |")
            out.append("")
        return "\n".join(out)


def main():
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("docs", nargs="?", default="docs")
    parser.add_argument("--max-draft-days", type=int, default=60)
    parser.add_argument("--index", action="store_true", help="print a generated index and exit")
    args = parser.parse_args()
    root = Path(args.docs).resolve()
    if not root.is_dir():
        print(f"not a directory: {root}", file=sys.stderr)
        return 2
    linter = Linter(root, args.max_draft_days)
    linter.run()
    if args.index:
        print(linter.index())
        return 0
    for level, path, msg in sorted(linter.findings):
        print(f"{level.upper():5} {path}: {msg}")
    errors = sum(1 for f in linter.findings if f[0] == "error")
    warns = len(linter.findings) - errors
    print(f"{len(linter.docs)} docs, {errors} errors, {warns} warnings", file=sys.stderr)
    return 1 if errors else 0


if __name__ == "__main__":
    sys.exit(main())

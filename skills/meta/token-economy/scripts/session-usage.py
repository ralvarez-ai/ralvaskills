#!/usr/bin/env python3
"""Token usage of one Claude Code session, from its transcripts. Read-only.

    session-usage.py                 latest session of the current directory's project
    session-usage.py <session-id>    that session (prefix allowed) in the current project
    session-usage.py <file.jsonl>    that transcript

Prints one row per transcript (main thread and each subagent): requests, uncached input,
cache reads, cache writes (5m / 1h TTL), output (and thinking), and the largest context
sent. Then the bytes each hook injected. Tokens only: multiply by your model's prices.
"""

import json
import os
import sys
from collections import Counter
from pathlib import Path


def project_dir() -> Path:
    root = Path(os.environ.get("CLAUDE_CONFIG_DIR", Path.home() / ".claude")) / "projects"
    return root / str(Path.cwd()).replace("/", "-").replace(".", "-")


def resolve(arg: str | None) -> Path:
    if arg and arg.endswith(".jsonl"):
        return Path(arg)
    sessions = sorted(project_dir().glob("*.jsonl"), key=lambda p: p.stat().st_mtime)
    if arg:
        sessions = [p for p in sessions if p.stem.startswith(arg)]
    if not sessions:
        sys.exit(f"no transcript found in {project_dir()}")
    return sessions[-1]


def scan(path: Path, hooks: Counter) -> dict:
    by_id: dict[str, dict] = {}
    for line in path.open(encoding="utf-8"):
        try:
            rec = json.loads(line)
        except json.JSONDecodeError:
            continue
        if rec.get("type") == "assistant":
            msg = rec.get("message") or {}
            if msg.get("usage"):
                by_id[msg.get("id") or str(len(by_id))] = msg["usage"]  # repeated per block
        elif rec.get("type") == "attachment":
            att = rec.get("attachment") or {}
            if att.get("type") == "hook_success":
                size = len(str(att.get("content") or "")) + len(str(att.get("stdout") or ""))
                hooks[att.get("hookName") or att.get("hookEvent") or "?"] += size
    row = Counter(reqs=len(by_id))
    for u in by_id.values():
        cc = u.get("cache_creation") or {}
        sent = u.get("input_tokens", 0) + u.get("cache_read_input_tokens", 0) + u.get("cache_creation_input_tokens", 0)
        row["input"] += u.get("input_tokens", 0)
        row["read"] += u.get("cache_read_input_tokens", 0)
        row["w5m"] += cc.get("ephemeral_5m_input_tokens", 0)
        row["w1h"] += cc.get("ephemeral_1h_input_tokens", 0)
        row["out"] += u.get("output_tokens", 0)
        row["think"] += (u.get("output_tokens_details") or {}).get("thinking_tokens", 0)
        row["ctx_max"] = max(row["ctx_max"], sent)
    return row


def fmt(n: int) -> str:
    return f"{n / 1e6:.2f}M" if n >= 1e6 else f"{n / 1e3:.0f}k" if n >= 1e3 else str(n)


def main() -> None:
    main_file = resolve(sys.argv[1] if len(sys.argv) > 1 else None)
    files = [main_file, *sorted((main_file.parent / main_file.stem).glob("**/*.jsonl"))]
    hooks: Counter = Counter()
    cols = ("reqs", "input", "read", "w5m", "w1h", "out", "think", "ctx_max")
    print(f"session {main_file.stem}")
    print(f"{'transcript':<24}" + "".join(f"{c:>9}" for c in cols))
    total: Counter = Counter()
    for f in files:
        row = scan(f, hooks)
        name = "main" if f == main_file else f.stem[:24]
        print(f"{name:<24}" + "".join(f"{fmt(row[c]):>9}" for c in cols))
        for c in cols:
            total[c] = max(total[c], row[c]) if c == "ctx_max" else total[c] + row[c]
    print(f"{'TOTAL':<24}" + "".join(f"{fmt(total[c]):>9}" for c in cols))
    if hooks:
        print("\nhook output entering the context (chars, ~4 per token):")
        for name, size in hooks.most_common(8):
            print(f"  {size:>9,}  {name}")


if __name__ == "__main__":
    main()

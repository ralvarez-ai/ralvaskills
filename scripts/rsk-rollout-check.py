#!/usr/bin/env python3 -I
"""Read-only roll-out check for rsk 3.10.0 (docs/plans/2026-11-22-consumer-rollout-rsk-3-10.md).

Walks a root (default ~/Dev), finds every project with .rsk/rsk.lock and prints:
  - the survey numbers (locks, official-skill projects, v2 projects, pins by version);
  - one line per FAIL: a v2 project whose mise.toml pin (in or above it) is older than 3.10.0.
Unpinned projects are counted, never failures. Exit 1 when any FAIL line was printed.
Use --list official|v2|old-pin to list the projects behind a number.
"""
import os
import re
import sys

MIN = (3, 10, 0)
SKIP = {"node_modules", ".git", "vendor", ".venv", "dist", "target", ".cache"}
PIN = re.compile(r'^\s*"?github:ralvarez-ai/ralvaskills"?\s*=\s*(?:"([^"]+)"|\{[^}]*version\s*=\s*"([^"]+)")', re.M)
V2 = re.compile(r"^version\s*=\s*2", re.M)
OFFICIAL = re.compile(r'^\s*source\s*=\s*"official"', re.M)


def ver(s):
    return tuple(int(x) for x in re.findall(r"\d+", s)[:3])


def pin_in(d):
    p = os.path.join(d, "mise.toml")
    if os.path.isfile(p):
        m = PIN.search(open(p).read())
        if m:
            return m.group(1) or m.group(2)
    return None


def pin_above(d, root):
    while d.startswith(root):
        v = pin_in(d)
        if v:
            return v
        d = os.path.dirname(d)
    return None


def read(p):
    return open(p).read() if os.path.isfile(p) else ""


def main():
    args = sys.argv[1:]
    listing = None
    if "--list" in args:
        i = args.index("--list")
        listing = args[i + 1]
        del args[i : i + 2]
    root = os.path.abspath(os.path.expanduser(args[0] if args else "~/Dev"))
    locks, pins = [], {}
    for dp, dn, fn in os.walk(root):
        dn[:] = [d for d in dn if d not in SKIP]
        if os.path.basename(dp) == ".rsk" and "rsk.lock" in fn:
            locks.append(os.path.dirname(dp))
        if "mise.toml" in fn:
            v = pin_in(dp)
            if v:
                pins.setdefault(v, []).append(dp)
    official, v2, fails, unpinned, old = [], [], [], [], []
    for p in sorted(locks):
        text = read(p + "/.rsk/rsk.lock") + "\n" + read(p + "/.rsk/rsk.mod")
        pin = pin_above(p, root)
        if OFFICIAL.search(text):
            official.append(p)
        if pin is None:
            unpinned.append(p)
        elif ver(pin) < MIN:
            old.append(p)
        if V2.search(text):
            v2.append(p)
            if pin is not None and ver(pin) < MIN:
                fails.append((p, pin))
    if listing:
        for p in {"official": official, "v2": v2, "old-pin": old}[listing]:
            print(p)
        return 0
    print(f"locks={len(locks)} official={len(official)} v2={len(v2)} unpinned={len(unpinned)} old-pin={len(old)}")
    print("pins " + " ".join(f"{k}:{len(v)}" for k, v in sorted(pins.items(), key=lambda kv: ver(kv[0]))))
    for p, pin in fails:
        print(f"FAIL v2 under pin {pin}: {p}")
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())

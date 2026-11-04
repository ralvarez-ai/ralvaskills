"""Build the rsk landing proposals: python3 build.py [a b c d ...]"""

from __future__ import annotations

import importlib
import sys

import data as D

PAGES = {"a": ("dir_a", "direction-1.html"), "b": ("dir_b", "direction-2.html"), "c": ("dir_c", "direction-3.html"), "d": ("combined", "direction-4-combined.html")}


def main() -> None:
    d = D.load()
    want = sys.argv[1:] or list(PAGES)
    for k in want:
        mod, name = PAGES[k]
        html = importlib.import_module(mod).build(d)
        out = D.OUT / name
        out.write_text(html)
        print(f"{name}: {len(html.encode()) / 1024:.1f} KB")


if __name__ == "__main__":
    main()

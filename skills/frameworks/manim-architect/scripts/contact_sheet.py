#!/usr/bin/env python3
"""One image with a frame from every scene, so a version can be reviewed without playing it.

    python3 contact_sheet.py video.mp4 video.sections.json            # -> video.sheet.png
    python3 contact_sheet.py video.mp4 video.sections.json out.png --at 0.85 --cols 3 --width 640
    python3 contact_sheet.py video.mp4 --every 5                       # no sections file: a frame every 5 s

Frame position is a fraction of each scene (default 0.85: the finished drawing, before the fade-out).
Run it for every version and keep the sheets next to the videos (`v1.sheet.png`, `v2.sheet.png`):
two sheets side by side show what a rework changed. Each tile is labelled with its scene id when
ffmpeg has drawtext; otherwise the order and ids are printed instead.
Standard library and ffmpeg only.
"""
from __future__ import annotations

import argparse
import json
import math
import subprocess
import sys
import tempfile
from pathlib import Path

FONTS = ["/usr/share/fonts/dejavu-sans-fonts/DejaVuSans.ttf", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
         "/System/Library/Fonts/Supplemental/Arial.ttf", "C:/Windows/Fonts/arial.ttf"]


def run(cmd: list[str]) -> subprocess.CompletedProcess:
    return subprocess.run(cmd, capture_output=True, text=True, check=False)


def duration(video: Path) -> float:
    out = run(["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", str(video)])
    return float(out.stdout.strip())


def moments(video: Path, sections: Path | None, at: float, every: float) -> list[tuple[str, float]]:
    if sections:
        rows = json.loads(sections.read_text(encoding="utf-8"))
        return [(r["id"], r["start"] + at * (r["end"] - r["start"])) for r in rows]
    total = duration(video)
    return [(f"{t:g}s", t) for t in (i * every for i in range(math.ceil(total / every))) if t < total]


def label_args(text: str, width: int) -> list[str]:
    font = next((f for f in FONTS if Path(f).exists()), None)
    face = f"fontfile={font}:" if font else ""  # no file: ffmpeg falls back to fontconfig's default
    safe = text.replace("\\", "").replace(":", " ").replace("'", "")
    return ["-vf", f"scale={width}:-1,drawtext={face}text='{safe}':x=8:y=8:fontsize=22:"
            "fontcolor=white:box=1:boxcolor=black@0.6:boxborderw=4"]


def frame(video: Path, t: float, target: Path, name: str, width: int) -> None:
    base = ["ffmpeg", "-v", "error", "-y", "-ss", f"{t:.2f}", "-i", str(video), "-frames:v", "1"]
    labelled = base + label_args(name, width) + [str(target)]
    if label_args(name, width) and run(labelled).returncode == 0:
        return
    plain = run(base + ["-vf", f"scale={width}:-1", str(target)])  # no drawtext in this ffmpeg
    if plain.returncode != 0:
        sys.exit(plain.stderr.strip() or f"ffmpeg failed at {t:.2f}s")


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("video", type=Path)
    p.add_argument("sections", nargs="?", type=Path, help="sections json: [{id,start,end}, ...]")
    p.add_argument("out", nargs="?", type=Path)
    p.add_argument("--at", type=float, default=0.85, help="position inside each scene, 0..1")
    p.add_argument("--cols", type=int, default=3)
    p.add_argument("--width", type=int, default=640, help="tile width in pixels")
    p.add_argument("--every", type=float, default=5.0, help="seconds between frames when there is no sections file")
    a = p.parse_args()
    if a.sections and a.sections.suffix not in {".json"}:   # `video out.png` form
        a.out, a.sections = a.sections, None
    out = a.out or a.video.with_suffix(".sheet.png")
    items = moments(a.video, a.sections, a.at, a.every)
    if not items:
        sys.exit("no scenes")
    with tempfile.TemporaryDirectory() as tmp:
        for i, (name, t) in enumerate(items):
            frame(a.video, t, Path(tmp) / f"{i:03d}.png", name, a.width)
        rows = math.ceil(len(items) / a.cols)
        tile = run(["ffmpeg", "-v", "error", "-y", "-framerate", "1", "-i", f"{tmp}/%03d.png",
                    "-vf", f"tile={a.cols}x{rows}:padding=6:color=0x202020", "-frames:v", "1", str(out)])
        if tile.returncode != 0:
            sys.exit(tile.stderr.strip())
    for i, (name, t) in enumerate(items):
        print(f"{i + 1:>2}  {name:<24} {t:7.2f}s")
    print(out)
    return 0


if __name__ == "__main__":
    sys.exit(main())

#!/usr/bin/env python3
"""Check a narration script against the narration-architect rules.

Read-only. Stdlib only. Exit code 1 when any error is found (or any warning with --strict).

  narration_lint.py SCRIPT.json [--audio-dir DIR] [--strict]

The script file is JSON:

  {
    "rate": 2.5, "leadIn": 0.5, "tail": 0.5,
    "limits": {"sentence": 20, "paragraph": 6},
    "glossary": [{"term": "scan cycle", "avoid": ["scan loop"], "say": "scan cycle"}],
    "scenes": [{"id": "intro", "seconds": 10, "narration": "One idea per sentence."}]
  }

Only the heuristics here are checked. There is no approved-word dictionary, so a clean run does not
mean the text is ASD-STE100 compliant. Confirm the limits against the official specification.
"""

from __future__ import annotations

import argparse
import json
import math
import re
import shutil
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path

DEFAULTS = {"rate": 2.5, "leadIn": 0.5, "tail": 0.5}
DEFAULT_LIMITS = {"sentence": 20, "paragraph": 6}
AUDIO_SUFFIXES = (".mp3", ".wav", ".m4a", ".ogg", ".flac")

SENTENCE_SPLIT = re.compile(r"(?<=[.!?])\s+")
CONTRACTION = re.compile(r"\b\w+(?:n't|'re|'ve|'ll|'d|'m)\b|\b(?:it|that|there|what|let|here|who)'s\b", re.I)
PASSIVE = re.compile(r"\b(?:is|are|was|were|be|been|being)\s+(?:\w+ly\s+)?\w+(?:ed|en)\b", re.I)
UNSPEAKABLE = re.compile(r"[→←↔≥≤≠≈`_/\\{}<>|=]")


@dataclass(frozen=True)
class Finding:
    level: str  # "error" or "warning"
    where: str
    message: str

    def __str__(self) -> str:
        return f"{self.level.upper():7} {self.where}: {self.message}"


def words(text: str) -> list[str]:
    """Whitespace-separated tokens that contain a letter or a digit."""
    return [token for token in text.split() if re.search(r"\w", token)]


def sentences(text: str) -> list[str]:
    return [part.strip() for part in SENTENCE_SPLIT.split(text.strip()) if part.strip()]


def budget_words(seconds: float, lead_in: float, tail: float, rate: float) -> int:
    return max(0, math.floor((seconds - lead_in - tail) * rate))


def audio_seconds(path: Path) -> float | None:
    """Length of an audio file in seconds, or None when ffprobe is missing or fails."""
    if shutil.which("ffprobe") is None:
        return None
    result = subprocess.run(
        ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", str(path)],
        capture_output=True,
        text=True,
        check=False,
    )
    try:
        return float(result.stdout.strip())
    except ValueError:
        return None


def clip(text: str, size: int = 60) -> str:
    return text if len(text) <= size else text[: size - 3] + "..."


def check_scene(scene: dict, cfg: dict, glossary: list[dict]) -> tuple[list[Finding], int, int]:
    """Findings, words used and words allowed for one scene."""
    scene_id = scene.get("id", "?")
    text = scene.get("narration", "")
    findings: list[Finding] = []
    limits = cfg["limits"]
    total = 0

    for p_index, paragraph in enumerate((p for p in re.split(r"\n\s*\n", text) if p.strip()), start=1):
        parts = sentences(paragraph)
        if len(parts) > limits["paragraph"]:
            findings.append(
                Finding("warning", f"scene {scene_id!r} paragraph {p_index}", f"{len(parts)} sentences (limit {limits['paragraph']})")
            )
        for s_index, sentence in enumerate(parts, start=1):
            where = f"scene {scene_id!r} sentence {s_index}"
            count = len(words(sentence))
            total += count
            if count > limits["sentence"]:
                findings.append(Finding("error", where, f"{count} words (limit {limits['sentence']}): {clip(sentence)!r}"))
            if CONTRACTION.search(sentence):
                findings.append(Finding("error", where, f"contraction: {CONTRACTION.search(sentence).group(0)!r}"))
            if UNSPEAKABLE.search(sentence):
                findings.append(Finding("error", where, f"symbol that cannot be spoken: {UNSPEAKABLE.search(sentence).group(0)!r}"))
            if PASSIVE.search(sentence):
                findings.append(Finding("warning", where, f"possible passive voice: {PASSIVE.search(sentence).group(0)!r}"))

    lowered = text.lower()
    for entry in glossary:
        for avoided in entry.get("avoid", []):
            if re.search(rf"\b{re.escape(avoided.lower())}\b", lowered):
                findings.append(
                    Finding("error", f"scene {scene_id!r}", f"uses {avoided!r}; the glossary term is {entry['term']!r}")
                )

    seconds = scene.get("seconds")
    allowed = 0
    if seconds is not None:
        allowed = budget_words(seconds, cfg["leadIn"], cfg["tail"], cfg["rate"])
        if total > allowed:
            findings.append(Finding("error", f"scene {scene_id!r}", f"{total} words, budget {allowed} for {seconds} s"))
    return findings, total, allowed


def check_audio(scene: dict, cfg: dict, audio_dir: Path) -> list[Finding]:
    scene_id = scene.get("id", "?")
    seconds = scene.get("seconds")
    match = next((audio_dir / f"{scene_id}{suffix}" for suffix in AUDIO_SUFFIXES if (audio_dir / f"{scene_id}{suffix}").exists()), None)
    if match is None:
        return [Finding("warning", f"scene {scene_id!r}", f"no audio file in {audio_dir}")]
    length = audio_seconds(match)
    if length is None:
        return [Finding("warning", f"scene {scene_id!r}", f"could not measure {match.name} (is ffprobe installed?)")]
    if seconds is not None and length > seconds - cfg["tail"]:
        return [Finding("error", f"scene {scene_id!r}", f"audio is {length:.2f} s, scene allows {seconds - cfg['tail']:.2f} s: cut words")]
    return []


def main() -> int:
    parser = argparse.ArgumentParser(description="Check a narration script.")
    parser.add_argument("script", type=Path)
    parser.add_argument("--audio-dir", type=Path, help="folder with one audio file per scene id")
    parser.add_argument("--strict", action="store_true", help="treat warnings as errors")
    args = parser.parse_args()

    data = json.loads(args.script.read_text())
    cfg = {**DEFAULTS, **{k: data[k] for k in DEFAULTS if k in data}}
    cfg["limits"] = {**DEFAULT_LIMITS, **data.get("limits", {})}
    glossary = data.get("glossary", [])

    findings: list[Finding] = []
    for scene in data.get("scenes", []):
        scene_findings, used, allowed = check_scene(scene, cfg, glossary)
        findings += scene_findings
        if args.audio_dir is not None:
            findings += check_audio(scene, cfg, args.audio_dir)
        print(f"scene {scene.get('id', '?')!r}: {used} words" + (f" of {allowed}" if allowed else ""))

    for finding in findings:
        print(finding)
    errors = sum(1 for f in findings if f.level == "error")
    warnings = len(findings) - errors
    print(f"{errors} error(s), {warnings} warning(s)")
    return 1 if errors or (args.strict and warnings) else 0


if __name__ == "__main__":
    sys.exit(main())

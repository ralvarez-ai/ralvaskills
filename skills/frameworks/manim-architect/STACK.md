# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| manim (Community) | 0.22.1 | Scenes, animations, rendering |
| Python | 3.13+ | Runtime; run with `uv` |
| ffmpeg | any recent | Encoding, `adelay`/`amix` audio muxing, `scripts/contact_sheet.py` (needs `drawtext` for labels) |
| LaTeX + dvisvgm | TeX Live | `MathTex` formulas |
| LaTeX packages | standalone, xkeyval, preview, amsmath, amssymb | Needed by the minimal formula template |
| Cairo and Pango dev files | system | Text rendering (`Text`) |

## Notes

- Pin `manim` exactly in `pyproject.toml` (`manim==0.22.1`) and commit the lock file.
- Without root: install the TeX pieces into the user's home with `tlmgr --usermode`, or use a TinyTeX install; with root on Fedora: `texlive-dvisvgm texlive-standalone texlive-xkeyval texlive-preview`.
- The default Manim TeX template loads many packages that a small install lacks. Use a minimal template (`amsmath`, `amssymb`) and add packages one at a time.
- Check a tool before claiming it is missing: `command -v latex dvisvgm ffmpeg`.

_Last reviewed: 2026-11-09_
_Skill version at last review: 1.0.0_

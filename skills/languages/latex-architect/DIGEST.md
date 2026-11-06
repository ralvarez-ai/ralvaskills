# latex-architect (digest)

LaTeX standards for documents you control: TeX Live 2026, pdfLaTeX default, latexmk as the only build entry point.

## Always-on rules
- Build only with `latexmk`; never hand-run `pdflatex` twice.
- Commit `.latexmkrc` at the repo root; build into `build/`; set `$pdf_mode` and `$bibtex_use` explicitly.
- Root `main.tex` holds preamble and `\input` list only; one file per section, no preamble in sections.
- `\input` over `\include` unless per-chapter `.aux` files are needed.
- `\includegraphics` without extension; relative forward-slash paths; vector figures for plots.
- Package order is load-bearing: fontenc/lmodern, geometry/setspace, amsmath then mathtools, graphicx/booktabs/siunitx, biblatex, hyperref second-to-last, cleveref last. `microtype` always.
- Semantic macros; `\newcommand` over `\def`.
- Cross-refs via `\cref`/`\Cref`; `\label` right after `\caption`; prefixes `fig: tab: eq: sec: alg: app:`.
- `siunitx` for every number with a unit; `booktabs` rules only; `~` before every ref and cite.
- Floats `[t]` or `[tb]`; caption above tables, below figures; one `\centering` per float.
- Display math via `\[ ... \]`, `equation` or `align`; starred forms for unreferenced equations.
- Default `biblatex` + `biber`, `backend=biber`; protect capitals in titles with braces; cite published versions.
- One sentence per line in source; never reflow paragraphs.
- Quality gates: chktex and latexindent; CI builds the PDF on every push.

## Hard constraints
- Never `$$`, `eqnarray`, `\hline`, `subfigure`, `\bf`/`\it`/`\rm`/`\sc`, `\centerline`, or `\sloppy` to fix overfull boxes.
- Zero tolerance for `Overfull \hbox` over ~5pt in the final pass; fix by rewording.
- Never `\vspace`/`\hspace` for layout; never absolute paths in `\includegraphics`.
- Never convert a venue's natbib template to biblatex; keep natbib and `bibtex` there.
- Never suppress chktex by disabling it; use a committed `.chktexrc`.

Invoke the full skill when: writing, structuring, building or reviewing a LaTeX paper, report or thesis, or setting up its latexmk, bibliography or CI config.

Full guide: invoke the latex-architect skill.

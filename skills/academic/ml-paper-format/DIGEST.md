# ml-paper-format digest

Formatting and structure for NeurIPS, ICML, ICLR and AAAI submissions: style files, anonymity, page limits, natbib, reproducibility, appendix, desk-reject checks. Sits on latex-architect.

## Always-on rules

- Read the current year's CFP and author kit before applying anything; where they disagree, the CFP wins.
- Use the official style file for the target year; never copy a preamble from a previous paper.
- Confirm before drafting: venue and year, track, stage (submission, camera-ready, preprint), and whether code and data will be released.
- Anonymize: no author names, affiliations, emails or acknowledgements in submission mode; cite own work in third person; anonymize links; strip PDF metadata; check figure provenance.
- Page limit applies to main text. Cut content to fit; the appendix is unlimited.
- Keep natbib and the venue .bst. Brace-protect capitals in .bib. Cite published versions, not arXiv.
- Skeleton: abstract (under ~200 words), introduction ending in a contributions list, related work, method, experiments (setup, results, ablations), limitations, conclusion, broader impact, references, appendix.
- Every abstract and intro claim maps to an experiment. Takeaway-led captions. Report variance and seed counts.
- Reproducibility checklist answered honestly; report full setup and compute cost; anonymized code at submission.
- AAAI is the outlier: read its kit in full; do not assume NeurIPS habits transfer.
- Main text must stand alone; appendix goes after references, sections lettered A, B.

## Hard constraints

- Never shrink the layout to fit: no smaller font, baselinestretch, margins, negative vspace, or small body text.
- Never fabricate a number, baseline or reference.
- Never edit the official style file.
- Never let a load-bearing claim depend on the appendix.

Invoke the full skill when: writing, formatting or submitting an ML conference paper; need the venue matrix, the pre-submission checklist or the skeleton recipes.

Full guide: invoke the ml-paper-format skill.

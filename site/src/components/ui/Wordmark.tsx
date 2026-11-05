const ROWS = [
  '██████╗  ███████╗ ██╗  ██╗',
  '██╔══██╗ ██╔════╝ ██║ ██╔╝',
  '██████╔╝ ███████╗ █████╔╝',
  '██╔══██╗ ╚════██║ ██╔═██╗',
  '██║  ██║ ███████║ ██║  ██╗',
  '╚═╝  ╚═╝ ╚══════╝ ╚═╝  ╚═╝',
];

/** The real rsk banner wordmark: one SVG text row per line, CLI row gradient. */
export function Wordmark() {
  return (
    <svg className="wordmark" viewBox="0 -16 312 134" role="img" aria-label="rsk">
      {ROWS.map((row, i) => (
        <text
          key={row}
          className={`wordmark__row-${i}`}
          x={0}
          y={i * 22}
          fontFamily="DejaVu Sans Mono,Menlo,Consolas,monospace"
          fontSize={20}
          textLength={row.trimEnd().length * 12}
          lengthAdjust="spacingAndGlyphs"
          xmlSpace="preserve"
        >
          {row.trimEnd()}
        </text>
      ))}
    </svg>
  );
}

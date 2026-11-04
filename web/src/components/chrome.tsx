import type { Theme } from '../lib/useMotion';
import { sources } from '../data/facts';
import { data } from '../data/types';
import { Badge, Mark } from './ui';

export interface NavItem {
  label: string;
  href: string;
  /** Hidden below 760 px. */
  optional?: boolean;
}

const WORDMARK = [
  '██████╗  ███████╗ ██╗  ██╗',
  '██╔══██╗ ██╔════╝ ██║ ██╔╝',
  '██████╔╝ ███████╗ █████╔╝',
  '██╔══██╗ ╚════██║ ██╔═██╗',
  '██║  ██║ ███████║ ██║  ██╗',
  '╚═╝  ╚═╝ ╚══════╝ ╚═╝  ╚═╝',
];

/** The real rsk banner wordmark, one SVG text row per line, CLI row gradient. */
export function Wordmark() {
  return (
    <svg className="k-wm" viewBox="0 -16 312 134" role="img" aria-label="rsk">
      {WORDMARK.map((row, i) => (
        <text
          key={i}
          className={`wm${i}`}
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

export function Header({
  nav,
  theme,
  mounted,
  onToggle,
}: {
  nav: NavItem[];
  theme: Theme;
  mounted: boolean;
  onToggle: () => void;
}) {
  return (
    <>
      <a className="skip" href="#main">
        Skip to content
      </a>
      <header className="k-top">
        <div className="wrap">
          <a className="k-brand" href="#top">
            <Mark />rsk <small>{data.version}</small>
          </a>
          <nav className="k-nav" aria-label="Page">
            {nav.map((item) => (
              <a key={item.href} href={item.href} data-opt={item.optional ? '' : undefined}>
                {item.label}
              </a>
            ))}
            {mounted ? (
            <button
              className="k-btn ghost k-theme"
              type="button"
              aria-pressed={theme === 'dark'}
              onClick={onToggle}
            >
              <span className="ico" aria-hidden="true">
                ◐
              </span>
              <span className="t">Dark theme</span>
            </button>
            ) : null}
          </nav>
          <div className="k-prog" aria-hidden="true">
            <i />
          </div>
        </div>
      </header>
    </>
  );
}

export function Rail({ items, active }: { items: { id: string; label: string }[]; active: string }) {
  return (
    <nav className="k-rail" aria-label="Sections">
      <ol>
        {items.map((item, i) => (
          <li key={item.id}>
            <a href={`#${item.id}`} aria-current={active === item.id ? 'location' : undefined}>
              <span className="n">{String(i + 1).padStart(2, '0')}</span>
              <span className="lb">{item.label}</span>
            </a>
          </li>
        ))}
      </ol>
    </nav>
  );
}

export function Footer() {
  return (
    <footer className="k-foot">
      <div className="wrap">
        <div>
          <Wordmark />
          <p style={{ marginTop: 14 }}>
            rsk {data.version} · GPL-3.0 ·{' '}
            <a href="https://github.com/ralvarezdev/ralvaskills">github.com/ralvarezdev/ralvaskills</a>
          </p>
        </div>
        <details>
          <summary>Where every number on this page comes from</summary>
          <p>
            Counted in the repository on 2026-11-04. Anything marked <Badge kind="sa">sample</Badge> is illustrative.
          </p>
          <ul className="k-src">
            {sources.map((s) => (
              <li key={s.what}>
                <span>{s.what}</span>
                <code>{s.where}</code>
              </li>
            ))}
          </ul>
        </details>
      </div>
    </footer>
  );
}

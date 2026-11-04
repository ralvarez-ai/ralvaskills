import { forwardRef, useState } from 'react';
import type { CSSProperties, ReactNode } from 'react';

import type { FilmState } from '../lib/useFilm';

export type BadgeKind = 'ro' | 'de' | 'hu' | 'sa';

/** The brand-kit square and dot, the family mark. */
export function Mark() {
  return (
    <span className="fd" aria-hidden="true">
      <i />
      <i />
    </span>
  );
}

export function Kicker({ children, reveal = true, className }: { children: ReactNode; reveal?: boolean; className?: string }) {
  return (
    <p className={['k-kick', reveal ? 'rv' : '', className ?? ''].filter(Boolean).join(' ')}>
      <Mark /> {children}
    </p>
  );
}

export function Badge({ kind = 'ro', children }: { kind?: BadgeKind; children: ReactNode }) {
  return <span className={`k-badge ${kind}`}>{children}</span>;
}

export function Btn({
  href,
  primary,
  children,
  style,
}: {
  href: string;
  primary?: boolean;
  children: ReactNode;
  style?: CSSProperties;
}) {
  return (
    <a className={`k-btn${primary ? ' primary' : ''}`} href={href} style={style}>
      {children}
    </a>
  );
}

export function CodeWindow({
  title,
  badge,
  children,
  className,
  bodyClassName,
}: {
  title: ReactNode;
  badge?: ReactNode;
  children: ReactNode;
  className?: string;
  bodyClassName?: string;
}) {
  return (
    <div className={`k-win${className ? ` ${className}` : ''}`}>
      <div className="k-win-h">
        {title}
        {badge ? <span style={{ marginLeft: 'auto' }}>{badge}</span> : null}
      </div>
      <div className={bodyClassName}>{children}</div>
    </div>
  );
}

export const FilmButton = forwardRef<HTMLButtonElement, { state: FilmState; onToggle: () => void }>(
  function FilmButton({ state, onToggle }, ref) {
    const label = state === 'playing' ? 'Pause' : state === 'paused' ? 'Play' : 'Replay';
    const icon = state === 'playing' ? '❚❚' : state === 'paused' ? '▶' : '↻';
    return (
      <button ref={ref} className="k-film-ctl" type="button" data-state={state} onClick={onToggle}>
        <span className="ic" aria-hidden="true">
          {icon}
        </span>
        <span className="l">{label}</span>
        <span className="bar" aria-hidden="true">
          <i />
        </span>
        <span className="sr"> the hero animation</span>
      </button>
    );
  },
);

/** A copy block: comments above each command, commands only to the clipboard. */
export function CopyBlock({ id, lines, label }: { id: string; lines: { cmd: string; note?: string }[]; label: string }) {
  const [status, setStatus] = useState<'idle' | 'ok' | 'fail'>('idle');
  const raw = lines
    .map((l) => l.cmd)
    .filter(Boolean)
    .join('\n');

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(raw);
      setStatus('ok');
    } catch {
      setStatus('fail');
    }
    setTimeout(() => setStatus('idle'), 1800);
  };

  return (
    <div className="k-copy">
      <pre id={id} data-raw={raw} aria-label={label} role="region" tabIndex={0}>
        {lines.map((line, i) => (
          <span key={i}>
            {line.note ? <span className="k-cm">{`# ${line.note}`}</span> : null}
            {line.cmd ? <span className="k-c">{line.cmd}</span> : null}
            {i < lines.length - 1 ? '\n' : null}
          </span>
        ))}
      </pre>
      <button className="k-btn" type="button" onClick={copy}>
        <span className="l">{status === 'ok' ? 'Copied' : status === 'fail' ? 'Copy failed' : 'Copy'}</span>
        <span className="sr"> {label}</span>
      </button>
      <p className="sr" aria-live="polite">
        {status === 'ok' ? 'Commands copied to the clipboard' : status === 'fail' ? 'Copy failed; select the text instead' : ''}
      </p>
    </div>
  );
}

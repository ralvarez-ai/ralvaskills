import { useState } from 'react';

import { Button } from './Button';

/** A copy block: notes above each command, commands only to the clipboard. */
export function CopyBlock({
  id,
  lines,
  label,
}: {
  id: string;
  lines: { cmd: string; note?: string }[];
  label: string;
}) {
  const [status, setStatus] = useState<'idle' | 'ok' | 'fail'>('idle');
  const raw = lines
    .map((line) => line.cmd)
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
    <div className="copy-block">
      <pre id={id} data-raw={raw} aria-label={label} role="region" tabIndex={0}>
        {lines.map((line, i) => (
          <span key={i}>
            {line.note ? <span className="copy-block__comment">{`# ${line.note}`}</span> : null}
            {line.cmd ? <span className="copy-block__command">{line.cmd}</span> : null}
            {i < lines.length - 1 ? '\n' : null}
          </span>
        ))}
      </pre>
      <Button onClick={copy}>
        <span>{status === 'ok' ? 'Copied' : status === 'fail' ? 'Copy failed' : 'Copy'}</span>
        <span className="sr-only"> {label}</span>
      </Button>
      <p className="sr-only" aria-live="polite">
        {status === 'ok' ? 'Commands copied to the clipboard' : status === 'fail' ? 'Copy failed; select the text instead' : ''}
      </p>
    </div>
  );
}

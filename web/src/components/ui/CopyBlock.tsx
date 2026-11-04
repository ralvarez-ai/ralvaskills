import { useState } from 'react';

import { Button } from './Button';
import { CodeRegion } from './CodeRegion';

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
  const rows = lines.flatMap((line) => [
    ...(line.note ? [{ kind: 'note' as const, text: `# ${line.note}` }] : []),
    ...(line.cmd ? [{ kind: 'command' as const, text: line.cmd }] : []),
  ]);
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
      <CodeRegion id={id} data-raw={raw} label={label}>
        {rows.map((row, i) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: a static list of lines that never reorders, and two lines can have the same text
          <span key={i}>
            <span className={row.kind === 'note' ? 'copy-block__comment' : 'copy-block__command'}>{row.text}</span>
            {i < rows.length - 1 ? '\n' : null}
          </span>
        ))}
      </CodeRegion>
      <Button onClick={copy}>
        <span>{status === 'ok' ? 'Copied' : status === 'fail' ? 'Copy failed' : 'Copy'}</span>
        <span className="sr-only"> {label}</span>
      </Button>
      <p className="sr-only" aria-live="polite">
        {status === 'ok'
          ? 'Commands copied to the clipboard'
          : status === 'fail'
            ? 'Copy failed; select the text instead'
            : ''}
      </p>
    </div>
  );
}

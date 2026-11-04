import type { ReactNode } from 'react';

import { cx } from './cx';

export type TagTone = 'readonly' | 'destructive' | 'hint' | 'sample';

/** A small stamped label: read-only, writes, sample, from the code. */
export function Tag({ tone = 'readonly', className, children }: { tone?: TagTone; className?: string; children: ReactNode }) {
  return (
    <span className={cx('tag', className)} data-tone={tone}>
      {children}
    </span>
  );
}

import type { CSSProperties, ReactNode } from 'react';

import { cx } from './cx';

export type SectionTone = 'paper' | 'sunken' | 'ink';

/** A page section. The tone is the band colour, not a one-off style prop. */
export function Section({
  id,
  labelledBy,
  tone = 'paper',
  className,
  style,
  children,
}: {
  id: string;
  labelledBy?: string;
  tone?: SectionTone;
  className?: string;
  style?: CSSProperties;
  children: ReactNode;
}) {
  return (
    <section
      className={cx('section', className)}
      id={id}
      aria-labelledby={labelledBy}
      data-tone={tone === 'paper' ? undefined : tone}
      style={style}
    >
      {children}
    </section>
  );
}

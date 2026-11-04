import type { CSSProperties, ReactNode } from 'react';
import { BrandMark } from './BrandMark';
import { cx } from './cx';

/** The mono, uppercase kicker above a section title. */
export function Eyebrow({
  children,
  reveal = true,
  mark = false,
  className,
  style,
}: {
  children: ReactNode;
  reveal?: boolean;
  mark?: boolean;
  className?: string;
  style?: CSSProperties;
}) {
  return (
    <p className={cx('eyebrow', reveal && 'reveal', className)} style={style}>
      {mark ? <BrandMark /> : null}
      {children}
    </p>
  );
}

import type { CSSProperties, ReactNode } from 'react';

import { cx } from './cx';

/** A bordered surface with a title bar: code windows, file blocks, tables. */
export function Panel({
  title,
  action,
  className,
  bodyClassName,
  bodyStyle,
  flush = false,
  children,
}: {
  title: ReactNode;
  action?: ReactNode;
  className?: string;
  bodyClassName?: string;
  bodyStyle?: CSSProperties;
  /** The child is the body (a pre, a list): do not add a padding wrapper. */
  flush?: boolean;
  children: ReactNode;
}) {
  return (
    <div className={cx('panel', className)}>
      <div className="panel__header">
        {title}
        {action ? <span style={{ marginLeft: 'auto' }}>{action}</span> : null}
      </div>
      {flush ? children : <div className={bodyClassName} style={bodyStyle}>{children}</div>}
    </div>
  );
}

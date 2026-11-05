import type { ButtonHTMLAttributes, ReactNode } from 'react';

import { cx } from './cx';

export type ButtonVariant = 'primary' | 'ghost';

interface Props {
  variant?: ButtonVariant;
  className?: string;
  children: ReactNode;
}

/** A link styled as a button (the page's CTAs are all links). */
export function ButtonLink({ href, variant, className, children }: Props & { href: string }) {
  return (
    <a className={cx('button', className)} href={href} data-variant={variant}>
      {children}
    </a>
  );
}

/** A real button, for anything with behaviour. */
export function Button({
  variant,
  className,
  children,
  type = 'button',
  ...rest
}: Props & ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button className={cx('button', className)} type={type} data-variant={variant} {...rest}>
      {children}
    </button>
  );
}

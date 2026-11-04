import type { Theme } from '../lib/useMotion';
import { data } from '../data/types';
import { BrandMark, Button } from './ui';

export interface NavItem {
  label: string;
  href: string;
  /** Hidden below 760 px. */
  optional?: boolean;
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
  onToggle: (event: React.MouseEvent<HTMLButtonElement>) => void;
}) {
  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <header className="site-header">
        <div className="container">
          <a className="brand" href="#top">
            <BrandMark />rsk <small>{data.version}</small>
          </a>
          <nav className="site-nav" aria-label="Page">
            {nav.map((item) => (
              <a key={item.href} href={item.href} data-opt={item.optional ? '' : undefined}>
                {item.label}
              </a>
            ))}
            {mounted ? (
              <Button variant="ghost" className="theme-toggle" aria-pressed={theme === 'dark'} onClick={onToggle}>
                <span className="theme-toggle__icon" aria-hidden="true">
                  ◐
                </span>
                <span className="theme-toggle__label">Dark theme</span>
              </Button>
            ) : null}
          </nav>
          <div className="reading-progress" aria-hidden="true">
            <i />
          </div>
        </div>
      </header>
    </>
  );
}

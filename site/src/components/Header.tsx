import { data } from '../data/types';
import type { Theme } from '../lib/useMotion';
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
  active,
  onToggle,
}: {
  nav: NavItem[];
  /** Id of the section in view; the nav link pointing at it is marked current. */
  active: string;
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
            <BrandMark />
            rsk <small>{data.version}</small>
          </a>
          <nav className="site-nav" aria-label="Page">
            {nav.map((item) => (
              <a
                key={item.href}
                href={item.href}
                data-opt={item.optional ? '' : undefined}
                aria-current={active === item.href.slice(1) ? 'location' : undefined}
              >
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

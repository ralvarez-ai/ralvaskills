import { useCallback, useState } from 'react';

import { Footer, Header, Rail } from './components/chrome';
import type { NavItem } from './components/chrome';
import { useCounts, useReveals, useScrollKit, useTheme } from './lib/useMotion';
import { Catalogue } from './sections/Catalogue';
import { Commands } from './sections/Commands';
import { Faq } from './sections/Faq';
import { Files } from './sections/Files';
import { Hero } from './sections/Hero';
import { QuickStart } from './sections/QuickStart';
import { Signals } from './sections/Signals';
import { ThisRepo } from './sections/ThisRepo';
import { Told } from './sections/Told';
import { Tools } from './sections/Tools';
import { Trade } from './sections/Trade';
import { Writes } from './sections/Writes';

const NAV: NavItem[] = [
  { label: 'How it works', href: '#top', optional: true },
  { label: 'Catalogue', href: '#catalogue', optional: true },
  { label: 'Trade-off', href: '#trade', optional: true },
  { label: 'Install', href: '#start' },
];

const SECTIONS = [
  { id: 'top', label: 'The loop' },
  { id: 'writes', label: 'What it writes' },
  { id: 'files', label: 'rsk new' },
  { id: 'catalogue', label: 'Catalogue' },
  { id: 'tools', label: 'Tools' },
  { id: 'profile', label: 'Detection' },
  { id: 'this-repo', label: 'Real file' },
  { id: 'commands', label: 'Commands' },
  { id: 'trade', label: 'Trade-off' },
  { id: 'start', label: 'Install' },
  { id: 'faq', label: 'FAQ' },
];

export default function App() {
  const { theme, toggle, mounted } = useTheme();
  const [active, setActive] = useState('top');
  const onActive = useCallback((id: string) => setActive(id), []);

  useReveals();
  useCounts();
  useScrollKit(onActive);

  return (
    <>
      <Header nav={NAV} theme={theme} mounted={mounted} onToggle={toggle} />
      <Rail items={SECTIONS} active={active} />
      <main id="main">
        <Hero />
        <Writes />
        <Files />
        <Catalogue />
        <Tools />
        <Signals />
        <Told />
        <ThisRepo />
        <Commands />
        <Trade />
        <QuickStart />
        <Faq />
      </main>
      <Footer />
    </>
  );
}

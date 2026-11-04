/**
 * Interactive states, mockup vs build: rest, hover, keyboard focus and (for the
 * FAQ) open, captured per element. Writes PNGs and the element's computed
 * style; scripts/compare-states.py diffs them.
 *
 *   pnpm build && node scripts/compare-states.mjs [outDir]
 */
import { mkdirSync, readdirSync, readFileSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
import { chromium } from 'playwright-core';
import { preview } from 'vite';

const out = resolve(process.argv[2] ?? '/tmp/cmp-states');
mkdirSync(out, { recursive: true });
const base = join(homedir(), '.cache/ms-playwright');
const exe =
  process.env.CHROME ??
  readdirSync(base)
    .filter((d) => /^chromium-\d+$/.test(d))
    .sort()
    .reverse()
    .map((d) => join(base, d, 'chrome-linux64/chrome'))
    .find((p) => {
      try {
        readFileSync(p);
        return true;
      } catch {
        return false;
      }
    });

// [name, mock selector, web selector, index]
const TARGETS = [
  ['btn-primary', '.k-btn.primary', '.button[data-variant="primary"]', 0],
  ['btn-secondary', 'a.k-btn:not(.primary)', 'a.button:not([data-variant="primary"])', 0],
  ['btn-theme', '.k-theme', '.theme-toggle', 0],
  ['film-ctl', '#c-film .k-film-ctl', '#c-film .film-control', 0],
  ['chip', '#b-chips button', '.filter__chip', 3],
  ['chip-toggle', '#b-inc', '.filter__chip--toggle', 0],
  ['faq', '#faq summary', '#faq summary', 1],
  ['nav-link', '.k-nav a', '.site-nav a', 1],
  ['brand', '.k-top a[href="#top"]', '.site-header .brand', 0],
  ['copy', '[data-copy]', '.copy-block button, [class*=copy] button', 0],
  ['skill-tag', '.b-tag', '.skill', 5],
];
const STATES = ['rest', 'hover', 'focus', 'pressed'];

const server = await preview({ logLevel: 'error', preview: { port: 4332, strictPort: true } });
const browser = await chromium.launch({ executablePath: exe });
const URLS = {
  mock: pathToFileURL(resolve(`../mockups/ui-proposals/2026-11-04/${process.env.MOCK ?? 'direction-4-combined.html'}`))
    .href,
  web: server.resolvedUrls.local[0],
};
const styleProps = [
  'color',
  'backgroundColor',
  'borderTopColor',
  'borderTopWidth',
  'outlineStyle',
  'outlineWidth',
  'outlineColor',
  'outlineOffset',
  'boxShadow',
  'transform',
  'opacity',
  'textDecorationLine',
  'cursor',
  'paddingLeft',
  'height',
];

for (const theme of ['light', 'dark']) {
  for (const [who, url] of Object.entries(URLS)) {
    const ctx = await browser.newContext({
      viewport: { width: 1440, height: 900 },
      colorScheme: theme,
      reducedMotion: 'reduce',
    });
    await ctx.addInitScript((t) => localStorage.setItem('rsk-theme', t), theme);
    const page = await ctx.newPage();
    await page.goto(url, { waitUntil: 'networkidle' });
    await page.addStyleTag({
      content:
        'html{scroll-behavior:auto!important}*{transition-duration:0s!important;animation-duration:0s!important}',
    });
    const styles = {};
    for (const [name, mockSel, webSel, idx] of TARGETS) {
      const sel = who === 'mock' ? mockSel : webSel;
      const loc = page.locator(sel).nth(idx);
      if (!(await loc.count())) {
        styles[name] = 'MISSING';
        continue;
      }
      // Header controls are sticky: capture them at the top, or whatever
      // scrolls behind the translucent bar is part of the diff.
      if (['nav-link', 'brand', 'btn-theme'].includes(name)) await page.evaluate(() => scrollTo(0, 0));
      else await loc.evaluate((e) => e.scrollIntoView({ block: 'center' }));
      for (const state of STATES) {
        await page.mouse.move(2, 2);
        await page.evaluate(() => document.activeElement?.blur());
        if (name === 'faq' && state === 'pressed')
          await loc.click(); // open
        else if (state === 'pressed') {
          if (!name.startsWith('chip') && name !== 'film-ctl') continue;
        }
        if (state === 'hover') await loc.hover();
        if (state === 'focus') {
          await page.keyboard.press('Shift');
          await loc.focus();
        }
        if (state === 'pressed' && name !== 'faq') {
          await page.mouse.move(2, 2);
        } // pressed chips show aria-pressed rest style
        await page.waitForTimeout(60);
        const box = await loc.boundingBox();
        const pad = 12;
        await page.screenshot({
          path: join(out, `${theme}-${name}-${state}-${who}.png`),
          clip: {
            x: Math.max(0, box.x - pad),
            y: Math.max(0, box.y - pad),
            width: Math.min(box.width + pad * 2, 520),
            height: Math.min(box.height + pad * 2, 140),
          },
        });
        styles[name] ??= {};
        styles[name][state] = await loc.evaluate(
          (e, props) => Object.fromEntries(props.map((p) => [p, getComputedStyle(e)[p]])),
          styleProps,
        );
        if (name === 'faq' && state === 'pressed') await loc.click(); // close again
      }
    }
    writeFileSync(join(out, `${theme}-${who}.json`), JSON.stringify(styles));
    await ctx.close();
  }
}
await browser.close();
await new Promise((r) => server.httpServer.close(r));
console.log(out);

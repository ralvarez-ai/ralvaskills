/**
 * The deployment gate, run against the built site.
 *
 *   pnpm build && node scripts/verify.mjs
 *
 * Serves dist/ with `vite preview`, then for 390x844 and 1440x800 x light/dark
 * x motion/reduced: runs axe-core (wcag2a/2aa/21a/21aa/22aa + best-practice),
 * checks horizontal overflow, and fails if any reveal never fired or any
 * count-up did not settle. Also renders with JavaScript off, which the
 * prerender makes meaningful.
 */
import { readdirSync, readFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { join } from 'node:path';

import { chromium } from 'playwright-core';
import { preview } from 'vite';

const AXE = readFileSync(new URL('../node_modules/axe-core/axe.min.js', import.meta.url), 'utf8');
const TAGS = ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa', 'best-practice'];
const VIEWPORTS = { 390: { width: 390, height: 844 }, 1440: { width: 1440, height: 800 } };

// Rules that must actually reach the page. A pruned or renamed selector can
// leave the markup intact while the styles silently stop applying (the mono
// stack and the theme toggle both did), and axe cannot see that.
const WIRED = [
  ['code', 'fontFamily', 'JetBrains Mono'],
  ['.agent-hero__line > span', 'transitionDuration', '0.95s'],
  ['.button', 'transitionDuration', '0.16s'],
  ['.theme-toggle', 'paddingLeft', '12px'],
  ['.film-control', 'display', 'flex'],
  ['.panel', 'borderRadius', '6px'],
  ['.rule', 'backgroundColor', 'rgb'],
  // Optional 4th item: minimum viewport width. The rail once sat under a later
  // `display:none` and never showed, and no other check noticed.
  ['.section-rail', 'display', 'block', 1400],
  ['.section-rail', 'position', 'fixed', 1400],
];

function chromePath() {
  if (process.env.CHROME) return process.env.CHROME;
  const base = join(homedir(), '.cache/ms-playwright');
  for (const dir of readdirSync(base)
    .filter((d) => /^chromium-\d+$/.test(d))
    .sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))) {
    const candidate = join(base, dir, 'chrome-linux64/chrome');
    try {
      readFileSync(candidate);
      return candidate;
    } catch {
      // keep looking
    }
  }
  return undefined;
}

async function settle(page) {
  await page.addStyleTag({ content: 'html{scroll-behavior:auto!important}' });
  const height = await page.evaluate(() => document.documentElement.scrollHeight);
  for (let y = 0; y < height; y += 400) {
    await page.evaluate((to) => window.scrollTo(0, to), y);
    await page.waitForTimeout(40);
  }
  // Films started by the scroll-through must finish first: axe measures
  // contrast, and a step mid-fade is partly transparent.
  await page.waitForFunction(() => !document.querySelector('.film-control[data-state="playing"]'), null, {
    timeout: 12000,
  });
  await page.waitForTimeout(1200);
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.waitForTimeout(150);
}

const failures = [];
const fail = (msg) => {
  failures.push(msg);
  console.log(`FAIL ${msg}`);
};

const server = await preview({ logLevel: 'error', preview: { port: 4321, strictPort: true } });
const base = server.resolvedUrls?.local?.[0] ?? 'http://localhost:4173/';
const browser = await chromium.launch({ executablePath: chromePath() });

try {
  for (const [width, viewport] of Object.entries(VIEWPORTS)) {
    for (const theme of ['light', 'dark']) {
      for (const reduced of [false, true]) {
        const ctx = await browser.newContext({
          viewport,
          deviceScaleFactor: 1,
          colorScheme: theme,
          reducedMotion: reduced ? 'reduce' : 'no-preference',
        });
        await ctx.addInitScript((t) => {
          try {
            localStorage.setItem('rsk-theme', t);
          } catch {
            // ignore
          }
        }, theme);
        const page = await ctx.newPage();
        const errors = [];
        page.on('pageerror', (e) => errors.push(String(e)));
        page.on('console', (m) => {
          if (m.type() === 'error') errors.push(m.text());
        });
        await page.goto(base, { waitUntil: 'networkidle' });
        await settle(page);

        const run = `${width} ${theme} ${reduced ? 'reduced' : 'motion'}`;
        await page.addScriptTag({ content: AXE });
        const violations = await page.evaluate(async (tags) => {
          const result = await window.axe.run(document, { runOnly: { type: 'tag', values: tags } });
          return result.violations.map((v) => ({
            id: v.id,
            nodes: v.nodes.length,
            targets: v.nodes.slice(0, 3).map((n) => String(n.target)),
          }));
        }, TAGS);
        const layout = await page.evaluate(() => ({
          scroll: document.documentElement.scrollWidth,
          client: document.documentElement.clientWidth,
          hidden: document.querySelectorAll('.reveal:not(.is-visible), .rule:not(.is-visible)').length,
          counts: [...document.querySelectorAll('[data-count]')].filter((el) => el.textContent !== el.dataset.count)
            .length,
        }));

        const wired = reduced
          ? []
          : await page.evaluate((rules) => {
              const out = [];
              for (const [sel, prop, want, minWidth] of rules) {
                if (minWidth && innerWidth < minWidth) continue;
                const el = document.querySelector(sel);
                if (!el) {
                  out.push(`${sel}: missing`);
                  continue;
                }
                const got = String(getComputedStyle(el)[prop]);
                const bad = want === 'rgb' ? got === 'rgba(0, 0, 0, 0)' : !got.includes(want);
                if (bad) out.push(`${sel} ${prop}: ${got} (want ${want})`);
              }
              return out;
            }, WIRED);
        const ok = !violations.length && layout.scroll <= layout.client && !errors.length && !wired.length;
        if (!ok)
          fail(
            `${run}: ${violations.map((v) => `${v.id}(${v.nodes})`).join(',') || 'ok'} overflow=${layout.scroll}>${layout.client} errors=${errors.length}`,
          );
        else console.log(`PASS ${run}`);
        for (const v of violations) fail(`  ${v.id}: ${v.targets.join(' | ')}`);
        for (const w of wired) fail(`${run}: ${w}`);
        if (reduced === false && layout.hidden) fail(`${run}: ${layout.hidden} reveal(s) never fired`);
        if (layout.counts) fail(`${run}: ${layout.counts} count-up(s) did not settle`);
        await ctx.close();
      }
    }
  }

  // Films: each figure must start when scrolled into view, hide steps while it
  // runs, and end with every step shown. A figure whose ref never attached
  // exited the engine silently and sat in its final state; no other check sees it.
  for (const reduced of [false, true]) {
    const ctx = await browser.newContext({
      viewport: VIEWPORTS[1440],
      reducedMotion: reduced ? 'reduce' : 'no-preference',
    });
    const page = await ctx.newPage();
    await page.goto(base, { waitUntil: 'networkidle' });
    await page.addStyleTag({ content: 'html{scroll-behavior:auto!important}' });
    for (const id of ['c-film', 'a-film']) {
      const run = `film ${id} ${reduced ? 'reduced' : 'motion'}`;
      await page.evaluate((i) => document.getElementById(i)?.scrollIntoView({ block: 'center' }), id);
      await page.waitForTimeout(600);
      const mid = await page.evaluate((i) => {
        const r = document.getElementById(i);
        return {
          state: r?.querySelector('.film-control')?.dataset.state,
          hidden: r?.querySelectorAll('.is-hidden').length ?? 0,
        };
      }, id);
      if (reduced) {
        if (mid.state !== 'done' || mid.hidden) fail(`${run}: state=${mid.state} hidden=${mid.hidden} (want done, 0)`);
        else console.log(`PASS ${run}`);
        continue;
      }
      if (mid.state !== 'playing' || !mid.hidden) {
        fail(`${run}: did not start (state=${mid.state} hidden=${mid.hidden})`);
        continue;
      }
      const ended = await page
        .waitForFunction(
          (i) => document.getElementById(i).querySelector('.film-control')?.dataset.state === 'done',
          id,
          { timeout: 9000 },
        )
        .then(
          () => true,
          () => false,
        );
      const end = await page.evaluate((i) => {
        const r = document.getElementById(i);
        return {
          hidden: r.querySelectorAll('.is-hidden').length,
          typed: [...r.querySelectorAll('[data-t]')].filter((e) => e.textContent !== e.dataset.t).length,
          pin: r.querySelector('#a-pin')?.textContent ?? 'n/a',
        };
      }, id);
      if (!ended || end.hidden || end.typed || (id === 'a-film' && !/pinned = \["go-architect"\]/.test(end.pin))) {
        fail(`${run}: ended=${ended} hidden=${end.hidden} untyped=${end.typed} pin=${end.pin.trim()}`);
      } else console.log(`PASS ${run}`);
    }
    // Current section: header nav and rail follow the scroll position.
    await page.evaluate(() => document.getElementById('catalogue')?.scrollIntoView());
    await page.waitForTimeout(500);
    const cur = await page.evaluate(() => ({
      nav: document.querySelector('.site-nav a[aria-current]')?.textContent?.trim(),
      rail: document.querySelector('.section-rail a[aria-current] .section-rail__label')?.textContent?.trim(),
    }));
    if (cur.nav !== 'Catalogue' || cur.rail !== 'Catalogue')
      fail(`current ${reduced ? 'reduced' : 'motion'}: nav=${cur.nav} rail=${cur.rail} (want Catalogue)`);
    else console.log(`PASS current section ${reduced ? 'reduced' : 'motion'}`);
    await ctx.close();
  }

  // JavaScript off: the prerendered document must read on its own.
  for (const [width, viewport] of Object.entries(VIEWPORTS)) {
    const ctx = await browser.newContext({ viewport, javaScriptEnabled: false, deviceScaleFactor: 1 });
    const page = await ctx.newPage();
    await page.goto(base, { waitUntil: 'load' });
    const r = await page.evaluate(() => ({
      h1: document.querySelector('h1')?.innerText.trim() ?? '',
      sections: document.querySelectorAll('main section').length,
      text: document.querySelector('main')?.innerText.length ?? 0,
      scroll: document.documentElement.scrollWidth,
      client: document.documentElement.clientWidth,
    }));
    if (!r.h1 || r.text < 6000 || r.scroll > r.client)
      fail(`${width} js-off: h1="${r.h1.slice(0, 30)}" text=${r.text} overflow=${r.scroll}>${r.client}`);
    else console.log(`PASS ${width} js-off (${r.sections} sections, ${r.text} chars)`);
    await ctx.close();
  }
} finally {
  await browser.close();
  await new Promise((resolve) => server.httpServer.close(resolve));
}

if (failures.length) {
  console.log(`\n${failures.length} FAILED`);
  process.exit(1);
}
console.log('\nALL CLEAN');

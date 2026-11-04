/**
 * Interactive behaviour, mockup vs build: filter chips, rail highlight, theme
 * toggle, count-ups, reading progress, copy button.
 *   pnpm build && node scripts/compare-behaviour.mjs
 */
import { readdirSync, readFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
import { chromium } from 'playwright-core';
import { preview } from 'vite';

const base = join(homedir(), '.cache/ms-playwright');
const exe = readdirSync(base)
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
const server = await preview({ logLevel: 'error', preview: { port: 4328, strictPort: true } });
const browser = await chromium.launch({ executablePath: exe });
const T = {
  mock: {
    url: pathToFileURL(resolve('../mockups/ui-proposals/2026-11-04/direction-4-combined.html')).href,
    chips: '#b-chips button',
    tags: '.b-tag',
    on: '.b-tag.on',
    note: '#b-note',
    cnt: '#b-count',
    cntl: '#b-count-l',
    rail: '.k-rail a',
    theme: '.k-theme',
    prog: '.k-prog i',
    copy: '[data-copy]',
    copyLbl: '.l',
  },
  web: {
    url: server.resolvedUrls.local[0],
    chips: '.filter__chip',
    tags: '.skill',
    on: '.skill.is-match',
    note: '.catalogue__note',
    cnt: '.filter-count',
    cntl: '.filter-count__label',
    rail: '.section-rail a',
    theme: '.theme-toggle',
    prog: '.reading-progress i',
    copy: '.copy-block button, [class*=copy] button',
    copyLbl: 'span',
  },
};
const out = {};
for (const [who, c] of Object.entries(T)) {
  const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 }, colorScheme: 'light' });
  await ctx.grantPermissions(['clipboard-read', 'clipboard-write']);
  const page = await ctx.newPage();
  const errs = [];
  page.on('pageerror', (e) => errs.push(e.message));
  await page.goto(c.url, { waitUntil: 'networkidle' });
  await page.addStyleTag({ content: 'html{scroll-behavior:auto!important}' });
  const r = (out[who] = {});
  r.errors = errs;

  // filter chips
  await page.evaluate(() => document.getElementById('catalogue').scrollIntoView());
  const chips = page.locator(c.chips);
  const n = await chips.count();
  r.chipCount = n;
  r.filter = [];
  for (let i = 0; i < n; i += 1) {
    await chips.nth(i).click();
    await page.waitForTimeout(80);
    r.filter.push(
      await page.evaluate(
        (c2) => ({
          label:
            document.querySelectorAll(c2.chips)[0] && [...document.querySelectorAll(c2.chips)].map(() => '').length,
          pressed: [...document.querySelectorAll(c2.chips)]
            .map((b) => (b.getAttribute('aria-pressed') === 'true' ? 1 : 0))
            .join(''),
          matched: document.querySelectorAll(c2.on).length,
          visibleTags: [...document.querySelectorAll(c2.tags)].filter((t) => t.offsetParent !== null).length,
          count: (document.querySelector(c2.cnt)?.textContent ?? '').replace(/\s+/g, ' ').trim(),
          note: (document.querySelector(c2.note)?.textContent ?? '').replace(/\s+/g, ' ').trim().slice(0, 70),
        }),
        c,
      ),
    );
  }

  // rail highlight
  r.rail = [];
  const ids = await page.$$eval(c.rail, (as) => as.map((a) => a.getAttribute('href').slice(1)));
  for (const id of ids) {
    await page.evaluate((i) => document.getElementById(i).scrollIntoView(), id);
    await page.waitForTimeout(250);
    r.rail.push(
      `${id}→` +
        (await page.$$eval(c.rail, (as) =>
          as
            .filter((a) => a.hasAttribute('aria-current'))
            .map((a) => a.textContent.trim().slice(0, 14))
            .join('|'),
        )),
    );
  }

  // reading progress at mid page
  await page.evaluate(() => scrollTo(0, (document.documentElement.scrollHeight - innerHeight) / 2));
  await page.waitForTimeout(200);
  r.progressMid = await page.$eval(
    c.prog,
    (e) =>
      getComputedStyle(e).getPropertyValue('--p') || e.style.getPropertyValue('--p') || getComputedStyle(e).transform,
  );

  // count-ups: first frame after scrolling a counter into view
  await page.evaluate(() => scrollTo(0, 0));
  await page.waitForTimeout(200);
  r.counts = [];
  const nCounts = await page.locator('[data-count]').count();
  for (let i = 0; i < nCounts; i += 1) {
    const loc = page.locator('[data-count]').nth(i);
    await loc.scrollIntoViewIfNeeded();
    await page.waitForTimeout(250);
    const mid = (await loc.textContent()).trim();
    await page.waitForTimeout(1500);
    r.counts.push(`${mid}→${(await loc.textContent()).trim()}`);
  }

  // theme toggle
  await page.evaluate(() => scrollTo(0, 0));
  const before = await page.evaluate(() => document.documentElement.dataset.theme);
  await page.locator(c.theme).first().click();
  await page.waitForTimeout(800);
  r.theme = await page
    .evaluate(
      (t) => ({
        from: document.documentElement.dataset.theme,
        bg: getComputedStyle(document.body).backgroundColor,
        pressed: document.querySelector(t).getAttribute('aria-pressed'),
        stored: localStorage.getItem('rsk-theme'),
      }),
      c.theme,
    )
    .then((x) => ({ ...x, before }));

  // copy
  const copy = page.locator(c.copy).first();
  await copy.scrollIntoViewIfNeeded();
  await copy.click();
  await page.waitForTimeout(300);
  r.copy = {
    label: (await copy.textContent()).trim(),
    clip: (await page.evaluate(() => navigator.clipboard.readText())).slice(0, 40),
  };
  await ctx.close();
}
await browser.close();
await new Promise((res) => server.httpServer.close(res));

const keys = new Set([...Object.keys(out.mock), ...Object.keys(out.web)]);
let bad = 0;
for (const k of keys) {
  const a = JSON.stringify(out.mock[k]);
  const b = JSON.stringify(out.web[k]);
  if (a === b) {
    console.log(`= ${k}`);
    continue;
  }
  bad += 1;
  console.log(`≠ ${k}\n   mock ${a?.slice(0, 900)}\n   web  ${b?.slice(0, 900)}`);
}
process.exitCode = bad ? 1 : 0;

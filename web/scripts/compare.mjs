/**
 * Side-by-side capture of the built site and the combined mockup, for review.
 *
 *   pnpm build && node scripts/compare.mjs [outDir]
 *
 * Writes, per width x theme: <name>-{mock,web}.png (full page), and
 * <name>-{mock,web}.json (a dump of every text-bearing element: section, text,
 * rect, computed typography/colour). scripts/compare.py turns these into
 * side-by-sides and a diff report.
 */
import { mkdirSync, readdirSync, readFileSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

import { chromium } from 'playwright-core';
import { preview } from 'vite';

function chromePath() {
  if (process.env.CHROME) return process.env.CHROME;
  const base = join(homedir(), '.cache/ms-playwright');
  return readdirSync(base)
    .filter((d) => /^chromium-\d+$/.test(d))
    .sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))
    .map((d) => join(base, d, 'chrome-linux64/chrome'))
    .find((p) => {
      try {
        readFileSync(p);
        return true;
      } catch {
        return false;
      }
    });
}

const out = resolve(process.argv[2] ?? '/tmp/cmp');
mkdirSync(out, { recursive: true });
const mockFile = resolve('../mockups/ui-proposals/2026-11-04/direction-4-combined.html');

// WEB_URL points the "web" side at a running server (e.g. `pnpm dev`) instead of the build.
const server = process.env.WEB_URL
  ? null
  : await preview({ logLevel: 'error', preview: { port: 4324, strictPort: true } });
const browser = await chromium.launch({ executablePath: chromePath() });
const targets = {
  mock: pathToFileURL(mockFile).href,
  web: process.env.WEB_URL ?? server.resolvedUrls.local[0],
};

/** Runs in the page: dump text-bearing elements with their section and style. */
function dump() {
  const props = [
    'fontFamily',
    'fontSize',
    'fontWeight',
    'letterSpacing',
    'lineHeight',
    'textTransform',
    'color',
    'backgroundColor',
    'borderTopColor',
    'borderTopWidth',
    'borderRadius',
    'paddingTop',
    'paddingLeft',
    'display',
    'opacity',
  ];
  const rows = [];
  const counts = new Map();
  for (const el of document.body.querySelectorAll('*')) {
    if (['SCRIPT', 'STYLE', 'SVG', 'PATH'].includes(el.tagName.toUpperCase())) continue;
    const own = [...el.childNodes]
      .filter((n) => n.nodeType === 3)
      .map((n) => n.textContent.replace(/\s+/g, ' ').trim())
      .join(' ')
      .trim();
    if (!own) continue;
    const r = el.getBoundingClientRect();
    if (r.width === 0 && r.height === 0) continue;
    const sec = el.closest('section[id], header, footer, nav');
    const secId = sec ? sec.id || sec.tagName.toLowerCase() + '.' + sec.className.split(' ')[0] : '';
    const key = `${secId}|${own.slice(0, 60)}`;
    const n = counts.get(key) ?? 0;
    counts.set(key, n + 1);
    const cs = getComputedStyle(el);
    const style = {};
    for (const p of props) style[p] = cs[p];
    rows.push({
      key: `${key}#${n}`,
      tag: el.tagName.toLowerCase(),
      cls: el.className && el.className.baseVal === undefined ? String(el.className) : '',
      x: Math.round(r.left + scrollX),
      y: Math.round(r.top + scrollY),
      w: Math.round(r.width),
      h: Math.round(r.height),
      style,
    });
  }
  const sections = [...document.querySelectorAll('main > section, header, footer')].map((s) => {
    const r = s.getBoundingClientRect();
    return { id: s.id || s.tagName.toLowerCase(), y: Math.round(r.top + scrollY), h: Math.round(r.height) };
  });
  return { rows, sections, height: document.documentElement.scrollHeight };
}

const widths = (process.env.WIDTHS ?? '1440,390').split(',').map(Number);
const themes = (process.env.THEMES ?? 'light,dark').split(',');
for (const width of widths) {
  for (const theme of themes) {
    for (const [who, url] of Object.entries(targets)) {
      const ctx = await browser.newContext({
        viewport: { width, height: width === 390 ? 844 : 800 },
        deviceScaleFactor: 1,
        colorScheme: theme,
      });
      await ctx.addInitScript((t) => localStorage.setItem('rsk-theme', t), theme);
      const page = await ctx.newPage();
      await page.goto(url, { waitUntil: 'networkidle' });
      await page.waitForTimeout(7500);
      await page.addStyleTag({ content: 'html{scroll-behavior:auto!important}' });
      const height = await page.evaluate(() => document.documentElement.scrollHeight);
      for (let y = 0; y < height; y += 400) {
        await page.evaluate((to) => window.scrollTo(0, to), y);
        await page.waitForTimeout(60);
      }
      // Let everything the scroll-through started finish (films, reveals, rule
      // draws), or the capture depends on timing and two runs disagree.
      await page.waitForFunction(() => !document.querySelector('[data-state="playing"]'), null, { timeout: 15000 });
      await page.waitForTimeout(1500);
      await page.evaluate(() => window.scrollTo(0, 0));
      await page.waitForFunction(() => window.scrollY === 0);
      await page.waitForTimeout(600);
      const name = join(out, `${width}-${theme}-${who}`);
      await page.screenshot({ path: `${name}.png`, fullPage: true });
      writeFileSync(`${name}.json`, JSON.stringify(await page.evaluate(dump)));
      console.log(name);
      await ctx.close();
    }
  }
}

await browser.close();
if (server) await new Promise((r) => server.httpServer.close(r));

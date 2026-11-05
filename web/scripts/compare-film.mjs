/**
 * Film behaviour, mockup vs build: scroll each figure into view, then sample its
 * button label and the count of visible [data-f] steps at fixed times.
 *   node scripts/compare-film.mjs
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
const server = await preview({ logLevel: 'error', preview: { port: 4325, strictPort: true } });
const browser = await chromium.launch({ executablePath: exe });
const targets = {
  mock: pathToFileURL(resolve(`../mockups/ui-proposals/2026-11-04/${process.env.MOCK ?? 'direction-4-combined.html'}`))
    .href,
  web: server.resolvedUrls.local[0],
};
const times = [0, 700, 1500, 2500, 3500, 4500, 5500, 6500, 7500];
for (const id of ['c-film', 'a-film', 'b-film']) {
  const res = {};
  for (const [who, url] of Object.entries(targets)) {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 800 } });
    const page = await ctx.newPage();
    await page.goto(url, { waitUntil: 'networkidle' });
    await page.addStyleTag({ content: 'html{scroll-behavior:auto!important}' });
    await page.evaluate((i) => document.getElementById(i)?.scrollIntoView({ block: 'center' }), id);
    const t0 = Date.now();
    const rows = [];
    for (const t of times) {
      const wait = t - (Date.now() - t0);
      if (wait > 0) await page.waitForTimeout(wait);
      rows.push(
        await page.evaluate((i) => {
          const r = document.getElementById(i);
          if (!r) return 'MISSING';
          const hidden = r.querySelectorAll('.f-h,.is-hidden').length;
          const total = r.querySelectorAll('[data-f]').length;
          const btn = r.querySelector('.k-film-ctl,[data-film-control],button[class*=film]');
          const typed = [...r.querySelectorAll('[data-t]')].map((e) => e.textContent.length).join(',');
          return `${(btn?.textContent ?? '?').trim().padEnd(8)} shown ${total - hidden}/${total} typed[${typed}]`;
        }, id),
      );
    }
    res[who] = rows;
    await ctx.close();
  }
  console.log(`\n== ${id}`);
  times.forEach((t, i) => {
    const same = res.mock[i] === res.web[i] ? ' ' : '≠';
    console.log(`${same} t=${String(t).padStart(4)}  mock: ${res.mock[i]}\n            web : ${res.web[i]}`);
  });
}
await browser.close();
await new Promise((r) => server.httpServer.close(r));

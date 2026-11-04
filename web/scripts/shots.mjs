/**
 * Screenshots of the built site, for review: settled first viewport and full
 * page, both themes at 1440 and 390.
 *
 *   pnpm build && node scripts/shots.mjs [outDir]
 */
import { mkdirSync, readFileSync, readdirSync } from 'node:fs';
import { homedir } from 'node:os';
import { join, resolve } from 'node:path';

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

const out = resolve(process.argv[2] ?? 'shots');
mkdirSync(out, { recursive: true });

const server = await preview({ logLevel: 'error', preview: { port: 4323, strictPort: true } });
const browser = await chromium.launch({ executablePath: chromePath() });
const base = server.resolvedUrls.local[0];

for (const width of [1440, 390]) {
  for (const theme of ['light', 'dark']) {
    const ctx = await browser.newContext({
      viewport: { width, height: width === 390 ? 844 : 800 },
      deviceScaleFactor: 1,
      colorScheme: theme,
    });
    await ctx.addInitScript((t) => localStorage.setItem('rsk-theme', t), theme);
    const page = await ctx.newPage();
    await page.goto(base, { waitUntil: 'networkidle' });
    await page.waitForTimeout(7500);
    // Scroll through once so every reveal fires; a fullPage shot does not scroll.
    const height = await page.evaluate(() => document.documentElement.scrollHeight);
    for (let y = 0; y < height; y += 500) {
      await page.evaluate((to) => window.scrollTo(0, to), y);
      await page.waitForTimeout(60);
    }
    await page.evaluate(() => window.scrollTo(0, 0));
    await page.waitForTimeout(400);
    const name = join(out, `${width}-${theme}`);
    await page.screenshot({ path: `${name}.png` });
    await page.screenshot({ path: `${name}-full.png`, fullPage: true });
    console.log(`${name}.png`);
    await ctx.close();
  }
}

await browser.close();
await new Promise((r) => server.httpServer.close(r));

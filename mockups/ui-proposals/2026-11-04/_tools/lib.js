// Shared helpers for verify/smoke/shots. Uses the Chromium already installed
// for Playwright on this machine; override with CHROME=/path/to/chrome.
const fs = require('fs');
const os = require('os');
const path = require('path');
const { chromium } = require('playwright');

const DIR = path.resolve(__dirname, '..');
const AXE = fs.readFileSync(require.resolve('axe-core/axe.min.js'), 'utf8');

function chromePath() {
  if (process.env.CHROME) return process.env.CHROME;
  const base = path.join(os.homedir(), '.cache/ms-playwright');
  const dirs = fs.readdirSync(base).filter((d) => /^chromium-\d+$/.test(d)).sort((a, b) => +b.split('-')[1] - +a.split('-')[1]);
  for (const d of dirs) {
    const p = path.join(base, d, 'chrome-linux64/chrome');
    if (fs.existsSync(p)) return p;
  }
  return undefined;
}

async function launch() {
  return chromium.launch({ executablePath: chromePath() });
}

const VIEWPORTS = { 390: { width: 390, height: 844 }, 1440: { width: 1440, height: 800 } };

async function open(browser, file, { width, theme, reduced, js = true, blockImages = false }) {
  const ctx = await browser.newContext({
    viewport: VIEWPORTS[width],
    deviceScaleFactor: 1,
    reducedMotion: reduced ? 'reduce' : 'no-preference',
    colorScheme: theme,
    javaScriptEnabled: js,
  });
  if (js) await ctx.addInitScript((t) => { try { localStorage.setItem('rsk-theme', t); } catch (e) {} }, theme);
  if (blockImages) await ctx.route('**/*.{png,jpg,jpeg,webp,avif,gif,svg}', (r) => r.abort());
  const page = await ctx.newPage();
  const errors = [];
  page.on('pageerror', (e) => errors.push(String(e)));
  page.on('console', (m) => { if (m.type() === 'error') errors.push('console: ' + m.text()); });
  await page.goto('file://' + path.join(DIR, file), { waitUntil: 'load' });
  return { ctx, page, errors };
}

// Scroll the whole page so one-shot reveals fire, then wait for every film to settle.
async function settle(page, { reduced }) {
  if (!reduced) {
    await page.evaluate(async () => {
      document.documentElement.style.scrollBehavior = 'auto';
      const step = innerHeight * 0.5;
      for (let y = 0; y < document.documentElement.scrollHeight; y += step) {
        scrollTo(0, y);
        await new Promise((r) => setTimeout(r, 90));
      }
      scrollTo(0, document.documentElement.scrollHeight);
      await new Promise((r) => setTimeout(r, 200));
      scrollTo(0, 0);
      document.documentElement.style.scrollBehavior = '';
    });
    await page.waitForFunction(() => [...document.querySelectorAll('.k-film-ctl')].every((b) => b.dataset.state === 'done'), null, { timeout: 15000 }).catch(() => {});
    // films only start when visible: force any still idle ones to their final state
    await page.evaluate(() => document.querySelectorAll('[data-film]').forEach((r) => { if (r.kFilm && r.kFilm.state() !== 'done') r.kFilm.finish(); }));
    // count-ups are view-triggered too; wait, then force whatever is left
    await page.waitForFunction(() => [...document.querySelectorAll('[data-count]')].every((e) => e.dataset.done === '1'), null, { timeout: 6000 }).catch(() => {});
    await page.evaluate(() => window.kCountsDone && window.kCountsDone());
    await page.waitForTimeout(1100);
  } else {
    await page.waitForTimeout(250);
  }
}

async function axe(page) {
  await page.addScriptTag({ content: AXE });
  return page.evaluate(async () => {
    const r = await window.axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa', 'best-practice'] }, resultTypes: ['violations'] });
    return r.violations.map((v) => ({ id: v.id, impact: v.impact, n: v.nodes.length, help: v.help, targets: v.nodes.slice(0, 4).map((n) => n.target.join(' ') + ' :: ' + (n.failureSummary || '').split('\n').slice(0, 2).join(' ')) }));
  });
}

// Content left hidden by a reveal would be skipped by axe: count it as a failure.
async function unrevealed(page) {
  return page.evaluate(() => [...document.querySelectorAll('.rv:not(.in)')].length + (document.documentElement.classList.contains('anim') ? 0 : 0));
}

async function overflow(page) {
  return page.evaluate(() => {
    const w = document.documentElement.clientWidth;
    const sw = document.documentElement.scrollWidth;
    const culprits = [];
    if (sw > w) document.querySelectorAll('body *').forEach((e) => { const r = e.getBoundingClientRect(); if (r.right > w + 1 && culprits.length < 6 && getComputedStyle(e).position !== 'fixed') culprits.push(e.tagName + '.' + e.className + ' ' + Math.round(r.right)); });
    return { overflow: sw > w, sw, w, culprits };
  });
}

module.exports = { DIR, launch, open, settle, axe, overflow, unrevealed, VIEWPORTS };

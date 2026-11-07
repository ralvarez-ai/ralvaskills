// Renders public/og.png (1200x630), the social card, in the Teal & Sand palette.
// Run by hand when the palette or the tagline changes: `node scripts/og.mjs`.
// Uses the same Chrome lookup as verify.mjs (CHROME, then Playwright's cache).
import { readdirSync } from 'node:fs';
import { homedir } from 'node:os';
import { join } from 'node:path';
import { chromium } from 'playwright-core';

function chromePath() {
  if (process.env.CHROME) return process.env.CHROME;
  const base = join(homedir(), '.cache/ms-playwright');
  const dir = readdirSync(base)
    .filter((d) => d.startsWith('chromium-'))
    .sort()
    .pop();
  return join(base, dir, 'chrome-linux64/chrome');
}

const html = `<!doctype html><meta charset="utf-8"><style>
*{margin:0;box-sizing:border-box}
body{width:1200px;height:630px;background:#f6f1e7;color:#1a1d1c;font:600 64px/1.1 system-ui,-apple-system,"Segoe UI",Roboto,sans-serif;padding:84px 96px;display:flex;flex-direction:column;justify-content:space-between;border-top:14px solid #0b6e6a}
.mark{font:700 44px ui-monospace,"JetBrains Mono",Menlo,monospace;color:#0b6e6a;letter-spacing:-.02em}
h1{font-size:76px;letter-spacing:-.02em;max-width:980px}
.you{color:#a83e16}
p{font:400 30px/1.4 system-ui,sans-serif;color:#4f564e}
</style><div class="mark">rsk</div>
<h1>Your agent asks.<br><span class="you">You approve.</span><br>rsk installs.</h1>
<p>Versioned skills for Claude Code and OpenCode · rsk.ralvarez.dev</p>`;

const browser = await chromium.launch({ executablePath: chromePath() });
const page = await browser.newPage({ viewport: { width: 1200, height: 630 } });
await page.setContent(html);
await page.screenshot({ path: new URL('../public/og.png', import.meta.url).pathname });
await browser.close();

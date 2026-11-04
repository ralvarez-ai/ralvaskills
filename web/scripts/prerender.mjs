/**
 * Prerender the app into dist/index.html after `vite build`.
 *
 * The page is a single document that must read without JavaScript: search
 * engines, link previews and a JS-off reader should see the copy, the tables
 * and the code blocks, not an empty #root. Rendering once with react-dom/server
 * gives that for free, and main.tsx hydrates it.
 *
 * Effects do not run here, so the films, reveals and theme toggle stay client
 * concerns; the markup is already the films' final state.
 */
import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import { createElement } from 'react';
import { renderToString } from 'react-dom/server';
import { createServer } from 'vite';

const WEB = resolve(dirname(fileURLToPath(import.meta.url)), '..');

const server = await createServer({
  root: WEB,
  appType: 'custom',
  logLevel: 'error',
  server: { middlewareMode: true },
});

try {
  const { default: App } = await server.ssrLoadModule('/src/App.tsx');
  const markup = renderToString(createElement(App));
  const file = resolve(WEB, 'dist/index.html');
  const doc = readFileSync(file, 'utf8');
  if (!doc.includes('<div id="root"></div>')) throw new Error('dist/index.html has no empty #root to fill');
  writeFileSync(file, doc.replace('<div id="root"></div>', `<div id="root">${markup}</div>`));
  console.log(`prerendered ${(markup.length / 1024).toFixed(1)} KB of markup into dist/index.html`);
} finally {
  await server.close();
}

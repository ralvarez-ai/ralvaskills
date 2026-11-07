import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// The landing page is a single static document: one HTML entry, no router, no
// API. Everything the page says about rsk comes from src/data/catalog.json,
// generated from the repository by scripts/gen-data.py — the build never reads
// a live registry.
// The public host is configuration: SITE_URL (https://<SITE_HOST>) is set by `task site:deploy`
// and the deploy workflow. Without it the document carries no canonical, og:url or og:image,
// so a build cannot advertise a URL nobody owns.
declare const process: { env: Record<string, string | undefined> };
const siteUrl = (process.env.SITE_URL ?? '').trim().replace(/\/+$/, '');

/** Canonical, og:url and the absolute og:image, only when a host is configured. */
function deployMeta() {
  return {
    name: 'deploy-meta',
    transformIndexHtml() {
      if (!siteUrl) return [];
      const url = `${siteUrl}/`;
      const image = `${siteUrl}/og.png`;
      return [
        { tag: 'link', attrs: { rel: 'canonical', href: url }, injectTo: 'head' as const },
        { tag: 'meta', attrs: { property: 'og:url', content: url }, injectTo: 'head' as const },
        { tag: 'meta', attrs: { property: 'og:image', content: image }, injectTo: 'head' as const },
        { tag: 'meta', attrs: { name: 'twitter:image', content: image }, injectTo: 'head' as const },
      ];
    },
  };
}

export default defineConfig({
  plugins: [react(), tailwindcss(), deployMeta()],
  build: {
    target: 'es2022',
    cssCodeSplit: false,
    sourcemap: false,
  },
  preview: { port: 4173 },
  server: { port: 5173 },
});

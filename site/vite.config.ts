import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// The landing page is a single static document: one HTML entry, no router, no
// API. Everything the page says about rsk comes from src/data/catalog.json,
// generated from the repository by scripts/gen-data.py — the build never reads
// a live registry.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    target: 'es2022',
    cssCodeSplit: false,
    sourcemap: false,
  },
  preview: { port: 4173 },
  server: { port: 5173 },
});

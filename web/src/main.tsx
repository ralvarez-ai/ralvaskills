import { createRoot, hydrateRoot } from 'react-dom/client';

import App from './App';
import './styles/tokens.css';
import './styles/base.css';
import './styles/a.css';
import './styles/b.css';
import './styles/c.css';
import './styles/d.css';

const root = document.getElementById('root');
if (!root) throw new Error('#root is missing from index.html');

// No StrictMode: the film engine and the one-shot reveal observers are DOM
// side effects, and the double mount in development would replay both.
// The build prerenders the page into #root (scripts/prerender.mjs), so hydrate
// it when present and fall back to a client render in `vite dev`.
if (root.hasChildNodes()) hydrateRoot(root, <App />);
else createRoot(root).render(<App />);

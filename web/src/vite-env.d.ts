/// <reference types="vite/client" />

/** View Transitions are in Chromium; the page treats them as progressive. */
interface Document {
  startViewTransition?: (update: () => void) => { ready: Promise<void>; finished: Promise<void> };
}

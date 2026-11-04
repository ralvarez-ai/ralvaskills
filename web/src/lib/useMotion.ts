import { useCallback, useEffect, useState } from 'react';

const THEMES = ['light', 'dark'] as const;
export type Theme = (typeof THEMES)[number];

/** True when the boot script in index.html allowed motion. */
export function motionAllowed(): boolean {
  return document.documentElement.classList.contains('has-motion');
}

/** Theme lives on <html data-theme>, set before paint; this keeps it in sync. */
export function useTheme() {
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof document === 'undefined') return 'dark';
    return document.documentElement.dataset.theme === 'light' ? 'light' : 'dark';
  });
  // The server render cannot know the stored theme, so the control appears
  // after mount; rendering it in SSR would be a hydration mismatch.
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  useEffect(() => {
    document.documentElement.dataset.theme = theme;
    try {
      localStorage.setItem('rsk-theme', theme);
    } catch {
      // Private mode: the page still switches for this visit.
    }
    document
      .querySelector('meta[name="theme-color"]:not([media])')
      ?.setAttribute('content', theme === 'light' ? '#f4f5f7' : '#0c0e11');
  }, [theme]);

  /**
   * Switch themes. With motion allowed and View Transitions available, the new
   * theme is revealed by a circle growing from the toggle, the same wipe the
   * proposals use; otherwise the attribute changes immediately.
   */
  const toggle = useCallback(
    (event?: { currentTarget: HTMLElement }) => {
      const next: Theme = theme === 'dark' ? 'light' : 'dark';
      const apply = () => setTheme(next);
      if (!motionAllowed() || typeof document.startViewTransition !== 'function') {
        apply();
        return;
      }
      const rect = event?.currentTarget.getBoundingClientRect();
      const x = rect ? rect.left + rect.width / 2 : window.innerWidth / 2;
      const y = rect ? rect.top + rect.height / 2 : 0;
      const radius = Math.hypot(Math.max(x, window.innerWidth - x), Math.max(y, window.innerHeight - y));
      const transition = document.startViewTransition(apply);
      transition.ready
        .then(() => {
          document.documentElement.animate(
            { clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${radius}px at ${x}px ${y}px)`] },
            { duration: 520, easing: 'cubic-bezier(.16,1,.3,1)', pseudoElement: '::view-transition-new(root)' },
          );
        })
        .catch(() => {
          // A cancelled transition leaves the applied theme in place.
        });
    },
    [theme],
  );

  return { theme, toggle, mounted };
}

/** One-shot reveals: add `.is-visible` when a section element enters the viewport. */
export function useReveals() {
  useEffect(() => {
    if (!motionAllowed() || !('IntersectionObserver' in window)) return;
    const targets = document.querySelectorAll<HTMLElement>('.reveal, .rule');
    const io = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue;
          entry.target.classList.add('is-visible');
          io.unobserve(entry.target);
        }
      },
      { rootMargin: '0px 0px -12% 0px', threshold: 0.01 },
    );
    targets.forEach((el) => {
      io.observe(el);
    });
    return () => io.disconnect();
  }, []);
}

/** Count-ups: the HTML already holds the final number, so JS-off stays correct. */
export function useCounts() {
  useEffect(() => {
    if (!motionAllowed() || !('IntersectionObserver' in window)) return;
    const io = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue;
          io.unobserve(entry.target);
          animate(entry.target as HTMLElement);
        }
      },
      { threshold: 0.55 },
    );
    document.querySelectorAll<HTMLElement>('[data-count]').forEach((el) => {
      io.observe(el);
    });
    return () => io.disconnect();
  }, []);
}

function animate(el: HTMLElement) {
  const to = Number(el.dataset.count ?? 0);
  if (!Number.isFinite(to)) return;
  const start = performance.now();
  const duration = 1100;
  const frame = (now: number) => {
    const p = Math.min(1, (now - start) / duration);
    const eased = 1 - (1 - p) ** 3;
    el.textContent = Math.round(to * eased).toLocaleString('en-US');
    if (p < 1) requestAnimationFrame(frame);
    else el.textContent = String(to);
  };
  requestAnimationFrame(frame);
}

/** Reading progress, the stuck header hairline and the active section. */
export function useScrollKit(onActive: (id: string) => void) {
  useEffect(() => {
    const bar = document.querySelector<HTMLElement>('.reading-progress i');
    const header = document.querySelector<HTMLElement>('.site-header');
    const onScroll = () => {
      const max = document.documentElement.scrollHeight - window.innerHeight;
      bar?.style.setProperty('--p', String(max > 0 ? Math.min(1, window.scrollY / max) : 0));
      header?.classList.toggle('is-stuck', window.scrollY > 8);
    };
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  useEffect(() => {
    if (!('IntersectionObserver' in window)) return;
    const io = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) if (entry.isIntersecting) onActive(entry.target.id);
      },
      { rootMargin: '-45% 0px -50% 0px' },
    );
    document.querySelectorAll<HTMLElement>('main section[id]').forEach((el) => {
      io.observe(el);
    });
    return () => io.disconnect();
  }, [onActive]);
}

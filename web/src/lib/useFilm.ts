import { useCallback, useEffect, useRef, useState } from 'react';

import { motionAllowed } from './useMotion';

export type FilmState = 'done' | 'playing' | 'paused';

export interface Step {
  at: number;
  dur?: number;
  fn: (p: number) => void;
  /** Last progress applied, so a step runs only when it changes. */
  _p?: number;
}

export interface Timeline {
  total: number;
  steps: Step[];
  /** Start state: hide the steps, reset the copy. */
  reset: () => void;
  /** Final state: the plain HTML the page ships. */
  end: () => void;
}

/**
 * The film engine: a pausable timeline that plays once when the figure is 30 %
 * visible. It mutates classes on the DOM the section already rendered, and the
 * final state is always the markup, so JS-off and reduced motion see the answer.
 */
export function useFilm(
  rootRef: React.RefObject<HTMLElement | null>,
  ctlRef: React.RefObject<HTMLButtonElement | null>,
  build: () => Timeline | null,
) {
  const [state, setState] = useState<FilmState>('done');
  const toggleRef = useRef<() => void>(() => {});
  const buildRef = useRef(build);
  buildRef.current = build;

  useEffect(() => {
    const root = rootRef.current;
    const ctl = ctlRef.current;
    if (!root || !ctl) return;
    const timeline = buildRef.current();
    if (!timeline) return;

    let t = 0;
    let last = 0;
    let raf = 0;
    let playing = false;
    let started = false;
    let done = false;

    const paint = () => ctl.style.setProperty('--p', String(Math.min(1, t / timeline.total)));

    const apply = () => {
      for (const step of timeline.steps) {
        if (t < step.at) continue;
        const p = step.dur ? Math.min(1, (t - step.at) / step.dur) : 1;
        if (step._p !== p) {
          step._p = p;
          step.fn(p);
        }
      }
      paint();
    };

    const reset = () => {
      t = 0;
      done = false;
      started = true;
      ctl.style.setProperty('--p', '0');
      timeline.steps.forEach((s) => {
        s._p = undefined;
      });
      root.classList.add('is-filming');
      timeline.reset();
    };

    const finish = () => {
      playing = false;
      done = true;
      cancelAnimationFrame(raf);
      timeline.steps.forEach((s) => {
        if (s._p !== 1) {
          s._p = 1;
          s.fn(1);
        }
      });
      timeline.end();
      root.classList.remove('is-filming');
      ctl.style.setProperty('--p', '1');
      setState('done');
    };

    const frame = (now: number) => {
      if (!playing) return;
      t += Math.min(64, now - last);
      last = now;
      apply();
      if (t >= timeline.total) {
        finish();
        return;
      }
      raf = requestAnimationFrame(frame);
    };

    const play = () => {
      if (!started || done) reset();
      playing = true;
      last = performance.now();
      setState('playing');
      raf = requestAnimationFrame(frame);
    };

    const pause = () => {
      playing = false;
      cancelAnimationFrame(raf);
      setState('paused');
    };

    toggleRef.current = () => (playing ? pause() : play());

    const onVisibility = () => {
      if (document.hidden && playing) pause();
    };
    document.addEventListener('visibilitychange', onVisibility);

    let io: IntersectionObserver | undefined;
    if (motionAllowed() && 'IntersectionObserver' in window) {
      io = new IntersectionObserver(
        (entries) => {
          if (entries[0]?.isIntersecting && !started) {
            io?.disconnect();
            play();
          }
        },
        { threshold: 0.3 },
      );
      io.observe(root);
    }

    return () => {
      playing = false;
      cancelAnimationFrame(raf);
      io?.disconnect();
      document.removeEventListener('visibilitychange', onVisibility);
    };
  }, [rootRef, ctlRef]);

  const toggle = useCallback(() => toggleRef.current(), []);

  return { state, toggle };
}

/** Film helpers that work on the rendered figure. */
export function show(root: HTMLElement, key: string) {
  root.querySelectorAll(`[data-f="${key}"]`).forEach((el) => el.classList.remove('is-hidden'));
}

export function hideAll(root: HTMLElement) {
  root.querySelectorAll('[data-f]').forEach((el) => el.classList.add('is-hidden'));
}

export function revealAll(root: HTMLElement) {
  root.querySelectorAll('.is-hidden').forEach((el) => el.classList.remove('is-hidden'));
  root.querySelectorAll<HTMLElement>('[data-t]').forEach((el) => {
    el.textContent = el.dataset.t ?? '';
    el.classList.remove('is-typing');
  });
}

export function type(el: HTMLElement | null, p: number) {
  if (!el) return;
  const full = el.dataset.t ?? '';
  el.textContent = full.slice(0, Math.round(full.length * p));
  el.classList.toggle('is-typing', p < 1);
}

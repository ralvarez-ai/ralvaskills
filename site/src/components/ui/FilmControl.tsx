import { forwardRef } from 'react';

import type { FilmState } from '../../lib/useFilm';

const LABEL: Record<FilmState, string> = { playing: 'Pause', paused: 'Play', done: 'Replay' };
const ICON: Record<FilmState, string> = { playing: '❚❚', paused: '▶', done: '↻' };

/** Pause / play / replay, with the timeline's progress in the bar. */
export const FilmControl = forwardRef<HTMLButtonElement, { state: FilmState; onToggle: () => void }>(
  function FilmControl({ state, onToggle }, ref) {
    return (
      <button ref={ref} className="film-control" type="button" data-state={state} onClick={onToggle}>
        <span className="film-control__icon" aria-hidden="true">
          {ICON[state]}
        </span>
        <span className="film-control__label">{LABEL[state]}</span>
        <span className="film-control__bar" aria-hidden="true">
          <i />
        </span>
        <span className="sr-only"> the hero animation</span>
      </button>
    );
  },
);

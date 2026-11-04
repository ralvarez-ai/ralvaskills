import { useRef } from 'react';
import type { ReactNode } from 'react';

import { FilmControl, Panel, Tag } from './ui';
import { rskNewOutput } from '../data/facts';
import { data } from '../data/types';
import { hideAll, revealAll, show, type, useFilm } from '../lib/useFilm';

interface ManifestProps {
  names: string[];
  second: string;
  agent?: boolean;
  caption: ReactNode;
}

const S = ({ value }: { value: string }) => <span className="syntax-string">{`"${value}"`}</span>;

function Cmd({ k, text }: { k: string; text: string }) {
  return (
    <li data-f={k}>
      <span className="command-list__prompt" aria-hidden="true">
        $
      </span>
      <code>
        <span data-t={text} aria-hidden="true">
          {text}
        </span>
        <span className="sr-only">{text}</span>
      </code>
    </li>
  );
}

/**
 * The manifest device: the terminal, the manifest assembling line by line and
 * the lock resolving beneath it. `second` is the line the agent's call stands
 * for when `agent` is set.
 */
export function Manifest({ names, second, agent = false, caption }: ManifestProps) {
  const fig = useRef<HTMLElement>(null);
  const ctl = useRef<HTMLButtonElement>(null);
  const pin = useRef<HTMLSpanElement>(null);
  const pinLine = useRef<HTMLLIElement>(null);

  const sorted = [...names].sort();
  const mod = [
    { k: 'new', body: <><span className="syntax-key">version</span> = <span className="syntax-number">1</span></> },
    {
      k: 'new',
      body: (
        <>
          <span className="syntax-key">tools</span> = [<S value="claude-code" />, <S value="opencode" />]
        </>
      ),
    },
    {
      k: 'new',
      body: (
        <>
          <span className="syntax-key">pinned</span> = <span className="manifest-file__pinned">[<S value="go-architect" />]</span>
        </>
      ),
      pinned: true,
    },
    { k: 'new', body: <>&nbsp;</> },
    { k: 'new', body: <span className="syntax-table">[skills]</span> },
    ...sorted.map((n) => ({
      k: `s${names.indexOf(n)}`,
      body: (
        <>
          {'  '}
          <span className="syntax-key">{n}</span> = <S value="*" />
        </>
      ),
    })),
  ];

  const { state, toggle } = useFilm(fig, ctl, () => {
    const root = fig.current;
    if (!root) return null;
    const pinEl = pin.current;
    const pinRow = pinLine.current;
    const rows = root.querySelectorAll('tr[data-f]').length;
    const steps: { at: number; dur?: number; fn: (p: number) => void }[] = [
      { at: 150, fn: () => show(root, 'c1') },
      { at: 200, dur: 550, fn: (p) => type(root.querySelector('[data-f="c1"] [data-t]'), p) },
      { at: 950, fn: () => show(root, 'o1') },
      { at: 1150, fn: () => show(root, 'new') },
      { at: 1700, fn: () => show(root, 'c2') },
      { at: 1750, dur: 900, fn: (p) => type(root.querySelector('[data-f="c2"] [data-t]'), p) },
    ];
    for (let i = 0; i < rows; i += 1) steps.push({ at: 2850 + i * 230, fn: () => show(root, `s${i}`) });
    steps.push({ at: 4200, fn: () => show(root, 'c3') });
    steps.push({ at: 4250, dur: 800, fn: (p) => type(root.querySelector('[data-f="c3"] [data-t]'), p) });
    steps.push({
      at: 5250,
      fn: () => {
        if (!pinEl) return;
        pinEl.dataset.full ??= pinEl.innerHTML;
        pinEl.innerHTML = pinEl.dataset.full;
        pinRow?.classList.add('is-changed');
      },
    });

    return {
      total: 6400,
      steps,
      reset: () => {
        hideAll(root);
        root.querySelectorAll<HTMLElement>('[data-t]').forEach((el) => {
          el.textContent = '';
        });
        if (pinEl) {
          pinEl.dataset.full ??= pinEl.innerHTML;
          pinEl.textContent = '[]';
        }
        pinRow?.classList.remove('is-changed');
      },
      end: () => {
        revealAll(root);
        if (pinEl?.dataset.full) pinEl.innerHTML = pinEl.dataset.full;
        if (pinRow) setTimeout(() => pinRow.classList.remove('is-changed'), 900);
      },
    };
  });

  return (
    <figure className="manifest" id="a-film" aria-describedby="a-cap">
      <div className="manifest__terminal surface-terminal">
        <div className="panel__header">
          <span className="panel__dots" aria-hidden="true">
            <i />
            <i />
            <i />
          </span>
          ~/demo — a Go project
        </div>
        <ol className="command-list" aria-label="Commands run">
          <Cmd k="c1" text="rsk new" />
          {rskNewOutput.map((line) => (
            <li className="command-list__output" data-f="o1" key={line}>
              {line}
            </li>
          ))}
          {agent ? (
            <li data-f="c2" className="command-list__agent">
              <span className="command-list__prompt" aria-hidden="true">
                ◆
              </span>
              <code>
                <span data-t={second} aria-hidden="true">
                  {second}
                </span>
                <span className="sr-only">{second}</span>
              </code>
            </li>
          ) : (
            <Cmd k="c2" text={second} />
          )}
          <Cmd k="c3" text="rsk pin go-architect" />
        </ol>
      </div>

      <div className="manifest__files">
        <Panel title=".rsk/rsk.mod" action={<Tag tone="readonly">commit</Tag>} flush>
          <ol className="manifest-file" aria-label=".rsk/rsk.mod contents">
            {mod.map((line, i) => (
              <li
                className="manifest-file__line"
                data-f={line.k}
                key={`${line.k}-${i}`}
                id={line.pinned ? 'a-pin' : undefined}
                ref={line.pinned ? pinLine : undefined}
              >
                <span className="gutter" aria-hidden="true">
                  +
                </span>
                <span>
                  {line.pinned ? <span ref={pin}>{line.body}</span> : line.body}
                </span>
              </li>
            ))}
          </ol>
        </Panel>
        <Panel
          title=".rsk/rsk.lock"
          action={<Tag tone="readonly">commit</Tag>}
          bodyStyle={{ padding: '6px 0 4px' }}
        >
            <table className="lock-table">
              <caption className="sr-only">rsk.lock resolved versions, all from the registry, shown as a table</caption>
              <thead>
                <tr>
                  <th scope="col">skill</th>
                  <th scope="col" className="lock-table__version">
                    version
                  </th>
                </tr>
              </thead>
              <tbody>
                {names.map((n, i) => (
                  <tr data-f={`s${i}`} key={n}>
                    <td>{n}</td>
                    <td className="lock-table__version">{data.latest[n]}</td>
                  </tr>
                ))}
              </tbody>
            </table>
        </Panel>
      </div>

      <figcaption className="manifest__caption" id="a-cap">
        <p>{caption}</p>
        <FilmControl state={state} onToggle={toggle} ref={ctl} />
      </figcaption>
    </figure>
  );
}

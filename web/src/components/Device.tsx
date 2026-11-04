import { useRef } from 'react';
import type { ReactNode } from 'react';

import { Badge, FilmButton } from './ui';
import { rskNewOutput } from '../data/facts';
import { data } from '../data/types';
import { hideAll, revealAll, show, type, useFilm } from '../lib/useFilm';

interface DeviceProps {
  names: string[];
  second: string;
  agent?: boolean;
  caption: ReactNode;
}

const S = ({ value }: { value: string }) => <span className="s">{`"${value}"`}</span>;

function Cmd({ k, text }: { k: string; text: string }) {
  return (
    <li data-f={k}>
      <span className="ps" aria-hidden="true">
        $
      </span>
      <code>
        <span data-t={text} aria-hidden="true">
          {text}
        </span>
        <span className="sr">{text}</span>
      </code>
    </li>
  );
}

/**
 * Direction A's device: the terminal, the manifest assembling line by line and
 * the lock resolving beneath it. `second` is the line the agent's call stands
 * for when `agent` is set.
 */
export function Device({ names, second, agent = false, caption }: DeviceProps) {
  const fig = useRef<HTMLElement>(null);
  const ctl = useRef<HTMLButtonElement>(null);
  const pin = useRef<HTMLSpanElement>(null);
  const pinLine = useRef<HTMLLIElement>(null);

  const sorted = [...names].sort();
  const mod = [
    { k: 'new', body: <><span className="k">version</span> = <span className="n">1</span></> },
    {
      k: 'new',
      body: (
        <>
          <span className="k">tools</span> = [<S value="claude-code" />, <S value="opencode" />]
        </>
      ),
    },
    {
      k: 'new',
      body: (
        <>
          <span className="k">pinned</span> = <span className="pin">[<S value="go-architect" />]</span>
        </>
      ),
      pinned: true,
    },
    { k: 'new', body: <>&nbsp;</> },
    { k: 'new', body: <span className="t">[skills]</span> },
    ...sorted.map((n) => ({
      k: `s${names.indexOf(n)}`,
      body: (
        <>
          {'  '}
          <span className="k">{n}</span> = <S value="*" />
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
        pinRow?.classList.add('chg');
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
        pinRow?.classList.remove('chg');
      },
      end: () => {
        revealAll(root);
        if (pinEl?.dataset.full) pinEl.innerHTML = pinEl.dataset.full;
        if (pinRow) setTimeout(() => pinRow.classList.remove('chg'), 900);
      },
    };
  });

  return (
    <figure className="a-dev" id="a-film" aria-describedby="a-cap">
      <div className="a-term tn">
        <div className="k-win-h">
          <span className="dots" aria-hidden="true">
            <i />
            <i />
            <i />
          </span>
          ~/demo — a Go project
        </div>
        <ol className="a-cmds" aria-label="Commands run">
          <Cmd k="c1" text="rsk new" />
          {rskNewOutput.map((line) => (
            <li className="out" data-f="o1" key={line}>
              {line}
            </li>
          ))}
          {agent ? (
            <li data-f="c2" className="ag">
              <span className="ps" aria-hidden="true">
                ◆
              </span>
              <code>
                <span data-t={second} aria-hidden="true">
                  {second}
                </span>
                <span className="sr">{second}</span>
              </code>
            </li>
          ) : (
            <Cmd k="c2" text={second} />
          )}
          <Cmd k="c3" text="rsk pin go-architect" />
        </ol>
      </div>

      <div className="a-files">
        <div className="k-win">
          <div className="k-win-h">
            .rsk/rsk.mod <Badge kind="ro">commit</Badge>
          </div>
          <ol className="a-mod" aria-label=".rsk/rsk.mod contents">
            {mod.map((line, i) => (
              <li
                className="add"
                data-f={line.k}
                key={`${line.k}-${i}`}
                id={line.pinned ? 'a-pin' : undefined}
                ref={line.pinned ? pinLine : undefined}
              >
                <span className="g" aria-hidden="true">
                  +
                </span>
                <span>
                  {line.pinned ? <span ref={pin}>{line.body}</span> : line.body}
                </span>
              </li>
            ))}
          </ol>
        </div>
        <div className="k-win">
          <div className="k-win-h">
            .rsk/rsk.lock <Badge kind="ro">commit</Badge>
          </div>
          <div className="k-win-b" style={{ padding: '6px 0 4px' }}>
            <table className="a-lock">
              <caption className="sr">rsk.lock resolved versions, all from the registry, shown as a table</caption>
              <thead>
                <tr>
                  <th scope="col">skill</th>
                  <th scope="col" className="v">
                    version
                  </th>
                </tr>
              </thead>
              <tbody>
                {names.map((n, i) => (
                  <tr data-f={`s${i}`} key={n}>
                    <td>{n}</td>
                    <td className="v">{data.latest[n]}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <figcaption className="a-cap" id="a-cap">
        <p>{caption}</p>
        <FilmButton state={state} onToggle={toggle} ref={ctl} />
      </figcaption>
    </figure>
  );
}

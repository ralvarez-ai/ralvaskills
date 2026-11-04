import { Fragment, useRef } from 'react';
import type { CSSProperties } from 'react';

import { Badge, Btn, FilmButton, Mark } from '../components/ui';
import { profile, profileNames } from '../data/facts';
import { data } from '../data/types';
import { hideAll, revealAll, show, useFilm } from '../lib/useFilm';

const NAMES = profileNames.join(', ');
const INSTALLED = profile.candidates.map((c) => `${c.skill} ${c.latest}`).join(' · ');
const PROPOSES = profile.candidates.map((c) => `${c.skill} ${c.latest}`).join(' · ');
const NAMES_JSON = `{"names": [${profileNames.map((n) => `"${n}"`).join(', ')}]}`;

export function Hero() {
  const fig = useRef<HTMLElement>(null);
  const ctl = useRef<HTMLButtonElement>(null);
  const you = useRef<HTMLDivElement>(null);

  const { state, toggle } = useFilm(fig, ctl, () => {
    const root = fig.current;
    const card = you.current;
    if (!root || !card) return null;
    return {
      total: 6800,
      steps: [
        { at: 250, fn: () => show(root, 'a1') },
        { at: 1150, fn: () => show(root, 'a2') },
        { at: 2450, fn: () => show(root, 'b1') },
        { at: 3050, fn: () => show(root, 'b2') },
        { at: 4250, fn: () => card.classList.remove('pend') },
        { at: 4850, fn: () => show(root, 'c1') },
        { at: 5750, fn: () => show(root, 'c2') },
      ],
      reset: () => {
        hideAll(root);
        card.classList.add('pend');
      },
      end: () => {
        revealAll(root);
        card.classList.remove('pend');
      },
    };
  });

  return (
    <section className="c-hero" id="top" aria-labelledby="c-h1">
      <div className="k-grid" aria-hidden="true" />
      <div className="wrap">
        <div className="c-top">
          <div>
            <p className="k-kick rv">
              <Mark /> rsk mcp · for Claude Code and OpenCode
            </p>
            <h1 className="c-h1 rv" id="c-h1" style={{ '--i': 1 } as CSSProperties}>
              <span className="ln">
                <span>Your agent asks.</span>
              </span>
              <span className="ln">
                <span className="you">You approve.</span>
              </span>
              <span className="ln">
                <span className="rsk">rsk installs.</span>
              </span>
            </h1>
          </div>
          <div>
            <p className="c-lede rv" style={{ '--i': 2 } as CSSProperties}>
              <code>rsk</code> gives your coding agent three tools. Two only read: your project’s files and a catalogue
              of {data.skills.length} skills. The one that writes, <code>install_skills</code>, makes Claude Code ask you
              on every call.
            </p>
            <div className="c-cta rv" style={{ '--i': 3 } as CSSProperties}>
              <Btn href="#start" primary>
                Install rsk
              </Btn>
              <Btn href="#tools">The three tools</Btn>
            </div>
          </div>
        </div>

        <figure className="c-fig" id="c-film" aria-describedby="c-cap">
          <ol className="c-seq" aria-label="The three calls, in order">
            <li className="c-lanes" aria-hidden="true">
              <span className="c-lane y">You</span>
              <span className="c-lane">Agent</span>
              <span className="c-lane">rsk mcp</span>
            </li>

            <li className="c-ex">
              <span className="c-no c-nod" aria-hidden="true">
                1
              </span>
              <p className="c-route">
                <span className="c-no">1</span>Agent<b>→</b>rsk mcp
              </p>
              <div className="c-m req" data-f="a1">
                <div className="c-arrow r" aria-hidden="true" />
                <div className="c-card">
                  <span className="m">project_profile</span>
                  <Badge kind="ro">read-only</Badge>
                  <span className="pf"> {'{}'}</span>
                </div>
              </div>
              <div className="c-m" data-f="a2">
                <div className="c-arrow lt back" aria-hidden="true" />
                <div className="c-card back">
                  <span className="sr">rsk mcp replies: </span>
                  <dl>
                    <dt>signals</dt>
                    <dd>
                      <dl className="sg" style={{ margin: 0 }}>
                        {profile.signals.map((s) => (
                          <Fragment key={s.kind}>
                            <dt>{s.kind}</dt>
                            <dd>
                              <span className="pf">←</span> {s.proof}
                            </dd>
                          </Fragment>
                        ))}
                      </dl>
                    </dd>
                    <dt>proposes</dt>
                    <dd>{PROPOSES}</dd>
                  </dl>
                </div>
              </div>
            </li>

            <li className="c-ex h">
              <span className="c-no c-nod" aria-hidden="true">
                2
              </span>
              <p className="c-route">
                <span className="c-no">2</span>Agent<b>→</b>You
              </p>
              <div className="c-m l req" data-f="b1">
                <div className="c-arrow lt ask" aria-hidden="true" />
                <div className="c-card">
                  <span className="m">install_skills</span>
                  <Badge kind="de">writes</Badge>
                  <span className="pf"> needs your approval</span>
                </div>
              </div>
              <div className="c-youcell" data-f="b2">
                <div className="c-you" ref={you}>
                  <p className="q">
                    Allow <code>install_skills</code> from rsk?
                  </p>
                  <p className="nm">{NAMES}</p>
                  <p className="acts" style={{ margin: 0 }}>
                    <span className="b">Deny</span>
                    <span className="b al">
                      <span className="p">Allow</span>
                      <span className="d">✓ Allowed</span>
                    </span>
                    <span className="lbl">you, every call</span>
                  </p>
                </div>
              </div>
            </li>

            <li className="c-ex">
              <span className="c-no c-nod" aria-hidden="true">
                3
              </span>
              <p className="c-route">
                <span className="c-no">3</span>Agent<b>→</b>rsk mcp
              </p>
              <div className="c-m req" data-f="c1">
                <div className="c-arrow r" aria-hidden="true" />
                <div className="c-card">
                  <span className="m">install_skills</span>
                  <Badge kind="de">writes</Badge> <span className="pf">{NAMES_JSON}</span>
                </div>
              </div>
              <div className="c-m" data-f="c2">
                <div className="c-arrow lt back" aria-hidden="true" />
                <div className="c-card back">
                  <span className="sr">rsk mcp replies: </span>
                  <dl>
                    <dt>installed</dt>
                    <dd>{INSTALLED}</dd>
                    <dt>failed</dt>
                    <dd>none</dd>
                  </dl>
                </div>
              </div>
            </li>
          </ol>
          <figcaption className="c-cap" id="c-cap">
            <p>
              <Badge kind="ro">real</Badge> Call 1 is a real <code>project_profile</code> response for a scratch Go folder
              (go.mod, Dockerfile, cmd/demo/main.go). <Badge kind="sa">sample</Badge> The approval is drawn, not a capture
              of the client dialog, and call 3’s reply is the shape from <code>internal/mcp/install_tool.go</code> with
              real latest versions.
            </p>
            <FilmButton state={state} onToggle={toggle} ref={ctl} />
          </figcaption>
        </figure>
      </div>
    </section>
  );
}

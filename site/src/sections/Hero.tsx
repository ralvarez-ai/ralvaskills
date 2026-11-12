import type { CSSProperties } from 'react';
import { Fragment, useRef } from 'react';

import { BrandMark, ButtonLink, Eyebrow, FilmControl, Tag } from '../components/ui';
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
        { at: 4250, fn: () => card.classList.remove('is-pending') },
        { at: 4850, fn: () => show(root, 'c1') },
        { at: 5750, fn: () => show(root, 'c2') },
      ],
      reset: () => {
        hideAll(root);
        card.classList.add('is-pending');
      },
      end: () => {
        revealAll(root);
        card.classList.remove('is-pending');
      },
    };
  });

  return (
    <section className="agent-hero" id="top" aria-labelledby="c-h1">
      <div className="grid-texture" aria-hidden="true" />
      <div className="container">
        <div className="agent-hero__top">
          <div>
            <Eyebrow>
              <BrandMark /> rsk mcp · for ten coding clients
            </Eyebrow>
            <h1 className="agent-hero__title reveal" id="c-h1" style={{ '--i': 1 } as CSSProperties}>
              <span className="agent-hero__line">
                <span>Your agent asks.</span>
              </span>
              <span className="agent-hero__line">
                <span className="is-human">You approve.</span>
              </span>
              <span className="agent-hero__line">
                <span className="is-product">rsk installs.</span>
              </span>
            </h1>
          </div>
          <div>
            <p className="agent-hero__lede reveal" style={{ '--i': 2 } as CSSProperties}>
              <code>rsk</code> gives your coding agent three tools. Two only read: your project’s files and a catalogue
              of {data.skills.length} skills. The one that writes, <code>install_skills</code>, makes Claude Code ask
              you on every call; the other clients get the read-only pair until their approval flow is verified.
            </p>
            <div className="agent-hero__actions reveal" style={{ '--i': 3 } as CSSProperties}>
              <ButtonLink href="#start" variant="primary">
                Install rsk
              </ButtonLink>
              <ButtonLink href="#tools">The three tools</ButtonLink>
            </div>
          </div>
        </div>

        <figure ref={fig} className="agent-loop" id="c-film" aria-describedby="c-cap">
          <ol className="agent-loop__list" aria-label="The three calls, in order">
            <li className="agent-loop__lanes" aria-hidden="true">
              <span className="lane is-human">You</span>
              <span className="lane">Agent</span>
              <span className="lane">rsk mcp</span>
            </li>

            <li className="call">
              <span className="call__number call__number--inline" aria-hidden="true">
                1
              </span>
              <p className="call__route">
                <span className="call__number">1</span>Agent<b>→</b>rsk mcp
              </p>
              <div className="call__message call__message--request" data-f="a1">
                <div className="call-arrow call-arrow--right" aria-hidden="true" />
                <div className="call-card">
                  <span className="call-card__method">project_profile</span>
                  <Tag tone="readonly">read-only</Tag>
                  <span className="call-card__note"> {'{}'}</span>
                </div>
              </div>
              <div className="call__message" data-f="a2">
                <div className="call-arrow call-arrow--left is-reply" aria-hidden="true" />
                <div className="call-card is-reply">
                  <span className="sr-only">rsk mcp replies: </span>
                  <dl>
                    <dt>signals</dt>
                    <dd>
                      <dl className="signal-pairs" style={{ margin: 0 }}>
                        {profile.signals.map((s) => (
                          <Fragment key={s.kind}>
                            <dt>{s.kind}</dt>
                            <dd>
                              <span className="call-card__note">←</span> {s.proof}
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

            <li className="call is-human">
              <span className="call__number call__number--inline" aria-hidden="true">
                2
              </span>
              <p className="call__route">
                <span className="call__number">2</span>Agent<b>→</b>You
              </p>
              <div className="call__message call__message--left call__message--request" data-f="b1">
                <div className="call-arrow call-arrow--left call-arrow--ask" aria-hidden="true" />
                <div className="call-card">
                  <span className="call-card__method">install_skills</span>
                  <Tag tone="destructive">writes</Tag>
                  <span className="call-card__note"> needs your approval</span>
                </div>
              </div>
              <div className="approval__cell" data-f="b2">
                <div className="approval" ref={you}>
                  <p className="approval__question">
                    Allow <code>install_skills</code> from rsk?
                  </p>
                  <p className="approval__names">{NAMES}</p>
                  <p className="approval__actions" style={{ margin: 0 }}>
                    <span className="approval__button">Deny</span>
                    <span className="approval__button approval__allow">
                      <span className="approval__label-pending">Allow</span>
                      <span className="approval__label-done">✓ Allowed</span>
                    </span>
                    <span className="approval__note">you, every call</span>
                  </p>
                </div>
              </div>
            </li>

            <li className="call">
              <span className="call__number call__number--inline" aria-hidden="true">
                3
              </span>
              <p className="call__route">
                <span className="call__number">3</span>Agent<b>→</b>rsk mcp
              </p>
              <div className="call__message call__message--request" data-f="c1">
                <div className="call-arrow call-arrow--right" aria-hidden="true" />
                <div className="call-card">
                  <span className="call-card__method">install_skills</span>
                  <Tag tone="destructive">writes</Tag> <span className="call-card__note">{NAMES_JSON}</span>
                </div>
              </div>
              <div className="call__message" data-f="c2">
                <div className="call-arrow call-arrow--left is-reply" aria-hidden="true" />
                <div className="call-card is-reply">
                  <span className="sr-only">rsk mcp replies: </span>
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
          <figcaption className="agent-loop__caption" id="c-cap">
            <p>
              <Tag tone="readonly">real</Tag> Call 1 is a real <code>project_profile</code> response for a scratch Go
              folder (go.mod, Dockerfile, cmd/demo/main.go). <Tag tone="sample">sample</Tag> The approval is drawn, not
              a capture of the client dialog, and call 3’s reply is the shape from{' '}
              <code>internal/mcp/install_tool.go</code> with real latest versions.
            </p>
            <FilmControl state={state} onToggle={toggle} ref={ctl} />
          </figcaption>
        </figure>
      </div>
    </section>
  );
}

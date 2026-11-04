import { Eyebrow, Section } from '../components/ui';
import type { ReactNode } from 'react';

import { data } from '../data/types';

const QA: [string, ReactNode][] = [
  [
    'Do I need the MCP server?',
    <>
      No. Every step works from the CLI: <code>rsk install</code>, <code>rsk pin</code>, <code>rsk catalog</code>. The
      server is there so the agent can propose skills from your files instead of you naming them.
    </>,
  ],
  [
    'Which clients does it support?',
    <>
      Claude Code and OpenCode. <code>rsk new --for claude-code|opencode|all</code> picks which config files rsk writes;
      installed skills are linked once into <code>.claude/skills/</code>.
    </>,
  ],
  [
    'Can I install a skill for every project?',
    <>
      Yes: <code>rsk install &lt;name&gt; --global</code> links it into <code>~/.claude/skills/</code> and{' '}
      <code>~/.config/opencode/skills/</code>. The <code>global</code> bundle is meant for that.
    </>,
  ],
  [
    'Can I define my own bundles?',
    <>
      Write <code>~/.config/rsk/catalog.toml</code>. A bundle with the same name replaces the built-in one; new names
      are added.
    </>,
  ],
  [
    'What are personal skills?',
    <>
      {data.personal} of the {data.skills.length} are marked personal in the registry index. <code>rsk catalog</code>{' '}
      hides them unless you pass <code>--include personal</code>.
    </>,
  ],
  [
    'How do updates work?',
    <>
      <code>rsk status</code> shows installed skills and version drift; <code>rsk update [name]</code> pulls the latest
      and re-links. Each skill is versioned and released on its own.
    </>,
  ],
  ['What does it cost?', <>Nothing. rsk and the catalogue are GPL-3.0, on GitHub.</>],
];

export function Faq() {
  return (
    <Section id="faq" labelledBy="d-faq-h">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>Questions</Eyebrow>
        <h2 id="d-faq-h" className="reveal">
          Before you install.
        </h2>
        <div className="faq reveal">
          {QA.map(([q, a]) => (
            <details key={q}>
              <summary>{q}</summary>
              <p>{a}</p>
            </details>
          ))}
        </div>
      </div>
    </Section>
  );
}

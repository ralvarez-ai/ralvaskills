import { Eyebrow, Section } from '../components/ui';
import { helpGroups } from '../data/facts';

export function Commands() {
  return (
    <Section id="commands" labelledBy="a-cmd-h">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>The whole surface</Eyebrow>
        <h2 id="a-cmd-h" className="reveal">
          Fourteen commands, grouped the way <code>rsk --help</code> groups them.
        </h2>
        <ul className="help-groups reveal">
          {helpGroups.map((group) => (
            <li key={group.group}>
              <h3>{group.group}</h3>
              <dl>
                {group.commands.map((cmd) => (
                  <div key={cmd.name}>
                    <dt>{cmd.name}</dt>
                    <dd>{cmd.desc}</dd>
                  </div>
                ))}
              </dl>
            </li>
          ))}
        </ul>
        <p style={{ marginTop: 20, color: 'var(--text-muted)' }}>
          Data on stdout, logs on stderr, <code>--output json</code> on the read-only views, <code>NO_COLOR</code>{' '}
          honoured.
        </p>
      </div>
    </Section>
  );
}

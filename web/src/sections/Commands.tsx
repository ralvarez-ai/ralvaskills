import { helpGroups } from '../data/facts';

export function Commands() {
  return (
    <section className="k-sec" id="commands" aria-labelledby="a-cmd-h">
      <div className="wrap">
        <div className="k-rule" aria-hidden="true" />
        <p className="k-kick rv">The whole surface</p>
        <h2 id="a-cmd-h" className="rv">
          Thirteen commands, grouped the way <code>rsk --help</code> groups them.
        </h2>
        <ul className="a-help rv">
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
        <p style={{ marginTop: 20, color: 'var(--muted)' }}>
          Data on stdout, logs on stderr, <code>--output json</code> on the read-only views, <code>NO_COLOR</code>{' '}
          honoured.
        </p>
      </div>
    </section>
  );
}

import { Eyebrow, Panel, Section, Tag } from '../components/ui';
import { data } from '../data/types';

const LOCK = [...data.repoLock].sort((a, b) => a.name.localeCompare(b.name));

export function ThisRepo() {
  return (
    <Section id="this-repo" labelledBy="a-real-h" tone="ink" className="manifest-ink">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>The real file</Eyebrow>
        <h2 id="a-real-h" className="reveal">
          The repository that ships rsk runs on rsk.
        </h2>
        <div className="real-manifest">
          <div className="reveal">
            <Panel
              title="ralvaskills/.rsk/rsk.mod"
              action={<Tag tone="readonly">real · working tree 2026-11-04</Tag>}
              flush
            >
              <pre className="panel__body" tabIndex={0} role="region" aria-label="ralvaskills .rsk/rsk.mod">
                {data.repoMod.trimEnd()}
              </pre>
            </Panel>
          </div>
          <div className="reveal">
            <ul className="note-list">
              <li>
                <b>
                  {LOCK.length} skills, {data.repoPinned.length} pinned.
                </b>{' '}
                Pinned skills are imported through <code>.rsk/CLAUDE.md</code> and <code>opencode.json</code>, so they
                load in every session.
              </li>
              <li>
                <b>
                  <code>"*"</code> means latest;
                </b>{' '}
                the lock records the exact version that resolved, so a clone gets the same ones.
              </li>
              <li>
                <b>Both tools from one list.</b> <code>tools = ["claude-code", "opencode"]</code> drives which adapter
                files rsk writes.
              </li>
            </ul>
            <table className="lock-table" style={{ marginTop: 28 }}>
              <caption>ralvaskills/.rsk/rsk.lock, shown as a table</caption>
              <thead>
                <tr>
                  <th scope="col">skill</th>
                  <th scope="col" className="lock-table__version">
                    version
                  </th>
                  <th scope="col">source</th>
                </tr>
              </thead>
              <tbody>
                {LOCK.map((entry) => (
                  <tr key={entry.name}>
                    <td>{entry.name}</td>
                    <td className="lock-table__version">{entry.version}</td>
                    <td>{entry.source}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Section>
  );
}

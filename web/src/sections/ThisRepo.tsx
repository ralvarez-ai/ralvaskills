import { Badge } from '../components/ui';
import { data } from '../data/types';

const LOCK = [...data.repoLock].sort((a, b) => a.name.localeCompare(b.name));

export function ThisRepo() {
  return (
    <section className="k-sec night a-night" id="this-repo" aria-labelledby="a-real-h">
      <div className="wrap">
        <div className="k-rule" aria-hidden="true" />
        <p className="k-kick rv">The real file</p>
        <h2 id="a-real-h" className="rv">
          The repository that ships rsk runs on rsk.
        </h2>
        <div className="a-real">
          <div className="rv">
            <div className="k-win">
              <div className="k-win-h">
                ralvaskills/.rsk/rsk.mod <Badge kind="ro">real · working tree 2026-11-04</Badge>
              </div>
              <pre className="k-win-b" tabIndex={0} role="region" aria-label="ralvaskills .rsk/rsk.mod">
                {data.repoMod.trimEnd()}
              </pre>
            </div>
          </div>
          <div className="rv">
            <ul className="a-notes">
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
            <table className="a-lock" style={{ marginTop: 28 }}>
              <caption>ralvaskills/.rsk/rsk.lock, shown as a table</caption>
              <thead>
                <tr>
                  <th scope="col">skill</th>
                  <th scope="col" className="v">
                    version
                  </th>
                  <th scope="col">source</th>
                </tr>
              </thead>
              <tbody>
                {LOCK.map((entry) => (
                  <tr key={entry.name}>
                    <td>{entry.name}</td>
                    <td className="v">{entry.version}</td>
                    <td>{entry.source}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  );
}

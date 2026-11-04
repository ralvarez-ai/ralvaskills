import { Badge } from '../components/ui';
import { profile, signals } from '../data/facts';

const HITS = new Set(profile.signals.map((s) => s.kind));
const JSON_TEXT = JSON.stringify(profile, null, 1);

export function Signals() {
  return (
    <section className="k-sec" id="profile" aria-labelledby="c-sig-h">
      <div className="wrap">
        <div className="k-rule" aria-hidden="true" />
        <p className="k-kick rv">From internal/mcp/signals.go</p>
        <h2 id="c-sig-h" className="rv">
          What <code>project_profile</code> looks for.
        </h2>
        <div className="c-two">
          <div className="rv">
            <table className="c-sig">
              <caption className="sr">Signals, the files that prove them, and the skills proposed</caption>
              <thead>
                <tr>
                  <th scope="col">Signal</th>
                  <th scope="col">Proved by</th>
                  <th scope="col">Proposes</th>
                </tr>
              </thead>
              <tbody>
                {signals.map((s) => (
                  <tr className={HITS.has(s.kind) ? 'hit' : undefined} key={s.kind}>
                    <td>{s.kind}</td>
                    <td>{s.proof}</td>
                    <td>{s.proposes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <p style={{ marginTop: 14, color: 'var(--muted)', fontSize: '.92rem' }}>
              It proposes, it does not rank. Refactoring skills are never proposed by a file; only when the agent is
              already refactoring. ● matched in the scratch folder.
            </p>
          </div>
          <div className="c-json rv">
            <p className="k-kick" style={{ marginBottom: 10 }}>
              The full response <Badge kind="ro">real</Badge>
            </p>
            <pre tabIndex={0} role="region" aria-label="project_profile JSON response">
              {JSON_TEXT}
            </pre>
          </div>
        </div>
      </div>
    </section>
  );
}

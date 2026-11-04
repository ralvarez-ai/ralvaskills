import { Eyebrow, Section, Tag } from '../components/ui';
import { profile, signals } from '../data/facts';

const HITS = new Set(profile.signals.map((s) => s.kind));
const JSON_TEXT = JSON.stringify(profile, null, 1);

export function Signals() {
  return (
    <Section id="profile" labelledBy="c-sig-h">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>From internal/mcp/signals.go</Eyebrow>
        <h2 id="c-sig-h" className="reveal">
          What <code>project_profile</code> looks for.
        </h2>
        <div className="split">
          <div className="reveal">
            <table className="signal-table">
              <caption className="sr-only">Signals, the files that prove them, and the skills proposed</caption>
              <thead>
                <tr>
                  <th scope="col">Signal</th>
                  <th scope="col">Proved by</th>
                  <th scope="col">Proposes</th>
                </tr>
              </thead>
              <tbody>
                {signals.map((s) => (
                  <tr className={HITS.has(s.kind) ? 'is-hit' : undefined} key={s.kind}>
                    <td>{s.kind}</td>
                    <td>{s.proof}</td>
                    <td>{s.proposes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <p style={{ marginTop: 14, color: 'var(--text-muted)', fontSize: '.92rem' }}>
              It proposes, it does not rank. Refactoring skills are never proposed by a file; only when the agent is
              already refactoring. ● matched in the scratch folder.
            </p>
          </div>
          <div className="json-block reveal">
            <Eyebrow reveal={false} style={{ marginBottom: 10 }}>
              The full response <Tag tone="readonly">real</Tag>
            </Eyebrow>
            <pre tabIndex={0} role="region" aria-label="project_profile JSON response">
              {JSON_TEXT}
            </pre>
          </div>
        </div>
      </div>
    </Section>
  );
}

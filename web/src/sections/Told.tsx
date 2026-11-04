import { Badge } from '../components/ui';
import { rskNewClaudeMd } from '../data/facts';

export function Told() {
  return (
    <section className="k-sec night" id="told" aria-labelledby="c-told-h">
      <div className="wrap c-told">
        <div className="k-rule" aria-hidden="true" style={{ gridColumn: '1/-1' }} />
        <div className="rv">
          <p className="k-kick">What the agent is told</p>
          <h2 id="c-told-h">rsk new writes the rule into CLAUDE.md.</h2>
          <blockquote>
            <p style={{ margin: 0 }}>“Never install without asking.”</p>
          </blockquote>
        </div>
        <div className="rv">
          <pre tabIndex={0} role="region" aria-label="CLAUDE.md written by rsk new: the rule for agents">
            {rskNewClaudeMd}
          </pre>
          <p style={{ marginTop: 14, color: 'var(--muted)', fontSize: '.92rem' }}>
            <Badge kind="ro">real</Badge> From a real <code>rsk new</code> run. It also registers the server in{' '}
            <code>.mcp.json</code> and <code>opencode.json</code>; <code>rsk destroy</code> removes both and the note.
          </p>
        </div>
      </div>
    </section>
  );
}

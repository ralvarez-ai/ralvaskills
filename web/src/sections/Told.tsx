import { Eyebrow, Section, Tag } from '../components/ui';
import { rskNewClaudeMd } from '../data/facts';

export function Told() {
  return (
    <Section id="told" labelledBy="c-told-h" tone="ink">
      <div className="container quote-block">
        <div className="rule" aria-hidden="true" style={{ gridColumn: '1/-1' }} />
        <div className="reveal">
          <Eyebrow reveal={false}>What the agent is told</Eyebrow>
          <h2 id="c-told-h">rsk new writes the rule into CLAUDE.md.</h2>
          <blockquote>
            <p style={{ margin: 0 }}>“Never install without asking.”</p>
          </blockquote>
        </div>
        <div className="reveal">
          <pre tabIndex={0} role="region" aria-label="CLAUDE.md written by rsk new: the rule for agents">
            {rskNewClaudeMd}
          </pre>
          <p style={{ marginTop: 14, color: 'var(--text-muted)', fontSize: '.92rem' }}>
            <Tag tone="readonly">real</Tag> From a real <code>rsk new</code> run. It also registers the server in{' '}
            <code>.mcp.json</code> and <code>opencode.json</code>; <code>rsk destroy</code> removes both and the note.
          </p>
        </div>
      </div>
    </Section>
  );
}

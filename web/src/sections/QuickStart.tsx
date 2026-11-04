import { CopyBlock, Eyebrow, Section } from '../components/ui';
import { quickStart } from '../data/facts';

export function QuickStart() {
  return (
    <Section id="start" labelledBy="c-start-h">
      <div className="container quick-start">
        <div className="rule" aria-hidden="true" style={{ gridColumn: '1/-1' }} />
        <div className="reveal">
          <Eyebrow reveal={false}>Quick start</Eyebrow>
          <h2 id="c-start-h">Set it up once. Then let the agent ask.</h2>
          <p className="lede">
            <code>rsk new</code> registers the MCP server for both clients. Skip the last two lines if you want the agent
            to propose them.
          </p>
        </div>
        <div className="reveal">
          <CopyBlock id="c-qs" lines={quickStart} label="quick-start commands" />
        </div>
      </div>
    </Section>
  );
}

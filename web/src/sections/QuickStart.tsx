import { CopyBlock } from '../components/ui';
import { quickStart } from '../data/facts';

export function QuickStart() {
  return (
    <section className="k-sec" id="start" aria-labelledby="c-start-h">
      <div className="wrap c-start">
        <div className="k-rule" aria-hidden="true" style={{ gridColumn: '1/-1' }} />
        <div className="rv">
          <p className="k-kick">Quick start</p>
          <h2 id="c-start-h">Set it up once. Then let the agent ask.</h2>
          <p className="k-lede">
            <code>rsk new</code> registers the MCP server for both clients. Skip the last two lines if you want the agent
            to propose them.
          </p>
        </div>
        <div className="rv">
          <CopyBlock id="c-qs" lines={quickStart} label="quick-start commands" />
        </div>
      </div>
    </section>
  );
}

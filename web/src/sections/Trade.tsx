import { Eyebrow, Section } from '../components/ui';
export function Trade() {
  return (
    <Section id="trade" labelledBy="a-trade-h" tone="sunken">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>The trade-off</Eyebrow>
        <h2 id="a-trade-h" className="reveal">
          One binary and one config file for every project. One dependency in exchange.
        </h2>
        <div className="tradeoff">
          <div className="tradeoff__give reveal">
            <h3>What you get</h3>
            <ul>
              <li>
                One binary and one machine config, <code>~/.config/rsk/config.json</code>, for every project.
              </li>
              <li>
                Versions in a lock you commit. A clone runs <code>rsk install</code> and matches.
              </li>
              <li>
                Claude Code and OpenCode from one install: skills are linked once into <code>.claude/skills/</code>.
              </li>
            </ul>
          </div>
          <div className="tradeoff__take reveal">
            <h3>What you accept</h3>
            <ul>
              <li>
                <b>The registry at skills.ralvarez.dev has to be reachable</b> when you install or update, or you point
                rsk at a local clone of the repository.
              </li>
              <li>
                The MCP server needs <code>rsk</code> on PATH on every machine that opens the project.
              </li>
              <li>The catalogue is one author’s standards. To change one, fork the repo and use your clone.</li>
            </ul>
          </div>
        </div>
      </div>
    </Section>
  );
}

import { Badge } from '../components/ui';
import { rskNewClaudeMd } from '../data/facts';

const NEW_FILES: [string, string][] = [
  ['.rsk/rsk.mod', 'The manifest: which tools, which skills, which are pinned. Commit it.'],
  ['.rsk/CLAUDE.md', 'Imports for pinned skills; empty until you pin one. Commit it.'],
  ['CLAUDE.md', 'Starts with @AGENTS.md, then @.rsk/CLAUDE.md, then a short note for the agent.'],
  ['AGENTS.md', 'A skeleton, only when the project has none. An existing one is never touched.'],
  ['.mcp.json', 'Registers rsk mcp for Claude Code.'],
  ['opencode.json', 'Registers rsk mcp for OpenCode.'],
];

const INSTALL_FILES: [string, string][] = [
  ['.rsk/rsk.lock', 'Exact versions and sources. Commit it, like go.sum.'],
  ['.claude/skills/<name>', 'A symlink into rsk’s cache under ~/.config/rsk/. OpenCode reads the same folder.'],
];

function Tree({ rows, label, children }: { rows?: [string, string][]; label: string; children?: React.ReactNode }) {
  return (
    <ul className="a-tree" aria-label={label}>
      {children}
      {rows?.map(([path, note]) => (
        <li key={path}>
          <span className="g" aria-hidden="true">
            +
          </span>
          <span>
            {path}
            <span className="sr"> (new)</span>
          </span>
          <span className="d">{note}</span>
        </li>
      ))}
    </ul>
  );
}

export function Files() {
  return (
    <section className="k-sec" id="files" aria-labelledby="a-files-h">
      <div className="wrap a-two">
        <div className="rv">
          <div className="k-rule" aria-hidden="true" />
          <p className="k-kick">Before and after</p>
          <h2 id="a-files-h">
            What <code>rsk new</code> actually writes.
          </h2>
          <p className="k-lede">
            Six files in the project root and <code>.rsk/</code>, from a real run on a three-file Go project. Installing
            adds the lock and the links. Nothing lands outside the project except the cache.
          </p>
          <p className="a-step">
            rsk new <Badge kind="ro">real run</Badge>
          </p>
          <Tree rows={NEW_FILES} label="Files after rsk new">
            {['go.mod', 'Dockerfile', 'cmd/demo/main.go'].map((path) => (
              <li key={path}>
                <span aria-hidden="true" />
                <span className="yours">{path}</span>
                <span className="d">yours, untouched</span>
              </li>
            ))}
          </Tree>
          <p className="a-step">
            rsk install go-grpc <Badge kind="sa">from the code</Badge>
          </p>
          <Tree rows={INSTALL_FILES} label="Files added by rsk install" />
        </div>
        <div className="rv">
          <div className="k-win">
            <div className="k-win-h">
              CLAUDE.md, as rsk new wrote it <Badge kind="ro">real</Badge>
            </div>
            <pre className="k-win-b" tabIndex={0} role="region" aria-label="CLAUDE.md written by rsk new: what it writes">
              {rskNewClaudeMd}
            </pre>
          </div>
          <p style={{ marginTop: 16, color: 'var(--muted)', fontSize: '.95rem' }}>
            <code>@AGENTS.md</code> stays the first line so the guide you wrote loads before anything rsk generates.{' '}
            <code>rsk destroy</code> removes the note and the MCP registration and leaves <code>AGENTS.md</code> alone.
          </p>
        </div>
      </div>
    </section>
  );
}

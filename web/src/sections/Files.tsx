import { Eyebrow, Panel, Section, Tag } from '../components/ui';
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
    <ul className="file-tree" aria-label={label}>
      {children}
      {rows?.map(([path, note]) => (
        <li key={path}>
          <span className="gutter" aria-hidden="true">
            +
          </span>
          <span>
            {path}
            <span className="sr-only"> (new)</span>
          </span>
          <span className="file-tree__note">{note}</span>
        </li>
      ))}
    </ul>
  );
}

export function Files() {
  return (
    <Section id="files" labelledBy="a-files-h">
      <div className="container manifest-section">
        <div className="reveal">
          <div className="rule" aria-hidden="true" />
          <Eyebrow reveal={false}>Before and after</Eyebrow>
          <h2 id="a-files-h">
            What <code>rsk new</code> actually writes.
          </h2>
          <p className="lede">
            Six files in the project root and <code>.rsk/</code>, from a real run on a three-file Go project. Installing
            adds the lock and the links. Nothing lands outside the project except the cache.
          </p>
          <p className="step-label">
            rsk new <Tag tone="readonly">real run</Tag>
          </p>
          <Tree rows={NEW_FILES} label="Files after rsk new">
            {['go.mod', 'Dockerfile', 'cmd/demo/main.go'].map((path) => (
              <li key={path}>
                <span aria-hidden="true" />
                <span className="file-tree__own">{path}</span>
                <span className="file-tree__note">yours, untouched</span>
              </li>
            ))}
          </Tree>
          <p className="step-label">
            rsk install go-grpc <Tag tone="sample">from the code</Tag>
          </p>
          <Tree rows={INSTALL_FILES} label="Files added by rsk install" />
        </div>
        <div className="reveal">
          <Panel title="CLAUDE.md, as rsk new wrote it" action={<Tag tone="readonly">real</Tag>} flush>
            <pre
              className="panel__body"
              tabIndex={0}
              role="region"
              aria-label="CLAUDE.md written by rsk new: what it writes"
            >
              {rskNewClaudeMd}
            </pre>
          </Panel>
          <p style={{ marginTop: 16, color: 'var(--text-muted)', fontSize: '.95rem' }}>
            <code>@AGENTS.md</code> stays the first line so the guide you wrote loads before anything rsk generates.{' '}
            <code>rsk destroy</code> removes the note and the MCP registration and leaves <code>AGENTS.md</code> alone.
          </p>
        </div>
      </div>
    </Section>
  );
}

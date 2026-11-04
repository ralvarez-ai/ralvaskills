import { tools } from '../data/facts';

const ANNOTATIONS = {
  ro: 'readOnlyHint: true\ndestructiveHint: false\nidempotentHint: true',
  de: 'readOnlyHint: false\ndestructiveHint: true\n_meta: anthropic/\n  requiresUserInteraction',
} as const;

export function Tools() {
  return (
    <section className="k-sec" id="tools" aria-labelledby="c-tools-h" style={{ background: 'var(--bg2)' }}>
      <div className="wrap">
        <div className="k-rule" aria-hidden="true" />
        <p className="k-kick rv">From internal/mcp/tools.go</p>
        <h2 id="c-tools-h" className="rv">
          Three tools. One writes, and it is marked to say so.
        </h2>
        <p className="k-lede rv">
          The descriptions are the ones the agent reads. <code>destructiveHint</code> is only advice to the client; the{' '}
          <code>requiresUserInteraction</code> flag is what makes Claude Code prompt you even in auto-approve modes.
        </p>
        <ul className="c-tools rv">
          {tools.map((tool) => (
            <li className={tool.kind === 'de' ? 'w' : undefined} key={tool.name}>
              <h3>{tool.name}</h3>
              <p>{tool.desc}</p>
              <pre aria-label={`${tool.name} annotations`}>{ANNOTATIONS[tool.kind]}</pre>
            </li>
          ))}
          <li>
            <h3>rsk://catalog</h3>
            <p>
              A resource, not a tool: every skill with its latest version, descriptions cut to 200 characters. Name{' '}
              <code>rsk-catalog</code>.
            </p>
            <pre aria-label="rsk://catalog details">mimeType:{'\n'} application/json</pre>
          </li>
        </ul>
      </div>
    </section>
  );
}

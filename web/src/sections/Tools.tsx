import { Eyebrow, Section } from '../components/ui';
import { tools } from '../data/facts';

const ANNOTATIONS = {
  ro: 'readOnlyHint: true\ndestructiveHint: false\nidempotentHint: true',
  de: 'readOnlyHint: false\ndestructiveHint: true\n_meta: anthropic/\n  requiresUserInteraction',
} as const;

export function Tools() {
  return (
    <Section id="tools" labelledBy="c-tools-h">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <Eyebrow>From internal/mcp/tools.go</Eyebrow>
        <h2 id="c-tools-h" className="reveal">
          Three tools. One writes, and it is marked to say so.
        </h2>
        <p className="lede reveal">
          The descriptions are the ones the agent reads. <code>destructiveHint</code> is only advice to the client; the{' '}
          <code>requiresUserInteraction</code> flag is what makes Claude Code prompt you even in auto-approve modes.
        </p>
        <ul className="tool-list reveal">
          {tools.map((tool) => (
            <li className={tool.kind === 'de' ? 'is-destructive' : undefined} key={tool.name}>
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
            <pre aria-label="rsk://catalog details">{'mimeType:\n  application/json'}</pre>
          </li>
        </ul>
      </div>
    </Section>
  );
}

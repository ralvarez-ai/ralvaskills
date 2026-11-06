import { sources } from '../data/facts';
import { data } from '../data/types';
import { FamilyMark, Tag, Wordmark } from './ui';

export function Footer() {
  return (
    <footer className="site-footer">
      <div className="container">
        <div>
          <Wordmark />
          <p style={{ marginTop: 14 }}>
            rsk {data.version} · GPL-3.0 ·{' '}
            <a href="https://github.com/ralvarezdev/ralvaskills">github.com/ralvarezdev/ralvaskills</a>
          </p>
          <a className="family-credit" href="https://github.com/ralvarez-ai">
            <FamilyMark />
            <span>A ralvarez-ai product</span>
          </a>
        </div>
        <details>
          <summary>Where every number on this page comes from</summary>
          <p>
            Counted in the repository on 2026-11-04. Anything marked <Tag tone="sample">sample</Tag> is illustrative.
          </p>
          <ul className="source-list">
            {sources.map((source) => (
              <li key={source.what}>
                <span>{source.what}</span>
                <code>{source.where}</code>
              </li>
            ))}
          </ul>
        </details>
      </div>
    </footer>
  );
}

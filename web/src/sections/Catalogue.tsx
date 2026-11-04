import { Section } from '../components/ui';
import { useMemo, useState } from 'react';

import { data } from '../data/types';
import type { Skill } from '../data/types';

export function Catalogue() {
  const domains = useMemo(
    () =>
      Object.entries(data.domains).sort(
        ([a, sa], [b, sb]) =>
          Number(a === 'personal') - Number(b === 'personal') || sb.length - sa.length || a.localeCompare(b),
      ),
    [],
  );

  const bundlesOf = useMemo(() => {
    const map = new Map<string, string[]>();
    for (const bundle of data.bundles) {
      for (const name of bundle.inCatalogue) map.set(name, [...(map.get(name) ?? []), bundle.name]);
    }
    return map;
  }, []);

  const notes = useMemo(() => {
    const out: Record<string, string> = {};
    for (const b of data.bundles) {
      let text = `${b.description}. rsk install ${b.name}`;
      if (b.official.length) {
        text += ` — includes ${b.official.length} official Anthropic skill${b.official.length > 1 ? 's' : ''} (${b.official.join(', ')}), fetched from Anthropic and not part of the ${data.skills.length}.`;
      }
      if (b.missing.length) {
        text += ` — also names ${b.missing.join(' and ')}, which are not in the catalogue (flagged in the notes).`;
      }
      out[b.name] = text;
    }
    return out;
  }, []);

  const [bundle, setBundle] = useState('');
  const [includePersonal, setIncludePersonal] = useState(false);

  const matches = (s: Skill) => (!bundle ? true : (bundlesOf.get(s.name) ?? []).includes(bundle));
  const shown = (s: Skill) => (includePersonal || !s.personal) && matches(s);
  const count = data.skills.filter(shown).length;
  const notesFor = (name: string) => `in ${name}, of ${data.skills.length}`;

  const note = bundle
    ? notes[bundle]
    : includePersonal
      ? `Showing all ${data.skills.length} skills, personal ones included.`
      : `Showing the ${count} skills rsk catalog lists by default; the ${data.personal} personal ones need --include personal.`;

  return (
    <Section id="catalogue" labelledBy="b-ex-h" tone="sunken">
      <div className="container">
        <div className="rule" aria-hidden="true" />
        <div className="catalogue__top">
          <div>
            <span className="section-figure" aria-hidden="true" data-count={data.bundles.length}>
              {data.bundles.length}
            </span>
            <h2 id="b-ex-h" style={{ marginTop: 0 }}>
              <span className="sr-only">{data.bundles.length} </span>
              bundles. Filter the catalogue the way <code>rsk catalog</code> does.
            </h2>
            <p className="lede">
              Grouped by the folder they live in. Pick a bundle to see what <code>rsk install &lt;bundle&gt;</code>{' '}
              brings; personal skills stay hidden unless you pass <code>--include personal</code>, as in the CLI.
            </p>
          </div>
          <p style={{ margin: 0, textAlign: 'right' }}>
            <span className="filter-count" aria-hidden="true">
              {count}
            </span>
            <br />
            <span className="filter-count__label" aria-hidden="true">
              {bundle ? notesFor(bundle) : `shown of ${data.skills.length}`}
            </span>
          </p>
        </div>

        <ul className="filter" aria-label="Filter by bundle">
          <li>
            <button className="filter__chip" type="button" aria-pressed={bundle === ''} onClick={() => setBundle('')}>
              all
            </button>
          </li>
          {data.bundles.map((b) => (
            <li key={b.name}>
              <button className="filter__chip" type="button" aria-pressed={bundle === b.name} onClick={() => setBundle(b.name)}>
                {b.name}
              </button>
            </li>
          ))}
          <li>
            <button
              className="filter__chip filter__chip--toggle"
              type="button"
              aria-pressed={includePersonal}
              onClick={() => setIncludePersonal((v) => !v)}
            >
              --include personal
            </button>
          </li>
        </ul>

        <p className="catalogue__note" aria-live="polite">
          {note}
        </p>

        <div className={`catalogue__field${bundle ? ' is-filtered' : ''}${includePersonal ? ' is-inclusive' : ''}`}>
          {domains.map(([domain, skills]) => (
            <section className="domain" aria-label={domain} key={domain}>
              <h3>
                <span>{domain}</span>
                <span>{skills.length}</span>
              </h3>
              <ul>
                {skills.map((s) => (
                  <li
                    className={`skill${s.personal ? ' is-personal' : ''}${bundle && matches(s) ? ' is-match' : ''}`}
                    data-b={(bundlesOf.get(s.name) ?? []).join(' ')}
                    key={s.name}
                  >
                    {s.name}
                    <span className="skill__version">{s.latest}</span>
                  </li>
                ))}
                {domain === 'personal' ? (
                  <li className="domain__placeholder">{skills.length} personal skills, hidden by default</li>
                ) : null}
              </ul>
            </section>
          ))}
        </div>
      </div>
    </Section>
  );
}

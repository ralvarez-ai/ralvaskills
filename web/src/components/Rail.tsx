export interface RailItem {
  id: string;
  label: string;
}

/** Section rail: a tick per section, numbered, active by scroll position. */
export function Rail({ items, active }: { items: RailItem[]; active: string }) {
  return (
    <nav className="section-rail" aria-label="Sections">
      <ol>
        {items.map((item, i) => (
          <li key={item.id}>
            <a href={`#${item.id}`} aria-current={active === item.id ? 'location' : undefined}>
              <span className="section-rail__number">{String(i + 1).padStart(2, '0')}</span>
              <span className="section-rail__label">{item.label}</span>
            </a>
          </li>
        ))}
      </ol>
    </nav>
  );
}

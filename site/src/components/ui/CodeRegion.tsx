import type { ComponentProps } from 'react';

type Props = Omit<ComponentProps<'pre'>, 'role' | 'tabIndex' | 'aria-label'> & {
  /** Accessible name: what this block is, since the text alone does not say. */
  label: string;
};

/**
 * A code block that scrolls sideways, so it is a focusable, named region: a
 * keyboard user can reach it and scroll it (axe's scrollable-region-focusable).
 */
export function CodeRegion({ label, ...rest }: Props) {
  return (
    // biome-ignore lint/a11y/useSemanticElements lint/a11y/noNoninteractiveTabindex: a <pre> cannot become a <section>; the tab stop is the point, it makes the scrollable block reachable by keyboard
    <pre role="region" tabIndex={0} aria-label={label} {...rest} />
  );
}

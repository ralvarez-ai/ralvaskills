/** The ralvarez-apps credit mark: the family hexagon with its square-and-dot cut out. */
export function FamilyMark() {
  return (
    <svg viewBox="0 0 64 64" aria-hidden="true">
      <defs>
        <mask id="family-cut-ai" maskUnits="userSpaceOnUse" x="0" y="0" width="64" height="64">
          <rect width="64" height="64" fill="#fff" />
          <rect x="18" y="18.5" width="12" height="12" rx="2.6" fill="#000" />
          <rect x="34" y="18.5" width="12" height="12" rx="2.6" fill="#000" />
          <rect x="18" y="34.5" width="12" height="12" rx="2.6" fill="#000" />
          <circle cx="40" cy="40.5" r="6" fill="#000" />
        </mask>
      </defs>
      <path
        fill="currentColor"
        mask="url(#family-cut-ai)"
        d="M56.25 16L33.34 2.8L32.88 2.64L32.4 2.55L31.92 2.53L31.12 2.64L30.66 2.8L30.26 3L7.52 16.15L6.91 16.68L6.44 17.33L6.14 18.08L6.02 18.92L6.02 45L6.05 45.44L6.14 45.92L6.3 46.37L6.52 46.81L7.02 47.44L7.38 47.76L7.75 48L30.26 61L30.66 61.2L31.12 61.36L31.92 61.47L32.4 61.45L32.88 61.36L33.34 61.2L33.74 61L56.31 47.96L56.98 47.44L57.48 46.81L57.7 46.37L57.86 45.92L57.98 45.08L57.98 19L57.95 18.56L57.86 18.08L57.56 17.33L57.3 16.93L56.98 16.56L56.62 16.24Z"
      />
    </svg>
  );
}

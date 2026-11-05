import catalog from './catalog.json';

export interface Skill {
  name: string;
  domain: string;
  latest: string;
  releases: number;
  personal: boolean;
  desc: string;
}

export interface Bundle {
  name: string;
  description: string;
  inCatalogue: string[];
  official: string[];
  missing: string[];
}

export interface LockEntry {
  name: string;
  version: string;
  source: string;
}

export interface Catalog {
  version: string;
  skills: Skill[];
  domains: Record<string, Skill[]>;
  bundles: Bundle[];
  releases: number;
  first: string;
  /** [ISO week start, releases] */
  hist: [string, number][];
  latest: Record<string, string>;
  repoMod: string;
  repoLock: LockEntry[];
  repoPinned: string[];
  personal: number;
}

/** Generated from the repository by scripts/gen-data.mjs; never edited by hand. */
export const data = catalog as unknown as Catalog;

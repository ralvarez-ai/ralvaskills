/**
 * Generate src/data/catalog.json from the repository.
 *
 * The page may only state numbers the tree can prove, so the data is derived,
 * never typed by hand. This runs before `dev` and `build` (see package.json) and
 * fails — not warns — when the tree drifts from the anchors the copy relies on:
 *
 *   58 skills · 15 bundles · 12 cobra AddCommand (+ `completion` = 13)
 *   3 MCP tools, exactly one mutating · rsk://catalog present
 *
 * Hand-checked constants that cannot be derived (a real `rsk new` run, the real
 * project_profile response) live in src/data/facts.ts.
 *
 *   pnpm data
 */
import { execFileSync } from 'node:child_process';
import { mkdirSync, readdirSync, readFileSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import { parse as parseToml } from 'smol-toml';

const WEB = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const REPO = resolve(WEB, '..');
const OUT = join(WEB, 'src/data/catalog.json');

const read = (p) => readFileSync(join(REPO, p), 'utf8');

function assert(ok, message) {
  if (!ok) {
    console.error(`gen-data: ${message}`);
    process.exit(1);
  }
}

/** Every folder that directly contains a SKILL.md, as [domain, name]. */
function skillDirs() {
  const out = [];
  const skillsRoot = join(REPO, 'skills');
  for (const domain of readdirSync(skillsRoot)) {
    const domainDir = join(skillsRoot, domain);
    if (!statSync(domainDir).isDirectory()) continue;
    for (const name of readdirSync(domainDir)) {
      const dir = join(domainDir, name);
      if (statSync(dir).isDirectory() && statSync(join(dir, 'SKILL.md')).isFile()) {
        out.push([domain, name]);
      }
    }
  }
  return out.sort((a, b) => a[0].localeCompare(b[0]) || a[1].localeCompare(b[1]));
}

function firstSentence(s) {
  const m = s.trim().match(/^(.+?[.!?])(\s|$)/);
  return (m ? m[1] : s).trim();
}

function main() {
  const index = JSON.parse(read('registry/index.json')).skills;
  const dirs = skillDirs();

  const skills = dirs.map(([domain, name]) => {
    const entry = index[name];
    assert(entry, `${name} is on disk but missing from registry/index.json`);
    return {
      name,
      domain,
      latest: entry.latest,
      releases: Object.keys(entry.versions).length,
      personal: Boolean(entry.personal),
      desc: firstSentence(entry.description),
    };
  });
  assert(
    skills.length === 58 && skills.length === Object.keys(index).length,
    `expected 58 skills, found ${skills.length}`,
  );

  const domains = {};
  for (const s of skills) {
    domains[s.domain] ??= [];
    domains[s.domain].push(s);
  }

  const bundles = parseToml(read('internal/config/catalog.toml')).bundle;
  assert(bundles.length === 15, `expected 15 bundles, found ${bundles.length}`);
  const names = new Set(skills.map((s) => s.name));
  const shaped = bundles.map((b) => ({
    name: b.name,
    description: b.description ?? '',
    inCatalogue: b.skills.filter((s) => names.has(s.name)).map((s) => s.name),
    official: b.skills.filter((s) => s.source === 'official').map((s) => s.name),
    missing: b.skills.filter((s) => s.source !== 'official' && !names.has(s.name)).map((s) => s.name),
  }));

  const setup = read('cmd/rsk/setup.go');
  assert((setup.match(/rootCmd\.AddCommand\(/g) ?? []).length === 12, 'expected 12 AddCommand calls in setup.go');

  const toolsGo = read('internal/mcp/tools.go');
  const tools = [...toolsGo.matchAll(/Name:\s+"([a-z_]+)"/g)].map((m) => m[1]);
  assert(tools.join() === 'project_profile,search_skills,install_skills', `unexpected MCP tools: ${tools.join()}`);
  assert((toolsGo.match(/mutatingAnnotations\("/g) ?? []).length === 1, 'expected exactly one mutating tool');

  const resource = read('internal/mcp/resource.go');
  assert(resource.includes('"rsk://catalog"') && resource.includes('"rsk-catalog"'), 'rsk://catalog is missing');

  const version = execFileSync('git', ['describe', '--tags', '--abbrev=0', '--match', 'v*'], { cwd: REPO })
    .toString()
    .trim()
    .replace(/^v/, '');

  const stamps = Object.values(index).flatMap((e) => Object.values(e.versions).map((v) => new Date(v.published_at)));
  const first = new Date(Math.min(...stamps.map((d) => d.getTime())));
  const week0 = new Date(first);
  week0.setUTCDate(week0.getUTCDate() - ((week0.getUTCDay() + 6) % 7));
  week0.setUTCHours(0, 0, 0, 0);
  const weeks = new Map();
  for (const d of stamps) {
    const i = Math.floor((d - week0) / (7 * 24 * 3600 * 1000));
    weeks.set(i, (weeks.get(i) ?? 0) + 1);
  }
  const nweeks = Math.max(...weeks.keys()) + 1;
  const hist = Array.from({ length: nweeks }, (_, i) => {
    const day = new Date(week0.getTime() + i * 7 * 24 * 3600 * 1000);
    return [day.toISOString().slice(0, 10), weeks.get(i) ?? 0];
  });

  const mod = read('.rsk/rsk.mod');
  const lock = parseToml(read('.rsk/rsk.lock')).skills;

  const catalog = {
    version,
    skills,
    domains,
    bundles: shaped,
    releases: stamps.length,
    first: first.toISOString().slice(0, 10),
    hist,
    latest: Object.fromEntries(skills.map((s) => [s.name, s.latest])),
    repoMod: mod,
    repoLock: lock,
    repoPinned: parseToml(mod).pinned,
    personal: skills.filter((s) => s.personal).length,
  };

  mkdirSync(dirname(OUT), { recursive: true });
  writeFileSync(OUT, `${JSON.stringify(catalog, null, 1)}\n`);
  console.log(
    `catalog.json <- ${skills.length} skills, ${shaped.length} bundles, ${stamps.length} releases, ` +
      `${Object.keys(domains).length} domains (rsk ${version})`,
  );
}

main();

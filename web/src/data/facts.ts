/**
 * Checked facts that cannot be derived from the tree.
 *
 * Every entry here came from a real run or a real file on 2026-11-04 and is
 * labelled on the page as such (real run / real call / from the code). Numbers
 * that *can* be derived live in catalog.json instead.
 */

export const rskNewOutput = [
  '✓  registered the rsk MCP server (.mcp.json, opencode.json)',
  '✓  initialized .rsk/ for all',
] as const;

/** `rsk new` on a scratch Go project; the block it appends to CLAUDE.md. */
export const rskNewClaudeMd = `@AGENTS.md
@.rsk/CLAUDE.md

<!-- rsk-mcp:start -->
## Project standards

Before scaffolding or a non-trivial change, call the rsk MCP \`project_profile\` tool to learn which skills apply here. Installed skills appear automatically in your skill list — invoke them with the Skill tool. Never install without asking.
<!-- rsk-mcp:end -->`;

export interface ProfileSignal {
  kind: string;
  proof: string;
}
export interface ProfileCandidate {
  because: string[];
  installed: boolean;
  latest: string;
  skill: string;
}
export interface Profile {
  candidates: ProfileCandidate[];
  installed: unknown[];
  manifest: { exists: boolean; pinned: number; skills: number };
  project_root: string;
  signals: ProfileSignal[];
}

export const profile: Profile = {
  candidates: [
    { because: ['language:go'], installed: false, latest: '1.14.0', skill: 'go-architect' },
    { because: ['infra:docker'], installed: false, latest: '1.3.0', skill: 'docker-architect' },
    { because: ['repo:cli'], installed: false, latest: '1.0.0', skill: 'cli-tool-architect' },
  ],
  installed: [],
  manifest: { exists: false, pinned: 0, skills: 0 },
  project_root: '/tmp/rsk-demo-go',
  signals: [
    { kind: 'language:go', proof: 'go.mod' },
    { kind: 'infra:docker', proof: 'Dockerfile' },
    { kind: 'repo:cli', proof: 'cmd/demo/main.go' },
  ],
};

/** The skills the profile above proposes, in order. */
export const profileNames = profile.candidates.map((c) => c.skill);

export interface Signal {
  kind: string;
  proof: string;
  proposes: string;
}

/** internal/mcp/signals.go — the signal kinds, their proof, the skills proposed. */
export const signals: Signal[] = [
  { kind: 'language:go', proof: 'go.mod', proposes: 'go-architect' },
  { kind: 'language:python', proof: 'pyproject.toml, requirements.txt, uv.lock', proposes: 'python-architect' },
  { kind: 'language:typescript', proof: 'tsconfig.json, package.json', proposes: '—' },
  {
    kind: 'framework:nextjs',
    proof: 'next.config.*, or next in package.json',
    proposes: 'nextjs-architect, react-architect',
  },
  { kind: 'framework:astro', proof: 'astro.config.*', proposes: 'astro-architect' },
  {
    kind: 'protocol:grpc',
    proof: 'any .proto, or google.golang.org/grpc in go.mod',
    proposes: 'grpc-architect, protobuf-architect',
  },
  { kind: 'protocol:rest', proof: 'openapi.yaml/.yml, swagger.json, any .http', proposes: 'rest-api-architect' },
  { kind: 'infra:docker', proof: 'Dockerfile, docker-compose.yml/.yaml', proposes: 'docker-architect' },
  { kind: 'infra:k8s', proof: 'Chart.yaml, kustomization.yaml/.yml', proposes: '—' },
  { kind: 'infra:ci', proof: 'files in .github/workflows', proposes: 'ci-cd-architect' },
  { kind: 'repo:tooling', proof: '.golangci.yml, mise.toml, Taskfile.yml', proposes: 'repo-tooling-architect' },
  { kind: 'repo:cli', proof: 'main.go, or cmd/*/main.go', proposes: 'cli-tool-architect' },
  { kind: 'data:postgres', proof: '.sql files with sqlx or psycopg', proposes: '—' },
];

export interface HelpGroup {
  group: string;
  commands: { name: string; desc: string }[];
}

/** `rsk --help` groups (cmd/rsk/setup.go; bin/rsk --help on 3.4.1). */
export const helpGroups: HelpGroup[] = [
  { group: 'Machine setup', commands: [{ name: 'init', desc: 'Set up rsk for this machine.' }] },
  {
    group: 'Project lifecycle',
    commands: [
      { name: 'new', desc: 'Initialize an rsk project in the current directory.' },
      { name: 'destroy', desc: 'Remove the rsk project from the current directory.' },
    ],
  },
  {
    group: 'Install / uninstall / update',
    commands: [
      { name: 'install', desc: 'Install bundles or skills.' },
      { name: 'uninstall', desc: 'Remove installed bundles or skills.' },
      { name: 'update', desc: 'Pull the latest skills and re-link.' },
    ],
  },
  {
    group: 'Project pinning',
    commands: [{ name: 'pin', desc: "Pin a skill so it is imported into each tool's project config." }],
  },
  {
    group: 'Tool configuration',
    commands: [{ name: 'tools', desc: 'Manage tool permissions (allow/deny/list/remove).' }],
  },
  {
    group: 'Read-only views',
    commands: [
      { name: 'catalog', desc: 'Browse the catalog of available skills and bundles.' },
      { name: 'list', desc: 'Show installed skills (project / --global).' },
      { name: 'status', desc: 'Show installed skills and version drift.' },
    ],
  },
  {
    group: 'Other',
    commands: [
      { name: 'mcp', desc: 'Run the MCP server over stdio.' },
      { name: 'completion', desc: 'Generate the autocompletion script for the specified shell.' },
    ],
  },
];

export type ToolKind = 'ro' | 'de';

export interface Tool {
  name: string;
  desc: string;
  kind: ToolKind;
}

/** internal/mcp/tools.go — the agent-facing descriptions, verbatim. */
export const tools: Tool[] = [
  {
    name: 'project_profile',
    desc: "Detect this project's stack from its files and propose the ralvaskills skills that apply. Call it before scaffolding or a non-trivial change.",
    kind: 'ro',
  },
  {
    name: 'search_skills',
    desc: 'Search the ralvaskills catalog by keyword. Use it to find a skill by name or intent when project_profile does not name one.',
    kind: 'ro',
  },
  {
    name: 'install_skills',
    desc: 'Install one or more skills or bundles into this project by name, resolving them through the ralvaskills catalog. Writes symlinks and the project manifest. The client asks the user to approve.',
    kind: 'de',
  },
];

export const quickStart: { cmd: string; note: string }[] = [
  { cmd: 'go install github.com/ralvarezdev/ralvaskills/v3/cmd/rsk@latest', note: 'or a binary from Releases' },
  { cmd: 'rsk init', note: 'once per machine' },
  { cmd: 'rsk new', note: 'in a project: writes .rsk/ and registers rsk mcp' },
  { cmd: 'rsk install go-grpc', note: 'a bundle or a skill' },
  { cmd: 'rsk pin go-architect', note: 'load it in every session' },
];

export const sources: { what: string; where: string }[] = [
  { what: '58 skills', where: 'registry/index.json · skills/*/*/SKILL.md' },
  { what: '17 domains', where: 'skills/<domain>/' },
  { what: '15 bundles', where: 'internal/config/catalog.toml' },
  { what: '13 commands', where: 'cmd/rsk/setup.go (+ cobra completion)' },
  { what: '3 MCP tools, 1 writes', where: 'internal/mcp/tools.go' },
  { what: 'rsk://catalog', where: 'internal/mcp/resource.go' },
  { what: 'Detection rules', where: 'internal/mcp/signals.go' },
  { what: 'Releases and dates', where: 'registry/index.json versions[].published_at' },
  { what: 'Version', where: 'git tag v3.4.1 · bin/rsk --version' },
  { what: 'Config and cache', where: 'internal/config/config.go' },
  { what: 'What install writes', where: 'internal/install/apply.go' },
  { what: 'rsk new output', where: 'real run, scratch folder, 2026-11-04' },
  { what: 'project_profile JSON', where: 'real call, scratch folder, 2026-11-04' },
  { what: 'Colours', where: 'termkit TokyoNight via internal/ui/theme.go' },
  { what: 'Wordmark', where: 'internal/ui/theme.go' },
];

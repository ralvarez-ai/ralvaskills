#!/usr/bin/env bash
# Read-only token audit for Claude Code and OpenCode.
# Usage: audit.sh [project-dir]   Prints OK / GAP / WARN / INFO lines; never changes anything.
set -uo pipefail

proj="$(cd "${1:-.}" && pwd)"
oc_global="${XDG_CONFIG_HOME:-$HOME/.config}/opencode"
claude_dir="$HOME/.claude"

say() { printf '%-5s %-28s %s\n' "$1" "$2" "$3"; }
has() { command -v "$1" >/dev/null 2>&1; }
# grep_any PATTERN FILE...  -> true if any existing file matches
grep_any() {
  local pat="$1"; shift
  local f
  for f in "$@"; do [ -e "$f" ] && grep -rqs -- "$pat" "$f" && return 0; done
  return 1
}
oc_files=("$proj/opencode.json" "$proj/opencode.jsonc" "$proj/.opencode" "$oc_global")
cc_files=("$claude_dir/settings.json" "$claude_dir/settings.local.json" "$HOME/.claude.json"
          "$proj/.claude/settings.json" "$proj/.claude/settings.local.json" "$proj/.mcp.json")

echo "== tools  ($proj)"
has claude && say INFO "claude" "$(claude --version 2>/dev/null | head -1)"
has opencode && say INFO "opencode" "$(opencode --version 2>/dev/null | head -1)"

# context-mode
if grep_any 'context-mode' "$claude_dir/plugins/installed_plugins.json" "${cc_files[@]}"; then
  say OK "context-mode/claude" "plugin or MCP entry found"
else
  say GAP "context-mode/claude" "not installed (Claude Code plugin marketplace)"
fi
if has opencode; then
  if grep_any 'context-mode' "${oc_files[@]}"; then say OK "context-mode/opencode" "configured"
  else say GAP "context-mode/opencode" "run: opencode plugin add context-mode"; fi
fi

# RTK
if has rtk && rtk gain >/dev/null 2>&1; then
  say OK "rtk" "$(rtk --version 2>/dev/null | head -1)"
  if grep_any 'rtk' "${cc_files[@]}"; then say OK "rtk/claude" "hook found"
  else say GAP "rtk/claude" "run: rtk init -g"; fi
  if has opencode; then
    if grep_any 'rtk' "${oc_files[@]}"; then say OK "rtk/opencode" "plugin found"
    else say GAP "rtk/opencode" "run: rtk init -g --opencode"; fi
  fi
elif has rtk; then
  say WARN "rtk" "binary present but 'rtk gain' fails: likely the unrelated crates.io rtk"
else
  say GAP "rtk" "not installed (github.com/rtk-ai/rtk)"
fi

# CodeGraph
if has codegraph; then
  say OK "codegraph" "$(codegraph --version 2>/dev/null | head -1)"
  if grep_any 'codegraph' "${cc_files[@]}"; then say OK "codegraph/claude" "MCP entry found"
  else say GAP "codegraph/claude" "run: codegraph install"; fi
  if grep -qs '"codegraph"' "$HOME/.claude.json" && ! grep -qs 'alwaysLoad' "$HOME/.claude.json"; then
    say WARN "codegraph/claude" "no alwaysLoad: tool may be deferred behind tool search"
  fi
  if has opencode; then
    if grep_any 'codegraph' "${oc_files[@]}"; then say OK "codegraph/opencode" "MCP entry found"
    else say GAP "codegraph/opencode" "run: codegraph install"; fi
  fi
  if [ -d "$proj/.codegraph" ]; then say OK "codegraph/index" ".codegraph/ present"
  else say GAP "codegraph/index" "run in project: codegraph init"; fi
else
  say GAP "codegraph" "not installed (npm i -g @colbymchenry/codegraph)"
fi

# Explore override
explore=""
for f in "$proj/.claude/agents/Explore.md" "$claude_dir/agents/Explore.md"; do
  [ -f "$f" ] && explore="$f" && break
done
if [ -n "$explore" ]; then
  say OK "explore-agent" "$(grep -m1 '^model:' "$explore" || echo 'no model: line') ($explore)"
else
  say GAP "explore-agent" "built-in Explore runs on the main model; add Explore.md with model: haiku"
fi

echo "== always-loaded instructions"
total=0
declare -A seen=()
# walk FILE DEPTH: count bytes, follow whole-line @imports (relative to the importing file)
walk() {
  local file="$1" depth="$2" line target bytes
  [ -f "$file" ] || return
  file="$(realpath "$file")"
  [ -n "${seen[$file]:-}" ] && return
  seen[$file]=1
  bytes=$(wc -c <"$file")
  total=$((total + bytes))
  say INFO "~$((bytes / 4)) tok" "${file#"$proj"/}"
  case "$file" in *SKILL.md) say WARN "skill body at start" "${file#"$proj"/} (use the skill index instead)";; esac
  [ "$depth" -ge 5 ] && return
  while IFS= read -r line; do
    target="${line#@}"; target="${target%%[[:space:]]*}"
    case "$target" in \~/*) target="$HOME/${target#\~/}";; /*) ;; *) target="$(dirname "$file")/$target";; esac
    walk "$target" $((depth + 1))
  done < <(grep -E '^@[^[:space:]]+' "$file")
}
walk "$claude_dir/CLAUDE.md" 0
walk "$proj/CLAUDE.md" 0
walk "$proj/AGENTS.md" 0
for cfg in "$proj/opencode.json" "$oc_global/opencode.json"; do
  { [ -f "$cfg" ] && has python3; } || continue
  while IFS= read -r inst; do
    [ -n "$inst" ] && walk "$proj/$inst" 1
  done < <(python3 -c 'import json,sys; [print(i) for i in json.load(open(sys.argv[1])).get("instructions",[])]' "$cfg" 2>/dev/null)
done
tok=$((total / 4))
if [ "$tok" -gt 4000 ]; then say WARN "always-loaded total" "~$tok tok per session start"
else say OK "always-loaded total" "~$tok tok per session start"; fi

echo "== MCP servers"
if has python3; then
  python3 - "$proj" "$HOME/.claude.json" <<'PY'
import json, os, sys
proj, cj = sys.argv[1], sys.argv[2]
def load(p):
    try:
        with open(p) as f: return json.load(f)
    except Exception: return {}
names = set(load(os.path.join(proj, ".mcp.json")).get("mcpServers", {}))
c = load(cj)
names |= set(c.get("mcpServers", {}))
names |= set(c.get("projects", {}).get(proj, {}).get("mcpServers", {}))
oc = set(load(os.path.join(proj, "opencode.json")).get("mcp", {}))
print(f"INFO  claude MCP servers         {len(names)}: {', '.join(sorted(names)) or '-'} (plugins not counted)")
print(f"INFO  opencode MCP servers       {len(oc)}: {', '.join(sorted(oc)) or '-'}")
if len(oc) > 5:
    print("WARN  opencode MCP servers       many servers: every tool schema loads each turn")
PY
fi
exit 0

package tool

import (
	"slices"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// actionShell is the OpenCode action name for command execution.
	actionShell = "shell"

	// actionRead, actionEdit and actionWebFetch mirror OpenCode's V2 actions.
	actionRead     = "read"
	actionEdit     = "edit"
	actionWebFetch = "webfetch"

	// effectAllow is the only effect rsk writes. A skill never restricts the
	// user's environment, so it never grants a deny.
	effectAllow = "allow"

	// claudeSkillDirVar is substituted by Claude Code in both the skill body
	// and the frontmatter rule. Clients that do not expand it need a glob.
	claudeSkillDirVar = "${CLAUDE_SKILL_DIR}"

	// claudeProjectDirVar is the other variable Claude Code substitutes.
	claudeProjectDirVar = "${CLAUDE_PROJECT_DIR}"
)

// actionByTool maps a declared tool name to the client-neutral action.
var actionByTool = map[string]string{
	"bash":     actionShell,
	"read":     actionRead,
	"write":    actionEdit,
	"edit":     actionEdit,
	"webfetch": actionWebFetch,
}

// Rule is a tool permission in client-neutral form: an action on a resource.
// An adapter translates it into its own config shape.
type Rule struct {
	Action   string
	Resource string
}

// Key returns a stable "action:resource" identity for bookkeeping in the lock
// and for comparing rules across a version bump.
func (r Rule) Key() string {
	return r.Action + ":" + r.Resource
}

// PermissionWriter is implemented by adapters whose client keeps tool
// permissions in a config file rsk can write. Clients that read the skill's
// allowed-tools frontmatter (Claude Code) do not implement it.
type PermissionWriter interface {
	// GrantPermissions appends the rules the config does not already carry and
	// returns the subset it wrote. Existing entries keep their order and their
	// meaning, and a rule the user added by hand is never reported as written.
	GrantPermissions(projectDir string, rules []Rule) ([]Rule, error)

	// RevokePermissions removes exactly rules and leaves every other entry
	// alone. It returns nil when nothing matched.
	RevokePermissions(projectDir string, rules []Rule) error
}

// WriterFor returns the permission writer for id, and whether the adapter has
// one.
//
//nolint:ireturn // returns the interface the adapter implements
func WriterFor(id ID) (PermissionWriter, bool) {
	t, ok := Get(id)
	if !ok {
		return nil, false
	}
	writer, ok := t.(PermissionWriter)
	return writer, ok
}

// Translate converts a declared Tool(specifier) entry into a client-neutral
// rule. The second return is false when the client has no equivalent, such as
// an MCP tool reference.
func Translate(entry string) (Rule, bool) {
	toolName, spec, ok := skill.SplitRule(entry)
	if !ok {
		return Rule{}, false
	}
	action, ok := actionByTool[strings.ToLower(toolName)]
	if !ok {
		return Rule{}, false
	}
	spec = strings.ReplaceAll(spec, claudeSkillDirVar, "*")
	spec = strings.ReplaceAll(spec, claudeProjectDirVar, "*")
	if action == actionWebFetch {
		spec = strings.TrimPrefix(spec, "domain:")
	}
	return Rule{Action: action, Resource: spec}, true
}

// SortRules returns the rules sorted by action and resource, with duplicates
// removed.
func SortRules(rules []Rule) []Rule {
	seen := make(map[string]struct{}, len(rules))
	kept := make([]Rule, 0, len(rules))
	for _, rule := range rules {
		if _, dup := seen[rule.Key()]; dup {
			continue
		}
		seen[rule.Key()] = struct{}{}
		kept = append(kept, rule)
	}
	slices.SortFunc(kept, func(a, b Rule) int {
		if a.Action != b.Action {
			return strings.Compare(a.Action, b.Action)
		}
		return strings.Compare(a.Resource, b.Resource)
	})
	return kept
}

// ParseRules turns keys back into rules. It is the inverse of Rule.Key and is
// used to rebuild a revocation set from the lock.
func ParseRules(keys []string) []Rule {
	rules := make([]Rule, 0, len(keys))
	for _, key := range keys {
		action, resource, ok := strings.Cut(key, ":")
		if !ok || action == "" || resource == "" {
			continue
		}
		rules = append(rules, Rule{Action: action, Resource: resource})
	}
	return SortRules(rules)
}

// RuleKeys returns the Key of every rule, sorted, for storage in the lock.
func RuleKeys(rules []Rule) []string {
	sorted := SortRules(rules)
	if len(sorted) == 0 {
		return nil
	}
	keys := make([]string, len(sorted))
	for i, rule := range sorted {
		keys[i] = rule.Key()
	}
	return keys
}

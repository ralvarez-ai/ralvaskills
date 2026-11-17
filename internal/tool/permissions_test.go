package tool

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestTranslate(t *testing.T) {
	t.Parallel()
	cases := map[string]Rule{
		"Bash(git status *)":                         {Action: actionShell, Resource: "git status *"},
		"Bash(${CLAUDE_SKILL_DIR}/scripts/notify *)": {Action: actionShell, Resource: "*/scripts/notify *"},
		"Read(~/docs/**)":                            {Action: actionRead, Resource: "~/docs/**"},
		"Write(~/out/**)":                            {Action: actionEdit, Resource: "~/out/**"},
		"Edit(~/out/**)":                             {Action: actionEdit, Resource: "~/out/**"},
		"WebFetch(domain:example.com)":               {Action: actionWebFetch, Resource: "example.com"},
		"Bash(${CLAUDE_PROJECT_DIR}/scripts/task *)": {Action: actionShell, Resource: "*/scripts/task *"},
	}
	for entry, want := range cases {
		got, ok := Translate(entry)
		if !ok {
			t.Errorf("Translate(%q) ok = false", entry)
			continue
		}
		if got != want {
			t.Errorf("Translate(%q) = %+v, want %+v", entry, got, want)
		}
	}

	for _, entry := range []string{"mcp__rsk__search_skills", "Bash", "Unknown(x)"} {
		if _, ok := Translate(entry); ok {
			t.Errorf("Translate(%q) ok = true, want false", entry)
		}
	}
}

func TestRuleKeysRoundTrip(t *testing.T) {
	t.Parallel()
	rules := []Rule{
		{Action: actionShell, Resource: "git status *"},
		{Action: actionRead, Resource: "~/docs/**"},
	}
	keys := RuleKeys(rules)
	if !slices.Equal(keys, []string{"read:~/docs/**", "shell:git status *"}) {
		t.Fatalf("RuleKeys = %q", keys)
	}
	if got := ParseRules(keys); !slices.Equal(got, []Rule{rules[1], rules[0]}) {
		t.Errorf("ParseRules = %+v, want the same rules back", got)
	}
}

// writeOpenCodeConfig seeds a project opencode.json and returns its directory.
func writeOpenCodeConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, OpenCodeFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write opencode.json: %v", err)
	}
	return dir
}

// readPermissionRules returns every permission entry, in file order.
func readPermissionRules(t *testing.T, dir string) []any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, OpenCodeFileName))
	if err != nil {
		t.Fatalf("read opencode.json: %v", err)
	}
	cfg := make(map[string]any)
	if err = json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse opencode.json: %v", err)
	}
	list, _ := cfg[OpenCodePermissionsKey].([]any)
	return list
}

func TestOpenCodeGrantPermissions(t *testing.T) {
	t.Parallel()
	const userRule = `{"action":"shell","resource":"git *","effect":"ask"}`
	dir := writeOpenCodeConfig(t, `{"permissions":[`+userRule+`],"instructions":["keep.md"]}`)

	oc := &openCodeTool{}
	rules := []Rule{
		{Action: actionShell, Resource: "*/scripts/notify *"},
		{Action: actionRead, Resource: "~/docs/**"},
	}
	added, err := oc.GrantPermissions(dir, rules)
	if err != nil {
		t.Fatalf("GrantPermissions: %v", err)
	}
	if len(added) != 2 {
		t.Fatalf("added = %+v, want both rules", added)
	}
	if added, err = oc.GrantPermissions(dir, rules); err != nil {
		t.Fatalf("GrantPermissions (second run): %v", err)
	}
	if len(added) != 0 {
		t.Fatalf("second run added = %+v, want none", added)
	}

	got := readPermissionRules(t, dir)
	if len(got) != 3 {
		t.Fatalf("permissions = %v, want the user rule plus two rsk rules", got)
	}
	first, _ := got[0].(map[string]any)
	if first["resource"] != "git *" || first["effect"] != "ask" {
		t.Errorf("the user rule was reordered or rewritten: %v", first)
	}
	last, _ := got[2].(map[string]any)
	if last["resource"] != "~/docs/**" || last["effect"] != effectAllow {
		t.Errorf("new rules must be appended (last match wins): %v", last)
	}
}

func TestOpenCodeGrantPermissionsSkipsExistingRule(t *testing.T) {
	t.Parallel()
	const ours = `{"action":"shell","resource":"*/scripts/notify *","effect":"allow"}`
	dir := writeOpenCodeConfig(t, `{"permissions":[`+ours+`]}`)

	added, err := (&openCodeTool{}).GrantPermissions(dir, []Rule{
		{Action: actionShell, Resource: "*/scripts/notify *"},
	})
	if err != nil {
		t.Fatalf("GrantPermissions: %v", err)
	}
	if len(added) != 0 {
		t.Fatalf("added = %+v, want none: a rule already in the file is not rsk's to revoke", added)
	}
	if got := readPermissionRules(t, dir); len(got) != 1 {
		t.Errorf("permissions = %v, want the single existing entry", got)
	}
}

func TestOpenCodeRevokePermissions(t *testing.T) {
	t.Parallel()
	dir := writeOpenCodeConfig(t, `{"permissions":[
		{"action":"shell","resource":"git *","effect":"ask"},
		{"action":"shell","resource":"*/scripts/notify *","effect":"allow"},
		{"action":"read","resource":"~/docs/**","effect":"allow"}
	],"instructions":["keep.md"]}`)

	oc := &openCodeTool{}
	if err := oc.RevokePermissions(dir, []Rule{{Action: actionShell, Resource: "*/scripts/notify *"}}); err != nil {
		t.Fatalf("RevokePermissions: %v", err)
	}

	got := readPermissionRules(t, dir)
	if len(got) != 2 {
		t.Fatalf("permissions = %v, want the unrelated entries kept", got)
	}
	for _, entry := range got {
		obj, _ := entry.(map[string]any)
		if obj["resource"] == "*/scripts/notify *" {
			t.Errorf("revoked rule still present: %v", obj)
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, OpenCodeFileName))
	if err != nil {
		t.Fatal(err)
	}
	cfg := make(map[string]any)
	if err = json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg[OpenCodeInstructionsKey] == nil {
		t.Error("revoking a permission dropped an unrelated key")
	}
}

func TestOpenCodeRevokePermissionsNoMatchWritesNothing(t *testing.T) {
	t.Parallel()
	const seed = `{"permissions":[{"action":"shell","resource":"git *","effect":"ask"}]}`
	dir := writeOpenCodeConfig(t, seed)

	if err := (&openCodeTool{}).RevokePermissions(dir, []Rule{{Action: actionShell, Resource: "other *"}}); err != nil {
		t.Fatalf("RevokePermissions: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, OpenCodeFileName))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != seed {
		t.Errorf("file changed with no matching rule:\n%s", data)
	}
}

func TestWriterFor(t *testing.T) {
	t.Parallel()
	if _, ok := WriterFor(OpenCodeID); !ok {
		t.Error("OpenCode must expose a permission writer")
	}
	if _, ok := WriterFor(ClaudeID); ok {
		t.Error("Claude Code applies the frontmatter itself and must not expose a writer")
	}
}

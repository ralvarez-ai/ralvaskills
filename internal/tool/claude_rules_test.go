package tool

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// claudeProject returns a project directory whose .claude folder exists, as it
// does in a real project.
func claudeProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ClaudeFolderName), 0o750); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestClaudeToolStatesAllowBeatsDeny(t *testing.T) {
	t.Parallel()
	states := ClaudeToolStates([]string{"Bash", "Read"}, []string{"Read", "Write"})
	got := map[string]Permission{}
	for _, s := range states {
		got[s.Name] = s.Permission
	}
	want := map[string]Permission{
		"Bash":  Allowed,
		"Read":  Allowed,
		"Write": Denied,
		"Edit":  Unconfigured,
	}
	for name, perm := range want {
		if got[name] != perm {
			t.Errorf("%s = %v, want %v", name, got[name], perm)
		}
	}
	if len(states) != len(AvailableClaudeTools()) {
		t.Errorf("states = %d, tools = %d", len(states), len(AvailableClaudeTools()))
	}
}

func TestAllowAndDenyMoveARuleBetweenTheLists(t *testing.T) {
	t.Parallel()
	dir, ct := claudeProject(t), &ClaudeTool{}

	changed, err := ct.DenyRule(dir, "Bash")
	if err != nil || !changed {
		t.Fatalf("deny: changed=%v err=%v", changed, err)
	}
	if changed, err = ct.DenyRule(dir, "Bash"); err != nil || changed {
		t.Fatalf("second deny: changed=%v err=%v, want unchanged", changed, err)
	}

	if changed, err = ct.AllowRule(dir, "Bash"); err != nil || !changed {
		t.Fatalf("allow: changed=%v err=%v", changed, err)
	}
	allow, deny, err := ct.ReadPermissions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(allow, []string{"Bash"}) || len(deny) != 0 {
		t.Errorf("allow=%v deny=%v, want Bash moved to allow", allow, deny)
	}
}

func TestRemoveRuleReportsWhetherAnythingChanged(t *testing.T) {
	t.Parallel()
	dir, ct := claudeProject(t), &ClaudeTool{}
	if changed, err := ct.RemoveRule(dir, "Bash"); err != nil || changed {
		t.Fatalf("remove on empty: changed=%v err=%v, want unchanged", changed, err)
	}
	if _, err := ct.AllowRule(dir, "Bash"); err != nil {
		t.Fatal(err)
	}
	if changed, err := ct.RemoveRule(dir, "Bash"); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if allow, deny, err := ct.ReadPermissions(dir); err != nil || len(allow)+len(deny) != 0 {
		t.Errorf("after remove: allow=%v deny=%v err=%v", allow, deny, err)
	}
}

package permission

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

const notifyRule = "Bash(${CLAUDE_SKILL_DIR}/scripts/notify *)"

// project creates a temp project root with an empty opencode.json.
func project(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, tool.OpenCodeFileName), []byte("{}"), 0o600); err != nil {
		t.Fatalf("seed opencode.json: %v", err)
	}
	return root
}

// declared returns a skill that declares entries.
func declared(name, version string, entries ...string) skill.Skill {
	return skill.Skill{Name: name, Version: version, AllowedTools: entries}
}

// lockWith returns a lock holding one entry.
func lockWith(permissions, granted []string) manifest.Lock {
	return manifest.Lock{Skills: []manifest.LockEntry{{
		Name:        "notify",
		Version:     "0.1.0",
		Permissions: permissions,
		Granted:     granted,
	}}}
}

// openCodeRules returns the rules opencode.json allows, in file order.
func openCodeRules(t *testing.T, root string) []tool.Rule {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, tool.OpenCodeFileName))
	if err != nil {
		t.Fatalf("read opencode.json: %v", err)
	}
	writer, _ := tool.WriterFor(tool.OpenCodeID)
	if writer == nil {
		t.Fatal("no OpenCode writer")
	}
	// GrantPermissions is the only reader exposed; revoking nothing would
	// rewrite the file, so read the file shape directly instead.
	return rulesFromJSON(t, data)
}

func baseOptions(root string, skills ...skill.Skill) Options {
	return Options{
		ProjectRoot: root,
		Tools:       []tool.ID{tool.OpenCodeID},
		Skills:      skills,
		Mode:        ModeApply,
	}
}

func TestReconcileFirstInstallWritesAndRecords(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith(nil, nil)

	next, result, err := Reconcile(lock, baseOptions(root, declared("notify", "0.1.0", notifyRule)))
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if len(result.Added) != 1 {
		t.Fatalf("added = %+v, want one rule", result.Added)
	}
	if got := openCodeRules(t, root); len(got) != 1 || got[0].Action != "shell" {
		t.Fatalf("opencode.json rules = %+v, want the shell rule", got)
	}

	entry, ok := manifest.FindLockEntry(next, "notify")
	if !ok {
		t.Fatal("lock entry missing")
	}
	if !slices.Equal(entry.Permissions, []string{notifyRule}) {
		t.Errorf("lock permissions = %q", entry.Permissions)
	}
	if len(entry.Granted) != 1 {
		t.Errorf("lock granted = %q, want the written rule", entry.Granted)
	}
}

func TestReconcileIsIdempotent(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith([]string{notifyRule}, []string{"shell:*/scripts/notify *"})

	next, result, err := Reconcile(lock, baseOptions(root, declared("notify", "0.1.0", notifyRule)))
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if !result.Empty() {
		t.Errorf("result = %+v, want nothing to do", result)
	}
	if got := openCodeRules(t, root); len(got) != 0 {
		t.Errorf("opencode.json rules = %+v, want none written", got)
	}
	if next.Skills[0].Granted[0] != "shell:*/scripts/notify *" {
		t.Errorf("reconcile dropped an existing grant: %+v", next.Skills[0])
	}
}

func TestReconcileNewVersionAddsEntry(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith([]string{notifyRule}, []string{"shell:*/scripts/notify *"})

	next, result, err := Reconcile(lock,
		baseOptions(root, declared("notify", "0.2.0", notifyRule, "Read(~/docs/**)")))
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if len(result.Added) != 1 || result.Added[0].Action != "read" {
		t.Fatalf("added = %+v, want the new read rule", result.Added)
	}
	if got := openCodeRules(t, root); len(got) != 1 || got[0].Action != "read" {
		t.Errorf("opencode.json rules = %+v, want only the new rule written", got)
	}
	entry, _ := manifest.FindLockEntry(next, "notify")
	if len(entry.Permissions) != 2 || len(entry.Granted) != 2 {
		t.Errorf("lock = %+v, want both entries recorded", entry)
	}
}

func TestReconcileRemovalIsOfferedNotApplied(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith([]string{notifyRule}, []string{"shell:*/scripts/notify *"})
	seedRules(t, root, tool.Rule{Action: "shell", Resource: "*/scripts/notify *"})

	next, result, err := Reconcile(lock, baseOptions(root, declared("notify", "0.2.0")))
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if result.Revoked {
		t.Fatal("revoked without --revoke-removed")
	}
	if len(result.Removed) != 1 {
		t.Fatalf("removed = %+v, want the withdrawal offered", result.Removed)
	}
	if got := openCodeRules(t, root); len(got) != 1 {
		t.Errorf("opencode.json rules = %+v, want the rule kept", got)
	}
	entry, _ := manifest.FindLockEntry(next, "notify")
	if len(entry.Granted) != 1 {
		t.Errorf("lock granted = %q, want it kept until withdrawn", entry.Granted)
	}
	if len(entry.Permissions) != 0 {
		t.Errorf("lock permissions = %q, want the new declaration", entry.Permissions)
	}
}

func TestReconcileRevokeRemoved(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith([]string{notifyRule}, []string{"shell:*/scripts/notify *"})
	seedRules(t, root, tool.Rule{Action: "shell", Resource: "*/scripts/notify *"})

	opts := baseOptions(root, declared("notify", "0.2.0"))
	opts.RevokeRemoved = true
	next, result, err := Reconcile(lock, opts)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if !result.Revoked || len(result.Removed) != 1 {
		t.Fatalf("result = %+v, want one withdrawn rule", result)
	}
	if got := openCodeRules(t, root); len(got) != 0 {
		t.Errorf("opencode.json rules = %+v, want the rule withdrawn", got)
	}
	entry, _ := manifest.FindLockEntry(next, "notify")
	if len(entry.Granted) != 0 {
		t.Errorf("lock granted = %q, want it dropped", entry.Granted)
	}
}

func TestReconcileSkipWritesNothing(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := lockWith(nil, nil)

	opts := baseOptions(root, declared("notify", "0.1.0", notifyRule))
	opts.Mode = ModeSkip
	next, result, err := Reconcile(lock, opts)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if !result.Skipped {
		t.Error("skipped = false, want true")
	}
	if got := openCodeRules(t, root); len(got) != 0 {
		t.Errorf("opencode.json rules = %+v, want none", got)
	}
	if entry, _ := manifest.FindLockEntry(next, "notify"); len(entry.Permissions) != 0 {
		t.Errorf("a skipped declaration must not be recorded: %+v", entry)
	}
}

func TestReconcileAskWithoutConfirmDeclines(t *testing.T) {
	t.Parallel()
	root := project(t)
	opts := baseOptions(root, declared("notify", "0.1.0", notifyRule))
	opts.Mode = ModeAsk

	_, result, err := Reconcile(manifest.Lock{}, opts)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if !result.Skipped {
		t.Error("a session with nobody to ask must decline")
	}
}

func TestReconcileClaudeOnlyRecordsWithoutWriting(t *testing.T) {
	t.Parallel()
	root := project(t)
	opts := baseOptions(root, declared("notify", "0.1.0", notifyRule))
	opts.Tools = []tool.ID{tool.ClaudeID}
	opts.Mode = ModeAsk // no prompt must happen: there is nothing to write

	next, result, err := Reconcile(manifest.Lock{}, opts)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if result.Skipped || len(result.Declared) != 1 {
		t.Errorf("result = %+v, want the declaration reported", result)
	}
	entry, ok := manifest.FindLockEntry(next, "notify")
	if !ok || !slices.Equal(entry.Permissions, []string{notifyRule}) {
		t.Errorf("lock entry = %+v, want the declaration recorded", entry)
	}
	if entry.Granted != nil {
		t.Errorf("granted = %q, want nothing: rsk wrote nothing", entry.Granted)
	}
}

func TestReconcileInvalidDeclarationFails(t *testing.T) {
	t.Parallel()
	_, _, err := Reconcile(manifest.Lock{},
		baseOptions(project(t), declared("bad", "1.0.0", "Bash(*)")))
	if err == nil {
		t.Fatal("Reconcile accepted Bash(*), want an error")
	}
}

func TestReconcileSharedRuleSurvivesOneRemoval(t *testing.T) {
	t.Parallel()
	root := project(t)
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "a", Version: "1.0.0", Permissions: []string{notifyRule}, Granted: []string{"shell:*/scripts/notify *"}},
		{Name: "b", Version: "1.0.0", Permissions: []string{notifyRule}, Granted: []string{"shell:*/scripts/notify *"}},
	}}
	seedRules(t, root, tool.Rule{Action: "shell", Resource: "*/scripts/notify *"})

	opts := baseOptions(root, declared("a", "2.0.0"))
	opts.RevokeRemoved = true
	_, result, err := Reconcile(lock, opts)
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if result.Revoked || len(result.Removed) != 0 {
		t.Errorf("result = %+v, want no withdrawal: b still declares the rule", result)
	}
	if got := openCodeRules(t, root); len(got) != 1 {
		t.Errorf("opencode.json rules = %+v, want the shared rule kept", got)
	}
}

func TestRevokeSkillsDropsOnlyUnsharedRules(t *testing.T) {
	t.Parallel()
	root := project(t)
	seedRules(t, root, tool.Rule{Action: "shell", Resource: "*/scripts/notify *"})
	lock := manifest.Lock{Skills: []manifest.LockEntry{
		{Name: "a", Version: "1.0.0", Granted: []string{"shell:*/scripts/notify *"}},
		{Name: "b", Version: "1.0.0", Granted: []string{"shell:*/scripts/notify *"}},
	}}

	next, revoked, err := RevokeSkills(lock, root, []tool.ID{tool.OpenCodeID}, []string{"a"})
	if err != nil {
		t.Fatalf("RevokeSkills: %v", err)
	}
	if len(revoked) != 0 {
		t.Errorf("revoked = %+v, want none: b still needs it", revoked)
	}
	if len(next.Skills) != 1 || next.Skills[0].Name != "b" {
		t.Errorf("lock = %+v, want only b left", next.Skills)
	}

	next, revoked, err = RevokeSkills(next, root, []tool.ID{tool.OpenCodeID}, []string{"b"})
	if err != nil {
		t.Fatalf("RevokeSkills (last owner): %v", err)
	}
	if len(revoked) != 1 {
		t.Errorf("revoked = %+v, want the rule withdrawn", revoked)
	}
	if len(next.Skills) != 0 {
		t.Errorf("lock = %+v, want empty", next.Skills)
	}
	if got := openCodeRules(t, root); len(got) != 0 {
		t.Errorf("opencode.json rules = %+v, want none", got)
	}
}

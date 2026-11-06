package tool

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func readInstructions(t *testing.T, path string) []any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read opencode.json: %v", err)
	}
	cfg := make(map[string]any)
	if err = json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse opencode.json: %v", err)
	}
	list, _ := cfg[OpenCodeInstructionsKey].([]any)
	return list
}

func TestOpenCodeSyncPinned_importsGeneratedIndexOnly(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, OpenCodeFileName)
	seed := `{"instructions":["docs/house-rules.md",".claude/skills/old-skill/SKILL.md"]}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	oc := &openCodeTool{}
	for range 2 { // the second run must not duplicate the entry
		if err := oc.SyncPinned(dir, []string{"tdd", "go-architect"}); err != nil {
			t.Fatalf("SyncPinned: %v", err)
		}
	}

	want := []any{"docs/house-rules.md", OpenCodeSkillsIndexEntry}
	if got := readInstructions(t, path); !reflect.DeepEqual(got, want) {
		t.Errorf("instructions = %v, want %v", got, want)
	}
}

func TestOpenCodeSyncPinned_nothingPinnedRemovesRskEntries(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, OpenCodeFileName)
	seed := `{"instructions":["docs/house-rules.md",".rsk/SKILLS.md"]}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := (&openCodeTool{}).SyncPinned(dir, nil); err != nil {
		t.Fatalf("SyncPinned: %v", err)
	}

	want := []any{"docs/house-rules.md"}
	if got := readInstructions(t, path); !reflect.DeepEqual(got, want) {
		t.Errorf("instructions = %v, want %v", got, want)
	}
}

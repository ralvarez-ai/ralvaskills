package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func TestUpdatePairNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []updatePair
		want []string
	}{
		{name: "empty", in: nil, want: []string{}},
		{
			name: "several",
			in: []updatePair{
				{name: "go-architect", installed: "1.0.0", latest: "1.1.0", newSkill: skill.Skill{}},
				{name: "tdd", installed: "1.0.0", latest: "1.0.1", newSkill: skill.Skill{}},
			},
			want: []string{"go-architect", "tdd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := updatePairNames(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updatePairNames(%v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

// writeLinkedSkill seeds the linked skill a pinned entry points at: a SKILL.md
// and the DIGEST.md it ships.
func writeLinkedSkill(t *testing.T, dir, version, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillMD := "---\nname: tdd\nversion: " + version + "\ndescription: red-green-refactor\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, digest.FileName), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestRefreshPinnedDigest covers the bug where `rsk update` re-linked a skill
// but left .rsk/SKILLS.md holding the previous digest until the next `rsk pin`:
// a changed DIGEST.md must land in the generated file, and the call must be a
// silent no-op outside a project.
//
//nolint:paralleltest // t.Chdir changes the process-wide working directory
func TestRefreshPinnedDigest(t *testing.T) {
	projDir := t.TempDir()
	t.Chdir(projDir)

	var buf bytes.Buffer
	refreshPinnedDigest(&buf)
	if got := buf.String(); got != "" {
		t.Errorf("outside a project: expected no output, got %q", got)
	}

	if err := manifest.WriteMod(filepath.Join(projDir, ".rsk"), manifest.Mod{Pinned: []string{"tdd"}}); err != nil {
		t.Fatalf("write rsk.mod: %v", err)
	}
	skillDir := filepath.Join(projDir, ".claude", "skills", "tdd")
	index := filepath.Join(projDir, ".rsk", digest.IndexFileName)

	writeLinkedSkill(t, skillDir, "1.0.0", "Rules\n- old rule\n")
	refreshPinnedDigest(&buf)
	if got := buf.String(); got != "" {
		t.Errorf("expected no warning, got %q", got)
	}
	assertIndexHas(t, index, "## tdd v1.0.0", "- old rule")

	// A pull rewrites DIGEST.md in place; the next refresh must carry the new
	// rule into SKILLS.md without a second `rsk pin`.
	writeLinkedSkill(t, skillDir, "1.0.1", "Rules\n- new rule\n")
	refreshPinnedDigest(&buf)
	assertIndexHas(t, index, "## tdd v1.0.1", "- new rule")

	data, err := os.ReadFile(index)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("- old rule")) {
		t.Errorf("SKILLS.md still has the previous digest:\n%s", data)
	}
}

// assertIndexHas fails when SKILLS.md lacks any of wants.
func assertIndexHas(t *testing.T, path string, wants ...string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	for _, want := range wants {
		if !bytes.Contains(data, []byte(want)) {
			t.Errorf("SKILLS.md lacks %q:\n%s", want, data)
		}
	}
}

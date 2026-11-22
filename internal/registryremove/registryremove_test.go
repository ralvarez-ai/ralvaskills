package registryremove

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const testIndex = `{
  "version": 2,
  "generated_at": "2026-01-01T00:00:00Z",
  "skills": {
    "demo": {
      "name": "demo",
      "latest": "1.2.0",
      "versions": {
        "1.0.0": {"version": "1.0.0"},
        "1.1.0": {"version": "1.1.0"},
        "1.2.0": {"version": "1.2.0"}
      }
    }
  }
}`

func demoEntry() *skillEntry {
	return &skillEntry{Versions: map[string]*versionEntry{"1.0.0": {}, "1.1.0": {}, "1.2.0": {}}}
}

func TestSelectVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		version      string
		minV, maxV   string
		want         []string
		wantErr      bool
		wantWholeSet bool
	}{
		{name: "exact", version: "1.1.0", want: []string{"1.1.0"}},
		{name: "missing exact", version: "9.9.9", wantErr: true},
		{name: "range", minV: "1.1.0", maxV: "1.2.0", want: []string{"1.1.0", "1.2.0"}},
		{name: "open lower bound", maxV: "1.0.0", want: []string{"1.0.0"}},
		{name: "whole skill", want: []string{"1.0.0", "1.1.0", "1.2.0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := selectVersions(demoEntry(), tt.version, tt.minV, tt.maxV)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			slices.Sort(got)
			if !tt.wantErr && !slices.Equal(got, tt.want) {
				t.Errorf("selectVersions = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunValidatesBeforeTouchingAnything(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(path, []byte(testIndex), 0o600); err != nil {
		t.Fatal(err)
	}

	bad := []Options{
		{IndexPath: path},
		{IndexPath: path, Skill: "demo", Version: "1.0.0", Min: "1.0.0"},
		{IndexPath: path, Skill: "demo", Version: "1.0"},
		{IndexPath: path, Skill: "ghost"},
	}
	for _, o := range bad {
		if err := Run(o); err == nil {
			t.Errorf("Run(%+v) = nil, want an error", o)
		}
	}
}

func TestRunDryRunLeavesTheIndexUntouched(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(path, []byte(testIndex), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Run(Options{IndexPath: path, Skill: "demo", Version: "1.1.0", DryRun: true}); err != nil {
		t.Fatalf("dry run: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != testIndex {
		t.Errorf("dry run changed the index:\n%s", got)
	}
}

func TestRunSkipGitHubRemovesTheVersionFromTheIndex(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(path, []byte(testIndex), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Run(Options{IndexPath: path, Skill: "demo", Version: "1.1.0", SkipGitHub: true}); err != nil {
		t.Fatalf("run: %v", err)
	}
	idx, err := loadIndex(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := idx.Skills["demo"].Versions["1.1.0"]; ok {
		t.Error("1.1.0 is still in the index")
	}
	if len(idx.Skills["demo"].Versions) != 2 {
		t.Errorf("versions = %d, want 2", len(idx.Skills["demo"].Versions))
	}
}

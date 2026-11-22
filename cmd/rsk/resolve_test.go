package main

import (
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func TestDedupSkills(t *testing.T) {
	t.Parallel()
	input := []skill.Skill{
		{Name: "a"},
		{Name: "b"},
		{Name: "a"},
		{Name: "c"},
		{Name: "b"},
	}
	got := dedupSkills(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 unique skills, got %d: %v", len(got), got)
	}
	names := map[string]bool{}
	for _, s := range got {
		if names[s.Name] {
			t.Errorf("duplicate skill %q in dedup output", s.Name)
		}
		names[s.Name] = true
	}
}

func TestJoinKeys(t *testing.T) {
	t.Parallel()
	// joinKeys returns a comma-separated string of map keys.
	// Order is non-deterministic so we test with a single-key map.
	m := map[string]string{"claude-code": "/path"}
	got := joinKeys(m)
	if got != "claude-code" {
		t.Errorf("joinKeys = %q, want %q", got, "claude-code")
	}
}

func TestParseNameVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       string
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{"go-architect", "go-architect", "", false},
		{"go-architect@1.0.0", "go-architect", "1.0.0", false},
		{"foo@", "foo", "", false},
		{"@1.0.0", "", "", true},
		{"@", "", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			name, pinVersion, err := parseNameVersion(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("parseNameVersion(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseNameVersion(%q): unexpected error: %v", tc.input, err)
			}
			if name != tc.wantName {
				t.Errorf("name = %q, want %q", name, tc.wantName)
			}
			if pinVersion != tc.wantVersion {
				t.Errorf("version = %q, want %q", pinVersion, tc.wantVersion)
			}
		})
	}
}

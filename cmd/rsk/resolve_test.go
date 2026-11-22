package main

import (
	"testing"
)

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

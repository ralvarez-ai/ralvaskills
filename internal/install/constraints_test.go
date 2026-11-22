package install_test

import (
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func TestConstraintsFromArgsKeepsOnlyVersioned(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
		want map[string]string
	}{
		{
			name: "no version constraints",
			args: []string{"go-grpc", "docs"},
			want: map[string]string{},
		},
		{
			name: "single skill with version",
			args: []string{"go-architect@1.0.0"},
			want: map[string]string{"go-architect": "1.0.0"},
		},
		{
			name: "mixed bundle and pinned skill",
			args: []string{"go-grpc", "go-architect@1.0.0"},
			want: map[string]string{"go-architect": "1.0.0"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := install.ConstraintsFromArgs(tc.args)
			if err != nil {
				t.Fatalf("install.ConstraintsFromArgs(%v): unexpected error: %v", tc.args, err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("install.ConstraintsFromArgs(%v) = %v, want %v", tc.args, got, tc.want)
			}
			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("install.ConstraintsFromArgs(%v)[%q] = %q, want %q", tc.args, k, got[k], v)
				}
			}
		})
	}

	if _, err := install.ConstraintsFromArgs([]string{"@1.0.0"}); err == nil {
		t.Error("install.ConstraintsFromArgs with invalid name@version: expected error, got nil")
	}
}

func TestDedupSkills(t *testing.T) {
	t.Parallel()
	input := []skill.Skill{
		{Name: "a"},
		{Name: "b"},
		{Name: "a"},
		{Name: "c"},
		{Name: "b"},
	}
	got := install.DedupSkills(input)
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

func TestCheckNoPersonalNeedsTheFlag(t *testing.T) {
	t.Parallel()
	skills := []skill.Skill{{Name: "a"}, {Name: "mine", IsPersonal: true}}
	if err := install.CheckNoPersonal(skills, false); err == nil {
		t.Error("a personal skill without --include personal: expected an error")
	}
	if err := install.CheckNoPersonal(skills, true); err != nil {
		t.Errorf("with the flag: unexpected error %v", err)
	}
	if err := install.CheckNoPersonal([]skill.Skill{{Name: "a"}}, false); err != nil {
		t.Errorf("no personal skill: unexpected error %v", err)
	}
}

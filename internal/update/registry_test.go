package update_test

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/update"
)

func TestNamesOfPairs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []update.Pair
		want []string
	}{
		{name: "empty", in: nil, want: []string{}},
		{
			name: "several",
			in: []update.Pair{
				{Name: "go-architect", Installed: "1.0.0", Latest: "1.1.0", Skill: skill.Skill{}},
				{Name: "tdd", Installed: "1.0.0", Latest: "1.0.1", Skill: skill.Skill{}},
			},
			want: []string{"go-architect", "tdd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := update.Names(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("update.Names(%v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNeedsUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		installed, locked, latest string
		want                      bool
	}{
		{name: "link behind", installed: "1.0.0", locked: "1.0.0", latest: "1.1.0", want: true},
		{name: "lock behind, link current", installed: "1.1.0", locked: "1.0.0", latest: "1.1.0", want: true},
		{name: "both current", installed: "1.1.0", locked: "1.1.0", latest: "1.1.0", want: false},
		{name: "untracked, link current", installed: "1.1.0", locked: "", latest: "1.1.0", want: false},
		{name: "untracked, link behind", installed: "1.0.0", locked: "", latest: "1.1.0", want: true},
		{name: "lock ahead of latest", installed: "1.1.0", locked: "1.2.0", latest: "1.1.0", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := update.NeedsUpdate(tt.installed, tt.locked, tt.latest); got != tt.want {
				t.Errorf("update.NeedsUpdate(%q, %q, %q) = %v, want %v",
					tt.installed, tt.locked, tt.latest, got, tt.want)
			}
		})
	}
}

// TestWithLockedNames covers the second half of the same bug: a lock entry
// whose link went missing was never checked, so it could not be repaired.
func TestWithLockedNames(t *testing.T) {
	t.Parallel()

	if got := update.WithLockedNames([]string{"tdd"}, nil); !reflect.DeepEqual(got, []string{"tdd"}) {
		t.Errorf("no lock: got %v, want the names unchanged", got)
	}
	got := update.WithLockedNames([]string{"tdd"}, map[string]string{"tdd": "1.0.0", "rsk-guide": "0.5.3"})
	if len(got) != 2 || got[0] != "tdd" || got[1] != "rsk-guide" {
		t.Errorf("got %v, want [tdd rsk-guide]", got)
	}
}

func TestLinkedNamesDeduplicatesAcrossTargets(t *testing.T) {
	t.Parallel()
	a, b := t.TempDir(), t.TempDir()
	for dir, names := range map[string][]string{a: {"tdd", "docs"}, b: {"docs", "unslop"}} {
		for _, n := range names {
			if err := os.Mkdir(filepath.Join(dir, n), 0o750); err != nil {
				t.Fatal(err)
			}
		}
	}
	got := update.LinkedNames([]string{a, b, filepath.Join(a, "missing")})
	slices.Sort(got)
	if want := []string{"docs", "tdd", "unslop"}; !slices.Equal(got, want) {
		t.Errorf("LinkedNames = %v, want %v", got, want)
	}
}

func TestLockedVersionsIsEmptyWithoutALock(t *testing.T) {
	t.Parallel()
	if got := update.LockedVersions(t.TempDir()); len(got) != 0 {
		t.Errorf("no lock: %v, want no versions", got)
	}
}

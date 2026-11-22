package skill

import (
	"path/filepath"
	"testing"
)

const commit40 = "0123456789abcdef0123456789abcdef01234567"

func TestDetectSource(t *testing.T) {
	t.Parallel()
	repo := filepath.FromSlash("/home/user/ralvaskills")
	official := filepath.FromSlash("/home/user/.config/rsk/cache/anthropic")
	registry := filepath.FromSlash("/home/user/.config/rsk/cache/registry")

	tests := []struct {
		name      string
		skillPath string
		want      Source
	}{
		{
			name:      "local repo skill",
			skillPath: filepath.Join(repo, "skills", "go-architect"),
			want:      SourceLocal,
		},
		{
			name:      "official cache skill",
			skillPath: filepath.Join(official, "skills", "claude-code"),
			want:      SourceOfficial,
		},
		{
			name:      "registry cache skill",
			skillPath: filepath.Join(registry, "go-architect", "1.0.0"),
			want:      SourceRegistry,
		},
		{
			name:      "unknown path defaults to local",
			skillPath: filepath.FromSlash("/some/random/path/skill"),
			want:      SourceLocal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DetectSource(tc.skillPath, repo, official, registry)
			if got != tc.want {
				t.Errorf("DetectSource(%q) = %v, want %v", tc.skillPath, got, tc.want)
			}
		})
	}
}

func TestParseGitLink(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dest := filepath.Join(root, "sources", "acme", "skills", "alpha", commit40)
	id, name, commit, ok := ParseGitLink(dest, root)
	if !ok || id != "@acme/skills" || name != "alpha" || commit != commit40 {
		t.Errorf("ParseGitLink = %q %q %q %v", id, name, commit, ok)
	}
	for _, bad := range []string{
		filepath.Join(root, "registry", "alpha", "1.0.0"),
		filepath.Join(root, "sources", "acme", "skills", "alpha"),
		filepath.Join(root, "sources", "acme", "skills", "alpha", commit40, "extra"),
		"/elsewhere/alpha",
	} {
		if _, _, _, ok = ParseGitLink(bad, root); ok {
			t.Errorf("ParseGitLink(%s) accepted a path that is not a git extraction", bad)
		}
	}
}

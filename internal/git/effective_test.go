package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/git"
)

func TestEffectiveURLShowsWhatGitsURLRewritesMakeOfIt(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "gitconfig")
	rules := "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/\n"
	if err := os.WriteFile(cfg, []byte(rules), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", cfg)

	https, err := git.ParseRemoteURL("https://github.com/acme/skills.git")
	if err != nil {
		t.Fatal(err)
	}
	got, err := git.EffectiveURL(t.Context(), https)
	if err != nil || got != "git@github.com:acme/skills.git" {
		t.Errorf("EffectiveURL = %q, %v: want the rewritten SSH URL", got, err)
	}

	ssh, _ := git.ParseRemoteURL("git@github.com:acme/skills.git")
	if got, err = git.EffectiveURL(t.Context(), ssh); err != nil || got != "git@github.com:acme/skills.git" {
		t.Errorf("an unaffected URL came back as %q, %v", got, err)
	}
}

func TestEffectiveURLWithoutRulesIsTheURLItself(t *testing.T) {
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_SYSTEM", "/dev/null")
	u, _ := git.ParseRemoteURL("https://github.com/acme/skills.git")
	if got, err := git.EffectiveURL(t.Context(), u); err != nil || got != u.String() {
		t.Errorf("EffectiveURL = %q, %v", got, err)
	}
}

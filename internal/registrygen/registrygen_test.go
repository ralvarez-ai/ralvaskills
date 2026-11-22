package registrygen

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, version string) Options {
	t.Helper()

	root := t.TempDir()
	skillDir := filepath.Join(root, "skills", "demo", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: demo-skill\nversion: " + version + "\ndescription: A demo skill.\n---\n\n# Demo\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(md), 0o600); err != nil {
		t.Fatal(err)
	}
	aliases := filepath.Join(root, "aliases.toml")
	if err := os.WriteFile(aliases, []byte("[aliases]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return Options{
		SkillsDir:  filepath.Join(root, "skills"),
		OutputDir:  filepath.Join(root, "dist"),
		GitHubRepo: GitHubRepoDefault,
		Aliases:    aliases,
	}
}

func TestRunWritesTheTarballAndTheIndex(t *testing.T) {
	t.Parallel()

	o := writeFixture(t, "1.0.0")
	if err := Run(o); err != nil {
		t.Fatalf("Run: %v", err)
	}
	for _, name := range []string{"demo-skill-v1.0.0.tar.gz", indexFileName, newVersionsFileName} {
		if _, err := os.Stat(filepath.Join(o.OutputDir, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
}

func TestRunSkipsAVersionThatIsAlreadyInTheIndex(t *testing.T) {
	t.Parallel()

	o := writeFixture(t, "1.0.0")
	if err := Run(o); err != nil {
		t.Fatalf("first run: %v", err)
	}

	o.ExistingIndex = filepath.Join(o.OutputDir, indexFileName)
	o.OutputDir = filepath.Join(filepath.Dir(o.OutputDir), "dist2")
	if err := Run(o); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(o.OutputDir, "demo-skill-v1.0.0.tar.gz")); !os.IsNotExist(err) {
		t.Errorf("a published version was packed again (stat err = %v)", err)
	}
}

func TestCheckWritesNothingAndRejectsABadVersion(t *testing.T) {
	t.Parallel()

	o := writeFixture(t, "1.0.0")
	o.Check = true
	if err := Run(o); err != nil {
		t.Fatalf("check on a valid version: %v", err)
	}
	if _, err := os.Stat(o.OutputDir); !os.IsNotExist(err) {
		t.Errorf("--check created the output dir (stat err = %v)", err)
	}

	bad := writeFixture(t, "1.0")
	bad.Check = true
	if err := Run(bad); err == nil {
		t.Error("check accepted a version that is not X.Y.Z")
	}
}

func TestRunFailsOnAMissingAliasTable(t *testing.T) {
	t.Parallel()

	o := writeFixture(t, "1.0.0")
	o.Aliases = filepath.Join(t.TempDir(), "absent.toml")
	if err := Run(o); err == nil {
		t.Error("Run succeeded without an alias table")
	}
}

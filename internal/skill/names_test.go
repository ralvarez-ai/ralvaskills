package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

func TestValidateName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		valid bool
	}{
		{"go-architect", true},
		{"tdd", true},
		{"hexagonal-architecture", true},
		{"go2", true},
		{"", false},
		{"Go-Architect", false},
		{"go--architect", false},
		{"-go-architect", false},
		{"go-architect-", false},
		{"go architect", false},
		{"go_architect", false},
		{strings.Repeat("a", MaxNameLen), true},
		{strings.Repeat("a", MaxNameLen+1), false},
	}
	for _, c := range cases {
		err := ValidateName(c.name)
		if c.valid && err != nil {
			t.Errorf("ValidateName(%q) = %v, want nil", c.name, err)
		}
		if !c.valid && err == nil {
			t.Errorf("ValidateName(%q) = nil, want an error", c.name)
		}
	}
}

func TestCheckRetirements(t *testing.T) {
	t.Parallel()

	retired := Retirements{
		"hexagonal-arch":      "hexagonal-architecture",
		"narration-architect": "narration-writer",
	}
	rejected := []string{
		"hexagonal-arch",
		"narration-architect",
		"hexagonal-arch-v2",
		"go-hexagonal-arch",
		"my-hexagonal-arch-notes",
	}
	for _, name := range rejected {
		if err := CheckRetirements(name, retired); err == nil {
			t.Errorf("CheckRetirements(%q) = nil, want an error", name)
		}
	}
	accepted := []string{
		"hexagonal-architecture",
		"narration-writer",
		"hexagonal",
		"architect",
		"narration",
	}
	for _, name := range accepted {
		if err := CheckRetirements(name, retired); err != nil {
			t.Errorf("CheckRetirements(%q) = %v, want nil", name, err)
		}
	}
}

func TestValidateRetirements(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		retired Retirements
		valid   bool
	}{
		{"empty", Retirements{}, true},
		{"single hop", Retirements{"old-name": "new-name"}, true},
		{"chain", Retirements{"a-b": "c-d", "c-d": "e-f"}, false},
		{"self", Retirements{"a-b": "a-b"}, false},
		{"empty target", Retirements{"a-b": ""}, false},
		{"bad key", Retirements{"A_B": "c-d"}, false},
	}
	for _, c := range cases {
		err := ValidateRetirements(c.retired)
		if c.valid && err != nil {
			t.Errorf("%s: ValidateRetirements() = %v, want nil", c.name, err)
		}
		if !c.valid && err == nil {
			t.Errorf("%s: ValidateRetirements() = nil, want an error", c.name)
		}
	}
}

// TestCatalogueNames is the guard: no skill in the repository may take a
// retired name, and every name must be legal, declared, and unique.
func TestCatalogueNames(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..")
	retired, err := LoadRetirements(filepath.Join(root, RetirementsPath))
	if err != nil {
		t.Fatalf("load %s: %v", RetirementsPath, err)
	}
	if validErr := ValidateRetirements(retired); validErr != nil {
		t.Fatalf("retirements table: %v", validErr)
	}

	issues, err := CatalogueIssues(filepath.Join(root, "skills"), retired)
	if err != nil {
		t.Fatalf("CatalogueIssues: %v", err)
	}
	for _, issue := range issues {
		t.Error(issue)
	}
}

func TestCatalogueIssuesCatchesRetiredName(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	dir := filepath.Join(root, "hexagonal-arch")
	if err := os.MkdirAll(dir, fsperm.Dir); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "---\nname: hexagonal-arch\nversion: 1.0.0\ndescription: Test skill.\n---\n\n# body\n"
	if err := os.WriteFile(filepath.Join(dir, SkillFileName), []byte(body), fsperm.File); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	issues, err := CatalogueIssues(root, Retirements{"hexagonal-arch": "hexagonal-architecture"})
	if err != nil {
		t.Fatalf("CatalogueIssues: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("CatalogueIssues() = no issues, want a retired-name error")
	}
	if !strings.Contains(issues[0], "retired") {
		t.Errorf("issue = %q, want it to mention a retired name", issues[0])
	}
}

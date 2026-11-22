package source

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitScanReportsAClashWhereAllRefusesTheSource(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.entries = append(repo.entries, tree("other/alpha"), blob("other/alpha/SKILL.md"))
	repo.blobs["other/alpha/SKILL.md"] = skillMD("alpha", "9.0.0")
	h := newHarness(t, nil, repo)

	if _, err := h.g.All(t.Context()); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("All err = %v, want ErrDuplicate", err)
	}
	res, err := h.g.Scan(t.Context())
	if err != nil {
		t.Fatalf("Scan err = %v: a scan exists to show the clash, not to fail on it", err)
	}
	if len(res.Clashes) != 1 || res.Clashes[0].Name != "alpha" ||
		strings.Join(res.Clashes[0].Paths, ",") != "other/alpha,skills/alpha" {
		t.Errorf("clashes = %+v", res.Clashes)
	}
	names := map[string]bool{}
	for _, s := range res.Skills {
		names[s.Name] = true
	}
	if !names["beta"] || names["alpha"] {
		t.Errorf("skills = %+v: a clashing name must not be listed as resolvable", res.Skills)
	}
}

func TestGitScanFlagsFrontmatterThatDisagreesWithTheDirectory(t *testing.T) {
	t.Parallel()
	repo := baseRepo()
	repo.blobs["skills/beta/SKILL.md"] = "---\nname: not-beta\ndescription: x\n---\n"
	h := newHarness(t, nil, repo)

	res, err := h.g.Scan(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	issues := map[string][]string{}
	for _, s := range res.Skills {
		issues[s.Name] = s.Issues
	}
	if len(issues["alpha"]) != 0 {
		t.Errorf("a clean skill has issues: %v", issues["alpha"])
	}
	joined := strings.Join(issues["beta"], "; ")
	if !strings.Contains(joined, "not-beta") || !strings.Contains(joined, "version") {
		t.Errorf("beta issues = %q, want the name mismatch and the missing version", joined)
	}
}

func TestPathScanReportsClashesAndSkippedWithoutFailing(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, dir := range []string{"a/same", "b/same", "fine", "Bad_Name"} {
		full := filepath.Join(root, filepath.FromSlash(dir))
		if err := os.MkdirAll(full, 0o750); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: " + filepath.Base(dir) + "\nversion: 1.0.0\ndescription: x\n---\n"
		if err := os.WriteFile(filepath.Join(full, "SKILL.md"), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	p, err := newPath(Spec{ID: "@me/scratch", Type: KindPath, URL: root, Declared: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, allErr := p.All(t.Context()); !errors.Is(allErr, ErrDuplicate) {
		t.Fatalf("All err = %v, want ErrDuplicate", allErr)
	}
	res, err := p.Scan(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Clashes) != 1 || res.Clashes[0].Name != "same" || len(res.Skills) != 1 || res.Skills[0].Name != "fine" ||
		len(res.Skipped) != 1 || res.Skipped[0].Path != "Bad_Name" {
		t.Errorf("scan = %+v", res)
	}
}

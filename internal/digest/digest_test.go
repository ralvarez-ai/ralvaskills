package digest

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSkill(t *testing.T, projectDir, name, version, digest string) {
	t.Helper()
	dir := filepath.Join(projectDir, ".claude", "skills", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillMD := "---\nname: " + name + "\nversion: " + version + "\ndescription: does " + name + "\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillMD), 0o644); err != nil {
		t.Fatal(err)
	}
	if digest != "" {
		if err := os.WriteFile(filepath.Join(dir, FileName), []byte(digest), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSync_readsNativeSkillDirectories(t *testing.T) {
	t.Parallel()
	for _, base := range []string{".opencode", ".agents"} {
		t.Run(base, func(t *testing.T) {
			t.Parallel()
			project := t.TempDir()
			dir := filepath.Join(project, base, "skills", "example")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(project, ".rsk"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(
				filepath.Join(dir, "SKILL.md"),
				[]byte("---\nname: example\nversion: 1.0.0\ndescription: Example\n---\n"),
				0o644,
			); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, FileName), []byte("- native rule"), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := Sync(project, []string{"example"}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(project, ".rsk", IndexFileName))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "- native rule") {
				t.Fatalf("missing digest: %s", data)
			}
		})
	}
}

func TestSync_writesDigestsAndFallbacks(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rsk"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSkill(
		t,
		dir,
		"tdd",
		"1.1.0",
		"# tdd digest\n\nPurpose.\n\n## Rules\n- red first\n\nFull guide: invoke the tdd skill.\n",
	)
	writeSkill(t, dir, "no-digest", "0.2.0", "")

	tokens, err := Sync(dir, []string{"tdd", "no-digest", "ghost"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if tokens <= 0 {
		t.Errorf("expected a positive token estimate, got %d", tokens)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".rsk", IndexFileName))
	if err != nil {
		t.Fatalf("read SKILLS.md: %v", err)
	}
	got := string(data)
	for _, want := range []string{
		"## tdd v1.1.0",
		"### Rules",
		"- red first",
		"Full guide: invoke the tdd skill.",
		"## no-digest v0.2.0",
		"- no-digest: does no-digest",
		"## ghost\n",
		"is it installed?",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("SKILLS.md lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "# tdd digest") {
		t.Errorf("the digest's own H1 should be dropped:\n%s", got)
	}
}

func TestSync_idempotent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rsk"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSkill(t, dir, "tdd", "1.1.0", "Rules\n")
	path := filepath.Join(dir, ".rsk", IndexFileName)

	if _, err := Sync(dir, []string{"tdd"}); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(path)
	if _, err := Sync(dir, []string{"tdd"}); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(path)
	if !bytes.Equal(first, second) {
		t.Errorf("second run changed the file:\n%s\n---\n%s", first, second)
	}
}

func TestSync_removesFileWhenNothingPinned(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rsk"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSkill(t, dir, "tdd", "1.1.0", "Rules\n")
	path := filepath.Join(dir, ".rsk", IndexFileName)
	if _, err := Sync(dir, []string{"tdd"}); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(dir, nil); err != nil {
		t.Fatalf("Sync empty: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("SKILLS.md should be gone, stat err = %v", err)
	}
	if _, err := Sync(dir, nil); err != nil {
		t.Errorf("removing an absent file should be a no-op: %v", err)
	}
}

func TestMeasure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rsk"), 0o755); err != nil {
		t.Fatal(err)
	}

	names := []string{"biggest", "big", "medium", "smaller", "tiny", "tiny2"}
	for i, name := range names {
		writeSkill(t, dir, name, "1.0.0", strings.Repeat("always-on rule\n", 40*(len(names)-i)))
	}

	size := Measure(dir, names)
	if !size.OverBudget() {
		t.Errorf("six digests should pass the %d-token budget, measured %d", BudgetTokens, size.Tokens)
	}
	if len(size.Largest) != len(names) {
		t.Fatalf("Largest has %d entries, want %d", len(size.Largest), len(names))
	}
	if size.Largest[0].Name != "biggest" {
		t.Errorf("Largest[0] = %q, want the biggest digest", size.Largest[0].Name)
	}
	for i := 1; i < len(size.Largest); i++ {
		if size.Largest[i].Tokens > size.Largest[i-1].Tokens {
			t.Errorf("Largest is not ordered by size: %v", size.Largest)
		}
	}

	// Measure predicts what Sync writes, and writes nothing itself.
	if _, err := os.Stat(filepath.Join(dir, ".rsk", IndexFileName)); !os.IsNotExist(err) {
		t.Errorf("Measure must not create %s (stat err = %v)", IndexFileName, err)
	}
	tokens, err := Sync(dir, names)
	if err != nil {
		t.Fatal(err)
	}
	if tokens != size.Tokens {
		t.Errorf("Sync measured %d tokens, Measure predicted %d", tokens, size.Tokens)
	}
	if empty := Measure(dir, nil); empty.Tokens != 0 || empty.OverBudget() || len(empty.Largest) != 0 {
		t.Errorf("nothing pinned should measure zero, got %+v", empty)
	}
}

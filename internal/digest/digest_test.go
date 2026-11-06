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

package skill

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// writeSkillFile writes a SKILL.md in a fresh directory and returns that
// directory.
func writeSkillFile(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, SkillFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}
	return dir
}

func TestReadFrontmatterAllowedToolsList(t *testing.T) {
	t.Parallel()
	dir := writeSkillFile(t, `---
name: notify-on-finish
version: 0.1.0
allowed-tools:
  - Bash(${CLAUDE_SKILL_DIR}/scripts/notify-on-finish *)
  - Bash(notify-on-finish *)
description: sends notifications
---

Body.
`)

	fm, err := ReadFrontmatter(dir)
	if err != nil {
		t.Fatalf("ReadFrontmatter: %v", err)
	}
	want := []string{
		"Bash(${CLAUDE_SKILL_DIR}/scripts/notify-on-finish *)",
		"Bash(notify-on-finish *)",
	}
	if !slices.Equal(fm.AllowedTools, want) {
		t.Errorf("allowed-tools = %q, want %q", fm.AllowedTools, want)
	}
	if fm.Description != "sends notifications" {
		t.Errorf("description = %q; the block list must end at the next key", fm.Description)
	}
}

func TestReadFrontmatterAllowedToolsInline(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"allowed-tools: Read Grep\n":                {"Read", "Grep"},
		"allowed-tools: Read, Bash(git status *)\n": {"Read", "Bash(git status *)"},
		`allowed-tools: "Bash(ls *)"` + "\n":        {"Bash(ls *)"},
		"allowed-tools:\n":                          nil,
		"allowed-tools: Bash(a b) Bash(c d)\n":      {"Bash(a b)", "Bash(c d)"},
	}
	for value, want := range cases {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			dir := writeSkillFile(t, "---\nname: x\nversion: 1.0.0\n"+value+"---\n")
			fm, err := ReadFrontmatter(dir)
			if err != nil {
				t.Fatalf("ReadFrontmatter: %v", err)
			}
			if !slices.Equal(fm.AllowedTools, want) {
				t.Errorf("allowed-tools = %q, want %q", fm.AllowedTools, want)
			}
		})
	}
}

func TestValidateAllowedTools(t *testing.T) {
	t.Parallel()
	valid := [][]string{
		nil,
		{"Bash(git status *)"},
		{"Read(~/docs/**)", "Write(~/out/**)"},
		{"Bash(${CLAUDE_SKILL_DIR}/scripts/x *)"},
	}
	for _, entries := range valid {
		if err := ValidateAllowedTools(entries); err != nil {
			t.Errorf("ValidateAllowedTools(%q) = %v, want nil", entries, err)
		}
	}

	invalid := map[string][]string{
		"too many":       {"Bash(a)", "Bash(b)", "Bash(c)", "Bash(d)", "Bash(e)", "Bash(f)"},
		"no parentheses": {"Bash"},
		"bare tool":      {"Read"},
		"empty spec":     {"Bash()"},
		"whole wildcard": {"Bash(*)"},
		"deep wildcard":  {"Bash(**)"},
		"leading glob":   {"Bash(*notify-on-finish*)"},
		"too long":       {"Bash(" + strings.Repeat("a", MaxAllowedToolLength) + ")"},
	}
	for name, entries := range invalid {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := ValidateAllowedTools(entries); err == nil {
				t.Errorf("ValidateAllowedTools(%q) = nil, want an error", entries)
			}
		})
	}
}

func TestNormalizeAllowedTools(t *testing.T) {
	t.Parallel()
	got := NormalizeAllowedTools([]string{"  Bash(b) ", "Bash(a)", "Bash(b)", ""})
	want := []string{"Bash(a)", "Bash(b)"}
	if !slices.Equal(got, want) {
		t.Errorf("NormalizeAllowedTools = %q, want %q", got, want)
	}
}

func TestSplitRule(t *testing.T) {
	t.Parallel()
	tool, spec, ok := SplitRule("Bash(git status *)")
	if !ok || tool != "Bash" || spec != "git status *" {
		t.Errorf("SplitRule = (%q, %q, %v), want (Bash, git status *, true)", tool, spec, ok)
	}
	for _, entry := range []string{"Read", "Bash", "()", "Bash(", ""} {
		if _, _, parsed := SplitRule(entry); parsed {
			t.Errorf("SplitRule(%q) ok = true, want false", entry)
		}
	}
}

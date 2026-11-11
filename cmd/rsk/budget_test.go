package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/digest"
)

// bigDigest stands in for one skill's DIGEST.md: enough always-on rules that a
// handful of them pass the 2000-token budget (the estimator counts about four
// bytes per token).
var bigDigest = strings.Repeat("Keep the early return; name the magic value; keep the branch flat.\n", 40)

// seedPinned writes name's linked skill with body as its DIGEST.md.
func seedPinned(t *testing.T, projectDir, name, body string) {
	t.Helper()
	writeLinkedSkill(t, filepath.Join(projectDir, ".claude", "skills", name), name, "1.0.0", body)
}

// TestWarnSkillsIndexBudget checks the over-budget warning: silent with nothing
// pinned or with a set that fits, one stderr line naming the largest digests
// once it passes.
func TestWarnSkillsIndexBudget(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	var buf bytes.Buffer

	warnSkillsIndexBudget(&buf, dir, nil)
	if buf.Len() != 0 {
		t.Errorf("nothing pinned should warn nothing, got %q", buf.String())
	}

	seedPinned(t, dir, "small", "one rule\n")
	buf.Reset()
	warnSkillsIndexBudget(&buf, dir, []string{"small"})
	if buf.Len() != 0 {
		t.Errorf("a small set should fit the budget, got %q", buf.String())
	}

	names := []string{"biggest", "big", "medium", "smaller", "tiny", "tiny2"}
	for i, name := range names {
		seedPinned(t, dir, name, strings.Repeat(bigDigest, len(names)-i))
	}
	buf.Reset()
	warnSkillsIndexBudget(&buf, dir, names)

	got := buf.String()
	if !strings.Contains(got, "tokens of the") || !strings.Contains(got, "budget") {
		t.Fatalf("expected a budget warning, got %q", got)
	}
	if !strings.Contains(got, "largest: biggest ") {
		t.Errorf("the warning should name the biggest digest first, got %q", got)
	}
	if n := strings.Count(got, "\n"); n != 1 {
		t.Errorf("expected one line, got %d:\n%s", n, got)
	}
}

// TestSkillsIndexLine pins the exact wording, and that the largest-digest list
// is capped at three names.
func TestSkillsIndexLine(t *testing.T) {
	t.Parallel()

	fit := digest.Size{Tokens: 840, Largest: []digest.EntrySize{{Name: "a", Tokens: 500}}}
	if got, want := skillsIndexLine(fit), ".rsk/SKILLS.md ~840 tokens of the 2000-token budget"; got != want {
		t.Errorf("skillsIndexLine(fit) = %q, want %q", got, want)
	}

	over := digest.Size{
		Tokens: 2400,
		Largest: []digest.EntrySize{
			{Name: "docs-architect", Tokens: 780},
			{Name: "go-architect", Tokens: 610},
			{Name: "tdd", Tokens: 420},
			{Name: "caveman", Tokens: 90},
		},
	}
	got := skillsIndexLine(over)
	want := ".rsk/SKILLS.md ~2400 tokens of the 2000-token budget — largest: docs-architect 780, go-architect 610, tdd 420"
	if got != want {
		t.Errorf("skillsIndexLine(over) = %q, want %q", got, want)
	}
}

// TestPrintSkillsIndexStatus checks the `rsk status` line: a plain muted line
// while the file fits, a warning once it does not.
func TestPrintSkillsIndexStatus(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	seedPinned(t, dir, "small", "one rule\n")

	var buf bytes.Buffer
	printSkillsIndexStatus(&buf, dir, []string{"small"})
	if !strings.Contains(buf.String(), "tokens of the 2000-token budget") {
		t.Fatalf("expected a size line, got %q", buf.String())
	}
	if strings.Contains(buf.String(), "largest: ") {
		t.Errorf("a set inside the budget should not list digests, got %q", buf.String())
	}

	for _, name := range []string{"a", "b", "c", "d", "e", "f"} {
		seedPinned(t, dir, name, bigDigest)
	}
	buf.Reset()
	printSkillsIndexStatus(&buf, dir, []string{"small", "a", "b", "c", "d", "e", "f"})
	if !strings.Contains(buf.String(), "largest: ") {
		t.Errorf("an over-budget set should name the largest digests, got %q", buf.String())
	}
}

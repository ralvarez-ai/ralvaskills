package main

import (
	"fmt"
	"io"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// largestShown is how many skills the over-budget line names.
const largestShown = 3

// skillsIndexLine renders the generated file's size, with the biggest digests
// appended once the estimate passes the budget.
func skillsIndexLine(size digest.Size) string {
	line := fmt.Sprintf(".rsk/%s ~%d tokens of the %d-token budget",
		digest.IndexFileName, size.Tokens, digest.BudgetTokens)
	if !size.OverBudget() {
		return line
	}

	names := make([]string, 0, largestShown)
	for _, e := range size.Largest {
		if len(names) == largestShown {
			break
		}
		names = append(names, fmt.Sprintf("%s %d", e.Name, e.Tokens))
	}
	return line + " — largest: " + strings.Join(names, ", ")
}

// warnSkillsIndexBudget writes one warning to errOut when the project's
// .rsk/SKILLS.md estimate passes digest.BudgetTokens. It goes to stderr so a
// JSON run of the same command stays clean, and it is a no-op with nothing
// pinned.
func warnSkillsIndexBudget(errOut io.Writer, projectDir string, pinned []string) {
	if len(pinned) == 0 {
		return
	}
	size := digest.Measure(projectDir, pinned)
	if !size.OverBudget() {
		return
	}
	ui.Warn(errOut, skillsIndexLine(size))
}

// printSkillsIndexStatus renders the same estimate under `rsk status`: muted
// while inside the budget, a warning once it passes.
func printSkillsIndexStatus(out io.Writer, projectDir string, pinned []string) {
	if len(pinned) == 0 {
		return
	}
	size := digest.Measure(projectDir, pinned)
	if size.OverBudget() {
		ui.Warn(out, skillsIndexLine(size))
		return
	}
	ui.Indent(out, skillsIndexLine(size))
}

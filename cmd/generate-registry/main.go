// generate-registry walks the skills/ directory, creates versioned tarballs for
// any skill whose version is not yet in the existing index, and writes an
// updated index.json plus a new-versions.json for the CI publish step.
package main

import (
	"flag"
	"fmt"
	"os"

	"go.ralva.dev/ralvaskills/v3/internal/registrygen"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func main() {
	var o registrygen.Options
	flag.StringVar(&o.SkillsDir, "skills-dir", skill.SkillsFolderName, "Path to skills directory")
	flag.StringVar(&o.OutputDir, "output-dir", "dist", "Output directory for tarballs and index")
	flag.StringVar(&o.ExistingIndex, "existing-index", "", "Path to existing index.json to merge with")
	flag.StringVar(
		&o.GitHubRepo, "github-repo", registrygen.GitHubRepoDefault,
		"GitHub repo (owner/name) used to build release asset URLs",
	)
	flag.BoolVar(&o.Check, "check", false, "Preview what would be published and validate versions; no files written")
	flag.StringVar(
		&o.Aliases, "aliases", skill.RetirementsPath,
		"Path to the retired-name table, relative to the repo root",
	)
	flag.Parse()

	if err := registrygen.Run(o); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

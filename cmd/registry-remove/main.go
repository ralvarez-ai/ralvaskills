// registry-remove edits registry/index.json to drop a skill, a single
// version, or an inclusive version range, then deletes the matching
// "<skill>@vX.Y.Z" GitHub tags and releases so the registry and GitHub stay
// in sync. Without --yes it only prints the plan.
package main

import (
	"flag"
	"fmt"
	"os"

	"go.ralva.dev/ralvaskills/v3/internal/registryremove"
)

func main() {
	var o registryremove.Options
	flag.StringVar(&o.IndexPath, "index", "registry/index.json", "Path to the registry index to edit")
	flag.StringVar(&o.Skill, "skill", "", "Name of the skill to remove from (required)")
	flag.StringVar(&o.Version, "version", "", "Remove only this exact version (mutually exclusive with --min/--max)")
	flag.StringVar(&o.Min, "min", "", "Remove versions >= this version (inclusive; use with --max for a range)")
	flag.StringVar(&o.Max, "max", "", "Remove versions <= this version (inclusive; use with --min for a range)")
	flag.StringVar(
		&o.Repo,
		"repo",
		registryremove.GitHubRepoDefault,
		"GitHub repo (owner/name) whose tags/releases get deleted",
	)
	flag.BoolVar(&o.SkipGitHub, "skip-github", false, "Only edit the index; leave GitHub tags/releases untouched")
	flag.BoolVar(
		&o.DryRun,
		"dry-run",
		false,
		"Print the removal plan and exit; index.json and GitHub are left untouched",
	)
	flag.Parse()

	if err := registryremove.Run(o); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
